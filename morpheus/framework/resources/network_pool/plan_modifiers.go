package network_pool

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// unknownOnIPRangesChange returns an Int64 plan modifier for the Computed count
// attributes (ip_count, free_count). Those counts are recomputed by the API
// whenever the pool's ip_ranges change, so pinning them to the prior state
// value (as UseStateForUnknown does) produces an "inconsistent result after
// apply" error on update. This modifier keeps the prior value when ip_ranges is
// unchanged (avoiding a spurious diff) and marks the value unknown when
// ip_ranges changes (so the API-computed value is accepted).
func unknownOnIPRangesChange() planmodifier.Int64 {
	return ipRangesCountModifier{}
}

type ipRangesCountModifier struct{}

func (m ipRangesCountModifier) Description(_ context.Context) string {
	return "Recomputed when ip_ranges changes; otherwise retains the prior value."
}

func (m ipRangesCountModifier) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

func (m ipRangesCountModifier) PlanModifyInt64(
	ctx context.Context,
	req planmodifier.Int64Request,
	resp *planmodifier.Int64Response,
) {
	// Create: no prior state, leave the planned value (known after apply).
	if req.State.Raw.IsNull() {
		return
	}
	// Destroy: no plan, nothing to do.
	if req.Plan.Raw.IsNull() {
		return
	}

	var stateRanges, planRanges types.Object
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("ip_ranges"), &stateRanges)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("ip_ranges"), &planRanges)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if stateRanges.Equal(planRanges) {
		// ip_ranges unchanged: keep the known prior value to avoid a spurious diff.
		resp.PlanValue = req.StateValue

		return
	}

	// ip_ranges changed: the API will recompute the count.
	resp.PlanValue = types.Int64Unknown()
}
