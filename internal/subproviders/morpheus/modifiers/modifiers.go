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

// NullableStringUpdateModifier can be used when the desired state of
// a string is null and the current state is non-null. Usually
// terraform plan will not treat this as something that should
// trigger an update. But using this modifier will cause plan
// to trigger an update, eg "foo" -> null
type NullableStringUpdateModifier struct{}

func (m NullableStringUpdateModifier) Description(_ context.Context) string {
	return "Force diff when config changes from non-null to null" // nolint: goconst
}

func (m NullableStringUpdateModifier) MarkdownDescription(_ context.Context) string {
	return "Force diff when config changes from non-null to null"
}

func (m NullableStringUpdateModifier) PlanModifyString(
	_ context.Context,
	req planmodifier.StringRequest,
	resp *planmodifier.StringResponse,
) {
	if req.ConfigValue.IsNull() && !req.StateValue.IsNull() {
		resp.PlanValue = types.StringNull()
	}
}

// NullableInt64UpdateModifier can be used when the desired state of
// an int64 is null and the current state is non-null. Usually
// terraform plan will not treat this as something that should
// trigger an update. But using this modifier will cause plan
// to trigger an update, eg 100 -> null
type NullableInt64UpdateModifier struct{}

func (m NullableInt64UpdateModifier) Description(_ context.Context) string {
	return "Force diff when config changes from non-null to null"
}

func (m NullableInt64UpdateModifier) MarkdownDescription(_ context.Context) string {
	return "Force diff when config changes from non-null to null"
}

func (m NullableInt64UpdateModifier) PlanModifyInt64(
	_ context.Context,
	req planmodifier.Int64Request,
	resp *planmodifier.Int64Response,
) {
	if req.ConfigValue.IsNull() && !req.StateValue.IsNull() {
		resp.PlanValue = types.Int64Null()
	}
}
