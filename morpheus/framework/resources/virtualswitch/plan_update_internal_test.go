// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package virtualswitch

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func uplinkSet(t *testing.T, ups []UplinksValue) types.Set {
	t.Helper()

	ctx := context.Background()
	elems := make([]attr.Value, 0, len(ups))
	for _, u := range ups {
		elems = append(elems, u)
	}

	s, d := types.SetValue(UplinksValue{}.Type(ctx), elems)
	if d.HasError() {
		t.Fatalf("uplink set: %v", d)
	}

	return s
}

func newUplink(serverID int64, nic1 string, nic2 *string) UplinksValue {
	n2 := types.StringNull()
	if nic2 != nil {
		n2 = types.StringValue(*nic2)
	}

	return UplinksValue{
		Id:          types.Int64Null(),
		Nic1:        types.StringValue(nic1),
		Nic2:        n2,
		ServerId:    types.Int64Value(serverID),
		ServerName:  types.StringNull(),
		SyncMessage: types.StringNull(),
		SyncStatus:  types.StringNull(),
		state:       attr.ValueStateKnown,
	}
}

func TestComputeUpdatePlan(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	vlan100 := int64(100)
	vlan200 := int64(200)

	state := VirtualSwitchModel{
		Name:     tfstr("old"),
		BondMode: types.StringNull(),
		Mtu:      tfint(1500),
		Segments: buildSegmentsMapForTest(t, map[string]SegmentsValue{
			"mgmt": newSeg(t, 1, "management", nil),
			"nfs":  newSeg(t, 2, "data-nfs", &vlan100),
			"old":  newSeg(t, 3, "vm", nil),
		}),
		Uplinks: uplinkSet(t, []UplinksValue{newUplink(7, "eth2", nil)}),
	}

	plan := VirtualSwitchModel{
		Name:     tfstr("new"),
		BondMode: types.StringNull(),
		Mtu:      tfint(1500),
		Segments: buildSegmentsMapForTest(t, map[string]SegmentsValue{
			"mgmt": newSeg(t, 1, "management", nil),
			"nfs":  newSeg(t, 2, "data-nfs", &vlan200),  // vlan edited
			"new2": newSeg(t, 0, "live-migration", nil), // new key
			// "old" removed
		}),
		Uplinks: uplinkSet(t, []UplinksValue{newUplink(7, "eth2", nil)}),
	}

	up, diags := computeUpdatePlan(ctx, state, plan)
	if diags.HasError() {
		t.Fatalf("diags: %v", diags)
	}

	if up.renameTo == nil || *up.renameTo != "new" {
		t.Fatalf("rename: %v", up.renameTo)
	}

	if len(up.removedSegIDs) != 1 || up.removedSegIDs[0] != 3 {
		t.Fatalf("removed: %v", up.removedSegIDs)
	}

	if len(up.newSegments) != 1 || up.newSegments[0].key != "new2" {
		t.Fatalf("new segments: %v", up.newSegments)
	}

	if len(up.editedSegments) != 1 || up.editedSegments[0].segID != 2 {
		t.Fatalf("edited segments: %v", up.editedSegments)
	}

	if up.uplinksChanged {
		t.Fatal("uplinks should not have changed")
	}
}

func TestUplinkRemovalDetected(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	state := VirtualSwitchModel{
		Uplinks: uplinkSet(t, []UplinksValue{newUplink(7, "eth2", nil), newUplink(8, "eth2", nil)}),
	}
	plan := VirtualSwitchModel{
		Uplinks: uplinkSet(t, []UplinksValue{newUplink(7, "eth2", nil)}),
	}

	removed, diags := uplinkRemovals(ctx, state, plan)
	if diags.HasError() {
		t.Fatalf("diags: %v", diags)
	}

	if len(removed) != 1 || removed[0] != 8 {
		t.Fatalf("removed: %v", removed)
	}
}

func TestTrafficTypeChangeRequiresReplace(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	state := VirtualSwitchModel{
		Segments: buildSegmentsMapForTest(t, map[string]SegmentsValue{
			"a": newSeg(t, 1, "vm", nil),
		}),
	}
	plan := VirtualSwitchModel{
		Segments: buildSegmentsMapForTest(t, map[string]SegmentsValue{
			"a": newSeg(t, 1, "management", nil),
		}),
	}

	if got, _ := requiresReplaceForSegments(ctx, state, plan); !got {
		t.Fatal("expected RequiresReplace for traffic_type change")
	}
}

