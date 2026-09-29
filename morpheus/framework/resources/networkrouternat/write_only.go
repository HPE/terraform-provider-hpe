// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkrouternat

import (
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/wodigest"
)

// Change detection for the write-only attributes action, firewall, service: see
// wodigest. A change marks write_only_revision unknown, which plans an in-place update in which
// Update re-sends the configured values. The deprecated *_version markers keep working for the
// whole 2.2 line.
var wo = wodigest.Tracker{ //nolint:gosec // attribute names, not credentials
	Resource: "network_router_nat",
	Attributes: map[string]string{
		"action":   "action_version",
		"firewall": "firewall_version",
		"service":  "service_version",
	},
}
