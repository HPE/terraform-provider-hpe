// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package host

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/configure"
)

// Interface compliance assertions.
var (
	_ resource.Resource                = &Resource{}
	_ resource.ResourceWithImportState = &Resource{}
)

// defaultCreateTimeout is used when no create timeout is configured. Adding a
// host (especially bare metal) is asynchronous and can take a while.
const defaultCreateTimeout = 30 * time.Minute

// defaultDeleteTimeout is used when no delete timeout is configured.
const defaultDeleteTimeout = 15 * time.Minute

// Resource implements the hpe_morpheus_host resource.
type Resource struct {
	configure.ResourceWithMorpheusConfigure
	resource.Resource
}

func NewResource() resource.Resource {
	return &Resource{}
}

// Metadata implements resource.Resource.
func (r *Resource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_" + "host"
}

// Schema implements resource.Resource.
func (r *Resource) Schema(
	ctx context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = HostResourceSchema(ctx)
}
