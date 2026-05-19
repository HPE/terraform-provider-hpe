package virtual_image

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type virtualImageModel struct {
	ID           types.Int64  `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	ImageType    types.String `tfsdk:"image_type"`
	OsTypeId     types.Int64  `tfsdk:"os_type_id"`
	IsCloudInit  types.Bool   `tfsdk:"is_cloud_init"`
	InstallAgent types.Bool   `tfsdk:"install_agent"`
	MinRam       types.Int64  `tfsdk:"min_ram"`
	MinDisk      types.Int64  `tfsdk:"min_disk"`
	Labels       types.List   `tfsdk:"labels"`
}

func VirtualImageSchema(_ context.Context) schema.Schema {
	return schema.Schema{
		Description: "Manages a Morpheus Virtual Image resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:    true,
				Description: "The ID of the virtual image.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the virtual image.",
			},
			"image_type": schema.StringAttribute{
				Required:    true,
				Description: "Code of image type (e.g. vmdk, qcow2, raw, vmware).",
			},
			"os_type_id": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "The ID of the OS type.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"is_cloud_init": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Whether Cloud Init is enabled.",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"install_agent": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Whether to install the agent.",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"min_ram": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "Minimum RAM in bytes.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"min_disk": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "Minimum disk in bytes.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"labels": schema.ListAttribute{
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				Description: "Labels for the virtual image.",
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}
