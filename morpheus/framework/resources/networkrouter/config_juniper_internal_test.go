// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkrouter

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// The SDN router blocks are sent under the plugin's option type keys — the
// BGP ones dotted — on create (typed) and on update (the free-form map), and
// read back from the same keys on import.
func TestUnitJuniperRouterConfigRoundTrip(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	const uid1, uid2 = "0be6b79b-757c-4d92-affc-8084e328cb4c", "1c2d3e4f-5678-4a9b-bcde-0123456789ab"

	sdn, d := NewConfigJuniperSdnValue(ConfigJuniperSdnValue{}.AttributeTypes(ctx), map[string]attr.Value{
		"resource_pool_id": types.Int64Value(153047),
		"virtual_network_uids": types.SetValueMust(types.StringType,
			[]attr.Value{types.StringValue(uid1), types.StringValue(uid2)}),
		"floating_ip_pool_ids": types.SetNull(types.Int64Type),
	})
	if d.HasError() {
		t.Fatal(d)
	}

	cfg, d := juniperSdnConfig(ctx, sdn)
	if d.HasError() {
		t.Fatal(d)
	}

	b, _ := json.Marshal(cfg)
	want := `{"juniperResourcePool":153047,"juniperVirtualNetworks":` +
		`["` + uid1 + `","` + uid2 + `"]}`
	if string(b) != want {
		t.Errorf("create body = %s, want %s", b, want)
	}

	m, _ := juniperSdnUpdateConfig(ctx, sdn)
	if m["juniperResourcePool"] != int64(153047) || len(m["juniperVirtualNetworks"].([]string)) != 2 {
		t.Errorf("update map = %v", m)
	}

	// On update an unset list is sent as an explicit empty list: an omitted
	// key leaves the appliance's value unchanged, so a removed set would never
	// clear (and it reads back as null, matching the unset attribute).
	if pools, ok := m["juniperFloatingIpPools"].([]int64); !ok || len(pools) != 0 {
		t.Errorf("an unset list must be sent as an empty list on update, got %v", m["juniperFloatingIpPools"])
	}

	back, d := juniperSdnConfigFromMap(ctx, map[string]any{
		"juniperResourcePool":    float64(153047),
		"juniperVirtualNetworks": []any{uid1, uid2},
	})
	if d.HasError() || !back.ResourcePoolId.Equal(types.Int64Value(153047)) ||
		len(back.VirtualNetworkUids.Elements()) != 2 || !back.FloatingIpPoolIds.IsNull() {
		t.Errorf("import: %v %v", back, d)
	}

	bgp, d := NewConfigJuniperBgpValue(ConfigJuniperBgpValue{}.AttributeTypes(ctx), map[string]attr.Value{
		"vendor":           types.StringValue("juniper"),
		"asn":              types.Int64Value(65000),
		"peer_address_1":   types.StringValue("10.0.0.1"),
		"peer_address_2":   types.StringNull(),
		"peer_address_3":   types.StringNull(),
		"peer_address_4":   types.StringNull(),
		"address_families": types.SetValueMust(types.StringType, []attr.Value{types.StringValue("inet")}),
		"port":             types.Int64Value(179),
		"hold_time":        types.Int64Value(90),
		"admin_down":       types.StringValue("off"),
		"md5_key_wo":       types.StringNull(), // write-only: the plan never carries it
	})
	if d.HasError() {
		t.Fatal(d)
	}

	// The MD5 key comes from the configuration, not the plan block.
	cfg, d = juniperBgpConfig(ctx, bgp, types.StringValue("s3cret"))
	if d.HasError() {
		t.Fatal(d)
	}

	b, _ = json.Marshal(cfg)

	var got map[string]any
	_ = json.Unmarshal(b, &got)

	for key, want := range map[string]any{
		"juniper.sdn.vendor": "juniper", "juniper.sdn.asnnumber": float64(65000), "juniper.sdn.address.1": "10.0.0.1",
		"juniper.sdn.bgp.port": float64(179), "juniper.sdn.bgp.holdTime": float64(90), "juniper.sdn.bgp.adminDown": "off",
		"juniper.sdn.bgp.md5Key": "s3cret",
	} {
		if got[key] != want {
			t.Errorf("create %s = %v, want %v", key, got[key], want)
		}
	}

	if _, ok := got["juniper.sdn.address.2"]; ok {
		t.Error("an unset peer address must be omitted")
	}

	// Update: the key is sent only when its write-only digest changed — the
	// configured value, or "" to clear a removed key — and omitted otherwise.
	m, _ = juniperBgpUpdateConfig(ctx, bgp, types.StringValue("s3cret"), true)
	if m["juniper.sdn.vendor"] != "juniper" || m["juniper.sdn.bgp.md5Key"] != "s3cret" || m["juniper.sdn.bgp.port"] != int64(179) {
		t.Errorf("update map = %v", m)
	}

	m, _ = juniperBgpUpdateConfig(ctx, bgp, types.StringValue(""), true)
	if v, ok := m["juniper.sdn.bgp.md5Key"]; !ok || v != "" {
		t.Errorf("a configured \"\" must be sent to clear the md5 key, got %v (present %v)", v, ok)
	}

	m, _ = juniperBgpUpdateConfig(ctx, bgp, types.StringValue("s3cret"), false)
	if _, ok := m["juniper.sdn.bgp.md5Key"]; ok {
		t.Errorf("an unchanged md5 key must not be re-sent, got %v", m["juniper.sdn.bgp.md5Key"])
	}

	back2, d := juniperBgpConfigFromMap(ctx, map[string]any{
		"juniper.sdn.vendor": "juniper", "juniper.sdn.asnnumber": float64(65000), "juniper.sdn.address.1": "10.0.0.1",
		"bgpAddressFamilies": []any{"inet"}, "juniper.sdn.bgp.port": float64(179),
	})
	if d.HasError() || back2.Vendor.ValueString() != "juniper" || !back2.Md5KeyWo.IsNull() || !back2.PeerAddress2.IsNull() {
		t.Errorf("import: %v %v", back2, d)
	}
}

