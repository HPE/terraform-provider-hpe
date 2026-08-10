// (C) Copyright 2025-2026 Hewlett Packard Enterprise Development LP

package instance

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

// TestUnitGetChildNetworksReadsSubnetID verifies that getChildNetworks reads subnet_id
// back from a child server interface's subnet association.
//
// When provisioning via subnet_id the Morpheus API resolves the subnet to its
// parent network (reported as network_id) but also returns the subnet itself
// (containerDetails.server.interfaces[].subnet, see _computeServerInterface.gson),
// so subnet_id round-trips on refresh rather than being write-only.
func TestUnitGetChildNetworksReadsSubnetID(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	parentID := int64(1)
	subIntfMap := map[int64][]int64{parentID: {2}}
	serverIntfsMap := map[int64]sdk.InstanceContainerServerInterfacesInner1{
		2: {
			Id:      sdk.PtrInt64(2),
			Network: &sdk.InstanceContainerServerInterfacesInner1Network{Id: sdk.PtrInt64(10)},
			Subnet:  &sdk.InstanceContainerServerInterfacesInner1Subnet{Id: sdk.PtrInt64(5)},
		},
	}

	list, diags := getChildNetworks(ctx, &parentID, subIntfMap, serverIntfsMap)
	if diags.HasError() {
		t.Fatalf("getChildNetworks returned diagnostics: %v", diags)
	}

	var children []ChildVirtualNetworksValue
	if d := list.ElementsAs(ctx, &children, false); d.HasError() {
		t.Fatalf("ElementsAs returned diagnostics: %v", d)
	}

	if len(children) != 1 {
		t.Fatalf("expected 1 child interface, got %d", len(children))
	}

	// subnet_id is read back from the subnet association...
	if got := children[0].SubnetId; !got.Equal(types.Int64Value(5)) {
		t.Errorf("child subnet_id = %v, want 5", got)
	}

	// ...and network_id reflects the subnet's resolved parent network.
	if got := children[0].NetworkId; !got.Equal(types.Int64Value(10)) {
		t.Errorf("child network_id = %v, want 10", got)
	}
}

// TestUnitGetChildNetworksNoSubnet verifies that a child interface with no subnet
// association yields a null subnet_id (e.g. when network_id was used directly).
func TestUnitGetChildNetworksNoSubnet(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	parentID := int64(1)
	subIntfMap := map[int64][]int64{parentID: {2}}
	serverIntfsMap := map[int64]sdk.InstanceContainerServerInterfacesInner1{
		2: {
			Id:      sdk.PtrInt64(2),
			Network: &sdk.InstanceContainerServerInterfacesInner1Network{Id: sdk.PtrInt64(10)},
		},
	}

	list, diags := getChildNetworks(ctx, &parentID, subIntfMap, serverIntfsMap)
	if diags.HasError() {
		t.Fatalf("getChildNetworks returned diagnostics: %v", diags)
	}

	var children []ChildVirtualNetworksValue
	if d := list.ElementsAs(ctx, &children, false); d.HasError() {
		t.Fatalf("ElementsAs returned diagnostics: %v", d)
	}

	if len(children) != 1 {
		t.Fatalf("expected 1 child interface, got %d", len(children))
	}

	if got := children[0].SubnetId; !got.IsNull() {
		t.Errorf("child subnet_id = %v, want null", got)
	}
}

