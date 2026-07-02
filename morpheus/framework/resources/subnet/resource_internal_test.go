// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package subnet

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
)

// TestMapResponseToModel_ComputedFieldsKnownAfterApply is a regression test for
// MORPH-13600. When a subnet is created with only network_id/type_id/name and no
// CIDR, the API returns nil for cidr, netmask and subnet_address. Those
// attributes are Computed with UseStateForUnknown, so on create their planned
// value is unknown. mapResponseToModel must resolve them to a known (null) value
// rather than leaving the unknown plan value in place, otherwise Terraform fails
// with "Provider returned invalid result object after apply ... still indicated
// an unknown value".
func TestMapResponseToModel_ComputedFieldsKnownAfterApply(t *testing.T) {
	t.Parallel()

	// Simulate the create plan: computed attributes start out unknown.
	model := &SubnetModel{
		Name:          types.StringUnknown(),
		Cidr:          types.StringUnknown(),
		Netmask:       types.StringUnknown(),
		SubnetAddress: types.StringUnknown(),
		Active:        types.BoolUnknown(),
		DhcpServer:    types.BoolUnknown(),
	}

	// Simulate the API response for a subnet created without a CIDR: the derived
	// fields are omitted (nil).
	subnet := &sdk.GetSubnet200ResponseSubnet{
		Id:   sdk.PtrInt64(42),
		Name: sdk.PtrString("qatf-basic-subnet"),
		Type: &sdk.GetSubnet200ResponseSubnetType{Id: sdk.PtrInt64(8)},
		Network: &sdk.GetSubnet200ResponseSubnetNetwork{
			Id: sdk.PtrInt64(1),
		},
		// Cidr, Netmask, SubnetAddress, Active, DhcpServer intentionally nil.
	}

	mapResponseToModel(model, subnet)

	// Every Computed attribute must be known after apply. Null is fine; unknown
	// is the bug.
	assertKnown := func(name string, unknown bool) {
		t.Helper()
		if unknown {
			t.Errorf("%s is still unknown after mapResponseToModel; "+
				"Computed attributes must be known after apply", name)
		}
	}
	assertKnown("cidr", model.Cidr.IsUnknown())
	assertKnown("netmask", model.Netmask.IsUnknown())
	assertKnown("subnet_address", model.SubnetAddress.IsUnknown())
	assertKnown("name", model.Name.IsUnknown())
	assertKnown("active", model.Active.IsUnknown())
	assertKnown("dhcp_server", model.DhcpServer.IsUnknown())

	// The nil API fields specifically should become null.
	if !model.Cidr.IsNull() {
		t.Errorf("cidr: expected null for nil API value, got %q", model.Cidr.ValueString())
	}
	if !model.Netmask.IsNull() {
		t.Errorf("netmask: expected null for nil API value, got %q", model.Netmask.ValueString())
	}
	if !model.SubnetAddress.IsNull() {
		t.Errorf("subnet_address: expected null for nil API value, got %q", model.SubnetAddress.ValueString())
	}
}

// TestMapResponseToModel_ComputedFieldsPopulated verifies the happy path: when
// the API returns values for the derived fields (e.g. an Azure subnet created
// from config.subnetCidr), they are mapped through into state.
func TestMapResponseToModel_ComputedFieldsPopulated(t *testing.T) {
	t.Parallel()

	model := &SubnetModel{
		Cidr:          types.StringUnknown(),
		Netmask:       types.StringUnknown(),
		SubnetAddress: types.StringUnknown(),
	}

	subnet := &sdk.GetSubnet200ResponseSubnet{
		Id:            sdk.PtrInt64(7),
		Name:          sdk.PtrString("azure-subnet"),
		Cidr:          sdk.PtrString("10.0.250.0/24"),
		Netmask:       sdk.PtrString("255.255.255.0"),
		SubnetAddress: sdk.PtrString("10.0.250.0"),
		Active:        sdk.PtrBool(true),
		DhcpServer:    sdk.PtrBool(false),
		Type:          &sdk.GetSubnet200ResponseSubnetType{Id: sdk.PtrInt64(12)},
		Network:       &sdk.GetSubnet200ResponseSubnetNetwork{Id: sdk.PtrInt64(88)},
	}

	mapResponseToModel(model, subnet)

	if got := model.Cidr.ValueString(); got != "10.0.250.0/24" {
		t.Errorf("cidr: expected %q, got %q", "10.0.250.0/24", got)
	}
	if got := model.Netmask.ValueString(); got != "255.255.255.0" {
		t.Errorf("netmask: expected %q, got %q", "255.255.255.0", got)
	}
	if got := model.SubnetAddress.ValueString(); got != "10.0.250.0" {
		t.Errorf("subnet_address: expected %q, got %q", "10.0.250.0", got)
	}
}
