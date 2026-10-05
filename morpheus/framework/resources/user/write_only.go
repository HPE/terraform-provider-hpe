// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package user

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/wodigest"
)

// Change detection for the write-only attributes password_wo,
// linux_password_wo and windows_password_wo: see wodigest. A change to one of
// them marks write_only_revision unknown, which plans an in-place update in
// which Update re-sends the configured values. The deprecated *_wo_version
// markers keep working for the whole 2.2 line.
var wo = wodigest.Tracker{
	Resource: "user",
	Attributes: map[string]string{
		"password_wo":         "password_wo_version",
		"linux_password_wo":   "linux_password_wo_version",
		"windows_password_wo": "windows_password_wo_version",
	},
}

// ModifyPlan plans an in-place update when a write-only value changed (or a
// deprecated marker did). With require_password_change set, the re-send of
// password_wo also marks the password expired, so the computed
// password_expired is planned as unknown then; a write-only change leaves the
// proposed state equal to the prior one, so the framework would otherwise
// keep the prior value and reject the apply as inconsistent (MORPH-18374).
func (r *Resource) ModifyPlan(
	ctx context.Context,
	req resource.ModifyPlanRequest,
	resp *resource.ModifyPlanResponse,
) {
	var alsoUnknown []path.Path

	if !req.Plan.Raw.IsNull() {
		var require types.Bool
		resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("require_password_change"), &require)...)

		if require.ValueBool() {
			alsoUnknown = append(alsoUnknown, path.Root("password_expired"))
		}
	}

	wo.ModifyPlan(ctx, req, resp, alsoUnknown...)
}
