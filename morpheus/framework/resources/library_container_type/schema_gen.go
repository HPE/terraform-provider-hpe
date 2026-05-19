package library_container_type

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type libraryContainerTypeModel struct {
	ID                types.Int64  `tfsdk:"id"`
	Name              types.String `tfsdk:"name"`
	Labels            types.List   `tfsdk:"labels"`
	ShortName         types.String `tfsdk:"short_name"`
	Description       types.String `tfsdk:"description"`
	ContainerVersion  types.String `tfsdk:"container_version"`
	ProvisionTypeCode types.String `tfsdk:"provision_type_code"`
	Scripts           types.List   `tfsdk:"scripts"`
	Templates         types.List   `tfsdk:"templates"`
	VirtualImageID    types.Int64  `tfsdk:"virtual_image_id"`
	StatTypeCode      types.String `tfsdk:"stat_type_code"`
	LogTypeCode       types.String `tfsdk:"log_type_code"`
	ServerType        types.String `tfsdk:"server_type"`
}

func LibraryContainerTypeSchema(_ context.Context) schema.Schema {
	return schema.Schema{
		Description: "Manages a Morpheus Library Container Type (node type) resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:    true,
				Description: "The ID of the container type.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the container type.",
			},
			"labels": schema.ListAttribute{
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				Description: "Labels for the container type.",
			},
			"short_name": schema.StringAttribute{
				Required:    true,
				Description: "The short name (no spaces) for display in container lists.",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "The description of the container type.",
			},
			"container_version": schema.StringAttribute{
				Required:    true,
				Description: "The version of the container type.",
			},
			"provision_type_code": schema.StringAttribute{
				Required:    true,
				Description: "Provision type code (e.g. docker, amazon, vmware).",
			},
			"scripts": schema.ListAttribute{
				Optional:    true,
				Computed:    true,
				ElementType: types.Int64Type,
				Description: "Array of script IDs.",
			},
			"templates": schema.ListAttribute{
				Optional:    true,
				Computed:    true,
				ElementType: types.Int64Type,
				Description: "Array of file template IDs.",
			},
			"virtual_image_id": schema.Int64Attribute{
				Optional:    true,
				Description: "Virtual image ID.",
			},
			"stat_type_code": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Stat type code.",
				PlanModifiers: []planmodifier.String{
					stringUseStateForUnknown{},
				},
			},
			"log_type_code": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Log type code.",
				PlanModifiers: []planmodifier.String{
					stringUseStateForUnknown{},
				},
			},
			"server_type": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Server type (typically 'vm').",
				PlanModifiers: []planmodifier.String{
					stringUseStateForUnknown{},
				},
			},
		},
	}
}

// stringUseStateForUnknown is a plan modifier that uses the state value for unknown values.
type stringUseStateForUnknown struct{}

func (m stringUseStateForUnknown) Description(_ context.Context) string {
	return "Use state value for unknown."
}

func (m stringUseStateForUnknown) MarkdownDescription(_ context.Context) string {
	return "Use state value for unknown."
}

func (m stringUseStateForUnknown) PlanModifyString(_ context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	if !req.PlanValue.IsUnknown() {
		return
	}
	if req.StateValue.IsNull() {
		return
	}
	resp.PlanValue = req.StateValue
}
