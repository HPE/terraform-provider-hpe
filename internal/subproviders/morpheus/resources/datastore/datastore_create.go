// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

//go:build experimental

package datastore

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"slices"
	"time"

	"github.com/cenkalti/backoff/v5"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/sdk"
)

const (
	cloudRefType   = "ComputeZone"
	clusterRefType = "ComputeServerGroup"

	associatedResourceTypeCluster = "Cluster"
	associatedResourceTypeCloud   = "Cloud"
)

var (
	CreateTargetStatuses = []string{
		"provisioned",
	}

	CreateErrorStatuses = []string{
		"failed",
		"warning",
	}
)

func checkStatusDone(status string, targetStatuses []string, errorStatuses []string) error {
	switch {
	case slices.Contains(errorStatuses, status):
		return backoff.Permanent(errors.New("reached error status: " + status))
	case slices.Contains(targetStatuses, status):
		return nil
	default:
		return backoff.RetryAfter(5)
	}
}

func (r *Resource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	var plan, cfg DatastoreModel

	tflog.Info(ctx, "creating datastore")

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		tflog.Info(ctx, "error getting plan")
		return
	}

	// Read configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &cfg)...)
	if resp.Diagnostics.HasError() {
		tflog.Info(ctx, "error getting config")
		return
	}

	name := plan.Name.ValueString()
	datastoreType := plan.DatastoreType
	associatedResourceType := plan.AssociatedResourceType.ValueString()
	associatedResourceId := plan.AssociatedResourceId.ValueInt64()

	var config DatastoreModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create a Morpheus API client
	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"create cloud resource",
			"cloud "+name+": failed to create client: "+err.Error(),
		)

		return
	}

	var id int64
	switch associatedResourceType {
	case associatedResourceTypeCloud:
		id = datastoreCreateDatastore(
			ctx,
			datastoreType,
			associatedResourceType, name,
			associatedResourceId,
			client,
			plan,
			resp,
		)

		if resp.Diagnostics.HasError() {
			return
		}
	case associatedResourceTypeCluster:
		id = datastoreCreateCluster(
			ctx,
			datastoreType,
			name,
			associatedResourceId,
			client,
			plan,
			resp,
		)

		if resp.Diagnostics.HasError() {
			return
		}
	default:
		resp.Diagnostics.AddError(
			"create datastore resource",
			"datastore "+name+": invalid associated_resource_type "+associatedResourceType+", must be 'Cloud' or 'Cluster'",
		)
		return
	}

	// Set the resource ID
	plan.Id = types.Int64Value(id)

	// write id as soon as possible
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Wait for the datastore to be ready
	waitForReady := func() (*sdk.GetDatastores200Response, error) {
		response, hresp, err := client.DatastoresAPI.GetDatastores(ctx, plan.Id.ValueInt64()).Execute()
		if err != nil || hresp.StatusCode != http.StatusOK {
			return nil, backoff.Permanent(err)
		}

		status := response.GetDatastore().Status

		return response, checkStatusDone(
			status,
			CreateTargetStatuses,
			CreateErrorStatuses,
		)
	}

	if r, err := backoff.Retry(
		ctx,
		waitForReady,
		backoff.WithBackOff(backoff.NewConstantBackOff(5*time.Second)),
		backoff.WithMaxElapsedTime(5*time.Minute),
	); err != nil {
		var status string

		if r != nil {
			status = r.GetDatastore().Status
		}

		resp.Diagnostics.AddError(
			"create datastore resource",
			fmt.Sprintf(
				"datastore %d: provisioning failed current status is: %v",
				plan.Id.ValueInt64(),
				status,
			),
		)
	}

	state, pdiags := getDatastoreAsState(ctx, plan.Id.ValueInt64(), plan, client)
	if pdiags.HasError() {
		resp.Diagnostics.Append(pdiags...)
		resp.Diagnostics.AddError(
			"create datastore resource",
			fmt.Sprintf("datastore %d: failed to read from api", plan.Id.ValueInt64()),
		)

		return
	}

	// TODO reenable this when we add datastore_update.go
	/*
		// For now, we need to call update to set resourcePermission and tenantPermissions
		// because the API does not set these on create, even if provided.
		state, pdiags = updateDatastoreFunc(ctx, plan.Id.ValueInt64(), plan, state, client)
		if pdiags.HasError() {
			resp.Diagnostics.Append(pdiags...)
			resp.Diagnostics.AddError(
				"create datastore resource",
				fmt.Sprintf("datastore %d: failed to update", plan.Id.ValueInt64()),
			)

			return
		}

	*/

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

}
