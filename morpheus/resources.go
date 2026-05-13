// (C) Copyright 2025-2026 Hewlett Packard Enterprise Development LP

package morpheus

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/cloud"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/cluster"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/datastore"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/group"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/image"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/instance"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/loadbalancer"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/loadbalancermonitor"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/network"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/networkdhcpserver"
	networkrouter "github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/networkrouter"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/networkrouterbgpneighbor"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/ostype"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/ostypeimage"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/policy"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/role"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/serviceplan"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/task"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/user"
)

func (s SubProvider) GetResources(
	_ context.Context,
) []func() resource.Resource {
	resources := []func() resource.Resource{
		cloud.NewResource,
		datastore.NewResource,
		group.NewResource,
		image.NewResource,
		loadbalancer.NewResource,
		loadbalancermonitor.NewResource,
		network.NewResource,
		networkdhcpserver.NewResource,
		networkrouter.NewResource,
		networkrouterbgpneighbor.NewResource,
		ostype.NewResource,
		ostypeimage.NewResource,
		user.NewResource,
		role.NewResource,
		serviceplan.NewResource,
		task.NewResource,
		instance.NewResource,
		policy.NewResource,
		cluster.NewResource,
	}

	return resources
}
