package morpheus

import (
	"fmt"
	"time"
)

var (
	// InstancesPath is the API endpoint for instances
	InstancesPath = "/api/instances"
)

// Instance structures for use in request and response payloads
type Instance struct {
	ID        int64  `json:"id"`
	UUID      string `json:"uuid"`
	AccountId int64  `json:"accountId"`
	Tenant    struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	} `json:"tenant"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	DisplayName  string `json:"displayName"`
	InstanceType struct {
		ID       int64  `json:"id"`
		Name     string `json:"name"`
		Code     string `json:"code"`
		Category string `json:"category"`
		Image    string `json:"image"`
	} `json:"instanceType"`
	Layout struct {
		ID                int64  `json:"id"`
		Name              string `json:"name"`
		ProvisionTypeId   int64  `json:"provisionTypeId"`
		ProvisionTypeCode string `json:"provisionTypeCode"`
	} `json:"layout"`
	Group struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	} `json:"group"`
	Cloud struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
		Type string `json:"type"`
	} `json:"cloud"`
	Cluster struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
		Type struct {
			ID   int64  `json:"id"`
			Name string `json:"name"`
			Code string `json:"code"`
		} `json:"type"`
	} `json:"cluster"`
	Containers []int64 `json:"containers"`
	Servers    []int64 `json:"servers"`
	Resources  []struct {
		ID              int64    `json:"id"`
		UUID            string   `json:"uuid"`
		Code            string   `json:"code"`
		Category        string   `json:"category"`
		Name            string   `json:"name"`
		DisplayName     string   `json:"displayName"`
		Labels          []string `json:"labels"`
		ResourceVersion string   `json:"resourceVersion"`
		ResourceContext string   `json:"resourceContext"`
		Owner           struct {
			Id   int64  `json:"id"`
			Name string `json:"name"`
		} `json:"owner"`
		ResourceType string `json:"resourceType"`
		ResourceIcon string `json:"resourceIcon"`
		Type         struct {
			Id   int64  `json:"id"`
			Code string `json:"code"`
			Name string `json:"name"`
		} `json:"type"`
		Status     string `json:"status"`
		Enabled    bool   `json:"enabled"`
		ExternalId string `json:"externalId"`
	} `json:"resources"`
	ConnectionInfo []struct {
		Ip   string `json:"ip"`
		Port int64  `json:"port"`
		Name string `json:"name"`
	} `json:"connectionInfo"`
	Environment string                 `json:"environment"`
	Plan        InstancePlan           `json:"plan"`
	Config      map[string]interface{} `json:"config"`
	Labels      []string               `json:"labels"`
	Version     string                 `json:"instanceVersion"`
	Status      string                 `json:"status"`
	Owner       Owner                  `json:"owner"`
	Volumes     []struct {
		ID                   interface{} `json:"id"`
		Name                 string      `json:"name"`
		ShortName            string      `json:"shortName"`
		Description          string      `json:"description"`
		ControllerId         int64       `json:"controllerId"`
		ControllerMountPoint string      `json:"controllerMountPoint"`
		Resizeable           interface{} `json:"resizeable"`
		PlanResizable        interface{} `json:"planResizable"`
		Size                 interface{} `json:"size"`
		StorageType          interface{} `json:"storageType"`
		RootVolume           interface{} `json:"rootVolume"`
		UnitNumber           string      `json:"unitNumber"`
		DeviceName           string      `json:"deviceName"`
		DeviceDisplayName    string      `json:"deviceDisplayName"`
		Type                 struct {
			ID   int64  `json:"id"`
			Code string `json:"code"`
			Name string `json:"name"`
		} `json:"type"`
		TypeId           int64  `json:"typeId"`
		Category         string `json:"category"`
		Status           string `json:"status"`
		StatusMessage    string `json:"statusMessage"`
		ConfigurableIOPS bool   `json:"configurableIOPS"`
		MaxStorage       int64  `json:"maxStorage"`
		DisplayOrder     int64  `json:"displayOrder"`
		MaxIOPS          string `json:"maxIOPS"`
		Uuid             string `json:"uuid"`
		Active           bool   `json:"active"`
		Zone             struct {
			ID   int64  `json:"id"`
			Name string `json:"name"`
		} `json:"zone"`
		ZoneId    int64 `json:"zoneId"`
		Datastore struct {
			ID   int64  `json:"id"`
			Name string `json:"name"`
		} `json:"datastore"`
		DatastoreId   interface{} `json:"datastoreId"`
		StorageGroup  string      `json:"storageGroup"`
		Namespace     string      `json:"namespace"`
		StorageServer string      `json:"storageServer"`
		Source        string      `json:"source"`
		Owner         struct {
			ID   int64  `json:"id"`
			Name string `json:"name"`
		} `json:"owner"`
	} `json:"volumes"`
	Interfaces  []NetworkInterface `json:"interfaces"`
	Controllers []StorageControler `json:"controllers"`
	Tags        []struct {
		Id    int64       `json:"id"`
		Name  string      `json:"name"`
		Value interface{} `json:"value"`
	} `json:"tags"`
	Metadata             []map[string]interface{} `json:"metadata"`
	EnvironmentVariables []struct {
		Name   string      `json:"name"`
		Value  interface{} `json:"value"`
		Export bool        `json:"export"`
		Masked bool        `json:"masked"`
	} `json:"evars"`
	CustomOptions   map[string]interface{} `json:"customOptions"`
	PendingRequests []struct {
		Id           int64  `json:"id"`
		Type         string `json:"type"`
		Status       string `json:"status"`
		Name         string `json:"name"`
		ExternalName string `json:"externalName"`
		DateCreated  string `json:"dateCreated"`
		DateApproved string `json:"dateApproved"`
		DateDenied   string `json:"dateDenied"`
		ApprovedBy   string `json:"approvedBy"`
		DeniedBy     string `json:"deniedBy"`
	} `json:"pendingRequests"`
	MaxMemory      int64   `json:"maxMemory"`
	MaxStorage     int64   `json:"maxStorage"`
	MaxCores       int64   `json:"maxCores"`
	CoresPerSocket int64   `json:"coresPerSocket"`
	MaxCpu         int64   `json:"maxCpu"`
	HourlyCost     float64 `json:"hourlyCost"`
	HourlyPrice    float64 `json:"hourlyPrice"`
	InstancePrice  struct {
		Price    float64 `json:"price"`
		Cost     float64 `json:"cost"`
		Currency string  `json:"currency"`
		Unit     string  `json:"unit"`
	} `json:"instancePrice"`
	DateCreated   string `json:"dateCreated"`
	LastUpdated   string `json:"lastUpdated"`
	HostName      string `json:"hostName"`
	DomainName    string `json:"domainName"`
	NetworkDomain struct {
		Id int64 `json:"id"`
	} `json:"networkDomain"`
	EnvironmentPrefix   string `json:"environmentPrefix"`
	FirewallEnabled     bool   `json:"firewallEnabled"`
	NetworkLevel        string `json:"networkLevel"`
	AutoScale           bool   `json:"autoScale"`
	InstanceContext     string `json:"instanceContext"`
	Locked              bool   `json:"locked"`
	IsScalable          bool   `json:"isScalable"`
	ExpireCount         int64  `json:"expireCount"`
	ExpireDate          string `json:"expireDate"`
	ExpireWarningDate   string `json:"expireWarningDate"`
	ExpireWarningSent   bool   `json:"expireWarningSent"`
	ShutdownDays        int64  `json:"shutdownDays"`
	ShutdownRenewDays   int64  `json:"shutdownRenewDays"`
	ShutdownCount       int64  `json:"shutdownCount"`
	ShutdownDate        string `json:"shutdownDate"`
	ShutdownWarningDate string `json:"shutdownWarningDate"`
	ShutdownWarningSent bool   `json:"shutdownWarningSent"`
	CreatedBy           struct {
		Id       int64  `json:"id"`
		Username string `json:"username"`
	} `json:"createdBy"`
	Notes string `json:"notes"`
	Stats struct {
		UsedStorage  int64   `json:"usedStorage"`
		MaxStorage   int64   `json:"maxStorage"`
		UsedMemory   int64   `json:"usedMemory"`
		MaxMemory    int64   `json:"maxMemory"`
		UsedCpu      float64 `json:"usedCpu"`
		CpuUsage     float64 `json:"cpuUsage"`
		CpuUsagePeak float64 `json:"cpuUsagePeak"`
		CpuUsageAvg  float64 `json:"cpuUsageAvg"`
	} `json:"stats"`
	IsBusy bool `json:"isBusy"`
	Apps   []struct {
		Id   int64  `json:"id"`
		Name string `json:"name"`
		Tier string `json:"tier"`
	} `json:"apps"`
}

// Only used as an optional struct until all possible options are validated
type InstanceConfig struct {
	CreateUser       bool     `json:"createUser"`
	IsEC2            bool     `json:"isEC2"`
	AllowExisting    bool     `json:"allowExisting"`
	IsVpcSelectable  bool     `json:"isVpcSelectable"`
	Username         string   `json:"username"`
	Host             string   `json:"host"`
	Port             int64    `json:"port"`
	Password         string   `json:"password"`
	PasswordHash     string   `json:"passwordHash"`
	NoAgent          bool     `json:"noAgent"`
	RootPassword     string   `json:"rootPassword"`
	RootPasswordHash string   `json:"rootPasswordHash"`
	ServicePassword  string   `json:"servicePassword"`
	ServiceUsername  string   `json:"serviceUsername"`
	SmbiosAssetTag   string   `json:"smbiosAssetTag"`
	Cloud            string   `json:"cloud"`
	LoadBalancerId   int64    `json:"loadBalancerId"`
	Expose           []string `json:"expose"`
	RemovalOptions   struct {
		Force           bool  `json:"force"`
		KeepBackups     bool  `json:"keepBackups"`
		ReleaseEIPs     bool  `json:"releaseEIPs"`
		RemoveVolumes   bool  `json:"removeVolumes"`
		RemoveResources bool  `json:"removeResources"`
		UserId          int64 `json:"userId"`
	} `json:"removalOptions"`
	Attributes     map[string]interface{} `json:"attributes"`
	VmwareFolderId string                 `json:"vmwareFolderId"`
	SecurityId     string                 `json:"securityId"`
	KmsKeyId       string                 `json:"kmsKeyId"`
	ResourcePoolId string                 `json:"resourcePoolId"`
	PublicIpType   string                 `json:"publicIpType"`
	CatalogItem    int64                  `json:"catalogItem"`
	CreateBackup   bool                   `json:"createBackup"`
	MemoryDisplay  string                 `json:"memoryDisplay"`
	SecurityGroups []struct {
		ID string `json:"id"`
	}
	Backup struct {
		CreateBackup               bool   `json:"createBackup"`
		BackupRepository           int64  `json:"backupRepository"`
		JobAction                  string `json:"jobAction"`
		JobRetentionCount          string `json:"jobRetentionCount"`
		ProviderBackupType         int64  `json:"providerBackupType"`
		Enabled                    bool   `json:"enabled"`
		ShowScheduledBackupWarning bool   `json:"showScheduledBackupWarning"`
	} `json:"backup"`
	ReplicationGroup struct {
		ProviderMethod string `json:"providerMethod"`
	} `json:"replicationGroup"`
	ShutdownDays        string `json:"shutdownDays"`
	ExpireDays          string `json:"expireDays"`
	LayoutSize          int64  `json:"layoutSize"`
	ServicePasswordHash string `json:"servicePasswordHash"`
}

type StorageControler struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Type struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
		Code string `json:"code"`
	}
	MaxDevices         int64 `json:"maxDevices"`
	ReservedUnitNumber int64 `json:"reservedUnitNumber"`
}

type NetworkInterface struct {
	ID              string `json:"id"`
	Row             int64  `json:"row"`
	InternalId      string `json:"internalId"`
	ExternalId      string `json:"externalId"`
	UniqueId        string `json:"uniqueId"`
	PublicIpAddress string `json:"publicIpAddress"`
	Network         struct {
		ID         int64  `json:"id"`
		Subnet     int64  `json:"subnet"`
		Group      int64  `json:"group"`
		Name       string `json:"name"`
		DhcpServer bool   `json:"dhcpServer"`
		Pool       struct {
			ID   int64  `json:"id"`
			Name string `json:"name"`
		} `json:"pool"`
	} `json:"network"`
	IpAddress              string      `json:"ipAddress"`
	Ipv6Address            string      `json:"ipv6Address"`
	Description            string      `json:"description"`
	IpMode                 string      `json:"ipMode"`
	NetworkInterfaceTypeId interface{} `json:"networkInterfaceTypeId"`
	IsPrimary              bool        `json:"isPrimary"`
	Type                   struct {
		ID   int64  `json:"id"`
		Code string `json:"code"`
		Name string `json:"name"`
	} `json:"type"`
	PoolAssigned bool   `json:"poolAssigned"`
	MacAddress   string `json:"macAddress"`
	Active       bool   `json:"active"`
	DHCP         bool   `json:"dhcp"`
}

type InstancePlan struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"`
}

