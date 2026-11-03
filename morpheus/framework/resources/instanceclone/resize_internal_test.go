// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package instanceclone

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// These drive the resize pairing directly. The behaviour they encode was
// observed on an HVM appliance: every in-place update of a clone issued a
// resize, and the resize carried the volumes without their ids. Both were
// invisible while the clone had only a root disk; once it also had a CD-ROM
// (any non-root volume with an external id) the appliance answered the
// bookkeeping update that removes the deprecated source_instance_id_version
// marker with a 500, because its removed-disk pass dereferences the id of
// every submitted volume.

// vol is a known volume named "root" with the given id (0 for unknown, as the
// framework plans a computed id on update) and size; the computed-optional
// fields are unknown, as an unconfigured plan leaves them.
func vol(id, size int64) VolumesValue {
	return namedVol("root", id, size)
}

// namedVol is vol with a name.
func namedVol(name string, id, size int64) VolumesValue {
	v := VolumesValue{
		state:                attr.ValueStateKnown,
		Id:                   types.Int64Unknown(),
		Name:                 types.StringValue(name),
		Size:                 types.Int64Value(size),
		RootVolume:           types.BoolValue(true),
		DatastoreId:          types.Int64Unknown(),
		StorageType:          types.Int64Unknown(),
		ControllerMountPoint: types.StringUnknown(),
		SizeId:               types.Int64Null(),
		StorageProfile:       types.StringUnknown(),
	}
	if id != 0 {
		v.Id = types.Int64Value(id)
	}

	return v
}

// stateVol is vol as the prior state holds it: everything known.
func stateVol(id, size, datastore int64) VolumesValue {
	return namedStateVol("root", id, size, datastore)
}

// namedStateVol is stateVol with a name.
func namedStateVol(name string, id, size, datastore int64) VolumesValue {
	v := namedVol(name, id, size)
	v.DatastoreId = types.Int64Value(datastore)
	v.StorageType = types.Int64Value(1)
	v.ControllerMountPoint = types.StringValue("1:0:0:0")
	v.StorageProfile = types.StringValue("kvm-cache-none")

	return v
}

func volList(t *testing.T, vols ...VolumesValue) types.List {
	t.Helper()

	elems := make([]attr.Value, 0, len(vols))
	for _, v := range vols {
		elems = append(elems, v)
	}

	l, d := types.ListValue(VolumesValue{}.Type(context.Background()), elems)
	if d.HasError() {
		t.Fatalf("volume list: %v", d)
	}

	return l
}

func iface(id, network int64) NetworkInterfacesValue {
	ni := NetworkInterfacesValue{
		state:                  attr.ValueStateKnown,
		Id:                     types.Int64Unknown(),
		NetworkId:              types.Int64Value(network),
		NetworkInterfaceTypeId: types.Int64Unknown(),
		IpMode:                 types.StringValue(""),
		IpAddress:              types.StringUnknown(),
		MacAddress:             types.StringUnknown(),
		ChildVirtualNetworks:   types.ListUnknown(ChildVirtualNetworksValue{}.Type(context.Background())),
	}
	if id != 0 {
		ni.Id = types.Int64Value(id)
	}

	return ni
}

func stateIface(id, network int64) NetworkInterfacesValue {
	ni := iface(id, network)
	ni.NetworkInterfaceTypeId = types.Int64Value(3)
	ni.IpAddress = types.StringValue("10.0.0.5")
	ni.MacAddress = types.StringNull()
	ni.ChildVirtualNetworks = types.ListNull(ChildVirtualNetworksValue{}.Type(context.Background()))

	return ni
}

func ifaceList(t *testing.T, ifaces ...NetworkInterfacesValue) types.List {
	t.Helper()

	elems := make([]attr.Value, 0, len(ifaces))
	for _, ni := range ifaces {
		elems = append(elems, ni)
	}

	l, d := types.ListValue(NetworkInterfacesValue{}.Type(context.Background()), elems)
	if d.HasError() {
		t.Fatalf("interface list: %v", d)
	}

	return l
}

func TestUnitResizeVolumesCarriesIds(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	current := []VolumesValue{stateVol(290, 20, 10), stateVol(291, 5, 10)}

	// A planned update: ids unknown, a third disk added.
	plan := volList(t, vol(0, 20), vol(0, 5), vol(0, 8))

	got, changed, d := resizeVolumes(ctx, plan, current)
	if d.HasError() {
		t.Fatalf("diagnostics: %v", d)
	}

	if !changed {
		t.Fatal("an added volume is a change")
	}

	want := []int64{290, 291, -1}
	if len(got) != len(want) {
		t.Fatalf("got %d volumes, want %d", len(got), len(want))
	}

	for i, w := range want {
		if got[i].Id == nil {
			t.Fatalf("volume %d: id not sent; the appliance dereferences it", i)
		}

		if *got[i].Id != w {
			t.Errorf("volume %d: id %d, want %d", i, *got[i].Id, w)
		}
	}
}

