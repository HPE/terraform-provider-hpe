package morpheus

import (
	"fmt"
)

var (
	// CloudsPath is the API endpoint for clouds (zones)
	CloudsPath     = "/api/zones"
	CloudTypesPath = "/api/zone-types"
)

// Cloud structures for use in request and response payloads
type Cloud struct {
	ID         int64    `json:"id"`
	UUID       string   `json:"uuid"`
	ExternalID string   `json:"externalId"`
	Name       string   `json:"name"`
	Code       string   `json:"code"`
	Labels     []string `json:"labels"`
	Location   string   `json:"location"`
	Owner      struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	} `json:"owner"`
	AccountID int64 `json:"accountId"`
	Account   struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	} `json:"account"`
	Visibility           string    `json:"visibility"`
	Enabled              bool      `json:"enabled"`
	Status               string    `json:"status"`
	StatusMessage        string    `json:"statusMessage"`
	StatusDate           string    `json:"statusDate"`
	LastSync             string    `json:"lastSync"`
	NextRunDate          string    `json:"nextRunDate"`
	LastSyncDuration     int64     `json:"lastSyncDuration"`
	CostStatus           string    `json:"costStatus"`
	CostStatusMessage    string    `json:"costStatusMessage"`
	CostStatusDate       string    `json:"costStatusDate"`
	CostLastSyncDuration int64     `json:"costLastSyncDuration"`
	CostLastSync         string    `json:"costLastSync"`
	CloudType            CloudType `json:"zoneType"`
	CloudTypeID          int64     `json:"zoneTypeId"`
	GuidanceMode         string    `json:"guidanceMode"`
	StorageMode          string    `json:"storageMode"`
	AgentMode            string    `json:"agentMode"`
	UserDataLinux        string    `json:"userDataLinux"`
	UserDataWindows      string    `json:"userDataWindows"`
	ConsoleKeymap        string    `json:"consoleKeymap"`
	ContainerMode        string    `json:"containerMode"`
	CostingMode          string    `json:"costingMode"`
	ServiceVersion       string    `json:"serviceVersion"`
	SecurityMode         string    `json:"securityMode"`
	InventoryLevel       string    `json:"inventoryLevel"`
	TimeZone             string    `json:"timezone"`
	NetworkDomain        struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	} `json:"networkDomain"`
	DomainName            string `json:"domainName"`
	RegionCode            string `json:"regionCode"`
	AutoRecoverPowerState bool   `json:"autoRecoverPowerState"`
	ScalePriority         int64  `json:"scalePriority"`
	Config                struct {
		// AWS
		Endpoint             string `json:"endpoint"`
		IsVpc                string `json:"isVpc"`
		ImageStoreId         string `json:"imageStoreId"`
		EbsEncryption        string `json:"ebsEncryption"`
		CostingReport        string `json:"costingReport"`
		CostingRegion        string `json:"costingRegion"`
		CostingSecretKeyHash string `json:"costingSecretKeyHash"`
		SecretKeyHash        string `json:"secretKeyHash"`
		AccessKey            string `json:"accessKey"`
		SecretKey            string `json:"secretKey"`
		VPC                  string `json:"vpc"`
		StsAssumeRole        string `json:"stsAssumeRole"`
		UseHostCredentials   string `json:"useHostCredentials"`
		CostingAccessKey     string `json:"costingAccessKey"`
		CostingBucketName    string `json:"costingBucketName"`
		CostingFolder        string `json:"costingFolder"`
		CostingReportName    string `json:"costingReportName"`
		CostingSecretKey     string `json:"costingSecretKey"`

		// vSphere
		APIUrl                     string `json:"apiUrl"`
		APIVersion                 string `json:"apiVersion"`
		Datacenter                 string `json:"datacenter"`
		Cluster                    string `json:"cluster"`
		DiskStorageType            string `json:"diskStorageType"`
		DatacenterID               string `json:"datacenterId"`
		EnableVNC                  string `json:"enableVnc"`
		DiskEncryption             string `json:"diskEncryption"`
		EnableDiskTypeSelection    string `json:"enableDiskTypeSelection"`
		EnableStorageTypeSelection string `json:"enableStorageTypeSelection"`
		EnableNetworkTypeSelection string `json:"enableNetworkTypeSelection"`
		ResourcePool               string `json:"resourcePool"`
		ResourcePoolId             string `json:"resourcePoolId"`
		HideHostSelection          string `json:"hideHostSelection"`

		// Azure
		AzureCostingMode    string `json:"azureCostingMode"`
		SubscriberID        string `json:"subscriberId"`
		TenantID            string `json:"tenantId"`
		ClientID            string `json:"clientId"`
		ClientSecret        string `json:"clientSecret"`
		ClientSecretHash    string `json:"clientSecretHash"`
		ResourceGroup       string `json:"resourceGroup"`
		CSPCustomer         string `json:"cspCustomer"`
		CSPTenantID         string `json:"cspTenantId"`
		CSPClientID         string `json:"cspClientId"`
		CSPClientSecret     string `json:"cspClientSecret"`
		CSPClientSecretHash string `json:"cspClientSecretHash"`

		// GCP
		GoogleRegionID string `json:"googleRegionId"`
		GoogleBucket   string `json:"googleBucket"`
		PrivateKey     string `json:"privateKey"`
		PrivateKeyHash string `json:"privateKeyHash"`
		ClientEmail    string `json:"clientEmail"`

		// Hyperv
		Provider    string `json:"provider"`
		Host        string `json:"host"`
		WorkingPath string `json:"workingPath"`
		VMPath      string `json:"vmPath"`
		DiskPath    string `json:"diskPath"`

		// OpenStack
		IdentityApi                  string `json:"identityApi"`
		DomainId                     string `json:"domainId"`
		ProjectName                  string `json:"projectName"`
		OsRelease                    string `json:"osRelease"`
		DiskMode                     string `json:"diskMode"`
		IdentityVersion              string `json:"identityVersion"`
		ComputeApi                   string `json:"computeApi"`
		ComputeVersion               string `json:"computeVersion"`
		ImageApi                     string `json:"imageApi"`
		ImageVersion                 string `json:"imageVersion"`
		StorageApi                   string `json:"storageApi"`
		StorageVersion               string `json:"storageVersion"`
		NetworkApi                   string `json:"networkApi"`
		NetworkVersion               string `json:"networkVersion"`
		ProjectId                    string `json:"projectId"`
		ApiProjectId                 string `json:"apiProjectId"`
		ApiTokenExpiresAt            string `json:"apiTokenExpiresAt"`
		LbaasType                    string `json:"lbaasType"`
		ApiDomainId                  string `json:"apiDomainId"`
		ApiUserId                    string `json:"apiUserId"`
		ProvisionMethod              string `json:"provisionMethod"`
		ComputeMicroVersion          string `json:"computeMicroVersion"`
		ImageMicroVersion            string `json:"imageMicroVersion"`
		StorageMicroVersion          string `json:"storageMicroVersion"`
		NetworkMicroVersion          string `json:"networkMicroVersion"`
		LoadBalancerApi              string `json:"loadBalancerApi"`
		LoadBalancerVersion          string `json:"loadBalancerVersion"`
		LoadBalancerMicroVersion     string `json:"loadBalancerMicroVersion"`
		LoadBalancerV1Api            string `json:"loadBalancerV1Api"`
		LoadBalancerV1Version        string `json:"loadBalancerV1Version"`
		LoadBalancerV1MicroVersion   string `json:"loadBalancerV1MicroVersion"`
		ObjectStorageApi             string `json:"objectStorageApi"`
		ObjectStorageVersion         string `json:"objectStorageVersion"`
		ObjectStorageMicroVersion    string `json:"objectStorageMicroVersion"`
		SharedFileSystemApi          string `json:"sharedFileSystemApi"`
		SharedFileSystemVersion      string `json:"sharedFileSystemVersion"`
		SharedFileSystemMicroVersion string `json:"sharedFileSystemMicroVersion"`
		Region                       string `json:"region"`
		KubeUrl                      string `json:"kubeUrl"`

		// MaaS
		ReleasePoolName string `json:"releasePoolName"`
		ReleaseMode     string `json:"releaseMode"`
		AvailableFilter string `json:"availableFilter"`

		// VCD
		ServiceToken          string `json:"serviceToken"`
		OrgID                 string `json:"orgId"`
		VDCID                 string `json:"vdcId"`
		VCDVersion            string `json:"vcdVersion"`
		DefaultStorageProfile string `json:"defaultStorageProfile"`
		Catalog               string `json:"catalog"`

		// General
		ClusterRef      string `json:"clusterRef"`
		ApplianceUrl    string `json:"applianceUrl"`
		DatacenterName  string `json:"datacenterName"`
		ImportExisting  string `json:"importExisting"`
		InventoryLevel  string `json:"inventoryLevel"`
		NetworkServerID string `json:"networkServer.id"`
		KubernetesToken string `json:"kubernetes-token"`
		NetworkServer   struct {
			ID   string `json:"id"`
			Name string `json:"name"`
			Type string `json:"type"`
		} `json:"networkServer"`
		SecurityMode        string `json:"securityMode"`
		CertificateProvider string `json:"certificateProvider"`
		BackupMode          string `json:"backupMode"`
		ReplicationMode     string `json:"replicationMode"`
		DnsIntegrationID    string `json:"dnsIntegrationId"`
		ServiceRegistryID   string `json:"serviceRegistryId"`
		ConfigManagementID  string `json:"configManagementId"`
		ConfigCmdbID        string `json:"configCmdbId"`
		SecurityServer      string `json:"securityServer"`
		CloudType           string `json:"cloudType"`
		AccountType         string `json:"accountType"`
		RPCMode             string `json:"rpcMode"`
		EncryptionSet       string `json:"encryptionSet"`
		ConfigCmID          string `json:"configCmId"`
		CostingProjectID    string `json:"costingProjectId"`
		ConfigCMDBDiscovery bool   `json:"configCmdbDiscovery"`
		CostingDatasetID    string `json:"costingDatasetId"`
		Username            string `json:"username"`
		Password            string `json:"password"`
		MasterAddress       string `json:"masterAddress"`
		DistributedWorkerId string `json:"distributedWorkerId"`
		PasswordHash        string `json:"passwordHash"`
	} `json:"config"`
	Credential struct {
		Type string `json:"type"`
		ID   int64  `json:"id"`
		Name string `json:"name"`
	} `json:"credential"`
	ImagePath     string `json:"imagePath"`
	DarkImagePath string `json:"darkImagePath"`
	DateCreated   string `json:"dateCreated"`
	LastUpdated   string `json:"lastUpdated"`
	Groups        []struct {
		ID        int64  `json:"id"`
		Name      string `json:"name"`
		AccountID int64  `json:"accountId"`
	} `json:"groups"`
}

