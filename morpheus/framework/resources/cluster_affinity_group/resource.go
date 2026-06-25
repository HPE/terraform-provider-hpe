package cluster_affinity_group

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	sdk "github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

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
	resp.TypeName = req.ProviderTypeName + "_morpheus_cluster_affinity_group"
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
	ag.Visibility = plan.Visibility.ValueStringPointer()

	if !plan.TenantIds.IsNull() && !plan.TenantIds.IsUnknown() {
		var tenantIDs []types.Int64
		resp.Diagnostics.Append(plan.TenantIds.ElementsAs(ctx, &tenantIDs, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		for _, idVal := range tenantIDs {
			if !idVal.IsNull() {
				tenantID := idVal.ValueInt64()
				ag.Tenants = append(ag.Tenants, sdk.SaveClusterAffinityGroupRequestAffinityGroupTenantsInner{Id: &tenantID})
			}
		}
	}

	if !plan.ResourcePermissions.IsNull() && !plan.ResourcePermissions.IsUnknown() {
		rp := &sdk.SaveClusterAffinityGroupRequestAffinityGroupResourcePermissions{
			All: plan.ResourcePermissions.All.ValueBoolPointer(),
		}
		if !plan.ResourcePermissions.GroupIds.IsNull() && !plan.ResourcePermissions.GroupIds.IsUnknown() {
			var groupIDs []types.Int64
			resp.Diagnostics.Append(plan.ResourcePermissions.GroupIds.ElementsAs(ctx, &groupIDs, false)...)
			if resp.Diagnostics.HasError() {
				return
			}
			for _, g := range groupIDs {
				if !g.IsNull() {
					gid := g.ValueInt64()
					rp.Sites = append(rp.Sites, map[string]interface{}{"id": gid})
				}
			}
		}
		ag.ResourcePermissions = rp
	}

	body := sdk.SaveClusterAffinityGroupRequest{
		AffinityGroup: &ag,
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
	resp.Diagnostics.Append(mapGetResponseToModel(ctx, &plan, readAg)...)
	if resp.Diagnostics.HasError() {
		return
	}

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

	resp.Diagnostics.Append(mapGetResponseToModel(ctx, &state, ag)...)
	if resp.Diagnostics.HasError() {
		return
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

	ag := sdk.UpdateCloudAffinityGroupRequestAffinityGroup{}
	if !plan.Name.IsNull() {
		v := plan.Name.ValueString()
		ag.Name = &v
	}
	ag.Active = plan.Active.ValueBoolPointer()
	ag.Visibility = plan.Visibility.ValueStringPointer()

	if !plan.TenantIds.IsNull() && !plan.TenantIds.IsUnknown() {
		var tenantIDs []types.Int64
		resp.Diagnostics.Append(plan.TenantIds.ElementsAs(ctx, &tenantIDs, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		for _, idVal := range tenantIDs {
			if !idVal.IsNull() {
				tenantID := idVal.ValueInt64()
				ag.Tenants = append(ag.Tenants, sdk.UpdateCloudAffinityGroupRequestAffinityGroupTenantsInner{Id: &tenantID})
			}
		}
	}

	if !plan.ResourcePermissions.IsNull() && !plan.ResourcePermissions.IsUnknown() {
		rp := &sdk.UpdateCloudAffinityGroupRequestAffinityGroupResourcePermissions{
			All: plan.ResourcePermissions.All.ValueBoolPointer(),
		}
		if !plan.ResourcePermissions.GroupIds.IsNull() && !plan.ResourcePermissions.GroupIds.IsUnknown() {
			var groupIDs []types.Int64
			resp.Diagnostics.Append(plan.ResourcePermissions.GroupIds.ElementsAs(ctx, &groupIDs, false)...)
			if resp.Diagnostics.HasError() {
				return
			}
			for _, g := range groupIDs {
				if !g.IsNull() {
					gid := g.ValueInt64()
					rp.Sites = append(rp.Sites, map[string]interface{}{"id": gid})
				}
			}
		}
		ag.ResourcePermissions = rp
	}

	body := sdk.UpdateClusterAffinityGroupRequest{
		AffinityGroup: &ag,
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
	resp.Diagnostics.Append(mapGetResponseToModel(ctx, &plan, readAg)...)
	if resp.Diagnostics.HasError() {
		return
	}

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
	ctx context.Context,
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

	// tenant_ids
	model.TenantIds = types.SetNull(types.Int64Type)
	if len(ag.Tenants) > 0 {
		var tenantValues []attr.Value
		for _, t := range ag.Tenants {
			if t.Id != nil {
				tenantValues = append(tenantValues, types.Int64Value(*t.Id))
			}
		}
		if len(tenantValues) > 0 {
			tenantSet, d := types.SetValue(types.Int64Type, tenantValues)
			diags.Append(d...)
			if diags.HasError() {
				return diags
			}
			model.TenantIds = tenantSet
		}
	}

	// resource_permissions
	if ag.ResourcePermissions != nil {
		rp, d := convertAffinityGroupResourcePermissions(ctx, ag.ResourcePermissions)
		diags.Append(d...)
		if diags.HasError() {
			return diags
		}
		model.ResourcePermissions = rp
	} else {
		model.ResourcePermissions = NewResourcePermissionsValueNull()
	}

	return diags
}

func convertAffinityGroupResourcePermissions(
	ctx context.Context,
	rp *sdk.GetClusterAffinityGroup200ResponseAffinityGroupResourcePermissions,
) (ResourcePermissionsValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	var groupValues []attr.Value
	for _, site := range rp.Sites {
		if rawID, ok := site["id"]; ok {
			switch v := rawID.(type) {
			case float64:
				groupValues = append(groupValues, types.Int64Value(int64(v)))
			case int64:
				groupValues = append(groupValues, types.Int64Value(v))
			case *int64:
				if v != nil {
					groupValues = append(groupValues, types.Int64Value(*v))
				}
			}
		}
	}

	var groupIDsSet attr.Value
	if len(groupValues) > 0 {
		groupIDsSet, _ = types.SetValue(types.Int64Type, groupValues)
	} else {
		groupIDsSet = types.SetNull(types.Int64Type)
	}

	result, d := NewResourcePermissionsValue(
		ResourcePermissionsValue{}.AttributeTypes(ctx),
		map[string]attr.Value{
			"all":       types.BoolValue(rp.All != nil && *rp.All),
			"group_ids": groupIDsSet,
		},
	)
	diags.Append(d...)

	return result, diags
}