type Owner struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
}

type ContainerDetails struct {
	ID               int64  `json:"id"`
	UUID             string `json:"uuid"`
	Name             string `json:"name"`
	IP               string `json:"ip"`
	InternalIp       string `json:"internalIp"`
	InternalHostname string `json:"internalHostname"`
	ExternalHostname string `json:"externalHostname"`
	ExternalDomain   string `json:"externalDomain"`
	ExternalFqdn     string `json:"externalFqdn"`
	AccountId        int64  `json:"accountId"`
	Instance         struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	} `json:"instance"`
	ContainerType struct {
		ID       int64  `json:"id"`
		Name     string `json:"name"`
		Code     string `json:"code"`
		Category string `json:"category"`
	} `json:"containerType"`
	Server struct {
		ID               int64              `json:"id"`
		UUID             string             `json:"uuid"`
		ExternalId       string             `json:"externalId"`
		InternalId       string             `json:"internalId"`
		ExternalUniqueId string             `json:"externalUniqueId"`
		Name             string             `json:"name"`
		ExternalName     string             `json:"externalName"`
		Hostname         string             `json:"hostname"`
		AccountId        int64              `json:"accountId"`
		SshHost          string             `json:"sshHost"`
		ExternalIp       string             `json:"externalIp"`
		InternalIp       string             `json:"internalIp"`
		Platform         string             `json:"platform"`
		PlatformVersion  string             `json:"platformVersion"`
		AgentInstalled   bool               `json:"agentInstalled"`
		AgentVersion     string             `json:"agentVersion"`
		Interfaces       []NetworkInterface `json:"interfaces"`
		SourceImage      struct {
			ID   int64  `json:"id"`
			Name string `json:"name"`
			Code string `json:"code"`
		} `json:"sourceImage"`
	} `json:"server"`
}

