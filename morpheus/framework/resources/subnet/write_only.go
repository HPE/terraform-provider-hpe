// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package subnet

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/wodigest"
)

// Change detection for the write-only attribute config: see
// wodigest. Morpheus takes the configuration map only at creation, so a change plans a replacement
// -- the only way the new value can take effect -- rather than an update that would silently do
// nothing. The deprecated *_version marker keeps working for the
// whole 2.2 line.
var wo = wodigest.Tracker{ //nolint:gosec // attribute names, not credentials
	Resource: "subnet",
	Attributes: map[string]string{
		"config": "config_version",
	},
	Immutable: []string{"config"},
}

// ModifyPlan plans a replacement when a write-only value changed (or a
// deprecated marker did).
func (r *subnetResource) ModifyPlan(
	ctx context.Context,
	req resource.ModifyPlanRequest,
	resp *resource.ModifyPlanResponse,
) {
	wo.ModifyPlan(ctx, req, resp)
}
