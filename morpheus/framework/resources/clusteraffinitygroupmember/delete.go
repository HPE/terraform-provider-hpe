// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package clusteraffinitygroupmember

import (
	"context"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/affinitylock"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
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

	clusterID := state.ClusterID.ValueInt64()
	groupID := state.AffinityGroupID.ValueInt64()
	serverID := state.ServerID.ValueInt64()

	defer affinitylock.Acquire(affinitylock.Cluster, clusterID, groupID)()

	_, httpResp, err := client.ClustersAPI.
		GetClusterAffinityGroup(ctx, clusterID, groupID).Execute()

	// The group is already gone, so the membership is too.
	if errfmt.IsNotFound(httpResp) {
		return
	}

	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(
			&resp.Diagnostics, errfmt.OpDelete, "cluster_affinity_group_member", "", err, httpResp,
		)

		return
	}

	servers, ok := readMembership(ctx, client, clusterID, groupID, &resp.Diagnostics)
	if !ok {
		return
	}

	// Already absent: nothing to do, and rewriting the list would touch other
	// members for no reason.
	if !containsServer(servers, serverID) {
		return
	}

	if !writeMembership(
		ctx, client, clusterID, groupID, withoutServer(servers, serverID), &resp.Diagnostics,
	) {
		return
	}

	after, ok := readMembership(ctx, client, clusterID, groupID, &resp.Diagnostics)
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

// Update exists only to satisfy the interface.
//
// Every attribute forces replacement, so Terraform never calls this: a
// membership cannot be changed, only created or destroyed.
func (r *Resource) Update(
	_ context.Context,
	_ resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	resp.Diagnostics.AddError(
		"update not supported",
		"Every attribute of this resource forces replacement, so it should never "+
			"be updated in place. Please report this as a bug.",
	)
}

// ImportState accepts cluster_id/affinity_group_id/server_id.
//
// A membership has no identifier of its own -- it is the pairing that exists,
// not an object -- so the import id has to carry all three parts.
func (r *Resource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	parts := strings.Split(req.ID, "/")

	const wantParts = 3

	if len(parts) != wantParts {
		resp.Diagnostics.AddAttributeError(
			importIDPath,
			"unexpected import identifier",
			"Expected cluster_id/affinity_group_id/server_id, for example 15056/6221/572308, "+
				"but got: "+req.ID,
		)

		return
	}

	names := []string{"cluster_id", "affinity_group_id", "server_id"}
	values := make([]int64, wantParts)

	for i, part := range parts {
		v, err := strconv.ParseInt(part, 10, 64)
		if err != nil {
			resp.Diagnostics.AddAttributeError(
				importIDPath,
				"unexpected import identifier",
				"Expected "+names[i]+" to be a number in cluster_id/affinity_group_id/server_id, "+
					"but got: "+part,
			)

			return
		}

		values[i] = v
	}

	state := memberModel{
		ID:              types.StringValue(membershipID(values[0], values[1], values[2])),
		ClusterID:       types.Int64Value(values[0]),
		AffinityGroupID: types.Int64Value(values[1]),
		ServerID:        types.Int64Value(values[2]),
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