type Snapshot struct {
	ID              int64  `json:"id"`
	Name            string `json:"name"`
	Description     string `json:"description"`
	ExternalId      string `json:"externalId"`
	Status          string `json:"status"`
	State           string `json:"state"`
	SnapshotType    string `json:"snapshotType"`
	SnapshotCreated string `json:"snapshotCreated"`
	Zone            struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	} `json:"zone"`
	Datastore       string `json:"datastore"`
	ParentSnapshot  string `json:"parentSnapshot"`
	CurrentlyActive bool   `json:"currentlyActive"`
	DateCreated     string `json:"dateCreated"`
}

type InstanceSchedule struct {
	ID             int64          `json:"id"`
	ScheduleType   string         `json:"scheduleType"`
	StartDayOfWeek int64          `json:"startDayOfWeek"`
	StartTime      string         `json:"startTime"`
	EndDayOfWeek   int64          `json:"endDayOfWeek"`
	EndTime        string         `json:"endTime"`
	StartDisplay   string         `json:"startDisplay"`
	EndDisplay     string         `json:"endDisplay"`
	Threshold      ScaleThreshold `json:"threshold"`
	DateCreated    string         `json:"dateCreated"`
	LastUpdated    string         `json:"lastUpdated"`
}

// ListInstancesResult structure parses the list instances response payload
type ListInstancesResult struct {
	Instances *[]Instance `json:"instances"`
	Meta      *MetaResult `json:"meta"`
}

