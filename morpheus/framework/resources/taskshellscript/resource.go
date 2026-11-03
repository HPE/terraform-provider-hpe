// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package taskshellscript

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"

	"github.com/HPE/terraform-provider-hpe/morpheus/configure"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/taskscript"
)

// kind describes the shell script task: task-type code "script", the sudo flag
// mapped to taskOptions["shell.sudo"], and the local_repository_* attributes.
var kind = taskscript.Kind{
	TypeName:           "task_shell_script",
	TaskTypeCode:       "script",
	FlagAttr:           "sudo",
	FlagKey:            "shell.sudo",
	HasLocalRepository: true,
}

type Resource struct {
	configure.ResourceWithMorpheusConfigure
}

var (
	_ resource.Resource                 = &Resource{}
	_ resource.ResourceWithConfigure    = &Resource{}
	_ resource.ResourceWithImportState  = &Resource{}
	_ resource.ResourceWithModifyPlan   = &Resource{}
	_ resource.ResourceWithUpgradeState = &Resource{}
)

func NewResource() resource.Resource {
	return &Resource{}
}

func (r *Resource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_" + kind.TypeName
}

func (r *Resource) Schema(
	ctx context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	s := TaskShellScriptResourceSchema(ctx)
	// Bump the version so Terraform invokes UpgradeState when reading state
	// written by the legacy SDKv2 resource.
	s.Version = taskscript.SchemaVersion

	// Attach the whitespace-insensitive semantic-equality plan modifier to
	// script_content, reproducing the SDKv2 TrimSpace DiffSuppressFunc.
	if attr, ok := s.Attributes["script_content"].(schema.StringAttribute); ok {
		attr.PlanModifiers = append(
			attr.PlanModifiers,
			[]planmodifier.String{taskscript.ScriptContentSemanticEquality()}...,
		)
		s.Attributes["script_content"] = attr
	}

	resp.Schema = s
}
