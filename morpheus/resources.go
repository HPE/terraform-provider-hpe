// (C) Copyright 2025-2026 Hewlett Packard Enterprise Development LP

package morpheus

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/backup"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/backup_job"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/budget"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/certificate"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/cloud"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/cluster"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/cluster_affinity_group"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/cluster_namespace"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/container_script"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/datastore"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/deployment"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/group"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/image"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/instance"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/library_container_type"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/library_layout"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/library_option_type"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/loadbalancer"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/loadbalancermonitor"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/loadbalancervirtualserver"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/monitoring_alert"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/monitoring_check"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/monitoring_contact"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/monitoring_group"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/network"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/network_group"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/network_pool"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/network_pool_server"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/network_router_firewall_rule"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/network_router_nat"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/network_router_route"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/networkdhcpserver"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/networkfirewallrule"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/networkfirewallrulegroup"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/networkrouter"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/networkrouterbgpneighbor"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/option_list"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/ostype"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/ostypeimage"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/policy"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/power_schedule"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/provisioning_license"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/role"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/security_group"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/security_group_rule"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/serviceplan"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/storage_bucket"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/storage_server"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/storage_volume"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/subnet"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/task"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/user"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/user_source"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/vdi_app"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/vdi_gateway"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/vdi_pool"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/whitelabel_settings"
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
		loadbalancervirtualserver.NewResource,
		network.NewResource,
		networkfirewallrule.NewResource,
		networkfirewallrulegroup.NewResource,
		networkdhcpserver.NewResource,
		networkrouter.NewResource,
		networkrouterbgpneighbor.NewResource,
		network_router_route.NewResource,
		ostype.NewResource,
		ostypeimage.NewResource,
		user.NewResource,
		role.NewResource,
		serviceplan.NewResource,
		task.NewResource,
		instance.NewResource,
		policy.NewResource,
		cluster.NewResource,

		// Sprint 1: Simple resources
		certificate.NewResource,
		power_schedule.NewResource,
		vdi_app.NewResource,
		vdi_gateway.NewResource,
		container_script.NewResource,

		// Sprint 2: Networking
		network_group.NewResource,
		network_pool.NewResource,
		network_pool_server.NewResource,
		network_router_nat.NewResource,
		network_router_firewall_rule.NewResource,
		subnet.NewResource,
		security_group.NewResource,
		security_group_rule.NewResource,

		// Sprint 3: Automation & Orchestration
		deployment.NewResource,

		// Sprint 4: Infrastructure & Compute
		cluster_namespace.NewResource,
		cluster_affinity_group.NewResource,
		storage_server.NewResource,
		storage_volume.NewResource,
		storage_bucket.NewResource,

		// Sprint 5: Monitoring & Operations
		monitoring_check.NewResource,
		monitoring_alert.NewResource,
		monitoring_contact.NewResource,
		monitoring_group.NewResource,
		budget.NewResource,
		backup.NewResource,
		backup_job.NewResource,

		// Sprint 6: Library & Provisioning
		library_layout.NewResource,
		library_container_type.NewResource,
		library_option_type.NewResource,
		option_list.NewResource,
		provisioning_license.NewResource,

		// Sprint 7: Identity, VDI & Governance
		vdi_pool.NewResource,
		user_source.NewResource,
		whitelabel_settings.NewResource,
	}

	return resources
}
