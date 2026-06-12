// (C) Copyright 2025-2026 Hewlett Packard Enterprise Development LP

package instance

import (
	"context"
	"testing"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
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
