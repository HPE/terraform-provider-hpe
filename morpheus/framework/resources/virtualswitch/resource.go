// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package virtualswitch

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/configure"
)

type Resource struct {
	configure.ResourceWithMorpheusConfigure
}

var (
	_ resource.Resource                   = &Resource{}
	_ resource.ResourceWithImportState    = &Resource{}
	_ resource.ResourceWithModifyPlan     = &Resource{}
	_ resource.ResourceWithValidateConfig = &Resource{}
)

func NewResource() resource.Resource {
	return &Resource{}
}

func (r *Resource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_" + "virtual_switch"
}

func (r *Resource) Schema(
	ctx context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = VirtualSwitchResourceSchema(ctx)
}
