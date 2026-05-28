package option_type

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type optionTypeModel struct {
	ID           types.Int64  `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	Description  types.String `tfsdk:"description"`
	FieldName    types.String `tfsdk:"field_name"`
	Type         types.String `tfsdk:"type"`
	FieldLabel   types.String `tfsdk:"field_label"`
	Placeholder  types.String `tfsdk:"placeholder"`
	DefaultValue types.String `tfsdk:"default_value"`
	Required     types.Bool   `tfsdk:"required"`
	ExportMeta   types.Bool   `tfsdk:"export_meta"`
}

func OptionTypeSchema(_ context.Context) schema.Schema {
	return schema.Schema{
		Description: "Manages a Morpheus Option Type resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:    true,
				Description: "The ID of the option type.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the option type.",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "The description of the option type.",
			},
			"field_name": schema.StringAttribute{
				Required:    true,
				Description: "The field name of the option type.",
			},
			"type": schema.StringAttribute{
				Required:    true,
				Description: "The type of the option type (input, select, checkbox, etc.).",
			},
			"field_label": schema.StringAttribute{
				Optional:    true,
				Description: "The field label of the option type.",
			},
			"placeholder": schema.StringAttribute{
				Optional:    true,
				Description: "The placeholder text for the option type.",
			},
			"default_value": schema.StringAttribute{
				Optional:    true,
				Description: "The default value of the option type.",
			},
			"required": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "Whether the option type is required.",
			},
			"export_meta": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "Whether to export as tag.",
			},
		},
	}
}