// ListCloudsResult structure parses the list clouds response payload
type ListCloudsResult struct {
	Clouds *[]Cloud    `json:"zones"`
	Meta   *MetaResult `json:"meta"`
}

type GetCloudResult struct {
	Cloud *Cloud `json:"zone"`
}

type CreateCloudResult struct {
	Success bool              `json:"success"`
	Message string            `json:"msg"`
	Errors  map[string]string `json:"errors"`
	Cloud   *Cloud            `json:"zone"`
}

type UpdateCloudResult struct {
	CreateCloudResult
}

type UpdateCloudLogoResult struct {
	StandardResult
}

type RefreshCloudResult struct {
	StandardResult
}

type DeleteCloudResult struct {
	DeleteResult
}

// Datastores
type Datastore struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Zone struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	} `json:"zone"`
	Type       string `json:"type"`
	FreeSpace  int64  `json:"freeSpace"`
	Online     bool   `json:"online"`
	Active     bool   `json:"active"`
	Visibility string `json:"visibility"`
	Tenants    []struct {
		ID            int    `json:"id"`
		Name          string `json:"name"`
		DefaultStore  bool   `json:"defaultStore"`
		DefaultTarget bool   `json:"defaultTarget"`
	} `json:"tenants"`
	ResourcePermission struct {
		All   bool `json:"all"`
		Sites []struct {
			ID      int    `json:"id"`
			Name    string `json:"name"`
			Default bool   `json:"default"`
		} `json:"sites"`
		AllPlans bool          `json:"allPlans"`
		Plans    []interface{} `json:"plans"`
	} `json:"resourcePermission"`
}

