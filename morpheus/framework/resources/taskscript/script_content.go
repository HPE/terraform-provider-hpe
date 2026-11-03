// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package taskscript

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ScriptContentPlanModifier keeps the prior state value of script_content when
// it is only whitespace-different from the planned value. The SDKv2 resource
// compared script_content with strings.TrimSpace (a DiffSuppressFunc and a
// StateFunc), so a heredoc with a trailing newline did not produce a perpetual
// diff. This plan modifier reproduces that semantic equality: when
// TrimSpace(plan) == TrimSpace(state) the planned value is replaced with the
// state value, so Terraform plans no change.
type ScriptContentPlanModifier struct{}

// ScriptContentSemanticEquality returns the plan modifier for script_content.
func ScriptContentSemanticEquality() planmodifier.String {
	return ScriptContentPlanModifier{}
}

func (m ScriptContentPlanModifier) Description(_ context.Context) string {
	return "Suppresses a plan diff on script_content when the planned and prior " +
		"values differ only in leading or trailing whitespace."
}

func (m ScriptContentPlanModifier) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

func (m ScriptContentPlanModifier) PlanModifyString(
	_ context.Context,
	req planmodifier.StringRequest,
	resp *planmodifier.StringResponse,
) {
	// Only relevant on an update: both state and plan must be known values.
	if req.StateValue.IsNull() || req.StateValue.IsUnknown() {
		return
	}

	if req.PlanValue.IsNull() || req.PlanValue.IsUnknown() {
		return
	}

	if strings.TrimSpace(req.PlanValue.ValueString()) ==
		strings.TrimSpace(req.StateValue.ValueString()) {
		resp.PlanValue = types.StringValue(req.StateValue.ValueString())
	}
}
