package execute_schedule

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type executeScheduleModel struct {
	ID               types.Int64  `tfsdk:"id"`
	Name             types.String `tfsdk:"name"`
	Description      types.String `tfsdk:"description"`
	ScheduleType     types.String `tfsdk:"schedule_type"`
	ScheduleTimezone types.String `tfsdk:"schedule_timezone"`
	Cron             types.String `tfsdk:"cron"`
	Enabled          types.Bool   `tfsdk:"enabled"`
}

func ExecuteScheduleSchema(_ context.Context) schema.Schema {
	return schema.Schema{
		Description: "Manages a Morpheus Execute Schedule resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:    true,
				Description: "The ID of the execute schedule.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the execute schedule.",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "The description of the execute schedule.",
			},
			"schedule_type": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("execute"),
				Description: "The schedule type of the execute schedule.",
			},
			"schedule_timezone": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("UTC"),
				Description: "The timezone of the execute schedule.",
			},
			"cron": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("0 0 * * *"),
				Description: "The cron expression for the execute schedule.",
			},
			"enabled": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
				Description: "Whether the execute schedule is enabled.",
			},
		},
	}
}
