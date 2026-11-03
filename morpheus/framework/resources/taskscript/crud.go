// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package taskscript

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/cenkalti/backoff/v5"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/wodigest"
	"github.com/HPE/terraform-provider-hpe/utils/cleanup"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

// knownStr reports whether a string plan/config value is set (not null, not
// unknown) and non-empty.
func knownStr(v types.String) bool {
	return !v.IsNull() && !v.IsUnknown() && v.ValueString() != ""
}

// fileFields are the request "file" fields a plan supplies, mirroring the SDKv2
// resource: the content for a local script, the path for a URL, the repository
// id and ref for a repository. A value the plan does not set is omitted.
type fileFields struct {
	content, contentPath, contentRef *string
	repositoryID                     *int64
}

func planFileFields(plan Flat) fileFields {
	var f fileFields
	if knownStr(plan.ScriptContent) {
		f.content = plan.ScriptContent.ValueStringPointer()
	}
	if knownStr(plan.ScriptPath) {
		f.contentPath = plan.ScriptPath.ValueStringPointer()
	}
	if knownStr(plan.VersionRef) {
		f.contentRef = plan.VersionRef.ValueStringPointer()
	}
	if !plan.RepositoryId.IsNull() && !plan.RepositoryId.IsUnknown() {
		f.repositoryID = plan.RepositoryId.ValueInt64Pointer()
	}

	return f
}

func (f fileFields) toAdd(sourceType string) *sdk.AddTasksRequestTaskFile {
	file := &sdk.AddTasksRequestTaskFile{
		SourceType: sourceType, Content: f.content, ContentPath: f.contentPath, ContentRef: f.contentRef,
	}
	if f.repositoryID != nil {
		file.Repository = &sdk.AddTasksRequestTaskFileRepository{Id: f.repositoryID}
	}

	return file
}

func (f fileFields) toUpdate(sourceType string) *sdk.UpdateTasksRequestTaskFile {
	file := &sdk.UpdateTasksRequestTaskFile{
		SourceType: sourceType, Content: f.content, ContentPath: f.contentPath, ContentRef: f.contentRef,
	}
	if f.repositoryID != nil {
		file.Repository = &sdk.UpdateTasksRequestTaskFileRepository{Id: f.repositoryID}
	}

	return file
}

// buildTaskOptions builds the taskOptions map from the flattened plan. password
// is added by the caller (Create/Update) from config, which is where the
// write-only value lives. The flag (shell.sudo / winrm.elevated) is sent as
// "on" when set and omitted otherwise, exactly as the SDKv2 resource did.
func (k Kind) buildTaskOptions(plan Flat) map[string]any {
	opts := map[string]any{}

	if plan.Flag.ValueBool() {
		opts[k.FlagKey] = "on"
	}

	if knownStr(plan.RemoteTargetHost) {
		opts["host"] = plan.RemoteTargetHost.ValueString()
	}

	if knownStr(plan.RemoteTargetPort) {
		opts["port"] = plan.RemoteTargetPort.ValueString()
	}

	if knownStr(plan.RemoteTargetUsername) {
		opts["username"] = plan.RemoteTargetUsername.ValueString()
	}

	if k.HasLocalRepository {
		if knownStr(plan.LocalRepositoryId) {
			opts["localScriptGitId"] = plan.LocalRepositoryId.ValueString()
		}

		if knownStr(plan.LocalRepositoryRef) {
			opts["localScriptGitRef"] = plan.LocalRepositoryRef.ValueString()
		}
	}

	return opts
}

