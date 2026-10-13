// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package taskshellscript

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/taskscript"
)

var _ resource.ResourceWithModifyPlan = &Resource{}

// ModifyPlan handles the write-only value: a changed remote_target_password_wo
// (or its deprecated marker) marks write_only_revision unknown, planning an
// in-place update that re-sends the configured value.
func (r *Resource) ModifyPlan(
	ctx context.Context,
	req resource.ModifyPlanRequest,
	resp *resource.ModifyPlanResponse,
) {
	taskscript.WO.ModifyPlan(ctx, req, resp)
}
