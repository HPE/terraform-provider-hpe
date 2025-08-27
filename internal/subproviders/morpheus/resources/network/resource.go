// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package network

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/configure"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &Resource{}
	_ resource.ResourceWithImportState = &Resource{}
)

func NewResource() resource.Resource {
	return &Resource{}
}

// Resource defines the resource implementation.
type Resource struct {
	configure.ResourceWithMorpheusConfigure
	resource.Resource
}

func (r *Resource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_morpheus_network"
}

func (r *Resource) Schema(
	ctx context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = NetworkResourceSchema(ctx)
}
