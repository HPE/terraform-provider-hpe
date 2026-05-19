package library_instance_type

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type libraryInstanceTypeModel struct {
	ID          types.Int64  `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Code        types.String `tfsdk:"code"`
	Description types.String `tfsdk:"description"`
	Category    types.String `tfsdk:"category"`
	Visibility  types.String `tfsdk:"visibility"`
	Featured    types.Bool   `tfsdk:"featured"`
	Labels      types.List   `tfsdk:"labels"`
}

func LibraryInstanceTypeSchema(_ context.Context) schema.Schema {
	return schema.Schema{
		Description: "Manages a Morpheus Library Instance Type resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:    true,
				Description: "The ID of the instance type.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the instance type.",
			},
			"code": schema.StringAttribute{
				Optional:    true,
				Description: "The code of the instance type.",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "The description of the instance type.",
			},
			"category": schema.StringAttribute{
				Optional:    true,
				Description: "The category of the instance type.",
			},
			"visibility": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("private"),
				Description: "The visibility of the instance type.",
			},
			"featured": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "Whether the instance type is featured.",
			},
			"labels": schema.ListAttribute{
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				Description: "Labels for the instance type.",
			},
		},
	}
}
