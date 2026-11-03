// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package taskscript

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

func parseInt64(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}

// taskOpt reads a taskOptions string value by key, preferring the typed union
// members and falling back to the free-form map. The union decoder may land in
// MapmapOfStringAny, so both are handled.
func (k Kind) taskOpt(opts *sdk.GetTasks200ResponseAllOfTaskTaskOptions, key string) *string {
	if opts == nil {
		return nil
	}

	if opts.MapmapOfStringAny != nil {
		if v, ok := (*opts.MapmapOfStringAny)[key]; ok {
			if s, ok := v.(string); ok {
				return &s
			}
		}
	}

	if opts.ShellTaskConfig3 != nil {
		if v := shellOpt(opts.ShellTaskConfig3, key); v != nil {
			return v
		}
	}

	if opts.PowerShellTaskConfig3 != nil {
		if v := powershellOpt(opts.PowerShellTaskConfig3, key); v != nil {
			return v
		}
	}

	return nil
}

func shellOpt(c *sdk.ShellTaskConfig3, key string) *string {
	switch key {
	case "host":
		return c.Host.Get()
	case "port":
		return c.Port.Get()
	case "username":
		return c.Username.Get()
	case "localScriptGitId":
		return c.LocalScriptGitId.Get()
	case "localScriptGitRef":
		return c.LocalScriptGitRef.Get()
	case "shell.sudo":
		return c.ShellSudo
	}

	return nil
}

func powershellOpt(c *sdk.PowerShellTaskConfig3, key string) *string {
	switch key {
	case "host":
		return c.Host.Get()
	case "port":
		return c.Port.Get()
	case "username":
		return c.Username.Get()
	case "localScriptGitId":
		return c.LocalScriptGitId.Get()
	case "localScriptGitRef":
		return c.LocalScriptGitRef.Get()
	case "winrm.elevated":
		return c.WinrmElevated
	}

	return nil
}

// getAsState reads the task and maps it to a Flat. plan supplies the user's
// label casing so the API's casing changes do not churn the plan.
func (k Kind) getAsState(
	ctx context.Context,
	id int64,
	client *sdk.APIClient,
	plan Flat,
) (Flat, bool, diag.Diagnostics) {
	var state Flat
	var diags diag.Diagnostics

	taskResp, httpResp, err := client.AutomationAPI.GetTasks(ctx, id).Execute()
	if errfmt.IsNotFound(httpResp) {
		return state, false, diags
	}
	if err != nil || httpResp.StatusCode != http.StatusOK {
		diags.AddError(
			"populate "+k.TypeName+" resource",
			fmt.Sprintf("%s %d GET failed: ", k.TypeName, id)+errfmt.ErrMsg(err, httpResp),
		)

		return state, false, diags
	}

	task := taskResp.Task
	if task == nil {
		diags.AddError("populate "+k.TypeName+" resource", fmt.Sprintf("%s %d: response task is nil", k.TypeName, id))

		return state, false, diags
	}

	state.Id = convert.Int64ToType(task.Id)
	state.Name = convert.StrToType(task.Name)
	state.Code = convert.StrToType(task.Code.Get())
	state.ResultType = convert.StrToType(task.ResultType.Get())
	state.ExecuteTarget = convert.StrToType(task.ExecuteTarget)
	state.Visibility = convert.StrToType(task.Visibility)
	state.Retryable = convert.BoolToType(task.Retryable)
	state.RetryCount = convert.Int64ToType(task.RetryCount)
	state.RetryDelaySeconds = convert.Int64ToType(task.RetryDelaySeconds)
	state.AllowCustomConfig = convert.BoolToType(task.AllowCustomConfig)

	// labels: preserve the user's casing.
	respLabels := task.Labels
	wantLabels, lerr := convert.SetToStrSlice(plan.Labels)
	if lerr == nil {
		for _, want := range wantLabels {
			for i, got := range respLabels {
				if want != got && strings.EqualFold(want, got) {
					respLabels[i] = want
				}
			}
		}
	}
	state.Labels = convert.StrSliceToSet(respLabels)

	// file-derived attributes.
	if task.File != nil {
		state.SourceType = convert.StrToType(task.File.SourceType.Get())
		state.ScriptContent = convert.StrToType(task.File.Content.Get())
		state.ScriptPath = convert.StrToType(task.File.ContentPath.Get())
		state.VersionRef = convert.StrToType(task.File.ContentRef.Get())

		if task.File.Repository != nil && task.File.Repository.Id != nil {
			state.RepositoryId = convert.Int64ToType(task.File.Repository.Id)
		} else {
			state.RepositoryId = carryRepositoryID(plan)
		}
	} else {
		state.SourceType = types.StringNull()
		state.ScriptContent = types.StringNull()
		state.ScriptPath = types.StringNull()
		state.VersionRef = types.StringNull()
		// Same carry-forward as above: the API never returns the repository
		// id, whether or not it returns a file object at all.
		state.RepositoryId = carryRepositoryID(plan)
	}

	// taskOptions-derived attributes.
	opts := task.TaskOptions
	state.RemoteTargetHost = convert.StrToType(k.taskOpt(opts, "host"))
	state.RemoteTargetPort = convert.StrToType(k.taskOpt(opts, "port"))
	state.RemoteTargetUsername = convert.StrToType(k.taskOpt(opts, "username"))

	if k.HasLocalRepository {
		state.LocalRepositoryId = convert.StrToType(k.taskOpt(opts, "localScriptGitId"))
		state.LocalRepositoryRef = convert.StrToType(k.taskOpt(opts, "localScriptGitRef"))
	} else {
		state.LocalRepositoryId = types.StringNull()
		state.LocalRepositoryRef = types.StringNull()
	}

	flag := k.taskOpt(opts, k.FlagKey)
	state.Flag = types.BoolValue(flag != nil && *flag == "on")

	// remote_target_password_wo is write-only: never read from the API.
	state.RemoteTargetPasswordWo = types.StringNull()

	return state, true, diags
}

// carryRepositoryID returns the configured repository id, or null. The Morpheus
// task API never echoes the repository id back (file.repository is always null
// on read, and file itself may be absent), so the configured value is carried
// forward instead of being nulled and planning a change nobody made.
func carryRepositoryID(plan Flat) types.Int64 {
	if !plan.RepositoryId.IsNull() && !plan.RepositoryId.IsUnknown() {
		return plan.RepositoryId
	}

	return types.Int64Null()
}
