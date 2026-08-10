// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package cloudaffinitygroup

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/configure"
)

// gatedFeature names this resource in the appliance version gate diagnostic.
// Phrased as a plural noun so the message reads "Cloud affinity groups
// require ...".
const gatedFeature = "Cloud affinity groups"

var (
	_ resource.Resource                = &cloudAffinityGroupResource{}
	_ resource.ResourceWithConfigure   = &cloudAffinityGroupResource{}
	_ resource.ResourceWithImportState = &cloudAffinityGroupResource{}
)

type cloudAffinityGroupResource struct {
	configure.ResourceWithMorpheusConfigure
}

func NewResource() resource.Resource {
	return &cloudAffinityGroupResource{}
}

func (r *cloudAffinityGroupResource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_" + "cloud_affinity_group"
}

func (r *cloudAffinityGroupResource) Schema(
	ctx context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = CloudAffinityGroupResourceSchema(ctx)
}

// Ensure unused imports are satisfied.
var _ *http.Response
