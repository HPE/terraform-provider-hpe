package job

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type jobModel struct {
	ID           types.Int64  `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	WorkflowID   types.Int64  `tfsdk:"workflow_id"`
	ScheduleMode types.String `tfsdk:"schedule_mode"`
	Enabled      types.Bool   `tfsdk:"enabled"`
	TargetType   types.String `tfsdk:"target_type"`
	CustomConfig types.String `tfsdk:"custom_config"`
}

func JobSchema(_ context.Context) schema.Schema {
	return schema.Schema{
		Description: "Manages a Morpheus Job resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:    true,
				Description: "The ID of the job.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the job.",
			},
			"workflow_id": schema.Int64Attribute{
				Optional:    true,
				Description: "The workflow ID to associate with the job.",
			},
			"schedule_mode": schema.StringAttribute{
				Optional:    true,
				Description: "The schedule mode for the job.",
			},
			"enabled": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
				Description: "Whether the job is enabled.",
			},
			"target_type": schema.StringAttribute{
				Optional:    true,
				Description: "The target type for the job.",
			},
			"custom_config": schema.StringAttribute{
				Optional:    true,
				Description: "Custom configuration in JSON format.",
			},
		},
	}
}
