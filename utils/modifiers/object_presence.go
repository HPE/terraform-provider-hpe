// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package modifiers

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

// ObjectRequiresReplaceOnPresenceChange requires the resource to be replaced
// when a nested object attribute is added to or removed from the
// configuration — the block goes from null to set, or from set to null — while
// leaving changes inside a block that stays present to be applied as an update.
//
// It suits a resource whose type is chosen by which of several typed blocks
// is set: swapping one block for another changes the type, which the API
// cannot do in place, while editing the block's fields is an ordinary update.
// Nothing is required on create or destroy, where there is no prior or no
// planned resource to compare with.
func ObjectRequiresReplaceOnPresenceChange() planmodifier.Object {
	return objectPresenceModifier{}
}

type objectPresenceModifier struct{}

func (objectPresenceModifier) Description(_ context.Context) string {
	return "Adding or removing this block replaces the resource; changes within it are applied in place."
}

func (m objectPresenceModifier) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

func (objectPresenceModifier) PlanModifyObject(
	_ context.Context, req planmodifier.ObjectRequest, resp *planmodifier.ObjectResponse,
) {
	// Create or destroy: nothing to replace.
	if req.State.Raw.IsNull() || req.Plan.Raw.IsNull() {
		return
	}

	if req.PlanValue.IsUnknown() {
		return
	}

	if req.StateValue.IsNull() != req.PlanValue.IsNull() {
		resp.RequiresReplace = true
	}
}
