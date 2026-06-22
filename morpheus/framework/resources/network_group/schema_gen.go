package network_group

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type networkGroupModel struct {
	ID                          types.Int64  `tfsdk:"id"`
	Name                        types.String `tfsdk:"name"`
	Description                 types.String `tfsdk:"description"`
	Visibility                  types.String `tfsdk:"visibility"`
	Active                      types.Bool   `tfsdk:"active"`
	TenantIds                   types.Set    `tfsdk:"tenant_ids"`
	ResourcePermissionGroupsAll types.Bool   `tfsdk:"resource_permission_groups_all"`
	ResourcePermissionGroupIds  types.Set    `tfsdk:"resource_permission_group_ids"`
}

func NetworkGroupSchema(_ context.Context) schema.Schema {
	return schema.Schema{
		Description: "Manages a Morpheus Network Group resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:    true,
				Description: "The ID of the network group.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the network group.",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "The description of the network group.",
			},
			"visibility": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The visibility of the network group.",
				Default:     stringdefault.StaticString("private"),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf("public", "private"),
				},
			},
			"active": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Whether the network group is active.",
				Default:     booldefault.StaticBool(true),
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"tenant_ids": schema.SetAttribute{
				ElementType: types.Int64Type,
				Optional:    true,
				Computed:    true,
				Description: "List of tenant account IDs that are allowed access. Master-account only.",
			},
			"resource_permission_groups_all": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Allow access to all groups (sites).",
			},
			"resource_permission_group_ids": schema.SetAttribute{
				ElementType: types.Int64Type,
				Optional:    true,
				Computed:    true,
				Description: "List of group (site) IDs that are allowed access.",
			},
		},
	}
}
