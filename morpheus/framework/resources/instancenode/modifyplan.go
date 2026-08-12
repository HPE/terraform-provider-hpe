// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package instancenode

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// ModifyPlan implements resource.ResourceWithModifyPlan.
//
// Best-effort pre-flight. instance_id is unknown on the run that creates the
// instance and its nodes together, so this cannot be the only check — Create
// repeats it authoritatively. Where the instance already exists, this surfaces
// the error at plan time, before anything is modified.
//
// The check only fires when resource_pool_id is set. If resource_pool_id is
// null (virtual-instance path), no metal check is needed.
func (r *Resource) ModifyPlan(
	ctx context.Context,
	req resource.ModifyPlanRequest,
	resp *resource.ModifyPlanResponse,
) {
	// Destroy — nothing to check.
	if req.Plan.Raw.IsNull() {
		return
	}

	var plan instanceNodeModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// First apply: the instance does not exist yet.
	if plan.InstanceID.IsUnknown() || plan.InstanceID.IsNull() {
		return
	}

	// No resource_pool_id — virtual path, no metal check needed.
	if plan.ResourcePoolID.IsNull() {
		return
	}

	// resource_pool_id is unknown (e.g. from a data source) — cannot
	// determine yet whether the check applies; defer to Create.
	if plan.ResourcePoolID.IsUnknown() {
		return
	}

	client, err := r.NewClient(ctx)
	if err != nil {
		// Do not fail the plan on a client error.
		tflog.Debug(ctx, "instance provision-type pre-flight skipped: client error",
			map[string]any{
				"instance_id": plan.InstanceID.ValueInt64(),
				"error":       err.Error(),
			},
		)

		return
	}

	instanceID := plan.InstanceID.ValueInt64()

	getResp, _, getErr := client.InstancesAPI.GetInstance(ctx, instanceID).Execute()
	if getErr != nil {
		// Do not fail the plan on a lookup failure. An unreachable or
		// unauthenticated appliance must not make `terraform plan` unusable,
		// and Create performs this check authoritatively.
		tflog.Debug(ctx, "instance provision-type pre-flight skipped",
			map[string]any{
				"instance_id": instanceID,
				"error":       getErr.Error(),
			},
		)

		return
	}

	if getResp == nil || getResp.Instance == nil {
		return
	}

	code, ok := provisionTypeCode(getResp.Instance)
	if ok && code != metalProvisionTypeCode {
		resp.Diagnostics.AddError(
			notMetalSummary,
			notMetalDetail(instanceID, code),
		)
	}
}