// TestUnitJuniperBgpConfigReadsTheNestedApiShape proves the BGP block is
// hydrated from the config as the API actually returns it (nested objects),
// which is what import and refresh see.
func TestUnitJuniperBgpConfigReadsTheNestedApiShape(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	api := map[string]any{
		"juniper": map[string]any{"sdn": map[string]any{
			"vendor":    "juniper",
			"asnnumber": float64(65000),
			"address":   map[string]any{"1": "10.0.0.1", "2": nil, "3": nil, "4": nil},
			"bgp":       map[string]any{"holdTime": float64(120), "port": float64(179), "adminDown": "off", "md5Key": nil},
		}},
		"bgpAddressFamilies": []any{"inet"},
	}

	cfg, d := juniperBgpConfigFromMap(ctx, api)
	if d.HasError() {
		t.Fatal(d)
	}
	if cfg.Vendor.ValueString() != "juniper" || cfg.Asn.ValueInt64() != 65000 ||
		cfg.PeerAddress1.ValueString() != "10.0.0.1" || cfg.HoldTime.ValueInt64() != 120 || cfg.Port.ValueInt64() != 179 {
		t.Fatalf("nested shape not read: %+v", cfg)
	}
	if !cfg.PeerAddress2.IsNull() {
		t.Fatalf("null peer address must stay null, got %v", cfg.PeerAddress2)
	}
}

// TestUnitConfigInt64ReadsNumbersAndNumericStrings pins that a numeric
// option-type value is read whether the appliance returns it as a JSON
// number or, as values stored from the UI often are, as a string — and that
// resource_pool_id, which now replaces the router on change, is never read as
// null for a string value.
func TestUnitConfigInt64ReadsNumbersAndNumericStrings(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		in   map[string]interface{}
		want types.Int64
	}{
		"number":         {map[string]interface{}{"juniperResourcePool": float64(9)}, types.Int64Value(9)},
		"numeric string": {map[string]interface{}{"juniperResourcePool": "9"}, types.Int64Value(9)},
		"padded string":  {map[string]interface{}{"juniperResourcePool": " 9 "}, types.Int64Value(9)},
		"other string":   {map[string]interface{}{"juniperResourcePool": "nine"}, types.Int64Null()},
		"absent":         {map[string]interface{}{}, types.Int64Null()},
		"null":           {map[string]interface{}{"juniperResourcePool": nil}, types.Int64Null()},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			if got := configInt64(c.in, "juniperResourcePool"); !got.Equal(c.want) {
				t.Errorf("got %v, want %v", got, c.want)
			}
		})
	}
}
