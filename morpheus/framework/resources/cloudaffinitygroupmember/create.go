// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package cloudaffinitygroupmember

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/affinitylock"
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
func readMembership(
	ctx context.Context,
	client *sdk.APIClient,
	cloudID, groupID int64,
	diags *diag.Diagnostics,
) ([]int64, bool) {
	result, httpResp, err := client.CloudsAPI.
		GetCloudAffinityGroup(ctx, cloudID, groupID).Execute()
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(
			diags, errfmt.OpRead, "cloud_affinity_group_member", "", err, httpResp,
		)

		return nil, false
	}

	if result == nil || result.AffinityGroup == nil {
		diags.AddError("API returned nil", "AffinityGroup is nil in the response")

		return nil, false
	}

	servers := make([]int64, 0, len(result.AffinityGroup.Servers))
	for _, s := range result.AffinityGroup.Servers {
		if s.Id != nil {
			servers = append(servers, *s.Id)
		}
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
	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(
			diags, errfmt.OpUpdate, "cloud_affinity_group_member", "", err, httpResp,
		)

		return false
	}

	return true
}
