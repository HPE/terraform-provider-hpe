package task

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

func (r *Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan TaskModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError("creating client failed", err.Error())

		return
	}

	addTaskReq := sdk.NewAddTasksRequestWithDefaults()

	// allow_custom_config
	if !plan.AllowCustomConfig.IsNull() && !plan.AllowCustomConfig.IsUnknown() {
		addTaskReq.Task.SetAllowCustomConfig(plan.AllowCustomConfig.ValueBool())
	}

	// code
	if !plan.Code.IsNull() && !plan.Code.IsUnknown() {
		addTaskReq.Task.SetCode(plan.Code.ValueString())
	}

	taskOptionsSet := false

	// config
	taskOptions := &sdk.AddTasksRequestTaskTaskOptions{}
	if !plan.Config.IsNull() && !plan.Config.IsUnknown() {
		configValue := plan.Config.UnderlyingValue()
		configAny, err := convert.ValueToAny(ctx, configValue)
		if err != nil {
			resp.Diagnostics.AddError(
				"create task resource",
				"task: failed to convert config: "+
					err.Error(),
			)

			return
		}

		configDataMap, ok := configAny.(map[string]any)
		if !ok {
			resp.Diagnostics.AddError(
				"error creating task",
				"could not parse config value",
			)

			return
		}

		taskOptions.MapmapOfStringAny = &configDataMap
		taskOptionsSet = true
	}

	// config_conditional_workflow
	if !plan.ConfigConditionalWorkflow.IsNull() && !plan.ConfigConditionalWorkflow.IsUnknown() {
		conditionalWorkflow := &sdk.ConditionalWorkflowTaskConfig{}
		trimmed := plan.ConfigConditionalWorkflow.ConditionalScript.ValueString()
		conditionalWorkflow.ConditionalScript = &trimmed

		conditionalWorkflow.IfOperationalWorkflowId = plan.ConfigConditionalWorkflow.
			IfOperationalWorkflowId.ValueInt64Pointer()

		conditionalWorkflow.IfOperationalWorkflowName = plan.ConfigConditionalWorkflow.
			IfOperationalWorkflowName.ValueStringPointer()

		conditionalWorkflow.ElseOperationalWorkflowId = plan.ConfigConditionalWorkflow.
			ElseOperationalWorkflowId.ValueInt64Pointer()

		conditionalWorkflow.ElseOperationalWorkflowName = plan.ConfigConditionalWorkflow.
			ElseOperationalWorkflowName.ValueStringPointer()

		taskOptions.ConditionalWorkflowTaskConfig = conditionalWorkflow

		taskOptionsSet = true
	}

	if taskOptionsSet {
		addTaskReq.Task.SetTaskOptions(*taskOptions)
	}

	// execute_target
	if !plan.ExecuteTarget.IsNull() && !plan.ExecuteTarget.IsUnknown() {
		addTaskReq.Task.SetExecuteTarget(plan.ExecuteTarget.ValueString())
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

		addTaskReq.Task.SetLabels(labels)
	}

	// name
	if !plan.Name.IsNull() && !plan.Name.IsUnknown() {
		addTaskReq.Task.SetName(plan.Name.ValueString())
	}

	// result_type
	if !plan.ResultType.IsNull() && !plan.ResultType.IsUnknown() {
		addTaskReq.Task.SetResultType(plan.ResultType.ValueString())
	}

	// retry_count
	if !plan.RetryCount.IsNull() && !plan.RetryCount.IsUnknown() {
		addTaskReq.Task.SetRetryCount(plan.RetryCount.ValueInt64())
	}

	// retry_delay_seconds
	if !plan.RetryDelaySeconds.IsNull() && !plan.RetryDelaySeconds.IsUnknown() {
		addTaskReq.Task.SetRetryDelaySeconds(plan.RetryCount.ValueInt64())
	}

	// retryable
	if !plan.Retryable.IsNull() && !plan.Retryable.IsUnknown() {
		addTaskReq.Task.SetRetryable(plan.Retryable.ValueBool())
	}

	// task_type_code
	if !plan.TaskTypeCode.IsNull() && !plan.TaskTypeCode.IsUnknown() {
		addTaskReq.Task.TaskType.SetCode(plan.TaskTypeCode.ValueString())
	}

	// visibility
	if !plan.Visibility.IsNull() && !plan.Visibility.IsUnknown() {
		addTaskReq.Task.SetVisibility(plan.Visibility.ValueString())
	}

	// send the API request here
	taskResp, httpResp, err := client.AutomationAPI.AddTasks(ctx).
		AddTasksRequest(*addTaskReq).Execute()
	if err != nil || httpResp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError("error creating task", errfmt.ErrMsg(err, httpResp))

		return
	}

	plan.Id = convert.Int64ToType(taskResp.Task.Id)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	state, diag := getTaskAsState(ctx, plan.Id.ValueInt64(), client, plan)
	if resp.Diagnostics.Append(diag...); resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
