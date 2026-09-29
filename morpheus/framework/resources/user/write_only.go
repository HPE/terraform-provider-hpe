// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package user

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"

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
// deprecated marker did). Nothing else the re-send changes is computed on this
// resource, so write_only_revision is the only attribute marked unknown.
func (r *Resource) ModifyPlan(
	ctx context.Context,
	req resource.ModifyPlanRequest,
	resp *resource.ModifyPlanResponse,
) {
	wo.ModifyPlan(ctx, req, resp)
}
