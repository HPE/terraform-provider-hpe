// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkpool

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/clearable"
)

var _ resource.ResourceWithModifyPlan = &networkPoolResource{}

// ModifyPlan plans a clear for each optional attribute removed from the
// configuration since it was last applied; see clearableAttributes.
func (r *networkPoolResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// Nothing to clear on destroy or create.
	if req.Plan.Raw.IsNull() || req.State.Raw.IsNull() {
		return
	}

	_, diags := clearable.Plan(ctx, req, resp, clearableAttributes)
	resp.Diagnostics.Append(diags...)

	resp.Diagnostics.Append(carryRangeIDs(ctx, req, resp)...)
}

// carryRangeIDs fills in each planned ip_address_ranges element's computed id
// from the prior state — by id when the plan already has one, else by exact
// addresses — so an unchanged range does not plan as "(known after apply)"
// and an edited range is visibly the same range. A genuinely new range keeps
// an unknown id. Also marks ip_address_ranges unknown when the deprecated
// ip_ranges is what the configuration sets, so its computed value is refreshed
// after the apply.
func carryRangeIDs(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) diag.Diagnostics {
	var diags diag.Diagnostics

	var plan, state, config NetworkPoolModel
	diags.Append(req.Plan.Get(ctx, &plan)...)
	diags.Append(req.State.Get(ctx, &state)...)
	diags.Append(req.Config.Get(ctx, &config)...)
	if diags.HasError() {
		return diags
	}

	if config.IpAddressRanges.IsNull() {
		// Deprecated single-range attribute in use: ip_address_ranges is
		// purely computed here; it changes whenever the single range does.
		planned, _ := singleRangeFromObject(plan.IpRanges)
		prior, _ := singleRangeFromObject(state.IpRanges)
		if planned != prior {
			diags.Append(resp.Plan.SetAttribute(ctx, path.Root("ip_address_ranges"),
				types.ListUnknown(IpAddressRangesValue{}.Type(ctx)))...)
		}

		return diags
	}

	if plan.IpAddressRanges.IsUnknown() {
		return diags
	}

	planned, d := rangesFromList(ctx, plan.IpAddressRanges)
	diags.Append(d...)
	prior, d := rangesFromList(ctx, state.IpAddressRanges)
	diags.Append(d...)
	if diags.HasError() {
		return diags
	}

	carried := withIDsFromPrior(planned, prior)
	list, d := rangesToList(ctx, carried)
	diags.Append(d...)
	if diags.HasError() {
		return diags
	}
	// Elements with no id stay unknown rather than null: Morpheus assigns one.
	elems := list.Elements()
	for i, c := range carried {
		if c.id == 0 {
			v := elems[i].(IpAddressRangesValue)
			v.Id = types.Int64Unknown()
			elems[i] = v
		}
	}
	list, d = types.ListValue(IpAddressRangesValue{}.Type(ctx), elems)
	diags.Append(d...)
	diags.Append(resp.Plan.SetAttribute(ctx, path.Root("ip_address_ranges"), list)...)

	return diags
}
