package network_pool

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type networkPoolModel struct {
	ID            types.Int64  `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	TypeID        types.Int64  `tfsdk:"type_id"`
	IpCount       types.Int64  `tfsdk:"ip_count"`
	FreeCount     types.Int64  `tfsdk:"free_count"`
	PoolEnabled   types.Bool   `tfsdk:"pool_enabled"`
	DnsDomain     types.String `tfsdk:"dns_domain"`
	DhcpServer    types.Bool   `tfsdk:"dhcp_server"`
	Gateway       types.String `tfsdk:"gateway"`
	Netmask       types.String `tfsdk:"netmask"`
	SubnetAddress types.String `tfsdk:"subnet_address"`
}

func NetworkPoolSchema(_ context.Context) schema.Schema {
	return schema.Schema{
		Description: "Manages a Morpheus Network Pool resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:    true,
				Description: "The ID of the network pool.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the network pool.",
			},
			"type_id": schema.Int64Attribute{
				Required:    true,
				Description: "The ID of the network pool type.",
			},
			"ip_count": schema.Int64Attribute{
				Computed:    true,
				Description: "The total number of IPs in the pool.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"free_count": schema.Int64Attribute{
				Computed:    true,
				Description: "The number of free IPs in the pool.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"pool_enabled": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether the pool is enabled.",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"dns_domain": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The DNS domain for the network pool.",
			},
			"dhcp_server": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Whether DHCP server is enabled.",
			},
			"gateway": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The gateway address.",
			},
			"netmask": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The netmask.",
			},
			"subnet_address": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The subnet address.",
			},
		},
	}
}
