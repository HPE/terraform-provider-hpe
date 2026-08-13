// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package clusteraffinitygroupmember

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
)

// Read checks whether the server is still in the group.
//
// This is a presence check, not a comparison against a recorded list. The group
// gains and loses other members for reasons that have nothing to do with this
// resource, and none of that is drift as far as this membership is concerned.
func (r *Resource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
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

	result, httpResp, err := client.ClustersAPI.
		GetClusterAffinityGroup(ctx, clusterID, groupID).Execute()

	// The group itself is gone, so the membership is too. Deleting a cluster or
	// its resource pool cascades to its affinity groups.
	if errfmt.IsNotFound(httpResp) {
		resp.State.RemoveResource(ctx)

		return
	}

	if err := errfmt.CheckResponse(err, httpResp); err != nil {
		errfmt.DiagError(
			&resp.Diagnostics, errfmt.OpRead, "cluster_affinity_group_member", "", err, httpResp,
		)

		return
	}

	if result == nil || result.AffinityGroup == nil {
		resp.Diagnostics.AddError("API returned nil", "AffinityGroup is nil in the response")

		return
	}

	found := false

	for _, s := range result.AffinityGroup.Servers {
		if s.Id != nil && *s.Id == serverID {
			found = true

			break
		}
	}

	// Removed by something else. Drop it from state so the next plan offers to
	// put it back, rather than reporting success for a membership that no
	// longer exists.
	if !found {
		resp.State.RemoveResource(ctx)

		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
