// (C) Copyright 2025-2026 Hewlett Packard Enterprise Development LP

package datastore

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"
	"github.com/cenkalti/backoff/v5"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
)

var (
	resourcePermissionsUpdateFunc = sdk.NewUpdateDatastoresRequestDatastoreResourcePermissionsWithDefaults
	sitesPermissionsUpdateFunc    = sdk.NewUpdateCloudFoldersRequestFolderResourcePermissionsSitesInnerWithDefaults
	plansPermissionsUpdateFunc    = sdk.NewUpdateCloudFoldersRequestFolderResourcePermissionsPlansInnerWithDefaults
)

func updateDatastore(
	ctx context.Context,
	id int64,
	plan DatastoreModel,
	state DatastoreModel,
	client *sdk.APIClient,
) (DatastoreModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	name := plan.Name.ValueString()

	updateDatastore := sdk.NewUpdateDatastoresRequestDatastoreWithDefaults()
	updateDatastore.AdditionalProperties = make(map[string]any)

	// If the plan has unknown value for Visibility, use the state values
	// We can't use a PlanModifier for this since we call updateDatastore during Create
	// and Update, and in Create the state is not available since it hasn't been
	// written-out yet
	switch plan.Visibility.IsUnknown() {
	case true:
		updateDatastore.Visibility = state.Visibility.ValueStringPointer()
	case false:
		updateDatastore.Visibility = plan.Visibility.ValueStringPointer()
	}

	// If the plan has unknown value for Active, use the state values
	// We can't use a PlanModifier for this since we call updateDatastore during Create
	// and Update, and in Create the state is not available since it hasn't been
	// written-out yet
	switch plan.Active.IsUnknown() {
	case true:
		updateDatastore.Active = state.Active.ValueBoolPointer()
	case false:
		updateDatastore.Active = plan.Active.ValueBoolPointer()
	}

	if !plan.ResourcePermissions.IsNull() {
		resourcePermissions := resourcePermissionsUpdateFunc()
		resourcePermissions.All = plan.ResourcePermissions.All.ValueBoolPointer()

		if !plan.ResourcePermissions.Groups.IsNull() && !plan.ResourcePermissions.Groups.IsUnknown() {
			var groupsValues []GroupsValue
			d := plan.ResourcePermissions.Groups.ElementsAs(ctx, &groupsValues, false)
			diags.Append(d...)
			if diags.HasError() {
				return DatastoreModel{}, diags
			}

			var sites []sdk.UpdateCloudFoldersRequestFolderResourcePermissionsSitesInner
			for _, groupsValue := range groupsValues {
				site := sitesPermissionsUpdateFunc()
				site.Id = groupsValue.Id.ValueInt64Pointer()
				sites = append(sites, *site)
			}

			resourcePermissions.Sites = sites
		}

		resourcePermissions.AllPlans = plan.ResourcePermissions.AllPlans.ValueBoolPointer()
		// nolint:duplicate
		if !plan.ResourcePermissions.Plans.IsNull() && !plan.ResourcePermissions.Plans.IsUnknown() {
			var plansValues []PlansValue
			d := plan.ResourcePermissions.Plans.ElementsAs(ctx, &plansValues, false)
			diags.Append(d...)
			if diags.HasError() {
				return DatastoreModel{}, diags
			}

			var plans []sdk.UpdateCloudFoldersRequestFolderResourcePermissionsPlansInner
			for _, plansValue := range plansValues {
				planItem := plansPermissionsUpdateFunc()
				planItem.Id = plansValue.Id.ValueInt64Pointer()
				plans = append(plans, *planItem)
			}

			resourcePermissions.Plans = plans
		}

		updateDatastore.ResourcePermissions = resourcePermissions
	}

	if !plan.Tenants.IsNull() && !plan.Tenants.IsUnknown() {
		var tenantsValues []TenantsValue
		d := plan.Tenants.ElementsAs(ctx, &tenantsValues, false)
		diags.Append(d...)
		if diags.HasError() {
			return DatastoreModel{}, diags
		}

		// tenantPermissions := sdk.NewSaveDatastoreRequestDatastoreTenantPermissionsWithDefaults()
		tenantPermissions := sdk.NewUpdateDatastoresRequestDatastoreTenantPermissionsWithDefaults()
		var accounts []int64
		for _, tenantsValue := range tenantsValues {
			accounts = append(accounts, tenantsValue.Id.ValueInt64())
		}
		tenantPermissions.Accounts = accounts
		updateDatastore.TenantPermissions = tenantPermissions
	}

	updateDatastoreReq := sdk.NewUpdateDatastoresRequestWithDefaults()
	updateDatastoreReq.Datastore = *updateDatastore

	response, hresp, err := client.DatastoresAPI.UpdateDatastores(ctx, id).
		UpdateDatastoresRequest(*updateDatastoreReq).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		diags.AddError(
			"update datastore resource",
			"datastore "+name+" PUT failed: "+errfmt.ErrMsg(err, hresp),
		)

		return DatastoreModel{}, diags
	}

	datastore := response.Datastore
	if datastore == nil {
		diags.AddError(
			"update datastore resource",
			"datastore "+name+": could not get datastore from response",
		)

		return DatastoreModel{}, diags
	}
	responseId := datastore.Id
	if responseId == 0 {
		diags.AddError(
			"update datastore resource",
			"datastore "+name+": could not get id",
		)

		return DatastoreModel{}, diags
	}

	if responseId != id {
		diags.AddError(
			"update datastore resource",
			"datastore "+name+": id mismatch "+fmt.Sprintf("%d != %d", id, responseId),
		)

		return DatastoreModel{}, diags
	}

	// Wait for the datastore to be ready
	waitForReady := func() (string, error) {
		response, hresp, err := client.DatastoresAPI.GetDatastores(ctx, id).Execute()
		if err != nil {
			if hresp == nil || hresp.StatusCode != http.StatusOK {
				return "", backoff.Permanent(err)
			}
		}

		if response == nil || response.Datastore == nil {
			return "", backoff.Permanent(fmt.Errorf("missing datastore in response"))
		}

		status := response.Datastore.Status

		return status, checkStatusDone(
			status,
			CreateTargetStatuses,
			CreateErrorStatuses,
		)
	}

	if status, err := backoff.Retry(
		ctx,
		waitForReady,
		backoff.WithBackOff(backoff.NewConstantBackOff(5*time.Second)),
		backoff.WithMaxElapsedTime(45*time.Minute),
	); err != nil {
		diags.AddError(
			"update datastore resource",
			fmt.Sprintf(
				"datastore %d: provisioning failed current status is: %v",
				id,
				status,
			),
		)
	}

	if diags.HasError() {
		return state, diags
	}

	state, pdiags := getDatastoreAsState(ctx, id, plan, client)
	if pdiags.HasError() {
		diags.Append(pdiags...)
		diags.AddError(
			"update datastore resource",
			fmt.Sprintf("datastore %d: failed to read from api", id),
		)

		return state, diags
	}

	return state, diags
}

func (r *Resource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	var plan, state DatastoreModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	name := plan.Name.ValueString()

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"update datastore resource",
			"datastore "+name+": failed to create client: "+err.Error(),
		)

		return
	}

	state, diags := updateDatastore(ctx, state.Id.ValueInt64(), plan, state, client)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)

		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
