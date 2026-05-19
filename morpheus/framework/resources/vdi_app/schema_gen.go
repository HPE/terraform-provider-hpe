package vdi_app

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type vdiAppModel struct {
	ID           types.Int64  `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	Description  types.String `tfsdk:"description"`
	LaunchPrefix types.String `tfsdk:"launch_prefix"`
}

func VdiAppSchema(_ context.Context) schema.Schema {
	return schema.Schema{
		Description: "Manages a Morpheus VDI App resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:    true,
				Description: "The ID of the VDI app.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the VDI app.",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "The description of the VDI app.",
			},
			"launch_prefix": schema.StringAttribute{
				Optional:    true,
				Description: "The RDS app name prefix.",
			},
		},
	}
}
