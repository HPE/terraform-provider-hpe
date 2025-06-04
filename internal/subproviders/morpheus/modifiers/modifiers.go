package modifiers

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// RequireOnCreateModifier can be used for corner cases where
// an attribute is optional for import, but required for create
type RequireOnCreateModifier struct{}

func (m RequireOnCreateModifier) Description(_ context.Context) string {
	return "Requires the attribute to be set during resource creation."
}

func (m RequireOnCreateModifier) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

func (m RequireOnCreateModifier) PlanModifyString(
	ctx context.Context,
	req planmodifier.StringRequest,
	resp *planmodifier.StringResponse,
) {
	var id types.Int64

	diags := req.State.GetAttribute(ctx, path.Root("id"), &id)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)

		return
	}

	resourceExists := !id.IsNull() && !id.IsUnknown()

	if !resourceExists && req.ConfigValue.IsNull() {
		name := req.Path.String()
		msg := "attribute '" + name + "' not set " +
			"(this attribute is optional for some operations, eg import, " +
			"but needed during create)"
		resp.Diagnostics.AddError(
			"missing attribute",
			msg,
		)
	}
}