// TestComputeUpdatePlanRekeyIsNoOp proves that re-keying a segment after import
// (state has synthesised keys, plan has the user's keys, same traffic_type and
// vlan_id) produces zero segment operations: matching is by identity, so the
// segment is recognised as the same one under a new key rather than being
// deleted and re-added.
func TestComputeUpdatePlanRekeyIsNoOp(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	vlan100 := int64(100)

	// State carries import-synthesised keys.
	state := VirtualSwitchModel{
		Name:     tfstr("sw"),
		BondMode: types.StringNull(),
		Mtu:      tfint(1500),
		Segments: buildSegmentsMapForTest(t, map[string]SegmentsValue{
			"vm":           newSeg(t, 1, "vm", nil),
			"data-nfs_100": newSeg(t, 2, "data-nfs", &vlan100),
		}),
		Uplinks: uplinkSet(t, []UplinksValue{newUplink(7, "eth2", nil)}),
	}

	// Plan uses the user's keys for the same segments.
	plan := VirtualSwitchModel{
		Name:     tfstr("sw"),
		BondMode: types.StringNull(),
		Mtu:      tfint(1500),
		Segments: buildSegmentsMapForTest(t, map[string]SegmentsValue{
			"vm":  newSeg(t, 1, "vm", nil),
			"nfs": newSeg(t, 2, "data-nfs", &vlan100),
		}),
		Uplinks: uplinkSet(t, []UplinksValue{newUplink(7, "eth2", nil)}),
	}

	up, diags := computeUpdatePlan(ctx, state, plan)
	if diags.HasError() {
		t.Fatalf("diags: %v", diags)
	}

	if len(up.newSegments) != 0 {
		t.Fatalf("expected no new segments, got %v", up.newSegments)
	}

	if len(up.removedSegIDs) != 0 {
		t.Fatalf("expected no removed segments, got %v", up.removedSegIDs)
	}

	if len(up.editedSegments) != 0 {
		t.Fatalf("expected no edited segments, got %v", up.editedSegments)
	}

	// A pure re-key must not look like a traffic_type change.
	if got, _ := requiresReplaceForSegments(ctx, state, plan); got {
		t.Fatal("re-key must not require replacement")
	}
}

// TestComputeUpdatePlanNewAndRemovedAfterRekey proves identity matching still
// classifies genuinely new and genuinely removed segments correctly alongside a
// re-keyed one.
func TestComputeUpdatePlanNewAndRemovedAfterRekey(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	vlan100 := int64(100)
	vlan200 := int64(200)

	state := VirtualSwitchModel{
		Name:     tfstr("sw"),
		BondMode: types.StringNull(),
		Mtu:      tfint(1500),
		Segments: buildSegmentsMapForTest(t, map[string]SegmentsValue{
			"data-nfs_100": newSeg(t, 2, "data-nfs", &vlan100), // re-keyed below
			"data-nfs_200": newSeg(t, 3, "data-nfs", &vlan200), // removed below
		}),
		Uplinks: uplinkSet(t, []UplinksValue{newUplink(7, "eth2", nil)}),
	}

	plan := VirtualSwitchModel{
		Name:     tfstr("sw"),
		BondMode: types.StringNull(),
		Mtu:      tfint(1500),
		Segments: buildSegmentsMapForTest(t, map[string]SegmentsValue{
			"nfs": newSeg(t, 2, "data-nfs", &vlan100), // same identity as data-nfs_100
			"vm":  newSeg(t, 0, "vm", nil),            // genuinely new
			// data-nfs_200 removed
		}),
		Uplinks: uplinkSet(t, []UplinksValue{newUplink(7, "eth2", nil)}),
	}

	up, diags := computeUpdatePlan(ctx, state, plan)
	if diags.HasError() {
		t.Fatalf("diags: %v", diags)
	}

	if len(up.newSegments) != 1 || up.newSegments[0].key != "vm" {
		t.Fatalf("expected one new segment (vm), got %v", up.newSegments)
	}

	if len(up.removedSegIDs) != 1 || up.removedSegIDs[0] != 3 {
		t.Fatalf("expected removal of seg id 3, got %v", up.removedSegIDs)
	}

	if len(up.editedSegments) != 0 {
		t.Fatalf("expected no edited segments (re-key only), got %v", up.editedSegments)
	}
}

