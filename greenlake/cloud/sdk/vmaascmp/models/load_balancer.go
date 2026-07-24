// (C) Copyright 2022 Hewlett Packard Enterprise Development LP

package models

import "encoding/json"

// Create LB Req
type CreateLoadBalancerRequest struct {
	NetworkLoadBalancer CreateNetworkLoadBalancerRequest `json:"loadBalancer"`
}

type CreateNetworkLoadBalancerRequest struct {
	ID                  int                       `json:"-" tf:"id,computed"`
	Name                string                    `json:"name" tf:"name"`
	Type                string                    `json:"type" tf:"lb_type,computed"`
	Description         string                    `json:"description" tf:"description"`
	NetworkServerID     int                       `json:"networkServerId" tf:"network_server_id,computed"`
	Enabled             bool                      `json:"enabled" tf:"enabled"`
	Config              *CreateConfig             `json:"config" tf:"config,sub"`
	ResourcePermissions EnableResourcePermissions `json:"resourcePermissions" tf:"group_access"`
}

type EnableResourcePermissions struct {
	All   bool              `json:"all" tf:"all"`
	Sites []PermissionSites `json:"sites" tf:"sites"`
}

type PermissionSites struct {
	ID      int  `json:"id" tf:"id"`
	Default bool `json:"default" tf:"default"`
}

type CreateConfig struct {
	AdminState bool   `json:"adminState" tf:"admin_state"`
	Loglevel   string `json:"loglevel" tf:"log_level"`
	Size       string `json:"size" tf:"size"`
	Tier1      string `json:"tier1" tf:"tier1_gateways"`
}

type GetLoadBalancerTypes struct {
	LoadBalancerTypes []LoadBalancerTypes `json:"loadBalancerTypes"`
}

type LoadBalancerTypes struct {
	ID          int    `json:"id"`
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Enabled     bool   `json:"enabled"`
}

// Create LB resp
type CreateNetworkLoadBalancerResp struct {
	Success                 bool                    `json:"success"`
	NetworkLoadBalancerResp NetworkLoadBalancerResp `json:"loadBalancer"`
}

type NetworkLoadBalancerResp struct {
	ID          int          `json:"id" tf:"id,computed"`
	Name        string       `json:"name"`
	AccountID   int          `json:"accountId"`
	Cloud       CloudInfo    `json:"cloud"`
	Type        Types        `json:"type"`
	Description string       `json:"description"`
	Port        int          `json:"port"`
	SSLEnabled  bool         `json:"sslEnabled"`
	Enabled     bool         `json:"enabled"`
	Visibility  string       `json:"visibility"`
	Config      CreateConfig `json:"config"`
	DateCreated string       `json:"dateCreated"`
	LastUpdated string       `json:"lastUpdated"`
}

type CloudInfo struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Types struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"`
}

// Get LBs
type GetNetworkLoadBalancers struct {
	GetNetworkLoadBalancerResp []GetNetworkLoadBalancerResp `json:"loadBalancers"`
}

type GetNetworkLoadBalancerResp struct {
	ID          int          `json:"id" tf:"id,computed"`
	Name        string       `json:"name"`
	AccountID   int          `json:"accountId"`
	Cloud       CloudInfo    `json:"cloud"`
	Type        Types        `json:"type"`
	Description string       `json:"description"`
	Port        int          `json:"port"`
	Host        string       `json:"host"`
	IP          string       `json:"ip"`
	SSLEnabled  bool         `json:"sslEnabled"`
	Enabled     bool         `json:"enabled"`
	Visibility  string       `json:"visibility"`
	Config      CreateConfig `json:"config"`
	DateCreated string       `json:"dateCreated"`
	LastUpdated string       `json:"lastUpdated"`
	Meta        MetaInfo     `json:"meta"`
}

type MetaInfo struct {
	Max    json.Number `json:"max"`
	Offset json.Number `json:"offset"`
	Size   json.Number `json:"size"`
	Total  json.Number `json:"total"`
}

// Get Specific LB
type GetSpecificNetworkLoadBalancer struct {
	GetSpecificNetworkLoadBalancerResp GetSpecificNetworkLoadBalancerResp `json:"loadBalancer"`
}

type GetSpecificNetworkLoadBalancerResp struct {
	ID          int          `json:"id" tf:"id,computed"`
	Name        string       `json:"name"`
	AccountID   int          `json:"accountId"`
	Cloud       CloudInfo    `json:"cloud"`
	Type        Types        `json:"type"`
	Description string       `json:"description"`
	Port        int          `json:"port"`
	Host        string       `json:"host"`
	IP          string       `json:"ip"`
	SSLEnabled  bool         `json:"sslEnabled"`
	Enabled     bool         `json:"enabled"`
	Visibility  string       `json:"visibility"`
	Config      CreateConfig `json:"config"`
	DateCreated string       `json:"dateCreated"`
	LastUpdated string       `json:"lastUpdated"`
	Meta        MetaInfo     `json:"meta"`
}

// Create LB Monitor Req
type CreateLBMonitor struct {
	CreateLBMonitorReq CreateLBMonitorReq `json:"loadBalancerMonitor"`
}

