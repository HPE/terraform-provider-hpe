// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package taskpowershellscript

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/taskscript"
)

// toFlat converts the generated model to the shared Flat struct. The PowerShell
// resource has no local_repository_* attributes, so those are null.
func toFlat(m TaskPowershellScriptModel) taskscript.Flat {
	return taskscript.Flat{
		Id:                            m.Id,
		Name:                          m.Name,
		Code:                          m.Code,
		Labels:                        m.Labels,
		ResultType:                    m.ResultType,
		SourceType:                    m.SourceType,
		ScriptContent:                 m.ScriptContent,
		ScriptPath:                    m.ScriptPath,
		RepositoryId:                  m.RepositoryId,
		VersionRef:                    m.VersionRef,
		ExecuteTarget:                 m.ExecuteTarget,
		LocalRepositoryId:             types.StringNull(),
		LocalRepositoryRef:            types.StringNull(),
		Visibility:                    m.Visibility,
		RemoteTargetHost:              m.RemoteTargetHost,
		RemoteTargetPort:              m.RemoteTargetPort,
		RemoteTargetUsername:          m.RemoteTargetUsername,
		RemoteTargetPasswordWo:        m.RemoteTargetPasswordWo,
		RemoteTargetPassword:          m.RemoteTargetPassword,
		RemoteTargetPasswordWoVersion: m.RemoteTargetPasswordWoVersion,
		Retryable:                     m.Retryable,
		RetryCount:                    m.RetryCount,
		RetryDelaySeconds:             m.RetryDelaySeconds,
		AllowCustomConfig:             m.AllowCustomConfig,
		Flag:                          m.ElevatedShell,
		WriteOnlyRevision:             m.WriteOnlyRevision,
	}
}

// fromFlat converts the shared Flat struct back to the generated model.
func fromFlat(f taskscript.Flat) TaskPowershellScriptModel {
	return TaskPowershellScriptModel{
		Id:                            f.Id,
		Name:                          f.Name,
		Code:                          f.Code,
		Labels:                        f.Labels,
		ResultType:                    f.ResultType,
		SourceType:                    f.SourceType,
		ScriptContent:                 f.ScriptContent,
		ScriptPath:                    f.ScriptPath,
		RepositoryId:                  f.RepositoryId,
		VersionRef:                    f.VersionRef,
		ExecuteTarget:                 f.ExecuteTarget,
		Visibility:                    f.Visibility,
		RemoteTargetHost:              f.RemoteTargetHost,
		RemoteTargetPort:              f.RemoteTargetPort,
		RemoteTargetUsername:          f.RemoteTargetUsername,
		RemoteTargetPasswordWo:        f.RemoteTargetPasswordWo,
		RemoteTargetPassword:          f.RemoteTargetPassword,
		RemoteTargetPasswordWoVersion: f.RemoteTargetPasswordWoVersion,
		Retryable:                     f.Retryable,
		RetryCount:                    f.RetryCount,
		RetryDelaySeconds:             f.RetryDelaySeconds,
		AllowCustomConfig:             f.AllowCustomConfig,
		ElevatedShell:                 f.Flag,
		WriteOnlyRevision:             f.WriteOnlyRevision,
	}
}

// setState writes a Flat back to a Terraform state.
func setState(dst *tfsdk.State) func(context.Context, taskscript.Flat) diag.Diagnostics {
	return func(ctx context.Context, f taskscript.Flat) diag.Diagnostics {
		return dst.Set(ctx, fromFlat(f))
	}
}
