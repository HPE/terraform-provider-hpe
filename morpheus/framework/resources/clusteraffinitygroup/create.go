// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package clusteraffinitygroup

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/constants"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/versioncheck"
)

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

	// MORPH-15506: refuse to operate against an appliance older than the first
	// release with stable affinity group semantics, so the practitioner gets a
	// diagnostic naming the required version instead of an opaque API error.
	//
	// The check sits at the top of each CRUD method rather than in Configure:
	// the framework calls Configure on every RPC for the type, including
	// ValidateResourceConfig and UpgradeResourceState, and neither should have
	// to reach the network. See versioncheck.Require for the full rationale,
	// including why an unreadable version fails open.
	resp.Diagnostics.Append(versioncheck.Require(
		ctx, client, gatedFeature, constants.AffinityGroupMinVersion,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var plan ClusterAffinityGroupModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	clusterID := plan.ClusterId.ValueInt64()

	ag := sdk.SaveClusterAffinityGroupRequestAffinityGroup{
		Name: plan.Name.ValueStringPointer(),
	}

	// CRITICAL BEHAVIOUR 1: active MUST ALWAYS be sent on create.
	// The API does `active = (params.active == 'on' || params.active == true)` unconditionally,
	// which overrides the domain default of true. If omitted, the group is created INACTIVE.
	if plan.Active.IsNull() || plan.Active.IsUnknown() {
		active := true
		ag.Active = &active
	} else {
		ag.Active = plan.Active.ValueBoolPointer()
	}

	// CRITICAL BEHAVIOUR 7: affinity_type is create-only (absent from update model).
	if !plan.AffinityType.IsNull() && !plan.AffinityType.IsUnknown() {
		ag.AffinityType = plan.AffinityType.ValueStringPointer()
	}

	if !plan.Visibility.IsNull() && !plan.Visibility.IsUnknown() {
		ag.Visibility = plan.Visibility.ValueStringPointer()
	}

	// CRITICAL BEHAVIOUR 6: pool is COMPUTED ONLY for clusters — the API force-assigns it.
	// Never send it on create or update; only read it back.

	// ResourcePermissions.
	if !plan.ResourcePermissions.IsNull() && !plan.ResourcePermissions.IsUnknown() {
		rp := sdk.SaveClusterAffinityGroupRequestAffinityGroupResourcePermissions{
			All: plan.ResourcePermissions.All.ValueBoolPointer(),
		}
		if !plan.ResourcePermissions.Groups.IsNull() && !plan.ResourcePermissions.Groups.IsUnknown() {
			var groups []GroupsValue
			resp.Diagnostics.Append(plan.ResourcePermissions.Groups.ElementsAs(ctx, &groups, false)...)
			if resp.Diagnostics.HasError() {
				return
			}
			rp.Sites = buildSitesPayload(groups)
		}
		ag.ResourcePermissions = &rp
	}

	body := sdk.SaveClusterAffinityGroupRequest{
		AffinityGroup: &ag,
	}

	// Tenants are sent as the request-root `tenantPermissions` wrapper
	// (`{"accounts": [<id>, ...]}`), NOT as the nested `affinityGroup.tenants`
	// (`[{"id": <id>}]`) form. This is DELIBERATE — do not "harmonise" it with the
	// cloud affinity group resource, which legitimately uses the nested form.
	//
	// AffinityGroupService (morpheus-core) resolves tenant permissions as:
	//
	//	params.tenantPermissions ?: params.tenantPermission
	//	    ?: params.affinityGroup?.tenantPermissions ?: params.affinityGroup?.tenantPermission
	//	    ?: params.affinityGroup?.tenants
	//
	// so the request-root form is checked FIRST, at the highest precedence. Both forms
	// are accepted, but they have different shapes and are parsed by different branches
	// of permissionService.parseTenantPermissions. This resource has shipped sending the
	// request-root form, so switching would change the on-the-wire payload for existing
	// users with no benefit.
	if !plan.TenantIds.IsNull() && !plan.TenantIds.IsUnknown() {
		var tenantIDs []int64
		resp.Diagnostics.Append(plan.TenantIds.ElementsAs(ctx, &tenantIDs, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		body.TenantPermissions = &sdk.SaveClusterAffinityGroupRequestTenantPermissions{
			Accounts: tenantIDs,
		}
	}

	result, httpResp, err := client.ClustersAPI.SaveClusterAffinityGroup(ctx, clusterID).
		SaveClusterAffinityGroupRequest(body).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(
			&resp.Diagnostics, errfmt.OpCreate, "cluster_affinity_group",
			plan.Name.ValueString(), err, httpResp,
		)

		return
	}

	if result.AffinityGroup == nil || result.AffinityGroup.Id == nil {
		resp.Diagnostics.AddError(
			"API returned nil ID", "AffinityGroup ID is nil in the create response",
		)

		return
	}

	id := *result.AffinityGroup.Id

	// Read-back to populate full state.
	readAg, ok := fetchAffinityGroup(ctx, client, clusterID, id, &resp.Diagnostics)
	if !ok {
		if !resp.Diagnostics.HasError() {
			resp.Diagnostics.AddError(
				"affinity group not found after create",
				"The affinity group was created but could not be read back.",
			)
		}

		return
	}

	resp.Diagnostics.Append(mapAndResolveResponse(ctx, &plan, readAg, clusterID)...)

	// The create response is authoritative for the ID, so fall back to it if the
	// read-back somehow omitted it. id is Computed and therefore UNKNOWN in the
	// plan, and an unknown left in post-apply state is rejected outright.
	if plan.Id.IsUnknown() {
		plan.Id = types.Int64Value(id)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}
