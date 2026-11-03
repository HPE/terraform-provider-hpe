package task

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/clearable"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

func (r *Resource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	var plan, state, config TaskModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError("creating client failed", err.Error())

		return
	}

	updateRequest := &sdk.UpdateTasksRequest{}
	// Attributes removed from the configuration since the last apply, to be
	// cleared. Read before Record, below, overwrites what it reads.
	pending, pdiags := clearable.Compute(ctx, req.Config, req.State, req.Private, clearableAttributes)
	resp.Diagnostics.Append(pdiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateTask := &sdk.UpdateTasksRequestTask{}

	// allow_custom_config
	if !plan.AllowCustomConfig.IsNull() && !plan.AllowCustomConfig.IsUnknown() {
		updateTask.AllowCustomConfig = plan.AllowCustomConfig.ValueBoolPointer()
	}

	// code
	clearable.NullableString(pending, "code", plan.Code, &updateTask.Code)

	taskOptions := sdk.UpdateTasksRequestTaskTaskOptions{}

	// config
	if !plan.Config.IsNull() && !plan.Config.IsUnknown() {
		configValue := plan.Config.UnderlyingValue()
		configAny, err := convert.ValueToAny(ctx, configValue)
		if err != nil {
			resp.Diagnostics.AddError(
				"update task resource",
				"task: failed to convert config: "+
					err.Error(),
			)

			return
		}

		configDataMap, ok := configAny.(map[string]any)
		if !ok {
			resp.Diagnostics.AddError(
				"error updating task",
				"could not parse config value",
			)

			return
		}

		taskOptions.MapmapOfStringAny = &configDataMap
	}

	// config_conditional_workflow
	if !plan.ConfigConditionalWorkflow.IsNull() && !plan.ConfigConditionalWorkflow.IsUnknown() {
		conditionalWorkflow := &sdk.ConditionalWorkflowTaskConfig3{}
		conditionalWorkflow.ConditionalScript = plan.ConfigConditionalWorkflow.
			ConditionalScript.ValueStringPointer()
		conditionalWorkflow.IfOperationalWorkflowId = plan.ConfigConditionalWorkflow.
			IfOperationalWorkflowId.ValueInt64Pointer()
		conditionalWorkflow.IfOperationalWorkflowName = plan.ConfigConditionalWorkflow.
			IfOperationalWorkflowName.ValueStringPointer()
		if !plan.ConfigConditionalWorkflow.ElseOperationalWorkflowName.IsNull() &&
			!plan.ConfigConditionalWorkflow.ElseOperationalWorkflowName.IsUnknown() {
			conditionalWorkflow.ElseOperationalWorkflowId = plan.ConfigConditionalWorkflow.
				ElseOperationalWorkflowId.ValueInt64Pointer()
		}
		if !plan.ConfigConditionalWorkflow.ElseOperationalWorkflowName.IsNull() &&
			!plan.ConfigConditionalWorkflow.ElseOperationalWorkflowName.IsUnknown() {
			conditionalWorkflow.ElseOperationalWorkflowName = plan.ConfigConditionalWorkflow.
				ElseOperationalWorkflowName.ValueStringPointer()
		}

		taskOptions.ConditionalWorkflowTaskConfig3 = conditionalWorkflow
	}

	updateTask.TaskOptions = &taskOptions

	// execute_target
	if !plan.ExecuteTarget.IsNull() && !plan.ExecuteTarget.IsUnknown() {
		updateTask.ExecuteTarget = plan.ExecuteTarget.ValueStringPointer()
	}

	// labels: the planned list, or an empty list when being cleared
	if labels, ok, ldiags := pending.Strings(ctx, "labels", plan.Labels); ok {
		resp.Diagnostics.Append(ldiags...)
		if resp.Diagnostics.HasError() {
			return
		}

		updateTask.Labels = labels
	}

	// name
	if !plan.Name.IsNull() && !plan.Name.IsUnknown() {
		updateTask.Name = plan.Name.ValueStringPointer()
	}

	// result_type
	clearable.NullableString(pending, "result_type", plan.ResultType, &updateTask.ResultType)

	// retry_count
	clearable.NullableInt64(pending, "retry_count", plan.RetryCount, &updateTask.RetryCount)

	// retry_delay_seconds
	clearable.NullableInt64(pending, "retry_delay_seconds", plan.RetryDelaySeconds, &updateTask.RetryDelaySeconds)

	// retryable
	if !plan.Retryable.IsNull() && !plan.Retryable.IsUnknown() {
		updateTask.Retryable = plan.Retryable.ValueBoolPointer()
	}

	// task_type_code
	if !plan.TaskTypeCode.IsNull() && !plan.TaskTypeCode.IsUnknown() {
		updateTask.TaskType = &sdk.UpdateTasksRequestTaskTaskType{
			Code: plan.TaskTypeCode.ValueString(),
		}
	}

	// visibility
	if !plan.Visibility.IsNull() && !plan.Visibility.IsUnknown() {
		updateTask.Visibility = plan.Visibility.ValueStringPointer()
	}

	updateRequest.Task = *updateTask
	taskResp, httpResp, err := client.AutomationAPI.UpdateTasks(ctx, state.Id.ValueInt64()).
		UpdateTasksRequest(*updateRequest).
		Execute()
	if err != nil || httpResp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError("error updating task", errfmt.ErrMsg(err, httpResp))

		return
	}

	if taskResp.Task == nil {
		resp.Diagnostics.AddError("API returned nil", "Task is nil in the response")

		return
	}

	if taskResp.Task.Id == nil {
		resp.Diagnostics.AddError("API returned nil", "Task ID is nil in the response")

		return
	}

	// set the ID value in state
	plan.Id = convert.Int64ToType(taskResp.Task.Id)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	state, found, diag := getTaskAsState(ctx, *taskResp.Task.Id, client, plan)
	if !found && !diag.HasError() {
		diag.AddError("read task", fmt.Sprintf("task %d not found after update", *taskResp.Task.Id))
	}
	if resp.Diagnostics.Append(diag...); resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(clearable.Record(ctx, req.Config, resp.Private, clearableAttributes)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
