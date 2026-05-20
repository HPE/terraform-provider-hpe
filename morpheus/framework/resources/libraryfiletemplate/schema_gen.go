package libraryfiletemplate

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type libraryFileTemplateModel struct {
	ID              types.Int64  `tfsdk:"id"`
	Name            types.String `tfsdk:"name"`
	Labels          types.List   `tfsdk:"labels"`
	FileName        types.String `tfsdk:"file_name"`
	FilePath        types.String `tfsdk:"file_path"`
	Category        types.String `tfsdk:"category"`
	TemplatePhase   types.String `tfsdk:"template_phase"`
	Template        types.String `tfsdk:"template"`
	SettingName     types.String `tfsdk:"setting_name"`
	SettingCategory types.String `tfsdk:"setting_category"`
}

func LibraryFileTemplateSchema(_ context.Context) schema.Schema {
	return schema.Schema{
		Description: "Manages a Morpheus Library File Template resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:    true,
				Description: "The ID of the library file template.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the library file template.",
			},
			"labels": schema.ListAttribute{
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				Description: "The labels of the library file template.",
			},
			"file_name": schema.StringAttribute{
				Required:    true,
				Description: "The file name of the template.",
			},
			"file_path": schema.StringAttribute{
				Optional:    true,
				Description: "The file path of the template.",
			},
			"category": schema.StringAttribute{
				Optional:    true,
				Description: "The category of the library file template.",
			},
			"template_phase": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("provision"),
				Description: "The phase of the template.",
			},
			"template": schema.StringAttribute{
				Optional:    true,
				Description: "The template content.",
			},
			"setting_name": schema.StringAttribute{
				Optional:    true,
				Description: "The setting name.",
			},
			"setting_category": schema.StringAttribute{
				Optional:    true,
				Description: "The setting category.",
			},
		},
	}
}
