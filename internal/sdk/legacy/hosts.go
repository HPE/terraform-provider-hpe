package morpheus

import (
	"fmt"
	"time"
)

var (
	// HostsPath is the API endpoint for hosts
	HostsPath     = "/api/servers"
	HostTypesPath = "/api/server-types"
)

// Server structures for use in request and response payloads
type Host struct {
	ID               int64  `json:"id"`
	UUID             string `json:"uuid"`
	ExternalId       string `json:"externalId"`
	InternalId       string `json:"internalId"`
	ExternalUniqueId string `json:"externalUniqueId"`
	Name             string `json:"name"`
	ExternalName     string `json:"externalName"`
	Hostname         string `json:"hostname"`
	ParentServer     struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	} `json:"parentServer"`
	AccountId int64 `json:"accountId"`
	Account   struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	} `json:"account"`
	Owner struct {
		ID       int64  `json:"id"`
		Username string `json:"username"`
	} `json:"owner"`
	Zone struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	} `json:"zone"`
	Plan struct {
		ID   int64  `json:"id"`
		Code string `json:"code"`
		Name string `json:"name"`
	} `json:"plan"`
	ComputeServerType struct {
		ID             int64  `json:"id"`
		Code           string `json:"code"`
		Name           string `json:"name"`
		Managed        bool   `json:"managed"`
		ExternalDelete bool   `json:"externalDelete"`
	} `json:"computeServerType"`
	Visibility      string      `json:"visibility"`
	Description     string      `json:"description"`
	ZoneId          int64       `json:"zoneId"`
	SiteId          int64       `json:"siteId"`
	ResourcePoolId  int64       `json:"resourcePoolId"`
	FolderId        int64       `json:"folderId"`
	SshHost         string      `json:"sshHost"`
	SshPort         int64       `json:"sshPort"`
	ExternalIp      string      `json:"externalIp"`
	InternalIp      string      `json:"internalIp"`
	VolumeId        interface{} `json:"volumeId"`
	Platform        string      `json:"platform"`
	PlatformVersion string      `json:"platformVersion"`
	SshUsername     string      `json:"sshUsername"`
	SshPassword     string      `json:"sshPassword"`
	SshPasswordHash string      `json:"sshPasswordHash"`
	OsDevice        string      `json:"osDevice"`
	OsType          string      `json:"osType"`
	DataDevice      string      `json:"dataDevice"`
	LVMEnabled      bool        `json:"lvmEnabled"`
	ApiKey          string      `json:"apiKey"`
	SoftwareRaid    bool        `json:"softwareRaid"`
	DateCreated     time.Time   `json:"dateCreated"`
	LastUpdated     time.Time   `json:"lastUpdated"`
	Stats           struct {
		Ts               string  `json:"ts"`
		FreeMemory       int64   `json:"freeMemory"`
		UsedMemory       int64   `json:"usedMemory"`
		FreeSwap         int64   `json:"freeSwap"`
		UsedSwap         int64   `json:"usedSwap"`
		CpuIdleTime      int64   `json:"cpuIdleTime"`
		CpuSystemTime    int64   `json:"cpuSystemTime"`
		CpuUserTime      int64   `json:"cpuUserTime"`
		CpuTotalTime     int64   `json:"cpuTotalTime"`
		CpuUsage         float64 `json:"cpuUsage"`
		MaxStorage       int64   `json:"maxStorage"`
		UsedStorage      int64   `json:"usedStorage"`
		ReservedStorage  int64   `json:"reservedStorage"`
		NetTxUsage       int64   `json:"netTxUsage"`
		NetRxUsage       int64   `json:"netRxUsage"`
		NetworkBandwidth int64   `json:"networkBandwidth"`
	} `json:"stats"`
	Status                 string      `json:"status"`
	StatusMessage          string      `json:"statusMessage"`
	ErrorMessage           string      `json:"errorMessage"`
	StatusDate             time.Time   `json:"statusDate"`
	StatusPercent          interface{} `json:"statusPercent"`
	StatusEta              interface{} `json:"statusEta"`
	PowerState             string      `json:"powerState"`
	AgentInstalled         bool        `json:"agentInstalled"`
	LastAgentUpdate        time.Time   `json:"lastAgentUpdate"`
	AgentVersion           string      `json:"agentVersion"`
	MaxCores               int64       `json:"maxCores"`
	CoresPerSocket         int64       `json:"coresPerSocket"`
	MaxMemory              int64       `json:"maxMemory"`
	MaxStorage             int64       `json:"maxStorage"`
	MaxCpu                 int64       `json:"maxCpu"`
	ManageInternalFirewall bool        `json:"manageInternalFirewall"`
	EnableLogs             bool        `json:"enableLogs"`
	HourlyCost             float64     `json:"hourlyCost"`
	HourlyPrice            float64     `json:"hourlyPrice"`
	SourceImage            struct {
		ID   int64  `json:"id"`
		Code string `json:"code"`
		Name string `json:"name"`
	} `json:"sourceImage"`
	ServerOs struct {
		ID          int64  `json:"id"`
		Code        string `json:"code"`
		Name        string `json:"name"`
		Description string `json:"description"`
		Vendor      string `json:"vendor"`
		Category    string `json:"category"`
		OsFamily    string `json:"osFamily"`
		OsVersion   string `json:"osVersion"`
		BitCount    int64  `json:"bitCount"`
		Platform    string `json:"platform"`
	} `json:"serverOs"`
	Volumes []struct {
		ID                   int64  `json:"id"`
		Name                 string `json:"name"`
		ControllerId         int64  `json:"controllerId"`
		ControllerMountPoint string `json:"controllerMountPoint"`
		Resizeable           bool   `json:"resizeable"`
		PlanResizable        bool   `json:"planResizable"`
		RootVolume           bool   `json:"rootVolume"`
		UnitNumber           string `json:"unitNumber"`
		TypeId               int64  `json:"typeId"`
		ConfigurableIOPS     bool   `json:"configurableIOPS"`
		DatastoreId          int64  `json:"datastoreId"`
		MaxStorage           int64  `json:"maxStorage"`
		DisplayOrder         int64  `json:"displayOrder"`
		MaxIOPS              int64  `json:"maxIOPS"`
		UUID                 string `json:"uuid"`
		UniqueId             string `json:"uniqueId"`
		ExternalId           string `json:"externalId"`
		InternalId           string `json:"internalId"`
	} `json:"volumes"`
	Controllers []struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
		Type struct {
			ID   int64  `json:"id"`
			Code string `json:"code"`
			Name string `json:"name"`
		} `json:"type"`
		MaxDevices         int64 `json:"maxDevices"`
		ReservedUnitNumber int64 `json:"reservedUnitNumber"`
	} `json:"controllers"`
	Interfaces []struct {
		ID                int64  `json:"id"`
		RefType           string `json:"refType"`
		RefId             int64  `json:"refId"`
		Name              string `json:"name"`
		InternalId        string `json:"internalId"`
		ExternalId        string `json:"externalId"`
		UniqueId          string `json:"uniqueId"`
		PublicIpAddress   string `json:"publicIpAddress"`
		PublicIpv6Address string `json:"publicIpv6Address"`
		IpAddress         string `json:"ipAddress"`
		Ipv6Address       string `json:"ipv6Address"`
		IpSubnet          string `json:"ipSubnet"`
		Ipv6Subnet        string `json:"ipv6Subnet"`
		Description       string `json:"description"`
		Dhcp              bool   `json:"dhcp"`
		Active            bool   `json:"active"`
		PoolAssigned      bool   `json:"poolAssigned"`
		PrimaryInterface  bool   `json:"primaryInterface"`
		Network           struct {
			ID   int64  `json:"id"`
			Name string `json:"name"`
		} `json:"network"`
		Subnet          string      `json:"subnet"`
		NetworkGroup    interface{} `json:"networkGroup"`
		NetworkPosition interface{} `json:"networkPosition"`
		NetworkPool     struct {
			ID   int64  `json:"id"`
			Name string `json:"name"`
		} `json:"networkPool"`
		NetworkDomain string `json:"networkDomain"`
		Type          struct {
			ID   int64  `json:"id"`
			Code string `json:"code"`
			Name string `json:"name"`
		} `json:"type"`
		IpMode     string `json:"ipMode"`
		MacAddress string `json:"macAddress"`
	} `json:"interfaces"`
	Labels       []string      `json:"labels"`
	Tags         []interface{} `json:"tags"`
	Enabled      bool          `json:"enabled"`
	TagCompliant bool          `json:"tagCompliant"`
	Containers   []int64       `json:"containers"`
	Config       struct {
		PoolProviderType     string      `json:"poolProviderType"`
		IsVpcSelectable      bool        `json:"isVpcSelectable"`
		SmbiosAssetTag       string      `json:"smbiosAssetTag"`
		IsEC2                bool        `json:"isEC2"`
		ResourcePoolId       int64       `json:"resourcePoolId"`
		HostId               interface{} `json:"hostId"`
		CreateUser           bool        `json:"createUser"`
		NestedVirtualization interface{} `json:"nestedVirtualization"`
		VmwareFolderId       string      `json:"vmwareFolderId"`
		NoAgent              bool        `json:"noAgent"`
		PowerScheduleType    interface{} `json:"powerScheduleType"`
	} `json:"config"`
	GuestConsolePreferred    bool   `json:"guestConsolePreferred"`
	GuestConsoleType         string `json:"guestConsoleType"`
	GuestConsoleUsername     string `json:"guestConsoleUsername"`
	GuestConsolePassword     string `json:"guestConsolePassword"`
	GuestConsolePasswordHash string `json:"guestConsolePasswordHash"`
	GuestConsolePort         string `json:"guestConsolePort"`
}