type GetInstanceResult struct {
	Instance *Instance         `json:"instance"`
	Success  bool              `json:"success"`
	Message  string            `json:"msg"`
	Errors   map[string]string `json:"errors"`
}

type CreateInstanceResult struct {
	Success  bool              `json:"success"`
	Message  string            `json:"msg"`
	Errors   map[string]string `json:"errors"`
	Instance *Instance         `json:"instance"`
}

type UpdateInstanceResult struct {
	CreateInstanceResult
}

type DeleteInstanceResult struct {
	DeleteResult
}

type ListInstancePlansResult struct {
	Plans *[]InstancePlan `json:"plans"`
	Meta  *MetaResult     `json:"meta"`
}

type GetInstancePlanResult struct {
	Plan *InstancePlan `json:"plan"`
}

type GetInstanceSecurityGroupsResult struct {
	SecurityGroups *[]SecurityGroup  `json:"securityGroups"`
	Success        bool              `json:"success"`
	Message        string            `json:"msg"`
	Errors         map[string]string `json:"errors"`
}

type UpdateInstanceSecurityGroupsResult struct {
	SecurityGroups *[]SecurityGroup  `json:"securityGroups"`
	Success        bool              `json:"success"`
	Message        string            `json:"msg"`
	Errors         map[string]string `json:"errors"`
}

