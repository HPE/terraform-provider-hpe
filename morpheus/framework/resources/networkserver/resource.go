// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkserver

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/configure"
)

// typeCodeJuniperSdn is the network server type of the HPE Morpheus SDN
// (Juniper CN2) controller of an HVM cluster. It is the only type that uses
// the typed config_juniper_sdn block; every other type carries its
// type-specific option values in the free-form config map.
const typeCodeJuniperSdn = "morpheus-junipersdn-network"

var (
	_ resource.Resource                = &networkServerResource{}
	_ resource.ResourceWithConfigure   = &networkServerResource{}
	_ resource.ResourceWithImportState = &networkServerResource{}
	_ resource.ResourceWithModifyPlan  = &networkServerResource{}
)

type networkServerResource struct {
	configure.ResourceWithMorpheusConfigure
}

func NewResource() resource.Resource {
	return &networkServerResource{}
}

func (r *networkServerResource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_" + "network_server"
}

func (r *networkServerResource) Schema(
	ctx context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = NetworkServerResourceSchema(ctx)
}
