// (C) Copyright 2021-2023 Hewlett Packard Enterprise Development LP

package models

type ServersResponse struct {
	Server      []Server            `json:"servers"`
	Stats       map[int]ServerStats `json:"stats"`
	MultiTenant bool                `json:"multiTenant"`
	Meta        *Meta               `json:"meta"`
}

type GetSpecificServerResponse struct {
	Server Server `json:"server"`
}

type Server struct {
	ID                int                `json:"id"`
	UUID              string             `json:"uuid"`
	ExternalID        string             `json:"externalId"`
	InternalID        string             `json:"internalId"`
	ExternalUniqueID  string             `json:"externalUniqueId"`
	Name              string             `json:"name"`
	ExternalName      string             `json:"externalName"`
	Hostname          string             `json:"hostname"`
	ParentServer      *ParentServer      `json:"parentServer"`
	AccountID         int                `json:"accountId"`
	Account           *Account           `json:"account"`
	Owner             *Owner             `json:"owner"`
	Zone              *Zone              `json:"zone"`
	Plan              *Plan              `json:"plan"`
	ComputeServerType *ComputeServerType `json:"computeServerType"`
	Visibility        string             `json:"visibility"`
	Description       interface{}        `json:"description"`
	ZoneID            int                `json:"zoneId"`
	SiteID            int                `json:"siteId"`
	ResourcePoolID    interface{}        `json:"resourcePoolId"`
	FolderID          int                `json:"folderId"`
	SSHHost           interface{}        `json:"sshHost"`
	SSHPort           int                `json:"sshPort"`
	ExternalIP        string             `json:"externalIp"`
	InternalIP        string             `json:"internalIp"`
	VolumeID          interface{}        `json:"volumeId"`
	Platform          interface{}        `json:"platform"`
	PlatformVersion   interface{}        `json:"platformVersion"`
	SSHUsername       interface{}        `json:"sshUsername"`
	SSHPassword       interface{}        `json:"sshPassword"`
	OsDevice          string             `json:"osDevice"`
	OsType            string             `json:"osType"`
	DataDevice        string             `json:"dataDevice"`
	LvmEnabled        bool               `json:"lvmEnabled"`
	APIKey            string             `json:"apiKey"`
	SoftwareRaid      bool               `json:"softwareRaid"`
	DateCreated       string             `json:"dateCreated"`
	LastUpdated       string             `json:"lastUpdated"`
	Stats             *ServerStats       `json:"stats"`
	Status            string             `json:"status"`
	StatusMessage     interface{}        `json:"statusMessage"`
	ErrorMessage      interface{}        `json:"errorMessage"`
	StatusDate        string             `json:"statusDate"`
	StatusPercent     interface{}        `json:"statusPercent"`
	StatusEta         interface{}        `json:"statusEta"`
	PowerState        string             `json:"powerState"`
	AgentInstalled    bool               `json:"agentInstalled"`
	LastAgentUpdate   interface{}        `json:"lastAgentUpdate"`
	AgentVersion      interface{}        `json:"agentVersion"`
	MaxCores          int                `json:"maxCores"`
	MaxMemory         int64              `json:"maxMemory"`
	MaxStorage        int64              `json:"maxStorage"`
	MaxCPU            interface{}        `json:"maxCpu"`
	HourlyCost        float64            `json:"hourlyCost"`
	HourlyPrice       float64            `json:"hourlyPrice"`
	SourceImage       *SourceImage       `json:"sourceImage"`
	ServerOs          *ServerOs          `json:"serverOs"`
	Interfaces        []Interfaces       `json:"interfaces"`
	Labels            []interface{}      `json:"labels"`
	Tags              []interface{}      `json:"tags"`
	Enabled           bool               `json:"enabled"`
	TagCompliant      interface{}        `json:"tagCompliant"`
	Containers        []int              `json:"containers"`
}

type ParentServer struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Account struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Owner struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
}

type Zone struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Plan struct {
	ID   int    `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}

type ComputeServerType struct {
	ID             int    `json:"id"`
	Code           string `json:"code"`
	Name           string `json:"name"`
	Managed        bool   `json:"managed"`
	ExternalDelete bool   `json:"externalDelete"`
}

type ServerStats struct {
	TS             string  `json:"ts"`
	MaxMemory      int64   `json:"maxMemory"`
	UsedMemory     int64   `json:"usedMemory"`
	MaxStorage     int64   `json:"maxStorage"`
	UsedStorage    int64   `json:"usedStorage"`
	CPUUsage       float64 `json:"cpuUsage"`
	FreeMemory     int64   `json:"freeMemory"`
	ReservedMemory int64   `json:"reservedMemory"`
}

type SourceImage struct {
	ID   int    `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}

type ServerOs struct {
	ID          int         `json:"id"`
	Code        string      `json:"code"`
	Name        string      `json:"name"`
	Description interface{} `json:"description"`
	Vendor      string      `json:"vendor"`
	Category    string      `json:"category"`
	OsFamily    interface{} `json:"osFamily"`
	OsVersion   string      `json:"osVersion"`
	BitCount    int         `json:"bitCount"`
	Platform    string      `json:"platform"`
}

type Interfaces struct {
	ID                int         `json:"id"`
	RefType           interface{} `json:"refType"`
	RefID             interface{} `json:"refId"`
	Name              string      `json:"name"`
	InternalID        string      `json:"internalId"`
	ExternalID        string      `json:"externalId"`
	UniqueID          string      `json:"uniqueId"`
	PublicIPAddress   string      `json:"publicIpAddress"`
	PublicIpv6Address interface{} `json:"publicIpv6Address"`
	IPAddress         string      `json:"ipAddress"`
	Ipv6Address       string      `json:"ipv6Address"`
	IPSubnet          interface{} `json:"ipSubnet"`
	Ipv6Subnet        interface{} `json:"ipv6Subnet"`
	Description       interface{} `json:"description"`
	Dhcp              bool        `json:"dhcp"`
	Active            bool        `json:"active"`
	PoolAssigned      bool        `json:"poolAssigned"`
	PrimaryInterface  bool        `json:"primaryInterface"`
	Network           interface{} `json:"network"`
	Subnet            interface{} `json:"subnet"`
	NetworkGroup      interface{} `json:"networkGroup"`
	NetworkPosition   interface{} `json:"networkPosition"`
	NetworkPool       interface{} `json:"networkPool"`
	NetworkDomain     interface{} `json:"networkDomain"`
	Type              *Type       `json:"type"`
	IPMode            string      `json:"ipMode"`
	MacAddress        string      `json:"macAddress"`
}

type Type struct {
	ID   int    `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}
