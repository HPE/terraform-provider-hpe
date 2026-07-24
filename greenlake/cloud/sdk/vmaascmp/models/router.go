// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package models

type GetAllNetworkRouter struct {
	NetworkRouters []GetNetworkRouter `json:"networkRouters"`
}

type GetSpecificRouterResp struct {
	NetworkRouter GetNetworkRouter `json:"networkRouter"`
}

type CreateRouterResp struct {
	Success bool `json:"success"`
	ID      int  `json:"id" tf:"id,computed"`
}

type GetNetworkRouter struct {
	ID            int                `json:"id" tf:"id,computed"`
	Code          string             `json:"code"`
	Name          string             `json:"name"`
	Category      string             `json:"category"`
	DateCreated   string             `json:"dateCreated"`
	LastUpdated   string             `json:"lastUpdated"`
	RouterType    string             `json:"routerType"`
	Status        string             `json:"status"`
	Enabled       bool               `json:"enabled"`
	ExternalID    string             `json:"externalId"`
	ProviderID    string             `json:"providerId" tf:"provider_id,computed"`
	Type          IDModel            `json:"type"`
	NetworkServer IDModel            `json:"networkServer"`
	Zone          IDModel            `json:"zone"`
	Interfaces    []RouterInterfaces `json:"interfaces" tf:"interfaces,computed"`
}

type RouterInterfaces struct {
	ID        int    `json:"id" tf:"id,computed"`
	IPAddress string `json:"ipAddress" tf:"source_addresses,computed"`
	Cidr      string `json:"cidr" tf:"cidr,computed"`
}

type CreateRouterRequest struct {
	NetworkRouter CreateRouterRequestRouter `json:"networkRouter"`
}

type CreateRouterRequestRouter struct {
	ID              int                       `json:"-" tf:"id,computed"`
	Name            string                    `json:"name" tf:"name"`
	Type            IDModel                   `json:"type,omitempty"`
	TypeID          int                       `json:"-" tf:"type_id,computed"`
	Enabled         bool                      `json:"enabled" tf:"enable"`
	Site            IDStringModel             `json:"site,omitempty"`
	GroupID         string                    `json:"-" tf:"group_id"`
	NetworkServer   IDModel                   `json:"networkServer,omitempty"`
	NetworkServerID int                       `json:"-" tf:"network_server_id,computed"`
	EnableBGP       bool                      `json:"enableBgp"`
	Config          CreateRouterRequestConfig `json:"config"`

	// for tftags parsing
	TfTier0Config *CreateRouterTier0Config `json:"-" tf:"tier0_config,sub"`
	TfTier1Config *CreateRouterTier1Config `json:"-" tf:"tier1_config,sub"`
}

type CreateRouterRequestConfig struct {
	EdgeCluster   string `json:"edgeCluster"`
	HaMode        string `json:"haMode,omitempty"`
	FailOver      string `json:"failOver,omitempty"`
	Tier0Gateways string `json:"tier0Gateway"`
	CreateRouterTier0Config
}

type CreateRouterTier0Config struct {
	Bgp
	RouteRedistributionTier0
	RouteRedistributionTier1
	// For tftag parsing
	TfEdgeCluster string                   `json:"-" tf:"edge_cluster"`
	TfHaMode      string                   `json:"-" tf:"ha_mode"`
	TfFailOver    string                   `json:"-" tf:"fail_over"`
	TfBGP         Bgp                      `json:"-" tf:"bgp,sub"`
	TfRRTier0     RouteRedistributionTier0 `json:"-" tf:"route_redistribution_tier0,sub"`
	TfRRTier1     RouteRedistributionTier1 `json:"-" tf:"route_redistribution_tier1,sub"`
}

type CreateRouterTier1Config struct {
	TfEdgeCluster        string             `json:"-" tf:"edge_cluster"`
	TfFailOver           string             `json:"-" tf:"fail_over"`
	TfTier0Gateways      string             `json:"-" tf:"tier0_gateway,omitempty"`
	TfRouteAdvertisement RouteAdvertisement `json:"-" tf:"route_advertisement,sub"`
}

