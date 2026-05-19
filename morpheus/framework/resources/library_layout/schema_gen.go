package library_layout

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type libraryLayoutModel struct {
	ID                types.Int64  `tfsdk:"id"`
	InstanceTypeID    types.Int64  `tfsdk:"instance_type_id"`
	Name              types.String `tfsdk:"name"`
	Labels            types.List   `tfsdk:"labels"`
	InstanceVersion   types.String `tfsdk:"instance_version"`
	Description       types.String `tfsdk:"description"`
	SortOrder         types.Int64  `tfsdk:"sort_order"`
	Creatable         types.Bool   `tfsdk:"creatable"`
	ProvisionTypeCode types.String `tfsdk:"provision_type_code"`
	MemoryRequirement types.String `tfsdk:"memory_requirement"`
	HasAutoScale      types.Bool   `tfsdk:"has_auto_scale"`
	ContainerTypes    types.List   `tfsdk:"container_types"`
	OptionTypes       types.List   `tfsdk:"option_types"`
	SpecTemplates     types.List   `tfsdk:"spec_templates"`
}

func LibraryLayoutSchema(_ context.Context) schema.Schema {
	return schema.Schema{
		Description: "Manages a Morpheus Library Layout (instance type layout) resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:    true,
				Description: "The ID of the layout.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"instance_type_id": schema.Int64Attribute{
				Required:    true,
				Description: "The ID of the instance type this layout belongs to.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the layout.",
			},
			"labels": schema.ListAttribute{
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				Description: "Labels for the layout.",
			},
			"instance_version": schema.StringAttribute{
				Required:    true,
				Description: "The version of the layout.",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "The description of the layout.",
			},
			"sort_order": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "Display order of the layout, higher to lower.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"creatable": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
				Description: "Whether the layout is creatable.",
			},
			"provision_type_code": schema.StringAttribute{
				Required:    true,
				Description: "The provision type code (e.g. docker, vmware, etc).",
			},
			"memory_requirement": schema.StringAttribute{
				Optional:    true,
				Description: "Memory requirement in megabytes.",
			},
			"has_auto_scale": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "Whether horizontal scaling is enabled.",
			},
			"container_types": schema.ListAttribute{
				Optional:    true,
				Computed:    true,
				ElementType: types.Int64Type,
				Description: "Array of layout node type IDs.",
			},
			"option_types": schema.ListAttribute{
				Optional:    true,
				Computed:    true,
				ElementType: types.Int64Type,
				Description: "Array of layout option type IDs.",
			},
			"spec_templates": schema.ListAttribute{
				Optional:    true,
				Computed:    true,
				ElementType: types.Int64Type,
				Description: "Array of layout spec template IDs.",
			},
		},
	}
}
