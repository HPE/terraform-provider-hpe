// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package loadbalancer

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/configure"
)

const (
	typeCodeHAProxy = "haproxyContainer"
	typeCodeNSXT    = "nsx-t"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                   = &Resource{}
	_ resource.ResourceWithImportState    = &Resource{}
	_ resource.ResourceWithValidateConfig = &Resource{}
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
	resp.TypeName = req.ProviderTypeName + "_morpheus_load_balancer"
}

func (r *Resource) Schema(
	ctx context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = LoadBalancerResourceSchema(ctx)
}

// ValidateConfig enforces NSX-T specific requirements.
func (r *Resource) ValidateConfig(
	ctx context.Context,
	req resource.ValidateConfigRequest,
	resp *resource.ValidateConfigResponse,
) {
	var config LoadBalancerModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	nsxtByBlock := !config.ConfigNsxt.IsNull() && !config.ConfigNsxt.IsUnknown()
	nsxtByTypeCode := !config.TypeCode.IsNull() && !config.TypeCode.IsUnknown() &&
		config.TypeCode.ValueString() == typeCodeNSXT

	if nsxtByBlock || nsxtByTypeCode {
		if config.NetworkServerId.IsNull() && !config.NetworkServerId.IsUnknown() {
			resp.Diagnostics.AddAttributeError(
				path.Root("network_server_id"),
				"Missing required attribute",
				fmt.Sprintf(
					"network_server_id is required for NSX-T load balancers. Set network_server_id when using type_code %q or config_nsxt.",
					typeCodeNSXT,
				),
			)
		}
	}
}
