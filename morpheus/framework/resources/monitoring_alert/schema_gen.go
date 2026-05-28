package monitoring_alert

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

type monitoringAlertModel struct {
	ID          types.Int64  `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	MinSeverity types.String `tfsdk:"min_severity"`
	MinDuration types.Int64  `tfsdk:"min_duration"`
	Active      types.Bool   `tfsdk:"active"`
	AllChecks   types.Bool   `tfsdk:"all_checks"`
	AllGroups   types.Bool   `tfsdk:"all_groups"`
}

func MonitoringAlertSchema(_ context.Context) schema.Schema {
	return schema.Schema{
		Description: "Manages a Morpheus Monitoring Alert resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:    true,
				Description: "The ID of the monitoring alert.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the monitoring alert.",
			},
			"min_severity": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("critical"),
				Description: "Severity level threshold for sending notifications.",
			},
			"min_duration": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(0),
				Description: "Duration in minutes of the delay before sending notifications.",
			},
			"active": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
				Description: "Whether the alert is active.",
			},
			"all_checks": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "Trigger for all checks.",
			},
			"all_groups": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "Trigger for all check groups.",
			},
		},
	}
}
