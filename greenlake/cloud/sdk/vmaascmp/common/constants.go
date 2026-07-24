// (C) Copyright 2021-2025 Hewlett Packard Enterprise Development LP

package common

const (
	// VmaasCmpAPIBasePath base cmp api path - currently used to talk to CMP directly
	VmaasCmpAPIBasePath = "api"
	// InstancesPath
	InstancesPath = "instances"
	// InstanceTypesPath
	InstanceTypesPath = "instance-types"
	// GroupsPath
	GroupsPath = "groups"
	// NetworksPath
	NetworksPath = "networks"
	// NetworkTypePath
	NetworkTypePath              = "network-types"
	NetworkPoolPath              = "pools"
	NetworkProxyPath             = "proxies"
	DomainPath                   = "domains"
	RoutersNatPath               = "nats"
	RoutersFirewallRuleGroupPath = "firewall-rule-groups"
	// LibraryLayoutPath
	LibraryLayoutPath = "library/layouts"
	// LibraryInstanceTypesPath
	LibraryInstanceTypesPath = "library/instance-types"
	// Service plans path
	ServicePlansPath = "service-plans"
	// Instance Service Plan for storage vol type
	InstancePlanPath = "instances/service-plans"
	// Storage Controller Types for a layout
	StorageControllerTypesPath = "options/storageControllerTypes"
	// CloudsPath
	CloudsPath = "clouds"
	// ZonePath
	ZonePath = "zones"
	// DatstorePath
	DatstorePath = "data-stores"
	// Virtual-Images Path
	ResourcePoolPath               = "resource-pools"
	NetworkRouterPath              = "routers"
	LoadBalancerPath               = "load-balancers"
	LoadBalancerTypePath           = "load-balancer-types"
	LoadBalancerMonitorPath        = "monitors"
	LoadBalancerProfilePath        = "profiles"
	LoadBalancerPoolPath           = "pools"
	LoadBalancerVirtualServersPath = "virtual-servers"
	NetworkRouterTypePath          = "network-router-types"
	NetworkServicePath             = "services"
	NetworkScopePath               = "scopes"
	ServerPath                     = "servers"
	VirtualImagePath               = "virtual-images"
	FolderPath                     = "folders"
	ConfigOptionPath               = "config-options"
	PowerSchedulePath              = "power-schedules"
	EnvironmentPath                = "environments"
	OptionsPath                    = "options"
	ZoneNetworkOptionsPath         = "zoneNetworkOptions"
	ProvisionTypesPath             = "provision-types"
	RouterRoutePath                = "routes"
	RouterBgpNeighborPath          = "bgp-neighbors"
	NetworkEdgeClusterPath         = "edge-clusters"
	DhcpServerPath                 = "dhcp-servers"
	RefreshPath                    = "refresh"
	// Whoami Path
	WhoamiPath            = "whoami"
	LBSSLCertificatesPath = "certificates"

	// headers
	ContentType = "application/json"

	// Morpheus version
	CMPSixZeroFiveVersion = "6.0.5"

	// Broker API paths
	CMPDetails = "vmaas/api/v1/cmp_details"
)
