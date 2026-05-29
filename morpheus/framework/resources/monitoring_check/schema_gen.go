package monitoring_check

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type monitoringCheckModel struct {
	ID            types.Int64  `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	CheckTypeID   types.Int64  `tfsdk:"check_type_id"`
	Description   types.String `tfsdk:"description"`
	CheckInterval types.Int64  `tfsdk:"check_interval"`
	InUptime      types.Bool   `tfsdk:"in_uptime"`
	Active        types.Bool   `tfsdk:"active"`
	Severity      types.String `tfsdk:"severity"`
}

func MonitoringCheckSchema(_ context.Context) schema.Schema {
	return schema.Schema{
		Description: "Manages a Morpheus Monitoring Check resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:    true,
				Description: "The ID of the monitoring check.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the monitoring check.",
			},
			"check_type_id": schema.Int64Attribute{
				Required:    true,
				Description: "The check type ID.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "The description of the monitoring check.",
			},
			"check_interval": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(300),
				Description: "Check interval in seconds.",
			},
			"in_uptime": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
				Description: "Whether the check affects account-wide availability calculations.",
			},
			"active": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
				Description: "Whether the check is active.",
			},
			"severity": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("critical"),
				Description: "Severity level threshold for sending notifications.",
			},
		},
	}
}
