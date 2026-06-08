// (C) Copyright 2025-2026 Hewlett Packard Enterprise Development LP

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

	"github.com/HPE/terraform-provider-hpe/utils/cleanup"
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

	// Set the resource ID
	plan.Id = types.Int64Value(id)

	// Helper to taint the resource state on an error after the POST request
	taintResourceState := func(id int64) {
		cleanup.TaintResourceState(ctx, cleanup.TaintResourceStateConfig{
			ResourceType: "datastore",
			ResourceID:   id,
			StateWriter:  &resp.State,
			Diagnostics:  &resp.Diagnostics,
		})
	}

	// Wait for the datastore to be ready
	waitForReady := func() (*sdk.GetDatastores200Response, error) {
		response, hresp, err := client.DatastoresAPI.GetDatastores(ctx, id).Execute()
		if err != nil {
			if hresp == nil || hresp.StatusCode != http.StatusOK {
				return nil, backoff.Permanent(err)
			}
		}

		if response == nil || response.Datastore == nil {
			return nil, backoff.Permanent(fmt.Errorf("missing datastore in response"))
		}

		status := response.Datastore.Status

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
			status = r.Datastore.Status
		}

		// Unwrap the error to get the API/SDK error message if present
		var errUnwrapped error
		if r == nil {
			errUnwrapped = errors.Unwrap(err)
		}

		if errUnwrapped != nil {
			resp.Diagnostics.AddError(
				"datastore provisioning failed",
				fmt.Sprintf("Datastore %d failed to reach provisioned status: %v", id, errUnwrapped),
			)
		} else {
			resp.Diagnostics.AddError(
				"datastore provisioning failed",
				fmt.Sprintf("Datastore %d failed to reach provisioned status. Current status: %s", id, status),
			)
		}
		taintResourceState(id)

		return
	}

	state, gdiags := getDatastoreAsState(ctx, id, plan, client)
	if gdiags.HasError() {
		resp.Diagnostics.Append(gdiags...)
		resp.Diagnostics.AddError(
			"failed to read datastore state",
			fmt.Sprintf("Datastore %d was created but could not be read", id),
		)
		taintResourceState(id)

		return
	}

	// For now, we need to call update to set resourcePermission and tenantPermissions
	// because the API does not set these on create, even if provided.
	state, gdiags = updateDatastore(ctx, id, plan, state, client)
	if gdiags.HasError() {
		resp.Diagnostics.Append(gdiags...)
		resp.Diagnostics.AddError(
			"failed to update datastore",
			fmt.Sprintf("Datastore %d was created but permissions could not be updated", id),
		)
		taintResourceState(id)

		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError(
			"failed to set datastore state",
			fmt.Sprintf("Datastore %d was created but state could not be saved", id),
		)
		taintResourceState(id)

		return
	}
}
