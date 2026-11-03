// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package cloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/clearable"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/tenancy"
)

var _ resource.ResourceWithModifyPlan = &Resource{}

// ModifyPlan rejects, at plan time, visibility = "public" set by a non-master
// caller. Morpheus silently coerces visibility to "private" for sub-tenant
// callers (MORPH-16419), which would otherwise leave state disagreeing with
// config. This covers both create and the private->public update path.
//
// Nothing to validate on destroy (Plan.Raw is null). The guard only calls
// whoami when visibility is a known "public" value, and fails open on an
// inability to determine the caller's tenancy.
func (r *Resource) ModifyPlan(
	ctx context.Context,
	req resource.ModifyPlanRequest,
	resp *resource.ModifyPlanResponse,
) {
	if req.Plan.Raw.IsNull() {
		return
	}

	// A changed write-only value (or deprecated marker) plans an in-place
	// update; see write_only.go. It runs first, as on every resource that
	// has both.
	wo.ModifyPlan(ctx, req, resp)

	// Attributes removed from the configuration since they were last applied
	// are planned as clears; see clearableAttributes.
	if !req.State.Raw.IsNull() {
		_, cdiags := clearable.Plan(ctx, req, resp, clearableAttributes)
		resp.Diagnostics.Append(cdiags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	// Surface a client construction failure immediately as a diagnostic. The
	// guard's whoami call still fails open on a tenancy-determination error.
	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create Morpheus client",
			err.Error(),
		)

		return
	}

	tenancy.GuardVisibilityPublic(ctx, client, req.Plan, &resp.Diagnostics)
}
