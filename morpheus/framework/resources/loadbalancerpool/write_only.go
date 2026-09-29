// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package loadbalancerpool

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/wodigest"
)

// Change detection for the write-only partition: see wodigest. A change marks
// write_only_revision unknown, which plans an in-place update in which Update
// sends the configured value. The resource has no marker attribute.
var wo = wodigest.Tracker{
	Resource:   "load_balancer_pool",
	Attributes: map[string]string{"partition": ""},
}

// ModifyPlan plans an in-place update when partition changed.
func (r *Resource) ModifyPlan(
	ctx context.Context,
	req resource.ModifyPlanRequest,
	resp *resource.ModifyPlanResponse,
) {
	wo.ModifyPlan(ctx, req, resp)
}
