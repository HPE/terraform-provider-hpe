package network_pool_server

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type networkPoolServerModel struct {
	ID              types.Int64  `tfsdk:"id"`
	Name            types.String `tfsdk:"name"`
	TypeID          types.Int64  `tfsdk:"type_id"`
	ServiceUrl      types.String `tfsdk:"service_url"`
	ServiceUsername types.String `tfsdk:"service_username"`
	ServicePassword types.String `tfsdk:"service_password"`
	IgnoreSsl       types.Bool   `tfsdk:"ignore_ssl"`
	Enabled         types.Bool   `tfsdk:"enabled"`
	Status          types.String `tfsdk:"status"`
}

func NetworkPoolServerSchema(_ context.Context) schema.Schema {
	return schema.Schema{
		Description: "Manages a Morpheus Network Pool Server (IPAM integration) resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:    true,
				Description: "The ID of the network pool server.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the network pool server.",
			},
			"type_id": schema.Int64Attribute{
				Required:    true,
				Description: "The ID of the network pool server type.",
			},
			"service_url": schema.StringAttribute{
				Optional:    true,
				Description: "The service URL for the IPAM integration.",
			},
			"service_username": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "The service username for authentication.",
			},
			"service_password": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "The service password for authentication.",
			},
			"ignore_ssl": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "Whether to ignore SSL certificate errors.",
			},
			"enabled": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
				Description: "Whether the network pool server is enabled.",
			},
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "The current status of the network pool server.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}
