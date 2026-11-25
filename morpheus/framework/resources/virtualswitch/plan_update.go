// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package virtualswitch

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

// updatePlan is the result of diffing prior state against the plan: the set of
// mutations Update must issue, in the order the API's sync guard requires.
type updatePlan struct {
	renameTo *string // non-nil when name changed

	removedSegIDs []int64 // segment ids to DELETE (removed map keys)

	// compositeEdit inputs
	bondModeChanged bool
	bondMode        *string // "none" when cleared; nil when not changed
	mtuChanged      bool
	mtu             *int32
	uplinksChanged  bool
	uplinks         []uplinkInput // full current list when any uplink changed

	newSegments    []segmentInput  // new map keys
	editedSegments []editedSegment // existing keys whose editable fields changed

	// keptSegIDs maps each plan segment key that matched a prior segment to that
	// segment's id. An existing segment keeps its id across an update, so the
	// plan can carry it as known; a network built on the segment references it,
	// and an unknown id would plan that network for replacement.
	keptSegIDs map[string]int64

	descriptionChanged bool
	description        *string // PUT-only; compositeEdit ignores it
}

type editedSegment struct {
	segID   int64
	segment segmentInput
}

func (p updatePlan) hasCompositeWork() bool {
	return p.bondModeChanged || p.mtuChanged || p.uplinksChanged ||
		len(p.newSegments) > 0 || len(p.editedSegments) > 0
}

func (p updatePlan) empty() bool {
	return p.renameTo == nil && len(p.removedSegIDs) == 0 &&
		!p.hasCompositeWork() && !p.descriptionChanged
}

// computeUpdatePlan diffs state against plan.
func computeUpdatePlan(ctx context.Context, state, plan VirtualSwitchModel) (updatePlan, diag.Diagnostics) {
	var up updatePlan
	var diags diag.Diagnostics

	if !plan.Name.Equal(state.Name) {
		v := plan.Name.ValueString()
		up.renameTo = &v
	}

	if !plan.Description.Equal(state.Description) {
		up.descriptionChanged = true
		up.description = strToPtr(plan.Description)
	}

	if !plan.BondMode.Equal(state.BondMode) {
		up.bondModeChanged = true
		if plan.BondMode.IsNull() || plan.BondMode.IsUnknown() {
			v := bondModeNone
			up.bondMode = &v
		} else {
			up.bondMode = strToPtr(plan.BondMode)
		}
	}

	if !plan.Mtu.Equal(state.Mtu) {
		up.mtuChanged = true
		up.mtu = int64ToInt32Ptr(plan.Mtu)
	}

	planUps, d := uplinkInputs(ctx, plan)
	if diags.Append(d...); diags.HasError() {
		return up, diags
	}

	stateUps, d := uplinkInputs(ctx, state)
	if diags.Append(d...); diags.HasError() {
		return up, diags
	}

	if uplinksDiffer(stateUps, planUps) {
		up.uplinksChanged = true
		up.uplinks = planUps
	}

	planSegs, d := segmentInputs(ctx, plan)
	if diags.Append(d...); diags.HasError() {
		return up, diags
	}

	stateSegs, d := segmentInputs(ctx, state)
	if diags.Append(d...); diags.HasError() {
		return up, diags
	}

	// Match plan segments to state segments by key first, then by identity
	// (traffic_type, vlan_id), which the server keeps unique per switch. Matching
	// by identity means a re-keyed segment (e.g. after import, where keys are
	// synthesised) is recognised as the SAME segment under a new key -- an edit
	// or no-op, not a destructive delete+add that would tear down the segment
	// and its host IPs on the appliance.
	segIDByKey := segmentIDsByKey(ctx, state)

	stateByKey := make(map[string]segmentInput, len(stateSegs))
	for _, s := range stateSegs {
		stateByKey[s.key] = s
	}

	stateMatched := make(map[string]bool, len(stateSegs))
	up.keptSegIDs = make(map[string]int64, len(planSegs))

	for _, ps := range planSegs {
		prior, matchKey, matched := matchStateSegment(ps, stateSegs, stateByKey, stateMatched)
		if !matched {
			up.newSegments = append(up.newSegments, ps)

			continue
		}

		stateMatched[matchKey] = true

		if id := segIDByKey[matchKey]; id != 0 {
			up.keptSegIDs[ps.key] = id
		}

		if segmentEditableChanged(prior, ps) {
			id := segIDByKey[matchKey]
			if id == 0 {
				diags.AddError(
					"virtual switch segment has no id",
					fmt.Sprintf("segment %q is in state without an id and cannot be edited; refresh the state and retry", ps.key),
				)

				return up, diags
			}

			up.editedSegments = append(up.editedSegments, editedSegment{
				segID:   id,
				segment: ps,
			})
		}
	}

	for _, ss := range stateSegs {
		if !stateMatched[ss.key] {
			if id := segIDByKey[ss.key]; id != 0 {
				up.removedSegIDs = append(up.removedSegIDs, id)
			}
		}
	}

	return up, diags
}