type CreateLBMonitorReq struct {
	ID                  int                         `json:"-" tf:"id,computed"`
	LbID                int                         `json:"-" tf:"lb_id"`
	Name                string                      `json:"name" tf:"name"`
	Description         string                      `json:"description" tf:"description"`
	Type                string                      `json:"monitorType" tf:"type"`
	Timeout             int                         `json:"monitorTimeout"`
	Interval            int                         `json:"monitorInterval"`
	RequestVersion      string                      `json:"sendVersion"`
	RequestMethod       string                      `json:"sendType"`
	ResponseStatusCodes string                      `json:"receiveCode"`
	ResponseData        string                      `json:"receiveData"`
	RequestURL          string                      `json:"monitorDestination"`
	RequestBody         string                      `json:"sendData"`
	AliasPort           int                         `json:"aliasPort"`
	RiseCount           int                         `json:"riseCount"`
	FallCount           int                         `json:"fallCount"`
	DataLength          int                         `json:"dataLength"`
	MaxFail             int                         `json:"maxRetry"`
	TfHTTPConfig        *CreateHTTPMonitorConfig    `json:"-" tf:"http_monitor,sub"`
	TfHTTPSConfig       *CreateHTTPSMonitorConfig   `json:"-" tf:"https_monitor,sub"`
	TfIcmpConfig        *CreateIcmpMonitorConfig    `json:"-" tf:"icmp_monitor,sub"`
	TfPassiveConfig     *CreatePassiveMonitorConfig `json:"-" tf:"passive_monitor,sub"`
	TfTCPConfig         *CreateTCPMonitorConfig     `json:"-" tf:"tcp_monitor,sub"`
	TfUDPConfig         *CreateUDPMonitorConfig     `json:"-" tf:"udp_monitor,sub"`
}

type CreateHTTPMonitorConfig struct {
	Timeout             int    `json:"-" tf:"timeout"`
	Interval            int    `json:"-" tf:"interval"`
	RequestVersion      string `json:"-" tf:"request_version"`
	RequestMethod       string `json:"-" tf:"request_method"`
	ResponseStatusCodes string `json:"-" tf:"response_status_codes"`
	ResponseData        string `json:"-" tf:"response_data"`
	RequestURL          string `json:"-" tf:"request_url"`
	RequestBody         string `json:"-" tf:"request_body"`
	AliasPort           int    `json:"-" tf:"monitor_port"`
	RiseCount           int    `json:"-" tf:"rise_count"`
	FallCount           int    `json:"-" tf:"fall_count"`
}

type CreateHTTPSMonitorConfig struct {
	Timeout             int    `json:"-" tf:"timeout"`
	Interval            int    `json:"-" tf:"interval"`
	RequestVersion      string `json:"-" tf:"request_version"`
	RequestMethod       string `json:"-" tf:"request_method"`
	ResponseStatusCodes string `json:"-" tf:"response_status_codes"`
	ResponseData        string `json:"-" tf:"response_data"`
	RequestURL          string `json:"-" tf:"request_url"`
	RequestBody         string `json:"-" tf:"request_body"`
	AliasPort           int    `json:"-" tf:"monitor_port"`
	RiseCount           int    `json:"-" tf:"rise_count"`
	FallCount           int    `json:"-" tf:"fall_count"`
}

type CreateIcmpMonitorConfig struct {
	FallCount  int `json:"-" tf:"fall_count"`
	Interval   int `json:"-" tf:"interval"`
	AliasPort  int `json:"-" tf:"monitor_port"`
	RiseCount  int `json:"-" tf:"rise_count"`
	DataLength int `json:"-" tf:"data_length"`
	Timeout    int `json:"-" tf:"timeout"`
}

type CreatePassiveMonitorConfig struct {
	Timeout int `json:"-" tf:"timeout"`
	MaxFail int `json:"-" tf:"max_fail"`
}

type CreateTCPMonitorConfig struct {
	FallCount    int    `json:"-" tf:"fall_count"`
	Interval     int    `json:"-" tf:"interval"`
	AliasPort    int    `json:"-" tf:"monitor_port"`
	RiseCount    int    `json:"-" tf:"rise_count"`
	Timeout      int    `json:"-" tf:"timeout"`
	RequestBody  string `json:"-" tf:"request_body"`
	ResponseData string `json:"-" tf:"response_data"`
}

type CreateUDPMonitorConfig struct {
	FallCount    int    `json:"-" tf:"fall_count"`
	Interval     int    `json:"-" tf:"interval"`
	AliasPort    int    `json:"-" tf:"monitor_port"`
	RiseCount    int    `json:"-" tf:"rise_count"`
	Timeout      int    `json:"-" tf:"timeout"`
	RequestBody  string `json:"-" tf:"request_body"`
	ResponseData string `json:"-" tf:"response_data"`
}

// Create LB Monitor Resp
type CreateLBMonitorResp struct {
	Success       bool          `json:"success"`
	LBMonitorResp LBMonitorResp `json:"loadBalancerMonitor"`
}

