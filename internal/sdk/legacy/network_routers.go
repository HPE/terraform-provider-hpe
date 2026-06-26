package morpheus

import (
	"fmt"
)

var (
	// NetworkRoutersPath is the API endpoint for network routers
	NetworkRoutersPath     = "/api/networks/routers"
	NetworkRouterTypesPath = "/api/network-router-types"
)

// NetworkRouter structures for use in request and response payloads
type NetworkRouter struct {
	ID            int64             `json:"id"`
	Name          string            `json:"name"`
	Code          string            `json:"code"`
	Description   string            `json:"description"`
	Category      string            `json:"category"`
	DateCreated   string            `json:"dateCreated"`
	LastUpdated   string            `json:"lastUpdated"`
	RouterType    string            `json:"routerType"`
	Status        string            `json:"status"`
	Enabled       bool              `json:"enabled"`
	EnableBgp     bool              `json:"enableBgp"`
	ExternalIp    string            `json:"externalIp"`
	ExternalId    string            `json:"externalId"`
	ProviderId    string            `json:"providerId"`
	Type          NetworkRouterType `json:"type"`
	NetworkServer struct {
		ID          int64  `json:"id"`
		Name        string `json:"name"`
		Integration struct {
			ID              int64  `json:"id"`
			Name            string `json:"name"`
			Enabled         bool   `json:"enabled"`
			Type            string `json:"type"`
			IntegrationType struct {
				ID   int64  `json:"id"`
				Name string `json:"name"`
				Code string `json:"code"`
			} `json:"integrationType"`
			URL              string `json:"url"`
			Port             string `json:"port"`
			Username         string `json:"username"`
			Password         string `json:"password"`
			PasswordHash     string `json:"passwordHash"`
			RefType          string `json:"refType"`
			RefId            string `json:"refId"`
			IsPlugin         bool   `json:"isPlugin"`
			Status           string `json:"status"`
			StatusDate       string `json:"statusDate"`
			StatusMessage    string `json:"statusMessage"`
			LastSync         string `json:"lastSync"`
			LastSyncDuration int64  `json:"lastSyncDuration"`
		} `json:"integration"`
	} `json:"networkServer"`
	Interfaces []struct {
		ID              int64       `json:"id"`
		Name            string      `json:"name"`
		Code            string      `json:"code"`
		InterfaceType   string      `json:"interfaceType"`
		NetworkPosition interface{} `json:"networkPosition"`
		IpAddress       string      `json:"ipAddress"`
		Cidr            string      `json:"cidr"`
		ExternalLink    string      `json:"externalLink"`
		Enabled         bool        `json:"enabled"`
		Network         struct {
			ID   int64  `json:"id"`
			Name string `json:"name"`
			Code string `json:"code"`
		} `json:"network"`
	}
	Zone struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
		Code string `json:"code"`
	} `json:"zone"`
	ExternalNetwork struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
		Code string `json:"code"`
	} `json:"externalNetwork"`
	Permissions struct {
		Visibility       string `json:"visibility"`
		TenantPermission struct {
			Accounts []int64 `json:"accounts"`
		} `json:"tenantPermissions"`
	} `json:"permissions"`
}

type NetworkRouterType struct {
	ID                       int64        `json:"id"`
	Name                     string       `json:"name"`
	Code                     string       `json:"code"`
	Description              string       `json:"description"`
	Enabled                  bool         `json:"enabled"`
	Creatable                bool         `json:"creatable"`
	Selectable               bool         `json:"selectable"`
	HasFirewall              bool         `json:"hasFirewall"`
	HasDhcp                  bool         `json:"hasDhcp"`
	HasRouting               bool         `json:"hasRouting"`
	HasNetworkServer         bool         `json:"hasNetworkServer"`
	HasNat                   bool         `json:"hasNat"`
	HasBgp                   bool         `json:"hasBgp"`
	HasFirewallGroups        bool         `json:"hasFirewallGroups"`
	HasSecurityGroupPriority bool         `json:"hasSecurityGroupPriority"`
	OptionTypes              []OptionType `json:"optionTypes"`
	RuleOptionTypes          []OptionType `json:"ruleOptionTypes"`
	NatOptionTypes           []OptionType `json:"natOptionTypes"`
	BgpNeighborOptionTypes   []OptionType `json:"bgpNeighborOptionTypes"`
}

