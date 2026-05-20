package scale_threshold

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type scaleThresholdModel struct {
	ID            types.Int64   `tfsdk:"id"`
	Name          types.String  `tfsdk:"name"`
	AutoUpscale   types.Bool    `tfsdk:"auto_upscale"`
	AutoDownscale types.Bool    `tfsdk:"auto_downscale"`
	MinCount      types.Int64   `tfsdk:"min_count"`
	MaxCount      types.Int64   `tfsdk:"max_count"`
	CPUEnabled    types.Bool    `tfsdk:"cpu_enabled"`
	MinCPU        types.Float64 `tfsdk:"min_cpu"`
	MaxCPU        types.Float64 `tfsdk:"max_cpu"`
	MemoryEnabled types.Bool    `tfsdk:"memory_enabled"`
	MinMemory     types.Float64 `tfsdk:"min_memory"`
	MaxMemory     types.Float64 `tfsdk:"max_memory"`
}

func ScaleThresholdSchema(_ context.Context) schema.Schema {
	return schema.Schema{
		Description: "Manages a Morpheus Scale Threshold resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:    true,
				Description: "The ID of the scale threshold.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the scale threshold.",
			},
			"auto_upscale": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "Enable auto upscale.",
			},
			"auto_downscale": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "Enable auto downscale.",
			},
			"min_count": schema.Int64Attribute{
				Optional:    true,
				Description: "The minimum number of nodes to scale down to.",
			},
			"max_count": schema.Int64Attribute{
				Optional:    true,
				Description: "The maximum number of nodes to scale up to.",
			},
			"cpu_enabled": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "Enable CPU threshold.",
			},
			"min_cpu": schema.Float64Attribute{
				Optional:    true,
				Description: "Min CPU (%).",
			},
			"max_cpu": schema.Float64Attribute{
				Optional:    true,
				Description: "Max CPU (%).",
			},
			"memory_enabled": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "Enable memory threshold.",
			},
			"min_memory": schema.Float64Attribute{
				Optional:    true,
				Description: "Min memory (%).",
			},
			"max_memory": schema.Float64Attribute{
				Optional:    true,
				Description: "Max memory (%).",
			},
		},
	}
}
