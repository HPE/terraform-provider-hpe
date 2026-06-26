package morpheus

import (
	"fmt"
	"time"
)

var (
	// ClustersPath is the API endpoint for clusters
	ClustersPath = "/api/clusters"
)

// Cluster structures for use in request and response payloads
type Cluster struct {
	ID                  int64    `json:"id"`
	Name                string   `json:"name"`
	Code                string   `json:"code"`
	Category            string   `json:"category"`
	Visibility          string   `json:"visibility"`
	Description         string   `json:"description"`
	Location            string   `json:"location"`
	Enabled             bool     `json:"enabled"`
	ServiceUrl          string   `json:"serviceUrl"`
	ServiceHost         string   `json:"serviceHost"`
	ServicePath         string   `json:"servicePath"`
	ServiceHostname     string   `json:"serviceHostname"`
	ServicePort         int64    `json:"servicePort"`
	ServiceUsername     string   `json:"serviceUsername"`
	ServicePassword     string   `json:"servicePassword"`
	ServicePasswordHash string   `json:"servicePasswordHash"`
	ServiceToken        string   `json:"serviceToken"`
	ServiceTokenHash    string   `json:"serviceTokenHash"`
	ServiceAccess       string   `json:"serviceAccess"`
	ServiceAccessHash   string   `json:"serviceAccessHash"`
	ServiceCert         string   `json:"serviceCert"`
	ServiceCertHash     string   `json:"serviceCertHash"`
	ServiceVersion      string   `json:"serviceVersion"`
	SearchDomains       string   `json:"searchDomains"`
	EnableInternalDns   bool     `json:"enableInternalDns"`
	InternalId          string   `json:"internalId"`
	ExternalId          string   `json:"externalId"`
	DatacenterId        string   `json:"datacenterId"`
	StatusMessage       string   `json:"statusMessage"`
	InventoryLevel      string   `json:"inventoryLevel"`
	LastSyncDuration    int64    `json:"lastSyncDuration"`
	Labels              []string `json:"labels"`
	Type                struct {
		Id   int64  `json:"id"`
		Name string `json:"name"`
	} `json:"type"`
	Layout struct {
		Id                int64  `json:"id"`
		Name              string `json:"name"`
		ProvisionTypeCode string `json:"provisionTypeCode"`
	} `json:"layout"`
	Group map[string]interface{} `json:"group"`
	Site  struct {
		Id   int64  `json:"id"`
		Name string `json:"name"`
	} `json:"site"`
	Zone struct {
		Id       int64  `json:"id"`
		Name     string `json:"name"`
		ZoneType struct {
			Id int64 `json:"id"`
		} `json:"zoneType"`
	} `json:"zone"`
	Servers      []Server `json:"servers"`
	Status       string   `json:"status"`
	Managed      bool     `json:"managed"`
	ServiceEntry string   `json:"serviceEntry"`
	CreatedBy    struct {
		Id       int64  `json:"id"`
		Username string `json:"username"`
	} `json:"createdBy"`
	UserGroup string `json:"userGroup"`
	Owner     struct {
		Id   int64  `json:"id"`
		Name string `json:"name"`
	} `json:"owner"`
	WorkerStats struct {
		UsedStorage  int64   `json:"usedStorage"`
		MaxStorage   int64   `json:"maxStorage"`
		UsedMemory   int64   `json:"usedMemory"`
		MaxMemory    int64   `json:"maxMemory"`
		UsedCpu      float64 `json:"usedCpu"`
		CpuUsage     float64 `json:"cpuUsage"`
		CpuUsagePeak float64 `json:"cpuUsagePeak"`
		CpuUsageAvg  float64 `json:"cpuUsageAvg"`
	}
	ContainersCount  int64                  `json:"containersCount"`
	DeploymentsCount int64                  `json:"deploymentsCount"`
	PodsCount        int64                  `json:"podsCount"`
	JobsCount        int64                  `json:"jobsCount"`
	VolumesCount     int64                  `json:"volumesCount"`
	NamespacesCount  int64                  `json:"namespacesCount"`
	WorkersCount     int64                  `json:"workersCount"`
	ServicesCount    int64                  `json:"servicesCount"`
	Config           map[string]interface{} `json:"config"`
}

type Server struct {
	Id      int64  `json:"id"`
	Name    string `json:"name"`
	TypeSet struct {
		Id   int64  `json:"id"`
		Code string `json:"code"`
		Name string `json:"name"`
	} `json:"typeSet"`
	ComputeServerType struct {
		Id       int64  `json:"id"`
		Code     string `json:"code"`
		NodeType string `json:"nodeType"`
	} `json:"computeServerType"`
}

