// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkpoolserver

import (
	"encoding/json"
	"testing"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
)

func superset() *sdk.InfobloxNetworkPoolServer {
	s := &sdk.InfobloxNetworkPoolServer{Type: "x", Name: "n"}
	s.ServiceUrl = *sdk.NewNullableString(sdk.PtrString("https://ipam.example.com"))
	s.NetworkFilter = *sdk.NewNullableString(sdk.PtrString("10.0.0.0/8"))
	s.ZoneFilter = *sdk.NewNullableString(sdk.PtrString("example.com"))
	s.ServiceMode = sdk.PtrString("static")
	s.Config = &sdk.InfobloxNetworkPoolServerConfig{InventoryExisting: sdk.PtrString("on")}
	s.Credential = &sdk.InfobloxNetworkPoolServerCredential{Type: sdk.PtrString("7")}

	return s
}

func body(t *testing.T, v interface{}) map[string]interface{} {
	t.Helper()

	raw, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	var m map[string]interface{}
	if err := json.Unmarshal(raw, &m); err != nil {
		t.Fatal(err)
	}

	return m
}

// TestUnitCreateRequestByType checks that each type's body carries exactly
// its own fields: appId only for phpIPAM, Infoblox-only fields only for
// Infoblox (and the superset fallback), no networkFilter for SolarWinds.
func TestUnitCreateRequestByType(t *testing.T) {
	t.Parallel()

	app := sdk.PtrString("morpheus")
	cases := map[string]struct {
		code                              string
		appID                             *string
		wantAppID                         interface{}
		wantZoneFilter, wantNetworkFilter bool
	}{
		"phpipam":    {typeCodePhpIPAM, app, "morpheus", false, true},
		"bluecat":    {typeCodeBluecat, nil, nil, false, true},
		"solarwinds": {typeCodeSolarWinds, nil, nil, false, false},
		"infoblox":   {typeCodeInfoblox, nil, nil, true, true},
		"fallback":   {"solidserver", nil, nil, true, true},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			m := body(t, createRequest(tc.code, superset(), tc.appID))
			cfg, _ := m["config"].(map[string]interface{})

			if got := cfg["appId"]; got != tc.wantAppID {
				t.Errorf("config.appId = %v, want %v", got, tc.wantAppID)
			}
			if cfg["inventoryExisting"] != "on" {
				t.Errorf("config.inventoryExisting = %v, want on", cfg["inventoryExisting"])
			}
			if _, ok := m["zoneFilter"]; ok != tc.wantZoneFilter {
				t.Errorf("zoneFilter present = %v, want %v", ok, tc.wantZoneFilter)
			}
			if _, ok := m["networkFilter"]; ok != tc.wantNetworkFilter {
				t.Errorf("networkFilter present = %v, want %v", ok, tc.wantNetworkFilter)
			}
			if m["name"] != "n" || m["serviceUrl"] != "https://ipam.example.com" {
				t.Errorf("common fields not carried: %v", m)
			}
			if cred, _ := m["credential"].(map[string]interface{}); cred["type"] != "7" {
				t.Errorf("credential = %v", m["credential"])
			}
		})
	}
}

// TestUnitUpdateRequestByType mirrors the create check for updates, including
// that a phpIPAM update without app_id or inventory_existing sends no config.
func TestUnitUpdateRequestByType(t *testing.T) {
	t.Parallel()

	s := &sdk.InfobloxNetworkPoolServerUpdate{Name: sdk.PtrString("n")}
	s.ZoneFilter = *sdk.NewNullableString(sdk.PtrString("example.com"))
	s.NetworkFilter = *sdk.NewNullableString(sdk.PtrString("10.0.0.0/8"))

	php := body(t, updateRequest(typeCodePhpIPAM, s, sdk.PtrString("morpheus")))
	if cfg, _ := php["config"].(map[string]interface{}); cfg["appId"] != "morpheus" {
		t.Errorf("phpipam config = %v", php["config"])
	}
	if _, ok := php["zoneFilter"]; ok {
		t.Errorf("phpipam update carries zoneFilter")
	}

	phpNoCfg := body(t, updateRequest(typeCodePhpIPAM, s, nil))
	if _, ok := phpNoCfg["config"]; ok {
		t.Errorf("phpipam update without app_id carries config: %v", phpNoCfg["config"])
	}

	sw := body(t, updateRequest(typeCodeSolarWinds, s, nil))
	if _, ok := sw["networkFilter"]; ok {
		t.Errorf("solarwinds update carries networkFilter")
	}

	ib := body(t, updateRequest(typeCodeInfoblox, s, nil))
	if ib["zoneFilter"] != "example.com" {
		t.Errorf("infoblox update dropped zoneFilter: %v", ib)
	}
}
