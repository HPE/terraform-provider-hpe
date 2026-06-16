package cluster_namespace

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
)

var (
	_ resource.Resource                = &clusterNamespaceResource{}
	_ resource.ResourceWithConfigure   = &clusterNamespaceResource{}
	_ resource.ResourceWithImportState = &clusterNamespaceResource{}
)

type clusterNamespaceResource struct {
	configure.ResourceWithMorpheusConfigure
}

func NewResource() resource.Resource {
	return &clusterNamespaceResource{}
}

func (r *clusterNamespaceResource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_morpheus_cluster_namespace"
}

func (r *clusterNamespaceResource) Schema(
	ctx context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = ClusterNamespaceSchema(ctx)
}

func (r *clusterNamespaceResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var plan clusterNamespaceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	clusterID := plan.ClusterID.ValueInt64()

	ns := sdk.AddClusterNamespaceRequestNamespace{
		Name: plan.Name.ValueString(),
	}
	if !plan.Description.IsNull() {
		v := plan.Description.ValueString()
		ns.Description = &v
	}
	if !plan.Active.IsNull() {
		ns.Active = plan.Active.ValueBoolPointer()
	}

	if !plan.ResourcePermissions.IsNull() && !plan.ResourcePermissions.IsUnknown() {
		attrs := plan.ResourcePermissions.Attributes()
		rp := sdk.AddClusterNamespaceRequestNamespaceResourcePermissions{}
		if allVal, ok := attrs["all"].(types.Bool); ok && !allVal.IsNull() {
			rp.All = allVal.ValueBoolPointer()
		}
		if allPlans, ok := attrs["all_plans"].(types.Bool); ok && !allPlans.IsNull() {
			rp.AllPlans = allPlans.ValueBoolPointer()
		}
		if siteIDsVal, ok := attrs["site_ids"].(types.Set); ok && !siteIDsVal.IsNull() {
			var siteIDs []int64
			d := siteIDsVal.ElementsAs(ctx, &siteIDs, false)
			resp.Diagnostics.Append(d...)
			if resp.Diagnostics.HasError() {
				return
			}
			sites := make([]sdk.AddClusterNamespaceRequestNamespaceResourcePermissionsSitesInner, 0, len(siteIDs))
			for _, sid := range siteIDs {
				sid := sid
				sites = append(sites, sdk.AddClusterNamespaceRequestNamespaceResourcePermissionsSitesInner{Id: &sid})
			}
			rp.Sites = sites
		}
		if planIDsVal, ok := attrs["plan_ids"].(types.Set); ok && !planIDsVal.IsNull() {
			var planIDs []int64
			d := planIDsVal.ElementsAs(ctx, &planIDs, false)
			resp.Diagnostics.Append(d...)
			if resp.Diagnostics.HasError() {
				return
			}
			plans := make([]sdk.AddClusterNamespaceRequestNamespaceResourcePermissionsPlansInner, 0, len(planIDs))
			for _, pid := range planIDs {
				pid := pid
				plans = append(plans, sdk.AddClusterNamespaceRequestNamespaceResourcePermissionsPlansInner{Id: &pid})
			}
			rp.Plans = plans
		}
		ns.ResourcePermissions = &rp
	}

	body := sdk.AddClusterNamespaceRequest{
		Namespace: &ns,
	}

	result, httpResp, err := client.ClustersAPI.AddClusterNamespace(ctx, clusterID).
		AddClusterNamespaceRequest(body).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpCreate, "cluster_namespace", plan.Name.ValueString(), err, httpResp)

		return
	}

	if result.Namespace == nil || result.Namespace.Id == nil {
		resp.Diagnostics.AddError("API returned nil ID", "Namespace ID is nil in the create response")

		return
	}

	id := *result.Namespace.Id

	readResult, httpResp, err := client.ClustersAPI.GetClusterNamespace(ctx, clusterID, id).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpRead, "cluster_namespace", plan.Name.ValueString(), err, httpResp)
		cleanup.TaintResourceState(ctx, cleanup.TaintResourceStateConfig{
			ResourceType: "cluster_namespace",
			ResourceID:   id,
			StateWriter:  &resp.State,
			Diagnostics:  &resp.Diagnostics,
		})

		return
	}

	readNs := readResult.Namespace
	if readNs == nil {
		resp.Diagnostics.AddError("API returned nil", "Namespace is nil in the response")

		return
	}
	mapGetResponseToModel(&plan, readNs, &resp.Diagnostics)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *clusterNamespaceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var state clusterNamespaceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	clusterID := state.ClusterID.ValueInt64()
	id := state.ID.ValueInt64()

	result, httpResp, err := client.ClustersAPI.GetClusterNamespace(ctx, clusterID, id).Execute()
	if errfmt.IsNotFound(httpResp) {
		resp.State.RemoveResource(ctx)

		return
	}
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpRead, "cluster_namespace", "", err, httpResp)

		return
	}

	ns := result.Namespace
	if ns == nil {
		resp.Diagnostics.AddError("API returned nil", "Namespace is nil in the response")

		return
	}

	mapGetResponseToModel(&state, ns, &resp.Diagnostics)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *clusterNamespaceResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var plan clusterNamespaceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	clusterID := plan.ClusterID.ValueInt64()
	id := plan.ID.ValueInt64()

	ns := sdk.UpdateClusterNamespaceRequestNamespace{}
	if !plan.Name.IsNull() {
		v := plan.Name.ValueString()
		ns.Name = &v
	}
	if !plan.Description.IsNull() {
		v := plan.Description.ValueString()
		ns.Description = &v
	}
	if !plan.Active.IsNull() {
		ns.Active = plan.Active.ValueBoolPointer()
	}

	if !plan.ResourcePermissions.IsNull() && !plan.ResourcePermissions.IsUnknown() {
		attrs := plan.ResourcePermissions.Attributes()
		rp := sdk.UpdateClusterNamespaceRequestNamespacePermissionsResourcePermissions{}
		if allVal, ok := attrs["all"].(types.Bool); ok && !allVal.IsNull() {
			rp.All = allVal.ValueBoolPointer()
		}
		if allPlans, ok := attrs["all_plans"].(types.Bool); ok && !allPlans.IsNull() {
			rp.AllPlans = allPlans.ValueBoolPointer()
		}
		if siteIDsVal, ok := attrs["site_ids"].(types.Set); ok && !siteIDsVal.IsNull() {
			var siteIDs []int64
			d := siteIDsVal.ElementsAs(ctx, &siteIDs, false)
			resp.Diagnostics.Append(d...)
			if resp.Diagnostics.HasError() {
				return
			}
			sites := make([]sdk.UpdateClusterNamespaceRequestNamespacePermissionsResourcePermissionsSitesInner, 0, len(siteIDs))
			for _, sid := range siteIDs {
				sid := sid
				sites = append(sites, sdk.UpdateClusterNamespaceRequestNamespacePermissionsResourcePermissionsSitesInner{Id: &sid})
			}
			rp.Sites = sites
		}
		if planIDsVal, ok := attrs["plan_ids"].(types.Set); ok && !planIDsVal.IsNull() {
			var planIDs []int64
			d := planIDsVal.ElementsAs(ctx, &planIDs, false)
			resp.Diagnostics.Append(d...)
			if resp.Diagnostics.HasError() {
				return
			}
			plans := make([]sdk.UpdateClusterNamespaceRequestNamespacePermissionsResourcePermissionsPlansInner, 0, len(planIDs))
			for _, pid := range planIDs {
				pid := pid
				plans = append(plans, sdk.UpdateClusterNamespaceRequestNamespacePermissionsResourcePermissionsPlansInner{Id: &pid})
			}
			rp.Plans = plans
		}
		ns.Permissions = &sdk.UpdateClusterNamespaceRequestNamespacePermissions{
			ResourcePermissions: &rp,
		}
	}

	body := sdk.UpdateClusterNamespaceRequest{
		Namespace: &ns,
	}

	_, httpResp, err := client.ClustersAPI.UpdateClusterNamespace(ctx, clusterID, id).
		UpdateClusterNamespaceRequest(body).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpUpdate, "cluster_namespace", plan.Name.ValueString(), err, httpResp)

		return
	}

	readResult, httpResp, err := client.ClustersAPI.GetClusterNamespace(ctx, clusterID, id).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpRead, "cluster_namespace", plan.Name.ValueString(), err, httpResp)

		return
	}

	readNs := readResult.Namespace
	if readNs == nil {
		resp.Diagnostics.AddError("API returned nil", "Namespace is nil in the response")

		return
	}
	mapGetResponseToModel(&plan, readNs, &resp.Diagnostics)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *clusterNamespaceResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var state clusterNamespaceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	clusterID := state.ClusterID.ValueInt64()
	id := state.ID.ValueInt64()

	_, httpResp, err := client.ClustersAPI.DeleteClusterNamespace(ctx, clusterID, id).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(&resp.Diagnostics, errfmt.OpDelete, "cluster_namespace", "", err, httpResp)

		return
	}
}

