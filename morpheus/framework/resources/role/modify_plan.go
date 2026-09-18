// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package role

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/tenancy"
)

var _ resource.ResourceWithModifyPlan = &Resource{}

// ModifyPlan rejects, at plan time, multitenant or multitenant_locked set true
// by a non-master caller. Both are master-tenant-only concepts: Morpheus
// silently ignores them for sub-tenant users, so the value the caller planned
// (true) is never stored and the read-back is false, yielding a confusing
// "Provider produced inconsistent result after apply" (was cty.True, but now
// cty.False). Failing fast avoids that (MORPH-17516).
//
// Both attributes are Optional+Computed with a booldefault.StaticBool(false),
// so at plan time they are always known and non-null; the bool guard fires only
// on a planned true, leaving the default (false) case untouched.
//
// Nothing to validate on destroy (Plan.Raw is null). Fail-open on an inability
// to determine the caller's tenancy.
func (r *Resource) ModifyPlan(
	ctx context.Context,
	req resource.ModifyPlanRequest,
	resp *resource.ModifyPlanResponse,
) {
	if req.Plan.Raw.IsNull() {
		return
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

	tenancy.GuardMasterOnlyBoolAttribute(
		ctx, client, req.Plan, path.Root("multitenant"), &resp.Diagnostics)
	tenancy.GuardMasterOnlyBoolAttribute(
		ctx, client, req.Plan, path.Root("multitenant_locked"), &resp.Diagnostics)
}