type LBMonitorResp struct {
	ID                 int       `json:"id"`
	Name               string    `json:"name"`
	Visibility         string    `json:"visibility"`
	Description        string    `json:"description"`
	MonitorType        string    `json:"monitorType"`
	MonitorInterval    int       `json:"monitorInterval"`
	MonitorTimeout     int       `json:"monitorTimeout"`
	SendVersion        string    `json:"sendVersion"`
	SendType           string    `json:"sendType"`
	ReceiveCode        string    `json:"receiveCode"`
	MonitorDestination string    `json:"monitorDestination"`
	MonitorReverse     bool      `json:"monitorReverse"`
	MonitorTransparent bool      `json:"monitorTransparent"`
	MonitorAdaptive    bool      `json:"monitorAdaptive"`
	AliasPort          int       `json:"aliasPort"`
	InternalID         string    `json:"internalId"`
	MonitorSource      string    `json:"monitorSource"`
	Status             string    `json:"status"`
	Enabled            bool      `json:"enabled"`
	FallCount          int       `json:"fallCount"`
	RiseCount          int       `json:"riseCount"`
	DateCreated        string    `json:"dateCreated"`
	LastUpdated        string    `json:"lastUpdated"`
	LoadBalancer       LBMonitor `json:"loadBalancer"`
}

type LBMonitor struct {
	ID   int           `json:"id"`
	Name string        `json:"name"`
	IP   string        `json:"ip"`
	Type LBMonitorType `json:"type"`
}

type LBMonitorType struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"`
}

// Create LB Profile Req
type CreateLBProfile struct {
	CreateLBProfileReq CreateLBProfileReq `json:"loadBalancerProfile"`
}

type CreateLBProfileReq struct {
	ID              int                         `json:"-" tf:"id,computed"`
	LbID            int                         `json:"-" tf:"lb_id"`
	Name            string                      `json:"name" tf:"name"`
	Description     string                      `json:"description" tf:"description"`
	ServiceType     string                      `json:"serviceType"`
	ProfileType     string                      `json:"" tf:"profile_type"`
	TfHTTPConfig    *CreateHTTPProfileConfig    `json:"-" tf:"http_profile,sub"`
	TfTCPConfig     *CreateTCPProfileConfig     `json:"-" tf:"tcp_profile,sub"`
	TfUDPConfig     *CreateUDPProfileConfig     `json:"-" tf:"udp_profile,sub"`
	TfCookieConfig  *CreateCookieProfileConfig  `json:"-" tf:"cookie_profile,sub"`
	TfGenericConfig *CreateGenericProfileConfig `json:"-" tf:"generic_profile,sub"`
	TfSourceConfig  *CreateSourceProfileConfig  `json:"-" tf:"sourceip_profile,sub"`
	TfClientConfig  *CreateClientProfileConfig  `json:"-" tf:"client_profile,sub"`
	TfServerConfig  *CreateServerProfileConfig  `json:"-" tf:"server_profile,sub"`
	ProfileConfig   *LBProfile                  `json:"config"  tf:"config,sub"`
}

type CreateHTTPProfileConfig struct {
	ServiceType        string `json:"-" tf:"service_type"`
	HTTPIdleTimeout    int    `json:"-" tf:"http_idle_timeout"`
	HTTPSRedirect      string `json:"-" tf:"redirection"`
	RequestHeaderSize  int    `json:"-" tf:"request_header_size"`
	ResponseHeaderSize int    `json:"-" tf:"response_header_size"`
	NtlmAuthentication bool   `json:"-" tf:"ntlm_authentication"`
	RequestBodySize    string `json:"-" tf:"request_body_size"`
	ResponseTimeout    int    `json:"-" tf:"response_timeout"`
	XForwardedFor      string `json:"-" tf:"x_forwarded_for"`
}

type LBProfile struct {
	ProfileType              string `json:"profileType"`
	FastTCPIdleTimeout       int    `json:"fastTcpIdleTimeout"`
	FastUDPIdleTimeout       int    `json:"fastUdpIdleTimeout"`
	HTTPIdleTimeout          int    `json:"httpIdleTimeout"`
	ConnectionCloseTimeout   int    `json:"connectionCloseTimeout"`
	HaFlowMirroring          bool   `json:"haFlowMirroring"`
	RequestHeaderSize        int    `json:"requestHeaderSize"`
	ResponseHeaderSize       int    `json:"responseHeaderSize"`
	HTTPSRedirect            string `json:"httpsRedirect"`
	XForwardedFor            string `json:"xForwardedFor"`
	RequestBodySize          string `json:"requestBodySize"`
	ResponseTimeout          int    `json:"responseTimeout"`
	NtlmAuthentication       bool   `json:"ntlmAuthentication"`
	SharePersistence         bool   `json:"sharePersistence"`
	CookieName               string `json:"cookieName"`
	CookieFallback           bool   `json:"cookieFallback"`
	CookieGarbling           bool   `json:"cookieGarbling"`
	CookieMode               string `json:"cookieMode"`
	CookieType               string `json:"cookieType"`
	CookieDomain             string `json:"cookieDomain"`
	CookiePath               string `json:"cookiePath"`
	MaxIdleTime              int    `json:"maxIdleTime"`
	HaPersistenceMirroring   bool   `json:"haPersistenceMirroring"`
	PersistenceEntryTimeout  int    `json:"persistenceEntryTimeout"`
	PurgeEntries             bool   `json:"purgeEntries"`
	SSLSuite                 string `json:"sslSuite"`
	SessionCache             bool   `json:"sessionCache"`
	SessionCacheEntryTimeout int    `json:"sessionCacheTimeout"`
	PreferServerCipher       bool   `json:"preferServerCipher"`
	Tag                      []Tags `json:"tags" tf:"tags"`
}

