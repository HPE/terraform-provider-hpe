// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkrouter

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
)

// The HPE Morpheus SDN network plugin for HPE VM Essentials registers two
// router types. Their configuration fields are option types with
// fieldContext "config", so the API carries them in the router's config map:
//
//   - juniperSdnRouter: a logical router joining SDN virtual networks
//     (juniperResourcePool, juniperVirtualNetworks, juniperFloatingIpPools);
//   - juniperBGPSdnRouter: a BGP peering (juniper.sdn.vendor,
//     juniper.sdn.asnnumber, juniper.sdn.address.1..4, bgpAddressFamilies,
//     juniper.sdn.bgp.port, holdTime, adminDown, md5Key).
//
// Unlike the NSX-T gateways, the plugin marks these fields editable, so
// Update sends them too (as the free-form map the update request takes; the
// keys are the same). resource_pool_id — the cluster the router lives on —
// replaces the router when it changes; Update still carries the unchanged
// value, which is harmless.
const (
	codeJuniperSdnRouter    = "juniperSdnRouter"
	codeJuniperBgpSdnRouter = "juniperBGPSdnRouter"
)

// juniperSdnConfig builds the typed create configuration for a
// juniperSdnRouter from the config_juniper_sdn block.
func juniperSdnConfig(
	ctx context.Context, cfg ConfigJuniperSdnValue,
) (sdk.CreateNetworkRouterRequestNetworkRouterConfig, diag.Diagnostics) {
	var diags diag.Diagnostics

	c := &sdk.JuniperSDNRouterConfig{}

	if !cfg.ResourcePoolId.IsNull() && !cfg.ResourcePoolId.IsUnknown() {
		c.JuniperResourcePool.Set(cfg.ResourcePoolId.ValueInt64Pointer())
	}

	uids, d := stringSet(ctx, cfg.VirtualNetworkUids)
	diags.Append(d...)
	if uids != nil {
		c.JuniperVirtualNetworks = uids
	}

	ids, d := int64Set(ctx, cfg.FloatingIpPoolIds)
	diags.Append(d...)
	if ids != nil {
		c.JuniperFloatingIpPools = ids
	}

	return sdk.CreateNetworkRouterRequestNetworkRouterConfig{JuniperSDNRouterConfig: c}, diags
}

// juniperBgpConfig builds the typed create configuration for a
// juniperBGPSdnRouter from the config_juniper_bgp block. vendor,
// peer_address_1 and address_families are required by the schema; asn, port,
// hold_time and admin_down carry the plugin's defaults.
func juniperBgpConfig(
	ctx context.Context, cfg ConfigJuniperBgpValue, md5Key types.String,
) (sdk.CreateNetworkRouterRequestNetworkRouterConfig, diag.Diagnostics) {
	var diags diag.Diagnostics

	c := &sdk.JuniperBGPSDNRouterConfig{
		JuniperSdnVendor:    cfg.Vendor.ValueString(),
		JuniperSdnAsnnumber: cfg.Asn.ValueInt64(),
		JuniperSdnAddress1:  cfg.PeerAddress1.ValueString(),
	}

	setIfKnown(&c.JuniperSdnAddress2, cfg.PeerAddress2)
	setIfKnown(&c.JuniperSdnAddress3, cfg.PeerAddress3)
	setIfKnown(&c.JuniperSdnAddress4, cfg.PeerAddress4)

	families, d := stringSet(ctx, cfg.AddressFamilies)
	diags.Append(d...)
	c.BgpAddressFamilies = families

	if !cfg.Port.IsNull() && !cfg.Port.IsUnknown() {
		c.JuniperSdnBgpPort = cfg.Port.ValueInt64Pointer()
	}

	if !cfg.HoldTime.IsNull() && !cfg.HoldTime.IsUnknown() {
		c.JuniperSdnBgpHoldTime = cfg.HoldTime.ValueInt64Pointer()
	}

	if !cfg.AdminDown.IsNull() && !cfg.AdminDown.IsUnknown() {
		c.JuniperSdnBgpAdminDown = cfg.AdminDown.ValueStringPointer()
	}

	// The MD5 key is write-only, so it comes from the configuration, not the plan.
	setIfKnown(&c.JuniperSdnBgpMd5Key, md5Key)

	return sdk.CreateNetworkRouterRequestNetworkRouterConfig{JuniperBGPSDNRouterConfig: c}, diags
}

