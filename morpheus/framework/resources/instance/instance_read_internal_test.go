// (C) Copyright 2025-2026 Hewlett Packard Enterprise Development LP

package instance

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
)

// TestGetChildNetworksReadsSubnetID verifies that getChildNetworks reads subnet_id
// back from a child server interface's subnet association.
//
// When provisioning via subnet_id the Morpheus API resolves the subnet to its
// parent network (reported as network_id) but also returns the subnet itself
// (containerDetails.server.interfaces[].subnet, see _computeServerInterface.gson),
// so subnet_id round-trips on refresh rather than being write-only.
func TestGetChildNetworksReadsSubnetID(t *testing.T) {
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

// TestGetChildNetworksNoSubnet verifies that a child interface with no subnet
// association yields a null subnet_id (e.g. when network_id was used directly).
func TestGetChildNetworksNoSubnet(t *testing.T) {
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

// TestRemoveExternalStorageVolumes verifies that storage-server (SAN) volumes —
// e.g. Alletra MP BMaaS LUNs exported to the instance's host by
// hpe_morpheus_storage_volume — are excluded from the instance's tracked
// volumes, while the instance's own provisioned disks (no storageServer) are
// retained in order. This prevents spurious drift on the instance's volumes.
func TestRemoveExternalStorageVolumes(t *testing.T) {
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
