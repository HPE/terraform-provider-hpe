// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package image

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/configure"
)

var (
	// All known statuses are:
	// Saving
	// Active
	// Converting
	// Downloading
	// Failed
	// Error
	CreateTargetStatuses = []string{
		"Active",
	}

	CreateErrorStatuses = []string{
		"Failed",
		"Error",
	}
)

type Resource struct {
	configure.ResourceWithMorpheusConfigure
}

var _ resource.Resource = &Resource{}

func NewResource() resource.Resource {
	return &Resource{}
}

func (g *Resource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_" + "image"
}

func (g *Resource) Schema(
	ctx context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = ImageResourceSchema(ctx)
}