type GetInstanceSnapshotResult struct {
	Snapshot *Snapshot         `json:"snapshot"`
	Success  bool              `json:"success"`
	Message  string            `json:"msg"`
	Errors   map[string]string `json:"errors"`
}

type RunWorkflowOnInstanceResult struct {
	ProcessId int64             `json:"processId"`
	Success   bool              `json:"success"`
	Message   string            `json:"msg"`
	Errors    map[string]string `json:"errors"`
}

type ApplyStateResult struct {
	ExecutionId string            `json:"executionId"`
	Success     bool              `json:"success"`
	Message     string            `json:"msg"`
	Errors      map[string]string `json:"errors"`
}

type ValidateInstanceStateApplyResult struct {
	ExecutionId string            `json:"executionId"`
	Success     bool              `json:"success"`
	Message     string            `json:"msg"`
	Errors      map[string]string `json:"errors"`
}

type UpdateInstanceScheduleResult struct {
	InstanceSchedule InstanceSchedule  `json:"instanceSchedule"`
	Success          bool              `json:"success"`
	Message          string            `json:"msg"`
	Errors           map[string]string `json:"errors"`
}

type GetInstanceScheduleResult struct {
	InstanceSchedule InstanceSchedule  `json:"instanceSchedule"`
	Success          bool              `json:"success"`
	Message          string            `json:"msg"`
	Errors           map[string]string `json:"errors"`
}

type ListInstanceScheduleResult struct {
	InstanceSchedules []InstanceSchedule `json:"instanceSchedules"`
	Success           bool               `json:"success"`
	Message           string             `json:"msg"`
	Errors            map[string]string  `json:"errors"`
}

// API endpoints
func (client *Client) ListInstances(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        InstancesPath,
		QueryParams: req.QueryParams,
		Result:      &ListInstancesResult{},
	})
}

func (client *Client) GetInstance(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        fmt.Sprintf("%s/%d", InstancesPath, id),
		QueryParams: req.QueryParams,
		Result:      &GetInstanceResult{},
	})
}

func (client *Client) CreateInstance(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "POST",
		Path:        InstancesPath,
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &CreateInstanceResult{},
	})
}

func (client *Client) UpdateInstance(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "PUT",
		Path:        fmt.Sprintf("%s/%d", InstancesPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &UpdateInstanceResult{},
	})
}

func (client *Client) DeleteInstance(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "DELETE",
		Path:        fmt.Sprintf("%s/%d", InstancesPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &DeleteInstanceResult{},
	})
}

func (client *Client) StartInstance(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "PUT",
		Path:        fmt.Sprintf("%s/%d/start", InstancesPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &UpdateInstanceResult{},
	})
}

func (client *Client) StopInstance(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "PUT",
		Path:        fmt.Sprintf("%s/%d/stop", InstancesPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &StandardResult{},
	})
}

func (client *Client) RestartInstance(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "PUT",
		Path:        fmt.Sprintf("%s/%d/restart", InstancesPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &StandardResult{},
	})
}

func (client *Client) SuspendInstance(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "PUT",
		Path:        fmt.Sprintf("%s/%d/suspend", InstancesPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &StandardResult{},
	})
}

func (client *Client) SnapshotInstance(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "PUT",
		Path:        fmt.Sprintf("%s/%d/snapshot", InstancesPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &StandardResult{},
	})
}

func (client *Client) BackupInstance(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "PUT",
		Path:        fmt.Sprintf("%s/%d/backup", InstancesPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &StandardResult{},
	})
}

func (client *Client) CancelInstanceExpiration(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "PUT",
		Path:        fmt.Sprintf("%s/%d/cancel-expiration", InstancesPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &StandardResult{},
	})
}

