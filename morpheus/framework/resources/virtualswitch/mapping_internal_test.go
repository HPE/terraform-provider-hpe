// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package virtualswitch

import (
	"context"
	"encoding/json"
	"testing"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
)

// fixture mirrors the controller's switchDetailModel / segmentModel /
// uplinkModel shape (rc-9.1.0).
const detailFixture = `{
  "id": 42,
  "name": "vsw-hvm",
  "code": null,
  "description": "primary switch",
  "baseType": "general",
  "bondMode": "active-backup",
  "mtu": 9000,
  "status": "available",
  "statusMessage": null,
  "active": true,
  "clusterId": 3,
  "nicCount": 2,
  "networkCount": 3,
  "segments": [
    {
      "id": 100,
      "trafficType": "management",
      "interfaceName": "bond0",
      "ipAddress": null,
      "netmask": null,
      "gateway": null,
      "vlanId": null,
      "mtu": null,
      "status": "available"
    },
    {
      "id": 101,
      "trafficType": "data-nfs",
      "interfaceName": "bond0.100",
      "ipAddress": null,
      "netmask": "255.255.255.0",
      "gateway": null,
      "vlanId": 100,
      "mtu": 9000,
      "ovsBridge": "ovsbr0",
      "status": "available",
      "hostIps": [
        {"serverId": 7, "ip": "10.20.0.7"},
        {"serverId": 8, "ip": "10.20.0.8"}
      ]
    }
  ],
  "uplinks": [
    {"id": 200, "serverId": 7, "serverName": "host-7", "nic1": "eth2", "nic2": "eth3", "syncStatus": "synced"},
    {"id": 201, "serverId": 8, "serverName": "host-8", "nic1": "eth2", "nic2": "eth3", "syncStatus": "synced"}
  ]
}`

func TestMapDetailToState(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	var vs sdk.GetVirtualSwitch200ResponseVirtualSwitch
	if err := json.Unmarshal([]byte(detailFixture), &vs); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	state, diags := mapDetailToState(ctx, &vs, VirtualSwitchModel{})
	if diags.HasError() {
		t.Fatalf("diags: %v", diags)
	}

	if state.Id.ValueInt64() != 42 {
		t.Fatalf("id: %d", state.Id.ValueInt64())
	}

	if state.Name.ValueString() != "vsw-hvm" {
		t.Fatalf("name: %s", state.Name.ValueString())
	}

	if state.BondMode.ValueString() != "active-backup" {
		t.Fatalf("bond_mode: %s", state.BondMode.ValueString())
	}

	if state.Mtu.ValueInt64() != 9000 {
		t.Fatalf("mtu: %d", state.Mtu.ValueInt64())
	}

	if state.ClusterId.ValueInt64() != 3 {
		t.Fatalf("cluster_id: %d", state.ClusterId.ValueInt64())
	}

	segsMap := make(map[string]SegmentsValue)
	if d := state.Segments.ElementsAs(ctx, &segsMap, false); d.HasError() {
		t.Fatalf("segments: %v", d)
	}

	if len(segsMap) != 2 {
		t.Fatalf("segment count: %d", len(segsMap))
	}

	// keys synthesised: "management" and "data-nfs_100".
	nfs, ok := segsMap["data-nfs_100"]
	if !ok {
		t.Fatalf("missing synthesised key, got keys %v", segsMap)
	}

	if nfs.VlanId.ValueInt64() != 100 {
		t.Fatalf("nfs vlan: %d", nfs.VlanId.ValueInt64())
	}

	if nfs.Netmask.ValueString() != "255.255.255.0" {
		t.Fatalf("nfs netmask: %s", nfs.Netmask.ValueString())
	}

	hostIps := make([]HostIpsValue, 0)
	if d := nfs.HostIps.ElementsAs(ctx, &hostIps, false); d.HasError() {
		t.Fatalf("host_ips: %v", d)
	}

	if len(hostIps) != 2 {
		t.Fatalf("host ip count: %d", len(hostIps))
	}

	ups := make([]UplinksValue, 0)
	if d := state.Uplinks.ElementsAs(ctx, &ups, false); d.HasError() {
		t.Fatalf("uplinks: %v", d)
	}

	if len(ups) != 2 {
		t.Fatalf("uplink count: %d", len(ups))
	}
}
