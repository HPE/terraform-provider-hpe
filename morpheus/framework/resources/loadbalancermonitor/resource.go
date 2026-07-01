// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package loadbalancermonitor

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/configure"
)

type Resource struct {
	configure.ResourceWithMorpheusConfigure
}

var _ resource.Resource = &Resource{}

func NewResource() resource.Resource {
	return &Resource{}
}

func (r *Resource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_" + "load_balancer_monitor"
}

func (r *Resource) Schema(
	ctx context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = LoadBalancerMonitorResourceSchema(ctx)
}