func (client *Client) CancelInstanceRemoval(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "PUT",
		Path:        fmt.Sprintf("%s/%d/cancel-removal", InstancesPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &StandardResult{},
	})
}

func (client *Client) CancelInstanceShutdown(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "PUT",
		Path:        fmt.Sprintf("%s/%d/cancel-shutdown", InstancesPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &StandardResult{},
	})
}

func (client *Client) EjectInstance(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "PUT",
		Path:        fmt.Sprintf("%s/%d/eject", InstancesPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &StandardResult{},
	})
}

func (client *Client) CloneInstanceToImage(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "PUT",
		Path:        fmt.Sprintf("%s/%d/clone-image", InstancesPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &StandardResult{},
	})
}

func (client *Client) CloneInstance(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "PUT",
		Path:        fmt.Sprintf("%s/%d/clone", InstancesPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &StandardResult{},
	})
}

func (client *Client) ImportInstanceSnapshot(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "PUT",
		Path:        fmt.Sprintf("%s/%d/import-snapshot", InstancesPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &StandardResult{},
	})
}

func (client *Client) CreateLinkedClone(id int64, snapshotId int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "PUT",
		Path:        fmt.Sprintf("%s/%d/linked-clone/%d", InstancesPath, id, snapshotId),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &StandardResult{},
	})
}

func (client *Client) ApplyInstanceState(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "POST",
		Path:        fmt.Sprintf("%s/%d/apply", InstancesPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &ApplyStateResult{},
	})
}

func (client *Client) RefreshInstanceState(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "POST",
		Path:        fmt.Sprintf("%s/%d/refresh", InstancesPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &StandardResult{},
	})
}

func (client *Client) RevertInstanceToSnapshot(id int64, snapshotId int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "PUT",
		Path:        fmt.Sprintf("%s/%d/revert-snapshot/%d", InstancesPath, id, snapshotId),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &StandardResult{},
	})
}

func (client *Client) DeleteInstanceSnapshot(snapshotId int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "DELETE",
		Path:        fmt.Sprintf("api/snapshots/%d", snapshotId),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &StandardResult{},
	})
}

func (client *Client) ResizeInstance(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "PUT",
		Path:        fmt.Sprintf("%s/%d/resize", InstancesPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &UpdateInstanceResult{},
	})
}

func (client *Client) LockInstance(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "PUT",
		Path:        fmt.Sprintf("%s/%d/lock", InstancesPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &UpdateInstanceResult{},
	})
}

func (client *Client) UnlockInstance(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "PUT",
		Path:        fmt.Sprintf("%s/%d/unlock", InstancesPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &UpdateInstanceResult{},
	})
}

func (client *Client) GetInstanceSecurityGroups(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        fmt.Sprintf("%s/%d/security-groups", InstancesPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &GetInstanceSecurityGroupsResult{},
	})
}

func (client *Client) UpdateInstanceSecurityGroups(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "POST",
		Path:        fmt.Sprintf("%s/%d/security-groups", InstancesPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &UpdateInstanceSecurityGroupsResult{},
	})
}

// helper functions

func (client *Client) FindInstanceByName(name string) (*Response, error) {
	// Find by name, then get by ID
	resp, err := client.ListInstances(&Request{
		QueryParams: map[string]string{
			"name": name,
		},
	})
	if err != nil {
		return resp, err
	}
	listResult := resp.Result.(*ListInstancesResult)
	instanceCount := len(*listResult.Instances)
	if instanceCount != 1 {
		return resp, fmt.Errorf("found %d Instances for %v", instanceCount, name)
	}
	firstRecord := (*listResult.Instances)[0]
	instanceId := firstRecord.ID
	return client.GetInstance(instanceId, &Request{})
}

// Plan fetching
// todo: this needs to be refactored soon

func (client *Client) ListInstancePlans(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        fmt.Sprintf("%s/service-plans", InstancesPath),
		QueryParams: req.QueryParams,
		Result:      &ListInstancePlansResult{},
	})
}

// todo: need this api endpoint still, and consolidate to /api/plans perhaps
func (client *Client) GetInstancePlan(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        fmt.Sprintf("%s/service-plans/%d", InstancesPath, id),
		QueryParams: req.QueryParams,
		Result:      &GetInstancePlanResult{},
	})
}

