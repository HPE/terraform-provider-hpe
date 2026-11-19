// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkrouter

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var _ resource.ResourceWithModifyPlan = &Resource{}

// ModifyPlan plans a replacement for a rename of an HPE Morpheus SDN router.
// The SDN router types (config_juniper_sdn / config_juniper_bgp) reject an
// update that changes the name ("Field name cannot be changed", verified live
// against plugin 1.1.0), while the other router types rename in place, so the
// rule is applied per type rather than through a schema-wide RequiresReplace.
func (r *Resource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// A changed write-only md5_key_wo plans an in-place update (write_only_revision).
	wo.ModifyPlan(ctx, req, resp)
	if resp.Diagnostics.HasError() {
		return
	}

	if req.Plan.Raw.IsNull() || req.State.Raw.IsNull() {
		return
	}

	var plan, state NetworkRouterModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	sdn := !plan.ConfigJuniperSdn.IsNull() || !plan.ConfigJuniperBgp.IsNull()
	if sdn && !plan.Name.IsUnknown() && !plan.Name.Equal(state.Name) {
		resp.RequiresReplace = append(resp.RequiresReplace, path.Root("name"))
	}
}