// ListNetworkRoutersResult structure parses the list network routers response payload
type ListNetworkRoutersResult struct {
	NetworkRouters *[]NetworkRouter `json:"networkRouters"`
	Meta           *MetaResult      `json:"meta"`
}

type GetNetworkRouterResult struct {
	NetworkRouter *NetworkRouter `json:"networkRouter"`
}

type CreateNetworkRouterResult struct {
	Success       bool              `json:"success"`
	Message       string            `json:"msg"`
	Errors        map[string]string `json:"errors"`
	NetworkRouter *NetworkRouter    `json:"networkRouter"`
}

type UpdateNetworkRouterResult struct {
	CreateNetworkRouterResult
}

type DeleteNetworkRouterResult struct {
	DeleteResult
}

type ListNetworkRouterTypesResult struct {
	NetworkRouterTypes *[]NetworkRouterType `json:"networkRouterTypes"`
	Meta               *MetaResult          `json:"meta"`
}

type GetNetworkRouterTypeResult struct {
	NetworkRouterType *NetworkRouterType `json:"networkRouterType"`
}

// ListNetworkRouters lists all network routers
func (client *Client) ListNetworkRouters(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        NetworkRoutersPath,
		QueryParams: req.QueryParams,
		Result:      &ListNetworkRoutersResult{},
	})
}

// GetNetworkRouter gets an existing network router
func (client *Client) GetNetworkRouter(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        fmt.Sprintf("%s/%d", NetworkRoutersPath, id),
		QueryParams: req.QueryParams,
		Result:      &GetNetworkRouterResult{},
	})
}

// CreateNetworkRouter creates a new network router
func (client *Client) CreateNetworkRouter(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "POST",
		Path:        NetworkRoutersPath,
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &CreateNetworkRouterResult{},
	})
}

// UpdateNetworkRouter updates an existing network router
func (client *Client) UpdateNetworkRouter(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "PUT",
		Path:        fmt.Sprintf("%s/%d", NetworkRoutersPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &UpdateNetworkRouterResult{},
	})
}

// DeleteNetworkRouter deletes an existing network router
func (client *Client) DeleteNetworkRouter(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "DELETE",
		Path:        fmt.Sprintf("%s/%d", NetworkRoutersPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &DeleteNetworkRouterResult{},
	})
}

// FindNetworkRouterByName gets an existing network router by name
func (client *Client) FindNetworkRouterByName(name string) (*Response, error) {
	// Find by name, then get by ID
	resp, err := client.ListNetworkRouters(&Request{
		QueryParams: map[string]string{
			"name": name,
		},
	})
	if err != nil {
		return resp, err
	}
	listResult := resp.Result.(*ListNetworkRoutersResult)
	networkRoutersCount := len(*listResult.NetworkRouters)
	if networkRoutersCount != 1 {
		return resp, fmt.Errorf("found %d Network Routers for %v", networkRoutersCount, name)
	}
	firstRecord := (*listResult.NetworkRouters)[0]
	networkRouterID := firstRecord.ID
	return client.GetNetworkRouter(networkRouterID, &Request{})
}

// Router Types

// ListNetworkRouterTypes fetches existing network router types
func (client *Client) ListNetworkRouterTypes(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        NetworkRouterTypesPath,
		QueryParams: req.QueryParams,
		Result:      &ListNetworkRouterTypesResult{},
	})
}

// GetNetworkRouterType fetches an existing network router type
func (client *Client) GetNetworkRouterType(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        fmt.Sprintf("%s/%d", NetworkRouterTypesPath, id),
		QueryParams: req.QueryParams,
		Result:      &GetNetworkRouterTypeResult{},
	})
}

func (client *Client) FindNetworkRouterTypeByName(name string) (*Response, error) {
	// Find by name, then get by ID
	resp, err := client.ListNetworkRouterTypes(&Request{
		QueryParams: map[string]string{
			"name": name,
		},
	})
	if err != nil {
		return resp, err
	}
	listResult := resp.Result.(*ListNetworkRouterTypesResult)
	networkRouterTypesCount := len(*listResult.NetworkRouterTypes)
	if networkRouterTypesCount != 1 {
		return resp, fmt.Errorf("found %d Network Routers Types for %v", networkRouterTypesCount, name)
	}
	firstRecord := (*listResult.NetworkRouterTypes)[0]
	networkRouterTypeId := firstRecord.ID
	return client.GetNetworkRouterType(networkRouterTypeId, &Request{})
}
