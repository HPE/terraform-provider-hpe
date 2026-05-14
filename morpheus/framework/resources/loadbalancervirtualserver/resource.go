// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package loadbalancervirtualserver

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/HPE/terraform-provider-hpe/morpheus/configure"
)

var (
	_ resource.Resource                = &Resource{}
	_ resource.ResourceWithImportState = &Resource{}
)

func NewResource() resource.Resource {
	return &Resource{}
}

type Resource struct {
	configure.ResourceWithMorpheusConfigure
	resource.Resource
}

func (r *Resource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_morpheus_load_balancer_virtual_server"
}

func (r *Resource) Schema(
	ctx context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = LoadBalancerVirtualServerResourceSchema(ctx)
}

// loadBalancerIDFromInt64 extracts an int64 from a types.Int64 attribute.
func loadBalancerIDFromInt64(n types.Int64) (int64, error) {
	if n.IsNull() || n.IsUnknown() {
		return 0, fmt.Errorf("load_balancer_id is nil")
	}

	return n.ValueInt64(), nil
}