// TestSegmentMtuReplaceOnlyWhenConfigured proves that segment mtu forces a
// replacement only for a configured value that Morpheus does not already hold:
// changed, or newly set where Morpheus held none (the API cannot edit it in
// place, and create does send it). Morpheus fills in the switch's mtu on a
// segment added through a composite edit, so a state value against a null plan
// value must not be read as a change.
func TestSegmentMtuReplaceOnlyWhenConfigured(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	withMtu := func(mtu *int64) SegmentsValue {
		s := newSeg(t, 1, "vm", nil)
		if mtu == nil {
			s.Mtu = types.Int64Null()
		} else {
			s.Mtu = types.Int64Value(*mtu)
		}

		return s
	}
	m1500, m9000 := int64(1500), int64(9000)

	cases := []struct {
		name        string
		state, plan *int64
		replace     bool
	}{
		{"backfilled by Morpheus, unset in config", &m9000, nil, false},
		{"unset on both sides", nil, nil, false},
		{"configured and unchanged", &m9000, &m9000, false},
		{"configured and changed", &m1500, &m9000, true},
		{"newly configured where Morpheus held none", nil, &m9000, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			state := VirtualSwitchModel{Segments: buildSegmentsMapForTest(t, map[string]SegmentsValue{"a": withMtu(tc.state)})}
			plan := VirtualSwitchModel{Segments: buildSegmentsMapForTest(t, map[string]SegmentsValue{"a": withMtu(tc.plan)})}
			if got, _ := requiresReplaceForSegments(ctx, state, plan); got != tc.replace {
				t.Fatalf("replace = %v, want %v", got, tc.replace)
			}
		})
	}
}

// TestUpdateKeepsExistingSegmentIDsKnown pins that an update of the switch
// leaves the ids of its existing segments known in the plan — under the plan's
// own keys, re-keyed or not — and only a new segment's id unknown. A network
// built on a segment references its id; an unknown id would plan that network
// for replacement on every edit of the switch.
func TestUpdateKeepsExistingSegmentIDsKnown(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	vlan100 := int64(100)
	vlan200 := int64(200)

	state := VirtualSwitchModel{
		Name:     tfstr("sw"),
		BondMode: types.StringNull(),
		Mtu:      tfint(1500),
		Segments: buildSegmentsMapForTest(t, map[string]SegmentsValue{
			"vm":           newSeg(t, 1, "vm", nil),
			"data-nfs_100": newSeg(t, 2, "data-nfs", &vlan100),
		}),
		Uplinks: uplinkSet(t, []UplinksValue{newUplink(7, "eth2", nil)}),
	}

	// The MTU changes, the nfs segment is re-keyed, and an iscsi segment is new.
	plan := VirtualSwitchModel{
		Name:     tfstr("sw"),
		BondMode: types.StringNull(),
		Mtu:      tfint(9000),
		Segments: buildSegmentsMapForTest(t, map[string]SegmentsValue{
			"vm":    newSeg(t, 1, "vm", nil),
			"nfs":   newSeg(t, 2, "data-nfs", &vlan100),
			"iscsi": newSeg(t, 0, "data-iscsi", &vlan200),
		}),
		Uplinks: uplinkSet(t, []UplinksValue{newUplink(7, "eth2", nil)}),
	}

	up, diags := computeUpdatePlan(ctx, state, plan)
	if diags.HasError() {
		t.Fatalf("diags: %v", diags)
	}

	if !up.mtuChanged || len(up.newSegments) != 1 {
		t.Fatalf("expected an MTU change and one new segment, got %+v", up)
	}

	if diags := markComputedUnknownForUpdate(ctx, &plan, up.keptSegIDs); diags.HasError() {
		t.Fatalf("diags: %v", diags)
	}

	var segs map[string]SegmentsValue
	if d := plan.Segments.ElementsAs(ctx, &segs, false); d.HasError() {
		t.Fatalf("diags: %v", d)
	}

	for key, want := range map[string]int64{"vm": 1, "nfs": 2} {
		if got := segs[key].Id; got.IsUnknown() || got.ValueInt64() != want {
			t.Errorf("segment %q: id = %v, want %d known", key, got, want)
		}
	}

	if got := segs["iscsi"].Id; !got.IsUnknown() {
		t.Errorf("new segment iscsi: id = %v, want unknown", got)
	}

	// Every other computed field is still unknown on the kept segments.
	if !segs["nfs"].Status.IsUnknown() || !segs["nfs"].InterfaceName.IsUnknown() {
		t.Error("kept segment: volatile computed fields must be unknown")
	}
}
