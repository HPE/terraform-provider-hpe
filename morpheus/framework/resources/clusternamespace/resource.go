package clusternamespace

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
	resp.TypeName = req.ProviderTypeName + "_" + "cluster_namespace"
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
	ns.Visibility = plan.Visibility.ValueStringPointer()
	if perms := buildNamespaceCreatePermissions(ctx, plan, &resp.Diagnostics); perms != nil {
		ns.Permissions = perms
	}
	if resp.Diagnostics.HasError() {
		return
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
	mapGetResponseToModel(&plan, readNs)

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

	mapGetResponseToModel(&state, ns)

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
	ns.Visibility = plan.Visibility.ValueStringPointer()
	if perms := buildNamespaceUpdatePermissions(ctx, plan, &resp.Diagnostics); perms != nil {
		ns.Permissions = perms
	}
	if resp.Diagnostics.HasError() {
		return
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
	mapGetResponseToModel(&plan, readNs)

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

func mapGetResponseToModel(model *ClusterNamespaceModel, ns *sdk.GetClusterNamespace200ResponseNamespace) {
	if ns.Id != nil {
		model.Id = types.Int64Value(*ns.Id)
	}
	if ns.Name != nil {
		model.Name = types.StringValue(*ns.Name)
	}
	if ns.Description != nil {
		model.Description = types.StringValue(*ns.Description)
	}
	if ns.Visibility != nil {
		model.Visibility = types.StringValue(*ns.Visibility)
	} else {
		model.Visibility = types.StringNull()
	}
	// Map tenant_ids from permissions.tenantPermissions.accounts
	if ns.Permissions != nil && ns.Permissions.TenantPermissions != nil {
		accts := ns.Permissions.TenantPermissions.Accounts
		tenantVals := make([]attr.Value, 0, len(accts))
		for _, a := range accts {
			if a.Id != nil {
				tenantVals = append(tenantVals, types.Int64Value(*a.Id))
			}
		}
		model.TenantIds = types.SetValueMust(types.Int64Type, tenantVals)
	} else {
		model.TenantIds = types.SetValueMust(types.Int64Type, []attr.Value{})
	}
	// Map resource_permissions from permissions.resourcePermissions
	if ns.Permissions != nil && ns.Permissions.ResourcePermissions != nil {
		rp := ns.Permissions.ResourcePermissions
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
		model.ResourcePermissions = NewResourcePermissionsValueMust(
			map[string]attr.Type{
				"all":       types.BoolType,
				"site_ids":  types.ListType{ElemType: types.Int64Type},
				"all_plans": types.BoolType,
				"plan_ids":  types.ListType{ElemType: types.Int64Type},
			},
			map[string]attr.Value{
				"all":       types.BoolPointerValue(rp.All),
				"site_ids":  types.ListValueMust(types.Int64Type, siteVals),
				"all_plans": types.BoolPointerValue(rp.AllPlans),
				"plan_ids":  types.ListValueMust(types.Int64Type, planVals),
			},
		)
	} else {
		model.ResourcePermissions = NewResourcePermissionsValueNull()
	}
	// NOTE: Active is not in the API GET at all. Config value is preserved in state.
}

func buildNamespaceCreatePermissions(
	ctx context.Context,
	plan ClusterNamespaceModel,
	diags *diag.Diagnostics,
) *sdk.AddClusterNamespaceRequestNamespacePermissions {
	hasTP := !plan.TenantIds.IsNull() && !plan.TenantIds.IsUnknown()
	hasRP := !plan.ResourcePermissions.IsNull() && !plan.ResourcePermissions.IsUnknown()
	if !hasTP && !hasRP {
		return nil
	}
	perms := &sdk.AddClusterNamespaceRequestNamespacePermissions{}
	if hasTP {
		var ids []int64
		diags.Append(plan.TenantIds.ElementsAs(ctx, &ids, false)...)
		if diags.HasError() {
			return nil
		}
		accts := make([]sdk.AddClusterNamespaceRequestNamespacePermissionsTenantPermissionsAccountsInner, 0, len(ids))
		for i := range ids {
			id := ids[i]
			accts = append(accts, sdk.AddClusterNamespaceRequestNamespacePermissionsTenantPermissionsAccountsInner{Id: &id})
		}
		perms.TenantPermissions = &sdk.AddClusterNamespaceRequestNamespacePermissionsTenantPermissions{Accounts: accts}
	}
	if hasRP {
		rp := sdk.AddClusterNamespaceRequestNamespacePermissionsResourcePermissions{}
		rp.All = plan.ResourcePermissions.All.ValueBoolPointer()
		rp.AllPlans = plan.ResourcePermissions.AllPlans.ValueBoolPointer()
		if !plan.ResourcePermissions.SiteIds.IsNull() && !plan.ResourcePermissions.SiteIds.IsUnknown() {
			var siteIDs []int64
			diags.Append(plan.ResourcePermissions.SiteIds.ElementsAs(ctx, &siteIDs, false)...)
			if diags.HasError() {
				return nil
			}
			sites := make([]sdk.AddClusterNamespaceRequestNamespacePermissionsResourcePermissionsSitesInner, 0, len(siteIDs))
			for i := range siteIDs {
				id := siteIDs[i]
				sites = append(sites, sdk.AddClusterNamespaceRequestNamespacePermissionsResourcePermissionsSitesInner{Id: &id})
			}
			rp.Sites = sites
		}
		if !plan.ResourcePermissions.PlanIds.IsNull() && !plan.ResourcePermissions.PlanIds.IsUnknown() {
			var planIDs []int64
			diags.Append(plan.ResourcePermissions.PlanIds.ElementsAs(ctx, &planIDs, false)...)
			if diags.HasError() {
				return nil
			}
			plans := make([]sdk.AddClusterNamespaceRequestNamespacePermissionsResourcePermissionsPlansInner, 0, len(planIDs))
			for i := range planIDs {
				id := planIDs[i]
				plans = append(plans, sdk.AddClusterNamespaceRequestNamespacePermissionsResourcePermissionsPlansInner{Id: &id})
			}
			rp.Plans = plans
		}
		perms.ResourcePermissions = &rp
	}

	return perms
}

func buildNamespaceUpdatePermissions(
	ctx context.Context,
	plan ClusterNamespaceModel,
	diags *diag.Diagnostics,
) *sdk.UpdateClusterNamespaceRequestNamespacePermissions {
	hasTP := !plan.TenantIds.IsNull() && !plan.TenantIds.IsUnknown()
	hasRP := !plan.ResourcePermissions.IsNull() && !plan.ResourcePermissions.IsUnknown()
	if !hasTP && !hasRP {
		return nil
	}
	perms := &sdk.UpdateClusterNamespaceRequestNamespacePermissions{}
	if hasTP {
		var ids []int64
		diags.Append(plan.TenantIds.ElementsAs(ctx, &ids, false)...)
		if diags.HasError() {
			return nil
		}
		accts := make([]sdk.UpdateClusterNamespaceRequestNamespacePermissionsTenantPermissionsAccountsInner, 0, len(ids))
		for i := range ids {
			id := ids[i]
			accts = append(accts, sdk.UpdateClusterNamespaceRequestNamespacePermissionsTenantPermissionsAccountsInner{Id: &id})
		}
		perms.TenantPermissions = &sdk.UpdateClusterNamespaceRequestNamespacePermissionsTenantPermissions{Accounts: accts}
	}
	if hasRP {
		rp := sdk.UpdateClusterNamespaceRequestNamespacePermissionsResourcePermissions{}
		rp.All = plan.ResourcePermissions.All.ValueBoolPointer()
		rp.AllPlans = plan.ResourcePermissions.AllPlans.ValueBoolPointer()
		if !plan.ResourcePermissions.SiteIds.IsNull() && !plan.ResourcePermissions.SiteIds.IsUnknown() {
			var siteIDs []int64
			diags.Append(plan.ResourcePermissions.SiteIds.ElementsAs(ctx, &siteIDs, false)...)
			if diags.HasError() {
				return nil
			}
			sites := make([]sdk.UpdateClusterNamespaceRequestNamespacePermissionsResourcePermissionsSitesInner, 0, len(siteIDs))
			for i := range siteIDs {
				id := siteIDs[i]
				sites = append(sites, sdk.UpdateClusterNamespaceRequestNamespacePermissionsResourcePermissionsSitesInner{Id: &id})
			}
			rp.Sites = sites
		}
		if !plan.ResourcePermissions.PlanIds.IsNull() && !plan.ResourcePermissions.PlanIds.IsUnknown() {
			var planIDs []int64
			diags.Append(plan.ResourcePermissions.PlanIds.ElementsAs(ctx, &planIDs, false)...)
			if diags.HasError() {
				return nil
			}
			plans := make([]sdk.UpdateClusterNamespaceRequestNamespacePermissionsResourcePermissionsPlansInner, 0, len(planIDs))
			for i := range planIDs {
				id := planIDs[i]
				plans = append(plans, sdk.UpdateClusterNamespaceRequestNamespacePermissionsResourcePermissionsPlansInner{Id: &id})
			}
			rp.Plans = plans
		}
		perms.ResourcePermissions = &rp
	}

	return perms
}
