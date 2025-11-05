// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package image

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/configure"
	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/constants"
)

var (
	// All known statuses are:
	// Saving
	// Active
	// Converting
	// Downloading
	// Failed
	CreateTargetStatuses = []string{
		"Active",
	}

	CreateErrorStatuses = []string{
		"Failed",
	}
)

type Resource struct {
	configure.ResourceWithMorpheusConfigure
}

var (
	_ resource.Resource = &Resource{}
)

func NewResource() resource.Resource {
	return &Resource{}
}

func (g *Resource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_" + constants.SubProviderName + "_image"
	resp.TypeName = strings.Join(
		[]string{req.ProviderTypeName, constants.SubProviderName, "image"},
		"_",
	)
}

func (g *Resource) Schema(
	ctx context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = ImageResourceSchema(ctx)
}