type HostStats struct {
	UsedStorage     int64   `json:"usedStorage"`
	ReservedStorage int64   `json:"reservedStorage"`
	MaxStorage      int64   `json:"maxStorage"`
	UsedMemory      int64   `json:"usedMemory"`
	ReservedMemory  int64   `json:"reservedMemory"`
	MaxMemory       int64   `json:"maxMemory"`
	CpuUsage        float64 `json:"cpuUsage"`
}

type HostType struct {
	ID                  int64  `json:"id"`
	Code                string `json:"code"`
	Name                string `json:"name"`
	Description         string `json:"description"`
	NodeType            string `json:"nodeType"`
	Platform            string `json:"platform"`
	Enabled             bool   `json:"enabled"`
	Selectable          bool   `json:"selectable"`
	ExternalDelete      bool   `json:"externalDelete"`
	Managed             bool   `json:"managed"`
	ControlPower        bool   `json:"controlPower"`
	ControlSuspend      bool   `json:"controlSuspend"`
	Creatable           bool   `json:"creatable"`
	HasAgent            bool   `json:"hasAgent"`
	VmHypervisor        bool   `json:"vmHypervisor"`
	ContainerHypervisor bool   `json:"containerHypervisor"`
	BareMetalHost       bool   `json:"bareMetalHost"`
	GuestVm             bool   `json:"guestVm"`
	HasAutomation       bool   `json:"hasAutomation"`
	ProvisionType       struct {
		ID           int64  `json:"id"`
		Code         string `json:"code"`
		Name         string `json:"name"`
		HasNetworks  bool   `json:"hasNetworks"`
		HasZonePools bool   `json:"hasZonePools"`
	} `json:"provisionType"`
	OptionTypes  []OptionType `json:"optionTypes"`
	DisplayOrder int64        `json:"displayOrder"`
}

