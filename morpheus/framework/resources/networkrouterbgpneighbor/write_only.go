// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkrouterbgpneighbor

import (
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/wodigest"
)

// Change detection for the write-only attribute password_wo: see
// wodigest. A change marks write_only_revision unknown, which plans an in-place
// update in which Update re-sends the configured value. The deprecated
// *_version marker keeps working for the whole 2.2 line.
var wo = wodigest.Tracker{
	Resource: "network_router_bgp_neighbor",
	Attributes: map[string]string{
		"password_wo": "password_wo_version",
	},
}
