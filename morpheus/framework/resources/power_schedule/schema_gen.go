package power_schedule

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type powerScheduleModel struct {
	ID                     types.Int64   `tfsdk:"id"`
	Name                   types.String  `tfsdk:"name"`
	Description            types.String  `tfsdk:"description"`
	ScheduleType           types.String  `tfsdk:"schedule_type"`
	ScheduleTimezone       types.String  `tfsdk:"schedule_timezone"`
	Enabled                types.Bool    `tfsdk:"enabled"`
	MondayOnTime           types.String  `tfsdk:"monday_on_time"`
	MondayOffTime          types.String  `tfsdk:"monday_off_time"`
	TuesdayOnTime          types.String  `tfsdk:"tuesday_on_time"`
	TuesdayOffTime         types.String  `tfsdk:"tuesday_off_time"`
	WednesdayOnTime        types.String  `tfsdk:"wednesday_on_time"`
	WednesdayOffTime       types.String  `tfsdk:"wednesday_off_time"`
	ThursdayOnTime         types.String  `tfsdk:"thursday_on_time"`
	ThursdayOffTime        types.String  `tfsdk:"thursday_off_time"`
	FridayOnTime           types.String  `tfsdk:"friday_on_time"`
	FridayOffTime          types.String  `tfsdk:"friday_off_time"`
	SaturdayOnTime         types.String  `tfsdk:"saturday_on_time"`
	SaturdayOffTime        types.String  `tfsdk:"saturday_off_time"`
	SundayOnTime           types.String  `tfsdk:"sunday_on_time"`
	SundayOffTime          types.String  `tfsdk:"sunday_off_time"`
	TotalMonthlyHoursSaved types.Float64 `tfsdk:"total_monthly_hours_saved"`
}

func PowerScheduleSchema(_ context.Context) schema.Schema {
	return schema.Schema{
		Description: "Manages a Morpheus Power Schedule resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:    true,
				Description: "The ID of the power schedule.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the power schedule.",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "The description of the power schedule.",
			},
			"schedule_type": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("power"),
				Description: "The schedule type of the power schedule.",
			},
			"schedule_timezone": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("UTC"),
				Description: "The timezone of the power schedule.",
			},
			"enabled": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
				Description: "Whether the power schedule is enabled.",
			},
		"monday_on_time": schema.StringAttribute{
			Optional:    true,
			Computed:    true,
			Description: "The Monday on time for the power schedule.",
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"monday_off_time": schema.StringAttribute{
			Optional:    true,
			Computed:    true,
			Description: "The Monday off time for the power schedule.",
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"tuesday_on_time": schema.StringAttribute{
			Optional:    true,
			Computed:    true,
			Description: "The Tuesday on time for the power schedule.",
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"tuesday_off_time": schema.StringAttribute{
			Optional:    true,
			Computed:    true,
			Description: "The Tuesday off time for the power schedule.",
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"wednesday_on_time": schema.StringAttribute{
			Optional:    true,
			Computed:    true,
			Description: "The Wednesday on time for the power schedule.",
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"wednesday_off_time": schema.StringAttribute{
			Optional:    true,
			Computed:    true,
			Description: "The Wednesday off time for the power schedule.",
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"thursday_on_time": schema.StringAttribute{
			Optional:    true,
			Computed:    true,
			Description: "The Thursday on time for the power schedule.",
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"thursday_off_time": schema.StringAttribute{
			Optional:    true,
			Computed:    true,
			Description: "The Thursday off time for the power schedule.",
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"friday_on_time": schema.StringAttribute{
			Optional:    true,
			Computed:    true,
			Description: "The Friday on time for the power schedule.",
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"friday_off_time": schema.StringAttribute{
			Optional:    true,
			Computed:    true,
			Description: "The Friday off time for the power schedule.",
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"saturday_on_time": schema.StringAttribute{
			Optional:    true,
			Computed:    true,
			Description: "The Saturday on time for the power schedule.",
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"saturday_off_time": schema.StringAttribute{
			Optional:    true,
			Computed:    true,
			Description: "The Saturday off time for the power schedule.",
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"sunday_on_time": schema.StringAttribute{
			Optional:    true,
			Computed:    true,
			Description: "The Sunday on time for the power schedule.",
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"sunday_off_time": schema.StringAttribute{
			Optional:    true,
			Computed:    true,
			Description: "The Sunday off time for the power schedule.",
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
			"total_monthly_hours_saved": schema.Float64Attribute{
				Computed:    true,
				Description: "The total monthly hours saved by the power schedule.",
			},
		},
	}
}
