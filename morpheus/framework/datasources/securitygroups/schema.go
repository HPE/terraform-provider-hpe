// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

// Hand-written schema for the plural hpe_morpheus_security_groups data source.
// Unlike the singular data source, this list/filter shape is not produced by the
// code-spec generator, so the schema is maintained here directly.

package securitygroups

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// securityGroupObjectAttrTypes is the single source of truth for the attribute
// types of each element in the security_groups list. The schema's nested object
// attributes and the values built in Read must match this exactly.
func securityGroupObjectAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"account_id":                     types.Int64Type,
		"active":                         types.BoolType,
		"cloud_id":                       types.Int64Type,
		"description":                    types.StringType,
		"enabled":                        types.StringType,
		"external_id":                    types.StringType,
		"group_source":                   types.StringType,
		"id":                             types.Int64Type,
		"name":                           types.StringType,
		"resource_permission_group_ids":  types.SetType{ElemType: types.Int64Type},
		"resource_permission_groups_all": types.BoolType,
		"sync_source":                    types.StringType,
		"tenant_ids":                     types.SetType{ElemType: types.Int64Type},
		"visibility":                     types.StringType,
	}
}

func SecurityGroupsDataSourceSchema(_ context.Context) schema.Schema {
	return schema.Schema{
		Description: "Retrieves a list of Morpheus security groups, optionally filtered by " +
			"name, phrase, cloud, visibility, or active state.",
		MarkdownDescription: "Retrieves a list of Morpheus security groups, optionally filtered by " +
			"name, phrase, cloud, visibility, or active state.",
		Attributes: map[string]schema.Attribute{
			// Filter inputs.
			"name": schema.StringAttribute{
				Optional:            true,
				Description:         "Filter by exact security group name (server-side).",
				MarkdownDescription: "Filter by exact security group name (server-side).",
			},
			"phrase": schema.StringAttribute{
				Optional: true,
				Description: "Filter by a search phrase matched against name or description " +
					"(server-side, partial match).",
				MarkdownDescription: "Filter by a search phrase matched against name or description " +
					"(server-side, partial match).",
			},
			"cloud_id": schema.Int64Attribute{
				Optional:            true,
				Description:         "Filter to security groups belonging to this cloud (zone) ID.",
				MarkdownDescription: "Filter to security groups belonging to this cloud (zone) ID.",
			},
			"visibility": schema.StringAttribute{
				Optional:            true,
				Description:         "Filter by visibility (e.g. public or private).",
				MarkdownDescription: "Filter by visibility (e.g. public or private).",
			},
			"active": schema.BoolAttribute{
				Optional:            true,
				Description:         "Filter by active state.",
				MarkdownDescription: "Filter by active state.",
			},
			// Result.
			"security_groups": schema.ListNestedAttribute{
				Computed:            true,
				Description:         "The list of security groups matching the supplied filters.",
				MarkdownDescription: "The list of security groups matching the supplied filters.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"account_id": schema.Int64Attribute{
							Computed:            true,
							Description:         "The tenant (account) ID that owns the security group.",
							MarkdownDescription: "The tenant (account) ID that owns the security group.",
						},
						"active": schema.BoolAttribute{
							Computed:            true,
							Description:         "Whether the security group is active.",
							MarkdownDescription: "Whether the security group is active.",
						},
						"cloud_id": schema.Int64Attribute{
							Computed:            true,
							Description:         "The cloud (zone) ID for the security group.",
							MarkdownDescription: "The cloud (zone) ID for the security group.",
						},
						"description": schema.StringAttribute{
							Computed:            true,
							Description:         "The description of the security group.",
							MarkdownDescription: "The description of the security group.",
						},
						"enabled": schema.StringAttribute{
							Computed:            true,
							Description:         "Whether the security group is enabled (on/off).",
							MarkdownDescription: "Whether the security group is enabled (on/off).",
						},
						"external_id": schema.StringAttribute{
							Computed:            true,
							Description:         "The external ID of the security group.",
							MarkdownDescription: "The external ID of the security group.",
						},
						"group_source": schema.StringAttribute{
							Computed:            true,
							Description:         "The source of the security group.",
							MarkdownDescription: "The source of the security group.",
						},
						"id": schema.Int64Attribute{
							Computed:            true,
							Description:         "The ID of the security group.",
							MarkdownDescription: "The ID of the security group.",
						},
						"name": schema.StringAttribute{
							Computed:            true,
							Description:         "The name of the security group.",
							MarkdownDescription: "The name of the security group.",
						},
						"resource_permission_group_ids": schema.SetAttribute{
							ElementType:         types.Int64Type,
							Computed:            true,
							Description:         "Set of group IDs that have access to the security group.",
							MarkdownDescription: "Set of group IDs that have access to the security group.",
						},
						"resource_permission_groups_all": schema.BoolAttribute{
							Computed:            true,
							Description:         "Whether all groups have access to the security group.",
							MarkdownDescription: "Whether all groups have access to the security group.",
						},
						"sync_source": schema.StringAttribute{
							Computed:            true,
							Description:         "The sync source of the security group.",
							MarkdownDescription: "The sync source of the security group.",
						},
						"tenant_ids": schema.SetAttribute{
							ElementType:         types.Int64Type,
							Computed:            true,
							Description:         "Set of tenant IDs that are allowed access to the security group.",
							MarkdownDescription: "Set of tenant IDs that are allowed access to the security group.",
						},
						"visibility": schema.StringAttribute{
							Computed:            true,
							Description:         "The visibility of the security group.",
							MarkdownDescription: "The visibility of the security group.",
						},
					},
				},
			},
		},
	}
}

type SecurityGroupsModel struct {
	Name           types.String `tfsdk:"name"`
	Phrase         types.String `tfsdk:"phrase"`
	CloudId        types.Int64  `tfsdk:"cloud_id"`
	Visibility     types.String `tfsdk:"visibility"`
	Active         types.Bool   `tfsdk:"active"`
	SecurityGroups types.List   `tfsdk:"security_groups"`
}
