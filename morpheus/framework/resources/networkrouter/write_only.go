// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkrouter

import (
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/wodigest"
)

// Change detection for the write-only attribute config_juniper_bgp.md5_key_wo:
// see wodigest. A change marks write_only_revision unknown, which plans an
// in-place update in which Update re-sends the configured value ("" clears
// the key; an attribute removed from the configuration is unmanaged, as for
// every write-only value in the provider).
var wo = wodigest.Tracker{
	Resource: "network_router",
	Attributes: map[string]string{ //nolint:gosec // an attribute name, not a credential
		"config_juniper_bgp.md5_key_wo": "",
	},
}

// configuredMd5Key returns the MD5 key from the configuration model: a
// write-only value is available there only, never in the plan or the state.
// Null when the BGP block or the key is not configured.
func configuredMd5Key(config NetworkRouterModel) types.String {
	if config.ConfigJuniperBgp.IsNull() || config.ConfigJuniperBgp.IsUnknown() {
		return types.StringNull()
	}

	return config.ConfigJuniperBgp.Md5KeyWo
}
