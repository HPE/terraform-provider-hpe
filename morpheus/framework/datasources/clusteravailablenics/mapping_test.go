// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package clusteravailablenics

import (
	"context"
	"encoding/json"
	"testing"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
)

// availableNicsFixture is a hand-derived response matching the availableNics
// shape: one host with two NICs, one of them already consumed.
const availableNicsFixture = `{
  "availableNics": [
    {
      "serverId": 11,
      "serverName": "host-a",
      "nics": [
        {"name": "eno1", "macAddress": "00:11:22:33:44:55", "active": true, "consumed": false},
        {"name": "eno2", "macAddress": null, "active": false, "consumed": true}
      ]
    }
  ]
}`

func TestAvailableNicsAsState(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	var resp sdk.ListClusterAvailableNics200Response
	if err := json.Unmarshal([]byte(availableNicsFixture), &resp); err != nil {
		t.Fatalf("unmarshal fixture: %v", err)
	}

	list, err := availableNicsAsState(ctx, resp.AvailableNics)
	if err != nil {
		t.Fatalf("availableNicsAsState: %v", err)
	}

	hosts := list.Elements()
	if len(hosts) != 1 {
		t.Fatalf("available_nics len = %d, want 1", len(hosts))
	}

	host := hosts[0].(AvailableNicsValue)
	if got := host.ServerId.ValueInt64(); got != 11 {
		t.Errorf("server_id = %d, want 11", got)
	}

	if got := host.ServerName.ValueString(); got != "host-a" {
		t.Errorf("server_name = %q, want host-a", got)
	}

	nics := host.Nics.Elements()
	if len(nics) != 2 {
		t.Fatalf("nics len = %d, want 2", len(nics))
	}

	nic0 := nics[0].(NicsValue)
	if got := nic0.Name.ValueString(); got != "eno1" {
		t.Errorf("nic[0] name = %q, want eno1", got)
	}

	if nic0.Consumed.ValueBool() {
		t.Errorf("nic[0] consumed = true, want false")
	}

	nic1 := nics[1].(NicsValue)
	if !nic1.MacAddress.IsNull() {
		t.Errorf("nic[1] mac_address = %v, want null", nic1.MacAddress)
	}

	if !nic1.Consumed.ValueBool() {
		t.Errorf("nic[1] consumed = false, want true")
	}
}
