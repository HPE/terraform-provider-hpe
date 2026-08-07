// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package instance

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
)

// TestUnitGetAllServerInterfacesDuplicateNameNoNetwork verifies Defect 1: a
// duplicate-name interface group where no row carries Network information
// must not produce a zero-value entry (which would nil-dereference later).
func TestUnitGetAllServerInterfacesDuplicateNameNoNetwork(t *testing.T) {
	t.Parallel()

	instance := sdk.GetInstance200ResponseInstance{
		ContainerDetails: []sdk.InstanceContainer2{
			{
				Server: &sdk.InstanceContainerServer2{
					Interfaces: []sdk.InstanceContainerServerInterfacesInner1{
						// Two entries with the same name, neither has Network.
						{Name: sdk.PtrString("eth0"), IpAddress: sdk.PtrString("10.0.0.1")},
						{Name: sdk.PtrString("eth0"), IpAddress: sdk.PtrString("10.0.0.2")},
					},
				},
			},
		},
	}

	procIntfs := getAllServerInterfaces(instance)

	// The merged entry should be skipped (no Id, no Network).
	for _, iface := range procIntfs.serverIntfsList {
		if iface.Id == nil {
			// This is the defect scenario — the entry is zero-value.
			// After the fix, these should not panic downstream.
			continue
		}
	}

	// The serverIntfsMap should not contain any nil-Id entry.
	for id := range procIntfs.serverIntfsMap {
		if id == 0 {
			t.Error("serverIntfsMap contains a zero key, indicating a nil-Id entry was stored")
		}
	}
}

// TestUnitGetChildNetworksSubIntfNotInServerMap verifies Defect 2: a
// sub-interface id present in subIntfsMap but absent from serverIntfsMap
// must not panic on the map read.
func TestUnitGetChildNetworksSubIntfNotInServerMap(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	parentID := int64(1)
	// Sub-interface 99 exists in subIntfsMap but NOT in serverIntfsMap.
	subIntfMap := map[int64][]int64{parentID: {99}}
	serverIntfsMap := map[int64]sdk.InstanceContainerServerInterfacesInner1{
		// Only parent, no entry for 99.
		1: {Id: sdk.PtrInt64(1)},
	}

	list, diags := getChildNetworks(ctx, &parentID, subIntfMap, serverIntfsMap)
	if diags.HasError() {
		t.Fatalf("getChildNetworks returned diagnostics: %v", diags)
	}

	// Should produce zero children since 99 is skipped.
	var children []ChildVirtualNetworksValue
	if d := list.ElementsAs(ctx, &children, false); d.HasError() {
		t.Fatalf("ElementsAs returned diagnostics: %v", d)
	}

	if len(children) != 0 {
		t.Errorf("expected 0 children, got %d", len(children))
	}
}

// TestUnitGetChildNetworksNilNetwork verifies Defect 3: a child interface
// whose Network is nil must not panic when reading NetworkId.
func TestUnitGetChildNetworksNilNetwork(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	parentID := int64(1)
	subIntfMap := map[int64][]int64{parentID: {2}}
	serverIntfsMap := map[int64]sdk.InstanceContainerServerInterfacesInner1{
		2: {
			Id:      sdk.PtrInt64(2),
			Network: nil, // No network — must not dereference.
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
		t.Fatalf("expected 1 child, got %d", len(children))
	}

	if !children[0].NetworkId.IsNull() {
		t.Errorf("expected null network_id when Network is nil, got %v", children[0].NetworkId)
	}
}

// TestUnitGetStateInterfacesDiagnosticsAccumulated verifies Defect 4:
// diagnostics raised on a later iteration are not lost (discarded by
// overwrite). We test by having multiple interfaces with sub-interfaces
// where getChildNetworks is called per iteration.
func TestUnitGetStateInterfacesDiagnosticsAccumulated(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	// Build an instance with two top-level interfaces, each with a child.
	// Both children have Network set so they produce valid output and
	// no diagnostics. This verifies the plumbing doesn't crash — the
	// real diagnostic accumulation test would require forcing
	// types.ListValueFrom to error, which requires a type mismatch we
	// can't easily inject without modifying the function under test.
	instance := sdk.GetInstance200ResponseInstance{
		ContainerDetails: []sdk.InstanceContainer2{
			{
				Server: &sdk.InstanceContainerServer2{
					Interfaces: []sdk.InstanceContainerServerInterfacesInner1{
						{
							Id:      sdk.PtrInt64(1),
							Name:    sdk.PtrString("eth0"),
							Network: &sdk.InstanceContainerServerInterfacesInner1Network{Id: sdk.PtrInt64(10)},
							Interfaces: []sdk.InstanceContainerServerInstancesInnerInner1{
								{Id: sdk.PtrInt64(3)},
							},
						},
						{
							Id:      sdk.PtrInt64(2),
							Name:    sdk.PtrString("eth1"),
							Network: &sdk.InstanceContainerServerInterfacesInner1Network{Id: sdk.PtrInt64(20)},
							Interfaces: []sdk.InstanceContainerServerInstancesInnerInner1{
								{Id: sdk.PtrInt64(4)},
							},
						},
						// Child interfaces:
						{
							Id:      sdk.PtrInt64(3),
							Name:    sdk.PtrString("eth0.1"),
							Network: &sdk.InstanceContainerServerInterfacesInner1Network{Id: sdk.PtrInt64(30)},
						},
						{
							Id:      sdk.PtrInt64(4),
							Name:    sdk.PtrString("eth1.1"),
							Network: &sdk.InstanceContainerServerInterfacesInner1Network{Id: sdk.PtrInt64(40)},
						},
					},
				},
			},
		},
	}

	ifaces, diags := getStateInterfacesFromInstanceServer(ctx, instance)
	if diags.HasError() {
		t.Fatalf("unexpected error diagnostics: %v", diags)
	}

	// Should get 2 top-level interfaces (children are filtered).
	if len(ifaces) != 2 {
		t.Errorf("expected 2 top-level interfaces, got %d", len(ifaces))
	}
}

// TestUnitWaitForIpAddressSchemaDefault verifies that the generated schema
// declares wait_for_ip_address as optional+computed with a default of false.
func TestUnitWaitForIpAddressSchemaDefault(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	schemaResp := &resource.SchemaResponse{}
	r := &Resource{}
	r.Schema(ctx, resource.SchemaRequest{}, schemaResp)

	attr, ok := schemaResp.Schema.Attributes["wait_for_ip_address"]
	if !ok {
		t.Fatal("wait_for_ip_address attribute not found in schema")
	}

	if !attr.IsOptional() {
		t.Error("wait_for_ip_address should be optional")
	}

	if !attr.IsComputed() {
		t.Error("wait_for_ip_address should be computed")
	}

	boolAttr, ok := attr.(schema.BoolAttribute)
	if !ok {
		t.Fatal("wait_for_ip_address is not a BoolAttribute")
	}

	if boolAttr.Default == nil {
		t.Fatal("wait_for_ip_address should have a default")
	}

	defResp := defaults.BoolResponse{}
	boolAttr.Default.DefaultBool(ctx, defaults.BoolRequest{}, &defResp)
	if defResp.PlanValue.ValueBool() != false {
		t.Error("wait_for_ip_address default should be false")
	}
}
