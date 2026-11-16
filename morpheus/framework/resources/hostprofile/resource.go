// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package hostprofile

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/configure"
)

var (
	_ resource.Resource                = &Resource{}
	_ resource.ResourceWithImportState = &Resource{}
)

// seedCode is the code of the immutable, system-seeded HVM default host
// profile. A profile carrying this code rejects update and delete on the
// API; the provider refuses those two operations up front with a friendly
// diagnostic, and the sweeper and acceptance tests never touch it.
const seedCode = "host-profile.hvm.default"

func NewResource() resource.Resource {
	return &Resource{}
}

type Resource struct {
	configure.ResourceWithMorpheusConfigure
}

func (r *Resource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_" + "host_profile"
}

func (r *Resource) Schema(
	ctx context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = HostProfileResourceSchema(ctx)
}