type ListCloudDatastoresResult struct {
	Datastores *[]Datastore `json:"datastores"`
	Meta       *MetaResult  `json:"meta"`
}

type GetCloudDatastoreResult struct {
	Datastore *Datastore `json:"datastore"`
}

type UpdateCloudDatastoreResult struct {
	Success   bool              `json:"success"`
	Message   string            `json:"msg"`
	Errors    map[string]string `json:"errors"`
	Datastore *Datastore        `json:"datastore"`
}

// Resource Folder
type Folder struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Zone struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	} `json:"zone"`
	Parent struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	} `json:"parent"`
	Type          string `json:"type"`
	ExternalId    string `json:"externalId"`
	Visibility    string `json:"visibility"`
	ReadOnly      bool   `json:"readOnly"`
	DefaultFolder bool   `json:"defaultFolder"`
	DefaultStore  bool   `json:"defaultStore"`
	Active        bool   `json:"active"`
	Tenants       []struct {
		ID            int64  `json:"id"`
		Name          string `json:"name"`
		DefaultStore  bool   `json:"defaultStore"`
		DefaultTarget bool   `json:"defaultTarget"`
	} `json:"tenants"`
	ResourcePermission struct {
		All      bool          `json:"all"`
		Sites    []interface{} `json:"sites"`
		AllPlans bool          `json:"allPlans"`
		Plans    []interface{} `json:"plans"`
	} `json:"resourcePermission"`
}

