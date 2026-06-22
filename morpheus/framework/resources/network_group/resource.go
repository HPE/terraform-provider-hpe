package network_group

import (
	"context"
	"fmt"
	"strconv"

	sdk "github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

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
	resp.TypeName = req.ProviderTypeName + "_morpheus_network_group"
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

	additionalProps := map[string]interface{}{}
	if !plan.Visibility.IsNull() && !plan.Visibility.IsUnknown() {
		additionalProps["visibility"] = plan.Visibility.ValueString()
	}
	if !plan.Active.IsNull() && !plan.Active.IsUnknown() {
		additionalProps["active"] = plan.Active.ValueBool()
	}
	if len(additionalProps) > 0 {
		body.AdditionalProperties = additionalProps
	}

	createReq := sdk.CreateNetworkGroupRequest{
		NetworkGroup: &body,
	}

	// tenant_ids and resource_permissions go at top-level of request (not inside networkGroup)
	reqAdditionalProps := map[string]interface{}{}

	if !plan.TenantIds.IsNull() && !plan.TenantIds.IsUnknown() {
		var tenantIDs []int64
		resp.Diagnostics.Append(plan.TenantIds.ElementsAs(ctx, &tenantIDs, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		reqAdditionalProps["tenantPermissions"] = map[string]interface{}{
			"accounts": tenantIDs,
		}
	}

	if !plan.ResourcePermissionGroupsAll.IsNull() && !plan.ResourcePermissionGroupsAll.IsUnknown() {
		rpMap := map[string]interface{}{
			"all": plan.ResourcePermissionGroupsAll.ValueBool(),
		}
		if !plan.ResourcePermissionGroupIds.IsNull() && !plan.ResourcePermissionGroupIds.IsUnknown() {
			var groupIDs []int64
			resp.Diagnostics.Append(plan.ResourcePermissionGroupIds.ElementsAs(ctx, &groupIDs, false)...)
			if resp.Diagnostics.HasError() {
				return
			}
			sites := make([]map[string]interface{}, 0, len(groupIDs))
			for _, gid := range groupIDs {
				sites = append(sites, map[string]interface{}{"id": gid})
			}
			rpMap["sites"] = sites
		}
		reqAdditionalProps["resourcePermissions"] = rpMap
	}

	if len(reqAdditionalProps) > 0 {
		createReq.AdditionalProperties = reqAdditionalProps
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

	additionalProps := map[string]interface{}{}
	if !plan.Visibility.IsNull() && !plan.Visibility.IsUnknown() {
		additionalProps["visibility"] = plan.Visibility.ValueString()
	}
	if !plan.Active.IsNull() && !plan.Active.IsUnknown() {
		additionalProps["active"] = plan.Active.ValueBool()
	}
	if len(additionalProps) > 0 {
		body.AdditionalProperties = additionalProps
	}

	updateReq := sdk.UpdateNetworkGroupRequest{
		NetworkGroup: &body,
	}

	// tenant_ids and resource_permissions go at top-level of request (not inside networkGroup)
	updateReqAdditionalProps := map[string]interface{}{}

	if !plan.TenantIds.IsNull() && !plan.TenantIds.IsUnknown() {
		var tenantIDs []int64
		resp.Diagnostics.Append(plan.TenantIds.ElementsAs(ctx, &tenantIDs, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		updateReqAdditionalProps["tenantPermissions"] = map[string]interface{}{
			"accounts": tenantIDs,
		}
	}

	if !plan.ResourcePermissionGroupsAll.IsNull() && !plan.ResourcePermissionGroupsAll.IsUnknown() {
		rpMap := map[string]interface{}{
			"all": plan.ResourcePermissionGroupsAll.ValueBool(),
		}
		if !plan.ResourcePermissionGroupIds.IsNull() && !plan.ResourcePermissionGroupIds.IsUnknown() {
			var groupIDs []int64
			resp.Diagnostics.Append(plan.ResourcePermissionGroupIds.ElementsAs(ctx, &groupIDs, false)...)
			if resp.Diagnostics.HasError() {
				return
			}
			sites := make([]map[string]interface{}, 0, len(groupIDs))
			for _, gid := range groupIDs {
				sites = append(sites, map[string]interface{}{"id": gid})
			}
			rpMap["sites"] = sites
		}
		updateReqAdditionalProps["resourcePermissions"] = rpMap
	}

	if len(updateReqAdditionalProps) > 0 {
		updateReq.AdditionalProperties = updateReqAdditionalProps
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

	// tenant_ids: returned as group.Tenants[*].Id
	model.TenantIds = types.SetNull(types.Int64Type)
	if len(group.Tenants) > 0 {
		vals := make([]attr.Value, 0, len(group.Tenants))
		for _, t := range group.Tenants {
			if t.Id != nil {
				vals = append(vals, types.Int64Value(*t.Id))
			}
		}
		if len(vals) > 0 {
			tenantSet, _ := types.SetValue(types.Int64Type, vals)
			model.TenantIds = tenantSet
		}
	}

	// resource_permissions: returned as group.AdditionalProperties["resourcePermission"] (singular)
	model.ResourcePermissionGroupsAll = types.BoolNull()
	model.ResourcePermissionGroupIds = types.SetNull(types.Int64Type)
	if rp, ok := group.AdditionalProperties["resourcePermission"]; ok {
		if rpMap, ok := rp.(map[string]interface{}); ok {
			if all, ok := rpMap["all"].(bool); ok {
				model.ResourcePermissionGroupsAll = types.BoolValue(all)
			}
			if sites, ok := rpMap["sites"].([]interface{}); ok {
				siteVals := make([]attr.Value, 0, len(sites))
				for _, s := range sites {
					if sMap, ok := s.(map[string]interface{}); ok {
						if id, ok := sMap["id"].(float64); ok {
							siteVals = append(siteVals, types.Int64Value(int64(id)))
						}
					}
				}
				if len(siteVals) > 0 {
					siteSet, _ := types.SetValue(types.Int64Type, siteVals)
					model.ResourcePermissionGroupIds = siteSet
				}
			}
		}
	}
}
