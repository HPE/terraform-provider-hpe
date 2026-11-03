// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkrouterroute

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/clearable"
)

var _ resource.ResourceWithModifyPlan = &Resource{}

// ModifyPlan plans a replacement for each optional attribute removed from the
// configuration since it was last applied: the route has no update, so the
// removal, like a change, recreates it; see clearableAttributes.
func (r *Resource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// Nothing to clear on destroy or create.
	if req.Plan.Raw.IsNull() || req.State.Raw.IsNull() {
		return
	}

	pending, diags := clearable.Plan(ctx, req, resp, clearableAttributes)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	for _, name := range pending.Names() {
		resp.RequiresReplace = append(resp.RequiresReplace, path.Root(name))
	}
}
