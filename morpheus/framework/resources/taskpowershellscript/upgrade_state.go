// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package taskpowershellscript

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// taskPowershellScriptModelV0 is the state shape produced by the retired
// terraform-plugin-sdk/v2 PowerShell script task resource (schema version 0).
// Its defining differences from the current model are that "id" is a string and
// that its remote_target_password held the API's password HASH (read back by the
// SDKv2 resource), never a configured value.
type taskPowershellScriptModelV0 struct {
	Id                            types.String `tfsdk:"id"`
	Name                          types.String `tfsdk:"name"`
	Code                          types.String `tfsdk:"code"`
	Labels                        types.Set    `tfsdk:"labels"`
	ResultType                    types.String `tfsdk:"result_type"`
	ElevatedShell                 types.Bool   `tfsdk:"elevated_shell"`
	SourceType                    types.String `tfsdk:"source_type"`
	ScriptContent                 types.String `tfsdk:"script_content"`
	ScriptPath                    types.String `tfsdk:"script_path"`
	RepositoryId                  types.Int64  `tfsdk:"repository_id"`
	VersionRef                    types.String `tfsdk:"version_ref"`
	ExecuteTarget                 types.String `tfsdk:"execute_target"`
	Visibility                    types.String `tfsdk:"visibility"`
	RemoteTargetHost              types.String `tfsdk:"remote_target_host"`
	RemoteTargetPort              types.String `tfsdk:"remote_target_port"`
	RemoteTargetUsername          types.String `tfsdk:"remote_target_username"`
	RemoteTargetPassword          types.String `tfsdk:"remote_target_password"`
	RemoteTargetPasswordWo        types.String `tfsdk:"remote_target_password_wo"`
	RemoteTargetPasswordWoVersion types.Int64  `tfsdk:"remote_target_password_wo_version"`
	Retryable                     types.Bool   `tfsdk:"retryable"`
	RetryCount                    types.Int64  `tfsdk:"retry_count"`
	RetryDelaySeconds             types.Int64  `tfsdk:"retry_delay_seconds"`
	AllowCustomConfig             types.Bool   `tfsdk:"allow_custom_config"`
}

// taskPowershellScriptSchemaV0 describes the SDKv2 state shape so the framework
// can decode prior state during an upgrade. Pinned by testdata/state_v0.json.
func taskPowershellScriptSchemaV0() schema.Schema {
	return schema.Schema{
		Version: 0,
		Attributes: map[string]schema.Attribute{
			"id":                                schema.StringAttribute{Computed: true},
			"name":                              schema.StringAttribute{Required: true},
			"code":                              schema.StringAttribute{Optional: true, Computed: true},
			"labels":                            schema.SetAttribute{ElementType: types.StringType, Optional: true, Computed: true},
			"result_type":                       schema.StringAttribute{Optional: true, Computed: true},
			"elevated_shell":                    schema.BoolAttribute{Optional: true, Computed: true},
			"source_type":                       schema.StringAttribute{Required: true},
			"script_content":                    schema.StringAttribute{Optional: true, Computed: true},
			"script_path":                       schema.StringAttribute{Optional: true, Computed: true},
			"repository_id":                     schema.Int64Attribute{Optional: true, Computed: true},
			"version_ref":                       schema.StringAttribute{Optional: true, Computed: true},
			"execute_target":                    schema.StringAttribute{Optional: true, Computed: true},
			"visibility":                        schema.StringAttribute{Optional: true},
			"remote_target_host":                schema.StringAttribute{Optional: true, Computed: true},
			"remote_target_port":                schema.StringAttribute{Optional: true, Computed: true},
			"remote_target_username":            schema.StringAttribute{Optional: true, Computed: true},
			"remote_target_password":            schema.StringAttribute{Optional: true, Computed: true, Sensitive: true},
			"remote_target_password_wo":         schema.StringAttribute{Optional: true, Sensitive: true, WriteOnly: true},
			"remote_target_password_wo_version": schema.Int64Attribute{Optional: true},
			"retryable":                         schema.BoolAttribute{Optional: true, Computed: true},
			"retry_count":                       schema.Int64Attribute{Optional: true},
			"retry_delay_seconds":               schema.Int64Attribute{Optional: true},
			"allow_custom_config":               schema.BoolAttribute{Optional: true, Computed: true},
		},
	}
}

func (r *Resource) UpgradeState(
	_ context.Context,
) map[int64]resource.StateUpgrader {
	priorSchema := taskPowershellScriptSchemaV0()

	return map[int64]resource.StateUpgrader{
		0: {
			PriorSchema:   &priorSchema,
			StateUpgrader: upgradeTaskPowershellScriptStateV0toV1,
		},
	}
}

// upgradeTaskPowershellScriptStateV0toV1 converts SDKv2-produced state (string
// id) to the current framework model (int64 id). It parses the id, resets
// remote_target_password (the legacy value is the API hash, not a configured
// value; the attribute itself stays accepted), drops remote_target_password_wo
// (write-only, never in state), keeps the deprecated marker, and carries the
// rest. write_only_revision is set null so the first
// plan after the upgrade is one "~ write_only_revision" update, clean after.
func upgradeTaskPowershellScriptStateV0toV1(
	ctx context.Context,
	req resource.UpgradeStateRequest,
	resp *resource.UpgradeStateResponse,
) {
	var old taskPowershellScriptModelV0
	resp.Diagnostics.Append(req.State.Get(ctx, &old)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if old.Id.IsNull() || old.Id.IsUnknown() || old.Id.ValueString() == "" {
		resp.Diagnostics.AddError(
			"upgrade task_powershell_script resource state",
			"the legacy task state has no id, so it cannot be upgraded; remove it "+
				"with 'terraform state rm' and import the task again",
		)

		return
	}

	parsed, err := strconv.ParseInt(old.Id.ValueString(), 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(
			"upgrade task_powershell_script resource state",
			"could not convert legacy string id "+old.Id.ValueString()+
				" to an integer: "+err.Error(),
		)

		return
	}

	upgraded := TaskPowershellScriptModel{
		Id:                            types.Int64Value(parsed),
		Name:                          old.Name,
		Code:                          old.Code,
		Labels:                        old.Labels,
		ResultType:                    old.ResultType,
		ElevatedShell:                 old.ElevatedShell,
		SourceType:                    old.SourceType,
		ScriptContent:                 old.ScriptContent,
		ScriptPath:                    old.ScriptPath,
		RepositoryId:                  old.RepositoryId,
		VersionRef:                    old.VersionRef,
		ExecuteTarget:                 old.ExecuteTarget,
		Visibility:                    old.Visibility,
		RemoteTargetHost:              old.RemoteTargetHost,
		RemoteTargetPort:              old.RemoteTargetPort,
		RemoteTargetUsername:          old.RemoteTargetUsername,
		RemoteTargetPasswordWoVersion: old.RemoteTargetPasswordWoVersion,
		Retryable:                     old.Retryable,
		RetryCount:                    old.RetryCount,
		RetryDelaySeconds:             old.RetryDelaySeconds,
		AllowCustomConfig:             old.AllowCustomConfig,
		RemoteTargetPasswordWo:        types.StringNull(),
		// The legacy state held the API's password *hash* here, never the
		// password, so it is not carried: a configuration that still sets the
		// deprecated attribute plans one update that sends the real value.
		RemoteTargetPassword: types.StringNull(),
		WriteOnlyRevision:    types.Int64Null(),
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, upgraded)...)
}
