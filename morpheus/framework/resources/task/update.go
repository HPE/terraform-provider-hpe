package task

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"

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

	updateRequest := sdk.NewUpdateTasksRequestWithDefaults()
	updateTask := sdk.NewUpdateTasksRequestTaskWithDefaults()

	// allow_custom_config
	if !plan.AllowCustomConfig.IsNull() && !plan.AllowCustomConfig.IsUnknown() {
		updateTask.SetAllowCustomConfig(plan.AllowCustomConfig.ValueBool())
	}

	// code
	if !plan.Code.IsNull() && !plan.Code.IsUnknown() {
		updateTask.SetCode(plan.Code.ValueString())
	}

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

	updateTask.SetTaskOptions(taskOptions)

	// execute_target
	if !plan.ExecuteTarget.IsNull() && !plan.ExecuteTarget.IsUnknown() {
		updateTask.SetExecuteTarget(plan.ExecuteTarget.ValueString())
	}

	// labels
	if !plan.Labels.IsNull() && !plan.Labels.IsUnknown() {
		labels, err := convert.SetToStrSlice(plan.Labels)
		if err != nil {
			resp.Diagnostics.AddError(
				"create task resource",
				"task "+plan.Name.ValueString()+": failed to parse label: "+err.Error(),
			)

			return
		}

		updateTask.SetLabels(labels)
	}

	// name
	if !plan.Name.IsNull() && !plan.Name.IsUnknown() {
		updateTask.SetName(plan.Name.ValueString())
	}

	// result_type
	if !plan.ResultType.IsNull() && !plan.ResultType.IsUnknown() {
		updateTask.SetResultType(plan.ResultType.ValueString())
	}

	// retry_count
	if !plan.RetryCount.IsNull() && !plan.RetryCount.IsUnknown() {
		updateTask.SetRetryCount(plan.RetryCount.ValueInt64())
	}

	// retry_delay_seconds
	if !plan.RetryDelaySeconds.IsNull() && !plan.RetryDelaySeconds.IsUnknown() {
		updateTask.SetRetryDelaySeconds(plan.RetryDelaySeconds.ValueInt64())
	}

	// retryable
	if !plan.Retryable.IsNull() && !plan.Retryable.IsUnknown() {
		updateTask.SetRetryable(plan.Retryable.ValueBool())
	}

	// task_type_code
	if !plan.TaskTypeCode.IsNull() && !plan.TaskTypeCode.IsUnknown() {
		updateTask.SetTaskType(*sdk.NewUpdateTasksRequestTaskTaskType(
			plan.TaskTypeCode.ValueString(),
		))
	}

	// visibility
	if !plan.Visibility.IsNull() && !plan.Visibility.IsUnknown() {
		updateTask.SetVisibility(
			plan.Visibility.ValueString(),
		)
	}

	updateRequest.SetTask(*updateTask)
	taskResp, httpResp, err := client.AutomationAPI.UpdateTasks(ctx, state.Id.ValueInt64()).
		UpdateTasksRequest(*updateRequest).
		Execute()
	if err != nil || httpResp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError("error updating task", errfmt.ErrMsg(err, httpResp))

		return
	}

	// set the ID value in state
	plan.Id = convert.Int64ToType(taskResp.Task.Id)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	state, diag := getTaskAsState(ctx, *taskResp.Task.Id, client, plan)
	if resp.Diagnostics.Append(diag...); resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
