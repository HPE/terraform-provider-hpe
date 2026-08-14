// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package cloudaffinitygroup

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"

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
	if !plan.TenantIds.IsNull() && !plan.TenantIds.IsUnknown() {
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

	result, httpResp, err := client.CloudsAPI.SaveCloudAffinityGroup(ctx, cloudID).
		SaveCloudAffinityGroupRequest(body).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(
			&resp.Diagnostics, errfmt.OpCreate, "cloud_affinity_group",
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
	readAg, ok := fetchAffinityGroup(ctx, client, cloudID, id, &resp.Diagnostics)
	if !ok {
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