// TestUnitReadServerUUIDFromOwnedContainer verifies that server_uuid is read
// back only from the container matching container_id, and that extra containers
// (added by hpe_morpheus_instance_node) do not change the value.
func TestUnitReadServerUUIDFromOwnedContainer(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		containers  []sdk.InstanceContainer2
		containerID types.Int64
		want        string
		wantNull    bool
	}{
		{
			name: "reads from owned container",
			containers: []sdk.InstanceContainer2{
				{Id: sdk.PtrInt64(100), Server: &sdk.InstanceContainerServer2{Uuid: sdk.PtrString("owned-uuid")}},
				{Id: sdk.PtrInt64(200), Server: &sdk.InstanceContainerServer2{Uuid: sdk.PtrString("node-uuid")}},
			},
			containerID: types.Int64Value(100),
			want:        "owned-uuid",
		},
		{
			name: "ignores other containers - scaling regression guard",
			containers: []sdk.InstanceContainer2{
				{Id: sdk.PtrInt64(100), Server: &sdk.InstanceContainerServer2{Uuid: sdk.PtrString("owned-uuid")}},
				{Id: sdk.PtrInt64(200), Server: &sdk.InstanceContainerServer2{Uuid: sdk.PtrString("different-uuid")}},
				{Id: sdk.PtrInt64(300), Server: &sdk.InstanceContainerServer2{Uuid: sdk.PtrString("third-uuid")}},
			},
			containerID: types.Int64Value(100),
			want:        "owned-uuid",
		},
		{
			name: "container not found returns null",
			containers: []sdk.InstanceContainer2{
				{Id: sdk.PtrInt64(200), Server: &sdk.InstanceContainerServer2{Uuid: sdk.PtrString("uuid-2")}},
			},
			containerID: types.Int64Value(999),
			wantNull:    true,
		},
		{
			name:        "null container_id returns null",
			containers:  []sdk.InstanceContainer2{{Id: sdk.PtrInt64(1), Server: &sdk.InstanceContainerServer2{Uuid: sdk.PtrString("x")}}},
			containerID: types.Int64Null(),
			wantNull:    true,
		},
		{
			name: "server with no uuid returns null",
			containers: []sdk.InstanceContainer2{
				{Id: sdk.PtrInt64(100), Server: &sdk.InstanceContainerServer2{Uuid: nil}},
			},
			containerID: types.Int64Value(100),
			wantNull:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := readServerUUIDFromOwnedContainer(tt.containers, tt.containerID)
			if tt.wantNull {
				if !got.IsNull() {
					t.Errorf("expected null, got %v", got)
				}

				return
			}
			if got.ValueString() != tt.want {
				t.Errorf("got %q, want %q", got.ValueString(), tt.want)
			}
		})
	}
}

// TestUnitRemoveExternalStorageVolumes verifies that storage-server (SAN) volumes —
// e.g. Alletra MP BMaaS LUNs exported to the instance's host by
// hpe_morpheus_storage_volume — are excluded from the instance's tracked
// volumes, while the instance's own provisioned disks (no storageServer) are
// retained in order. This prevents spurious drift on the instance's volumes.
func TestUnitRemoveExternalStorageVolumes(t *testing.T) {
	t.Parallel()

	volumes := []sdk.InstanceContainerServerVolume1{
		{Id: sdk.PtrInt64(1), Name: sdk.PtrString("root")},
		{Id: sdk.PtrInt64(2), Name: sdk.PtrString("data")},
		{
			Id:            sdk.PtrInt64(3),
			Name:          sdk.PtrString("alletra-bmaas-lun"),
			StorageServer: &sdk.InstanceContainerServerVolumeStorageServer1{},
		},
	}

	got := removeExternalStorageVolumes(volumes)

	if len(got) != 2 {
		t.Fatalf("expected 2 volumes after filtering, got %d", len(got))
	}
	for _, v := range got {
		if v.StorageServer != nil {
			t.Errorf("storage-server volume (id %v) was not filtered out", v.Id)
		}
	}
	if got[0].Id == nil || *got[0].Id != 1 || got[1].Id == nil || *got[1].Id != 2 {
		t.Errorf("provisioned disks not retained in order: got %v, %v", got[0].Id, got[1].Id)
	}
}

