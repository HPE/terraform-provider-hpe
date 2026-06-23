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
	resp.Schema = ClusterNamespaceResourceSchema(ctx)
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

	var plan ClusterNamespaceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	clusterID := plan.ClusterId.ValueInt64()

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

	// visibility goes under namespace.visibility
	ns.AdditionalProperties = map[string]interface{}{}
	if !plan.Visibility.IsNull() && !plan.Visibility.IsUnknown() {
		ns.AdditionalProperties["visibility"] = plan.Visibility.ValueString()
	}

	// resource_permissions and tenant_ids both go under namespace.permissions
	// (the typed ns.ResourcePermissions field uses the wrong JSON path "resourcePermissions"
	// directly on namespace; the API reads from namespace.permissions.resourcePermissions)
	permsMap := map[string]interface{}{}

	if !plan.ResourcePermissions.IsNull() && !plan.ResourcePermissions.IsUnknown() {
		rpMap := map[string]interface{}{}
		if !plan.ResourcePermissions.All.IsNull() {
			rpMap["all"] = plan.ResourcePermissions.All.ValueBool()
		}
		if !plan.ResourcePermissions.AllPlans.IsNull() {
			rpMap["allPlans"] = plan.ResourcePermissions.AllPlans.ValueBool()
		}
		if !plan.ResourcePermissions.SiteIds.IsNull() {
			var siteIDs []int64
			d := plan.ResourcePermissions.SiteIds.ElementsAs(ctx, &siteIDs, false)
			resp.Diagnostics.Append(d...)
			if resp.Diagnostics.HasError() {
				return
			}
			sites := make([]map[string]interface{}, 0, len(siteIDs))
			for _, sid := range siteIDs {
				sites = append(sites, map[string]interface{}{"id": sid})
			}
			rpMap["sites"] = sites
		}
		if !plan.ResourcePermissions.PlanIds.IsNull() {
			var planIDs []int64
			d := plan.ResourcePermissions.PlanIds.ElementsAs(ctx, &planIDs, false)
			resp.Diagnostics.Append(d...)
			if resp.Diagnostics.HasError() {
				return
			}
			plans := make([]map[string]interface{}, 0, len(planIDs))
			for _, pid := range planIDs {
				plans = append(plans, map[string]interface{}{"id": pid})
			}
			rpMap["plans"] = plans
		}
		permsMap["resourcePermissions"] = rpMap
	}

	if !plan.TenantIds.IsNull() && !plan.TenantIds.IsUnknown() {
		var tenantIDs []int64
		d := plan.TenantIds.ElementsAs(ctx, &tenantIDs, false)
		resp.Diagnostics.Append(d...)
		if resp.Diagnostics.HasError() {
			return
		}
		permsMap["tenantPermissions"] = map[string]interface{}{
			"accounts": tenantIDs,
		}
	}

	if len(permsMap) > 0 {
		ns.AdditionalProperties["permissions"] = permsMap
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
	mapGetResponseToModel(ctx, &plan, readNs, &resp.Diagnostics)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *clusterNamespaceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	var state ClusterNamespaceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	clusterID := state.ClusterId.ValueInt64()
	id := state.Id.ValueInt64()

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

	mapGetResponseToModel(ctx, &state, ns, &resp.Diagnostics)

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

	var plan ClusterNamespaceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	clusterID := plan.ClusterId.ValueInt64()
	id := plan.Id.ValueInt64()

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

	// visibility
	ns.AdditionalProperties = map[string]interface{}{}
	if !plan.Visibility.IsNull() && !plan.Visibility.IsUnknown() {
		ns.AdditionalProperties["visibility"] = plan.Visibility.ValueString()
	}

	// resource_permissions and tenant_ids both go under namespace.permissions
	if !plan.ResourcePermissions.IsNull() && !plan.ResourcePermissions.IsUnknown() {
		rp := sdk.UpdateClusterNamespaceRequestNamespacePermissionsResourcePermissions{}
		if !plan.ResourcePermissions.All.IsNull() {
			rp.All = plan.ResourcePermissions.All.ValueBoolPointer()
		}
		if !plan.ResourcePermissions.AllPlans.IsNull() {
			rp.AllPlans = plan.ResourcePermissions.AllPlans.ValueBoolPointer()
		}
		if !plan.ResourcePermissions.SiteIds.IsNull() {
			var siteIDs []int64
			d := plan.ResourcePermissions.SiteIds.ElementsAs(ctx, &siteIDs, false)
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
		if !plan.ResourcePermissions.PlanIds.IsNull() {
			var planIDs []int64
			d := plan.ResourcePermissions.PlanIds.ElementsAs(ctx, &planIDs, false)
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
			ResourcePermissions:  &rp,
			AdditionalProperties: map[string]interface{}{},
		}
	}

	if !plan.TenantIds.IsNull() && !plan.TenantIds.IsUnknown() {
		var tenantIDs []int64
		d := plan.TenantIds.ElementsAs(ctx, &tenantIDs, false)
		resp.Diagnostics.Append(d...)
		if resp.Diagnostics.HasError() {
			return
		}
		if ns.Permissions == nil {
			ns.Permissions = &sdk.UpdateClusterNamespaceRequestNamespacePermissions{
				AdditionalProperties: map[string]interface{}{},
			}
		}
		ns.Permissions.AdditionalProperties["tenantPermissions"] = map[string]interface{}{
			"accounts": tenantIDs,
		}
	}
	if !plan.Description.IsNull() {
		v := plan.Description.ValueString()
		ns.Description = &v
	}
	if !plan.Active.IsNull() {
		ns.Active = plan.Active.ValueBoolPointer()
	}

	if !plan.ResourcePermissions.IsNull() && !plan.ResourcePermissions.IsUnknown() {
		rp := sdk.UpdateClusterNamespaceRequestNamespacePermissionsResourcePermissions{}
		if !plan.ResourcePermissions.All.IsNull() {
			rp.All = plan.ResourcePermissions.All.ValueBoolPointer()
		}
		if !plan.ResourcePermissions.AllPlans.IsNull() {
			rp.AllPlans = plan.ResourcePermissions.AllPlans.ValueBoolPointer()
		}
		if !plan.ResourcePermissions.SiteIds.IsNull() {
			var siteIDs []int64
			d := plan.ResourcePermissions.SiteIds.ElementsAs(ctx, &siteIDs, false)
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
		if !plan.ResourcePermissions.PlanIds.IsNull() {
			var planIDs []int64
			d := plan.ResourcePermissions.PlanIds.ElementsAs(ctx, &planIDs, false)
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
	mapGetResponseToModel(ctx, &plan, readNs, &resp.Diagnostics)

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

	var state ClusterNamespaceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	clusterID := state.ClusterId.ValueInt64()
	id := state.Id.ValueInt64()

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

func mapGetResponseToModel(
	ctx context.Context,
	model *ClusterNamespaceModel,
	ns *sdk.GetClusterNamespace200ResponseNamespace,
	diags *diag.Diagnostics,
) {
	if ns.Id != nil {
		model.Id = types.Int64Value(*ns.Id)
	}
	if ns.Name != nil {
		model.Name = types.StringValue(*ns.Name)
	}
	if ns.Description != nil {
		model.Description = types.StringValue(*ns.Description)
	}
	// NOTE: Active is not in the API GET at all. Config value is preserved in state.

	// visibility
	if ns.Visibility != nil {
		model.Visibility = types.StringValue(*ns.Visibility)
	} else {
		model.Visibility = types.StringValue("private")
	}

	// tenant_ids: returned in permissions.tenantPermissions.accounts
	model.TenantIds = types.SetNull(types.Int64Type)
	if ns.Permissions != nil {
		if tp, ok := ns.Permissions.AdditionalProperties["tenantPermissions"]; ok {
			if tpMap, ok := tp.(map[string]interface{}); ok {
				if accounts, ok := tpMap["accounts"]; ok {
					if accountList, ok := accounts.([]interface{}); ok {
						vals := make([]attr.Value, 0, len(accountList))
						for _, a := range accountList {
							switch v := a.(type) {
							case float64:
								vals = append(vals, types.Int64Value(int64(v)))
							case int64:
								vals = append(vals, types.Int64Value(v))
							}
						}
						if len(vals) > 0 {
							tenantSet, d := types.SetValue(types.Int64Type, vals)
							diags.Append(d...)
							model.TenantIds = tenantSet
						}
					}
				}
			}
		}
	}

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

		permVal, d := NewResourcePermissionsValue(ResourcePermissionsValue{}.AttributeTypes(ctx), map[string]attr.Value{
			"all":       allVal,
			"all_plans": allPlansVal,
			"site_ids":  siteIDSet,
			"plan_ids":  planIDSet,
		})
		diags.Append(d...)
		model.ResourcePermissions = permVal
	} else {
		model.ResourcePermissions = NewResourcePermissionsValueNull()
	}
}
