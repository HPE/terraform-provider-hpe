// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package virtualswitch

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// ModifyPlan enforces two rules the API cannot express as an in-place change:
//
//   - A segment's traffic_type or mtu is not editable; changing either on an
//     existing map key forces the whole resource to be replaced.
//   - A host uplink cannot be removed from a switch in Morpheus (there is no
//     API path that deletes an uplink); dropping a server_id from the plan is a
//     plan-time error, not a silent no-op or an implicit replacement.
func (r *Resource) ModifyPlan(
	ctx context.Context,
	req resource.ModifyPlanRequest,
	resp *resource.ModifyPlanResponse,
) {
	// Nothing to compare on create or destroy.
	if req.Plan.Raw.IsNull() || req.State.Raw.IsNull() {
		return
	}

	var plan, state VirtualSwitchModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Replacement first: a plan that recreates the switch (a changed
	// cluster_id, base_type or code, or a traffic_type / mtu change on an
	// existing segment) naturally drops every uplink, so the uplink-removal
	// rule below does not apply to it.
	replace, d := requiresReplaceForSegments(ctx, state, plan)
	if resp.Diagnostics.Append(d...); resp.Diagnostics.HasError() {
		return
	}

	if replace {
		resp.RequiresReplace = append(resp.RequiresReplace, path.Root("segments"))

		return
	}

	if replacesSwitch(state, plan) {
		return
	}

	// Uplink removal is impossible in the API.
	removed, d := uplinkRemovals(ctx, state, plan)
	if resp.Diagnostics.Append(d...); resp.Diagnostics.HasError() {
		return
	}

	if len(removed) > 0 {
		resp.Diagnostics.AddAttributeError(
			path.Root("uplinks"),
			"Cannot remove a host uplink",
			fmt.Sprintf(
				"Morpheus cannot remove a host from a virtual switch (server_id %v). "+
					"Remove the host from the cluster, or recreate the switch, instead of "+
					"dropping its uplink from the configuration.", removed),
		)

		return
	}

	// Only when the update actually mutates the switch do we mark the volatile
	// server-side computed sub-attributes of uplinks and segments unknown, so
	// Terraform does not expect the planned (prior-state) values to survive the
	// async apply. Doing this on a no-op refresh plan would make the plan
	// spuriously non-empty. See markComputedUnknownForUpdate.
	up, d := computeUpdatePlan(ctx, state, plan)
	if resp.Diagnostics.Append(d...); resp.Diagnostics.HasError() {
		return
	}

	if up.empty() {
		return
	}

	if resp.Diagnostics.Append(markComputedUnknownForUpdate(ctx, &plan, up.keptSegIDs)...); resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.Plan.Set(ctx, &plan)...)
}

// requiresReplaceForSegments reports whether any segment that exists in both
// state and plan changes its traffic_type or mtu (neither is editable by the
// API). Matching is by key first, then by identity (traffic_type, vlan_id), the
// same rule computeUpdatePlan uses, so a pure re-key (e.g. adopting the user's
// key for an import-synthesised one) is NOT mistaken for a traffic_type change.
//
// A segment that fails to decode is reported, not treated as "no change": a
// silent false here could plan an impossible in-place edit.
func requiresReplaceForSegments(ctx context.Context, state, plan VirtualSwitchModel) (bool, diag.Diagnostics) {
	planSegs, diags := segmentInputs(ctx, plan)
	if diags.HasError() {
		return false, diags
	}

	stateSegs, d := segmentInputs(ctx, state)
	if diags.Append(d...); diags.HasError() {
		return false, diags
	}

	stateByKey := make(map[string]segmentInput, len(stateSegs))
	for _, s := range stateSegs {
		stateByKey[s.key] = s
	}

	stateMatched := make(map[string]bool, len(stateSegs))

	for _, ps := range planSegs {
		prior, matchKey, matched := matchStateSegment(ps, stateSegs, stateByKey, stateMatched)
		if !matched {
			continue
		}

		stateMatched[matchKey] = true

		// An unknown planned traffic_type (a deferred value) is not a change
		// yet; it is compared once known.
		if !ps.trafficTypeUnknown && prior.trafficType != ps.trafficType {
			return true, diags
		}

		// Segment mtu is stored by Morpheus but never reaches the hosts, and a
		// segment added through a composite edit comes back with the switch's
		// mtu filled in. Only a value the practitioner set forces a
		// replacement: one that differs from the stored value, or one set on a
		// segment Morpheus holds no mtu for (the API cannot edit a segment's
		// mtu in place, so there is no update that could apply it). An unset
		// mtu follows whatever Morpheus holds.
		if !ps.mtuUnknown && ps.mtu != nil && (prior.mtu == nil || *ps.mtu != *prior.mtu) {
			return true, diags
		}
	}

	return false, diags
}

// replacesSwitch reports whether an attribute with RequiresReplace differs
// between state and plan (an unknown planned value counts: it may differ).
func replacesSwitch(state, plan VirtualSwitchModel) bool {
	changed := func(s, p interface {
		IsUnknown() bool
		Equal(attr.Value) bool
	}, pv attr.Value,
	) bool {
		return p.IsUnknown() || !s.Equal(pv)
	}

	return changed(state.ClusterId, plan.ClusterId, plan.ClusterId) ||
		changed(state.BaseType, plan.BaseType, plan.BaseType) ||
		changed(state.Code, plan.Code, plan.Code)
}
