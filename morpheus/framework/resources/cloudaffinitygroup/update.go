// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package cloudaffinitygroup

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/affinitylock"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/affinityread"
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
	//
	// The API has no encoding for "leave membership alone": omitting servers and
	// sending an empty array both empty the group. Every update therefore has to
	// re-assert the current membership, even when the practitioner is only
	// renaming the group.
	//
	// That membership is read live rather than taken from Terraform state.
	// State records what this resource last saw, which is not the same thing: a
	// server added since then -- by an instance provisioned into the group, or
	// by a membership resource -- is absent from state, and echoing state back
	// would evict it. The read and the write happen under the group's lock so
	// they cannot interleave with a membership resource doing the same thing.
	defer affinitylock.Acquire(affinitylock.Cloud, cloudID, id)()

	currentServers, ok := readCurrentServers(ctx, client, cloudID, id, &resp.Diagnostics)
	if !ok {
		return
	}

	ag.Servers = currentServers

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

	// TODO(MORPH-15806): stop treating 500 as inconclusive once the appliance
	// defect is fixed.
	//
	// On a group that has been sent tenant permissions the update applies but
	// rendering the response fails, so a 500 comes back for a write that
	// succeeded. Failing here would make the group unmanageable. The read-back
	// below establishes the real outcome either way.
	renderFailed := affinityread.IsSingleItemRenderFailure(httpResp)

	if err := errfmt.CheckResponse(err, httpResp); err != nil && !renderFailed {
		errfmt.DiagError(
			&resp.Diagnostics, errfmt.OpUpdate, "cloud_affinity_group",
			plan.Name.ValueString(), err, httpResp,
		)

		return
	}

	// Read-back.
	readAg, ok := fetchAffinityGroup(ctx, client, cloudID, id, &resp.Diagnostics)
	if !ok {
		if !resp.Diagnostics.HasError() {
			resp.Diagnostics.AddError(
				"affinity group not found after update",
				"The affinity group was updated but could not be read back.",
			)
		}

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

// buildSitesPayload converts GroupsValue slice to []map[string]interface{} the SDK expects.
func buildSitesPayload(groups []GroupsValue) []map[string]interface{} {
	sites := make([]map[string]interface{}, 0, len(groups))
	for _, g := range groups {
		// Skip an entry with no usable id rather than sending 0, which the API
		// would take as a real group. id is required by the schema, so this
		// should not arise from configuration, but it can still be unknown
		// while a value it depends on is being computed.
		if g.Id.IsNull() || g.Id.IsUnknown() {
			continue
		}

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

// readCurrentServers returns the group's membership as it is right now.
//
// Used to re-assert membership on update. See CRITICAL BEHAVIOUR 2 in Update for
// why this cannot come from Terraform state.
func readCurrentServers(
	ctx context.Context,
	client *sdk.APIClient,
	cloudID, id int64,
	diags *diag.Diagnostics,
) ([]int64, bool) {
	result, httpResp, err := client.CloudsAPI.GetCloudAffinityGroup(ctx, cloudID, id).Execute()
	if err == nil && result != nil && result.AffinityGroup != nil {
		servers := make([]int64, 0, len(result.AffinityGroup.Servers))
		for _, s := range result.AffinityGroup.Servers {
			if s.Id != nil {
				servers = append(servers, *s.Id)
			}
		}

		return servers, true
	}

	// TODO(MORPH-15806): drop this fallback once the appliance defect is fixed.
	//
	// A group that has been sent tenant permissions returns 500 from its
	// single-item endpoint permanently, so without this an unrelated rename
	// would fail. The listing renders the same group correctly, and membership
	// is all this function needs, so the fallback loses nothing.
	if !affinityread.IsSingleItemRenderFailure(httpResp) {
		if err := errfmt.CheckResponse(err, httpResp); err != nil {
			errfmt.DiagError(diags, errfmt.OpRead, "cloud_affinity_group", "", err, httpResp)

			return nil, false
		}

		diags.AddError("API returned nil", "AffinityGroup is nil in the response")

		return nil, false
	}

	listResult, listResp, listErr := client.CloudsAPI.ListCloudAffinityGroups(ctx, cloudID).Execute()
	if err := errfmt.CheckResponse(listErr, listResp); err != nil {
		errfmt.DiagError(diags, errfmt.OpRead, "cloud_affinity_group", "", err, listResp)

		return nil, false
	}

	if listResult == nil {
		diags.AddError("API returned nil", "affinity group list is nil in the response")

		return nil, false
	}

	servers, found := affinityread.ServersFromList(
		listResult.AffinityGroups,
		id,
		func(a sdk.ListCloudAffinityGroups200ResponseAllOfAffinityGroupsInner) (int64, bool) {
			if a.Id == nil {
				return 0, false
			}

			return *a.Id, true
		},
		func(a sdk.ListCloudAffinityGroups200ResponseAllOfAffinityGroupsInner) []int64 {
			out := make([]int64, 0, len(a.Servers))
			for _, s := range a.Servers {
				if s.Id != nil {
					out = append(out, *s.Id)
				}
			}

			return out
		},
	)

	if !found {
		diags.AddError(
			"affinity group not found",
			"The affinity group could not be read from its single-item endpoint, "+
				"and does not appear in the group listing either.",
		)

		return nil, false
	}

	return servers, true
}
