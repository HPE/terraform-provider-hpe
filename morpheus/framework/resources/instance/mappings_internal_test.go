// (C) Copyright 2025-2026 Hewlett Packard Enterprise Development LP

package instance

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

// TestNetworkInterfaceMapperNetworkID verifies the network.id string sent to the
// API is built with the correct prefix for each of network_id, network_group_id
// and subnet_id.
//
// The "subnet-" prefix is required for subnet based provisioning (e.g. Azure):
// InstanceService.parseNetworks in morpheus-ui dispatches network.id by prefix
// ("network-", "subnet-", "networkGroup-", or a bare numeric id). The provider
// owns the prefix and accepts a plain integer subnet id from the user.
func TestNetworkInterfaceMapperNetworkID(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	childNull := types.ListNull(ChildVirtualNetworksValue{}.Type(ctx))

	tests := []struct {
		name           string
		networkID      types.Int64
		networkGroupID types.Int64
		subnetID       types.Int64
		want           string
	}{
		{
			name:           "network_id maps to a bare id",
			networkID:      types.Int64Value(28),
			networkGroupID: types.Int64Null(),
			subnetID:       types.Int64Null(),
			want:           "28",
		},
		{
			name:           "network_group_id maps to a networkGroup- prefix",
			networkID:      types.Int64Null(),
			networkGroupID: types.Int64Value(2),
			subnetID:       types.Int64Null(),
			want:           "networkGroup-2",
		},
		{
			name:           "subnet_id maps to a subnet- prefix",
			networkID:      types.Int64Null(),
			networkGroupID: types.Int64Null(),
			subnetID:       types.Int64Value(5),
			want:           "subnet-5",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			in := NetworkInterfacesValue{
				Id:                   types.Int64Null(),
				NetworkId:            tt.networkID,
				NetworkGroupId:       tt.networkGroupID,
				SubnetId:             tt.subnetID,
				IpPool:               types.Int64Null(),
				IpMode:               types.StringNull(),
				IpAddress:            types.StringNull(),
				NetworkTypeId:        types.Int64Null(),
				ChildVirtualNetworks: childNull,
			}

			if got := createNetworkInterfaceMapper(ctx)(in).Network.Id; got != tt.want {
				t.Errorf("create mapper network.id = %q, want %q", got, tt.want)
			}

			if got := updateNetworkInterfaceMapper(ctx)(in).Network.Id; got != tt.want {
				t.Errorf("update mapper network.id = %q, want %q", got, tt.want)
			}
		})
	}
}

// TestChildNetworkInterfaceMapperNetworkID verifies the network.id string sent to
// the API for child virtual networks is built with the correct prefix for each of
// network_id, network_group_id and subnet_id.
//
// Child virtual networks are parsed by the same (recursive) InstanceService.parseNetworks
// path in morpheus-ui as top-level interfaces, so subnet based provisioning (the
// "subnet-" prefix) applies equally to them.
func TestChildNetworkInterfaceMapperNetworkID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		networkID      types.Int64
		networkGroupID types.Int64
		subnetID       types.Int64
		want           string
	}{
		{
			name:           "network_id maps to a bare id",
			networkID:      types.Int64Value(28),
			networkGroupID: types.Int64Null(),
			subnetID:       types.Int64Null(),
			want:           "28",
		},
		{
			name:           "network_group_id maps to a networkGroup- prefix",
			networkID:      types.Int64Null(),
			networkGroupID: types.Int64Value(2),
			subnetID:       types.Int64Null(),
			want:           "networkGroup-2",
		},
		{
			name:           "subnet_id maps to a subnet- prefix",
			networkID:      types.Int64Null(),
			networkGroupID: types.Int64Null(),
			subnetID:       types.Int64Value(5),
			want:           "subnet-5",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			in := ChildVirtualNetworksValue{
				Id:             types.Int64Null(),
				NetworkId:      tt.networkID,
				NetworkGroupId: tt.networkGroupID,
				SubnetId:       tt.subnetID,
				IpPool:         types.Int64Null(),
				IpMode:         types.StringNull(),
				IpAddress:      types.StringNull(),
				NetworkTypeId:  types.Int64Null(),
			}

			if got := createChildNetworkInterfaceMapper(in).Network.Id; got != tt.want {
				t.Errorf("create child mapper network.id = %q, want %q", got, tt.want)
			}

			if got := updateChildNetworkInterfaceMapper(in).Network.Id; got != tt.want {
				t.Errorf("update child mapper network.id = %q, want %q", got, tt.want)
			}
		})
	}
}
