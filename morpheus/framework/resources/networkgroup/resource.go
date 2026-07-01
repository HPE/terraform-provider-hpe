package networkgroup

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"

	"github.com/HPE/terraform-provider-hpe/morpheus/configure"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/cleanup"
)

var (
	_ resource.Resource                = &networkGroupResource{}
	_ resource.ResourceWithConfigure   = &networkGroupResource{}
	_ resource.ResourceWithImportState = &networkGroupResource{}
)

type networkGroupResource struct {
	configure.ResourceWithMorpheusConfigure
}

func NewResource() resource.Resource {
	return &networkGroupResource{}
}

func (r *networkGroupResource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_" + "network_group"
}

func (r *networkGroupResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = NetworkGroupResourceSchema(ctx)
}

func (r *networkGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var plan NetworkGroupModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body := sdk.CreateNetworkGroupRequestNetworkGroup{
		Name: plan.Name.ValueStringPointer(),
	}
	if !plan.Description.IsNull() {
		body.Description = plan.Description.ValueStringPointer()
	}
	if !plan.Visibility.IsNull() && !plan.Visibility.IsUnknown() {
		body.Visibility = plan.Visibility.ValueStringPointer()
	}
	if !plan.Active.IsNull() && !plan.Active.IsUnknown() {
		body.Active = plan.Active.ValueBoolPointer()
	}

	createReq := sdk.CreateNetworkGroupRequest{NetworkGroup: &body}
	if !plan.TenantIds.IsNull() && !plan.TenantIds.IsUnknown() {
		var ids []int64
		resp.Diagnostics.Append(plan.TenantIds.ElementsAs(ctx, &ids, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		createReq.TenantPermissions = &sdk.CreateNetworkGroupRequestTenantPermissions{Accounts: ids}
	}
	if !plan.ResourcePermissions.IsNull() && !plan.ResourcePermissions.IsUnknown() {
		rp := sdk.CreateNetworkGroupRequestResourcePermissions{}
		rp.All = plan.ResourcePermissions.All.ValueBoolPointer()
		if !plan.ResourcePermissions.GroupIds.IsNull() && !plan.ResourcePermissions.GroupIds.IsUnknown() {
			var siteIDs []int64
			resp.Diagnostics.Append(plan.ResourcePermissions.GroupIds.ElementsAs(ctx, &siteIDs, false)...)
			if resp.Diagnostics.HasError() {
				return
			}
			sites := make([]sdk.CreateNetworkGroupRequestResourcePermissionsSitesInner, 0, len(siteIDs))
			for i := range siteIDs {
				id := siteIDs[i]
				sites = append(sites, sdk.CreateNetworkGroupRequestResourcePermissionsSitesInner{Id: &id})
			}
			rp.Sites = sites
		}
		if !plan.ResourcePermissions.AllPlans.IsNull() && !plan.ResourcePermissions.AllPlans.IsUnknown() {
			rp.AllPlans = plan.ResourcePermissions.AllPlans.ValueBoolPointer()
		}
		if !plan.ResourcePermissions.PlanIds.IsNull() && !plan.ResourcePermissions.PlanIds.IsUnknown() {
			var planIDs []int64
			resp.Diagnostics.Append(plan.ResourcePermissions.PlanIds.ElementsAs(ctx, &planIDs, false)...)
			if resp.Diagnostics.HasError() {
				return
			}
			plans := make([]sdk.CreateNetworkGroupRequestResourcePermissionsPlansInner, 0, len(planIDs))
			for i := range planIDs {
				id := planIDs[i]
				plans = append(plans, sdk.CreateNetworkGroupRequestResourcePermissionsPlansInner{Id: &id})
			}
			rp.Plans = plans
		}
		createReq.ResourcePermissions = &rp
	}

	result, httpResp, err := client.NetworksAPI.CreateNetworkGroup(ctx).
		CreateNetworkGroupRequest(createReq).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpCreate, "network_group", plan.Name.ValueString(), err, httpResp)

		return
	}

	// The API returns the id inside {"networkGroup":{"id":...}} but the SDK model
	// only exposes a top-level Id field. Extract from AdditionalProperties.
	var newID int64
	if ngData, ok := result.AdditionalProperties["networkGroup"]; ok {
		if ngMap, ok := ngData.(map[string]interface{}); ok {
			if id, ok := ngMap["id"].(float64); ok {
				newID = int64(id)
			}
		}
	}
	if newID == 0 && result.Id.IsSet() && result.Id.Get() != nil {
		newID = *result.Id.Get()
	}
	if newID == 0 {
		resp.Diagnostics.AddError("Failed to extract ID", "Could not determine network group ID from create response")

		return
	}

	// Re-read to get full object
	readResult, httpResp, err := client.NetworksAPI.GetNetworkGroup(ctx, newID).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpRead, "network_group", plan.Name.ValueString(), err, httpResp)
		cleanup.TaintResourceState(ctx, cleanup.TaintResourceStateConfig{
			ResourceType: "network_group",
			ResourceID:   newID,
			StateWriter:  &resp.State,
			Diagnostics:  &resp.Diagnostics,
		})

		return
	}

	group := readResult.NetworkGroup
	if group == nil {
		resp.Diagnostics.AddError("API returned nil", "NetworkGroup is nil in the response")

		return
	}
	mapResponseToModel(&plan, group)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *networkGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var state NetworkGroupModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Detect import: ImportState sets only id; name is null.
	// On normal refresh, name is always a known string from prior state.
	isImport := state.Name.IsNull()

	id := state.Id.ValueInt64()

	result, httpResp, err := client.NetworksAPI.GetNetworkGroup(ctx, id).Execute()
	if errfmt.IsNotFound(httpResp) {
		resp.State.RemoveResource(ctx)

		return
	}
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpRead, "network_group", "", err, httpResp)

		return
	}

	group := result.NetworkGroup
	if group == nil {
		resp.Diagnostics.AddError("API returned nil", "NetworkGroup is nil in the response")

		return
	}
	mapResponseToModel(&state, group)

	// mapResponseToModel never touches tenant_ids or resource_permissions, so on
	// a normal refresh those fields carry forward from the prior state naturally.
	// On import there is no prior state, so we explicitly populate them from the
	// API response.
	if isImport {
		state.TenantIds = networkGroupTenantIdsFromAPI(group.Tenants)
		state.ResourcePermissions = networkGroupResourcePermissionsFromAPI(group.ResourcePermission)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *networkGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var plan NetworkGroupModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := plan.Id.ValueInt64()

	body := sdk.UpdateNetworkGroupRequestNetworkGroup{
		Name: plan.Name.ValueStringPointer(),
	}
	if !plan.Description.IsNull() {
		body.Description = plan.Description.ValueStringPointer()
	}
	if !plan.Visibility.IsNull() && !plan.Visibility.IsUnknown() {
		body.Visibility = plan.Visibility.ValueStringPointer()
	}
	if !plan.Active.IsNull() && !plan.Active.IsUnknown() {
		body.Active = plan.Active.ValueBoolPointer()
	}

	updateReq := sdk.UpdateNetworkGroupRequest{NetworkGroup: &body}
	if !plan.TenantIds.IsNull() && !plan.TenantIds.IsUnknown() {
		var ids []int64
		resp.Diagnostics.Append(plan.TenantIds.ElementsAs(ctx, &ids, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		updateReq.TenantPermissions = &sdk.UpdateNetworkGroupRequestTenantPermissions{Accounts: ids}
	}
	if !plan.ResourcePermissions.IsNull() && !plan.ResourcePermissions.IsUnknown() {
		rp := sdk.UpdateNetworkGroupRequestResourcePermissions{}
		rp.All = plan.ResourcePermissions.All.ValueBoolPointer()
		if !plan.ResourcePermissions.GroupIds.IsNull() && !plan.ResourcePermissions.GroupIds.IsUnknown() {
			var siteIDs []int64
			resp.Diagnostics.Append(plan.ResourcePermissions.GroupIds.ElementsAs(ctx, &siteIDs, false)...)
			if resp.Diagnostics.HasError() {
				return
			}
			sites := make([]sdk.UpdateNetworkGroupRequestResourcePermissionsSitesInner, 0, len(siteIDs))
			for i := range siteIDs {
				sid := siteIDs[i]
				sites = append(sites, sdk.UpdateNetworkGroupRequestResourcePermissionsSitesInner{Id: &sid})
			}
			rp.Sites = sites
		}
		if !plan.ResourcePermissions.AllPlans.IsNull() && !plan.ResourcePermissions.AllPlans.IsUnknown() {
			rp.AllPlans = plan.ResourcePermissions.AllPlans.ValueBoolPointer()
		}
		if !plan.ResourcePermissions.PlanIds.IsNull() && !plan.ResourcePermissions.PlanIds.IsUnknown() {
			var planIDs []int64
			resp.Diagnostics.Append(plan.ResourcePermissions.PlanIds.ElementsAs(ctx, &planIDs, false)...)
			if resp.Diagnostics.HasError() {
				return
			}
			plans := make([]sdk.UpdateNetworkGroupRequestResourcePermissionsPlansInner, 0, len(planIDs))
			for i := range planIDs {
				pid := planIDs[i]
				plans = append(plans, sdk.UpdateNetworkGroupRequestResourcePermissionsPlansInner{Id: &pid})
			}
			rp.Plans = plans
		}
		updateReq.ResourcePermissions = &rp
	}

	_, httpResp, err := client.NetworksAPI.UpdateNetworkGroup(ctx, id).
		UpdateNetworkGroupRequest(updateReq).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpUpdate, "network_group", plan.Name.ValueString(), err, httpResp)

		return
	}

	// Re-read to get full object
	readResult, httpResp, err := client.NetworksAPI.GetNetworkGroup(ctx, id).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpRead, "network_group", "", err, httpResp)

		return
	}

	group := readResult.NetworkGroup
	if group == nil {
		resp.Diagnostics.AddError("API returned nil", "NetworkGroup is nil in the response")

		return
	}
	mapResponseToModel(&plan, group)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *networkGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var state NetworkGroupModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.Id.ValueInt64()

	_, httpResp, err := client.NetworksAPI.DeleteNetworkGroup(ctx, id).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpDelete, "network_group", "", err, httpResp)

		return
	}
}