func TestUnitResizeVolumesPrefersPlannedId(t *testing.T) {
	t.Parallel()

	got, _, d := resizeVolumes(
		context.Background(), volList(t, vol(77, 20)), []VolumesValue{stateVol(290, 20, 10)},
	)
	if d.HasError() {
		t.Fatalf("diagnostics: %v", d)
	}

	if got[0].Id == nil || *got[0].Id != 77 {
		t.Fatalf("a known planned id is the one to send, got %v", got[0].Id)
	}
}

func TestUnitResizeVolumesChangeDetection(t *testing.T) {
	t.Parallel()

	current := []VolumesValue{stateVol(290, 20, 10)}

	moved := vol(0, 20)
	moved.DatastoreId = types.Int64Value(11)

	profiled := vol(0, 20)
	profiled.StorageProfile = types.StringValue("kvm-cache-writeback")

	renamed := vol(0, 20)
	renamed.Name = types.StringValue("system")

	cases := []struct {
		name    string
		plan    types.List
		changed bool
	}{
		{"same shape, computed fields unknown", volList(t, vol(0, 20)), false},
		{"same shape, everything known", volList(t, stateVol(290, 20, 10)), false},
		{"rename only", volList(t, renamed), false},
		{"size grown", volList(t, vol(0, 30)), true},
		{"size shrunk", volList(t, vol(0, 10)), true},
		{"datastore moved", volList(t, moved), true},
		{"storage profile set", volList(t, profiled), true},
		{"volume added", volList(t, vol(0, 20), vol(0, 5)), true},
		{"volume removed", volList(t), true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			_, changed, d := resizeVolumes(context.Background(), tc.plan, current)
			if d.HasError() {
				t.Fatalf("diagnostics: %v", d)
			}

			if changed != tc.changed {
				t.Errorf("changed = %v, want %v", changed, tc.changed)
			}
		})
	}
}

func TestUnitResizeVolumesNullPlan(t *testing.T) {
	t.Parallel()

	got, changed, d := resizeVolumes(
		context.Background(),
		types.ListNull(VolumesValue{}.Type(context.Background())),
		[]VolumesValue{stateVol(290, 20, 10)},
	)
	if d.HasError() {
		t.Fatalf("diagnostics: %v", d)
	}

	if got != nil || changed {
		t.Fatalf("no planned volumes is nothing to send, got %v changed=%v", got, changed)
	}
}

func TestUnitResizeNetworkInterfaces(t *testing.T) {
	t.Parallel()

	current := []NetworkInterfacesValue{stateIface(40, 6)}

	cases := []struct {
		name    string
		plan    types.List
		changed bool
		ids     []*int64
	}{
		{"same interface, id unknown", ifaceList(t, iface(0, 6)), false, []*int64{ptr(40)}},
		{"same interface, id known", ifaceList(t, iface(40, 6)), false, []*int64{ptr(40)}},
		{"network changed", ifaceList(t, iface(0, 8)), true, []*int64{ptr(40)}},
		{"interface added", ifaceList(t, iface(0, 6), iface(0, 8)), true, []*int64{ptr(40), nil}},
		{"interface removed", ifaceList(t), true, nil},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got, changed, d := resizeNetworkInterfaces(context.Background(), tc.plan, current)
			if d.HasError() {
				t.Fatalf("diagnostics: %v", d)
			}

			if changed != tc.changed {
				t.Errorf("changed = %v, want %v", changed, tc.changed)
			}

			if len(got) != len(tc.ids) {
				t.Fatalf("got %d interfaces, want %d", len(got), len(tc.ids))
			}

			for i, want := range tc.ids {
				switch {
				case want == nil && got[i].Id != nil:
					t.Errorf("interface %d: id %d sent for a new interface", i, *got[i].Id)
				case want != nil && (got[i].Id == nil || *got[i].Id != *want):
					t.Errorf("interface %d: id %v, want %d", i, got[i].Id, *want)
				}
			}
		})
	}
}

