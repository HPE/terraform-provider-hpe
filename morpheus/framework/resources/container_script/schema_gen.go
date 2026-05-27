package container_script

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type containerScriptModel struct {
	ID            types.Int64  `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	Labels        types.List   `tfsdk:"labels"`
	Category      types.String `tfsdk:"category"`
	ScriptVersion types.String `tfsdk:"script_version"`
	ScriptPhase   types.String `tfsdk:"script_phase"`
	ScriptType    types.String `tfsdk:"script_type"`
	Script        types.String `tfsdk:"script"`
	RunAsUser     types.String `tfsdk:"run_as_user"`
	SudoUser      types.Bool   `tfsdk:"sudo_user"`
	FailOnError   types.Bool   `tfsdk:"fail_on_error"`
}

func ContainerScriptSchema(_ context.Context) schema.Schema {
	return schema.Schema{
		Description: "Manages a Morpheus Library Container Script resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:    true,
				Description: "The ID of the library container script.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the library container script.",
			},
			"labels": schema.ListAttribute{
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				Description: "The labels of the library container script.",
			},
			"category": schema.StringAttribute{
				Optional:    true,
				Description: "The category of the library container script.",
			},
			"script_version": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("1"),
				Description: "The version of the script.",
			},
			"script_phase": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("provision"),
				Description: "The phase of the script.",
			},
			"script_type": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("bash"),
				Description: "The type of the script.",
			},
			"script": schema.StringAttribute{
				Optional:    true,
				Description: "The script content.",
			},
			"run_as_user": schema.StringAttribute{
				Optional:    true,
				Description: "The user to run the script as.",
			},
			"sudo_user": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "Whether to run the script as sudo.",
			},
			"fail_on_error": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Whether the script fails on error.",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}
