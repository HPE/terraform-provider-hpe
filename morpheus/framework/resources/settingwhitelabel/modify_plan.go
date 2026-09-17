// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package settingwhitelabel

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/tenancy"
)

var _ resource.ResourceWithModifyPlan = &settingWhitelabelResource{}

// ModifyPlan rejects, at plan time, appliance_name set by a non-master caller.
// The whitelabel appliance_name is a master-account-only setting: Morpheus
// silently discards it for sub-tenant users, so the value the caller planned is
// never stored and the read-back is null (MORPH-16546). Failing fast avoids a
// confusing "inconsistent result" on apply.
//
// Nothing to validate on destroy (Plan.Raw is null). Fail-open on an inability
// to determine the caller's tenancy.
func (r *settingWhitelabelResource) ModifyPlan(
	ctx context.Context,
	req resource.ModifyPlanRequest,
	resp *resource.ModifyPlanResponse,
) {
	if req.Plan.Raw.IsNull() {
		return
	}

	tenancy.GuardMasterOnlyAttribute(
		ctx, r, req.Plan, path.Root("appliance_name"), &resp.Diagnostics)
}