type ListCloudResourceFoldersResult struct {
	Folders *[]Folder   `json:"folders"`
	Meta    *MetaResult `json:"meta"`
}

type GetCloudResourceFolderResult struct {
	Folder *Folder `json:"folder"`
}

type UpdateCloudResourceFolderResult struct {
	Success bool              `json:"success"`
	Message string            `json:"msg"`
	Errors  map[string]string `json:"errors"`
	Folder  *Folder           `json:"folder"`
}

// Resource Pool
type Pool struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Zone struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	} `json:"zone"`
	Parent struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	} `json:"parent"`
	Type        string      `json:"type"`
	ExternalId  string      `json:"externalId"`
	RegionCode  string      `json:"regionCode"`
	Visibility  string      `json:"visibility"`
	ReadOnly    bool        `json:"readOnly"`
	DefaultPool bool        `json:"defaultPool"`
	Active      bool        `json:"active"`
	Status      string      `json:"status"`
	Inventory   bool        `json:"inventory"`
	Config      interface{} `json:"config"`
	Tenants     []struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	} `json:"tenants"`
	ResourcePermission struct {
		All   bool `json:"all"`
		Sites []struct {
			ID      int64  `json:"id"`
			Name    string `json:"name"`
			Default bool   `json:"default"`
		} `json:"sites"`
		AllPlans bool `json:"allPlans"`
		Plans    []struct {
			ID      int64  `json:"id"`
			Name    string `json:"name"`
			Default bool   `json:"default"`
		} `json:"plans"`
	} `json:"resourcePermission"`
	Depth int64 `json:"depth"`
}

type ListCloudResourcePoolsResult struct {
	Pools *[]Pool     `json:"resourcePools"`
	Meta  *MetaResult `json:"meta"`
}

type GetCloudResourcePoolResult struct {
	Pool *Pool `json:"resourcePool"`
}

type CreateCloudResourcePoolResult struct {
	Success bool              `json:"success"`
	Message string            `json:"msg"`
	Errors  map[string]string `json:"errors"`
	Cloud   *Cloud            `json:"zone"`
}

type UpdateCloudResourcePoolResult struct {
	Success bool              `json:"success"`
	Message string            `json:"msg"`
	Errors  map[string]string `json:"errors"`
	Pool    *Pool             `json:"resourcePool"`
}

type DeleteCloudResourcePoolResult struct {
	DeleteResult
}

// Cloud Types
type CloudType struct {
	ID                            int64        `json:"id"`
	Name                          string       `json:"name"`
	Code                          string       `json:"code"`
	Enabled                       bool         `json:"enabled"`
	Provision                     bool         `json:"provision"`
	AutoCapacity                  bool         `json:"autoCapacity"`
	MigrationTarget               bool         `json:"migrationTarget"`
	HasDatastores                 bool         `json:"hasDatastores"`
	HasNetworks                   bool         `json:"hasNetworks"`
	HasResourcePools              bool         `json:"hasResourcePools"`
	HasSecurityGroups             bool         `json:"hasSecurityGroups"`
	HasContainers                 bool         `json:"hasContainers"`
	HasBareMetal                  bool         `json:"hasBareMetal"`
	HasServices                   bool         `json:"hasServices"`
	HasFunctions                  bool         `json:"hasFunctions"`
	HasJobs                       bool         `json:"hasJobs"`
	HasDiscovery                  bool         `json:"hasDiscovery"`
	HasCloudInit                  bool         `json:"hasCloudInit"`
	HasFolders                    bool         `json:"hasFolders"`
	HasFloatingIps                bool         `json:"hasFloatingIps"`
	HasMarketplace                bool         `json:"hasMarketplace"`
	CanCreateResourcePools        bool         `json:"canCreateResourcePools"`
	CanDeleteResourcePools        bool         `json:"canDeleteResourcePools"`
	CanCreateDatastores           bool         `json:"canCreateDatastores"`
	CanCreateNetworks             bool         `json:"canCreateNetworks"`
	CanChooseContainerMode        bool         `json:"canChooseContainerMode"`
	ProvisionRequiresResourcePool bool         `json:"provisionRequiresResourcePool"`
	SupportsDistributedWorker     bool         `json:"supportsDistributedWorker"`
	Cloud                         string       `json:"cloud"`
	ProvisionTypes                []int64      `json:"provisionTypes"`
	ZoneInstanceTypeLayoutId      int64        `json:"zoneInstanceTypeLayoutId"`
	ServerTypes                   []ServerType `json:"serverTypes"`
	OptionTypes                   []OptionType `json:"optionTypes"`
}

