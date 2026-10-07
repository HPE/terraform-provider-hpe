// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package plugin

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// jarSha256PlanModifier drives when jar_sha256 is treated as known vs unknown at
// plan time, so a changed jar surfaces as an update:
//
//   - configured by the user -> keep the configured value;
//   - no prior state (create) -> unknown;
//   - jar_path or jar_url changed vs state -> unknown (a re-upload is coming);
//   - otherwise -> the prior state value.
//
// Additionally, when jar_path is set and known and jar_sha256 is not configured,
// the local file is hashed during plan (when readable) and that digest becomes
// the planned value — so a rebuilt jar at the same path shows as an update even
// without filesha256(). If the file cannot be read at plan time the value is
// left per the rules above; Create/Update then error clearly on the missing file.
type jarSha256PlanModifier struct{}

func (m jarSha256PlanModifier) Description(_ context.Context) string {
	return "Recomputes jar_sha256 when the jar source changes so a rebuilt jar is detected."
}

func (m jarSha256PlanModifier) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

func (m jarSha256PlanModifier) PlanModifyString(
	ctx context.Context,
	req planmodifier.StringRequest,
	resp *planmodifier.StringResponse,
) {
	// The user set jar_sha256 explicitly: honour it verbatim.
	if !req.ConfigValue.IsNull() {
		resp.PlanValue = req.ConfigValue

		return
	}

	var plan, state PluginModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create: no prior state. Prefer hashing a known local jar_path so a plan
	// shows a concrete digest; otherwise leave it unknown.
	if req.State.Raw.IsNull() {
		resp.PlanValue = hashKnownJarPath(plan)

		return
	}

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// The jar source changed: a re-upload is coming, so the digest is unknown
	// unless the new local jar_path is readable now.
	if !plan.JarPath.Equal(state.JarPath) || !plan.JarUrl.Equal(state.JarUrl) {
		if h := hashKnownJarPath(plan); !h.IsUnknown() {
			resp.PlanValue = h
		} else {
			resp.PlanValue = types.StringUnknown()
		}

		return
	}

	// jar_path unchanged but its file may have been rebuilt in place: rehash it.
	if h := hashKnownJarPath(plan); !h.IsUnknown() && !h.IsNull() {
		resp.PlanValue = h

		return
	}

	resp.PlanValue = req.StateValue
}

// hashKnownJarPath returns the digest of the local jar_path when it is set,
// known and readable; otherwise an unknown value.
func hashKnownJarPath(plan PluginModel) types.String {
	if plan.JarPath.IsNull() || plan.JarPath.IsUnknown() {
		return types.StringUnknown()
	}

	sum, err := hashFile(plan.JarPath.ValueString())
	if err != nil {
		return types.StringUnknown()
	}

	return types.StringValue(sum)
}
