package morpheus

import (
	"fmt"
)

var (
	// NetworkPoolServersPath is the API endpoint for network pool servers
	NetworkPoolServersPath = "/api/networks/pool-servers"
)

// NetworkPoolServer structures for use in request and response payloads
type NetworkPoolServer struct {
	ID   int64 `json:"id"`
	Type struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
		Code string `json:"code"`
	} `json:"type"`
	Name                string `json:"name"`
	Enabled             bool   `json:"enabled"`
	ServiceUrl          string `json:"serviceUrl"`
	ServiceHost         string `json:"serviceHost"`
	ServicePort         int64  `json:"servicePort"`
	ServiceMode         string `json:"serviceMode"`
	ServiceUsername     string `json:"serviceUsername"`
	ServicePassword     string `json:"servicePassword"`
	ServicePasswordHash string `json:"servicePasswordHash"`
	Status              string `json:"status"`
	StatusMessage       string `json:"statusMessage"`
	StatusDate          string `json:"statusDate"`
	ServiceThrottleRate int64  `json:"serviceThrottleRate"`
	IgnoreSsl           bool   `json:"ignoreSsl"`
	Config              struct {
		AppId             string `json:"appId"`
		InventoryExisting string `json:"inventoryExisting"`
		ExtraAttributes   string `json:"extraAttributes"`
	} `json:"config"`
	NetworkFilter string `json:"networkFilter"`
	ZoneFilter    string `json:"zoneFilter"`
	TenantMatch   string `json:"tenantMatch"`
	DateCreated   string `json:"dateCreated"`
	LastUpdated   string `json:"lastUpdated"`
	Account       struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	} `json:"account"`
	Integration struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	} `json:"integration"`
	Pools []struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	} `json:"pools"`
	Credential struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
		Type string `json:"type"`
	}
}

// ListNetworkPoolsResult structure parses the list network pools response payload
type ListNetworkPoolServersResult struct {
	NetworkPoolServers *[]NetworkPoolServer `json:"networkPoolServers"`
	Meta               *MetaResult          `json:"meta"`
}

type GetNetworkPoolServerResult struct {
	NetworkPoolServer *NetworkPoolServer `json:"networkPoolServer"`
}

type CreateNetworkPoolServerResult struct {
	Success           bool               `json:"success"`
	Message           string             `json:"msg"`
	Errors            map[string]string  `json:"errors"`
	NetworkPoolServer *NetworkPoolServer `json:"networkPoolServer"`
}

type UpdateNetworkPoolServerResult struct {
	CreateNetworkPoolServerResult
}

type DeleteNetworkPoolServerResult struct {
	DeleteResult
}

// ListNetworkPoolServers lists all network pool servers
func (client *Client) ListNetworkPoolServers(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        NetworkPoolServersPath,
		QueryParams: req.QueryParams,
		Result:      &ListNetworkPoolServersResult{},
	})
}

// GetNetworkPool Server gets an existing network pool server
func (client *Client) GetNetworkPoolServer(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        fmt.Sprintf("%s/%d", NetworkPoolServersPath, id),
		QueryParams: req.QueryParams,
		Result:      &GetNetworkPoolServerResult{},
	})
}

// CreateNetworkPoolServer creates a new network pool server
func (client *Client) CreateNetworkPoolServer(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "POST",
		Path:        NetworkPoolServersPath,
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &CreateNetworkPoolServerResult{},
	})
}

// UpdateNetworkPoolServer updates an existing network pool server
func (client *Client) UpdateNetworkPoolServer(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "PUT",
		Path:        fmt.Sprintf("%s/%d", NetworkPoolServersPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &UpdateNetworkPoolServerResult{},
	})
}

// DeleteNetworkPoolServer deletes an existing network pool server
func (client *Client) DeleteNetworkPoolServer(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "DELETE",
		Path:        fmt.Sprintf("%s/%d", NetworkPoolServersPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &DeleteNetworkPoolServerResult{},
	})
}

// FindNetworkPoolServerByName gets an existing network pool server by name
func (client *Client) FindNetworkPoolServerByName(name string) (*Response, error) {
	// Find by name, then get by ID
	resp, err := client.ListNetworkPoolServers(&Request{
		QueryParams: map[string]string{
			"name": name,
		},
	})
	if err != nil {
		return resp, err
	}
	listResult := resp.Result.(*ListNetworkPoolServersResult)
	networkPoolServersCount := len(*listResult.NetworkPoolServers)
	if networkPoolServersCount != 1 {
		return resp, fmt.Errorf("found %d Network Pools Server for %v", networkPoolServersCount, name)
	}
	firstRecord := (*listResult.NetworkPoolServers)[0]
	networkPoolServerID := firstRecord.ID
	return client.GetNetworkPoolServer(networkPoolServerID, &Request{})
}