type CreateClientProfileConfig struct {
	ServiceType              string `json:"-" tf:"service_type"`
	SSLSuite                 string `json:"-" tf:"ssl_suite"`
	SessionCache             bool   `json:"-" tf:"session_cache"`
	SessionCacheEntryTimeout int    `json:"-" tf:"session_cache_entry_timeout"`
	PreferServerCipher       bool   `json:"-" tf:"prefer_server_cipher"`
}

type CreateServerProfileConfig struct {
	ServiceType  string `json:"-" tf:"service_type"`
	SSLSuite     string `json:"-" tf:"ssl_suite"`
	SessionCache bool   `json:"-" tf:"session_cache"`
}

type CreateSourceProfileConfig struct {
	ServiceType             string `json:"-" tf:"service_type"`
	HaPersistenceMirroring  bool   `json:"-" tf:"ha_persistence_mirroring"`
	PersistenceEntryTimeout int    `json:"-" tf:"persistence_entry_timeout"`
	PurgeEntries            bool   `json:"-" tf:"purge_entries_when_full"`
	SharePersistence        bool   `json:"-" tf:"share_persistence"`
}

type CreateGenericProfileConfig struct {
	ServiceType             string `json:"-" tf:"service_type"`
	HaPersistenceMirroring  bool   `json:"-" tf:"ha_persistence_mirroring"`
	PersistenceEntryTimeout int    `json:"-" tf:"persistence_entry_timeout"`
	SharePersistence        bool   `json:"-" tf:"share_persistence"`
}

type CreateCookieProfileConfig struct {
	ServiceType      string `json:"-" tf:"service_type"`
	CookieName       string `json:"-" tf:"cookie_name"`
	CookieFallback   bool   `json:"-" tf:"cookie_fallback"`
	CookieGarbling   bool   `json:"-" tf:"cookie_garbling"`
	CookieMode       string `json:"-" tf:"cookie_mode"`
	CookieType       string `json:"-" tf:"cookie_type"`
	CookieDomain     string `json:"-" tf:"cookie_domain"`
	CookiePath       string `json:"-" tf:"cookie_path"`
	MaxIdleTime      int    `json:"-" tf:"max_idle_time"`
	SharePersistence bool   `json:"-" tf:"share_persistence"`
}

type CreateTCPProfileConfig struct {
	ServiceType            string `json:"-" tf:"service_type"`
	ConnectionCloseTimeout int    `json:"-" tf:"connection_close_timeout"`
	FastTCPIdleTimeout     int    `json:"-" tf:"fast_tcp_idle_timeout"`
	HaFlowMirroring        bool   `json:"-" tf:"ha_flow_mirroring"`
}

type CreateUDPProfileConfig struct {
	ServiceType        string `json:"-" tf:"service_type"`
	FastUDPIdleTimeout int    `json:"-" tf:"fast_udp_idle_timeout"`
	HaFlowMirroring    bool   `json:"-" tf:"ha_flow_mirroring"`
}

type Tags struct {
	Tag   string `json:"tag" tf:"tag"`
	Scope string `json:"scope" tf:"scope"`
}

type PoolTags struct {
	Tag   string `json:"name" tf:"tag"`
	Scope string `json:"value" tf:"scope"`
}

// Create LB Profile Resp
type CreateLBProfileResp struct {
	Success       bool          `json:"success"`
	LBProfileResp LBProfileResp `json:"loadBalancerProfile"`
}

type LBProfileResp struct {
	ID                  int           `json:"id" tf:"id,computed"`
	Name                string        `json:"name"`
	Category            string        `json:"category"`
	ServiceType         string        `json:"serviceType"`
	ServiceTypeDisplay  string        `json:"serviceTypeDisplay"`
	Visibility          string        `json:"visibility"`
	Description         string        `json:"description"`
	InternalID          string        `json:"internalId"`
	ExternalID          string        `json:"externalId"`
	Enabled             bool          `json:"enabled"`
	InsertXforwardedFor bool          `json:"insertXforwardedFor"`
	Editable            bool          `json:"editable"`
	LBProfileConfig     profileConfig `json:"config"`
}

type profileConfig struct {
	HTTPIdleTimeout        int    `json:"httpIdleTimeout"`
	ResponseHeaderSize     int    `json:"responseHeaderSize"`
	SharePersistence       bool   `json:"sharePersistence"`
	RequestHeaderSize      int    `json:"requestHeaderSize"`
	HaPersistenceMirroring bool   `json:"haPersistenceMirroring"`
	PreferServerCipher     bool   `json:"preferServerCipher"`
	ProfileType            string `json:"profileType"`
	CookieGarbling         bool   `json:"cookieGarbling"`
	NtlmAuthentication     bool   `json:"ntlmAuthentication"`
	HaFlowMirroring        bool   `json:"haFlowMirroring"`
	SessionCache           bool   `json:"sessionCache"`
	CookieFallback         bool   `json:"cookieFallback"`
	ResponseTimeout        int    `json:"responseTimeout"`
	PurgeEntries           bool   `json:"purgeEntries"`
	ConnectionCloseTimeout int    `json:"connectionCloseTimeout"`
	FastTCPIdleTimeout     int    `json:"fastTcpIdleTimeout"`
}

