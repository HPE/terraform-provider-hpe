package price_set

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type priceSetModel struct {
	ID         types.Int64  `tfsdk:"id"`
	Name       types.String `tfsdk:"name"`
	Code       types.String `tfsdk:"code"`
	PriceUnit  types.String `tfsdk:"price_unit"`
	Type       types.String `tfsdk:"type"`
	RegionCode types.String `tfsdk:"region_code"`
}

func PriceSetSchema(_ context.Context) schema.Schema {
	return schema.Schema{
		Description: "Manages a Morpheus Price Set resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:    true,
				Description: "The ID of the price set.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the price set.",
			},
			"code": schema.StringAttribute{
				Required:    true,
				Description: "The code of the price set. Must be unique.",
			},
			"price_unit": schema.StringAttribute{
				Required:    true,
				Description: "The price unit.",
			},
			"type": schema.StringAttribute{
				Required:    true,
				Description: "The price set type.",
			},
			"region_code": schema.StringAttribute{
				Optional:    true,
				Description: "The region code.",
			},
		},
	}
}
