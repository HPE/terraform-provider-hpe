package morpheus

import (
	"fmt"
)

var (
	// NetworkProxiesPath is the API endpoint for network proxies
	NetworkProxiesPath = "/api/networks/proxies"
)

// NetworkProxy structures for use in request and response payloads
type NetworkProxy struct {
	ID               int64  `json:"id"`
	Name             string `json:"name"`
	ProxyHost        string `json:"proxyHost"`
	ProxyPort        int64  `json:"proxyPort"`
	ProxyUser        string `json:"proxyUser"`
	ProxyPassword    string `json:"proxyPassword"`
	ProxyWorkstation string `json:"proxyWorkstation"`
	ProxyDomain      string `json:"proxyDomain"`
	Visibility       string `json:"visibility"`
	Account          struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	} `json:"account"`
	Owner struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	} `json:"owner"`
}

// ListNetowrkProxiesResult structure parses the list network proxies response payload
type ListNetworkProxiesResult struct {
	NetworkProxies *[]NetworkProxy `json:"networkProxies"`
	Meta           *MetaResult     `json:"meta"`
}

type GetNetworkProxyResult struct {
	NetworkProxy *NetworkProxy `json:"networkProxy"`
}

type CreateNetworkProxyResult struct {
	Success      bool              `json:"success"`
	Message      string            `json:"msg"`
	Errors       map[string]string `json:"errors"`
	NetworkProxy *NetworkProxy     `json:"networkProxy"`
}

type UpdateNetworkProxyResult struct {
	CreateNetworkProxyResult
}

type DeleteNetworkProxyResult struct {
	DeleteResult
}

// ListNetworkProxies lists all network proxies
func (client *Client) ListNetworkProxies(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        NetworkProxiesPath,
		QueryParams: req.QueryParams,
		Result:      &ListNetworkProxiesResult{},
	})
}

// GetNetworkProxy gets an existing network proxy
func (client *Client) GetNetworkProxy(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        fmt.Sprintf("%s/%d", NetworkProxiesPath, id),
		QueryParams: req.QueryParams,
		Result:      &GetNetworkProxyResult{},
	})
}

// CreateNetworkProxy creates a new network proxy
func (client *Client) CreateNetworkProxy(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "POST",
		Path:        NetworkProxiesPath,
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &CreateNetworkProxyResult{},
	})
}

// UpdateNetworkProxy updates an existing network proxy
func (client *Client) UpdateNetworkProxy(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "PUT",
		Path:        fmt.Sprintf("%s/%d", NetworkProxiesPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &UpdateNetworkProxyResult{},
	})
}

// DeleteNetworkProxy deletes an existing network proxy
func (client *Client) DeleteNetworkProxy(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "DELETE",
		Path:        fmt.Sprintf("%s/%d", NetworkProxiesPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &DeleteNetworkProxyResult{},
	})
}

// FindNetworkProxyByName gets an existing network proxy by name
func (client *Client) FindNetworkProxyByName(name string) (*Response, error) {
	// Find by name, then get by ID
	resp, err := client.ListNetworkProxies(&Request{
		QueryParams: map[string]string{
			"name": name,
		},
	})
	if err != nil {
		return resp, err
	}
	listResult := resp.Result.(*ListNetworkProxiesResult)
	networkProxiesCount := len(*listResult.NetworkProxies)
	if networkProxiesCount != 1 {
		return resp, fmt.Errorf("found %d Network Proxies for %v", networkProxiesCount, name)
	}
	firstRecord := (*listResult.NetworkProxies)[0]
	networkProxyID := firstRecord.ID
	return client.GetNetworkProxy(networkProxyID, &Request{})
}