// Get LB Profile
type GetLBProfile struct {
	GetLBProfilesResp []GetLBProfilesResp `json:"loadBalancerProfiles"`
}

type GetLBProfilesResp struct {
	ID                  int    `json:"id" tf:"id,computed"`
	LbID                int    `json:"-" tf:"lb_id,computed"`
	Name                string `json:"name"`
	Category            string `json:"category"`
	ServiceType         string `json:"serviceType"`
	ServiceTypeDisplay  string `json:"serviceTypeDisplay"`
	Visibility          string `json:"visibility"`
	Description         string `json:"description"`
	InternalID          string `json:"internalId"`
	ExternalID          string `json:"externalId"`
	Enabled             bool   `json:"enabled"`
	InsertXforwardedFor bool   `json:"insertXforwardedFor"`
	Editable            bool   `json:"editable"`
	DateCreated         string `json:"dateCreated"`
	LastUpdated         string `json:"lastUpdated"`
}

type LBprofileConfig struct {
	HTTPIdleTimeout    int    `json:"httpIdleTimeout"`
	NtlmAuthentication bool   `json:"ntlmAuthentication"`
	RequestHeaderSize  int    `json:"requestHeaderSize"`
	ResponseHeaderSize int    `json:"responseHeaderSize"`
	ResponseTimeout    int    `json:"responseTimeout"`
	XForwardedFor      string `json:"xForwardedFor"`
	ProfileType        string `json:"profileType"`
	ResourceType       string `json:"resource_type"`
}

// Get LB Specific Profile
type GetLBSpecificProfile struct {
	GetLBSpecificProfilesResp GetLBSpecificProfilesResp `json:"loadBalancerProfile"`
}

type GetLBSpecificProfilesResp struct {
	ID                  int             `json:"-" tf:"id,computed"`
	LbID                int             `json:"-" tf:"lb_id"`
	Name                string          `json:"name"`
	Category            string          `json:"category"`
	ServiceType         string          `json:"serviceType"`
	ServiceTypeDisplay  string          `json:"serviceTypeDisplay"`
	Visibility          string          `json:"visibility"`
	Description         string          `json:"description"`
	InternalID          string          `json:"internalId"`
	ExternalID          string          `json:"externalId"`
	Enabled             bool            `json:"enabled"`
	InsertXforwardedFor bool            `json:"insertXforwardedFor"`
	Editable            bool            `json:"editable"`
	DateCreated         string          `json:"dateCreated"`
	LastUpdated         string          `json:"lastUpdated"`
	LBProfileConfig     LBprofileConfig `json:"config"`
}

// Get LB Monitors
type GetLBMonitors struct {
	GetLBMonitorsResp []GetLBMonitorsResp `json:"loadBalancerMonitors"`
}

type GetLBMonitorsResp struct {
	ID                 int       `json:"id" tf:"id,computed"`
	LbID               int       `json:"lb_id" tf:"lb_id,computed"`
	Name               string    `json:"name"`
	Visibility         string    `json:"visibility"`
	Description        string    `json:"description"`
	MonitorType        string    `json:"monitorType"`
	MonitorInterval    int       `json:"monitorInterval"`
	MonitorTimeout     int       `json:"monitorTimeout"`
	MonitorReverse     bool      `json:"monitorReverse"`
	MonitorTransparent bool      `json:"monitorTransparent"`
	MonitorAdaptive    bool      `json:"monitorAdaptive"`
	InternalID         string    `json:"internalId"`
	ExternalID         string    `json:"externalId"`
	MonitorSource      string    `json:"monitorSource"`
	Status             string    `json:"status"`
	Enabled            bool      `json:"enabled"`
	MaxRetry           int       `json:"maxRetry"`
	FallCount          int       `json:"fallCount"`
	RiseCount          int       `json:"riseCount"`
	DataLength         int       `json:"dataLength"`
	DateCreated        string    `json:"dateCreated"`
	LastUpdated        string    `json:"lastUpdated"`
	LoadBalancer       LBMonitor `json:"loadBalancer"`
}

// Get Specific LB Monitor
type GetSpecificLBMonitor struct {
	GetSpecificLBMonitorResp GetSpecificLBMonitorResp `json:"loadBalancerMonitor"`
}

