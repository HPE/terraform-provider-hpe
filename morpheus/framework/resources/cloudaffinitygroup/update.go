// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package cloudaffinitygroup

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/constants"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/versioncheck"
)

func (r *cloudAffinityGroupResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	client, err := r.NewClient(ctx)
	if err != nil {
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	// MORPH-15506 appliance version gate — see Create.
	resp.Diagnostics.Append(versioncheck.Require(
		ctx, client, gatedFeature, constants.AffinityGroupMinVersion,
	)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var plan, state CloudAffinityGroupModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	cloudID := plan.CloudId.ValueInt64()
	id := plan.Id.ValueInt64()

	ag := sdk.UpdateCloudAffinityGroupRequestAffinityGroup{
		Name:   plan.Name.ValueStringPointer(),
		Active: plan.Active.ValueBoolPointer(),
	}

	if !plan.Visibility.IsNull() && !plan.Visibility.IsUnknown() {
		ag.Visibility = plan.Visibility.ValueStringPointer()
	}

	// CRITICAL BEHAVIOUR 2: servers is WHOLESALE REPLACE on update.
	// resolveUpdateServers falls back to the membership recorded in STATE when the
	// planned value is unknown or null. Read its doc comment before changing this.
	servers, serverDiags := resolveUpdateServers(ctx, plan.Servers, state.Servers)
	resp.Diagnostics.Append(serverDiags...)
	if resp.Diagnostics.HasError() {
		return
	}
	ag.Servers = servers

	// Tenants.
	if !plan.TenantIds.IsNull() && !plan.TenantIds.IsUnknown() {
		var tenantIDs []int64
		resp.Diagnostics.Append(plan.TenantIds.ElementsAs(ctx, &tenantIDs, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		tenants := make(
			[]sdk.UpdateCloudAffinityGroupRequestAffinityGroupTenantsInner, 0, len(tenantIDs),
		)
		for _, tid := range tenantIDs {
			tid := tid
			tenants = append(tenants, sdk.UpdateCloudAffinityGroupRequestAffinityGroupTenantsInner{
				Id: &tid,
			})
		}
		ag.Tenants = tenants
	}

	// ResourcePermissions.
	if !plan.ResourcePermissions.IsNull() && !plan.ResourcePermissions.IsUnknown() {
		rp := sdk.UpdateCloudAffinityGroupRequestAffinityGroupResourcePermissions{
			All: plan.ResourcePermissions.All.ValueBoolPointer(),
		}
		if !plan.ResourcePermissions.Groups.IsNull() &&
			!plan.ResourcePermissions.Groups.IsUnknown() {
			var groups []GroupsValue
			resp.Diagnostics.Append(
				plan.ResourcePermissions.Groups.ElementsAs(ctx, &groups, false)...,
			)
			if resp.Diagnostics.HasError() {
				return
			}
			rp.Sites = buildSitesPayload(groups)
		}
		ag.ResourcePermissions = &rp
	}

	body := sdk.UpdateCloudAffinityGroupRequest{
		AffinityGroup: &ag,
	}

	_, httpResp, err := client.CloudsAPI.UpdateCloudAffinityGroup(ctx, cloudID, id).
		UpdateCloudAffinityGroupRequest(body).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(
			&resp.Diagnostics, errfmt.OpUpdate, "cloud_affinity_group",
			plan.Name.ValueString(), err, httpResp,
		)

		return
	}

	// Read-back.
	readResult, httpResp, err := client.CloudsAPI.GetCloudAffinityGroup(ctx, cloudID, id).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(
			&resp.Diagnostics, errfmt.OpRead, "cloud_affinity_group",
			plan.Name.ValueString(), err, httpResp,
		)

		return
	}

	readAg := readResult.AffinityGroup
	if readAg == nil {
		resp.Diagnostics.AddError("API returned nil", "AffinityGroup is nil in the response")

		return
	}

	// CRITICAL BEHAVIOUR 5: Refuse to manage sync-discovered groups (on update read-back).
	if readAg.Source != nil && *readAg.Source == "sync" {
		resp.Diagnostics.AddError(
			"Cannot manage sync-discovered affinity group",
			fmt.Sprintf(
				"Affinity group %d (cloud %d) has source=\"sync\" and is managed by "+
					"cloud sync. Use the hpe_morpheus_cloud_affinity_group data source instead.",
				id, cloudID,
			),
		)

		return
	}

	resp.Diagnostics.Append(mapAndResolveResponse(ctx, &plan, readAg, cloudID)...)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// resolveUpdateServers decides what to send for affinityGroup.servers on update.
//
// WHY THIS EXISTS — do not collapse the unknown/null branch back into an empty array.
//
// The Morpheus API treats servers on update as a WHOLESALE REPLACE.
// AffinityGroupService.updateAffinityGroup diffs the supplied collection against the
// current members and evicts everything absent from it; its else-branch removes every
// member. Because Groovy truthiness makes both a missing key and an empty list falsy,
// OMITTING servers and SENDING [] are identical to the API: both wipe the group.
//
// servers is Optional+Computed in the generated schema (schema_gen.go) and carries no
// UseStateForUnknown plan modifier. So whenever a practitioner has not configured
// servers and changes any other attribute — a rename, for instance — Terraform marks
// servers UNKNOWN in the plan. Treating that unknown as "send []" silently destroyed
// the group's entire membership, with no error and no warning.
//
// Resolution:
//
//	plan known, non-empty  -> the planned set   (the intended wholesale replace)
//	plan known, EMPTY      -> []                (practitioner wrote `servers = []`;
//	                                            a genuine, explicit "remove all")
//	plan UNKNOWN or NULL   -> the set in STATE  (re-asserts current membership, so the
//	                                            wholesale replace is a no-op and the
//	                                            members survive)
//
// The known-empty case must stay distinguishable from null/unknown: that distinction is
// the entire point of this function.
//
// If neither plan nor state carries a known set, the current membership is unknowable.
// The API has no encoding for "leave unchanged", so nil is returned and the key is
// omitted rather than asserting a membership we cannot substantiate.
func resolveUpdateServers(
	ctx context.Context,
	planServers types.Set,
	stateServers types.Set,
) ([]int64, diag.Diagnostics) {
	var diags diag.Diagnostics

	source := planServers
	if planServers.IsNull() || planServers.IsUnknown() {
		source = stateServers
	}

	if source.IsNull() || source.IsUnknown() {
		return nil, diags
	}

	var serverIDs []int64
	diags.Append(source.ElementsAs(ctx, &serverIDs, false)...)
	if diags.HasError() {
		return nil, diags
	}

	// ElementsAs always allocates for a known set, so a known-empty set yields []int64{}
	// rather than a nil slice. The SDK's ToMap only omits the key when the slice is nil,
	// so this is what makes an explicit `servers = []` serialise as "servers": [].
	return serverIDs, diags
}

// buildSitesPayload converts GroupsValue slice to []map[string]interface{} the SDK expects.
func buildSitesPayload(groups []GroupsValue) []map[string]interface{} {
	sites := make([]map[string]interface{}, 0, len(groups))
	for _, g := range groups {
		site := map[string]interface{}{
			"id": g.Id.ValueInt64(),
		}
		if !g.Default.IsNull() && !g.Default.IsUnknown() {
			site["default"] = g.Default.ValueBool()
		}
		sites = append(sites, site)
	}

	return sites
}
