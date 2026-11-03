// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package taskpowershellscript

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"

	"github.com/HPE/terraform-provider-hpe/morpheus/configure"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/taskscript"
)

// kind describes the PowerShell script task: task-type code "winrmTask" and the
// elevated_shell flag mapped to taskOptions["winrm.elevated"]. It has no
// local_repository_* attributes.
var kind = taskscript.Kind{
	TypeName:           "task_powershell_script",
	TaskTypeCode:       "winrmTask",
	FlagAttr:           "elevated_shell",
	FlagKey:            "winrm.elevated",
	HasLocalRepository: false,
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
	s := TaskPowershellScriptResourceSchema(ctx)
	s.Version = taskscript.SchemaVersion

	if attr, ok := s.Attributes["script_content"].(schema.StringAttribute); ok {
		attr.PlanModifiers = append(
			attr.PlanModifiers,
			[]planmodifier.String{taskscript.ScriptContentSemanticEquality()}...,
		)
		s.Attributes["script_content"] = attr
	}

	resp.Schema = s
}