type GetSpecificLBMonitorResp struct {
	ID                  int       `json:"-" tf:"id,computed"`
	LbID                int       `json:"-" tf:"lb_id"`
	Name                string    `json:"name"`
	Visibility          string    `json:"visibility"`
	Description         string    `json:"description"`
	Timeout             int       `json:"monitorTimeout"`
	Interval            int       `json:"monitorInterval"`
	RequestVersion      string    `json:"sendVersion"`
	RequestMethod       string    `json:"sendType"`
	ResponseStatusCodes string    `json:"receiveCode"`
	ResponseData        string    `json:"receiveData"`
	RequestURL          string    `json:"monitorDestination"`
	RequestBody         string    `json:"sendData"`
	AliasPort           int       `json:"aliasPort"`
	RiseCount           int       `json:"riseCount"`
	FallCount           int       `json:"fallCount"`
	DataLength          int       `json:"dataLength"`
	InternalID          string    `json:"internalId"`
	MonitorSource       string    `json:"monitorSource"`
	Status              string    `json:"status"`
	Enabled             bool      `json:"enabled"`
	DateCreated         string    `json:"dateCreated"`
	LastUpdated         string    `json:"lastUpdated"`
	LoadBalancer        LBMonitor `json:"loadBalancer"`
}

// Create LB Pool
type CreateLBPool struct {
	CreateLBPoolReq CreateLBPoolReq `json:"loadBalancerPool"`
}

type CreateLBPoolReq struct {
	ID          int         `json:"id" tf:"id,computed"`
	LbID        int         `json:"-" tf:"lb_id"`
	Name        string      `json:"name" tf:"name"`
	Description string      `json:"description" tf:"description"`
	VipBalance  string      `json:"vipBalance" tf:"algorithm"`
	MinActive   int         `json:"minActive" tf:"min_active_members"`
	PoolConfig  *PoolConfig `json:"config" tf:"config,sub"`
	Tag         []PoolTags  `json:"tags" tf:"tags,sub"`
}

type PoolConfig struct {
	SnatTranslationType   string       `json:"snatTranslationType" tf:"snat_translation_type"`
	PassiveMonitorPath    int          `json:"passiveMonitorPath" tf:"passive_monitor_path"`
	ActiveMonitorPaths    int          `json:"activeMonitorPaths" tf:"active_monitor_paths"`
	TCPMultiplexing       bool         `json:"tcpMultiplexing" tf:"tcp_multiplexing"`
	TCPMultiplexingNumber int          `json:"tcpMultiplexingNumber" tf:"tcp_multiplexing_number"`
	SnatIPAddress         string       `json:"snatIpAddress" tf:"snat_ip_address"`
	MemberGroup           *MemberGroup `json:"memberGroup" tf:"member_group,sub"`
}

type MemberGroup struct {
	Group            string `json:"path" tf:"group"`
	MaxIPListSize    int    `json:"MaxIPListSize" tf:"max_ip_list_size"`
	IPRevisionFilter string `json:"ipRevisionFilter" tf:"ip_revision_filter"`
	Port             int    `json:"port" tf:"port"`
}

// Create LB Pool Resp
type CreateLBPoolResp struct {
	Success    bool       `json:"success"`
	LBPoolResp LBPoolResp `json:"loadBalancerPool"`
}

type LBPoolResp struct {
	ID               int     `json:"id" tf:"id,computed"`
	LbID             int     `json:"lb_id" tf:"lb_id"`
	Name             string  `json:"name"`
	Category         string  `json:"category"`
	Visibility       string  `json:"visibility"`
	Description      string  `json:"description"`
	InternalID       string  `json:"internalId"`
	ExternalID       string  `json:"externalId"`
	Enabled          bool    `json:"enabled"`
	VipBalance       string  `json:"vipBalance"`
	MinActive        int     `json:"minActive"`
	NumberActive     int     `json:"numberActive"`
	NumberInService  int     `json:"numberInService"`
	HealthScore      float32 `json:"healthScore"`
	PerformanceScore float32 `json:"performanceScore"`
	HealthPenalty    float32 `json:"healthPenalty"`
	SecurityPenalty  float32 `json:"securityPenalty"`
	ErrorPenalty     float32 `json:"errorPenalty"`
	Status           string  `json:"status"`
	DateCreated      string  `json:"dateCreated"`
	LastUpdated      string  `json:"lastUpdated"`
}

// Get LB Pools
type GetLBPools struct {
	GetLBPoolsResp []GetLBPoolsResp `json:"loadBalancerPools"`
}

type GetLBPoolsResp struct {
	ID               int      `json:"id" tf:"id,computed"`
	LbID             int      `json:"-" tf:"lb_id,computed"`
	Name             string   `json:"name"`
	Visibility       string   `json:"visibility"`
	Description      string   `json:"description"`
	InternalID       string   `json:"internalId"`
	ExternalID       string   `json:"externalId"`
	Enabled          bool     `json:"enabled"`
	VipBalance       string   `json:"vipBalance"`
	MinActive        int      `json:"minActive"`
	NumberActive     int      `json:"numberActive"`
	NumberInService  int      `json:"numberInService"`
	HealthScore      float32  `json:"healthScore"`
	PerformanceScore float32  `json:"performanceScore"`
	HealthPenalty    float32  `json:"healthPenalty"`
	SecurityPenalty  float32  `json:"securityPenalty"`
	ErrorPenalty     float32  `json:"errorPenalty"`
	Status           string   `json:"status"`
	DateCreated      string   `json:"dateCreated"`
	LastUpdated      string   `json:"lastUpdated"`
	Meta             MetaInfo `json:"meta"`
}

