// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package provisioninglicense

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/wodigest"
)

// Change detection for the write-only attribute license_key_wo: see
// wodigest. A change marks write_only_revision unknown, which plans an in-place
// update in which Update re-sends the configured value. The deprecated
// *_version marker keeps working for the whole 2.2 line.
var wo = wodigest.Tracker{
	Resource: "provisioning_license",
	Attributes: map[string]string{
		"license_key_wo": "license_key_wo_version",
	},
}

// ModifyPlan plans an in-place update when a write-only value changed (or a
// deprecated marker did).
func (r *provisioningLicenseResource) ModifyPlan(
	ctx context.Context,
	req resource.ModifyPlanRequest,
	resp *resource.ModifyPlanResponse,
) {
	wo.ModifyPlan(ctx, req, resp)
}