// juniperSdnUpdateConfig is the juniperSdnRouter configuration as the update
// request's free-form map, with the plugin's keys.
func juniperSdnUpdateConfig(ctx context.Context, cfg ConfigJuniperSdnValue) (map[string]any, diag.Diagnostics) {
	var diags diag.Diagnostics

	m := map[string]any{}

	if !cfg.ResourcePoolId.IsNull() && !cfg.ResourcePoolId.IsUnknown() {
		m["juniperResourcePool"] = cfg.ResourcePoolId.ValueInt64()
	}

	// On update an omitted set is sent as an explicit empty list: the
	// appliance leaves an omitted key unchanged, and Morpheus reads an empty
	// option-type list back as null ([null]), which is what the omitted
	// attribute is in state. A configured empty set is refused at plan time
	// (SizeAtLeast(1)), since it would read back as null too.
	uids, d := stringSet(ctx, cfg.VirtualNetworkUids)
	diags.Append(d...)
	if uids == nil {
		uids = []string{}
	}
	m["juniperVirtualNetworks"] = uids

	ids, d := int64Set(ctx, cfg.FloatingIpPoolIds)
	diags.Append(d...)
	if ids == nil {
		ids = []int64{}
	}
	m["juniperFloatingIpPools"] = ids

	return m, diags
}

// juniperBgpUpdateConfig is the juniperBGPSdnRouter configuration as the
// update request's free-form map, with the plugin's dotted keys.
func juniperBgpUpdateConfig(
	ctx context.Context, cfg ConfigJuniperBgpValue, md5Key types.String, sendMd5Key bool,
) (map[string]any, diag.Diagnostics) {
	var diags diag.Diagnostics

	m := map[string]any{
		"juniper.sdn.vendor":    cfg.Vendor.ValueString(),
		"juniper.sdn.asnnumber": cfg.Asn.ValueInt64(),
		"juniper.sdn.address.1": cfg.PeerAddress1.ValueString(),
	}

	// Optional fields are always sent on update: a removed peer address goes
	// as "" so the plugin clears it rather than keeping the old value (it
	// treats a missing key as "unchanged").
	for key, v := range map[string]types.String{
		"juniper.sdn.address.2":     cfg.PeerAddress2,
		"juniper.sdn.address.3":     cfg.PeerAddress3,
		"juniper.sdn.address.4":     cfg.PeerAddress4,
		"juniper.sdn.bgp.adminDown": cfg.AdminDown,
	} {
		if v.IsUnknown() {
			continue
		}
		m[key] = v.ValueString() // "" when null
	}

	// The MD5 key is write-only: it is sent only when its digest says it
	// changed (a configured "" clears it on the plugin). Left out otherwise,
	// which the plugin treats as unchanged — including when the attribute was
	// removed from the configuration, which wodigest treats as unmanaged.
	if sendMd5Key && !md5Key.IsUnknown() {
		m["juniper.sdn.bgp.md5Key"] = md5Key.ValueString()
	}

	families, d := stringSet(ctx, cfg.AddressFamilies)
	diags.Append(d...)
	if families != nil {
		m["bgpAddressFamilies"] = families
	}

	if !cfg.Port.IsNull() && !cfg.Port.IsUnknown() {
		m["juniper.sdn.bgp.port"] = cfg.Port.ValueInt64()
	}

	if !cfg.HoldTime.IsNull() && !cfg.HoldTime.IsUnknown() {
		m["juniper.sdn.bgp.holdTime"] = cfg.HoldTime.ValueInt64()
	}

	return m, diags
}

// setIfKnown sets a nullable SDK string from a known planned value and leaves
// it unsent otherwise.
func setIfKnown(target *sdk.NullableString, v types.String) {
	if v.IsNull() || v.IsUnknown() {
		return
	}

	target.Set(v.ValueStringPointer())
}

func int64Set(ctx context.Context, v types.Set) ([]int64, diag.Diagnostics) {
	if v.IsNull() || v.IsUnknown() {
		return nil, nil
	}

	var out []int64

	diags := v.ElementsAs(ctx, &out, false)
	if out == nil {
		out = []int64{}
	}

	return out, diags
}

func stringSet(ctx context.Context, v types.Set) ([]string, diag.Diagnostics) {
	if v.IsNull() || v.IsUnknown() {
		return nil, nil
	}

	var out []string

	diags := v.ElementsAs(ctx, &out, false)
	if out == nil {
		out = []string{}
	}

	return out, diags
}

// juniperSdnConfigFromMap builds the config_juniper_sdn block from the config
// map the API returns for a juniperSdnRouter (on import and refresh).
func juniperSdnConfigFromMap(ctx context.Context, m map[string]any) (ConfigJuniperSdnValue, diag.Diagnostics) {
	attrs := map[string]attr.Value{
		"resource_pool_id":     configInt64(m, "juniperResourcePool"),
		"virtual_network_uids": configStringSet(m, "juniperVirtualNetworks"),
		"floating_ip_pool_ids": configInt64Set(m, "juniperFloatingIpPools"),
	}

	return NewConfigJuniperSdnValue(ConfigJuniperSdnValue{}.AttributeTypes(ctx), attrs)
}

