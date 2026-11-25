// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package clusteraffinitygroup

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/affinityshape"
)

var _ resource.ResourceWithModifyPlan = &clusterAffinityGroupResource{}

// ModifyPlan refuses, at plan time, what Morpheus 9.1 would reject at apply:
// creating or replacing an affinity group with no servers, or updating one to
// `servers = []` (MORPH-20083).
// servers is Optional+Computed: set, this resource owns the membership and the
// group is created with it; unset, membership belongs to the member
// resources, whose Delete handles the last server (see affinityshape). The appliance version is read here, once per plan of a new
// group; an unreadable version fails open.
func (r *clusterAffinityGroupResource) ModifyPlan(
	ctx context.Context,
	req resource.ModifyPlanRequest,
	resp *resource.ModifyPlanResponse,
) {
	if req.Plan.Raw.IsNull() {
		return // destroy
	}

	// The CONFIG value: servers is Optional+Computed, so an unconfigured
	// attribute is unknown in the plan but null in the configuration.
	var servers types.Set
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("servers"), &servers)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// A create or a replacement seeds a new group, so an unset or empty
	// servers is checked. An in-place update is checked only when servers is
	// configured empty -- `servers = []` would empty the group, which 9.1
	// refuses -- while an unset servers on update is fine: membership is then
	// left to the member resources and nothing is sent.
	update := !req.State.Raw.IsNull() && len(resp.RequiresReplace) == 0
	if update && (servers.IsNull() || servers.IsUnknown()) {
		return
	}

	if !r.Configured() {
		return // no provider data on this RPC; the CRUD gate reports real problems
	}
	client, err := r.NewClient(ctx)
	if err != nil {
		return
	}

	resp.Diagnostics.Append(affinityshape.RequireCreate(ctx, client, servers)...)
}