type ServerType struct {
	ID                  int           `json:"id"`
	Code                string        `json:"code"`
	Name                string        `json:"name"`
	Description         string        `json:"description"`
	NodeType            string        `json:"nodeType"`
	Platform            string        `json:"platform"`
	Enabled             bool          `json:"enabled"`
	Selectable          bool          `json:"selectable"`
	ExternalDelete      bool          `json:"externalDelete"`
	Managed             bool          `json:"managed"`
	ControlPower        bool          `json:"controlPower"`
	ControlSuspend      bool          `json:"controlSuspend"`
	Creatable           bool          `json:"creatable"`
	HasAgent            bool          `json:"hasAgent"`
	VmHypervisor        bool          `json:"vmHypervisor"`
	ContainerHypervisor bool          `json:"containerHypervisor"`
	BareMetalHost       bool          `json:"bareMetalHost"`
	GuestVm             bool          `json:"guestVm"`
	HasAutomation       bool          `json:"hasAutomation"`
	ProvisionType       ProvisionType `json:"provisionType"`
	OptionTypes         []OptionType  `json:"optionTypes"`
	DisplayOrder        int64         `json:"displayOrder"`
}

type ListCloudTypesResult struct {
	CloudTypes *[]CloudType `json:"zoneTypes"`
	Meta       *MetaResult  `json:"meta"`
}

type GetCloudTypeResult struct {
	CloudType *CloudType `json:"zoneType"`
}

// API endpoints
func (client *Client) ListClouds(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        CloudsPath,
		QueryParams: req.QueryParams,
		Result:      &ListCloudsResult{},
	})
}

func (client *Client) GetCloud(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        fmt.Sprintf("%s/%d", CloudsPath, id),
		QueryParams: req.QueryParams,
		Result:      &GetCloudResult{},
	})
}

// CreateCloud creates a new cloud
func (client *Client) CreateCloud(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "POST",
		Path:        CloudsPath,
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &CreateCloudResult{},
	})
}

// UpdateCloud updates an existing cloud
func (client *Client) UpdateCloud(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "PUT",
		Path:        fmt.Sprintf("%s/%d", CloudsPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &UpdateCloudResult{},
	})
}

func (client *Client) UpdateCloudLogo(id int64, filePayload []*FilePayload, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:         "POST",
		Path:           fmt.Sprintf("%s/%d/update-logo", CloudsPath, id),
		IsMultiPart:    true,
		MultiPartFiles: filePayload,
		Headers: map[string]string{
			"Content-Type": "application/x-www-form-urlencoded",
		},
		Result: &UpdateCloudLogoResult{},
	})
}

// DeleteCloud deletes an existing cloud
func (client *Client) DeleteCloud(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "DELETE",
		Path:        fmt.Sprintf("%s/%d", CloudsPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &DeleteCloudResult{},
	})
}

// RefreshCloud refreshes an existing cloud
func (client *Client) RefreshCloud(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "POST",
		Path:        fmt.Sprintf("%s/%d", CloudsPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &DeleteCloudResult{},
	})
}

func (client *Client) FindCloudByName(name string) (*Response, error) {
	// Find by name, then get by ID
	resp, err := client.ListClouds(&Request{
		QueryParams: map[string]string{
			"name": name,
		},
	})
	if err != nil {
		return resp, err
	}
	listResult := resp.Result.(*ListCloudsResult)
	cloudsCount := len(*listResult.Clouds)
	if cloudsCount != 1 {
		return resp, fmt.Errorf("found %d Clouds for %v", cloudsCount, name)
	}
	firstRecord := (*listResult.Clouds)[0]
	cloudId := firstRecord.ID
	return client.GetCloud(cloudId, &Request{})
}