// Get Specific LB Pools
type GetSpecificLBPool struct {
	GetSpecificLBPoolResp GetSpecificLBPoolResp `json:"loadBalancerPool"`
}

type GetSpecificLBPoolResp struct {
	ID               int     `json:"-" tf:"id,computed"`
	LbID             int     `json:"-" tf:"lb_id,computed"`
	Name             string  `json:"name"`
	Visibility       string  `json:"visibility"`
	Description      string  `json:"description"`
	InternalID       string  `json:"internalId"`
	ExternalID       string  `json:"externalId"`
	Enabled          bool    `json:"enabled"`
	VipBalance       string  `json:"vipBalance"`
	MinActive        int     `json:"minActive"`
	NumberActive     int     `json:"numberActive"`
	NumberInService  int     `json:"numberInService"`
	HealthScore      float32 `json:"healthScore"`
	PerformanceScore float32 `json:"performanceScore"`
	HealthPenalty    float32 `json:"healthPenalty"`
	SecurityPenalty  float32 `json:"securityPenalty"`
	ErrorPenalty     float32 `json:"errorPenalty"`
	Status           string  `json:"status"`
	DateCreated      string  `json:"dateCreated"`
	LastUpdated      string  `json:"lastUpdated"`
}

// CREATE LB Virtual servers

type CreateLBVirtualServers struct {
	CreateLBVirtualServersReq CreateLBVirtualServersReq `json:"loadBalancerInstance"`
}

type CreateLBVirtualServersReq struct {
	ID          int    `json:"id" tf:"id,computed"`
	LbID        int    `json:"-" tf:"lb_id"`
	Description string `json:"description" tf:"description"`
	VipName     string `json:"vipName" tf:"name"`
	VipAddress  string `json:"vipAddress" tf:"vip_address"`
	Pool        int    `json:"pool" tf:"pool"`
	VipPort     string `json:"vipPort" tf:"vip_port"`
	VipHostName string `json:"vipHostname" tf:"vip_host_name"`

	VipProtocol                  string                        `json:"vipProtocol" tf:"type"`
	TCPApplicationProfileConfig  *TCPApplicationProfileConfig  `json:"-" tf:"tcp_application_profile,sub"`
	UDPApplicationProfileConfig  *UDPApplicationProfileConfig  `json:"-" tf:"udp_application_profile,sub"`
	HTTPApplicationProfileConfig *HTTPApplicationProfileConfig `json:"-" tf:"http_application_profile,sub"`

	Persistence                      string                            `json:"-"  tf:"persistence"`
	CookiePersistenceProfileConfig   *CookiePersistenceProfileConfig   `json:"-" tf:"cookie_persistence_profile,sub"`
	SourceipPersistenceProfileConfig *SourceipPersistenceProfileConfig `json:"-" tf:"sourceip_persistence_profile,sub"`

	SSLServerCert   int              `json:"sslServerCert" tf:"ssl_server_cert"`
	SSLServerConfig *SSLServerConfig `json:"-" tf:"ssl_server_config,sub"`

	SSLCert         int              `json:"sslCert" tf:"ssl_client_cert"`
	SSLClientConfig *SSLClientConfig `json:"-" tf:"ssl_client_config,sub"`

	VirtualServerConfig VirtualServerConfig `json:"config" tf:"config,sub"`
}

type TCPApplicationProfileConfig struct {
	ApplicationProfile int `json:"-" tf:"application_profile"`
}

type UDPApplicationProfileConfig struct {
	ApplicationProfile int `json:"-" tf:"application_profile"`
}

type HTTPApplicationProfileConfig struct {
	ApplicationProfile int `json:"-" tf:"application_profile"`
}

type CookiePersistenceProfileConfig struct {
	PersistenceProfile int `json:"-" tf:"persistence_profile"`
}

type SourceipPersistenceProfileConfig struct {
	PersistenceProfile int `json:"-" tf:"persistence_profile"`
}

type SSLClientConfig struct {
	SSLClientProfile int `json:"-" tf:"ssl_client_profile"`
}

type SSLServerConfig struct {
	SSLServerProfile int `json:"-" tf:"ssl_server_profile"`
}

type VirtualServerConfig struct {
	Persistence        string `json:"persistence"`
	PersistenceProfile int    `json:"persistenceProfile"`
	ApplicationProfile int    `json:"applicationProfile"`
	SSLClientProfile   int    `json:"sslClientProfile"`
	SSLServerProfile   int    `json:"sslServerProfile"`
}

// CREATE LB Virtual Server Resp
type LBVirtualServersResp struct {
	CreateLBVirtualServersResp CreateLBVirtualServersResp `json:"loadBalancerInstance"`
	Success                    bool                       `json:"success"`
}

