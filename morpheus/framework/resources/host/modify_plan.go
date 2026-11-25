// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package host

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/clearable"
)

var _ resource.ResourceWithModifyPlan = &Resource{}

// ModifyPlan handles write-only changes and plans, as clears, the clearable attributes removed from the
// configuration since they were last applied (see clearableAttributes): a
// computed one such as ssh_username would otherwise be left to the appliance.
// Nothing to do on a create (null state) or a destroy (null plan).
func (r *Resource) ModifyPlan(
	ctx context.Context,
	req resource.ModifyPlanRequest,
	resp *resource.ModifyPlanResponse,
) {
	if req.Plan.Raw.IsNull() {
		return // destroy
	}

	// A changed write-only value plans an in-place update (or, for the
	// create-only variant passwords, a replacement); see write_only.go. The
	// appliance stamps last_updated on that update, and the framework has
	// not marked it unknown (nothing else in the plan differs from the
	// state), so it is marked here.
	wo.ModifyPlan(ctx, req, resp, path.Root("last_updated"))

	if req.State.Raw.IsNull() {
		return // create
	}

	_, diags := clearable.Plan(ctx, req, resp, clearableAttributes)
	resp.Diagnostics.Append(diags...)
}
