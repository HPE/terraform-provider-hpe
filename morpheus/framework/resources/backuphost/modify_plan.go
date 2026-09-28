// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package backuphost

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/clearable"
)

var _ resource.ResourceWithModifyPlan = &backupHostResource{}

// ModifyPlan plans a clear for each optional attribute removed from the
// configuration since it was last applied; see clearableAttributes.
func (r *backupHostResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// Nothing to clear on destroy or create.
	if req.Plan.Raw.IsNull() || req.State.Raw.IsNull() {
		return
	}

	_, diags := clearable.Plan(ctx, req, resp, clearableAttributes)
	resp.Diagnostics.Append(diags...)
}
