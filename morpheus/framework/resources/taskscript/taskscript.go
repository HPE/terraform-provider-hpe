// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

// Package taskscript holds the shared CRUD, mapping and state-upgrade logic for
// the two script-task framework resources, hpe_morpheus_task_shell_script and
// hpe_morpheus_task_powershell_script. The two resource packages are thin: each
// declares a Kind describing its differences (the Terraform type name, the
// Morpheus task-type code, and the type-specific flag attribute) and converts
// its generated model to and from the common Flat struct used here.
package taskscript

import (
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/wodigest"
)

// SchemaVersion is the current framework schema version for both script-task
// resources. It is 1 (not 0) because these resources replace the retired
// terraform-plugin-sdk/v2 implementations, whose state stored "id" as a string
// (the framework model uses int64) and whose remote_target_password held the
// API's password hash rather than a configured value; the upgrader resets it
// (the attribute itself remains accepted, deprecated). See
// upgrade_state.go in each resource package.
const SchemaVersion = 1

// Kind describes everything that differs between the shell and PowerShell
// script-task resources.
type Kind struct {
	// TypeName is the Terraform resource type name without the provider
	// prefix, e.g. "task_shell_script".
	TypeName string
	// TaskTypeCode is the Morpheus task-type code, e.g. "script" or
	// "winrmTask".
	TaskTypeCode string
	// FlagAttr is the type-specific boolean attribute name, e.g. "sudo" or
	// "elevated_shell".
	FlagAttr string
	// FlagKey is the taskOptions key the flag maps to, e.g. "shell.sudo" or
	// "winrm.elevated". The API stores "on" when the flag is set, and omits
	// the key otherwise.
	FlagKey string
	// HasLocalRepository is true when the resource exposes the
	// local_repository_id / local_repository_ref attributes (shell only).
	HasLocalRepository bool
}

// WO is the write-only tracker shared by both resources:
// remote_target_password_wo is sent to Morpheus but never stored in state, its
// changes detected through a private-state digest. The deprecated *_version
// marker keeps working for the whole 2.2 line.
var WO = wodigest.Tracker{
	Resource: "task_script",
	Attributes: map[string]string{
		"remote_target_password_wo": "remote_target_password_wo_version",
	},
}

// Flat is the common flattened representation of a script task, independent of
// the generated per-resource model. The resource packages convert their model
// to and from this.
type Flat struct {
	Id                     types.Int64
	Name                   types.String
	Code                   types.String
	Labels                 types.Set
	ResultType             types.String
	SourceType             types.String
	ScriptContent          types.String
	ScriptPath             types.String
	RepositoryId           types.Int64
	VersionRef             types.String
	ExecuteTarget          types.String
	LocalRepositoryId      types.String
	LocalRepositoryRef     types.String
	Visibility             types.String
	RemoteTargetHost       types.String
	RemoteTargetPort       types.String
	RemoteTargetUsername   types.String
	RemoteTargetPasswordWo types.String
	// Deprecated plaintext form, accepted until 2.6; sent on create and when
	// it changes, kept in state as configured, never read back.
	RemoteTargetPassword          types.String
	RemoteTargetPasswordWoVersion types.Int64
	Retryable                     types.Bool
	RetryCount                    types.Int64
	RetryDelaySeconds             types.Int64
	AllowCustomConfig             types.Bool
	// Flag is the type-specific boolean (sudo / elevated_shell).
	Flag types.Bool
	// WriteOnlyRevision is the computed digest-revision attribute.
	WriteOnlyRevision types.Int64
}