type RouteRedistributionTier0 struct {
	TIER0STATIC            bool `json:"TIER0_STATIC" tf:"tier0_static"`
	TIER0NAT               bool `json:"TIER0_NAT" tf:"tier0_nat"`
	TIER0IPSECLOCALIP      bool `json:"TIER0_IPSEC_LOCAL_IP" tf:"tier0_ipsec_local_ip"`
	TIER0DNSFORWARDERIP    bool `json:"TIER0_DNS_FORWARDER_IP" tf:"tier0_dns_forwarder_ip"`
	TIER0SERVICEINTERFACE  bool `json:"TIER0_SERVICE_INTERFACE" tf:"tier0_service_interface"`
	TIER0EXTERNALINTERFACE bool `json:"TIER0_EXTERNAL_INTERFACE" tf:"tier0_external_interface"`
	TIER0LOOPBACKINTERFACE bool `json:"TIER0_LOOPBACK_INTERFACE" tf:"tier0_loopback_interface"`
	TIER0SEGMENT           bool `json:"TIER0_SEGMENT" tf:"tier0_segment"`
}

type RouteRedistributionTier1 struct {
	TIER1SERVICEINTERFACE bool `json:"TIER1_SERVICE_INTERFACE," tf:"tier1_service_interface"`
	TIER1SEGMENT          bool `json:"TIER1_SEGMENT" tf:"tier1_segment"`
	RouteAdvertisement
	Tier1Connected          bool `json:"-" tf:"tier1_connected"`
	Tier1StaticRoutes       bool `json:"-" tf:"tier1_static_routes"`
	TIER1DNSFORWARDERIP     bool `json:"-" tf:"tier1_dns_forwarder_ip"`
	TIER1STATIC             bool `json:"-" tf:"tier1_static"`
	TIER1LBVIP              bool `json:"-" tf:"tier1_lb_vip"`
	TIER1NAT                bool `json:"-" tf:"tier1_nat"`
	TIER1LBSNAT             bool `json:"-" tf:"tier1_lb_snat"`
	TIER1IPSECLOCALENDPOINT bool `json:"-" tf:"tier1_ipsec_local_endpoint"`
}

type RouteAdvertisement struct {
	Tier1Connected          bool `json:"TIER1_CONNECTED" tf:"tier1_connected"`
	Tier1StaticRoutes       bool `json:"TIER1_STATIC_ROUTES" tf:"tier1_static_routes"`
	TIER1DNSFORWARDERIP     bool `json:"TIER1_DNS_FORWARDER_IP" tf:"tier1_dns_forwarder_ip"`
	TIER1STATIC             bool `json:"TIER1_STATIC" tf:"tier1_static"`
	TIER1LBVIP              bool `json:"TIER1_LB_VIP" tf:"tier1_lb_vip"`
	TIER1NAT                bool `json:"TIER1_NAT" tf:"tier1_nat"`
	TIER1LBSNAT             bool `json:"TIER1_LB_SNAT" tf:"tier1_lb_snat"`
	TIER1IPSECLOCALENDPOINT bool `json:"TIER1_IPSEC_LOCAL_ENDPOINT" tf:"tier1_ipsec_local_endpoint"`
}

type Bgp struct {
	LOCALASNUM     int    `json:"LOCAL_AS_NUM,omitempty" tf:"local_as_num"`
	ECMP           bool   `json:"ECMP" tf:"ecmp"`
	MULTIPATHRELAX bool   `json:"MULTIPATH_RELAX" tf:"multipath_relax"`
	INTERSRIBGP    bool   `json:"INTER_SR_IBGP" tf:"inter_sr_ibgp"`
	RESTARTMODE    string `json:"RESTART_MODE,omitempty" tf:"restart_mode"`
	RESTARTTIME    int    `json:"RESTART_TIME,omitempty" tf:"restart_time"`
	STALEROUTETIME int    `json:"STALE_ROUTE_TIME,omitempty" tf:"stale_route_time"`
	TfEnableBgp    bool   `json:"-" tf:"enable_bgp"`
}

type GetNetworlRouterTypes struct {
	NetworkRouterTypes []NetworkRouterTypes `json:"networkRouterTypes"`
}

type NetworkRouterTypes struct {
	ID          int    `json:"id"`
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Enabled     bool   `json:"enabled"`
}