// Create sends the POST, reads the created task back, records the write-only
// digest and sets state. setState converts the resulting Flat into the caller's
// model and writes it to resp.State.
func (k Kind) Create(
	ctx context.Context,
	client *sdk.APIClient,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
	plan, config Flat,
	setState func(context.Context, Flat) diag.Diagnostics,
) {
	addTask := sdk.AddTasksRequestTask{
		Name:          plan.Name.ValueString(),
		ExecuteTarget: plan.ExecuteTarget.ValueString(),
		TaskType:      sdk.AddTasksRequestTaskTaskType{Code: k.TaskTypeCode},
	}

	if knownStr(plan.Code) {
		addTask.Code = plan.Code.ValueStringPointer()
	}

	if !plan.Visibility.IsNull() && !plan.Visibility.IsUnknown() {
		addTask.Visibility = plan.Visibility.ValueStringPointer()
	}

	if !plan.ResultType.IsNull() && !plan.ResultType.IsUnknown() {
		addTask.ResultType.Set(plan.ResultType.ValueStringPointer())
	}

	if !plan.Retryable.IsNull() && !plan.Retryable.IsUnknown() {
		addTask.Retryable = plan.Retryable.ValueBoolPointer()
	}

	if !plan.RetryCount.IsNull() && !plan.RetryCount.IsUnknown() {
		addTask.RetryCount = plan.RetryCount.ValueInt64Pointer()
	}

	if !plan.RetryDelaySeconds.IsNull() && !plan.RetryDelaySeconds.IsUnknown() {
		addTask.RetryDelaySeconds = plan.RetryDelaySeconds.ValueInt64Pointer()
	}

	if !plan.AllowCustomConfig.IsNull() && !plan.AllowCustomConfig.IsUnknown() {
		addTask.AllowCustomConfig = plan.AllowCustomConfig.ValueBoolPointer()
	}

	if !plan.Labels.IsNull() && !plan.Labels.IsUnknown() {
		labels, err := convert.SetToStrSlice(plan.Labels)
		if err != nil {
			resp.Diagnostics.AddError("create "+k.TypeName, "failed to parse labels: "+err.Error())

			return
		}

		addTask.Labels = labels
	}

	addTask.File = planFileFields(plan).toAdd(plan.SourceType.ValueString())

	opts := k.buildTaskOptions(plan)
	// On create, send the password only when configured (write-only first;
	// the deprecated plaintext attribute otherwise — ConflictsWith keeps them
	// exclusive).
	switch {
	case knownStr(config.RemoteTargetPasswordWo):
		opts["password"] = config.RemoteTargetPasswordWo.ValueString()
	case knownStr(plan.RemoteTargetPassword):
		opts["password"] = plan.RemoteTargetPassword.ValueString()
	}

	optsMap := opts
	addTask.TaskOptions = &sdk.AddTasksRequestTaskTaskOptions{MapmapOfStringAny: &optsMap}

	createResp, httpResp, err := client.AutomationAPI.AddTasks(ctx).
		AddTasksRequest(sdk.AddTasksRequest{Task: addTask}).Execute()
	if err != nil || httpResp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError("error creating "+k.TypeName, errfmt.ErrMsg(err, httpResp))

		return
	}

	if createResp.Task == nil || createResp.Task.Id == nil {
		resp.Diagnostics.AddError("error creating "+k.TypeName, "task id is nil in the response")

		return
	}

	id := *createResp.Task.Id

	state, found, diags := k.getAsState(ctx, id, client, plan)
	if !found && !diags.HasError() {
		diags.AddError("read "+k.TypeName, "task not found immediately after create")
	}
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		// The task exists: record its id as tainted so the next apply replaces
		// it instead of leaving it outside state.
		cleanup.TaintResourceState(ctx, cleanup.TaintResourceStateConfig{
			ResourceType: k.TypeName,
			ResourceID:   id,
			StateWriter:  &resp.State,
			Diagnostics:  &resp.Diagnostics,
		})

		return
	}

	state.RemoteTargetPasswordWoVersion = plan.RemoteTargetPasswordWoVersion
	state.RemoteTargetPassword = plan.RemoteTargetPassword
	state.WriteOnlyRevision = wodigest.NewRevision()
	resp.Diagnostics.Append(WO.Store(ctx, resp.Private, req.Config)...)

	resp.Diagnostics.Append(setState(ctx, state)...)
}

