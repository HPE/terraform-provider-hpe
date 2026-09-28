// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package user

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/wodigest"
)

// Change detection for the write-only attributes password_wo,
// linux_password_wo and windows_password_wo.
//
// A write-only value is never in state, so a change to it cannot produce a
// plan diff. Digests of the configured values are kept in private state
// (wodigest); when they differ from what was last applied, ModifyPlan marks the
// computed write_only_revision unknown, which plans an in-place update, and
// Update re-sends the values. The deprecated *_wo_version markers keep working
// for the whole 2.2 line: a change to one is still treated as a change.

// woKey is the private-state key for this resource's write-only digests.
const woKey = "user.wo_digests"

// woValues returns the configured write-only values keyed by attribute name,
// as they would be sent to the API. A null value is nil, so it hashes
// differently from an empty string.
func woValues(config UserModel) map[string]any {
	str := func(v types.String) any {
		if v.IsNull() || v.IsUnknown() {
			return nil
		}

		return v.ValueString()
	}

	return map[string]any{
		"password_wo":         str(config.PasswordWo),
		"linux_password_wo":   str(config.LinuxPasswordWo),
		"windows_password_wo": str(config.WindowsPasswordWo),
	}
}

// woAttrs lists the write-only attributes, for the import ID.
var woAttrs = []string{"password_wo", "linux_password_wo", "windows_password_wo"}

// woDigests returns the digests of the configured write-only values. None of
// them requires replacement when changed, so there is no immutable subset.
func woDigests(config UserModel) (wodigest.Digests, error) {
	return wodigest.ComputeValues(woValues(config), nil)
}

// storeWoDigests records the digests of the configured write-only values,
// merged over what was stored before: an attribute left out of the
// configuration is not managed, so its last digest is kept rather than
// forgotten. Called by Create and, unconditionally, by Update, so that
// neither an import sentinel nor a missing digest survives an apply.
func storeWoDigests(ctx context.Context, private wodigest.PrivateState, config UserModel) diag.Diagnostics {
	var diags diag.Diagnostics

	current, err := woDigests(config)
	if err != nil {
		diags.AddError("user: write-only digests", err.Error())

		return diags
	}

	stored, ok, ldiags := wodigest.Load(ctx, private, woKey)
	diags.Append(ldiags...)

	if ok && stored.Sentinel == "" {
		merged := wodigest.Attrs{}
		for k, v := range stored.Attrs {
			merged[k] = v
		}

		for k, v := range current.Attrs {
			merged[k] = v
		}

		current = wodigest.Compute(merged, nil)
	}

	diags.Append(wodigest.Store(ctx, private, woKey, current)...)

	return diags
}

// markerChanged reports whether any deprecated *_wo_version marker differs
// between state and plan -- including a marker removed from the configuration,
// which re-sends the value once (harmlessly) and is then gone.
func markerChanged(plan, state UserModel) bool {
	return !plan.PasswordWoVersion.Equal(state.PasswordWoVersion) ||
		!plan.LinuxPasswordWoVersion.Equal(state.LinuxPasswordWoVersion) ||
		!plan.WindowsPasswordWoVersion.Equal(state.WindowsPasswordWoVersion)
}

// ModifyPlan plans an in-place update when a write-only value changed (or a
// deprecated marker did), by marking write_only_revision unknown. Create and
// destroy are left alone: Create always sends the values.
func (r *Resource) ModifyPlan(
	ctx context.Context,
	req resource.ModifyPlanRequest,
	resp *resource.ModifyPlanResponse,
) {
	if req.Plan.Raw.IsNull() || req.State.Raw.IsNull() {
		return
	}

	var config, plan, state UserModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	current, err := woDigests(config)
	if err != nil {
		resp.Diagnostics.AddError("user: write-only digests", err.Error())

		return
	}

	stored, ok, diags := wodigest.Load(ctx, req.Private, woKey)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	decision := wodigest.DecideConfigured(stored, ok, current.Attrs, nil)
	marker := markerChanged(plan, state)

	if ok && stored.Sentinel == wodigest.ImportSentinel && !marker {
		// Imported without digests, and no apply has recorded any yet: the
		// configured values are accepted as applied, and a change to them
		// cannot be seen. Say so on every plan until an update clears it.
		resp.Diagnostics.AddAttributeWarning(path.Root("password_wo"),
			"Write-only values of an imported user are taken as already applied",
			"This user was imported without the digests of its write-only values, so a change "+
				"to password_wo, linux_password_wo or windows_password_wo cannot be detected until the "+
				"user is next updated. To have changes detected from the start, import with the "+
				"values' sha256 digests in the import ID: \"<id>#password_wo=<sha256 hex>\". To send "+
				"the configured values now, change any other attribute, or set password_wo_version.")

		return
	}

	tflog.Debug(ctx, "user: write-only change detection", map[string]any{
		"digests_stored": ok,
		"sentinel":       stored.Sentinel,
		"decision":       decision.String(),
		"marker_changed": marker,
	})

	if decision == wodigest.Unchanged && !marker {
		return
	}

	// A change or a seed: plan the update. Nothing else the re-send changes
	// is computed on this resource, so write_only_revision is the only
	// attribute to mark unknown.
	resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx,
		path.Root("write_only_revision"), types.Int64Unknown())...)
}

// newRevision is the write_only_revision value for an apply that sent a
// write-only value.
func newRevision() types.Int64 {
	return types.Int64Value(time.Now().Unix())
}

// importWriteOnlyDigests records, at import, what is known about the write-only
// values. With digests in the import ID ("<id>#password_wo=<sha256 hex>&...",
// wodigest.ParseImportID) the resource stores the digests it would have
// computed from that configuration, so the plan after import is clean and the
// next change to a value is detected. Without them it stores the import
// sentinel: the plan after import is still clean, but a change to a write-only
// value is not detected until the next update (ModifyPlan warns).
func importWriteOnlyDigests(ctx context.Context, resp *resource.ImportStateResponse, attrs wodigest.Attrs) diag.Diagnostics {
	var diags diag.Diagnostics

	// write_only_revision is not in the API. Give the imported user one so
	// the plan after import does not show it as unknown.
	diags.Append(resp.State.SetAttribute(ctx, path.Root("write_only_revision"), newRevision())...)

	if attrs == nil {
		diags.Append(wodigest.StoreImportSentinel(ctx, resp.Private, woKey)...)

		return diags
	}

	diags.Append(wodigest.Store(ctx, resp.Private, woKey, wodigest.Compute(attrs, nil))...)

	return diags
}