// TestUnitIsExternallyAttachedVolume verifies which volumes are treated as
// externally attached to the instance's host.
//
// A storageServer reference alone is not sufficient: when a disk is provisioned
// onto a datastore that is itself backed by a storage server (e.g. an Alletra MP
// datastore), Morpheus copies the datastore's storage server onto the instance's
// own disk (assignVolumeDatastore sets volume.datastore and volume.storageServer
// together). Excluding those disks emptied the volume list and wrote a null
// `volumes` into state, surfacing as "Provider produced inconsistent result
// after apply: .volumes ... but now null".
func TestUnitIsExternallyAttachedVolume(t *testing.T) {
	t.Parallel()

	storageServer := &sdk.InstanceContainerServerVolumeStorageServer1{}

	tests := []struct {
		name   string
		volume sdk.InstanceContainerServerVolume1
		want   bool
	}{
		{
			name:   "ordinary provisioned disk",
			volume: sdk.InstanceContainerServerVolume1{Id: sdk.PtrInt64(1), Name: sdk.PtrString("data")},
			want:   false,
		},
		{
			name: "externally attached array LUN",
			volume: sdk.InstanceContainerServerVolume1{
				Id:            sdk.PtrInt64(2),
				Name:          sdk.PtrString("alletra-bmaas-lun"),
				StorageServer: storageServer,
			},
			want: true,
		},
		{
			name: "root volume on a storage-server-backed datastore",
			volume: sdk.InstanceContainerServerVolume1{
				Id:            sdk.PtrInt64(3),
				Name:          sdk.PtrString("root"),
				RootVolume:    sdk.PtrBool(true),
				DatastoreId:   sdk.PtrInt64(406),
				StorageServer: storageServer,
			},
			want: false,
		},
		{
			name: "data volume on a storage-server-backed datastore",
			volume: sdk.InstanceContainerServerVolume1{
				Id:            sdk.PtrInt64(4),
				Name:          sdk.PtrString("data"),
				DatastoreId:   sdk.PtrInt64(406),
				StorageServer: storageServer,
			},
			want: false,
		},
		{
			name: "root volume from a storage server with no datastore",
			volume: sdk.InstanceContainerServerVolume1{
				Id:            sdk.PtrInt64(5),
				Name:          sdk.PtrString("root"),
				RootVolume:    sdk.PtrBool(true),
				StorageServer: storageServer,
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := isExternallyAttachedVolume(tt.volume); got != tt.want {
				t.Errorf("isExternallyAttachedVolume() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestUnitRemoveExternalStorageVolumesKeepsDatastoreBackedDisks verifies the
// reported failure case: an instance whose disks are provisioned on a
// storage-server-backed datastore keeps those disks, while an array LUN exported
// to the same host is still excluded.
func TestUnitRemoveExternalStorageVolumesKeepsDatastoreBackedDisks(t *testing.T) {
	t.Parallel()

	storageServer := &sdk.InstanceContainerServerVolumeStorageServer1{}

	volumes := []sdk.InstanceContainerServerVolume1{
		{
			Id:            sdk.PtrInt64(1),
			Name:          sdk.PtrString("root"),
			RootVolume:    sdk.PtrBool(true),
			DatastoreId:   sdk.PtrInt64(406),
			StorageServer: storageServer,
		},
		{
			Id:            sdk.PtrInt64(2),
			Name:          sdk.PtrString("alletra-bmaas-lun"),
			StorageServer: storageServer,
		},
	}

	got := removeExternalStorageVolumes(volumes)

	if len(got) != 1 {
		t.Fatalf("expected 1 volume after filtering, got %d", len(got))
	}
	if got[0].Id == nil || *got[0].Id != 1 {
		t.Errorf("expected the instance's root disk (id 1) to be retained, got %v", got[0].Id)
	}
}

// TestUnitBoolFromConfig verifies the defensive coercion of untyped instance config
// values (native bool or the string encodings Morpheus may use) that back the
// config_bmaas.enforce_raid_boot_volume attribute, including the fallback to the
// supplied default when the value is absent or unrecognised.
func TestUnitBoolFromConfig(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   interface{}
		def  bool
		want types.Bool
	}{
		{"native true", true, false, types.BoolValue(true)},
		{"native false", false, true, types.BoolValue(false)},
		{"string on", "on", false, types.BoolValue(true)},
		{"string off", "off", true, types.BoolValue(false)},
		{"string true mixed case", "True", false, types.BoolValue(true)},
		{"string false", "false", true, types.BoolValue(false)},
		{"absent uses default true", nil, true, types.BoolValue(true)},
		{"absent uses default false", nil, false, types.BoolValue(false)},
		{"unrecognised uses default", "maybe", true, types.BoolValue(true)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := boolFromConfig(tt.in, tt.def); !got.Equal(tt.want) {
				t.Errorf("boolFromConfig(%v, %v) = %v, want %v", tt.in, tt.def, got, tt.want)
			}
		})
	}
}

// TestUnitSelectedHostsFromConfig verifies parsing of the baremetal plugin's
// selectedHosts config value into the config_bmaas.selected_hosts list. The plugin
// stores each entry as an object with a "value" id (host.value as Long); bare ids
// and the various JSON number/string encodings are also tolerated. An absent or
// otherwise unusable value yields a null list.
func TestUnitSelectedHostsFromConfig(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name     string
		in       interface{}
		wantNull bool
		want     []int64
	}{
		{
			name: "objects with value (plugin shape)",
			in: []interface{}{
				map[string]interface{}{"value": float64(12)},
				map[string]interface{}{"value": float64(34)},
			},
			want: []int64{12, 34},
		},
		{
			name: "bare numeric ids",
			in:   []interface{}{float64(7), float64(8)},
			want: []int64{7, 8},
		},
		{
			name: "string values",
			in:   []interface{}{map[string]interface{}{"value": "55"}},
			want: []int64{55},
		},
		{name: "absent yields null", in: nil, wantNull: true},
		{name: "empty slice yields null", in: []interface{}{}, wantNull: true},
		{name: "wrong type yields null", in: "nope", wantNull: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, diags := selectedHostsFromConfig(ctx, tt.in)
			if diags.HasError() {
				t.Fatalf("selectedHostsFromConfig returned diagnostics: %v", diags)
			}
			if tt.wantNull {
				if !got.IsNull() {
					t.Errorf("expected null list, got %v", got)
				}

				return
			}
			var ids []int64
			if d := got.ElementsAs(ctx, &ids, false); d.HasError() {
				t.Fatalf("ElementsAs returned diagnostics: %v", d)
			}
			if len(ids) != len(tt.want) {
				t.Fatalf("got %d ids %v, want %d %v", len(ids), ids, len(tt.want), tt.want)
			}
			for i := range ids {
				if ids[i] != tt.want[i] {
					t.Errorf("id[%d] = %d, want %d", i, ids[i], tt.want[i])
				}
			}
		})
	}
}

// TestUnitNumberToInt64 verifies coercion of the JSON-decoded number representations
// that may appear in the untyped instance config map.
func TestUnitNumberToInt64(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		in     interface{}
		want   int64
		wantOK bool
	}{
		{"float64", float64(42), 42, true},
		{"int", 43, 43, true},
		{"int64", int64(44), 44, true},
		{"numeric string", "45", 45, true},
		{"non-numeric string", "x", 0, false},
		{"nil", nil, 0, false},
		{"bool unsupported", true, 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, ok := numberToInt64(tt.in)
			if ok != tt.wantOK || got != tt.want {
				t.Errorf("numberToInt64(%v) = (%d, %v), want (%d, %v)", tt.in, got, ok, tt.want, tt.wantOK)
			}
		})
	}
}

