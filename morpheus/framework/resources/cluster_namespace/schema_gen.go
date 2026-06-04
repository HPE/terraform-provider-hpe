package cluster_namespace

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/HPE/terraform-provider-hpe/utils/validators"
)

type clusterNamespaceModel struct {
	ID          types.Int64  `tfsdk:"id"`
	ClusterID   types.Int64  `tfsdk:"cluster_id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Active      types.Bool   `tfsdk:"active"`
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
		},
	}
}