func TestUnitPlanIDChanged(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name        string
		plan, state types.Int64
		changed     bool
	}{
		{"unset leaves the plan alone", types.Int64Null(), types.Int64Value(20), false},
		{"unknown leaves the plan alone", types.Int64Unknown(), types.Int64Value(20), false},
		{"same plan", types.Int64Value(20), types.Int64Value(20), false},
		{"different plan", types.Int64Value(21), types.Int64Value(20), true},
		{"plan set where there was none", types.Int64Value(21), types.Int64Null(), true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if got := planIDChanged(tc.plan, tc.state); got != tc.changed {
				t.Errorf("planIDChanged = %v, want %v", got, tc.changed)
			}
		})
	}
}

func TestUnitCloneNeedsResizeGrowsOnly(t *testing.T) {
	t.Parallel()

	actual := []VolumesValue{stateVol(290, 20, 10)}

	cases := []struct {
		name string
		plan types.List
		need bool
	}{
		{"as cloned", volList(t, vol(0, 20)), false},
		{"smaller than cloned", volList(t, vol(0, 10)), false},
		{"larger than cloned", volList(t, vol(0, 30)), true},
		{"extra volume not yet present", volList(t, vol(0, 20), vol(0, 5)), false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			need, d := cloneNeedsResize(context.Background(), tc.plan, actual)
			if d.HasError() {
				t.Fatalf("diagnostics: %v", d)
			}

			if need != tc.need {
				t.Errorf("cloneNeedsResize = %v, want %v", need, tc.need)
			}
		})
	}
}

// TestUnitResizeRequestForUpdateBookkeepingOnly is the observed failure: an
// update that changes nothing the resize acts on must not issue one.
func TestUnitResizeRequestForUpdateBookkeepingOnly(t *testing.T) {
	t.Parallel()

	state := InstanceCloneModel{
		Id:                      types.Int64Value(103),
		PlanId:                  types.Int64Null(),
		SourceInstanceIdVersion: types.Int64Value(1),
		Volumes:                 volList(t, stateVol(290, 20, 10)),
		NetworkInterfaces:       ifaceList(t, stateIface(40, 6)),
	}

	plan := state
	plan.SourceInstanceIdVersion = types.Int64Null()
	plan.Volumes = volList(t, vol(0, 20))
	plan.NetworkInterfaces = ifaceList(t, iface(0, 6))

	req, resizing, d := resizeRequestForUpdate(context.Background(), plan, state)
	if d.HasError() {
		t.Fatalf("diagnostics: %v", d)
	}

	if resizing {
		t.Fatalf("a bookkeeping update must not resize, built %+v", req)
	}
}

func TestUnitResizeRequestForUpdateSendsOnlyWhatChanged(t *testing.T) {
	t.Parallel()

	state := InstanceCloneModel{
		Id:                types.Int64Value(103),
		PlanId:            types.Int64Value(20),
		Volumes:           volList(t, stateVol(290, 20, 10)),
		NetworkInterfaces: ifaceList(t, stateIface(40, 6)),
	}

	plan := state
	plan.Volumes = volList(t, vol(0, 30))
	plan.NetworkInterfaces = ifaceList(t, iface(0, 6))

	req, resizing, d := resizeRequestForUpdate(context.Background(), plan, state)
	if d.HasError() {
		t.Fatalf("diagnostics: %v", d)
	}

	if !resizing {
		t.Fatal("a grown volume is a resize")
	}

	if req.Instance == nil || req.Instance.Plan == nil || req.Instance.Plan.Id == nil ||
		*req.Instance.Plan.Id != 20 {
		t.Errorf("the plan is always carried when set, got %+v", req.Instance)
	}

	if len(req.Volumes) != 1 || req.Volumes[0].Id == nil || *req.Volumes[0].Id != 290 ||
		req.Volumes[0].Size == nil || *req.Volumes[0].Size != 30 {
		t.Errorf("volumes = %+v, want the grown volume with its id", req.Volumes)
	}

	if req.NetworkInterfaces != nil {
		t.Errorf("unchanged interfaces must not be sent, got %+v", req.NetworkInterfaces)
	}
}

