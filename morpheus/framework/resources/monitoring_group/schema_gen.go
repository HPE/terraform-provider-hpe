package monitoring_group

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type monitoringGroupModel struct {
	ID          types.Int64  `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	MinHappy    types.Int64  `tfsdk:"min_happy"`
	Severity    types.String `tfsdk:"severity"`
	InUptime    types.Bool   `tfsdk:"in_uptime"`
	Active      types.Bool   `tfsdk:"active"`
}

func MonitoringGroupSchema(_ context.Context) schema.Schema {
	return schema.Schema{
		Description: "Manages a Morpheus Monitoring Group (Check Group) resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:    true,
				Description: "The ID of the monitoring group.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the monitoring group.",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "The description of the monitoring group.",
			},
			"min_happy": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "Minimum number of checks that must be happy to keep the group healthy.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"severity": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Maximum severity level this group can incur on an incident when failing.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"in_uptime": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
				Description: "Whether the group affects account-wide availability calculations.",
			},
			"active": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
				Description: "Whether the group is active.",
			},
		},
	}
}
