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

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"

	"github.com/HPE/terraform-provider-hpe/internal/framework/customtypes"
	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/convert"
	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/errors"
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
			fmt.Sprintf("task %d GET failed: ", id)+errors.ErrMsg(err, httpResp),
		)

		return state, diags
	}

	task := taskResp.GetTask()
	// code
	state.Code = convert.StrToType(task.Code.Get())

	var typeCode string
	if task.GetTaskType().Code != nil {
		typeCode = *task.GetTaskType().Code
	}
	// config
	state.Config = basetypes.NewDynamicNull()
	if task.TaskOptions.MapmapOfStringAny != nil {
		o, err := convert.MapToDynamic(ctx, *task.TaskOptions.MapmapOfStringAny)
		if err != nil {
			diags.AddError("populate task resource", err.Error())
		}

		state.Config = o
	}

	// config_conditional_workflow_task
	if task.TaskOptions.ConditionalWorkflowTaskConfig != nil && typeCode == "conditionalWorkflow" {
		config := task.TaskOptions.ConditionalWorkflowTaskConfig

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
	respLabels := task.GetLabels()

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
	state.RetryDelaySeconds = convert.Int64ToType(task.RetryDelaySeconds)

	// retryable
	state.Retryable = convert.BoolToType(task.Retryable)

	// task_type_code
	state.TaskTypeCode = convert.StrToType(task.TaskType.Code)

	// visibility
	state.Visibility = convert.StrToType(task.Visibility)

	return state, diags
}