func (r *networkGroupResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	id, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Could not parse ID %q as integer: %s", req.ID, err))

		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}

func mapResponseToModel(model *NetworkGroupModel, group *sdk.GetNetworkGroup200ResponseNetworkGroup) {
	if group.Id != nil {
		model.Id = types.Int64Value(*group.Id)
	}
	if group.Name != nil {
		model.Name = types.StringValue(*group.Name)
	}
	if group.Description != nil {
		model.Description = types.StringValue(*group.Description)
	} else {
		model.Description = types.StringNull()
	}
	if group.Visibility != nil {
		model.Visibility = types.StringValue(*group.Visibility)
	}
	if group.Active != nil {
		model.Active = types.BoolValue(*group.Active)
	}
}

// networkGroupTenantIdsFromAPI converts the Tenants slice from a GET response
// into the types.Set used by the Terraform model. Used only on import.
func networkGroupTenantIdsFromAPI(
	tenants []sdk.GetNetworkGroup200ResponseNetworkGroupTenantsInner,
) types.Set {
	vals := make([]attr.Value, 0, len(tenants))
	for _, t := range tenants {
		if t.Id != nil {
			vals = append(vals, types.Int64Value(*t.Id))
		}
	}

	return types.SetValueMust(types.Int64Type, vals)
}

