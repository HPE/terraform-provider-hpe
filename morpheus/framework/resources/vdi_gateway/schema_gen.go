package vdi_gateway

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type vdiGatewayModel struct {
	ID          types.Int64  `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	GatewayUrl  types.String `tfsdk:"gateway_url"`
}

func VdiGatewaySchema(_ context.Context) schema.Schema {
	return schema.Schema{
		Description: "Manages a Morpheus VDI Gateway resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:    true,
				Description: "The ID of the VDI gateway.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the VDI gateway.",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "The description of the VDI gateway.",
			},
			"gateway_url": schema.StringAttribute{
				Required:    true,
				Description: "The URL of the VDI gateway.",
			},
		},
	}
}
