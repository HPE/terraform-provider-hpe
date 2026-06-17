// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

// Hand-written schema for the plural hpe_morpheus_security_groups data source.
// Unlike the singular data source, this list/filter shape is not produced by the
// code-spec generator, so the schema is maintained here directly.
//
// The filter input follows the SDKv2 plural convention used by data sources such
// as hpe_morpheus_clouds and hpe_morpheus_networks:
//
//	filter {
//	  name   = "name"
//	  values = ["<regex>", ...]
//	}
//
// implemented here as a terraform-plugin-framework SetNestedBlock.

package securitygroups

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Allow-listed field names usable in a filter block's `name` argument.
const (
	filterFieldName       = "name"
	filterFieldVisibility = "visibility"
	filterFieldCloudID    = "cloud_id"
	filterFieldActive     = "active"
)

// filterFieldNames returns the allow-list of fields a filter block may match on.
func filterFieldNames() []string {
	return []string{
		filterFieldName,
		filterFieldVisibility,
		filterFieldCloudID,
		filterFieldActive,
	}
}

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
		Description: "Retrieves a list of Morpheus security groups, optionally filtered using one " +
			"or more filter blocks.",
		MarkdownDescription: "Retrieves a list of Morpheus security groups, optionally filtered using " +
			"one or more filter blocks.",
		Attributes: map[string]schema.Attribute{
			"sort_ascending": schema.BoolAttribute{
				Optional: true,
				Description: "Whether to sort the returned security groups by id in ascending order. " +
					"Defaults to true.",
				MarkdownDescription: "Whether to sort the returned security groups by id in ascending " +
					"order. Defaults to true.",
			},
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
		Blocks: map[string]schema.Block{
			"filter": schema.SetNestedBlock{
				Description: "Filter block. Repeat to apply multiple filters (all are ANDed together). " +
					"Filter values are case-sensitive and support Go regular expressions " +
					"(https://regex101.com/).",
				MarkdownDescription: "Filter block. Repeat to apply multiple filters (all are ANDed " +
					"together). Filter values are case-sensitive and support Go regular expressions " +
					"(https://regex101.com/).",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Required: true,
							Description: "The name of the field to filter on. Valid names are: " +
								"name, visibility, cloud_id, active.",
							MarkdownDescription: "The name of the field to filter on. Valid names are: " +
								"`name`, `visibility`, `cloud_id`, `active`.",
							Validators: []validator.String{
								stringvalidator.OneOf(filterFieldNames()...),
							},
						},
						"values": schema.SetAttribute{
							ElementType: types.StringType,
							Required:    true,
							Description: "The filter values. A security group matches the block if the " +
								"chosen field matches ANY value (Go regular expression).",
							MarkdownDescription: "The filter values. A security group matches the block if " +
								"the chosen field matches ANY value (Go regular expression).",
							Validators: []validator.Set{
								setvalidator.SizeAtLeast(1),
							},
						},
					},
				},
			},
		},
	}
}

type SecurityGroupsModel struct {
	Filter         []securityGroupsFilterModel `tfsdk:"filter"`
	SortAscending  types.Bool                  `tfsdk:"sort_ascending"`
	SecurityGroups types.List                  `tfsdk:"security_groups"`
}

type securityGroupsFilterModel struct {
	Name   types.String `tfsdk:"name"`
	Values types.Set    `tfsdk:"values"`
}
