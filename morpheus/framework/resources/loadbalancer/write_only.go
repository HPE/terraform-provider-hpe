// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package loadbalancer

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/wodigest"
)

// Change detection for the write-only attributes group_id, network_server_id: see
// wodigest. Morpheus takes these only at creation, so a change plans a replacement -- the only way
// the new value can take effect -- rather than an update that would silently do nothing. The
// deprecated *_version markers keep working for the
// whole 2.2 line.
var wo = wodigest.Tracker{ //nolint:gosec // attribute names, not credentials
	Resource: "load_balancer",
	Attributes: map[string]string{
		"group_id":          "group_id_version",
		"network_server_id": "network_server_id_version",
	},
	Immutable: []string{"group_id", "network_server_id"},
}

// ModifyPlan plans a replacement when a write-only value changed (or a
// deprecated marker did).
func (r *Resource) ModifyPlan(
	ctx context.Context,
	req resource.ModifyPlanRequest,
	resp *resource.ModifyPlanResponse,
) {
	wo.ModifyPlan(ctx, req, resp)
}
