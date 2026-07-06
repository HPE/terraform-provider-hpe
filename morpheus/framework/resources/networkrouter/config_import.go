// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkrouter

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

// The Morpheus API persists a network router's config as an opaque JSON object
// and echoes it back on GET under the same keys/values that were sent. NSX-T
// checkbox fields round-trip as "on"/"off" strings (not booleans) under their
// uppercase keys (e.g. TIER1_CONNECTED); select/string fields round-trip under
// their camelCase keys (e.g. ipManagementType). Keys that were never set come
// back as null. These helpers reverse that map into the typed config objects so
// that importing a router hydrates config_nsxt_gateway_tier0 /
// config_nsxt_gateway_tier1 to match a freshly-applied resource.

// configString reads a string config value, returning a null string when the
// key is absent or not a string.
func configString(m map[string]interface{}, key string) types.String {
	if v, ok := m[key]; ok && v != nil {
		if s, ok := v.(string); ok {
			return types.StringValue(s)
		}
	}

	return types.StringNull()
}

// configInt64 reads a numeric config value (JSON numbers decode as float64),
// returning a null int64 when the key is absent or not numeric.
func configInt64(m map[string]interface{}, key string) types.Int64 {
	if v, ok := m[key]; ok && v != nil {
		if f, ok := v.(float64); ok {
			return types.Int64Value(int64(f))
		}
	}

	return types.Int64Null()
}

// configBoolOnOff reads an NSX-T checkbox config value using convert.StringToBool
// ("on"/"off" -> true/false). The typed schema fields are computed with a false
// default, so an absent/null/unrecognized value maps to false (not null) to keep
// import equivalent to a freshly-applied resource.
func configBoolOnOff(ctx context.Context, m map[string]interface{}, key string) types.Bool {
	if v, ok := m[key]; ok && v != nil {
		if s, ok := v.(string); ok {
			if b := convert.StringToBool(ctx, s); !b.IsNull() {
				return b
			}
		}
	}

	return types.BoolValue(false)
}

// tier1ConfigFromMap builds a typed NSX-T Tier-1 config object from the API
// config map returned by GET.
func tier1ConfigFromMap(
	ctx context.Context,
	m map[string]interface{},
) (ConfigNsxtGatewayTier1Value, diag.Diagnostics) {
	attrs := map[string]attr.Value{
		"edge_cluster":               configString(m, "edgeCluster"),
		"fail_over":                  configString(m, "failOver"),
		"ip_management_type":         configString(m, "ipManagementType"),
		"ip_server_id":               configString(m, "ipServerId"),
		"tier0_gateway":              configString(m, "tier0Gateway"),
		"tier1_connected":            configBoolOnOff(ctx, m, "TIER1_CONNECTED"),
		"tier1_dns_forwarder_ip":     configBoolOnOff(ctx, m, "TIER1_DNS_FORWARDER_IP"),
		"tier1_ipsec_local_endpoint": configBoolOnOff(ctx, m, "TIER1_IPSEC_LOCAL_ENDPOINT"),
		"tier1_lb_snat":              configBoolOnOff(ctx, m, "TIER1_LB_SNAT"),
		"tier1_lb_vip":               configBoolOnOff(ctx, m, "TIER1_LB_VIP"),
		"tier1_nat":                  configBoolOnOff(ctx, m, "TIER1_NAT"),
		"tier1_static_routes":        configBoolOnOff(ctx, m, "TIER1_STATIC_ROUTES"),
	}

	return NewConfigNsxtGatewayTier1Value(
		ConfigNsxtGatewayTier1Value{}.AttributeTypes(ctx), attrs,
	)
}

// tier0ConfigFromMap builds a typed NSX-T Tier-0 config object from the API
// config map returned by GET.
func tier0ConfigFromMap(
	ctx context.Context,
	m map[string]interface{},
) (ConfigNsxtGatewayTier0Value, diag.Diagnostics) {
	attrs := map[string]attr.Value{
		"ecmp":                       configBoolOnOff(ctx, m, "ECMP"),
		"edge_cluster":               configString(m, "edgeCluster"),
		"fail_over":                  configString(m, "failOver"),
		"ha_mode":                    configString(m, "haMode"),
		"inter_sr_ibgp":              configBoolOnOff(ctx, m, "INTER_SR_IBGP"),
		"ip_management_type":         configString(m, "ipManagementType"),
		"ip_server_id":               configString(m, "ipServerId"),
		"local_as_num":               configString(m, "LOCAL_AS_NUM"),
		"multipath_relax":            configBoolOnOff(ctx, m, "MULTIPATH_RELAX"),
		"restart_mode":               configString(m, "RESTART_MODE"),
		"restart_time":               configInt64(m, "RESTART_TIME"),
		"stale_route_time":           configInt64(m, "STALE_ROUTE_TIME"),
		"tier0_dns_forwarder_ip":     configBoolOnOff(ctx, m, "TIER0_DNS_FORWARDER_IP"),
		"tier0_external_interface":   configBoolOnOff(ctx, m, "TIER0_EXTERNAL_INTERFACE"),
		"tier0_ipsec_local_ip":       configBoolOnOff(ctx, m, "TIER0_IPSEC_LOCAL_IP"),
		"tier0_loopback_interface":   configBoolOnOff(ctx, m, "TIER0_LOOPBACK_INTERFACE"),
		"tier0_nat":                  configBoolOnOff(ctx, m, "TIER0_NAT"),
		"tier0_segment":              configBoolOnOff(ctx, m, "TIER0_SEGMENT"),
		"tier0_service_interface":    configBoolOnOff(ctx, m, "TIER0_SERVICE_INTERFACE"),
		"tier0_static":               configBoolOnOff(ctx, m, "TIER0_STATIC"),
		"tier1_dns_forwarder_ip":     configBoolOnOff(ctx, m, "TIER1_DNS_FORWARDER_IP"),
		"tier1_ipsec_local_endpoint": configBoolOnOff(ctx, m, "TIER1_IPSEC_LOCAL_ENDPOINT"),
		"tier1_lb_snat":              configBoolOnOff(ctx, m, "TIER1_LB_SNAT"),
		"tier1_lb_vip":               configBoolOnOff(ctx, m, "TIER1_LB_VIP"),
		"tier1_nat":                  configBoolOnOff(ctx, m, "TIER1_NAT"),
		"tier1_segment":              configBoolOnOff(ctx, m, "TIER1_SEGMENT"),
		"tier1_service_interface":    configBoolOnOff(ctx, m, "TIER1_SERVICE_INTERFACE"),
		"tier1_static":               configBoolOnOff(ctx, m, "TIER1_STATIC"),
	}

	return NewConfigNsxtGatewayTier0Value(
		ConfigNsxtGatewayTier0Value{}.AttributeTypes(ctx), attrs,
	)
}
