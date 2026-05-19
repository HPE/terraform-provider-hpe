package network_domain

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type networkDomainModel struct {
	ID               types.Int64  `tfsdk:"id"`
	Name             types.String `tfsdk:"name"`
	Description      types.String `tfsdk:"description"`
	PublicZone       types.Bool   `tfsdk:"public_zone"`
	Active           types.Bool   `tfsdk:"active"`
	DomainController types.Bool   `tfsdk:"domain_controller"`
	DomainUsername   types.String `tfsdk:"domain_username"`
	DomainPassword   types.String `tfsdk:"domain_password"`
	DcServer         types.String `tfsdk:"dc_server"`
	Fqdn             types.String `tfsdk:"fqdn"`
	Visibility       types.String `tfsdk:"visibility"`
}

func NetworkDomainSchema(_ context.Context) schema.Schema {
	return schema.Schema{
		Description: "Manages a Morpheus Network Domain resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:    true,
				Description: "The ID of the network domain.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the network domain.",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "The description of the network domain.",
			},
			"public_zone": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "Whether the domain is a public zone.",
			},
			"active": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
				Description: "Whether the network domain is active.",
			},
			"domain_controller": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "Whether this is a domain controller.",
			},
			"domain_username": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "The domain username.",
			},
			"domain_password": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "The domain password.",
			},
			"dc_server": schema.StringAttribute{
				Optional:    true,
				Description: "The domain controller server.",
			},
			"fqdn": schema.StringAttribute{
				Computed:    true,
				Description: "The fully qualified domain name.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"visibility": schema.StringAttribute{
				Computed:    true,
				Description: "The visibility of the network domain.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}