// ListClustersResult structure parses the list clusters response payload
type ListClustersResult struct {
	Clusters *[]Cluster  `json:"clusters"`
	Meta     *MetaResult `json:"meta"`
}

type GetClusterResult struct {
	Cluster *Cluster `json:"cluster"`
}

type GetClusterApiConfigResult struct {
	ServiceUrl          string `json:"serviceUrl"`
	ServiceHost         string `json:"serviceHost"`
	ServicePath         string `json:"servicePath"`
	ServiceHostname     string `json:"serviceHostname"`
	ServicePort         int64  `json:"servicePort"`
	ServiceUsername     string `json:"serviceUsername"`
	ServicePassword     string `json:"servicePassword"`
	ServicePasswordHash string `json:"servicePasswordHash"`
	ServiceToken        string `json:"serviceToken"`
	ServiceAccess       string `json:"serviceAccess"`
	ServiceCert         string `json:"serviceCert"`
	ServiceVersion      string `json:"serviceVersion"`
}

type ListClusterNamespacesResults struct {
	Namespaces []Namespaces `json:"namespaces"`
	Meta       *MetaResult  `json:"meta"`
}

type Namespaces struct {
	Id                 int64  `json:"id"`
	Name               string `json:"name"`
	Description        string `json:"description"`
	RegionCode         string `json:"regionCode"`
	ExternalId         string `json:"externalId"`
	Status             string `json:"status"`
	Visibility         string `json:"visibility"`
	Active             bool   `json:"active"`
	ResourcePermission struct {
		AllGroups            bool   `json:"allGroups"`
		DefaultStore         bool   `json:"defaultStore"`
		AllPlans             bool   `json:"allPlans"`
		DefaultTarget        bool   `json:"defaultTarget"`
		MorpheusResourceType string `json:"morpheusResourceType"`
		MorpheusResourceId   int64  `json:"morpheusResourceId"`
		CanManage            bool   `json:"canManage"`
		All                  bool   `json:"all"`
		Account              struct {
			ID int64 `json:"id"`
		} `json:"account"`
		Sites []struct {
			ID      int64  `json:"id"`
			Name    string `json:"name"`
			Default bool   `json:"default"`
		} `json:"sites"`
		Plans []struct {
			ID      int64  `json:"id"`
			Name    string `json:"name"`
			Default bool   `json:"default"`
		} `json:"plans"`
	} `json:"resourcePermission"`
}

