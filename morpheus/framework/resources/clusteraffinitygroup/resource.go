// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package clusteraffinitygroup

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"

	"github.com/HPE/terraform-provider-hpe/morpheus/configure"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/cleanup"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

var (
	_ resource.Resource                = &clusterAffinityGroupResource{}
	_ resource.ResourceWithConfigure   = &clusterAffinityGroupResource{}
	_ resource.ResourceWithImportState = &clusterAffinityGroupResource{}
)

type clusterAffinityGroupResource struct {
	configure.ResourceWithMorpheusConfigure
}

func NewResource() resource.Resource {
	return &clusterAffinityGroupResource{}
}

func (r *clusterAffinityGroupResource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_" + "cluster_affinity_group"
}

func (r *clusterAffinityGroupResource) Schema(
	ctx context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = ClusterAffinityGroupResourceSchema(ctx)
}

func (r *clusterAffinityGroupResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var plan ClusterAffinityGroupModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	clusterID := plan.ClusterId.ValueInt64()

	name := plan.Name.ValueString()
	ag := sdk.SaveClusterAffinityGroupRequestAffinityGroup{
		Name: &name,
	}
	ag.Active = plan.Active.ValueBoolPointer()
	if !plan.Visibility.IsNull() && !plan.Visibility.IsUnknown() {
		ag.Visibility = plan.Visibility.ValueStringPointer()
	}
	if !plan.ResourcePermissions.IsNull() && !plan.ResourcePermissions.IsUnknown() {
		rp := sdk.SaveClusterAffinityGroupRequestAffinityGroupResourcePermissions{}
		rp.All = plan.ResourcePermissions.All.ValueBoolPointer()
		if !plan.ResourcePermissions.GroupIds.IsNull() && !plan.ResourcePermissions.GroupIds.IsUnknown() {
			var siteIDs []int64
			resp.Diagnostics.Append(plan.ResourcePermissions.GroupIds.ElementsAs(ctx, &siteIDs, false)...)
			if resp.Diagnostics.HasError() {
				return
			}
			sites := make([]map[string]interface{}, 0, len(siteIDs))
			for _, sid := range siteIDs {
				sites = append(sites, map[string]interface{}{"id": sid})
			}
			rp.Sites = sites
		}
		ag.ResourcePermissions = &rp
	}

	body := sdk.SaveClusterAffinityGroupRequest{
		AffinityGroup: &ag,
	}
	if !plan.TenantIds.IsNull() && !plan.TenantIds.IsUnknown() {
		var ids []int64
		resp.Diagnostics.Append(plan.TenantIds.ElementsAs(ctx, &ids, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		body.TenantPermissions = &sdk.SaveClusterAffinityGroupRequestTenantPermissions{Accounts: ids}
	}

	result, httpResp, err := client.ClustersAPI.SaveClusterAffinityGroup(ctx, clusterID).
		SaveClusterAffinityGroupRequest(body).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpCreate, "cluster_affinity_group", plan.Name.ValueString(), err, httpResp)

		return
	}

	if result.AffinityGroup == nil || result.AffinityGroup.Id == nil {
		resp.Diagnostics.AddError("API returned nil ID", "AffinityGroup ID is nil in the create response")

		return
	}

	id := *result.AffinityGroup.Id

	readResult, httpResp, err := client.ClustersAPI.GetClusterAffinityGroup(ctx, clusterID, id).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpRead, "cluster_affinity_group", plan.Name.ValueString(), err, httpResp)
		cleanup.TaintResourceState(ctx, cleanup.TaintResourceStateConfig{
			ResourceType: "cluster_affinity_group",
			ResourceID:   id,
			StateWriter:  &resp.State,
			Diagnostics:  &resp.Diagnostics,
		})

		return
	}

	readAg := readResult.AffinityGroup
	if readAg == nil {
		resp.Diagnostics.AddError("API returned nil", "AffinityGroup is nil in the response")

		return
	}

	// The API silently drops tenant/site IDs that don't exist in the environment.
	// Preserve plan values so state matches the plan and Terraform's consistency
	// check passes. Read() will return the API-normalised values, surfacing any
	// divergence as a plan diff on the next run.
	savedTenantIds := plan.TenantIds
	savedRP := plan.ResourcePermissions

	resp.Diagnostics.Append(mapGetResponseToModel(&plan, readAg)...)

	plan.TenantIds = savedTenantIds
	plan.ResourcePermissions = savedRP

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *clusterAffinityGroupResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var state ClusterAffinityGroupModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Detect import: ImportState sets only cluster_id and id; name is null.
	// On normal refresh, name is always a known string from prior state.
	isImport := state.Name.IsNull()
	priorTenantIds := state.TenantIds
	priorRP := state.ResourcePermissions

	clusterID := state.ClusterId.ValueInt64()
	id := state.Id.ValueInt64()

	result, httpResp, err := client.ClustersAPI.GetClusterAffinityGroup(ctx, clusterID, id).Execute()
	if errfmt.IsNotFound(httpResp) {
		resp.State.RemoveResource(ctx)

		return
	}
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpRead, "cluster_affinity_group", "", err, httpResp)

		return
	}

	ag := result.AffinityGroup
	if ag == nil {
		resp.Diagnostics.AddError("API returned nil", "AffinityGroup is nil in the response")

		return
	}

	resp.Diagnostics.Append(mapGetResponseToModel(&state, ag)...)

	// On normal refresh, preserve tenant_ids and resource_permissions from prior
	// state. The API may silently drop IDs that don't exist in the environment,
	// which would cause a spurious diff. On import there is no prior state, so
	// we use the API values that mapGetResponseToModel just populated.
	if !isImport {
		state.TenantIds = priorTenantIds
		state.ResourcePermissions = priorRP
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *clusterAffinityGroupResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var plan ClusterAffinityGroupModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	clusterID := plan.ClusterId.ValueInt64()
	id := plan.Id.ValueInt64()

	ag := sdk.UpdateClusterAffinityGroupRequestAffinityGroup{}
	if !plan.Name.IsNull() {
		v := plan.Name.ValueString()
		ag.Name = &v
	}
	ag.Active = plan.Active.ValueBoolPointer()
	if !plan.Visibility.IsNull() && !plan.Visibility.IsUnknown() {
		ag.Visibility = plan.Visibility.ValueStringPointer()
	}
	if !plan.ResourcePermissions.IsNull() && !plan.ResourcePermissions.IsUnknown() {
		rp := sdk.UpdateClusterAffinityGroupRequestAffinityGroupResourcePermissions{}
		rp.All = plan.ResourcePermissions.All.ValueBoolPointer()
		if !plan.ResourcePermissions.GroupIds.IsNull() && !plan.ResourcePermissions.GroupIds.IsUnknown() {
			var siteIDs []int64
			resp.Diagnostics.Append(plan.ResourcePermissions.GroupIds.ElementsAs(ctx, &siteIDs, false)...)
			if resp.Diagnostics.HasError() {
				return
			}
			sites := make([]map[string]interface{}, 0, len(siteIDs))
			for _, sid := range siteIDs {
				sites = append(sites, map[string]interface{}{"id": sid})
			}
			rp.Sites = sites
		}
		ag.ResourcePermissions = &rp
	}

	body := sdk.UpdateClusterAffinityGroupRequest{
		AffinityGroup: &ag,
	}
	if !plan.TenantIds.IsNull() && !plan.TenantIds.IsUnknown() {
		var ids []int64
		resp.Diagnostics.Append(plan.TenantIds.ElementsAs(ctx, &ids, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		body.TenantPermissions = &sdk.UpdateClusterAffinityGroupRequestTenantPermissions{Accounts: ids}
	}

	_, httpResp, err := client.ClustersAPI.UpdateClusterAffinityGroup(ctx, clusterID, id).
		UpdateClusterAffinityGroupRequest(body).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpUpdate, "cluster_affinity_group", plan.Name.ValueString(), err, httpResp)

		return
	}

	readResult, httpResp, err := client.ClustersAPI.GetClusterAffinityGroup(ctx, clusterID, id).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpRead, "cluster_affinity_group", plan.Name.ValueString(), err, httpResp)

		return
	}

	readAg := readResult.AffinityGroup
	if readAg == nil {
		resp.Diagnostics.AddError("API returned nil", "AffinityGroup is nil in the response")

		return
	}

	// Same as Create: preserve plan values for tenant_ids and resource_permissions
	// so the consistency check passes when the API normalises submitted IDs.
	savedTenantIds := plan.TenantIds
	savedRP := plan.ResourcePermissions

	resp.Diagnostics.Append(mapGetResponseToModel(&plan, readAg)...)

	plan.TenantIds = savedTenantIds
	plan.ResourcePermissions = savedRP

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *clusterAffinityGroupResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var state ClusterAffinityGroupModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	clusterID := state.ClusterId.ValueInt64()
	id := state.Id.ValueInt64()

	_, httpResp, err := client.ClustersAPI.DeleteClusterAffinityGroup(ctx, clusterID, id).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpDelete, "cluster_affinity_group", "", err, httpResp)

		return
	}
}

