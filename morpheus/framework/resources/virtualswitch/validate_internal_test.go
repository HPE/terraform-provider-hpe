// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package virtualswitch

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func i32(v int32) *int32  { return &v }
func sp(v string) *string { return &v }

func TestValidateSegmentsAndUplinks(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		baseType types.String
		bondMode types.String
		segs     []segmentInput
		ups      []uplinkInput
		wantErr  bool
	}{
		{
			name:     "valid mgmt+vm+nfs",
			bondMode: types.StringValue("active-backup"),
			segs: []segmentInput{
				{key: "mgmt", trafficType: trafficManagement, ipAddress: sp("10.0.0.5"), netmask: sp("255.255.255.0")},
				{key: "vm", trafficType: trafficVM},
				{key: "nfs", trafficType: trafficDataNfs, vlanID: i32(100), netmask: sp("255.255.255.0"), ipAddress: sp("10.0.0.1")},
			},
			ups:     []uplinkInput{{serverID: 1, nic1: "eth2", nic2: sp("eth3")}},
			wantErr: false,
		},
		{
			name:     "iscsi forbids bond",
			bondMode: types.StringValue("802.3ad"),
			segs:     []segmentInput{{key: "i", trafficType: trafficDataIscsi}},
			ups:      []uplinkInput{{serverID: 1, nic1: "eth2"}},
			wantErr:  true,
		},
		{
			name:     "iscsi forbids nic2",
			bondMode: types.StringNull(),
			segs:     []segmentInput{{key: "i", trafficType: trafficDataIscsi}},
			ups:      []uplinkInput{{serverID: 1, nic1: "eth2", nic2: sp("eth3")}},
			wantErr:  true,
		},
		{
			name:     "iscsi exclusive with nfs",
			bondMode: types.StringNull(),
			segs: []segmentInput{
				{key: "i", trafficType: trafficDataIscsi},
				{key: "n", trafficType: trafficDataNfs},
			},
			ups:     []uplinkInput{{serverID: 1, nic1: "eth2"}},
			wantErr: true,
		},
		{
			name:     "sdn exclusive with nfs",
			bondMode: types.StringNull(),
			segs: []segmentInput{
				{key: "s", trafficType: trafficSDN},
				{key: "n", trafficType: trafficDataNfs},
			},
			ups:     []uplinkInput{{serverID: 1, nic1: "eth2"}},
			wantErr: true,
		},
		{
			name:     "live-migration on iscsi base",
			baseType: types.StringValue("iscsi"),
			bondMode: types.StringNull(),
			segs:     []segmentInput{{key: "lm", trafficType: trafficLiveMigration}},
			ups:      []uplinkInput{{serverID: 1, nic1: "eth2"}},
			wantErr:  true,
		},
		{
			name:     "live-migration with data-iscsi",
			bondMode: types.StringNull(),
			segs: []segmentInput{
				{key: "lm", trafficType: trafficLiveMigration},
				{key: "i", trafficType: trafficDataIscsi},
			},
			ups:     []uplinkInput{{serverID: 1, nic1: "eth2"}},
			wantErr: true,
		},
		{
			name:     "duplicate nfs vlan",
			bondMode: types.StringNull(),
			segs: []segmentInput{
				{key: "n1", trafficType: trafficDataNfs, vlanID: i32(100)},
				{key: "n2", trafficType: trafficDataNfs, vlanID: i32(100)},
			},
			ups:     []uplinkInput{{serverID: 1, nic1: "eth2"}},
			wantErr: true,
		},
		{
			name:     "two untagged nfs",
			bondMode: types.StringNull(),
			segs: []segmentInput{
				{key: "n1", trafficType: trafficDataNfs},
				{key: "n2", trafficType: trafficDataNfs},
			},
			ups:     []uplinkInput{{serverID: 1, nic1: "eth2"}},
			wantErr: true,
		},
		{
			name:     "ip without netmask",
			bondMode: types.StringNull(),
			segs:     []segmentInput{{key: "n", trafficType: trafficDataNfs, ipAddress: sp("10.0.0.1")}},
			ups:      []uplinkInput{{serverID: 1, nic1: "eth2"}},
			wantErr:  true,
		},
		{
			name:     "host_ips without netmask",
			bondMode: types.StringNull(),
			segs:     []segmentInput{{key: "n", trafficType: trafficDataNfs, hostIps: []hostIPInput{{serverID: 1, ip: "10.0.0.1"}}}},
			ups:      []uplinkInput{{serverID: 1, nic1: "eth2"}},
			wantErr:  true,
		},
		{
			name:     "two management segments",
			bondMode: types.StringNull(),
			segs: []segmentInput{
				{key: "m1", trafficType: trafficManagement},
				{key: "m2", trafficType: trafficManagement},
			},
			ups:     []uplinkInput{{serverID: 1, nic1: "eth2"}},
			wantErr: true,
		},
		{
			name:     "management without any address",
			bondMode: types.StringNull(),
			segs:     []segmentInput{{key: "mgmt", trafficType: trafficManagement}},
			ups:      []uplinkInput{{serverID: 1, nic1: "eth2"}},
			wantErr:  true,
		},
		{
			name:     "management with ip_address and netmask",
			bondMode: types.StringNull(),
			segs: []segmentInput{
				{key: "mgmt", trafficType: trafficManagement, ipAddress: sp("10.0.0.5"), netmask: sp("255.255.255.0")},
			},
			ups:     []uplinkInput{{serverID: 1, nic1: "eth2"}},
			wantErr: false,
		},
		{
			name:     "management with host_ips and netmask",
			bondMode: types.StringNull(),
			segs: []segmentInput{
				{
					key: "mgmt", trafficType: trafficManagement, netmask: sp("255.255.255.0"),
					hostIps: []hostIPInput{{serverID: 1, ip: "10.0.0.5"}},
				},
			},
			ups:     []uplinkInput{{serverID: 1, nic1: "eth2"}},
			wantErr: false,
		},
		// Unknown values must not be read as unset: ValidateConfig runs before
		// values from other resources or data sources are known.
		{
			name:     "management address not yet known",
			bondMode: types.StringNull(),
			segs: []segmentInput{
				{key: "mgmt", trafficType: trafficManagement, addressUnknown: true, netmaskUnknown: true},
			},
			ups:     []uplinkInput{{serverID: 1, nic1: "eth2"}},
			wantErr: false,
		},
		{
			name:     "netmask not yet known with an address",
			bondMode: types.StringNull(),
			segs: []segmentInput{
				{key: "nfs", trafficType: trafficDataNfs, vlanID: i32(5), ipAddress: sp("10.0.0.1"), netmaskUnknown: true},
			},
			ups:     []uplinkInput{{serverID: 1, nic1: "eth2"}},
			wantErr: false,
		},
		{
			name:     "two vlan_ids not yet known are not two untagged segments",
			bondMode: types.StringNull(),
			segs: []segmentInput{
				{key: "a", trafficType: trafficDataNfs, vlanUnknown: true},
				{key: "b", trafficType: trafficDataNfs, vlanUnknown: true},
			},
			ups:     []uplinkInput{{serverID: 1, nic1: "eth2"}},
			wantErr: false,
		},
		{
			name:     "two untagged segments are still refused",
			bondMode: types.StringNull(),
			segs: []segmentInput{
				{key: "a", trafficType: trafficDataNfs},
				{key: "b", trafficType: trafficDataNfs},
			},
			ups:     []uplinkInput{{serverID: 1, nic1: "eth2"}},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			diags := validateSegmentsAndUplinks(tt.baseType, tt.bondMode, tt.segs, tt.ups)
			if diags.HasError() != tt.wantErr {
				t.Fatalf("want err=%v, got diags=%v", tt.wantErr, diags)
			}
		})
	}
}
