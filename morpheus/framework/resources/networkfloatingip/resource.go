// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

// Package networkfloatingip implements hpe_morpheus_network_floating_ip: a
// floating IP allocated from a floating IP pool on an SDN network server. The
// floating IP pool and network server it is allocated from are not returned by
// the API, so they are kept in Terraform state; changing either replaces the
// resource.
package networkfloatingip

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/configure"
)

var (
	_ resource.Resource                = &networkFloatingIpResource{}
	_ resource.ResourceWithConfigure   = &networkFloatingIpResource{}
	_ resource.ResourceWithImportState = &networkFloatingIpResource{}
)

type networkFloatingIpResource struct {
	configure.ResourceWithMorpheusConfigure
}

func NewResource() resource.Resource {
	return &networkFloatingIpResource{}
}

func (r *networkFloatingIpResource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_" + "network_floating_ip"
}

func (r *networkFloatingIpResource) Schema(
	ctx context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = NetworkFloatingIpResourceSchema(ctx)
}
