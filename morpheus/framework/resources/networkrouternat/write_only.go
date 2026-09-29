// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkrouternat

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"

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

// ModifyPlan plans an in-place update when a write-only value changed (or a
// deprecated marker did).
func (r *Resource) ModifyPlan(
	ctx context.Context,
	req resource.ModifyPlanRequest,
	resp *resource.ModifyPlanResponse,
) {
	wo.ModifyPlan(ctx, req, resp)
}