// flattenDotted rewrites a nested config map into one keyed by dotted paths:
// Morpheus stores the BGP router's option values under dotted keys
// ("juniper.sdn.bgp.holdTime") but returns them as nested objects
// ({"juniper":{"sdn":{"bgp":{"holdTime":90}}}}), so the dotted readers below
// would otherwise see nothing on import or refresh. Keys that are already flat
// pass through, so a map in either shape reads the same.
func flattenDotted(m map[string]any) map[string]any {
	out := make(map[string]any, len(m))

	var walk func(prefix string, v any)
	walk = func(prefix string, v any) {
		if nested, ok := v.(map[string]any); ok {
			for k, nv := range nested {
				walk(prefix+"."+k, nv)
			}

			return
		}
		out[prefix] = v
	}

	for k, v := range m {
		if nested, ok := v.(map[string]any); ok {
			for nk, nv := range nested {
				walk(k+"."+nk, nv)
			}

			continue
		}
		out[k] = v
	}

	return out
}

// juniperBgpConfigFromMap builds the config_juniper_bgp block from the config
// map the API returns for a juniperBGPSdnRouter (on import and refresh). The MD5 key is
// never returned and reads as null.
func juniperBgpConfigFromMap(ctx context.Context, m map[string]any) (ConfigJuniperBgpValue, diag.Diagnostics) {
	m = flattenDotted(m)
	attrs := map[string]attr.Value{
		"vendor":           configString(m, "juniper.sdn.vendor"),
		"asn":              configInt64(m, "juniper.sdn.asnnumber"),
		"peer_address_1":   configString(m, "juniper.sdn.address.1"),
		"peer_address_2":   configString(m, "juniper.sdn.address.2"),
		"peer_address_3":   configString(m, "juniper.sdn.address.3"),
		"peer_address_4":   configString(m, "juniper.sdn.address.4"),
		"address_families": configStringSet(m, "bgpAddressFamilies"),
		"port":             configInt64(m, "juniper.sdn.bgp.port"),
		"hold_time":        configInt64(m, "juniper.sdn.bgp.holdTime"),
		"admin_down":       configString(m, "juniper.sdn.bgp.adminDown"),
		"md5_key_wo":       types.StringNull(), // write-only: never in state
	}

	return NewConfigJuniperBgpValue(ConfigJuniperBgpValue{}.AttributeTypes(ctx), attrs)
}

// configInt64Set reads a list of integers (or numeric strings) from a config
// map as a set; absent → null.
func configInt64Set(m map[string]any, key string) types.Set {
	raw, ok := m[key]
	if !ok || raw == nil {
		return types.SetNull(types.Int64Type)
	}

	items, ok := raw.([]any)
	if !ok {
		return types.SetNull(types.Int64Type)
	}

	vals := make([]attr.Value, 0, len(items))

	for _, it := range items {
		switch v := it.(type) {
		case float64:
			vals = append(vals, types.Int64Value(int64(v)))
		case int64:
			vals = append(vals, types.Int64Value(v))
		case string:
			if n, err := strconv.ParseInt(v, 10, 64); err == nil {
				vals = append(vals, types.Int64Value(n))
			}
		}
	}

	// Morpheus stores an unset option-type list as [null] (one empty entry),
	// so a list with no usable id is read as null, not as an empty set -- an
	// optional set the configuration omits must not come back as [] (an
	// inconsistent result after apply).
	if len(vals) == 0 {
		return types.SetNull(types.Int64Type)
	}

	return types.SetValueMust(types.Int64Type, vals)
}

// configStringSet reads a list of strings from a config map as a set; absent
// → null.
func configStringSet(m map[string]any, key string) types.Set {
	raw, ok := m[key]
	if !ok || raw == nil {
		return types.SetNull(types.StringType)
	}

	items, ok := raw.([]any)
	if !ok {
		return types.SetNull(types.StringType)
	}

	vals := make([]attr.Value, 0, len(items))

	for _, it := range items {
		if s, ok := it.(string); ok && s != "" {
			vals = append(vals, types.StringValue(s))
		}
	}

	// As for configInt64Set: Morpheus stores an unset option-type list as
	// [null], which must read as null rather than as an empty set.
	if len(vals) == 0 {
		return types.SetNull(types.StringType)
	}

	return types.SetValueMust(types.StringType, vals)
}
