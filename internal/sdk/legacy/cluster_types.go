package morpheus

var (
	ClusterTypesPath = "/api/cluster-types"
)

// ClusterType structures for use in request and response payloads
type ClusterType struct {
	ID                   int64        `json:"id"`
	DeployTargetService  string       `json:"deployTargetService"`
	ShortName            string       `json:"shortName"`
	ProviderType         string       `json:"providerType"`
	Code                 string       `json:"code"`
	HostService          string       `json:"hostService"`
	Managed              bool         `json:"managed"`
	HasMasters           bool         `json:"hasMasters"`
	HasWorkers           bool         `json:"hasWorkers"`
	ViewSet              string       `json:"viewSet"`
	ImageCode            string       `json:"imageCode"`
	KubeCtlLocal         bool         `json:"kubeCtlLocal"`
	HasDatastore         bool         `json:"hasDatastore"`
	SupportsCloudScaling bool         `json:"supportsCloudScaling"`
	Name                 string       `json:"name"`
	HasDefaultDataDisk   bool         `json:"hasDefaultDataDisk"`
	CanManage            bool         `json:"canManage"`
	HasCluster           bool         `json:"hasCluster"`
	Description          string       `json:"description"`
	OptionTypes          []OptionType `json:"optionTypes"`
	ControllerTypes      []struct {
		ID          int64  `json:"id"`
		Name        string `json:"name"`
		Code        string `json:"code"`
		Description string `json:"description"`
	} `json:"controllerTypes"`
	WorkerTypes []struct {
		ID          int64  `json:"id"`
		Name        string `json:"name"`
		Code        string `json:"code"`
		Description string `json:"description"`
	} `json:"workerTypes"`
}

// ListClusterTypeResult structure parses the list cluster types response payload
type ListClusterTypesResult struct {
	ClusterTypes *[]ClusterType `json:"clusterTypes"`
	Meta         *MetaResult    `json:"meta"`
}

// GetClusterTypeResult structure parses the get cluster type response payload
type GetClusterTypeResult struct {
	ClusterType *ClusterType `json:"clusterType"`
}

// API endpoints
// ListClusterTypes lists all cluster types
func (client *Client) ListClusterTypes(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        ClusterTypesPath,
		QueryParams: req.QueryParams,
		Result:      &ListClusterTypesResult{},
	})
}

// FindClusterTypeByName gets an existing Cluster type by name
func (client *Client) FindClusterTypeByName(name string) (*Response, error) {
	// Find by name, then get by ID
	resp, err := client.ListClusterTypes(&Request{
		QueryParams: map[string]string{
			"name": name,
		},
	})
	if err != nil {
		return resp, err
	}
	return resp, err
}
