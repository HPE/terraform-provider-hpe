// (C) Copyright 2025-2026 Hewlett Packard Enterprise Development LP

package morpheus

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/backuphost"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/backupinstance"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/backupjob"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/budget"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/certificate"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/cloud"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/cluster"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/clusteraffinitygroup"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/clusternamespace"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/containerscript"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/datastore"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/deployment"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/group"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/image"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/instance"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/instanceclone"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/instancepowerstate"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/instancesnapshot"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/loadbalancer"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/loadbalancermonitor"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/loadbalancerpool"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/loadbalancerprofile"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/loadbalancervirtualserver"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/monitoringalert"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/monitoringcheck"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/monitoringgroup"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/network"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/networkdhcpserver"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/networkfirewallrule"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/networkfirewallrulegroup"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/networkgroup"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/networkpool"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/networkpoolserver"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/networkrouter"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/networkrouterbgpneighbor"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/networkrouterfirewallrule"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/networkrouternat"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/networkrouterroute"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/optionlist"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/ostype"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/ostypeimage"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/policy"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/powerschedule"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/provisioninglicense"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/role"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/securitygroup"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/securitygrouprule"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/serviceplan"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/settingwhitelabel"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/storagebucket"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/storageserver"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/storagevolume"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/subnet"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/task"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/user"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/vdiapp"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/vdigateway"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/vdipool"
)

func (s SubProvider) GetResources(
	_ context.Context,
) []func() resource.Resource {
	resources := []func() resource.Resource{
		// Existing resources
		cloud.NewResource,
		datastore.NewResource,
		group.NewResource,
		image.NewResource,
		loadbalancer.NewResource,
		loadbalancermonitor.NewResource,
		loadbalancerpool.NewResource,
		loadbalancerprofile.NewResource,
		loadbalancervirtualserver.NewResource,
		network.NewResource,
		networkfirewallrule.NewResource,
		networkfirewallrulegroup.NewResource,
		networkdhcpserver.NewResource,
		networkrouter.NewResource,
		networkrouterbgpneighbor.NewResource,
		networkrouterroute.NewResource,
		ostype.NewResource,
		ostypeimage.NewResource,
		user.NewResource,
		role.NewResource,
		serviceplan.NewResource,
		task.NewResource,
		instance.NewResource,
		instancepowerstate.NewResource,
		instanceclone.NewResource,
		instancesnapshot.NewResource,
		policy.NewResource,
		cluster.NewResource,

		// Sprint 1: Simple resources
		certificate.NewResource,
		powerschedule.NewResource,
		vdiapp.NewResource,
		vdigateway.NewResource,
		containerscript.NewResource,

		// Sprint 2: Networking
		networkgroup.NewResource,
		networkpool.NewResource,
		networkpoolserver.NewResource,
		networkrouternat.NewResource,
		networkrouterfirewallrule.NewResource,
		subnet.NewResource,
		securitygroup.NewResource,
		securitygrouprule.NewResource,

		// Sprint 3: Automation & Orchestration
		deployment.NewResource,

		// Sprint 4: Infrastructure & Compute
		clusternamespace.NewResource,
		clusteraffinitygroup.NewResource,
		storageserver.NewResource,
		storagevolume.NewResource,
		storagebucket.NewResource,

		// Sprint 5: Monitoring & Operations
		monitoringcheck.NewResource,
		monitoringalert.NewResource,
		monitoringgroup.NewResource,
		budget.NewResource,
		backupjob.NewResource,
		backuphost.NewResource,
		backupinstance.NewResource,

		// Sprint 6: Library & Provisioning
		optionlist.NewResource,
		provisioninglicense.NewResource,

		// Sprint 7: Identity, VDI & Governance
		vdipool.NewResource,
		settingwhitelabel.NewResource,
	}

	return resources
}
