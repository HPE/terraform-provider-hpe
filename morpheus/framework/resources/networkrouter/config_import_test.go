// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkrouter

import (
	"context"
	"testing"
)

// TestUnitTier1ConfigFromMap mirrors the MORPH-12175 import scenario: a Tier1
// router applied with ip_management_type = dhcpLocal and all tier1_* checkboxes
// false round-trips from the API config map into the typed object.
func TestUnitTier1ConfigFromMap(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	m := map[string]interface{}{
		"ipManagementType":           "dhcpLocal",
		"tier0Gateway":               nil,
		"edgeCluster":                nil,
		"failOver":                   nil,
		"ipServerId":                 nil,
		"TIER1_CONNECTED":            "off",
		"TIER1_NAT":                  "off",
		"TIER1_STATIC_ROUTES":        "off",
		"TIER1_LB_VIP":               "off",
		"TIER1_LB_SNAT":              "off",
		"TIER1_DNS_FORWARDER_IP":     "off",
		"TIER1_IPSEC_LOCAL_ENDPOINT": "off",
	}

	got, diags := tier1ConfigFromMap(ctx, m)
	if diags.HasError() {
		t.Fatalf("unexpected diags: %v", diags)
	}

	if got.IpManagementType.ValueString() != "dhcpLocal" {
		t.Errorf("ip_management_type = %q, want dhcpLocal", got.IpManagementType.ValueString())
	}

	bools := map[string]bool{
		"tier1_connected":            got.Tier1Connected.ValueBool(),
		"tier1_nat":                  got.Tier1Nat.ValueBool(),
		"tier1_static_routes":        got.Tier1StaticRoutes.ValueBool(),
		"tier1_lb_vip":               got.Tier1LbVip.ValueBool(),
		"tier1_lb_snat":              got.Tier1LbSnat.ValueBool(),
		"tier1_dns_forwarder_ip":     got.Tier1DnsForwarderIp.ValueBool(),
		"tier1_ipsec_local_endpoint": got.Tier1IpsecLocalEndpoint.ValueBool(),
	}
	for name, v := range bools {
		if v {
			t.Errorf("%s = true, want false", name)
		}
	}

	// Unset string fields must be null so import matches a freshly-applied
	// resource.
	if !got.EdgeCluster.IsNull() || !got.FailOver.IsNull() ||
		!got.IpServerId.IsNull() || !got.Tier0Gateway.IsNull() {
		t.Errorf("expected unset string fields to be null")
	}
}

// TestUnitTier1ConfigFromMapOnOff verifies "on" maps to true and that a checkbox
// absent from the API map defaults to false (matching the computed schema
// default), rather than becoming null.
func TestUnitTier1ConfigFromMapOnOff(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	m := map[string]interface{}{
		"TIER1_CONNECTED": "on",
		"TIER1_NAT":       "on",
	}

	got, diags := tier1ConfigFromMap(ctx, m)
	if diags.HasError() {
		t.Fatalf("unexpected diags: %v", diags)
	}

	if !got.Tier1Connected.ValueBool() {
		t.Errorf("tier1_connected = false, want true")
	}
	if !got.Tier1Nat.ValueBool() {
		t.Errorf("tier1_nat = false, want true")
	}
	if got.Tier1LbVip.IsNull() || got.Tier1LbVip.ValueBool() {
		t.Errorf("tier1_lb_vip = %v, want known false (absent -> false)", got.Tier1LbVip)
	}
}

// TestUnitTier0ConfigFromMap exercises the string, numeric, and checkbox paths of
// the Tier0 reverse-mapper.
func TestUnitTier0ConfigFromMap(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	m := map[string]interface{}{
		"haMode":       "ACTIVE_ACTIVE",
		"LOCAL_AS_NUM": "65000",
		"RESTART_TIME": float64(120),
		"ECMP":         "on",
		"TIER0_STATIC": "off",
	}

	got, diags := tier0ConfigFromMap(ctx, m)
	if diags.HasError() {
		t.Fatalf("unexpected diags: %v", diags)
	}

	if got.HaMode.ValueString() != "ACTIVE_ACTIVE" {
		t.Errorf("ha_mode = %q, want ACTIVE_ACTIVE", got.HaMode.ValueString())
	}
	if got.LocalAsNum.ValueString() != "65000" {
		t.Errorf("local_as_num = %q, want 65000", got.LocalAsNum.ValueString())
	}
	if got.RestartTime.ValueInt64() != 120 {
		t.Errorf("restart_time = %d, want 120", got.RestartTime.ValueInt64())
	}
	if !got.Ecmp.ValueBool() {
		t.Errorf("ecmp = false, want true")
	}
	if got.Tier0Static.ValueBool() {
		t.Errorf("tier0_static = true, want false")
	}
	// An absent numeric field is null (not zero).
	if !got.StaleRouteTime.IsNull() {
		t.Errorf("stale_route_time = %v, want null when absent", got.StaleRouteTime)
	}
}

// TestUnitConfigHelpers checks the low-level extraction helpers, including the
// absent/null -> false default for checkboxes.
func TestUnitConfigHelpers(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	m := map[string]interface{}{
		"str": "value",
		"num": float64(7),
		"on":  "on",
		"off": "off",
		"nil": nil,
	}

	if got := configString(m, "str"); got.ValueString() != "value" {
		t.Errorf("configString(str) = %q, want value", got.ValueString())
	}
	if got := configString(m, "missing"); !got.IsNull() {
		t.Errorf("configString(missing) should be null")
	}
	if got := configInt64(m, "num"); got.ValueInt64() != 7 {
		t.Errorf("configInt64(num) = %d, want 7", got.ValueInt64())
	}
	if got := configInt64(m, "missing"); !got.IsNull() {
		t.Errorf("configInt64(missing) should be null")
	}
	if got := configBoolOnOff(ctx, m, "on"); !got.ValueBool() {
		t.Errorf("configBoolOnOff(on) = false, want true")
	}
	if got := configBoolOnOff(ctx, m, "off"); got.ValueBool() {
		t.Errorf("configBoolOnOff(off) = true, want false")
	}
	if got := configBoolOnOff(ctx, m, "nil"); got.IsNull() || got.ValueBool() {
		t.Errorf("configBoolOnOff(nil) = %v, want known false", got)
	}
	if got := configBoolOnOff(ctx, m, "missing"); got.IsNull() || got.ValueBool() {
		t.Errorf("configBoolOnOff(missing) = %v, want known false", got)
	}
}
