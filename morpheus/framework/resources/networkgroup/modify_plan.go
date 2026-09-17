// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkgroup

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/tenancy"
)

var _ resource.ResourceWithModifyPlan = &networkGroupResource{}

// ModifyPlan rejects, at plan time, visibility = "public" set by a non-master
// caller. Morpheus silently coerces visibility to "private" for sub-tenant
// callers (MORPH-16419), which would otherwise leave state disagreeing with
// config. This covers both create and the private->public update path.
//
// Nothing to validate on destroy (Plan.Raw is null). The guard only calls
// whoami when visibility is a known "public" value, and fails open on an
// inability to determine the caller's tenancy.
func (r *networkGroupResource) ModifyPlan(
	ctx context.Context,
	req resource.ModifyPlanRequest,
	resp *resource.ModifyPlanResponse,
) {
	if req.Plan.Raw.IsNull() {
		return
	}

	tenancy.GuardVisibilityPublic(ctx, r, req.Plan, &resp.Diagnostics)
}
