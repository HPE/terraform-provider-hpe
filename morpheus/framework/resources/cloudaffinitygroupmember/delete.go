// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package cloudaffinitygroupmember

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/affinitylock"
)

// Delete removes the server from the group.
//
// The mirror of Create: read, drop the one server, write the whole list back,
// then confirm it is gone.
func (r *Resource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	var state memberModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError("could not create sdk client", err.Error())

		return
	}

	cloudID := state.CloudID.ValueInt64()
	groupID := state.AffinityGroupID.ValueInt64()
	serverID := state.ServerID.ValueInt64()

	defer affinitylock.Acquire(affinitylock.Cloud, cloudID, groupID)()

	servers, ok := readMembership(ctx, client, cloudID, groupID, &resp.Diagnostics)
	if !ok {
		return
	}

	// Already absent: nothing to do, and rewriting the list would touch other
	// members for no reason.
	if !containsServer(servers, serverID) {
		return
	}

	if !writeMembership(
		ctx, client, cloudID, groupID, withoutServer(servers, serverID), &resp.Diagnostics,
	) {
		return
	}

	after, ok := readMembership(ctx, client, cloudID, groupID, &resp.Diagnostics)
	if !ok {
		return
	}

	if containsServer(after, serverID) {
		resp.Diagnostics.AddError(
			"affinity group membership was not removed",
			"Server "+state.ServerID.String()+" is still a member of affinity group "+
				state.AffinityGroupID.String()+" after the update was accepted. "+
				"The group's membership is replaced as a whole on every write, so "+
				"another process managing the same group concurrently can restore "+
				"the server. Re-apply, and avoid managing this group's servers "+
				"attribute at the same time as this resource.",
		)
	}
}