type GetNetworkServicesResp struct {
	NetworkServices []GetNetworkServices `json:"networkServices"`
}

type GetNetworkServices struct {
	ServiceType               string `json:"serviceType"`
	ServiceTypeName           string `json:"serviceTypeName"`
	Type                      string `json:"type"`
	TypeName                  string `json:"typeName"`
	Name                      string `json:"name"`
	ID                        int    `json:"id"`
	IntegrationID             int    `json:"integrationId"`
	CanEdit                   bool   `json:"canEdit"`
	CanDelete                 bool   `json:"canDelete"`
	Status                    string `json:"status"`
	LastUpdated               string `json:"lastUpdated"`
	SupportsTenantPermissions bool   `json:"supportsTenantPermissions"`
}

type CreateRouterNatRequest struct {
	CreateRouterNat CreateRouterNat `json:"networkRouterNAT"`
}

type CreateRouterNat struct {
	ID                 int                   `json:"-" tf:"id,computed"`
	Name               string                `json:"name,omitempty" tf:"name"`
	Description        string                `json:"description,omitempty" tf:"description"`
	Enabled            bool                  `json:"enabled" tf:"enabled"`
	Config             CreateRouterNatConfig `json:"config" tf:"config,sub"`
	SourceNetwork      string                `json:"sourceNetwork,omitempty" tf:"source_network"`
	DestinationNetwork string                `json:"destinationNetwork,omitempty" tf:"destination_network"`
	TranslatedNetwork  string                `json:"translatedNetwork,omitempty" tf:"translated_network"`
	TranslatedPorts    string                `json:"translatedPorts,omitempty" tf:"translated_ports"`
	Priority           int                   `json:"priority" tf:"priority"`
	RouterID           int                   `json:"-" tf:"router_id"`
	IsDeprecated       bool                  `json:"-" tf:"is_deprecated"`
}

type TfCreateRouterNat struct {
	NatRule []CreateRouterNat `tf:"nat_rule"`
}

type CreateRouterNatConfig struct {
	Action   string `json:"action,omitempty" tf:"action"`
	Service  string `json:"service,omitempty" tf:"service"`
	Firewall string `json:"firewall,omitempty" tf:"firewall"`
	Scope    string `json:"scope,omitempty" tf:"scope"`
	Logging  bool   `json:"logging" tf:"logging"`
}

type GetSpecificRouterNatResponse struct {
	GetSpecificRouterNat GetSpecificRouterNat `json:"networkRouterNAT"`
}

type GetSpecificRouterNat struct {
	ID                 int    `json:"id" tf:"id"`
	Name               string `json:"name"`
	Description        string `json:"description"`
	Enabled            bool   `json:"enabled"`
	SourceNetwork      string `json:"sourceNetwork"`
	DestinationNetwork string `json:"destinationNetwork"`
	TranslatedNetwork  string `json:"translatedNetwork"`
	SourcePorts        string `json:"sourcePorts"`
	DestinationPorts   string `json:"destinationPorts"`
	TranslatedPorts    string `json:"translatedPorts"`
	Priority           int    `json:"priority"`
	SyncSource         string `json:"syncSource"`
	InternalID         string `json:"internalId"`
	ExternalID         string `json:"externalId"`
	ProviderID         string `json:"providerId"`
	DateCreated        string `json:"dateCreated"`
	LastUpdated        string `json:"lastUpdated"`
	IsDeprecated       bool   `json:"-" tf:"is_deprecated,computed"`
}

type CreateRouterFirewallRuleGroupRequest struct {
	CreateRouterFirewallRuleGroup CreateRouterFirewallRuleGroup `json:"ruleGroup"`
}

type CreateRouterFirewallRuleGroup struct {
	ID           int    `json:"-" tf:"id,computed"`
	Name         string `json:"name,omitempty" tf:"name"`
	Description  string `json:"description,omitempty" tf:"description"`
	Priority     int    `json:"priority" tf:"priority"`
	ExternalType string `json:"externalType" tf:"external_type"`
	GroupLayer   string `json:"groupLayer" tf:"group_layer"`
	RouterID     int    `json:"-" tf:"router_id"`
	IsDeprecated bool   `json:"-" tf:"is_deprecated"`
}

