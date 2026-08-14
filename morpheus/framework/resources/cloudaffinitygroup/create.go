// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package cloudaffinitygroup

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

func (r *cloudAffinityGroupResource) Create(
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

	var plan CloudAffinityGroupModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	cloudID := plan.CloudId.ValueInt64()

	ag := sdk.SaveCloudAffinityGroupRequestAffinityGroup{
		Name: plan.Name.ValueStringPointer(),
	}

	// CRITICAL BEHAVIOUR 1: active MUST ALWAYS be sent on create.
	// The API does `active = (params.active == 'on' || params.active == true)` unconditionally,
	// which overrides the domain default of true. If omitted, the group is created INACTIVE.
	// Resolve the planned value, defaulting to true when null/unknown.
	if plan.Active.IsNull() || plan.Active.IsUnknown() {
		active := true
		ag.Active = &active
	} else {
		ag.Active = plan.Active.ValueBoolPointer()
	}

	if !plan.AffinityType.IsNull() && !plan.AffinityType.IsUnknown() {
		ag.AffinityType = plan.AffinityType.ValueStringPointer()
	}

	if !plan.Visibility.IsNull() && !plan.Visibility.IsUnknown() {
		ag.Visibility = plan.Visibility.ValueStringPointer()
	}

	// Pool — create-only field. Wraps the scalar pool_id into the nested object the API expects.
	if !plan.PoolId.IsNull() && !plan.PoolId.IsUnknown() {
		poolID := plan.PoolId.ValueInt64()
		ag.Pool = &sdk.SaveCloudAffinityGroupRequestAffinityGroupPool{
			Id: &poolID,
		}
	}

	// Tenants — mapped inside the affinityGroup body.
	sendsTenants := !plan.TenantIds.IsNull() && !plan.TenantIds.IsUnknown()
	if sendsTenants {
		var tenantIDs []int64
		resp.Diagnostics.Append(plan.TenantIds.ElementsAs(ctx, &tenantIDs, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		tenants := make(
			[]sdk.SaveCloudAffinityGroupRequestAffinityGroupTenantsInner, 0, len(tenantIDs),
		)
		for _, tid := range tenantIDs {
			tid := tid
			tenants = append(tenants, sdk.SaveCloudAffinityGroupRequestAffinityGroupTenantsInner{
				Id: &tid,
			})
		}
		ag.Tenants = tenants
	}

	// ResourcePermissions.
	if !plan.ResourcePermissions.IsNull() && !plan.ResourcePermissions.IsUnknown() {
		rp := sdk.SaveCloudAffinityGroupRequestAffinityGroupResourcePermissions{
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

	body := sdk.SaveCloudAffinityGroupRequest{
		AffinityGroup: &ag,
	}

	// MORPH-15806: a create carrying tenant permissions is stored but cannot be
	// rendered, so it answers 500 and the response arrives without the id.
	// Record what exists beforehand so the new group can be found by difference.
	//
	// Only done when tenants are actually sent: nothing else provokes the
	// defect, and the extra listing is not free.
	var priorIDs map[int64]struct{}
	if sendsTenants {
		priorIDs = listAffinityGroupIDs(ctx, client, cloudID)
	}

	result, httpResp, err := client.CloudsAPI.SaveCloudAffinityGroup(ctx, cloudID).
		SaveCloudAffinityGroupRequest(body).Execute()

	// The group is created even when the response cannot be rendered. Failing
	// here would abandon it: never in state, invisible to Terraform, and
	// removable only by hand.
	renderFailed := affinityread.IsSingleItemRenderFailure(httpResp)

	if err := errfmt.CheckResponse(err, httpResp); err != nil && !renderFailed {
		errfmt.DiagError(
			&resp.Diagnostics, errfmt.OpCreate, "cloud_affinity_group",
			plan.Name.ValueString(), err, httpResp,
		)

		return
	}

	id, ok := createdGroupID(
		ctx, client, cloudID, result, priorIDs, renderFailed, &resp.Diagnostics,
	)
	if !ok {
		return
	}

	// Read-back to populate full state.
	readAg, found := fetchAffinityGroup(ctx, client, cloudID, id, &resp.Diagnostics)
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
			ResourceType: "cloud_affinity_group",
			ResourceID:   id,
			StateWriter:  &resp.State,
			Diagnostics:  &resp.Diagnostics,
		})

		return
	}

	resp.Diagnostics.Append(mapAndResolveResponse(ctx, &plan, readAg, cloudID)...)

	// The create response is authoritative for the ID, so fall back to it if the
	// read-back somehow omitted it. id is Computed and therefore UNKNOWN in the
	// plan, and an unknown left in post-apply state is rejected outright.
	if plan.Id.IsUnknown() {
		plan.Id = types.Int64Value(id)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// listAffinityGroupIDs returns the ids of every affinity group on the cloud.
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
	cloudID int64,
) map[int64]struct{} {
	result, httpResp, err := client.CloudsAPI.ListCloudAffinityGroups(ctx, cloudID).Execute()
	if errfmt.CheckResponse(err, httpResp) != nil || result == nil {
		return nil
	}

	return affinityread.IDsFromList(
		result.AffinityGroups,
		func(e sdk.ListCloudAffinityGroups200ResponseAllOfAffinityGroupsInner) (int64, bool) {
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
	cloudID int64,
	result *sdk.SaveCloudAffinityGroup200Response,
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

	listResult, listResp, listErr := client.CloudsAPI.
		ListCloudAffinityGroups(ctx, cloudID).Execute()
	if err := errfmt.CheckResponse(listErr, listResp); err != nil {
		errfmt.DiagError(diags, errfmt.OpCreate, "cloud_affinity_group", "", err, listResp)

		return 0, false
	}

	if listResult == nil {
		diags.AddError("API returned nil", "affinity group list is nil in the response")

		return 0, false
	}

	id, ok := affinityread.NewIDFromList(
		listResult.AffinityGroups, priorIDs,
		func(e sdk.ListCloudAffinityGroups200ResponseAllOfAffinityGroupsInner) (int64, bool) {
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
				"on the same cloud at the same moment. The group exists on the "+
				"appliance and is not managed by Terraform. Import it or remove it "+
				"by hand.",
		)

		return 0, false
	}

	return id, true
}
