// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkconfig_test

import (
	"encoding/json"
	"os"
	"reflect"
	"testing"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/networkconfig"
)

// The live normaliser and the harvester must agree: a network type as the API
// returns it, normalised in Go, must equal the generated entry for the same
// type. The fixture is GET /api/network-types/{id} for ovsPortGroup on the
// appliance the table was harvested from, with the object ids removed.
func TestUnitFromAPIMatchesTheGeneratedTable(t *testing.T) {
	t.Parallel()

	raw, err := os.ReadFile("testdata/network_type_ovs_port_group.json")
	if err != nil {
		t.Fatal(err)
	}

	want, ok := networkconfig.Lookup("ovsPortGroup")
	if !ok {
		t.Fatal("the generated table must know ovsPortGroup")
	}

	// Through the SDK's show-response struct, as Create and ModifyPlan see it.
	var show sdk.GetNetworkType200Response
	if err := json.Unmarshal(raw, &show); err != nil {
		t.Fatal(err)
	}

	got, err := networkconfig.FromAPI(show.NetworkType)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("live normalisation differs from the generated entry\n got  %+v\n want %+v", got, want)
	}

	// Through the decoded JSON object, as the list endpoint's element would be.
	var wire struct {
		NetworkType map[string]any `json:"networkType"`
	}

	if err := json.Unmarshal(raw, &wire); err != nil {
		t.Fatal(err)
	}

	got, err = networkconfig.FromAPI(wire.NetworkType)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("live normalisation from JSON differs from the generated entry\n got  %+v\n want %+v", got, want)
	}
}

func TestUnitFromAPIDefaults(t *testing.T) {
	t.Parallel()

	// Absent creatable/editable are true, an absent context is the network
	// object, an absent input type is text; a type needs a code.
	spec, err := networkconfig.FromAPI(map[string]any{
		"code": "pluginType",
		"name": " Plugin Type ",
		"optionTypes": []any{
			map[string]any{"fieldName": "mystery", "displayOrder": 2},
			map[string]any{
				"fieldName": "hidden", "fieldContext": "config", "type": "hidden", "editable": false, "maxVal": 4094.0, "minVal": 0.5,
			},
			map[string]any{"fieldName": "", "fieldContext": "config"},
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	if spec.Name != "Plugin Type" || spec.Creatable || len(spec.Fields) != 2 {
		t.Fatalf("spec: %+v", spec)
	}

	mystery := spec.Fields[0]
	if mystery.Name != "mystery" || mystery.Context != networkconfig.ContextNetwork || mystery.Input != "text" ||
		!mystery.Creatable || !mystery.Editable || mystery.Required {
		t.Errorf("defaults: %+v", mystery)
	}

	hidden := spec.Fields[1]
	if hidden.Context != networkconfig.ContextConfig || hidden.Editable || hidden.Max != "4094" || hidden.Min != "0.5" {
		t.Errorf("hidden: %+v", hidden)
	}

	if _, err := networkconfig.FromAPI(map[string]any{"name": "no code"}); err == nil {
		t.Error("a type without a code must be rejected")
	}

	if _, err := networkconfig.FromAPI(nil); err == nil {
		t.Error("nil must be rejected")
	}
}
