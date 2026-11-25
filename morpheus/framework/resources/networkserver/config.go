// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkserver

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// parseJuniperSdnConfig builds the config_juniper_sdn block from the config
// map the API returns (on read and import). Morpheus stores the option type
// values as it received them, so the booleans may come back as real bools or
// as the strings "on"/"off"/"true"/"false"; parse tolerantly. When a boolean
// key is absent the value is false, matching the schema defaults.
//
// controllerAddress is the type's required option and the plugin keeps it on
// the server for its whole life, so it is always returned; were it ever
// missing, the prior value is kept rather than nulled, since the attribute
// requires replacement and a null read would plan a destroy of the network
// server — and with it the SDN onboarding — over a transient omission.
func parseJuniperSdnConfig(ctx context.Context, m map[string]interface{}, prior types.String) (ConfigJuniperSdnValue, error) {
	controllerAddress := configString(m, "controllerAddress")
	if controllerAddress.IsNull() && !prior.IsNull() && !prior.IsUnknown() {
		controllerAddress = prior
	}

	cfg, d := NewConfigJuniperSdnValue(
		ConfigJuniperSdnValue{}.AttributeTypes(ctx),
		map[string]attr.Value{
			"controller_address": controllerAddress,
			"show_introspect":    types.BoolValue(configBool(m, "showIntrospect")),
			"enable_sub_tenants": types.BoolValue(configBool(m, "enableSubTenants")),
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

// configBool coerces the truthy forms Morpheus may store an option-type
// checkbox as: a JSON bool, or the strings "on"/"true" (true) and
// "off"/"false"/"" (false). A missing key is false.
func configBool(m map[string]interface{}, key string) bool {
	switch v := m[key].(type) {
	case bool:
		return v
	case string:
		switch v {
		case "on", "true", "1":
			return true
		default:
			return false
		}
	default:
		return false
	}
}
