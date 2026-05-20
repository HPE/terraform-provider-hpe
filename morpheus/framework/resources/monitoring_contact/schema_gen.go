package monitoring_contact

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type monitoringContactModel struct {
	ID           types.Int64  `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	EmailAddress types.String `tfsdk:"email_address"`
	SmsAddress   types.String `tfsdk:"sms_address"`
	SlackHook    types.String `tfsdk:"slack_hook"`
}

func MonitoringContactSchema(_ context.Context) schema.Schema {
	return schema.Schema{
		Description: "Manages a Morpheus Monitoring Contact resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:    true,
				Description: "The ID of the monitoring contact.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the monitoring contact.",
			},
			"email_address": schema.StringAttribute{
				Optional:    true,
				Description: "Email notification address.",
			},
			"sms_address": schema.StringAttribute{
				Optional:    true,
				Description: "SMS notification address.",
			},
			"slack_hook": schema.StringAttribute{
				Optional:    true,
				Description: "Slack webhook URL.",
			},
		},
	}
}
