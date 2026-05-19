package price

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/float64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type priceModel struct {
	ID         types.Int64   `tfsdk:"id"`
	Name       types.String  `tfsdk:"name"`
	Code       types.String  `tfsdk:"code"`
	PriceType  types.String  `tfsdk:"price_type"`
	PriceUnit  types.String  `tfsdk:"price_unit"`
	Cost       types.Float64 `tfsdk:"cost"`
	MarkupType types.String  `tfsdk:"markup_type"`
	Markup     types.Float64 `tfsdk:"markup"`
	Currency   types.String  `tfsdk:"currency"`
}

func PriceSchema(_ context.Context) schema.Schema {
	return schema.Schema{
		Description: "Manages a Morpheus Price resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:    true,
				Description: "The ID of the price.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the price.",
			},
			"code": schema.StringAttribute{
				Required:    true,
				Description: "The code of the price. Must be unique.",
			},
			"price_type": schema.StringAttribute{
				Required:    true,
				Description: "The price type.",
			},
			"price_unit": schema.StringAttribute{
				Required:    true,
				Description: "The unit of pricing.",
			},
			"cost": schema.Float64Attribute{
				Required:    true,
				Description: "The cost.",
			},
			"markup_type": schema.StringAttribute{
				Optional:    true,
				Description: "The price adjustment type.",
			},
			"markup": schema.Float64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "The markup amount.",
				PlanModifiers: []planmodifier.Float64{
					float64planmodifier.UseStateForUnknown(),
				},
			},
			"currency": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("USD"),
				Description: "The ISO currency code.",
			},
		},
	}
}
