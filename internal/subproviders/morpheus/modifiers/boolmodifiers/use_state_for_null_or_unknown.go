// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package boolmodifiers

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

func UseStateForNullOrUnknown() planmodifier.Bool {
	return useStateForNullOrUnknownModifier{}
}

// useStateForNullOrUnknownModifier implements the plan modifier.
type useStateForNullOrUnknownModifier struct{}

// Description returns a human-readable description of the plan modifier.
func (m useStateForNullOrUnknownModifier) Description(_ context.Context) string {
	return "Once set, the value of this attribute in state will not change."
}

// MarkdownDescription returns a markdown description of the plan modifier.
func (m useStateForNullOrUnknownModifier) MarkdownDescription(_ context.Context) string {
	return "Once set, the value of this attribute in state will not change."
}

// PlanModifyBool implements the plan modification logic.
func (m useStateForNullOrUnknownModifier) PlanModifyBool(
	_ context.Context,
	req planmodifier.BoolRequest,
	resp *planmodifier.BoolResponse,
) {
	// Do nothing if there is a known planned value.
	if !req.PlanValue.IsUnknown() {
		return
	}

	// Do nothing if there is an unknown configuration value, otherwise interpolation gets messed up.
	if req.ConfigValue.IsUnknown() {
		return
	}

	// If we get to here, we have an unknown planned value OR a null planned value.
	resp.PlanValue = req.StateValue
}
