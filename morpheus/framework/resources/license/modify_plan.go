// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package license

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/wodigest"
)

// ModifyPlan detects a change to the write-only key from the digest kept in
// private state and, when it changed, plans an in-place update by marking the
// key id and every computed detail unknown (the new key has a different key id
// and different details, and Terraform checks plan/apply consistency on all of
// them).
func (r *Resource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// Nothing to do on destroy or create.
	if req.Plan.Raw.IsNull() || req.State.Raw.IsNull() {
		return
	}

	// key is write-only, so it is null in the plan: read it from config.
	var config LicenseModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// An unknown config key (an ephemeral/unresolved value) is treated as a
	// change: plan the update.
	configured := ""
	if !config.Key.IsUnknown() {
		configured = keyDigest(config.Key.ValueString())
	}

	stored, ok := storedDigest(ctx, req.Private, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	decision := wodigest.Update
	if !config.Key.IsUnknown() {
		decision = planDecision(stored, ok, configured)
	}

	tflog.Debug(ctx, "license write-only key plan decision", map[string]any{
		"decision": decision.String(),
	})

	if decision == wodigest.Unchanged {
		return
	}

	// Plan an in-place update: the new key has a different key id and details.
	resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("id"), types.StringUnknown())...)
	resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("product_tier"), types.StringUnknown())...)
	resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("account_name"), types.StringUnknown())...)
	resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("start_date"), types.StringUnknown())...)
	resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("end_date"), types.StringUnknown())...)
	resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("support_level"), types.StringUnknown())...)
	resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("max_instances"), types.Int64Unknown())...)
	resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("max_memory"), types.Int64Unknown())...)
	resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("max_storage"), types.Int64Unknown())...)
	resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("hard_limit"), types.BoolUnknown())...)
	resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("free_trial"), types.BoolUnknown())...)
	resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("multi_tenant"), types.BoolUnknown())...)
}