func (r *clusterNamespaceResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	parts := strings.Split(req.ID, ".")
	if len(parts) != 2 {
		resp.Diagnostics.AddError("Invalid Import ID",
			fmt.Sprintf("Expected format: cluster_id.namespace_id, got: %s", req.ID))

		return
	}

	clusterID, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Could not parse cluster_id %q: %s", parts[0], err))

		return
	}
	id, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Could not parse namespace_id %q: %s", parts[1], err))

		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("cluster_id"), clusterID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}

// Ensure unused imports are satisfied.
var _ *http.Response

func mapGetResponseToModel(model *clusterNamespaceModel, ns *sdk.GetClusterNamespace200ResponseNamespace, diags *diag.Diagnostics) {
	if ns.Id != nil {
		model.ID = types.Int64Value(*ns.Id)
	}
	if ns.Name != nil {
		model.Name = types.StringValue(*ns.Name)
	}
	if ns.Description != nil {
		model.Description = types.StringValue(*ns.Description)
	}
	// NOTE: Active is not in the API GET at all. Config value is preserved in state.

	// resource_permissions
	if ns.Permissions != nil && ns.Permissions.ResourcePermissions != nil {
		rp := ns.Permissions.ResourcePermissions

		allVal := types.BoolNull()
		if rp.AllGroups != nil {
			allVal = types.BoolValue(*rp.AllGroups)
		} else if rp.All != nil {
			allVal = types.BoolValue(*rp.All)
		}

		allPlansVal := types.BoolNull()
		if rp.AllPlans != nil {
			allPlansVal = types.BoolValue(*rp.AllPlans)
		}

		siteIDSet := types.SetNull(types.Int64Type)
		if len(rp.Sites) > 0 {
			siteVals := make([]attr.Value, 0, len(rp.Sites))
			for _, s := range rp.Sites {
				if s.Id != nil {
					siteVals = append(siteVals, types.Int64Value(*s.Id))
				}
			}
			if len(siteVals) > 0 {
				siteIDSet, _ = types.SetValue(types.Int64Type, siteVals)
			}
		}

		planIDSet := types.SetNull(types.Int64Type)
		if len(rp.Plans) > 0 {
			planVals := make([]attr.Value, 0, len(rp.Plans))
			for _, p := range rp.Plans {
				if p.Id != nil {
					planVals = append(planVals, types.Int64Value(*p.Id))
				}
			}
			if len(planVals) > 0 {
				planIDSet, _ = types.SetValue(types.Int64Type, planVals)
			}
		}

		permObj, d := types.ObjectValue(ClusterNsPermissionsAttrTypes, map[string]attr.Value{
			"all":      allVal,
			"all_plans": allPlansVal,
			"site_ids": siteIDSet,
			"plan_ids": planIDSet,
		})
		diags.Append(d...)
		model.ResourcePermissions = permObj
	} else {
		model.ResourcePermissions = types.ObjectNull(ClusterNsPermissionsAttrTypes)
	}
}
