// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkserver

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/wodigest"
)

// Change detection for the write-only attributes service_password_wo and
// service_token_wo: see wodigest. A change marks write_only_revision unknown,
// which plans an in-place update in which Update re-sends the configured
// values. The deprecated *_version marker of service_password_wo keeps working
// for the whole 2.x line; service_token_wo never had one.
var wo = wodigest.Tracker{
	Resource: "network_server",
	Attributes: map[string]string{ //nolint:gosec // attribute names, not credentials
		"service_password_wo": "service_password_wo_version",
		"service_token_wo":    "",
	},
}

// ModifyPlan plans an in-place update when a write-only value changed (or a
// deprecated marker did). A secret-only change leaves the proposed plan equal
// to the prior state, so the framework keeps every computed value known; the
// re-send is an API update that advances last_updated, which is therefore
// planned as unknown too (or the apply would be rejected as inconsistent).
// The clearable attributes on this resource are all Optional-only, so they
// need no clearable.Plan wiring here: the plan carries their null and Update
// sends it.
func (r *networkServerResource) ModifyPlan(
	ctx context.Context,
	req resource.ModifyPlanRequest,
	resp *resource.ModifyPlanResponse,
) {
	wo.ModifyPlan(ctx, req, resp, path.Root("last_updated"))
}