// networkGroupResourcePermissionsFromAPI converts the ResourcePermission object
// from a GET response into the ResourcePermissionsValue used by the Terraform
// model. Used only on import.
func networkGroupResourcePermissionsFromAPI(
	rp *sdk.GetNetworkGroup200ResponseNetworkGroupResourcePermission,
) ResourcePermissionsValue {
	if rp == nil {
		return NewResourcePermissionsValueNull()
	}

	siteVals := make([]attr.Value, 0, len(rp.Sites))
	for _, s := range rp.Sites {
		if s.Id != nil {
			siteVals = append(siteVals, types.Int64Value(*s.Id))
		}
	}

	planVals := make([]attr.Value, 0, len(rp.Plans))
	for _, p := range rp.Plans {
		if p.Id != nil {
			planVals = append(planVals, types.Int64Value(*p.Id))
		}
	}

	return NewResourcePermissionsValueMust(
		map[string]attr.Type{
			"all":       types.BoolType,
			"all_plans": types.BoolType,
			"group_ids": types.ListType{ElemType: types.Int64Type},
			"plan_ids":  types.ListType{ElemType: types.Int64Type},
		},
		map[string]attr.Value{
			"all":       types.BoolPointerValue(rp.All),
			"all_plans": types.BoolPointerValue(rp.AllPlans.Get()),
			"group_ids": types.ListValueMust(types.Int64Type, siteVals),
			"plan_ids":  types.ListValueMust(types.Int64Type, planVals),
		},
	)
}
