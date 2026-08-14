// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package cloudaffinitygroupmember

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/affinitylock"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/affinityread"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
)

// Create adds one server to the affinity group.
//
// There is no member-level API, so this reads the group, appends the server and
// writes the whole list back. The lock keeps concurrent memberships in the same
// apply from overwriting one another; the read-back afterwards catches the
// cases the lock cannot cover.
func (r *Resource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	var plan memberModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError("could not create sdk client", err.Error())

		return
	}

	cloudID := plan.CloudID.ValueInt64()
	groupID := plan.AffinityGroupID.ValueInt64()
	serverID := plan.ServerID.ValueInt64()

	defer affinitylock.Acquire(affinitylock.Cloud, cloudID, groupID)()

	servers, ok := readMembership(ctx, client, cloudID, groupID, &resp.Diagnostics)
	if !ok {
		return
	}

	// Already a member. Nothing to send, and sending anyway would rewrite the
	// whole list for no reason.
	if !containsServer(servers, serverID) {
		if !writeMembership(
			ctx, client, cloudID, groupID, append(servers, serverID), &resp.Diagnostics,
		) {
			return
		}

		// Read-back assertion.
		//
		// The update returning success is not proof the membership took: the
		// API replaces the whole list, so a concurrent writer that read before
		// this write can silently drop the server again. Confirm rather than
		// assume.
		after, ok := readMembership(ctx, client, cloudID, groupID, &resp.Diagnostics)
		if !ok {
			return
		}

		if !containsServer(after, serverID) {
			resp.Diagnostics.AddError(
				"affinity group membership was not applied",
				"Server "+plan.ServerID.String()+" is not a member of affinity group "+
					plan.AffinityGroupID.String()+" after the update was accepted. "+
					"The group's membership is replaced as a whole on every write, so "+
					"another process managing the same group concurrently can overwrite "+
					"this change. Re-apply, and avoid managing this group's servers "+
					"attribute at the same time as this resource.",
			)

			return
		}
	}

	plan.ID = types.StringValue(membershipID(cloudID, groupID, serverID))

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// readMembership returns the group's current server list.
//
// TODO(MORPH-15806): drop the list fallback once the appliance defect is fixed.
// An affinity group that has been sent tenant permissions returns HTTP 500 from
// its single-item endpoint from then on, permanently, even though the record is
// undamaged. Without the fallback one such group makes every plan fail for every
// membership that references it.
//
// The fallback is lossless here: the list omits only resourcePermissions and
// tenants, and this function needs neither.
func readMembership(
	ctx context.Context,
	client *sdk.APIClient,
	cloudID, groupID int64,
	diags *diag.Diagnostics,
) ([]int64, bool) {
	result, httpResp, err := client.CloudsAPI.
		GetCloudAffinityGroup(ctx, cloudID, groupID).Execute()
	if err == nil && result != nil && result.AffinityGroup != nil {
		servers := make([]int64, 0, len(result.AffinityGroup.Servers))
		for _, s := range result.AffinityGroup.Servers {
			if s.Id != nil {
				servers = append(servers, *s.Id)
			}
		}

		return servers, true
	}

	// The group is gone. Report absence without a diagnostic so the caller can
	// drop the membership from state rather than failing the plan: a group
	// deleted out of band should surface as drift, not as an error.
	if errfmt.IsNotFound(httpResp) {
		return nil, false
	}

	if !affinityread.IsSingleItemRenderFailure(httpResp) {
		if err := errfmt.CheckResponse(err, httpResp); err != nil {
			errfmt.DiagError(
				diags, errfmt.OpRead, "cloud_affinity_group_member", "", err, httpResp,
			)

			return nil, false
		}

		diags.AddError("API returned nil", "AffinityGroup is nil in the response")

		return nil, false
	}

	listResult, listResp, listErr := client.CloudsAPI.
		ListCloudAffinityGroups(ctx, cloudID).Execute()
	if err := errfmt.CheckResponse(listErr, listResp); err != nil {
		errfmt.DiagError(
			diags, errfmt.OpRead, "cloud_affinity_group_member", "", err, listResp,
		)

		return nil, false
	}

	if listResult == nil {
		diags.AddError("API returned nil", "affinity group list is nil in the response")

		return nil, false
	}

	servers, found := affinityread.ServersFromList(
		listResult.AffinityGroups,
		groupID,
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

// writeMembership replaces the group's server list.
func writeMembership(
	ctx context.Context,
	client *sdk.APIClient,
	cloudID, groupID int64,
	servers []int64,
	diags *diag.Diagnostics,
) bool {
	body := sdk.UpdateCloudAffinityGroupRequest{
		AffinityGroup: &sdk.UpdateCloudAffinityGroupRequestAffinityGroup{
			Servers: servers,
		},
	}

	_, httpResp, err := client.CloudsAPI.
		UpdateCloudAffinityGroup(ctx, cloudID, groupID).
		UpdateCloudAffinityGroupRequest(body).Execute()

	// TODO(MORPH-15806): stop treating 500 as inconclusive once the appliance
	// defect is fixed.
	//
	// On a group that has been sent tenant permissions, the update applies but
	// rendering the response fails, so a 500 comes back for a write that
	// succeeded. Failing here would make such a group unmanageable. The caller
	// verifies the outcome by reading membership back, which is the authority
	// either way, so an inconclusive result is passed through rather than
	// reported as a failure.
	if affinityread.IsSingleItemRenderFailure(httpResp) {
		return true
	}

	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(
			diags, errfmt.OpUpdate, "cloud_affinity_group_member", "", err, httpResp,
		)

		return false
	}

	return true
}
