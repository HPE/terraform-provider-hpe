package task

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
	"github.com/HPE/terraform-provider-hpe/utils/customtypes"
	"github.com/HPE/terraform-provider-hpe/utils/schemadefaults"
)

func (r *Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data TaskModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError("creating client failed", err.Error())

		return
	}

	state, diag := getTaskAsState(ctx, data.Id.ValueInt64(), client, data)
	if resp.Diagnostics.Append(diag...); resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Fill any schema-declared default the API omitted (null in state) so an
	// imported resource does not plan a change nobody made. MORPH-16192.
	resp.Diagnostics.Append(
		schemadefaults.Apply(ctx, TaskResourceSchema(ctx), &resp.State)...,
	)
}

func getTaskAsState(
	ctx context.Context,
	id int64,
	client *sdk.APIClient,
	plan TaskModel,
) (TaskModel, diag.Diagnostics) {
	var state TaskModel
	var diags diag.Diagnostics

	taskResp, httpResp, err := client.AutomationAPI.GetTasks(ctx, id).Execute()
	if err != nil || httpResp.StatusCode != http.StatusOK {
		diags.AddError(
			"populate task resource",
			fmt.Sprintf("task %d GET failed: ", id)+errfmt.ErrMsg(err, httpResp),
		)

		return state, diags
	}

	task := taskResp.Task
	if task == nil {
		diags.AddError("populate task resource", fmt.Sprintf("task %d: response task is nil", id))

		return state, diags
	}

	// allow_custom_config
	state.AllowCustomConfig = convert.BoolToType(task.AllowCustomConfig)

	// code
	state.Code = convert.StrToType(task.Code.Get())

	var typeCode string
	if task.TaskType != nil && task.TaskType.Code != nil {
		typeCode = *task.TaskType.Code
	}
	// config
	state.Config = basetypes.NewDynamicNull()
	if task.TaskOptions != nil && task.TaskOptions.MapmapOfStringAny != nil {
		o, err := convert.MapToDynamic(ctx, *task.TaskOptions.MapmapOfStringAny)
		if err != nil {
			diags.AddError("populate task resource", err.Error())
		}

		state.Config = o
	}

	// config_conditional_workflow_task
	if task.TaskOptions != nil && task.TaskOptions.ConditionalWorkflowTaskConfig2 != nil && typeCode == "conditionalWorkflow" {
		config := task.TaskOptions.ConditionalWorkflowTaskConfig2

		if config.ConditionalScript == nil {
			state.ConfigConditionalWorkflow.ConditionalScript = customtypes.NewTrimmedStringNull()
		} else {
			state.ConfigConditionalWorkflow.ConditionalScript = customtypes.NewTrimmedStringValue(
				strings.TrimSpace(*config.ConditionalScript),
			)
		}

		state.ConfigConditionalWorkflow.IfOperationalWorkflowId = convert.Int64ToType(
			config.IfOperationalWorkflowId,
		)

		state.ConfigConditionalWorkflow.IfOperationalWorkflowName = convert.StrToType(
			config.IfOperationalWorkflowName,
		)

		state.ConfigConditionalWorkflow.ElseOperationalWorkflowId = convert.Int64ToType(
			config.ElseOperationalWorkflowId,
		)

		state.ConfigConditionalWorkflow.ElseOperationalWorkflowName = convert.StrToType(
			config.ElseOperationalWorkflowName,
		)

		state.ConfigConditionalWorkflow.state = attr.ValueStateKnown
	}

	// execute_target
	state.ExecuteTarget = convert.StrToType(task.ExecuteTarget)

	// id
	state.Id = convert.Int64ToType(task.Id)

	// labels
	respLabels := task.Labels

	labels, err := convert.SetToStrSlice(plan.Labels)
	if err != nil {
		diags.AddError(
			"populate task resource",
			"could not parse a slice of labels",
		)

		return state, diags
	}

	// Morpheus API may change the casing of the labels, to avoid Terraform
	// throwing a gasket we convert the casing of labels to be as specified
	// by the user.
	for _, label := range labels {
		for i, respLabel := range respLabels {
			if strings.EqualFold(label, respLabel) {
				if label != respLabel {
					respLabels[i] = label
				}
			}
		}
	}

	state.Labels = convert.StrSliceToSet(respLabels)

	// name
	state.Name = convert.StrToType(task.Name)

	// result_type
	state.ResultType = convert.StrToType(task.ResultType.Get())

	// retry_count
	state.RetryCount = convert.Int64ToType(task.RetryCount)

	// retry_delay_seconds
	state.RetryDelaySeconds = retryDelaySecondsToState(plan.RetryDelaySeconds, task.RetryDelaySeconds)

	// retryable
	state.Retryable = convert.BoolToType(task.Retryable)

	// task_type_code
	state.TaskTypeCode = taskTypeCodeToState(task.TaskType)

	// visibility
	state.Visibility = convert.StrToType(task.Visibility)

	return state, diags
}

// retryDelaySecondsToState preserves a non-zero planned retry_delay_seconds when the
// API omits the field or returns a different value, so this Optional+Computed attribute
// does not drift to null on read of an out-of-band change. It falls back to the nil-safe
// API value otherwise. A nil apiVal is safe: the guard short-circuits before the deref.
func retryDelaySecondsToState(planned basetypes.Int64Value, apiVal *int64) basetypes.Int64Value {
	seconds := planned.ValueInt64()
	if seconds != 0 && (apiVal == nil || seconds != *apiVal) {
		return planned
	}

	return convert.Int64ToType(apiVal)
}

// taskTypeCodeToState reads the task type code without dereferencing task.TaskType,
// which the API may omit. The result is a null string when taskType (or its code) is
// absent.
func taskTypeCodeToState(taskType *sdk.GetTasks200ResponseAllOfTaskTaskType) basetypes.StringValue {
	var code *string
	if taskType != nil {
		code = taskType.Code
	}

	return convert.StrToType(code)
}
