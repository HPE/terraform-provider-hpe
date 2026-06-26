package morpheus

import (
	"fmt"
)

var (
	// StorageServersPath is the API endpoint for storage servers
	StorageServersPath = "/api/storage-servers"
)

// StorageServer structures for use in request and response payloads
type StorageServer struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Type struct {
		ID   int64  `json:"id"`
		Code string `json:"code"`
		Name string `json:"name"`
	} `json:"type"`
	Chassis             string      `json:"chassis"`
	Visibility          string      `json:"visibility"`
	Description         string      `json:"description"`
	InternalId          string      `json:"internalId"`
	ExternalId          string      `json:"externalId"`
	ServiceUrl          string      `json:"serviceUrl"`
	ServiceHost         string      `json:"serviceHost"`
	ServicePath         string      `json:"servicePath"`
	ServiceToken        string      `json:"serviceToken"`
	ServiceTokenHash    string      `json:"serviceTokenHash"`
	ServiceVersion      string      `json:"serviceVersion"`
	ServiceUsername     string      `json:"serviceUsername"`
	ServicePassword     string      `json:"servicePassword"`
	ServicePasswordHash string      `json:"servicePasswordHash"`
	InternalIp          string      `json:"internalIp"`
	ExternalIp          string      `json:"externalIp"`
	ApiPort             interface{} `json:"apiPort"`
	AdminPort           interface{} `json:"adminPort"`
	Config              struct {
		Permissions      []string `json:"permissions"`
		StorageUser      string   `json:"storageUser"`
		StorageGroup     string   `json:"storageGroup"`
		ReadPermissions  []string `json:"readPermissions"`
		AdminPermissions []string `json:"adminPermissions"`
	} `json:"config"`
	RefType       string      `json:"refType"`
	RefId         int64       `json:"refId"`
	Category      string      `json:"category"`
	ServerVendor  string      `json:"serverVendor"`
	ServerModel   interface{} `json:"serverModel"`
	SerialNumber  interface{} `json:"serialNumber"`
	Status        string      `json:"status"`
	StatusMessage string      `json:"statusMessage"`
	StatusDate    string      `json:"statusDate"`
	Errormessage  string      `json:"errorMessage"`
	MaxStorage    interface{} `json:"maxStorage"`
	UsedStorage   interface{} `json:"usedStorage"`
	DiskCount     interface{} `json:"diskCount"`
	DateCreated   string      `json:"dateCreated"`
	LastUpdated   string      `json:"lastUpdated"`
	Enabled       bool        `json:"enabled"`
	Groups        []struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	} `json:"groups"`
	HostGroups []struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	} `json:"hostGroups"`
	Hosts []struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	} `json:"hosts"`
	Tenants []interface{} `json:"tenants"`
	Owner   struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	} `json:"owner"`
	Credential struct {
		Type string `json:"type"`
	} `json:"credential"`
}

type ListStorageServersResult struct {
	StorageServers *[]StorageServer `json:"storageServers"`
	Meta           *MetaResult      `json:"meta"`
}

type GetStorageServerResult struct {
	StorageServer *StorageServer `json:"storageServer"`
}

type CreateStorageServerResult struct {
	Success       bool              `json:"success"`
	Message       string            `json:"msg"`
	Errors        map[string]string `json:"errors"`
	StorageServer *StorageServer    `json:"storageServer"`
}

type UpdateStorageServerResult struct {
	CreateStorageServerResult
}

type DeleteStorageServerResult struct {
	DeleteResult
}

// Client request methods
func (client *Client) ListStorageServers(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        StorageServersPath,
		QueryParams: req.QueryParams,
		Result:      &ListStorageServersResult{},
	})
}

func (client *Client) GetStorageServer(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        fmt.Sprintf("%s/%d", StorageServersPath, id),
		QueryParams: req.QueryParams,
		Result:      &GetStorageServerResult{},
	})
}

func (client *Client) CreateStorageServer(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "POST",
		Path:        StorageServersPath,
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &CreateStorageServerResult{},
	})
}

func (client *Client) UpdateStorageServer(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "PUT",
		Path:        fmt.Sprintf("%s/%d", StorageServersPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &UpdateStorageServerResult{},
	})
}

func (client *Client) DeleteStorageServer(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "DELETE",
		Path:        fmt.Sprintf("%s/%d", StorageServersPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &DeleteStorageServerResult{},
	})
}

// FindStorageServerByName gets an existing storageServer by name
func (client *Client) FindStorageServerByName(name string) (*Response, error) {
	// Find by name, then get by ID
	resp, err := client.ListStorageServers(&Request{
		QueryParams: map[string]string{
			"name": name,
		},
	})
	if err != nil {
		return resp, err
	}
	listResult := resp.Result.(*ListStorageServersResult)
	storageServerCount := len(*listResult.StorageServers)
	if storageServerCount != 1 {
		return resp, fmt.Errorf("found %d Storage Servers for %v", storageServerCount, name)
	}
	firstRecord := (*listResult.StorageServers)[0]
	storageServerID := firstRecord.ID
	return client.GetStorageServer(storageServerID, &Request{})
}