type CreateLBVirtualServersResp struct {
	ID                 int           `json:"id" tf:"id,computed"`
	Name               string        `json:"name"`
	Description        string        `json:"description"`
	Active             bool          `json:"active"`
	Sticky             bool          `json:"sticky"`
	SslEnabled         bool          `json:"sslEnabled"`
	ExternalAddress    bool          `json:"externalAddress"`
	VipName            string        `json:"vipName"`
	VipAddress         string        `json:"vipAddress"`
	VipProtocol        string        `json:"vipProtocol"`
	VipMode            string        `json:"vipMode"`
	VipPort            int           `json:"vipPort"`
	VipShared          bool          `json:"vipShared"`
	VipDirectAddress   string        `json:"vipDirectAddress"`
	Removing           bool          `json:"removing"`
	VirtualServiceName string        `json:"virtualServiceName"`
	VipSource          string        `json:"vipSource"`
	DateCreated        string        `json:"dateCreated"`
	LastUpdated        string        `json:"lastUpdated"`
	LoadBalancer       LBMonitor     `json:"loadBalancer"`
	VSConfig           VSConfig      `json:"config"`
	VSPool             VSPool        `json:"pool"`
	SSLCert            SSLCert       `json:"sslCert"`
	SSLServerCert      SSLServerCert `json:"sslServerCert"`
}

// GET LB Virtual servers

type GetLBVirtualServers struct {
	GetLBVirtualServersResp []GetLBVirtualServersResp `json:"loadBalancerInstances"`
}

type GetLBVirtualServersResp struct {
	ID                 int           `json:"id"`
	Name               string        `json:"name"`
	Description        string        `json:"description"`
	Active             bool          `json:"active"`
	Sticky             bool          `json:"sticky"`
	SslEnabled         bool          `json:"sslEnabled"`
	ExternalAddress    bool          `json:"externalAddress"`
	VipName            string        `json:"vipName"`
	VipAddress         string        `json:"vipAddress"`
	VipProtocol        string        `json:"vipProtocol"`
	VipMode            string        `json:"vipMode"`
	VipPort            int           `json:"vipPort"`
	VipShared          bool          `json:"vipShared"`
	VipDirectAddress   string        `json:"vipDirectAddress"`
	Removing           bool          `json:"removing"`
	VirtualServiceName string        `json:"virtualServiceName"`
	VipSource          string        `json:"vipSource"`
	DateCreated        string        `json:"dateCreated"`
	LastUpdated        string        `json:"lastUpdated"`
	LoadBalancer       LBMonitor     `json:"loadBalancer"`
	VSConfig           VSConfig      `json:"config"`
	VSPool             VSPool        `json:"pool"`
	SSLCert            SSLCert       `json:"sslCert"`
	SSLServerCert      SSLServerCert `json:"sslServerCert"`
	Meta               MetaInfo      `json:"meta"`
}

type VSConfig struct {
	Persistence        string   `json:"persistence"`
	PersistenceProfile int      `json:"persistenceProfile"`
	SslServerProfile   int      `json:"sslServerProfile"`
	SslClientProfile   int      `json:"sslClientProfile"`
	ApplicationProfile int      `json:"applicationProfile"`
	Monitors           []string `json:"monitors"`
}

type VSPool struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type SSLCert struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type SSLServerCert struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// GET LB Specific Virtual servers

type GetSpecificLBVirtualServers struct {
	GetSpecificLBVirtualServersResp GetSpecificLBVirtualServersResp `json:"loadBalancerInstance"`
}

type GetSpecificLBVirtualServersResp struct {
	ID                 int           `json:"id" tf:"id,computed"`
	Name               string        `json:"name"`
	Description        string        `json:"description"`
	Active             bool          `json:"active"`
	Sticky             bool          `json:"sticky"`
	SslEnabled         bool          `json:"sslEnabled"`
	ExternalAddress    bool          `json:"externalAddress"`
	VipName            string        `json:"vipName"`
	VipAddress         string        `json:"vipAddress"`
	VipProtocol        string        `json:"vipProtocol"`
	VipMode            string        `json:"vipMode"`
	VipPort            int           `json:"vipPort"`
	VipShared          bool          `json:"vipShared"`
	VipDirectAddress   string        `json:"vipDirectAddress"`
	Removing           bool          `json:"removing"`
	VirtualServiceName string        `json:"virtualServiceName"`
	VipSource          string        `json:"vipSource"`
	DateCreated        string        `json:"dateCreated"`
	LastUpdated        string        `json:"lastUpdated"`
	LoadBalancer       LBMonitor     `json:"loadBalancer"`
	VSConfig           VSConfig      `json:"config"`
	VSPool             VSPool        `json:"pool"`
	SSLCert            SSLCert       `json:"sslCert"`
	SSLServerCert      SSLServerCert `json:"sslServerCert"`
}

type GetMemeberGroupForPool struct {
	MemeberGroup []MemeberGroups `json:"groups"`
}

type MemeberGroups struct {
	ID          int         `json:"id" tf:"id,computed"`
	Name        string      `json:"name"`
	Description interface{} `json:"description"`
	InternalID  string      `json:"internalId"`
	ExternalID  string      `json:"externalId" tf:"external_id,computed"`
}

type GetSSLCertificates struct {
	Certificates []Certificates `json:"certificates"`
}

type Certificates struct {
	ID          int    `json:"id" tf:"id,computed"`
	Name        string `json:"name"`
	Description string `json:"description"`
}
