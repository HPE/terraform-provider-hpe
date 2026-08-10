// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package clusteraffinitygroup

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/configure"
)

// gatedFeature names this resource in the appliance version gate diagnostic.
// Phrased as a plural noun so the message reads "Cluster affinity groups
// require ...".
const gatedFeature = "Cluster affinity groups"

var (
	_ resource.Resource                = &clusterAffinityGroupResource{}
	_ resource.ResourceWithConfigure   = &clusterAffinityGroupResource{}
	_ resource.ResourceWithImportState = &clusterAffinityGroupResource{}
)

type clusterAffinityGroupResource struct {
	configure.ResourceWithMorpheusConfigure
}

func NewResource() resource.Resource {
	return &clusterAffinityGroupResource{}
}

func (r *clusterAffinityGroupResource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_" + "cluster_affinity_group"
}

func (r *clusterAffinityGroupResource) Schema(
	ctx context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = ClusterAffinityGroupResourceSchema(ctx)
}

// Ensure unused imports are satisfied.
var _ *http.Response
