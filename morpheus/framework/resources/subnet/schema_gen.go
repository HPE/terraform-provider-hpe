package subnet

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type subnetModel struct {
	ID            types.Int64  `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	TypeID        types.Int64  `tfsdk:"type_id"`
	Cidr          types.String `tfsdk:"cidr"`
	Gateway       types.String `tfsdk:"gateway"`
	Netmask       types.String `tfsdk:"netmask"`
	SubnetAddress types.String `tfsdk:"subnet_address"`
	Active        types.Bool   `tfsdk:"active"`
	DhcpServer    types.Bool   `tfsdk:"dhcp_server"`
	Visibility    types.String `tfsdk:"visibility"`
	Labels        types.List   `tfsdk:"labels"`
	Config        types.Map    `tfsdk:"config"`
}

func SubnetSchema(_ context.Context) schema.Schema {
	return schema.Schema{
		Description: "Manages a Morpheus Subnet resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:    true,
				Description: "The ID of the subnet.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Computed:    true,
				Description: "The name of the subnet.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"type_id": schema.Int64Attribute{
				Required:    true,
				Description: "The type ID of the subnet.",
			},
			"cidr": schema.StringAttribute{
				Computed:    true,
				Description: "The CIDR of the subnet.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"gateway": schema.StringAttribute{
				Computed:    true,
				Description: "The gateway of the subnet.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"netmask": schema.StringAttribute{
				Computed:    true,
				Description: "The netmask of the subnet.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"subnet_address": schema.StringAttribute{
				Computed:    true,
				Description: "The subnet address.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"active": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether the subnet is active.",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"dhcp_server": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether DHCP server is enabled.",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"visibility": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("private"),
				Description: "The visibility of the subnet.",
			},
			"labels": schema.ListAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Description: "Labels for the subnet.",
			},
			"config": schema.MapAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Description: "Configuration map for the subnet.",
			},
		},
	}
}