// Update sends the PUT, re-sending the password only when ModifyPlan marked a
// resend. For tasks, omitting password on update does NOT clear it (the SDKv2
// resource omitted it except on marker bump), so the key is simply omitted when
// not resending.
func (k Kind) Update(
	ctx context.Context,
	client *sdk.APIClient,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
	plan, currentState, config Flat,
	setState func(context.Context, Flat) diag.Diagnostics,
) {
	id := currentState.Id.ValueInt64()

	code := sdk.NewNullableString(nil)
	if knownStr(plan.Code) {
		code = sdk.NewNullableString(plan.Code.ValueStringPointer())
	}

	updTask := sdk.UpdateTasksRequestTask{
		Name:     plan.Name.ValueStringPointer(),
		Code:     *code,
		TaskType: &sdk.UpdateTasksRequestTaskTaskType{Code: k.TaskTypeCode},
	}

	if !plan.ExecuteTarget.IsNull() && !plan.ExecuteTarget.IsUnknown() {
		updTask.ExecuteTarget = plan.ExecuteTarget.ValueStringPointer()
	}

	if !plan.Visibility.IsNull() && !plan.Visibility.IsUnknown() {
		updTask.Visibility = plan.Visibility.ValueStringPointer()
	}

	if !plan.ResultType.IsNull() && !plan.ResultType.IsUnknown() {
		updTask.ResultType.Set(plan.ResultType.ValueStringPointer())
	}

	if !plan.Retryable.IsNull() && !plan.Retryable.IsUnknown() {
		updTask.Retryable = plan.Retryable.ValueBoolPointer()
	}

	if !plan.RetryCount.IsNull() && !plan.RetryCount.IsUnknown() {
		updTask.RetryCount.Set(plan.RetryCount.ValueInt64Pointer())
	}

	if !plan.RetryDelaySeconds.IsNull() && !plan.RetryDelaySeconds.IsUnknown() {
		updTask.RetryDelaySeconds.Set(plan.RetryDelaySeconds.ValueInt64Pointer())
	}

	if !plan.AllowCustomConfig.IsNull() && !plan.AllowCustomConfig.IsUnknown() {
		updTask.AllowCustomConfig = plan.AllowCustomConfig.ValueBoolPointer()
	}

	if !plan.Labels.IsNull() && !plan.Labels.IsUnknown() {
		labels, err := convert.SetToStrSlice(plan.Labels)
		if err != nil {
			resp.Diagnostics.AddError("update "+k.TypeName, "failed to parse labels: "+err.Error())

			return
		}

		updTask.Labels = labels
	}

	updTask.File = planFileFields(plan).toUpdate(plan.SourceType.ValueString())

	opts := k.buildTaskOptions(plan)

	resend := wodigest.Resend(ctx, req.Plan)
	switch {
	case resend && knownStr(config.RemoteTargetPasswordWo):
		opts["password"] = config.RemoteTargetPasswordWo.ValueString()
	case knownStr(plan.RemoteTargetPassword) && !plan.RemoteTargetPassword.Equal(currentState.RemoteTargetPassword):
		// Deprecated plaintext attribute: sent when it changes, as the SDKv2
		// resource did. Omitting password on a task update leaves it unchanged.
		opts["password"] = plan.RemoteTargetPassword.ValueString()
	}

	optsMap := opts
	updTask.TaskOptions = &sdk.UpdateTasksRequestTaskTaskOptions{MapmapOfStringAny: &optsMap}

	_, httpResp, err := client.AutomationAPI.UpdateTasks(ctx, id).
		UpdateTasksRequest(sdk.UpdateTasksRequest{Task: updTask}).Execute()
	if err != nil || httpResp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError("error updating "+k.TypeName, errfmt.ErrMsg(err, httpResp))

		return
	}

	// The PUT succeeded, so whatever was sent is what the appliance now holds:
	// record the digests before the read-back, so a failed read cannot leave a
	// stale digest that would re-send the password on the next plan.
	resp.Diagnostics.Append(WO.Store(ctx, resp.Private, req.Config)...)

	state, found, diags := k.getAsState(ctx, id, client, plan)
	if !found && !diags.HasError() {
		diags.AddError("read "+k.TypeName, "task not found immediately after update")
	}
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	state.RemoteTargetPasswordWoVersion = plan.RemoteTargetPasswordWoVersion
	state.RemoteTargetPassword = plan.RemoteTargetPassword
	state.WriteOnlyRevision = currentState.WriteOnlyRevision
	if resend {
		state.WriteOnlyRevision = wodigest.NewRevision()
	}

	resp.Diagnostics.Append(setState(ctx, state)...)
}