// TestUnitResizeVolumesPairsByNameWhenCountsDiffer is the reviewer's case:
// dropping the middle disk of three must leave the third disk its own id, not
// hand it the dropped disk's, or the appliance would keep the wrong one.
func TestUnitResizeVolumesPairsByNameWhenCountsDiffer(t *testing.T) {
	t.Parallel()

	current := []VolumesValue{
		namedStateVol("root", 290, 20, 10),
		namedStateVol("data", 291, 5, 10),
		namedStateVol("logs", 292, 8, 10),
	}

	cases := []struct {
		name string
		plan types.List
		ids  []int64
	}{
		{
			"middle removed",
			volList(t, namedVol("root", 0, 20), namedVol("logs", 0, 8)),
			[]int64{290, 292},
		},
		{
			"first removed",
			volList(t, namedVol("data", 0, 5), namedVol("logs", 0, 8)),
			[]int64{291, 292},
		},
		{
			"added in the middle",
			volList(t, namedVol("root", 0, 20), namedVol("scratch", 0, 4), namedVol("data", 0, 5), namedVol("logs", 0, 8)),
			[]int64{290, -1, 291, 292},
		},
		{
			"removed and one renamed: the renamed one is new",
			volList(t, namedVol("root", 0, 20), namedVol("audit", 0, 8)),
			[]int64{290, -1},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got, changed, d := resizeVolumes(context.Background(), tc.plan, current)
			if d.HasError() {
				t.Fatalf("diagnostics: %v", d)
			}

			if !changed {
				t.Fatal("a changed volume count is a change")
			}

			if len(got) != len(tc.ids) {
				t.Fatalf("got %d volumes, want %d", len(got), len(tc.ids))
			}

			for i, want := range tc.ids {
				if got[i].Id == nil || *got[i].Id != want {
					t.Errorf("volume %d (%s): id %v, want %d", i, *got[i].Name, got[i].Id, want)
				}
			}
		})
	}
}

// TestUnitResizeVolumesPositionalWhenCountsEqual: the same disks in the same
// order pair by position, as on the instance resource, and are not a change.
// A reordered list never gets this far: ModifyPlan refuses it (see
// TestUnitVolumeNameChanges), because paired by position it would send one
// disk's id under another's name and size.
func TestUnitResizeVolumesPositionalWhenCountsEqual(t *testing.T) {
	t.Parallel()

	current := []VolumesValue{namedStateVol("root", 290, 20, 10), namedStateVol("data", 291, 5, 10)}
	plan := volList(t, namedVol("root", 0, 20), namedVol("data", 0, 5))

	got, changed, d := resizeVolumes(context.Background(), plan, current)
	if d.HasError() {
		t.Fatalf("diagnostics: %v", d)
	}

	if changed {
		t.Fatal("the same two disks are not a change")
	}

	if *got[0].Id != 290 || *got[1].Id != 291 {
		t.Errorf("ids %d, %d; want 290, 291", *got[0].Id, *got[1].Id)
	}
}

func TestUnitVolumeNameChanges(t *testing.T) {
	t.Parallel()

	state := []VolumesValue{namedStateVol("root", 290, 20, 10), namedStateVol("data", 291, 5, 10)}

	cases := []struct {
		name    string
		plan    []VolumesValue
		changes int
		reorder bool
	}{
		{"same names", []VolumesValue{namedVol("root", 0, 20), namedVol("data", 0, 5)}, 0, false},
		{"one renamed", []VolumesValue{namedVol("root", 0, 20), namedVol("logs", 0, 5)}, 1, false},
		{"both renamed", []VolumesValue{namedVol("system", 0, 20), namedVol("logs", 0, 5)}, 2, false},
		// The reviewer's case: the same disks in the other order is reported at
		// both positions and recognised as a reorder, not a rename.
		{"reordered", []VolumesValue{namedVol("data", 0, 5), namedVol("root", 0, 20)}, 2, true},
		{"a size change is not a rename", []VolumesValue{namedVol("root", 0, 30), namedVol("data", 0, 5)}, 0, false},
		{"counts differ: pairing is by name, not a rename", []VolumesValue{namedVol("logs", 0, 5)}, 0, false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got, reorder := volumeNameChanges(tc.plan, state)
			if len(got) != tc.changes {
				t.Errorf("volumeNameChanges = %v, want %d changes", got, tc.changes)
			}

			if reorder != tc.reorder {
				t.Errorf("reorder = %v, want %v", reorder, tc.reorder)
			}
		})
	}
}

// TestUnitResizeNetworkInterfacesPairsByNetworkWhenCountsDiffer: removing the
// middle NIC of three leaves the third its id.
func TestUnitResizeNetworkInterfacesPairsByNetworkWhenCountsDiffer(t *testing.T) {
	t.Parallel()

	current := []NetworkInterfacesValue{stateIface(40, 6), stateIface(41, 8), stateIface(42, 9)}
	plan := ifaceList(t, iface(0, 6), iface(0, 9))

	got, changed, d := resizeNetworkInterfaces(context.Background(), plan, current)
	if d.HasError() {
		t.Fatalf("diagnostics: %v", d)
	}

	if !changed {
		t.Fatal("a removed interface is a change")
	}

	if got[0].Id == nil || *got[0].Id != 40 || got[1].Id == nil || *got[1].Id != 42 {
		t.Errorf("ids %v, %v; want 40, 42", got[0].Id, got[1].Id)
	}
}

func ptr(v int64) *int64 { return &v }