// ListHostsResult structure parses the list hosts response payload
type ListHostsResult struct {
	Hosts       *[]Host     `json:"servers"`
	MultiTenant bool        `json:"multiTenant"`
	Meta        *MetaResult `json:"meta"`
}

type ListHostTypesResult struct {
	HostTypes *[]HostType `json:"serverTypes"`
	Meta      *MetaResult `json:"meta"`
}

type GetHostResult struct {
	Host  *Host      `json:"server"`
	Stats *HostStats `json:"stats"`
}

type GetHostTypeResult struct {
	HostType *HostType `json:"serverType"`
}

type CreateHostResult struct {
	Success bool              `json:"success"`
	Message string            `json:"msg"`
	Errors  map[string]string `json:"errors"`
	Host    *Host             `json:"server"`
}

type UpdateHostResult struct {
	CreateHostResult
}

type ConvertToManagedResult struct {
	CreateHostResult
	PublicKey string `json:"publicKey"`
}

type DeleteHostResult struct {
	DeleteResult
}

type StandardHostResult struct {
	Success bool              `json:"success"`
	Message string            `json:"msg"`
	Errors  map[string]string `json:"errors"`
}

// Client request methods
func (client *Client) ListHosts(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        HostsPath,
		QueryParams: req.QueryParams,
		Result:      &ListHostsResult{},
	})
}

