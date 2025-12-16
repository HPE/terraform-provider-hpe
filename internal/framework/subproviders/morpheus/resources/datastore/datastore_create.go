// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package datastore

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"slices"
	"time"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"
	"github.com/cenkalti/backoff/v5"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/HPE/terraform-provider-hpe/internal/framework/utils"
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

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &cfg)...)
	if resp.Diagnostics.HasError() {
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
			"create datastore resource",
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

	// Set the resource ID locally but NOT in state yet
	plan.Id = types.Int64Value(id)

	// Helper function to delete the datastore if anything goes wrong
	deleteOnError := func() {
		utils.CleanupResourceOnError(ctx, utils.CleanupConfig{
			ResourceType: "datastore",
			ResourceID:   id,
			DeleteFunc: func(ctx context.Context, id int64) (*http.Response, error) {
				_, resp, err := client.DatastoresAPI.DeleteDatastores(ctx, id).Execute()
				return resp, err
			},
			GetFunc: func(ctx context.Context, id int64) (*http.Response, error) {
				_, resp, err := client.DatastoresAPI.GetDatastores(ctx, id).Execute()
				return resp, err
			},
			Diagnostics: &resp.Diagnostics,
		})
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
		deleteOnError()
		return
	}

	state, gdiags := getDatastoreAsState(ctx, plan.Id.ValueInt64(), plan, client)
	if gdiags.HasError() {
		resp.Diagnostics.Append(gdiags...)
		resp.Diagnostics.AddError(
			"create datastore resource",
			fmt.Sprintf("datastore %d: failed to read from api", plan.Id.ValueInt64()),
		)
		deleteOnError()
		return
	}

	// For now, we need to call update to set resourcePermission and tenantPermissions
	// because the API does not set these on create, even if provided.
	state, gdiags = updateDatastore(ctx, plan.Id.ValueInt64(), plan, state, client)
	if gdiags.HasError() {
		resp.Diagnostics.Append(gdiags...)
		resp.Diagnostics.AddError(
			"create datastore resource",
			fmt.Sprintf("datastore %d: failed to update", plan.Id.ValueInt64()),
		)
		deleteOnError()
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		deleteOnError()
		return
	}
}
