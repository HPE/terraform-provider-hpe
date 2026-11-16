// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package virtualswitches

import (
	"context"
	"encoding/json"
	"testing"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
)

// virtualSwitchesListFixture is a hand-derived list response matching the
// summary switchModel shape: segments carry only trafficType.
const virtualSwitchesListFixture = `{
  "virtualSwitches": [
    {
      "id": 42,
      "name": "vswitch01",
      "code": "vswitch01",
      "description": "management switch",
      "baseType": "general",
      "bondMode": null,
      "mtu": 9000,
      "status": "ok",
      "statusMessage": null,
      "active": true,
      "clusterId": 7,
      "nicCount": 2,
      "networkCount": 2,
      "uplinkNics": ["eno1", "eno2"],
      "segments": [
        {"trafficType": "management"},
        {"trafficType": "vm"}
      ]
    }
  ]
}`

func TestVirtualSwitchesAsState(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	var resp sdk.ListVirtualSwitches200Response
	if err := json.Unmarshal([]byte(virtualSwitchesListFixture), &resp); err != nil {
		t.Fatalf("unmarshal fixture: %v", err)
	}

	list, err := virtualSwitchesAsState(ctx, resp.VirtualSwitches)
	if err != nil {
		t.Fatalf("virtualSwitchesAsState: %v", err)
	}

	elems := list.Elements()
	if len(elems) != 1 {
		t.Fatalf("virtual_switches len = %d, want 1", len(elems))
	}

	sw := elems[0].(VirtualSwitchesValue)
	if got := sw.Id.ValueInt64(); got != 42 {
		t.Errorf("id = %d, want 42", got)
	}

	if got := sw.Mtu.ValueInt64(); got != 9000 {
		t.Errorf("mtu = %d, want 9000", got)
	}

	if !sw.BondMode.IsNull() {
		t.Errorf("bond_mode = %v, want null", sw.BondMode)
	}

	segs := sw.Segments.Elements()
	if len(segs) != 2 {
		t.Fatalf("segments len = %d, want 2", len(segs))
	}

	seg := segs[1].(SegmentsValue)
	if got := seg.TrafficType.ValueString(); got != "vm" {
		t.Errorf("segment[1] traffic_type = %q, want vm", got)
	}
}
