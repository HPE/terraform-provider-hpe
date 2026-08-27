package morpheus

import (
	"fmt"
)

var (
	// LoadBalancerVirtualServersPath is the API endpoint for load balancer virtual servers
	LoadBalancerVirtualServersPath = "/api/load-balancers"
)

// LoadBalancerVirtualServer structures for use in request and response payloads
type LoadBalancerVirtualServer struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	LoadBalancer struct {
		ID   int64 `json:"id"`
		Type struct {
			ID   int64  `json:"id"`
			Name string `json:"name"`
			Code string `json:"code"`
		} `json:"type"`
		Name string `json:"name"`
		IP   string `json:"ip"`
	} `json:"loadBalancer"`
	Code            string      `json:"code"`
	Category        string      `json:"category"`
	Visibility      string      `json:"visibility"`
	Instance        interface{} `json:"instance"`
	Description     string      `json:"description"`
	InternalId      string      `json:"internalId"`
	ExternalId      string      `json:"externalId"`
	Active          bool        `json:"active"`
	Sticky          bool        `json:"sticky"`
	SslEnabled      interface{} `json:"sslEnabled"`
	ExternalAddress bool        `json:"externalAddress"`
	BackendPort     interface{} `json:"backendPort"`
	VipType         interface{} `json:"vipType"`
	VipAddress      string      `json:"vipAddress"`
	VipHostname     string      `json:"vipHostname"`
	VipProtocol     string      `json:"vipProtocol"`
	VipScheme       interface{} `json:"vipScheme"`
	VipMode         string      `json:"vipMode"`
	VipName         string      `json:"vipName"`
	VipPort         int         `json:"vipPort"`
	VipSticky       interface{} `json:"vipSticky"`
	VipBalance      interface{} `json:"vipBalance"`
	ServicePort     interface{} `json:"servicePort"`
	SourceAddress   interface{} `json:"sourceAddress"`
	SslCert         struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"sslCert"`
	SslMode          interface{} `json:"sslMode"`
	SslRedirectMode  interface{} `json:"sslRedirectMode"`
	VipShared        bool        `json:"vipShared"`
	VipDirectAddress interface{} `json:"vipDirectAddress"`
	ServerName       interface{} `json:"serverName"`
	PoolName         interface{} `json:"poolName"`
	Removing         bool        `json:"removing"`
	VipSource        string      `json:"vipSource"`
	ExtraConfig      interface{} `json:"extraConfig"`
	ServiceAccess    interface{} `json:"serviceAccess"`
	NetworkId        interface{} `json:"networkId"`
	SubnetId         interface{} `json:"subnetId"`
	ExternalPortId   interface{} `json:"externalPortId"`
	Status           string      `json:"status"`
	VipStatus        string      `json:"vipStatus"`
	DateCreated      string      `json:"dateCreated"`
	LastUpdated      string      `json:"lastUpdated"`
}

// ListLoadBalancerVirtualServersResult structure parses the list load balancers response payload
type ListLoadBalancerVirtualServersResult struct {
	LoadBalancerVirtualServers *[]LoadBalancerVirtualServer `json:"loadBalancerInstances"`
	Meta                       *MetaResult                  `json:"meta"`
}

type GetLoadBalancerVirtualServerResult struct {
	LoadBalancerVirtualServer *LoadBalancerProfile `json:"loadBalancerInstance"`
}

type CreateLoadBalancerVirtualServerResult struct {
	Success                   bool                       `json:"success"`
	Message                   string                     `json:"msg"`
	Errors                    map[string]string          `json:"errors"`
	LoadBalancerVirtualServer *LoadBalancerVirtualServer `json:"loadBalancerInstance"`
}

type UpdateLoadBalancerVirtualServerResult struct {
	CreateLoadBalancerVirtualServerResult
}

type DeleteLoadBalancerVirtualServerResult struct {
	DeleteResult
}

// ListLoadBalancerVirtualServers lists all load balancer virtual servers
func (client *Client) ListLoadBalancerVirtualServers(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        fmt.Sprintf("%s/%d/virtual-servers", LoadBalancerVirtualServersPath, id),
		QueryParams: req.QueryParams,
		Result:      &ListLoadBalancerVirtualServersResult{},
	})
}

// GetLoadBalancerVirtualServer gets an existing load balancer virtual server
func (client *Client) GetLoadBalancerVirtualServer(loadBalancerId int64, id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        fmt.Sprintf("%s/%d/virtual-servers/%d", LoadBalancerVirtualServersPath, loadBalancerId, id),
		QueryParams: req.QueryParams,
		Result:      &GetLoadBalancerVirtualServerResult{},
	})
}

// CreateLoadBalancerVirtualServer creates a new load balancer virtual server
func (client *Client) CreateLoadBalancerVirtualServer(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "POST",
		Path:        fmt.Sprintf("%s/virtual-servers", LoadBalancerVirtualServersPath),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &CreateLoadBalancerVirtualServerResult{},
	})
}

// UpdateLoadBalancerVirtualServer updates an existing load balancer virtual server
func (client *Client) UpdateLoadBalancerVirtualServer(loadBalancerId int64, id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "PUT",
		Path:        fmt.Sprintf("%s/%d/virtual-servers/%d", LoadBalancerVirtualServersPath, loadBalancerId, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &UpdateLoadBalancerVirtualServerResult{},
	})
}

// DeleteLoadBalancerVirtualServer deletes an existing load balancer virtual server
func (client *Client) DeleteLoadBalancerVirtualServer(loadBalancerId int64, id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "DELETE",
		Path:        fmt.Sprintf("%s/%d/virtual-servers/%d", LoadBalancerVirtualServersPath, loadBalancerId, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &DeleteLoadBalancerVirtualServerResult{},
	})
}

// FindLoadBalancerVirtualServerByName gets an existing load balancer virtual server by name
func (client *Client) FindLoadBalancerVirtualServerByName(loadBalancerId int64, name string) (*Response, error) {
	// Find by name, then get by ID
	resp, err := client.ListLoadBalancerVirtualServers(loadBalancerId, &Request{
		QueryParams: map[string]string{
			"name": name,
		},
	})
	if err != nil {
		return resp, err
	}
	listResult := resp.Result.(*ListLoadBalancerVirtualServersResult)
	loadBalancerVirtualServersCount := len(*listResult.LoadBalancerVirtualServers)
	if loadBalancerVirtualServersCount != 1 {
		return resp, fmt.Errorf("found %d Load Balancer VirtualServers for %v", loadBalancerVirtualServersCount, name)
	}
	firstRecord := (*listResult.LoadBalancerVirtualServers)[0]
	loadBalancerVirtualServerID := firstRecord.ID
	return client.GetLoadBalancerProfile(loadBalancerId, loadBalancerVirtualServerID, &Request{})
}
