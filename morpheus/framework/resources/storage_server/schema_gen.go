package storage_server

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type storageServerModel struct {
	ID              types.Int64  `tfsdk:"id"`
	Name            types.String `tfsdk:"name"`
	Type            types.String `tfsdk:"type"`
	Description     types.String `tfsdk:"description"`
	Enabled         types.Bool   `tfsdk:"enabled"`
	ServiceUrl      types.String `tfsdk:"service_url"`
	ServiceUsername types.String `tfsdk:"service_username"`
	ServicePassword types.String `tfsdk:"service_password"`
}

func StorageServerSchema(_ context.Context) schema.Schema {
	return schema.Schema{
		Description: "Manages a Morpheus Storage Server resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:    true,
				Description: "The ID of the storage server.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the storage server.",
			},
			"type": schema.StringAttribute{
				Required:    true,
				Description: "The storage type code or ID.",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "The description of the storage server.",
			},
			"enabled": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
				Description: "Whether the storage server is enabled.",
			},
			"service_url": schema.StringAttribute{
				Optional:    true,
				Description: "The service URL of the storage server.",
			},
			"service_username": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "The service username for the storage server.",
			},
			"service_password": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "The service password for the storage server.",
			},
		},
	}
}