type TfCreateRouterFirewallRuleGroup struct {
	FirewallRuleGroup []CreateRouterFirewallRuleGroup `tf:"firewall_rule_group"`
}

type GetSpecificRouterFirewallRuleGroupResponse struct {
	GetSpecificRouterFirewallRuleGroup GetSpecificRouterFirewallRuleGroup `json:"ruleGroup"`
}

type GetSpecificRouterFirewallRuleGroup struct {
	ID           int    `json:"id" tf:"id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	ExternalID   string `json:"externalId"`
	Status       string `json:"status"`
	Priority     int    `json:"priority"`
	GroupLayer   string `json:"groupLayer"`
	IsDeprecated bool   `json:"-" tf:"is_deprecated,computed"`
}

type CreateRouterRoute struct {
	NetworkRoute RouterRouteBody `json:"networkRoute"`
}

type RouterRouteBody struct {
	ID           int    `json:"-" tf:"id,computed"`
	RouterID     int    `json:"-" tf:"router_id"`
	IsDeprecated bool   `json:"-" tf:"is_deprecated,computed"`
	Name         string `json:"name" tf:"name"`
	Description  string `json:"description" tf:"description"`
	Enabled      bool   `json:"enabled" tf:"enabled"`
	DefaultRoute bool   `json:"defaultRoute" tf:"default_route"`
	Source       string `json:"source" tf:"network"`
	Destination  string `json:"destination" tf:"next_hop"`
	NetworkMtu   int    `json:"networkMtu" tf:"mtu"`
	Priority     int    `json:"priority" tf:"priority"`
}

type GetSpecificRouterRoute struct {
	NetworkRoute GetSpecificRouterRouteBody `json:"networkRoute"`
}

type GetSpecificRouterRouteBody struct {
	ID              int    `json:"id" tf:"id,computed"`
	Name            string `json:"name"`
	Code            string `json:"code" tf:"code,computed"`
	Description     string `json:"description"`
	Priority        int    `json:"priority"`
	RouteType       string `json:"routeType" tf:"route_type,computed"`
	Source          string `json:"source"`
	SourceType      string `json:"sourceType" tf:"source_type,computed"`
	Destination     string `json:"destination"`
	DestinationType string `json:"destinationType"`
	DefaultRoute    bool   `json:"defaultRoute"`
	NetworkMtu      int    `json:"networkMtu"`
	ExternalID      string `json:"externalId" tf:"external_id,computed"`
	ProviderID      string `json:"providerId" tf:"provider_id,computed"`
	Enabled         bool   `json:"enabled"`
	Visible         bool   `json:"visible"`
}

type CreateNetworkRouterBgpNeighborRequest struct {
	NetworkRouterBgpNeighbor CreateRouterRequestBgpNeighborBody `json:"networkRouterBgpNeighbor"`
}

type Config struct {
	SourceAddresses []string `json:"sourceAddresses" tf:"source_addresses"`
}

type CreateRouterRequestBgpNeighborBody struct {
	ID                 int    `json:"-" tf:"id,computed"`
	IPAddress          string `json:"ipAddress" tf:"ip_address"`
	RemoteAs           int    `json:"remoteAs" tf:"remote_as"`
	KeepAlive          int    `json:"keepAlive" tf:"keepalive"`
	HoldDown           int    `json:"holdDown" tf:"holddown"`
	RouteFilteringType string `json:"routeFilteringType" tf:"router_filtering_type"`
	RouteFilteringIn   string `json:"routeFilteringIn" tf:"router_filtering_in"`
	RouteFilteringOut  string `json:"routeFilteringOut" tf:"router_filtering_out"`
	BfdEnabled         bool   `json:"bfdEnabled" tf:"bfd_enabled"`
	BfdInterval        int    `json:"bfdInterval" tf:"bfd_interval"`
	BfdMultiple        int    `json:"bfdMultiple" tf:"bfd_multiple"`
	AllowAsIn          bool   `json:"allowAsIn" tf:"allow_as_in"`
	HopLimit           int    `json:"hopLimit" tf:"hop_limit"`
	RestartMode        string `json:"restartMode" tf:"restart_mode"`
	RouterID           int    `json:"-" tf:"router_id"`
	Config             Config `json:"config" tf:"config,sub"`
}

type GetSpecificNetworkRouterBgpNeighbor struct {
	NetworkRouterBgpNeighbor NetworkRouterBgpNeighborBody `json:"networkRouterBgpNeighbor"`
}

type NetworkRouterBgpNeighborBody struct {
	ID                 int    `json:"id" tf:"id"`
	IPAddress          string `json:"ipAddress"`
	RemoteAs           string `json:"remoteAs"`
	Weight             int    `json:"weight"`
	KeepAlive          int    `json:"keepAlive"`
	HoldDown           int    `json:"holdDown"`
	RouteFilteringType string `json:"routeFilteringType"`
	BfdEnabled         bool   `json:"bfdEnabled"`
	BfdInterval        int    `json:"bfdInterval"`
	BfdMultiple        int    `json:"bfdMultiple"`
	AllowAsIn          bool   `json:"allowAsIn"`
	RestartMode        string `json:"restartMode"`
	ProviderID         string `json:"providerId"`
	HopLimit           int    `json:"hopLimit"`
	SyncSource         string `json:"syncSource"`
	ExternalID         string `json:"externalId"`
	Config             Config `json:"config"`
	DateCreated        string `json:"dateCreated"`
	LastUpdated        string `json:"lastUpdated"`
}

type TransportZonesResp struct {
	NetworkScopes []NetworkScope `json:"networkScopes"`
	Meta          Meta           `json:"meta"`
}

type NetworkScope struct {
	ID            int                `json:"id" tf:"id,computed"`
	ServiceID     int                `json:"-" tf:"service_id,computed"`
	InternalID    string             `json:"internalId" tf:"internal_id,computed"`
	Visibility    string             `json:"visibility"`
	DateCreated   string             `json:"dateCreated"`
	ProviderID    string             `json:"providerId" tf:"provider_id,computed"`
	LastUpdated   string             `json:"lastUpdated"`
	Active        bool               `json:"active"`
	StreamType    string             `json:"streamType"`
	DisplayName   string             `json:"displayName"`
	StatusMessage string             `json:"statusMessage,omitempty"`
	Name          string             `json:"name" tf:"name"`
	Status        string             `json:"status"`
	Enabled       bool               `json:"enabled"`
	ExternalID    string             `json:"externalId" tf:"external_id,computed"`
	Config        NetworkScopeConfig `json:"config"`
	Owner         IDModel            `json:"owner"`
	NetworkServer IDModel            `json:"networkServer"`
	Zone          IDModel            `json:"zone"`
}

type NetworkScopeConfig struct {
	NvdsName               string `json:"nvdsName"`
	HostMembershipCriteria string `json:"hostMembershipCriteria"`
}

type NetworkEdgeClustersResp struct {
	NetworkEdgeClusters []NetworkEdgeClusters `json:"networkEdgeClusters"`
}

type EdgeClusterConfig struct {
	ClusterProfileBindings string `json:"clusterProfileBindings"`
	Members                int    `json:"members"`
	DeploymentType         string `json:"deploymentType"`
	MemberNodeType         string `json:"memberNodeType"`
}

type NetworkEdgeClusters struct {
	ID            int               `json:"id" tf:"id,computed"`
	InternalID    string            `json:"internalId" tf:"internal_id,computed"`
	Visibility    string            `json:"visibility"`
	DateCreated   string            `json:"dateCreated"`
	ProviderID    string            `json:"providerId" tf:"provider_id,computed"`
	LastUpdated   string            `json:"lastUpdated"`
	Active        bool              `json:"active"`
	DisplayName   string            `json:"displayName"`
	Name          string            `json:"name" tf:"name"`
	Enabled       bool              `json:"enabled"`
	Description   string            `json:"description"`
	ExternalID    string            `json:"externalId" tf:"external_id,computed"`
	Config        EdgeClusterConfig `json:"config"`
	Owner         IDModel           `json:"owner"`
	NetworkServer IDModel           `json:"networkServer"`
	Zone          IDModel           `json:"zone"`
}
