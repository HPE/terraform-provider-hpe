package modifiers

import (
	"context"

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

// RetainWhenStateSatisfiesRequest returns an Int64 plan modifier for an
// attribute that records a *request* which the platform may satisfy with a
// larger value.
//
// Instance volume size is the motivating case. Morpheus rounds a requested size
// up when the image needs more room than it allows — a request that clears the
// image's minimum disk but falls short of the image's own size is grown to fit.
// A request below the minimum is rejected outright rather than grown, so it is
// the rounding, not the minimum, that produces a value larger than was asked
// for.
//
// Where prior state already meets or exceeds the configured value, state is
// retained rather than planning a reduction the platform would refuse. This
// matters most after import: import reads the provisioned size, so without
// this a freshly imported instance whose disk had been grown would plan a
// shrink on every run.
//
// A configured value greater than state is a genuine request to grow and is
// left alone, so the resize path picks it up. A configured value below state
// is therefore a no-op while state exceeds it; Morpheus generally cannot
// shrink a disk in place, and a server-grown volume cannot be told apart from
// a deliberate reduction — both present as config < state, and Terraform
// supplies prior state and current config but never prior config.
//
// The attribute must be Computed as well as Optional. Terraform requires the
// planned value of a non-computed attribute to equal its configured value, so
// this modifier cannot legally take effect otherwise.
func RetainWhenStateSatisfiesRequest() planmodifier.Int64 {
	return retainWhenStateSatisfiesRequestModifier{}
}

type retainWhenStateSatisfiesRequestModifier struct{}

func (m retainWhenStateSatisfiesRequestModifier) Description(
	_ context.Context,
) string {
	return "Retain the prior value when it already satisfies the configured request"
}

func (m retainWhenStateSatisfiesRequestModifier) MarkdownDescription(
	ctx context.Context,
) string {
	return m.Description(ctx)
}

func (m retainWhenStateSatisfiesRequestModifier) PlanModifyInt64(
	_ context.Context,
	req planmodifier.Int64Request,
	resp *planmodifier.Int64Response,
) {
	// Create (no prior state) or destroy (no plan): nothing to retain.
	if req.State.Raw.IsNull() || req.Plan.Raw.IsNull() {
		return
	}

	if req.StateValue.IsNull() || req.StateValue.IsUnknown() {
		return
	}

	// Configuration is silent, so the value is whatever the platform chose.
	// Keep it rather than letting it go unknown on an unrelated change.
	if req.ConfigValue.IsNull() {
		resp.PlanValue = req.StateValue

		return
	}

	if req.ConfigValue.IsUnknown() {
		return
	}

	// State already satisfies the request; keep it. A larger configured value
	// falls through untouched and is treated as a request to grow.
	if req.ConfigValue.ValueInt64() <= req.StateValue.ValueInt64() {
		resp.PlanValue = req.StateValue
	}
}
