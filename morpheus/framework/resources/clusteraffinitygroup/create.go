// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package clusteraffinitygroup

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/affinityread"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/constants"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/versioncheck"
	"github.com/HPE/terraform-provider-hpe/utils/cleanup"
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

	// MORPH-15806: a create carrying tenant permissions is stored but cannot be
	// rendered, so it answers 500 and the response arrives without the id.
	// Record what exists beforehand so the new group can be found by difference.
	//
	// Only done when tenants are actually sent: nothing else provokes the
	// defect, and the extra listing is not free.
	var priorIDs map[int64]struct{}
	if body.TenantPermissions != nil {
		priorIDs = listAffinityGroupIDs(ctx, client, clusterID)
	}

	result, httpResp, err := client.ClustersAPI.SaveClusterAffinityGroup(ctx, clusterID).
		SaveClusterAffinityGroupRequest(body).Execute()

	// The group is created even when the response cannot be rendered. Failing
	// here would abandon it: never in state, invisible to Terraform, and
	// removable only by hand.
	renderFailed := affinityread.IsSingleItemRenderFailure(httpResp)

	if err := errfmt.CheckResponse(err, httpResp); err != nil && !renderFailed {
		errfmt.DiagError(
			&resp.Diagnostics, errfmt.OpCreate, "cluster_affinity_group",
			plan.Name.ValueString(), err, httpResp,
		)

		return
	}

	id, ok := createdGroupID(
		ctx, client, clusterID, result, priorIDs, renderFailed, &resp.Diagnostics,
	)
	if !ok {
		return
	}

	// Read-back to populate full state.
	readAg, found := fetchAffinityGroup(ctx, client, clusterID, id, &resp.Diagnostics)
	if !found {
		if !resp.Diagnostics.HasError() {
			resp.Diagnostics.AddError(
				"affinity group not found after create",
				"The affinity group was created but could not be read back.",
			)
		}

		// The group exists in Morpheus but never reached state. Taint so the
		// next apply replaces it rather than leaving it orphaned and invisible
		// to Terraform.
		cleanup.TaintResourceState(ctx, cleanup.TaintResourceStateConfig{
			ResourceType: "cluster_affinity_group",
			ResourceID:   id,
			StateWriter:  &resp.State,
			Diagnostics:  &resp.Diagnostics,
		})

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

// listAffinityGroupIDs returns the ids of every affinity group on the cluster.
//
// TODO(MORPH-15806): remove alongside the rest of the workaround.
//
// A nil result means the listing could not be read. Errors are deliberately
// swallowed: this runs before the create as a precaution, and a failure to take
// the precaution must not stop a create that would otherwise succeed. If the
// precaution turns out to have been needed, createdGroupID reports it then,
// where the consequence is real rather than hypothetical.
func listAffinityGroupIDs(
	ctx context.Context,
	client *sdk.APIClient,
	clusterID int64,
) map[int64]struct{} {
	result, httpResp, err := client.ClustersAPI.ListClusterAffinityGroups(ctx, clusterID).Execute()
	if errfmt.CheckResponse(err, httpResp) != nil || result == nil {
		return nil
	}

	return affinityread.IDsFromList(
		result.AffinityGroups,
		func(e sdk.ListClusterAffinityGroups200ResponseAllOfAffinityGroupsInner) (int64, bool) {
			if e.Id == nil {
				return 0, false
			}

			return *e.Id, true
		},
	)
}

// createdGroupID determines the id of the group that was just created.
//
// TODO(MORPH-15806): remove alongside the rest of the workaround.
//
// Normally the create response carries it. When the response could not be
// rendered it does not, and the id is recovered by listing the groups and
// taking the one that was not there before.
func createdGroupID(
	ctx context.Context,
	client *sdk.APIClient,
	clusterID int64,
	result *sdk.SaveClusterAffinityGroup200Response,
	priorIDs map[int64]struct{},
	renderFailed bool,
	diags *diag.Diagnostics,
) (int64, bool) {
	if result != nil && result.AffinityGroup != nil && result.AffinityGroup.Id != nil {
		return *result.AffinityGroup.Id, true
	}

	if !renderFailed {
		diags.AddError(
			"API returned nil ID", "AffinityGroup ID is nil in the create response",
		)

		return 0, false
	}

	if priorIDs == nil {
		diags.AddError(
			"affinity group created but could not be identified",
			"The affinity group was created, but the API could not render the "+
				"response and the existing groups could not be listed beforehand, "+
				"so its ID is unknown. The group exists on the appliance and is "+
				"not managed by Terraform. Import it or remove it by hand.",
		)

		return 0, false
	}

	listResult, listResp, listErr := client.ClustersAPI.
		ListClusterAffinityGroups(ctx, clusterID).Execute()
	if err := errfmt.CheckResponse(listErr, listResp); err != nil {
		errfmt.DiagError(diags, errfmt.OpCreate, "cluster_affinity_group", "", err, listResp)

		return 0, false
	}

	if listResult == nil {
		diags.AddError("API returned nil", "affinity group list is nil in the response")

		return 0, false
	}

	id, ok := affinityread.NewIDFromList(
		listResult.AffinityGroups, priorIDs,
		func(e sdk.ListClusterAffinityGroups200ResponseAllOfAffinityGroupsInner) (int64, bool) {
			if e.Id == nil {
				return 0, false
			}

			return *e.Id, true
		},
	)
	if !ok {
		diags.AddError(
			"affinity group created but could not be identified",
			"The affinity group was created, but the API could not render the "+
				"response and the group could not be singled out from the listing "+
				"afterwards. This happens when another affinity group is created "+
				"on the same cluster at the same moment. The group exists on the "+
				"appliance and is not managed by Terraform. Import it or remove it "+
				"by hand.",
		)

		return 0, false
	}

	return id, true
}
