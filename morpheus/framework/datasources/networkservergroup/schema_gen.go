// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

// Hand-written schema — this data source has no codegen config.yaml.

package networkservergroup

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func NetworkServerGroupDataSourceSchema(_ context.Context) schema.Schema {
	return schema.Schema{
		Description:         "Look up a Network Server Group (e.g. NSX-T IP Pool / NSGROUP) by name.",
		MarkdownDescription: "Look up a Network Server Group (e.g. NSX-T IP Pool / NSGROUP) by name.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Required:            true,
				Description:         "The name of the network server group",
				MarkdownDescription: "The name of the network server group",
			},
			"id": schema.Int64Attribute{
				Computed:            true,
				Description:         "The ID of the network server group",
				MarkdownDescription: "The ID of the network server group",
			},
			"network_server_id": schema.Int64Attribute{
				Optional:            true,
				Computed:            true,
				Description:         "The ID of the parent network server (integration). If omitted, the first NSX-T server is used.",
				MarkdownDescription: "The ID of the parent network server (integration). If omitted, the first NSX-T server is used.",
			},
			"description": schema.StringAttribute{
				Computed:            true,
				Description:         "Description",
				MarkdownDescription: "Description",
			},
			"external_id": schema.StringAttribute{
				Computed:            true,
				Description:         "External ID",
				MarkdownDescription: "External ID",
			},
			"internal_id": schema.StringAttribute{
				Computed:            true,
				Description:         "Internal ID",
				MarkdownDescription: "Internal ID",
			},
			"visibility": schema.StringAttribute{
				Computed:            true,
				Description:         "Visibility",
				MarkdownDescription: "Visibility",
			},
		},
	}
}

type NetworkServerGroupModel struct {
	Name            types.String `tfsdk:"name"`
	Id              types.Int64  `tfsdk:"id"`
	NetworkServerId types.Int64  `tfsdk:"network_server_id"`
	Description     types.String `tfsdk:"description"`
	ExternalId      types.String `tfsdk:"external_id"`
	InternalId      types.String `tfsdk:"internal_id"`
	Visibility      types.String `tfsdk:"visibility"`
}
