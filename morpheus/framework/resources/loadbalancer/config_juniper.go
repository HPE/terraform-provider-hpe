// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package loadbalancer

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
)

// The HPE Morpheus SDN network plugin's load balancer type (juniper-sdn-lb)
// takes its configuration through option types with fieldContext "config":
// juniperResourcePool, virtualNetwork, virtualIp, sourcePort, destPort and
// protocol. None is editable after create, so the block replaces the load
// balancer when it changes and Update never sends it.

// juniperSdnConfig builds the typed create configuration from the
// config_juniper_sdn block. resource_pool_id, virtual_network_uid and
// virtual_ip are required by the schema.
func juniperSdnConfig(cfg ConfigJuniperSdnValue) *sdk.JuniperSDNLoadBalancerConfig {
	c := &sdk.JuniperSDNLoadBalancerConfig{
		JuniperResourcePool: cfg.ResourcePoolId.ValueInt64(),
		VirtualNetwork:      cfg.VirtualNetworkUid.ValueString(),
		VirtualIp:           cfg.VirtualIp.ValueString(),
	}

	if !cfg.SourcePort.IsNull() && !cfg.SourcePort.IsUnknown() {
		c.SourcePort.Set(cfg.SourcePort.ValueInt64Pointer())
	}

	if !cfg.DestinationPort.IsNull() && !cfg.DestinationPort.IsUnknown() {
		c.DestPort.Set(cfg.DestinationPort.ValueInt64Pointer())
	}

	if !cfg.Protocol.IsNull() && !cfg.Protocol.IsUnknown() {
		c.Protocol.Set(cfg.Protocol.ValueStringPointer())
	}

	return c
}

// parseJuniperSdnConfig builds the config_juniper_sdn block from the config
// map the API returns (on read and import). Morpheus stores the option type
// values as it received them, so ids may come back as numbers or as strings.
func parseJuniperSdnConfig(ctx context.Context, m map[string]interface{}) (ConfigJuniperSdnValue, error) {
	cfg, d := NewConfigJuniperSdnValue(
		ConfigJuniperSdnValue{}.AttributeTypes(ctx),
		map[string]attr.Value{
			"resource_pool_id":    configInt64(m, "juniperResourcePool"),
			"virtual_network_uid": configString(m, "virtualNetwork"),
			"virtual_ip":          configString(m, "virtualIp"),
			"source_port":         configInt64(m, "sourcePort"),
			"destination_port":    configInt64(m, "destPort"),
			"protocol":            configString(m, "protocol"),
		},
	)
	if d.HasError() {
		return ConfigJuniperSdnValue{}, fmt.Errorf("failed to build config_juniper_sdn value: %s", d.Errors())
	}

	return cfg, nil
}

func configString(m map[string]interface{}, key string) types.String {
	if v, ok := m[key].(string); ok && v != "" {
		return types.StringValue(v)
	}

	return types.StringNull()
}

func configInt64(m map[string]interface{}, key string) types.Int64 {
	switch v := m[key].(type) {
	case float64:
		return types.Int64Value(int64(v))
	case int64:
		return types.Int64Value(v)
	case string:
		// Whole numeric strings only, as the router's reader: Sscanf would
		// accept the leading digits of "9abc".
		if n, err := strconv.ParseInt(strings.TrimSpace(v), 10, 64); err == nil {
			return types.Int64Value(n)
		}
	}

	return types.Int64Null()
}