func (client *Client) FindInstancePlanByName(name string, req *Request) (*Response, error) {
	// Find by name, then get by ID
	resp, err := client.ListInstancePlans(&Request{
		QueryParams: map[string]string{
			//"name": name, // this is not even supported..
			"zoneId":   req.QueryParams["zoneId"], // todo: use cloudId
			"layoutId": req.QueryParams["layoutId"],
			"siteId":   req.QueryParams["siteId"], // todo: use groupId
		},
	})
	if err != nil {
		return resp, err
	}
	listResult := resp.Result.(*ListInstancePlansResult)
	planCount := len(*listResult.Plans)
	// need to filter the list ourselves for now..
	var matchingPlans []InstancePlan
	for i := 0; i < planCount; i++ {
		plan := (*listResult.Plans)[i] // .(InstancePlan)
		if plan.Name == name || plan.Code == name || string(rune(plan.ID)) == name {
			matchingPlans = append(matchingPlans, plan)
		}
	}
	matchingPlanCount := len(matchingPlans)
	if matchingPlanCount != 1 {
		return resp, fmt.Errorf("found %d Plans for '%v'", matchingPlanCount, name)
	}
	firstRecord := matchingPlans[0]

	// planId := firstRecord.ID
	// return client.GetInstancePlan(planId, &Request{})

	// for now just return a fake response until endpoint is ready
	var result = &GetInstancePlanResult{
		Plan: &firstRecord,
	}
	mockResp := &Response{
		//RestyResponse: restyResponse,
		Success:    true,
		StatusCode: 200,
		Status:     "200 OK",
		ReceivedAt: time.Now(),
		// Size: restyResponse.Size(),
		// Body: restyResponse.Body(), // byte[]
		Result: result,
	}
	return mockResp, nil
}

// this should work by code or name
// it also requires zoneId AND layoutId??
func (client *Client) FindInstancePlanByCode(code string, req *Request) (*Response, error) {
	return client.FindInstancePlanByName(code, req)
}

// ListInstanceSchedules fetches existing instance scaling schedules
func (client *Client) ListInstanceSchedules(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        fmt.Sprintf("%s/%d/schedules", InstancesPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &ListInstanceScheduleResult{},
	})
}

// CreateInstanceSchedule creates a new instance scaling schedule
func (client *Client) CreateInstanceSchedule(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "PUT",
		Path:        fmt.Sprintf("%s/%d/schedules", InstancesPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &UpdateInstanceScheduleResult{},
	})
}

// GetInstanceSchedule fetches an existing instance scaling schedule
func (client *Client) GetInstanceSchedule(id int64, scheduleId int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        fmt.Sprintf("%s/%d/schedules/%d", InstancesPath, id, scheduleId),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &GetInstanceScheduleResult{},
	})
}

// UpdateInstanceSchedule updates an existing instance scaling schedule
func (client *Client) UpdateInstanceSchedule(id int64, scheduleId int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "PUT",
		Path:        fmt.Sprintf("%s/%d/schedules/%d", InstancesPath, id, scheduleId),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &UpdateInstanceScheduleResult{},
	})
}

// DeleteInstanceSchedule deletes a specified instance scaling schedule
func (client *Client) DeleteInstanceSchedule(id int64, scheduleId int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "DELETE",
		Path:        fmt.Sprintf("%s/%d/schedules/%d", InstancesPath, id, scheduleId),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &StandardResult{},
	})
}

// ValidateInstanceStateApply validates instance configuration and templateParameter variables before executing the apply
func (client *Client) ValidateInstanceStateApply(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "POST",
		Path:        fmt.Sprintf("%s/%d/validate-apply", InstancesPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &ValidateInstanceStateApplyResult{},
	})
}

// RunWorkflowOnInstance executes a workflow on an existing instance
func (client *Client) RunWorkflowOnInstance(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "PUT",
		Path:        fmt.Sprintf("%s/%d/workflow", InstancesPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &RunWorkflowOnInstanceResult{},
	})
}

// RemoveInstanceFromControl removes an instance from Morpheus control
func (client *Client) RemoveInstanceFromControl(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "POST",
		Path:        fmt.Sprintf("%s/remove-from-control", InstancesPath),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &StandardResult{},
	})
}

// GetInstanceSnapshot fetches an existing instance snapshot
func (client *Client) GetInstanceSnapshot(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        fmt.Sprintf("api/snapshots/%d", id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &StandardResult{},
	})
}