// TestUnitContainerIdPreservedAcrossReads verifies that container_id is
// preserved from the plan when already set (UseStateForUnknown), and populated
// from the first container when null (import/create).
func TestUnitContainerIdPreservedAcrossReads(t *testing.T) {
	t.Parallel()

	containers := []sdk.InstanceContainer2{
		{Id: sdk.PtrInt64(100), Server: &sdk.InstanceContainerServer2{Uuid: sdk.PtrString("uuid-a")}},
		{Id: sdk.PtrInt64(200), Server: &sdk.InstanceContainerServer2{Uuid: sdk.PtrString("uuid-b")}},
	}

	// Case 1: plan has container_id set => preserved
	plan := InstanceModel{}
	plan.ContainerId = types.Int64Value(100)

	var state InstanceModel
	if !plan.ContainerId.IsNull() && !plan.ContainerId.IsUnknown() {
		state.ContainerId = plan.ContainerId
	} else if len(containers) > 0 && containers[0].Id != nil {
		state.ContainerId = convert.Int64ToType(containers[0].Id)
	}

	if state.ContainerId.ValueInt64() != 100 {
		t.Errorf("expected container_id=100, got %d", state.ContainerId.ValueInt64())
	}

	// Case 2: plan container_id is null => populated from first container
	plan2 := InstanceModel{}
	plan2.ContainerId = types.Int64Null()

	var state2 InstanceModel
	if !plan2.ContainerId.IsNull() && !plan2.ContainerId.IsUnknown() {
		state2.ContainerId = plan2.ContainerId
	} else if len(containers) > 0 && containers[0].Id != nil {
		state2.ContainerId = convert.Int64ToType(containers[0].Id)
	}

	if state2.ContainerId.ValueInt64() != 100 {
		t.Errorf("expected container_id=100 from first container, got %d", state2.ContainerId.ValueInt64())
	}
}

// TestUnitContainerIdPreservedWhenAbsent verifies that when the recorded
// container is no longer present in the instance, the stored container_id is
// kept rather than being nulled or erroring.
func TestUnitContainerIdPreservedWhenAbsent(t *testing.T) {
	t.Parallel()

	// The recorded container (100) is gone; only 200 exists.
	containers := []sdk.InstanceContainer2{
		{Id: sdk.PtrInt64(200), Server: &sdk.InstanceContainerServer2{Uuid: sdk.PtrString("uuid-b")}},
	}

	plan := InstanceModel{}
	plan.ContainerId = types.Int64Value(100)

	// The logic preserves plan when non-null (UseStateForUnknown behaviour)
	var state InstanceModel
	if !plan.ContainerId.IsNull() && !plan.ContainerId.IsUnknown() {
		state.ContainerId = plan.ContainerId
	} else if len(containers) > 0 && containers[0].Id != nil {
		state.ContainerId = convert.Int64ToType(containers[0].Id)
	}

	if state.ContainerId.ValueInt64() != 100 {
		t.Errorf("expected container_id=100 preserved, got %d", state.ContainerId.ValueInt64())
	}

	// server_uuid should be null since container 100 is not found
	uuid := readServerUUIDFromOwnedContainer(containers, state.ContainerId)
	if !uuid.IsNull() {
		t.Errorf("expected null server_uuid when container absent, got %q", uuid.ValueString())
	}
}