// segmentIdentity is a segment's server-unique identity (traffic_type, vlan_id);
// Morpheus keeps it unique per switch, so it is a stable cross-key match.
type segmentIdentity struct {
	trafficType string
	vlan        int32
	hasVlan     bool
}

func identityOf(s segmentInput) segmentIdentity {
	id := segmentIdentity{trafficType: s.trafficType}
	if s.vlanID != nil {
		id.vlan = *s.vlanID
		id.hasVlan = true
	}

	return id
}

// matchStateSegment finds the prior-state segment a plan segment corresponds to:
// first by map key, then -- when the keys differ (e.g. import-synthesised keys
// vs the user's config keys) -- by identity. Returns the matched state segment,
// its state key, and whether a match was found. Already-matched state segments
// are skipped so two plan segments cannot claim the same state segment.
func matchStateSegment(
	ps segmentInput,
	stateSegs []segmentInput,
	stateByKey map[string]segmentInput,
	stateMatched map[string]bool,
) (segmentInput, string, bool) {
	if prior, ok := stateByKey[ps.key]; ok && !stateMatched[ps.key] {
		return prior, ps.key, true
	}

	want := identityOf(ps)
	for _, ss := range stateSegs {
		if stateMatched[ss.key] {
			continue
		}

		if identityOf(ss) == want {
			return ss, ss.key, true
		}
	}

	return segmentInput{}, "", false
}

func segmentIDsByKey(ctx context.Context, state VirtualSwitchModel) map[string]int64 {
	out := map[string]int64{}

	if state.Segments.IsNull() || state.Segments.IsUnknown() {
		return out
	}

	segsMap := make(map[string]SegmentsValue, len(state.Segments.Elements()))
	if d := state.Segments.ElementsAs(ctx, &segsMap, false); d.HasError() {
		return out
	}

	for key, seg := range segsMap {
		if !seg.Id.IsNull() && !seg.Id.IsUnknown() {
			out[key] = seg.Id.ValueInt64()
		}
	}

	return out
}

// segmentEditableChanged reports whether an existing segment's API-editable
// fields (vlan_id, ip_address, netmask, gateway, host_ips) changed. traffic_type
// and mtu are not editable and are handled in ModifyPlan (RequiresReplace).
func segmentEditableChanged(a, b segmentInput) bool {
	if !int32PtrEqual(a.vlanID, b.vlanID) ||
		!strPtrEqual(a.ipAddress, b.ipAddress) ||
		!strPtrEqual(a.netmask, b.netmask) ||
		!strPtrEqual(a.gateway, b.gateway) {
		return true
	}

	return !hostIPsEqual(a.hostIps, b.hostIps)
}

func uplinksDiffer(a, b []uplinkInput) bool {
	if len(a) != len(b) {
		return true
	}

	am := make(map[int64]uplinkInput, len(a))
	for _, u := range a {
		am[u.serverID] = u
	}

	for _, u := range b {
		prior, ok := am[u.serverID]
		if !ok || prior.nic1 != u.nic1 || !strPtrEqual(prior.nic2, u.nic2) {
			return true
		}
	}

	return false
}

func hostIPsEqual(a, b []hostIPInput) bool {
	if len(a) != len(b) {
		return false
	}

	am := make(map[int64]string, len(a))
	for _, h := range a {
		am[h.serverID] = h.ip
	}

	for _, h := range b {
		if ip, ok := am[h.serverID]; !ok || ip != h.ip {
			return false
		}
	}

	return true
}

func int32PtrEqual(a, b *int32) bool {
	if a == nil || b == nil {
		return a == b
	}

	return *a == *b
}

func strPtrEqual(a, b *string) bool {
	if a == nil || b == nil {
		return a == b
	}

	return *a == *b
}

// uplinkRemovals returns the server_ids present in state but absent from plan;
// removing a host from a switch is impossible in the API (plan-time error).
func uplinkRemovals(ctx context.Context, state, plan VirtualSwitchModel) ([]int64, diag.Diagnostics) {
	var diags diag.Diagnostics

	stateUps, d := uplinkInputs(ctx, state)
	if diags.Append(d...); diags.HasError() {
		return nil, diags
	}

	planUps, known, d := uplinkInputsKnown(ctx, plan)
	if diags.Append(d...); diags.HasError() {
		return nil, diags
	}

	// With the set or a server_id still unknown (a deferred data source) the
	// removal cannot be judged at plan time; Update checks it once known.
	if !known {
		return nil, diags
	}

	planIDs := make(map[int64]bool, len(planUps))
	for _, u := range planUps {
		planIDs[u.serverID] = true
	}

	var removed []int64
	for _, u := range stateUps {
		if !planIDs[u.serverID] {
			removed = append(removed, u.serverID)
		}
	}

	return removed, diags
}
