package modifiers

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
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

// Int64UseStateForUnknownUnless returns an Int64 plan modifier that copies the
// prior state value into the plan for an unknown (computed) value, like
// int64planmodifier.UseStateForUnknown, except when any of the given trigger
// attributes differ between state and plan. In that case the value is left
// unknown so a value the API recomputes from those attributes (for example a
// count derived from address ranges) is accepted without an "inconsistent
// result after apply" error.
func Int64UseStateForUnknownUnless(triggers ...path.Path) planmodifier.Int64 {
	return int64UseStateForUnknownUnlessModifier{triggers: triggers}
}

type int64UseStateForUnknownUnlessModifier struct {
	triggers []path.Path
}

func (m int64UseStateForUnknownUnlessModifier) Description(_ context.Context) string {
	return "Use prior state for unknown values unless a trigger attribute changes."
}

func (m int64UseStateForUnknownUnlessModifier) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

func (m int64UseStateForUnknownUnlessModifier) PlanModifyInt64(
	ctx context.Context,
	req planmodifier.Int64Request,
	resp *planmodifier.Int64Response,
) {
	// Create (no prior state) or destroy (no plan): nothing to carry forward.
	if req.State.Raw.IsNull() || req.Plan.Raw.IsNull() {
		return
	}

	// Only act when the value would otherwise be unknown (the framework marks
	// computed attributes unknown during update). If it is already known, leave
	// it so a no-op plan stays empty.
	if !req.PlanValue.IsUnknown() {
		return
	}

	// If any trigger attribute changes, leave the value unknown so the API can
	// recompute it.
	for _, p := range m.triggers {
		var stateVal, planVal attr.Value
		resp.Diagnostics.Append(req.State.GetAttribute(ctx, p, &stateVal)...)
		resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, p, &planVal)...)
		if resp.Diagnostics.HasError() {
			return
		}

		if !stateVal.Equal(planVal) {
			return
		}
	}

	// No trigger changed: carry the prior value forward.
	resp.PlanValue = req.StateValue
}

// NormalizeLineEndingsModifier normalizes CRLF ("\r\n") and lone CR ("\r") line
// endings to LF ("\n") in the planned string value. Some Morpheus APIs store
// text (for example script content) with LF line endings, so a configuration
// authored with CRLF line endings (common on Windows) would otherwise produce
// an "inconsistent result after apply" error because the planned value ("\r\n")
// differs from the value the API returns ("\n").
type NormalizeLineEndingsModifier struct{}

func (m NormalizeLineEndingsModifier) Description(_ context.Context) string {
	return "Normalizes CRLF and CR line endings to LF to match the API."
}

func (m NormalizeLineEndingsModifier) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

func (m NormalizeLineEndingsModifier) PlanModifyString(
	_ context.Context,
	req planmodifier.StringRequest,
	resp *planmodifier.StringResponse,
) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	original := req.ConfigValue.ValueString()
	normalized := strings.ReplaceAll(original, "\r\n", "\n")
	normalized = strings.ReplaceAll(normalized, "\r", "\n")

	if normalized != original {
		resp.PlanValue = types.StringValue(normalized)
	}
}
