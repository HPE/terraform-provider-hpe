package morpheus

import (
	"fmt"
)

var (
	// NetworkSubnetsPath is the API endpoint for network subnets
	NetworkSubnetsPath = "/api/subnets"
)

// NetworkSubnet structures for use in request and response payloads
type NetworkSubnet struct {
	ID             int64       `json:"id"`
	Name           string      `json:"name"`
	Code           string      `json:"code"`
	Labels         []string    `json:"labels"`
	Active         bool        `json:"active"`
	Description    string      `json:"description"`
	ExternalId     string      `json:"externalId"`
	UniqueId       string      `json:"uniqueId"`
	AddressPrefix  string      `json:"addressPrefix"`
	Cidr           string      `json:"cidr"`
	Gateway        string      `json:"gateway"`
	Netmask        string      `json:"netmask"`
	SubnetAddress  string      `json:"subnetAddress"`
	TftpServer     string      `json:"tftpServer"`
	BootFile       string      `json:"bootFile"`
	Pool           interface{} `json:"pool"`
	Dhcpserver     bool        `json:"dhcpServer"`
	Hasfloatingips bool        `json:"hasFloatingIps"`
	DhcpIp         string      `json:"dhcpIp"`
	DnsPrimary     string      `json:"dnsPrimary"`
	DnsSecondary   string      `json:"dnsSecondary"`
	DhcpStart      string      `json:"dhcpStart"`
	DhcpEnd        string      `json:"dhcpEnd"`
	DhcpRange      interface{} `json:"dhcpRange"`
	NetworkSubnet  interface{} `json:"networkSubnet"`
	NetworkDomain  interface{} `json:"networkDomain"`
	SearchDomains  interface{} `json:"searchDomains"`
	DefaultNetwork bool        `json:"defaultNetwork"`
	AssignPublicIp bool        `json:"assignPublicIp"`
	Status         struct {
		Name     string `json:"name"`
		EnumType string `json:"enumType"`
	} `json:"status"`
	Network struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	} `json:"network"`
	Type struct {
		ID   int64  `json:"id"`
		Code string `json:"code"`
		Name string `json:"name"`
	} `json:"type"`
	Account struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	} `json:"account"`
	Securitygroups []interface{} `json:"securityGroups"`
	Tenants        []struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	} `json:"tenants"`
	ResourcePermission struct {
		All      bool          `json:"all"`
		Sites    []interface{} `json:"sites"`
		AllPlans bool          `json:"allPlans"`
		Plans    []interface{} `json:"plans"`
	} `json:"resourcePermission"`
	Visibility string `json:"visibility"`
}

// ListNetowrkSubnetsResult structure parses the list network subnets response payload
type ListNetworkSubnetsResult struct {
	NetworkSubnets *[]NetworkSubnet `json:"subnets"`
	Meta           *MetaResult      `json:"meta"`
}

// ListNetowrkSubnetsResult structure parses the list network subnets response payload
type ListNetworkSubnetsByNetworkResult struct {
	NetworkSubnets *[]NetworkSubnet `json:"subnets"`
	Meta           *MetaResult      `json:"meta"`
}

type GetNetworkSubnetResult struct {
	NetworkSubnet *NetworkSubnet `json:"subnet"`
}

type CreateNetworkSubnetResult struct {
	Success       bool              `json:"success"`
	Message       string            `json:"msg"`
	Errors        map[string]string `json:"errors"`
	NetworkSubnet *NetworkSubnet    `json:"subnet"`
}

type UpdateNetworkSubnetResult struct {
	CreateNetworkSubnetResult
}

type DeleteNetworkSubnetResult struct {
	DeleteResult
}

// ListNetworkSubnets lists all network subnets
func (client *Client) ListNetworkSubnets(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        NetworkSubnetsPath,
		QueryParams: req.QueryParams,
		Result:      &ListNetworkSubnetsResult{},
	})
}

// ListNetworkSubnetsByNetwork lists all network subnets for a given network
func (client *Client) ListNetworkSubnetsByNetwork(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        fmt.Sprintf("/api/networks/%d/subnets", id),
		QueryParams: req.QueryParams,
		Result:      &ListNetworkSubnetsByNetworkResult{},
	})
}

// GetNetworkSubnet gets an existing network subnet
func (client *Client) GetNetworkSubnet(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        fmt.Sprintf("%s/%d", NetworkSubnetsPath, id),
		QueryParams: req.QueryParams,
		Result:      &GetNetworkSubnetResult{},
	})
}

// CreateNetworkSubnet creates a new network subnet
func (client *Client) CreateNetworkSubnet(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "POST",
		Path:        NetworkSubnetsPath,
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &CreateNetworkSubnetResult{},
	})
}

// UpdateNetworkSubnet updates an existing network subnet
func (client *Client) UpdateNetworkSubnet(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "PUT",
		Path:        fmt.Sprintf("%s/%d", NetworkSubnetsPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &UpdateNetworkSubnetResult{},
	})
}

// DeleteNetworkSubnet deletes an existing network subnet
func (client *Client) DeleteNetworkSubnet(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "DELETE",
		Path:        fmt.Sprintf("%s/%d", NetworkSubnetsPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &DeleteNetworkSubnetResult{},
	})
}

// FindNetworkSubnetByName gets an existing network subnet by name
func (client *Client) FindNetworkSubnetByName(name string) (*Response, error) {
	// Find by name, then get by ID
	resp, err := client.ListNetworkSubnets(&Request{
		QueryParams: map[string]string{
			"name": name,
		},
	})
	if err != nil {
		return resp, err
	}
	listResult := resp.Result.(*ListNetworkSubnetsResult)
	networkSubnetsCount := len(*listResult.NetworkSubnets)
	if networkSubnetsCount != 1 {
		return resp, fmt.Errorf("found %d Network Subnets for %v", networkSubnetsCount, name)
	}
	firstRecord := (*listResult.NetworkSubnets)[0]
	networkSubnetID := firstRecord.ID
	return client.GetNetworkSubnet(networkSubnetID, &Request{})
}
