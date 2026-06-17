// (C) Copyright 2025-2026 Hewlett Packard Enterprise Development LP

package morpheus

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"

	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/backup"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/backuptype"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/cloud"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/cluster"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/datastore"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/environment"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/group"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/image"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/instance"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/instancesnapshot"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/instancetypelayout"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/loadbalancer"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/loadbalancermonitor"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/loadbalancervirtualserver"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/network"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/networkdhcpserver"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/networkdomain"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/networkedgecluster"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/networkfirewallrule"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/networkfirewallrulegroup"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/networkpool"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/networkrouter"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/networkrouterbgpneighbor"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/networkrouterroute"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/networkserver"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/networktransportzone"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/networktype"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/ostype"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/ostypeimage"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/policy"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/role"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/securitygroup"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/securitygroups"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/serviceplan"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/user"
)

func (SubProvider) GetDataSources(
	_ context.Context,
) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		backup.NewDataSource,
		backuptype.NewDataSource,
		cluster.NewDataSource,
		cloud.NewDataSource,
		datastore.NewDataSource,
		environment.NewDataSource,
		group.NewDataSource,
		instance.NewDataSource,
		instancetypelayout.NewDataSource,
		image.NewDataSource,
		loadbalancer.NewDataSource,
		loadbalancermonitor.NewDataSource,
		loadbalancervirtualserver.NewDataSource,
		network.NewDataSource,
		networkdhcpserver.NewDataSource,
		networkdomain.NewDataSource,
		networkfirewallrule.NewDataSource,
		networkfirewallrulegroup.NewDataSource,
		networkrouter.NewDataSource,
		networkrouterbgpneighbor.NewDataSource,
		networkrouterroute.NewDataSource,
		ostype.NewDataSource,
		ostypeimage.NewDataSource,
		policy.NewDataSource,
		role.NewDataSource,
		securitygroup.NewDataSource,
		securitygroups.NewDataSource,
		serviceplan.NewDataSource,
		user.NewDataSource,
		// hpegl VMaaS parity data sources
		networkserver.NewDataSource,
		networkedgecluster.NewDataSource,
		networktransportzone.NewDataSource,
		networkpool.NewDataSource,
		networktype.NewDataSource,
		instancesnapshot.NewDataSource,
	}
}
