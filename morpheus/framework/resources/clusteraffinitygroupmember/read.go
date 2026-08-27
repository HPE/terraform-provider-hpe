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
		errfmt.DiagClientError(&resp.Diagnostics, err)

		return
	}

	clusterID := state.ClusterID.ValueInt64()
	groupID := state.AffinityGroupID.ValueInt64()
	serverID := state.ServerID.ValueInt64()

	// Uses readMembership so a group whose single-item endpoint is broken still
	// resolves via the listing. TODO(MORPH-15806).
	servers, ok := readMembership(ctx, client, clusterID, groupID, &resp.Diagnostics)
	if !ok {
		// The group itself is gone, so the membership is too. Deleting a cluster
		// or its resource pool cascades to its affinity groups.
		if resp.Diagnostics.HasError() {
			return
		}

		resp.State.RemoveResource(ctx)

		return
	}

	found := false

	for _, id := range servers {
		if id == serverID {
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
