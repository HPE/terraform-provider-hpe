// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkdomain

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var _ resource.ResourceWithModifyPlan = &Resource{}

// ModifyPlan handles the write-only value: a changed domain_password_wo (or its
// deprecated marker) marks write_only_revision unknown, planning an in-place
// update that re-sends the configured value. See write_only.go.
func (r *Resource) ModifyPlan(
	ctx context.Context,
	req resource.ModifyPlanRequest,
	resp *resource.ModifyPlanResponse,
) {
	wo.ModifyPlan(ctx, req, resp)
}