// Read reads the task back and sets state. Carries the write-only revision and
// marker from the prior state; a 404 removes the resource from state.
func (k Kind) Read(
	ctx context.Context,
	client *sdk.APIClient,
	resp *resource.ReadResponse,
	prior Flat,
	setState func(context.Context, Flat) diag.Diagnostics,
) {
	state, found, diags := k.getAsState(ctx, prior.Id.ValueInt64(), client, prior)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	if !found {
		tflog.Info(ctx, fmt.Sprintf("%s %d not found, removing from state", k.TypeName, prior.Id.ValueInt64()))
		resp.State.RemoveResource(ctx)

		return
	}

	state.RemoteTargetPasswordWoVersion = prior.RemoteTargetPasswordWoVersion
	state.RemoteTargetPassword = prior.RemoteTargetPassword
	state.WriteOnlyRevision = prior.WriteOnlyRevision

	resp.Diagnostics.Append(setState(ctx, state)...)
}

// Delete removes the task and polls until it is gone. A 404 is treated as
// success.
func (k Kind) Delete(
	ctx context.Context,
	client *sdk.APIClient,
	resp *resource.DeleteResponse,
	state Flat,
) {
	id := state.Id.ValueInt64()

	result, httpResp, err := client.AutomationAPI.RemoveTasks(ctx, id).Execute()
	if err := errfmt.CheckDeleted(result, err, httpResp); err != nil {
		resp.Diagnostics.AddError(
			"error deleting "+k.TypeName,
			fmt.Sprintf("%s %d DELETE failed: ", k.TypeName, id)+err.Error(),
		)

		return
	}

	waitForDeleted := func() (struct{}, error) {
		_, httpResp, err := client.AutomationAPI.GetTasks(ctx, id).Execute()
		if err != nil {
			if httpResp != nil && httpResp.StatusCode == http.StatusNotFound {
				return struct{}{}, nil
			}

			return struct{}{}, backoff.Permanent(err)
		}

		return struct{}{}, fmt.Errorf("%s %d still exists", k.TypeName, id)
	}

	if _, err := backoff.Retry(
		ctx,
		waitForDeleted,
		backoff.WithBackOff(backoff.NewConstantBackOff(5*time.Second)),
		backoff.WithMaxElapsedTime(45*time.Minute),
	); err != nil {
		resp.Diagnostics.AddError(
			"error deleting "+k.TypeName,
			fmt.Sprintf("%s %d: DELETE confirmation failed: ", k.TypeName, id)+err.Error(),
		)
	}
}

// Import handles the write-only import IDs then the numeric id.
func (k Kind) Import(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	rawID, ok := WO.Import(ctx, req, resp)
	if !ok {
		return
	}

	id, err := parseInt64(rawID)
	if err != nil {
		resp.Diagnostics.AddError(
			"import "+k.TypeName+" resource",
			fmt.Sprintf("invalid %s ID %q: %s", k.TypeName, rawID, err),
		)

		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}
