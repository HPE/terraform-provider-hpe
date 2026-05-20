package budget

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type budgetModel struct {
	ID          types.Int64  `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Year        types.Int64  `tfsdk:"year"`
	Interval    types.String `tfsdk:"interval"`
	Scope       types.String `tfsdk:"scope"`
	StartDate   types.String `tfsdk:"start_date"`
	EndDate     types.String `tfsdk:"end_date"`
	Enabled     types.Bool   `tfsdk:"enabled"`
}

func BudgetSchema(_ context.Context) schema.Schema {
	return schema.Schema{
		Description: "Manages a Morpheus Budget resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:    true,
				Description: "The ID of the budget.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the budget.",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "The description of the budget.",
			},
			"year": schema.Int64Attribute{
				Optional:    true,
				Description: "The year of the budget.",
			},
			"interval": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The budget interval (year, quarter, month).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"scope": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The scope of the budget.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"start_date": schema.StringAttribute{
				Optional:    true,
				Description: "The start date of the budget.",
			},
			"end_date": schema.StringAttribute{
				Optional:    true,
				Description: "The end date of the budget.",
			},
			"enabled": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
				Description: "Whether the budget is enabled.",
			},
		},
	}
}