type ClusterWorker struct {
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
		Type string `json:"type"`
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
	LvmEnabled      bool        `json:"lvmEnabled"`
	ApiKey          string      `json:"apiKey"`
	SoftwareRaid    bool        `json:"softwareRaid"`
	DateCreated     time.Time   `json:"dateCreated"`
	LastUpdated     time.Time   `json:"lastUpdated"`
	Stats           struct {
		UsedStorage     int64   `json:"usedStorage"`
		ReservedStorage int64   `json:"reservedStorage"`
		MaxStorage      int64   `json:"maxStorage"`
		UsedMemory      int64   `json:"usedMemory"`
		ReservedMemory  int64   `json:"reservedMemory"`
		MaxMemory       int64   `json:"maxMemory"`
		CpuUsage        float64 `json:"cpuUsage"`
	} `json:"stats"`
	Status                 string      `json:"status"`
	StatusMessage          string      `json:"statusMessage"`
	ErrorMessage           string      `json:"errorMessage"`
	StatusDate             string      `json:"statusDate"`
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
	MaxCpu                 interface{} `json:"maxCpu"`
	ManageInternalFirewall bool        `json:"manageInternalFirewall"`
	EnableLogs             bool        `json:"enableLogs"`
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
	Volumes     []StorageVolume    `json:"volumes"`
	Controllers []StorageControler `json:"controllers"`
	Interfaces  []struct {
		ID                int         `json:"id"`
		Reftype           interface{} `json:"refType"`
		Refid             interface{} `json:"refId"`
		Name              string      `json:"name"`
		Internalid        string      `json:"internalId"`
		Externalid        string      `json:"externalId"`
		Uniqueid          interface{} `json:"uniqueId"`
		Publicipaddress   string      `json:"publicIpAddress"`
		Publicipv6Address interface{} `json:"publicIpv6Address"`
		Ipaddress         string      `json:"ipAddress"`
		Ipv6Address       interface{} `json:"ipv6Address"`
		Ipsubnet          interface{} `json:"ipSubnet"`
		Ipv6Subnet        interface{} `json:"ipv6Subnet"`
		Description       interface{} `json:"description"`
		Dhcp              bool        `json:"dhcp"`
		Active            bool        `json:"active"`
		Poolassigned      bool        `json:"poolAssigned"`
		Primaryinterface  bool        `json:"primaryInterface"`
		Network           struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"network"`
		Subnet          interface{} `json:"subnet"`
		Networkgroup    interface{} `json:"networkGroup"`
		Networkposition interface{} `json:"networkPosition"`
		Networkpool     interface{} `json:"networkPool"`
		Networkdomain   interface{} `json:"networkDomain"`
		Type            struct {
			ID   int    `json:"id"`
			Code string `json:"code"`
			Name string `json:"name"`
		} `json:"type"`
		Ipmode     string `json:"ipMode"`
		Macaddress string `json:"macAddress"`
	} `json:"interfaces"`
	Labels                   []interface{} `json:"labels"`
	Tags                     []interface{} `json:"tags"`
	Enabled                  bool          `json:"enabled"`
	TagCompliant             interface{}   `json:"tagCompliant"`
	Containers               []int64       `json:"containers"`
	GuestConsolePreferred    bool          `json:"guestConsolePreferred"`
	GuestConsoleType         string        `json:"guestConsoleType"`
	GuestConsoleUsername     string        `json:"guestConsoleUsername"`
	GuestConsolePassword     string        `json:"guestConsolePassword"`
	GuestConsolePasswordHash string        `json:"guestConsolePasswordHash"`
	GuestConsolePort         interface{}   `json:"guestConsolePort"`
}

type ClusterStorageVolume struct {
	ID                int64  `json:"id"`
	DisplayOrder      int64  `json:"displayOrder"`
	Active            bool   `json:"active"`
	UsedStorage       int64  `json:"usedStorage"`
	Resizeable        bool   `json:"resizeable"`
	Online            bool   `json:"online"`
	DeviceDisplayName string `json:"deviceDisplayName"`
	RefType           string `json:"refType"`
	Name              string `json:"name"`
	ClaimName         string `json:"claimName"`
	VolumeType        string `json:"volumeType"`
	DeviceName        string `json:"deviceName"`
	Removable         bool   `json:"removable"`
	PoolName          string `json:"poolName"`
	ReadOnly          bool   `json:"readOnly"`
	ZoneId            int64  `json:"zoneId"`
	RootVolume        bool   `json:"rootVolume"`
	RefId             int64  `json:"refId"`
	Category          string `json:"category"`
	Status            string `json:"status"`
	MaxStorage        int64  `json:"maxStorage"`
	Account           struct {
		ID int64 `json:"id"`
	} `json:"account"`
	Type struct {
		ID int64 `json:"id"`
	} `json:"type"`
}

type ClusterContainer struct {
	ID            int64       `json:"id"`
	UUID          string      `json:"uuid"`
	AccountId     int64       `json:"accountId"`
	Instance      interface{} `json:"instance"`
	ContainerType struct {
		ID       int64  `json:"id"`
		Code     string `json:"code"`
		Category string `json:"category"`
		Name     string `json:"name"`
	} `json:"containerType"`
	ContainerTypeSet struct {
		ID       int64  `json:"id"`
		Code     string `json:"code"`
		Category string `json:"category"`
	} `json:"containerTypeSet"`
	Server struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	} `json:"server"`
	Cloud struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
		Type string `json:"type"`
	} `json:"cloud"`
	Name             string        `json:"name"`
	IP               string        `json:"ip"`
	InternalIp       string        `json:"internalIp"`
	InternalHostname string        `json:"internalHostname"`
	ExternalHostname string        `json:"externalHostname"`
	ExternalDomain   string        `json:"externalDomain"`
	ExternalFqdn     string        `json:"externalFqdn"`
	Ports            []interface{} `json:"ports"`
	Plan             struct {
		ID   int64  `json:"id"`
		Code string `json:"code"`
		Name string `json:"name"`
	} `json:"plan"`
	DateCreated       time.Time   `json:"dateCreated"`
	LastUpdated       time.Time   `json:"lastUpdated"`
	StatsEnabled      bool        `json:"statsEnabled"`
	Status            string      `json:"status"`
	UserStatus        interface{} `json:"userStatus"`
	EnvironmentPrefix interface{} `json:"environmentPrefix"`
	ConfigGroup       interface{} `json:"configGroup"`
	ConfigId          interface{} `json:"configId"`
	ConfigRole        interface{} `json:"configRole"`
	Stats             struct {
		Ts             time.Time `json:"ts"`
		Running        bool      `json:"running"`
		UserCpuUsage   float64   `json:"userCpuUsage"`
		SystemCpuUsage float64   `json:"systemCpuUsage"`
		UsedMemory     int64     `json:"usedMemory"`
		MaxMemory      int64     `json:"maxMemory"`
		CacheMemory    int64     `json:"cacheMemory"`
		MaxStorage     int64     `json:"maxStorage"`
		UsedStorage    int64     `json:"usedStorage"`
		ReadIOPS       int64     `json:"readIOPS"`
		WriteIOPS      int64     `json:"writeIOPS"`
		TotalIOPS      int64     `json:"totalIOPS"`
		NetTxUsage     int64     `json:"netTxUsage"`
		NetRxUsage     int64     `json:"netRxUsage"`
	} `json:"stats"`
	RuntimeInfo struct {
	} `json:"runtimeInfo"`
	ContainerVersion string      `json:"containerVersion"`
	RepositoryImage  string      `json:"repositoryImage"`
	PlanCategory     interface{} `json:"planCategory"`
	Hostname         string      `json:"hostname"`
	DomainName       string      `json:"domainName"`
	VolumeCreated    bool        `json:"volumeCreated"`
	ContainerCreated bool        `json:"containerCreated"`
	MaxStorage       int64       `json:"maxStorage"`
	MaxMemory        int64       `json:"maxMemory"`
	MaxCores         int64       `json:"maxCores"`
	MaxCpu           interface{} `json:"maxCpu"`
	HourlyPrice      float64     `json:"hourlyPrice"`
	AvailableActions []struct {
		Code string `json:"code"`
		Name string `json:"name"`
	} `json:"availableActions"`
}

type ListClusterContainersResults struct {
	Containers *[]ClusterContainer `json:"containers"`
	Meta       *MetaResult         `json:"meta"`
}

type ListClusterUpgradeVersionsResults struct {
	Versions       *[]string   `json:"versions"`
	CurrentVersion string      `json:"currentVersion"`
	Meta           *MetaResult `json:"meta"`
}

type ListClusterVolumesResults struct {
	Volumes *[]ClusterStorageVolume `json:"volumes"`
	Meta    *MetaResult             `json:"meta"`
}

type ListClusterWorkersResults struct {
	Workers *[]ClusterWorker `json:"workers"`
	Meta    *MetaResult      `json:"meta"`
}

type ApplyTemplateToClusterResult struct {
	ExecutionId string            `json:"executionId"`
	Success     bool              `json:"success"`
	Message     string            `json:"msg"`
	Errors      map[string]string `json:"errors"`
}

type CreateClusterResult struct {
	Success bool              `json:"success"`
	Message string            `json:"msg"`
	Errors  map[string]string `json:"errors"`
	Cluster *Cluster          `json:"cluster"`
}

type AddClusterWorkerResult struct {
	Success bool              `json:"success"`
	Message string            `json:"msg"`
	Errors  map[string]string `json:"errors"`
	Workers *[]ClusterWorker  `json:"workers"`
}

type UpdateClusterResult struct {
	CreateClusterResult
}

type DeleteClusterResult struct {
	DeleteResult
}

type DeleteClusterWorkerResult struct {
	DeleteResult
}

// API endpoints
func (client *Client) ListClusters(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        ClustersPath,
		QueryParams: req.QueryParams,
		Result:      &ListClustersResult{},
	})
}

func (client *Client) ListClusterNamespaces(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        fmt.Sprintf("%s/%d/namespaces", ClustersPath, id),
		QueryParams: req.QueryParams,
		Result:      &ListClusterNamespacesResults{},
	})
}

func (client *Client) GetCluster(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        fmt.Sprintf("%s/%d", ClustersPath, id),
		QueryParams: req.QueryParams,
		Result:      &GetClusterResult{},
	})
}

// CreateCluster creates a new cluster
func (client *Client) CreateCluster(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "POST",
		Path:        ClustersPath,
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &CreateClusterResult{},
	})
}

// UpdateCluster updates an existing cluster
func (client *Client) UpdateCluster(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "PUT",
		Path:        fmt.Sprintf("%s/%d", ClustersPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &UpdateClusterResult{},
	})
}

// DeleteCluster deletes an existing cluster
func (client *Client) DeleteCluster(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "DELETE",
		Path:        fmt.Sprintf("%s/%d", ClustersPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &DeleteClusterResult{},
	})
}

// GetClusterApiConfig gets the api configuration for an existing cluster
func (client *Client) GetClusterApiConfig(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        fmt.Sprintf("%s/%d/api-config", ClustersPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &GetClusterApiConfigResult{},
	})
}

// ApplyTemplateToCluster applies a template to an existing cluster
func (client *Client) ApplyTemplateToCluster(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "POST",
		Path:        fmt.Sprintf("%s/%d/apply-template", ClustersPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &ApplyTemplateToClusterResult{},
	})
}

// Containers

// ListClusterContainers
func (client *Client) ListClusterContainers(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        fmt.Sprintf("%s/%d/containers", ClustersPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &ListClusterContainersResults{},
	})
}

// RestartClusterContainer restarts a container on an existing cluster
func (client *Client) RestartClusterContainer(id int64, containerId int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "PUT",
		Path:        fmt.Sprintf("%s/%d/containers/%d/restart", ClustersPath, id, containerId),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &StandardResult{},
	})
}

// DeleteClusterContainer removes a container from an existing cluster
func (client *Client) DeleteClusterContainer(id int64, containerId int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "DELETE",
		Path:        fmt.Sprintf("%s/%d/containers/%d", ClustersPath, id, containerId),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &StandardResult{},
	})
}

func (client *Client) DeleteClusterWorker(id int64, workerId int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "DELETE",
		Path:        fmt.Sprintf("%s/%d/servers/%d", ClustersPath, id, workerId),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &DeleteClusterWorkerResult{},
	})
}

func (client *Client) DeleteClusterService(id int64, serviceId int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "DELETE",
		Path:        fmt.Sprintf("%s/%d/services/%d", ClustersPath, id, serviceId),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &StandardResult{},
	})
}

func (client *Client) DeleteClusterStatefulSet(id int64, statefulSetId int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "DELETE",
		Path:        fmt.Sprintf("%s/%d/statefulsets/%d", ClustersPath, id, statefulSetId),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &StandardResult{},
	})
}

func (client *Client) RestartClusterStatefulSet(id int64, statefulSetId int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "PUT",
		Path:        fmt.Sprintf("%s/%d/statefulsets/%d/restart", ClustersPath, id, statefulSetId),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &StandardResult{},
	})
}

// ListClusterUpgradeVersions lists all available versions to upgrade a cluster to
func (client *Client) ListClusterUpgradeVersions(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        fmt.Sprintf("%s/%d/upgrade-cluster", ClustersPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &ListClusterUpgradeVersionsResults{},
	})
}

// UpgradeCluster updates the kubectl and kudeadm versions on a Kubernetes cluster to the specified version
func (client *Client) UpgradeCluster(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "POST",
		Path:        fmt.Sprintf("%s/%d/upgrade-cluster", ClustersPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &StandardResult{},
	})
}

// ListClusterWorkers lists all the workers for an existing cluster
func (client *Client) ListClusterVolumes(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        fmt.Sprintf("%s/%d/volumes", ClustersPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &ListClusterVolumesResults{},
	})
}

// DeleteClusterVolume removes a storage volume from an existing cluster
func (client *Client) DeleteClusterVolume(id int64, volumeId int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "DELETE",
		Path:        fmt.Sprintf("%s/%d/volumes/%d", ClustersPath, id, volumeId),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &StandardResult{},
	})
}

// UpdateClusterWorkerCount updates the number of workers for Azure AKS, Google GKE, and Amazon EKS clusters
func (client *Client) UpdateClusterWorkerCount(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "PUT",
		Path:        fmt.Sprintf("%s/%d/worker-count", ClustersPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &StandardResult{},
	})
}

// ListClusterWorkers lists all the workers for an existing cluster
func (client *Client) ListClusterWorkers(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        fmt.Sprintf("%s/%d/workers", ClustersPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &ListClusterWorkersResults{},
	})
}

// AddClusterWorker adds a new cluster worker to an existing cluster
func (client *Client) AddClusterWorker(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "POST",
		Path:        fmt.Sprintf("%s/%d/servers", ClustersPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &AddClusterWorkerResult{},
	})
}

// RefreshCluster triggers a refresh of an existing cluster
func (client *Client) RefreshCluster(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        fmt.Sprintf("%s/%d/refresh", ClustersPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &StandardResult{},
	})
}

func (client *Client) FindClusterByName(name string) (*Response, error) {
	// Find by name, then get by ID
	resp, err := client.ListClusters(&Request{
		QueryParams: map[string]string{
			"name": name,
		},
	})
	if err != nil {
		return resp, err
	}
	listResult := resp.Result.(*ListClustersResult)
	clustereCount := len(*listResult.Clusters)
	if clustereCount != 1 {
		return resp, fmt.Errorf("found %d Clusters for %v", clustereCount, name)
	}
	firstRecord := (*listResult.Clusters)[0]
	clustereId := firstRecord.ID
	return client.GetCluster(clustereId, &Request{})
}
