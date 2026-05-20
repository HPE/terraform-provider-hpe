package cluster_datastore

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type clusterDatastoreModel struct {
	ID         types.Int64  `tfsdk:"id"`
	ClusterID  types.Int64  `tfsdk:"cluster_id"`
	Name       types.String `tfsdk:"name"`
	Active     types.Bool   `tfsdk:"active"`
	Visibility types.String `tfsdk:"visibility"`
}

func ClusterDatastoreSchema(_ context.Context) schema.Schema {
	return schema.Schema{
		Description: "Manages a Morpheus Cluster Datastore resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:    true,
				Description: "The ID of the cluster datastore.",
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
				Computed:    true,
				Description: "The name of the datastore.",
			},
			"active": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
				Description: "Whether the datastore is active.",
			},
			"visibility": schema.StringAttribute{
				Optional:    true,
				Description: "Visibility for the datastore (private or public).",
			},
		},
	}
}