func (client *Client) ListHostTypes(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        HostTypesPath,
		QueryParams: req.QueryParams,
		Result:      &ListHostTypesResult{},
	})
}

func (client *Client) GetHost(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        fmt.Sprintf("%s/%d", HostsPath, id),
		QueryParams: req.QueryParams,
		Result:      &GetHostResult{},
	})
}

func (client *Client) GetHostType(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        fmt.Sprintf("%s/%d", HostTypesPath, id),
		QueryParams: req.QueryParams,
		Result:      &GetHostTypeResult{},
	})
}

func (client *Client) CreateHost(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "POST",
		Path:        HostsPath,
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &CreateHostResult{},
	})
}

func (client *Client) UpdateHost(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "PUT",
		Path:        fmt.Sprintf("%s/%d", HostsPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &UpdateHostResult{},
	})
}

func (client *Client) DeleteHost(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "DELETE",
		Path:        fmt.Sprintf("%s/%d", HostsPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &DeleteHostResult{},
	})
}

func (client *Client) AssignHostToTenant(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "PUT",
		Path:        fmt.Sprintf("%s/%d/assign-account", HostsPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &UpdateHostResult{},
	})
}

func (client *Client) ResizeHost(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "PUT",
		Path:        fmt.Sprintf("%s/%d/resize", HostsPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &UpdateHostResult{},
	})
}

func (client *Client) StartHost(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "PUT",
		Path:        fmt.Sprintf("%s/%d/start", HostsPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &UpdateHostResult{},
	})
}

func (client *Client) StopHost(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "PUT",
		Path:        fmt.Sprintf("%s/%d/stop", HostsPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &UpdateHostResult{},
	})
}

func (client *Client) RestartHost(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "PUT",
		Path:        fmt.Sprintf("%s/%d/restart", HostsPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &UpdateHostResult{},
	})
}

func (client *Client) UpgradeHostAgent(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "PUT",
		Path:        fmt.Sprintf("%s/%d/upgrade", HostsPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &UpdateHostResult{},
	})
}

func (client *Client) ConvertHostToManaged(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "PUT",
		Path:        fmt.Sprintf("%s/%d/make-managed", HostsPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &ConvertToManagedResult{},
	})
}

func (client *Client) InstallHostAgent(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "PUT",
		Path:        fmt.Sprintf("%s/%d/install-agent", HostsPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &ConvertToManagedResult{},
	})
}

func (client *Client) EnableHostMaintenance(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "PUT",
		Path:        fmt.Sprintf("%s/%d/maintenance", HostsPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &StandardHostResult{},
	})
}

func (client *Client) LeaveHostMaintenance(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "PUT",
		Path:        fmt.Sprintf("%s/%d/leave-maintenance", HostsPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &StandardHostResult{},
	})
}

func (client *Client) ManageHostPlacement(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "PUT",
		Path:        fmt.Sprintf("%s/%d/placement", HostsPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &StandardHostResult{},
	})
}

// helper functions
func (client *Client) FindHostByName(name string) (*Response, error) {
	// Find by name, then get by ID
	resp, err := client.ListHosts(&Request{
		QueryParams: map[string]string{
			"name": name,
		},
	})
	if err != nil {
		return resp, err
	}
	listResult := resp.Result.(*ListHostsResult)
	hostCount := len(*listResult.Hosts)
	if hostCount != 1 {
		return resp, fmt.Errorf("found %d Hosts for %v", hostCount, name)
	}
	firstRecord := (*listResult.Hosts)[0]
	hostID := firstRecord.ID
	return client.GetHost(hostID, &Request{})
}
