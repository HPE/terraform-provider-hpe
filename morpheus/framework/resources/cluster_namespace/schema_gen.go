package cluster_namespace

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/HPE/terraform-provider-hpe/utils/validators"
)

// ClusterNsPermissionsAttrTypes defines the attr.Type map for the resource_permissions object.
var ClusterNsPermissionsAttrTypes = map[string]attr.Type{
	"all":      types.BoolType,
	"all_plans": types.BoolType,
	"site_ids": types.SetType{ElemType: types.Int64Type},
	"plan_ids": types.SetType{ElemType: types.Int64Type},
}

type clusterNamespaceModel struct {
	ID                  types.Int64  `tfsdk:"id"`
	ClusterID           types.Int64  `tfsdk:"cluster_id"`
	Name                types.String `tfsdk:"name"`
	Description         types.String `tfsdk:"description"`
	Active              types.Bool   `tfsdk:"active"`
	ResourcePermissions types.Object `tfsdk:"resource_permissions"`
}

func ClusterNamespaceSchema(_ context.Context) schema.Schema {
	return schema.Schema{
		Description: "Manages a Morpheus Cluster Namespace resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:    true,
				Description: "The ID of the cluster namespace.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"cluster_id": schema.Int64Attribute{
				Required:    true,
				Description: "The ID of the cluster.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the namespace. Must be 63 characters or less. Must be lower case.",
				Validators: []validator.String{
					stringvalidator.LengthAtMost(63),
					validators.Lowercase(),
				},
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "The description of the namespace.",
			},
			"active": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
				Description: "Whether the namespace is active.",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"resource_permissions": schema.SingleNestedAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Resource permissions controlling group and plan access for this namespace.",
				Attributes: map[string]schema.Attribute{
					"all": schema.BoolAttribute{
						Optional:    true,
						Computed:    true,
						Description: "Pass true to allow access to all groups",
					},
					"all_plans": schema.BoolAttribute{
						Optional:    true,
						Computed:    true,
						Description: "Pass true to allow access to all plans",
					},
					"site_ids": schema.SetAttribute{
						ElementType: types.Int64Type,
						Optional:    true,
						Computed:    true,
						Description: "Array of group (site) IDs that are allowed access",
					},
					"plan_ids": schema.SetAttribute{
						ElementType: types.Int64Type,
						Optional:    true,
						Computed:    true,
						Description: "Array of service plan IDs that are allowed access",
					},
				},
			},
		},
	}
}
