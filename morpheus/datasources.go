// (C) Copyright 2025-2026 Hewlett Packard Enterprise Development LP

package morpheus

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"

	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/backup"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/backupjob"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/backuptype"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/cloud"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/cluster"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/clusterlayout"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/datastore"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/datastores"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/datastoretypes"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/environment"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/group"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/image"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/instance"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/instancesnapshot"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/instancetypelayout"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/loadbalancer"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/loadbalancermonitor"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/loadbalancerpool"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/loadbalancerprofile"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/loadbalancervirtualserver"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/monitoringchecktype"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/network"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/networkdhcpserver"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/networkdomain"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/networkedgecluster"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/networkfirewallrule"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/networkfirewallrulegroup"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/networkpool"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/networkpoolservertype"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/networkrouter"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/networkrouterbgpneighbor"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/networkrouterfirewallrulegroups"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/networkrouterroute"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/networkroutertype"
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
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/storageserver"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/storageservers"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/storagevolume"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/storagevolumes"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/subnettype"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/user"

	// missing-data-sources — new data sources (Groups A/B/C)
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/backuphost"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/backupinstance"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/certificate"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/cloudaffinitygroup"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/cloudaffinitygroups"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/clusteraffinitygroup"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/clusteraffinitygroups"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/clusternamespace"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/computeserver"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/computeservers"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/containerscript"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/deployment"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/monitoringalert"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/monitoringcheck"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/monitoringgroup"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/networkpoolserver"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/networkrouterfirewallrule"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/networkrouterfirewallrulegroup"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/networkrouternat"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/provisioninglicense"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/securitygrouprule"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/vdiapp"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/vdigateway"
)

func (p *MorpheusProvider) DataSources(
	_ context.Context,
) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		backup.NewDataSource,
		backuptype.NewDataSource,
		cluster.NewDataSource,
		cloud.NewDataSource,
		datastore.NewDataSource,
		datastores.NewDataSource,
		datastoretypes.NewDataSource,
		environment.NewDataSource,
		group.NewDataSource,
		instance.NewDataSource,
		instancetypelayout.NewDataSource,
		image.NewDataSource,
		loadbalancer.NewDataSource,
		loadbalancermonitor.NewDataSource,
		loadbalancerpool.NewDataSource,
		loadbalancerprofile.NewDataSource,
		loadbalancervirtualserver.NewDataSource,
		network.NewDataSource,
		networkdhcpserver.NewDataSource,
		networkdomain.NewDataSource,
		networkfirewallrule.NewDataSource,
		networkfirewallrulegroup.NewDataSource,
		networkrouter.NewDataSource,
		networkrouterbgpneighbor.NewDataSource,
		networkrouterfirewallrulegroups.NewDataSource,
		networkrouterroute.NewDataSource,
		ostype.NewDataSource,
		ostypeimage.NewDataSource,
		policy.NewDataSource,
		role.NewDataSource,
		securitygroup.NewDataSource,
		securitygroups.NewDataSource,
		serviceplan.NewDataSource,
		storageserver.NewDataSource,
		storageservers.NewDataSource,
		storagevolume.NewDataSource,
		storagevolumes.NewDataSource,
		user.NewDataSource,
		// hpegl VMaaS parity data sources
		networkserver.NewDataSource,
		networkedgecluster.NewDataSource,
		networktransportzone.NewDataSource,
		networkpool.NewDataSource,
		networktype.NewDataSource,
		instancesnapshot.NewDataSource,
		// resource-reference lookup data sources (avoid hard-coded IDs)
		backupjob.NewDataSource,
		clusterlayout.NewDataSource,
		monitoringchecktype.NewDataSource,
		networkpoolservertype.NewDataSource,
		networkroutertype.NewDataSource,
		subnettype.NewDataSource,
		// missing-data-sources — child resource data sources (Group B)
		networkrouternat.NewDataSource,
		networkrouterfirewallrule.NewDataSource,
		networkrouterfirewallrulegroup.NewDataSource,
		securitygrouprule.NewDataSource,
		clusternamespace.NewDataSource,
		clusteraffinitygroup.NewDataSource,
		clusteraffinitygroups.NewDataSource,
		cloudaffinitygroup.NewDataSource,
		cloudaffinitygroups.NewDataSource,
		computeserver.NewDataSource,
		computeservers.NewDataSource,
		// missing-data-sources — standalone data sources (Group A)
		certificate.NewDataSource,
		vdiapp.NewDataSource,
		vdigateway.NewDataSource,
		containerscript.NewDataSource,
		deployment.NewDataSource,
		provisioninglicense.NewDataSource,
		monitoringalert.NewDataSource,
		monitoringgroup.NewDataSource,
		// missing-data-sources — config-bearing data sources (Group C)
		monitoringcheck.NewDataSource,
		networkpoolserver.NewDataSource,
		// missing-data-sources — backup data sources (Group D)
		backuphost.NewDataSource,
		backupinstance.NewDataSource,
	}
}