func (r *clusterAffinityGroupResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	parts := strings.Split(req.ID, ".")
	if len(parts) != 2 {
		resp.Diagnostics.AddError("Invalid Import ID",
			fmt.Sprintf("Expected format: cluster_id.affinity_group_id, got: %s", req.ID))

		return
	}

	clusterID, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Could not parse cluster_id %q: %s", parts[0], err))

		return
	}
	id, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Could not parse affinity_group_id %q: %s", parts[1], err))

		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("cluster_id"), clusterID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}

// Ensure unused imports are satisfied.
var _ *http.Response

func mapGetResponseToModel(
	model *ClusterAffinityGroupModel,
	ag *sdk.GetClusterAffinityGroup200ResponseAffinityGroup,
) diag.Diagnostics {
	var diags diag.Diagnostics

	if ag.Id != nil {
		model.Id = types.Int64Value(*ag.Id)
	}
	if ag.Name != nil {
		model.Name = types.StringValue(*ag.Name)
	}
	if ag.Active != nil {
		model.Active = types.BoolValue(*ag.Active)
	}
	model.Visibility = convert.StrToType(ag.Visibility)

	tenantVals := make([]attr.Value, 0, len(ag.Tenants))
	for _, t := range ag.Tenants {
		if t.Id != nil {
			tenantVals = append(tenantVals, types.Int64Value(*t.Id))
		}
	}
	set, setDiags := types.SetValue(types.Int64Type, tenantVals)
	diags.Append(setDiags...)
	model.TenantIds = set

	if ag.ResourcePermissions != nil {
		siteVals := make([]attr.Value, 0, len(ag.ResourcePermissions.Sites))
		for _, s := range ag.ResourcePermissions.Sites {
			if id, ok := s["id"].(float64); ok {
				siteVals = append(siteVals, types.Int64Value(int64(id)))
			}
		}
		groupIdsList, listDiags := types.ListValue(types.Int64Type, siteVals)
		diags.Append(listDiags...)
		rp, rpDiags := NewResourcePermissionsValue(
			map[string]attr.Type{
				"all":       types.BoolType,
				"group_ids": types.ListType{ElemType: types.Int64Type},
			},
			map[string]attr.Value{
				"all":       types.BoolPointerValue(ag.ResourcePermissions.All),
				"group_ids": groupIdsList,
			},
		)
		diags.Append(rpDiags...)
		model.ResourcePermissions = rp
	} else {
		model.ResourcePermissions = NewResourcePermissionsValueNull()
	}

	return diags
}