func (client *Client) ListCloudDatastores(zoneId int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        fmt.Sprintf("%s/%d/data-stores", CloudsPath, zoneId),
		QueryParams: req.QueryParams,
		Result:      &ListCloudDatastoresResult{},
	})
}

func (client *Client) GetCloudDatastore(zoneId int64, id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        fmt.Sprintf("%s/%d/data-stores/%d", CloudsPath, zoneId, id),
		QueryParams: req.QueryParams,
		Result:      &GetCloudDatastoreResult{},
	})
}

func (client *Client) UpdateCloudDatastore(zoneId int64, id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "PUT",
		Path:        fmt.Sprintf("%s/%d/data-stores/%d", CloudsPath, zoneId, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &UpdateCloudDatastoreResult{},
	})
}

func (client *Client) ListCloudResourceFolders(zoneId int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        fmt.Sprintf("%s/%d/folders", CloudsPath, zoneId),
		QueryParams: req.QueryParams,
		Result:      &ListCloudResourceFoldersResult{},
	})
}

func (client *Client) GetCloudResourceFolder(zoneId int64, id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        fmt.Sprintf("%s/%d/folders/%d", CloudsPath, zoneId, id),
		QueryParams: req.QueryParams,
		Result:      &GetCloudResourceFolderResult{},
	})
}

func (client *Client) UpdateCloudResourceFolder(zoneId int64, id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "PUT",
		Path:        fmt.Sprintf("%s/%d/folders/%d", CloudsPath, zoneId, id),
		QueryParams: req.QueryParams,
		Result:      &UpdateCloudResourceFolderResult{},
	})
}

// ListCloudResourcePools fetches all existing cloud resource pools
func (client *Client) ListCloudResourcePools(zoneId int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        fmt.Sprintf("%s/%d/resource-pools", CloudsPath, zoneId),
		QueryParams: req.QueryParams,
		Result:      &ListCloudResourcePoolsResult{},
	})
}

// GetCloudResourcePool fetches an existing cloud resource pool
func (client *Client) GetCloudResourcePool(zoneId int64, id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        fmt.Sprintf("%s/%d/resource-pools/%d", CloudsPath, zoneId, id),
		QueryParams: req.QueryParams,
		Result:      &GetCloudResourcePoolResult{},
	})
}

// CreateCloudResourcePool creates a new cloud resource pool
func (client *Client) CreateCloudResourcePool(zoneId int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "POST",
		Path:        fmt.Sprintf("%s/%d/resource-pools", CloudsPath, zoneId),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &CreateCloudResourcePoolResult{},
	})
}

// UpdateCloudResourcePool updates an existing cloud resource pool
func (client *Client) UpdateCloudResourcePool(zoneId int64, id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "PUT",
		Path:        fmt.Sprintf("%s/%d/resource-pools/%d", CloudsPath, zoneId, id),
		QueryParams: req.QueryParams,
		Result:      &UpdateCloudResourcePoolResult{},
	})
}

// DeleteCloudResourcePool deletes an existing cloud resource pool
func (client *Client) DeleteCloudResourcePool(zoneId int64, id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "DELETE",
		Path:        fmt.Sprintf("%s/%d/resource-pools/%d", CloudsPath, zoneId, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &DeleteCloudResourcePoolResult{},
	})
}

// Cloud Types

// ListCloudTypes fetches existing cloud types
func (client *Client) ListCloudTypes(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        CloudTypesPath,
		QueryParams: req.QueryParams,
		Result:      &ListCloudTypesResult{},
	})
}

// GetCloudResourcePool fetches an existing cloud type
func (client *Client) GetCloudType(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        fmt.Sprintf("%s/%d", CloudTypesPath, id),
		QueryParams: req.QueryParams,
		Result:      &GetCloudTypeResult{},
	})
}

func (client *Client) FindCloudTypeByName(name string) (*Response, error) {
	// Find by name, then get by ID
	resp, err := client.ListCloudTypes(&Request{
		QueryParams: map[string]string{
			"name": name,
		},
	})
	if err != nil {
		return resp, err
	}
	listResult := resp.Result.(*ListCloudTypesResult)
	cloudTypesCount := len(*listResult.CloudTypes)
	fmt.Println(cloudTypesCount)
	if cloudTypesCount != 1 {
		return resp, fmt.Errorf("found %d Clouds Types for %v", cloudTypesCount, name)
	}
	firstRecord := (*listResult.CloudTypes)[0]
	cloudTypeId := firstRecord.ID
	fmt.Println(cloudTypeId)
	return client.GetCloudType(cloudTypeId, &Request{})
}
