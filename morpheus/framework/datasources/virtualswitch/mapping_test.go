// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package virtualswitch

import (
	"context"
	"encoding/json"
	"testing"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
)

// virtualSwitchDetailFixture is a hand-derived response body matching the
// switchDetailModel shape: a switch with a segment carrying hostIps and an
// uplink whose nic2 is null (single-NIC mode).
const virtualSwitchDetailFixture = `{
  "virtualSwitch": {
    "id": 42,
    "name": "vswitch01",
    "code": "vswitch01",
    "description": "management switch",
    "baseType": "general",
    "bondMode": "active-backup",
    "mtu": 1500,
    "status": "ok",
    "statusMessage": null,
    "active": true,
    "clusterId": 7,
    "nicCount": 2,
    "networkCount": 1,
    "uplinkNics": ["eno1", "eno2"],
    "segments": [
      {
        "id": 100,
        "trafficType": "management",
        "interfaceName": "bond0",
        "ipAddress": "10.0.0.5",
        "netmask": "255.255.255.0",
        "gateway": "10.0.0.1",
        "vlanId": null,
        "mtu": 1500,
        "ovsBridge": null,
        "ovsUuid": null,
        "ovsLibvirtNetwork": null,
        "status": "ok",
        "statusMessage": null,
        "hostIps": [
          {"serverId": 11, "ip": "10.0.0.5"},
          {"serverId": 12, "ip": "10.0.0.6"}
        ]
      }
    ],
    "uplinks": [
      {
        "id": 200,
        "serverId": 11,
        "serverName": "host-a",
        "nic1": "eno1",
        "nic2": null,
        "syncStatus": "ok",
        "syncMessage": null,
        "lastSyncDate": "2026-11-12T10:15:30Z"
      }
    ]
  }
}`

func TestVirtualSwitchAsState(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	var resp sdk.GetVirtualSwitch200Response
	if err := json.Unmarshal([]byte(virtualSwitchDetailFixture), &resp); err != nil {
		t.Fatalf("unmarshal fixture: %v", err)
	}

	state, err := virtualSwitchAsState(ctx, resp.VirtualSwitch)
	if err != nil {
		t.Fatalf("virtualSwitchAsState: %v", err)
	}

	if got := state.Id.ValueInt64(); got != 42 {
		t.Errorf("id = %d, want 42", got)
	}

	if got := state.Name.ValueString(); got != "vswitch01" {
		t.Errorf("name = %q, want vswitch01", got)
	}

	if got := state.BondMode.ValueString(); got != "active-backup" {
		t.Errorf("bond_mode = %q, want active-backup", got)
	}

	if got := state.ClusterId.ValueInt64(); got != 7 {
		t.Errorf("cluster_id = %d, want 7", got)
	}

	if !state.StatusMessage.IsNull() {
		t.Errorf("status_message = %v, want null", state.StatusMessage)
	}

	segs := state.Segments.Elements()
	if len(segs) != 1 {
		t.Fatalf("segments len = %d, want 1", len(segs))
	}

	seg := segs[0].(SegmentsValue)
	if got := seg.TrafficType.ValueString(); got != "management" {
		t.Errorf("segment traffic_type = %q, want management", got)
	}

	if !seg.VlanId.IsNull() {
		t.Errorf("segment vlan_id = %v, want null", seg.VlanId)
	}

	hostIps := seg.HostIps.Elements()
	if len(hostIps) != 2 {
		t.Fatalf("host_ips len = %d, want 2", len(hostIps))
	}

	hip := hostIps[0].(HostIpsValue)
	if got := hip.ServerId.ValueInt64(); got != 11 {
		t.Errorf("host_ip server_id = %d, want 11", got)
	}

	if got := hip.Ip.ValueString(); got != "10.0.0.5" {
		t.Errorf("host_ip ip = %q, want 10.0.0.5", got)
	}

	ups := state.Uplinks.Elements()
	if len(ups) != 1 {
		t.Fatalf("uplinks len = %d, want 1", len(ups))
	}

	up := ups[0].(UplinksValue)
	if got := up.Nic1.ValueString(); got != "eno1" {
		t.Errorf("uplink nic1 = %q, want eno1", got)
	}

	if !up.Nic2.IsNull() {
		t.Errorf("uplink nic2 = %v, want null", up.Nic2)
	}

	if got := up.LastSyncDate.ValueString(); got != "2026-11-12T10:15:30Z" {
		t.Errorf("uplink last_sync_date = %q, want the RFC 3339 timestamp", got)
	}

	if got := up.ServerName.ValueString(); got != "host-a" {
		t.Errorf("uplink server_name = %q, want host-a", got)
	}
}
