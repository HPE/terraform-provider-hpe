// (C) Copyright 2025-2026 Hewlett Packard Enterprise Development LP

package morpheus

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"

	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/cloud"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/cluster"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/datastore"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/environment"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/group"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/image"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/instance"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/instancetypelayout"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/loadbalancer"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/network"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/networkdhcpserver"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/networkdomain"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/networkrouterbgpneighbor"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/ostype"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/ostypeimage"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/policy"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/role"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/serviceplan"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/user"
)

func (SubProvider) GetDataSources(
	_ context.Context,
) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		cluster.NewDataSource,
		cloud.NewDataSource,
		datastore.NewDataSource,
		environment.NewDataSource,
		group.NewDataSource,
		instance.NewDataSource,
		instancetypelayout.NewDataSource,
		image.NewDataSource,
		loadbalancer.NewDataSource,
		network.NewDataSource,
		networkdhcpserver.NewDataSource,
		networkdomain.NewDataSource,
		networkrouterbgpneighbor.NewDataSource,
		ostype.NewDataSource,
		ostypeimage.NewDataSource,
		policy.NewDataSource,
		role.NewDataSource,
		serviceplan.NewDataSource,
		user.NewDataSource,
	}
}
