package morpheus

import (
	"fmt"
)

var (
	// NetworkStaticRoutesPath is the API endpoint for network static routes
	NetworkStaticRoutesPath = "/api/networks"
)

// NetworkRoute structures for use in request and response payloads
type NetworkRoute struct {
	ID                int64       `json:"id"`
	Name              string      `json:"name"`
	Code              string      `json:"code"`
	Description       string      `json:"description"`
	Priority          string      `json:"priority"`
	RouteType         string      `json:"routeType"`
	Source            string      `json:"source"`
	SourceType        string      `json:"sourceType"`
	Destination       string      `json:"destination"`
	DestinationType   string      `json:"destinationType"`
	DefaultRoute      bool        `json:"defaultRoute"`
	NetworkMtu        interface{} `json:"networkMtu"`
	ExternalInterface string      `json:"externalInterface"`
	InternalId        string      `json:"internalId"`
	UniqueId          string      `json:"uniqueId"`
	ExternalType      string      `json:"externalType"`
	Enabled           bool        `json:"enabled"`
	Visible           bool        `json:"visible"`
}

// ListNetworkStaticRoutesResult structure parses the list network static routes response payload
type ListNetworkStaticRoutesResult struct {
	NetworkRoutes *[]NetworkRoute `json:"networkRoutes"`
	Meta          *MetaResult     `json:"meta"`
}

type GetNetworkStaticRouteResult struct {
	NetworkRoute *NetworkRoute `json:"networkRoute"`
}

type CreateNetworkStaticRouteResult struct {
	Success      bool              `json:"success"`
	Message      string            `json:"msg"`
	Errors       map[string]string `json:"errors"`
	NetworkRoute *NetworkRoute     `json:"networkRoute"`
}

type UpdateNetworkStaticRouteResult struct {
	CreateNetworkStaticRouteResult
}

type DeleteNetworkStaticRouteResult struct {
	DeleteResult
}

// ListNetworkStaticRoutes lists all network static routes
func (client *Client) ListNetworkStaticRoutes(networkId int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        fmt.Sprintf("%s/%d/routes", NetworkStaticRoutesPath, networkId),
		QueryParams: req.QueryParams,
		Result:      &ListNetworkStaticRoutesResult{},
	})
}

// GetNetworkStaticRoute Server gets an existing network static route
func (client *Client) GetNetworkStaticRoute(networkId int64, routeId int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        fmt.Sprintf("%s/%d/routes/%d", NetworkStaticRoutesPath, networkId, routeId),
		QueryParams: req.QueryParams,
		Result:      &GetNetworkStaticRouteResult{},
	})
}

// CreateNetworkStaticRoute creates a new network static route
func (client *Client) CreateNetworkStaticRoute(networkId int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "POST",
		Path:        fmt.Sprintf("%s/%d/routes", NetworkStaticRoutesPath, networkId),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &CreateNetworkStaticRouteResult{},
	})
}

// UpdateNetworkStaticRoute updates an existing network static route
func (client *Client) UpdateNetworkStaticRoute(networkId int64, routeId int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "PUT",
		Path:        fmt.Sprintf("%s/%d/routes/%d", NetworkStaticRoutesPath, networkId, routeId),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &UpdateNetworkStaticRouteResult{},
	})
}

// DeleteNetworkStaticRoute deletes an existing network static route
func (client *Client) DeleteNetworkStaticRoute(networkId int64, routeId int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "DELETE",
		Path:        fmt.Sprintf("%s/%d/routes/%d", NetworkStaticRoutesPath, networkId, routeId),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &DeleteNetworkStaticRouteResult{},
	})
}

// FindNetworkStaticRouteByName gets an existing network static route by name
func (client *Client) FindNetworkStaticRouteByName(networkId int64, name string) (*Response, error) {
	// Find by name, then get by ID
	resp, err := client.ListNetworkStaticRoutes(networkId, &Request{
		QueryParams: map[string]string{
			"name": name,
		},
	})
	if err != nil {
		return resp, err
	}
	listResult := resp.Result.(*ListNetworkStaticRoutesResult)
	networkStaticRoutesCount := len(*listResult.NetworkRoutes)
	if networkStaticRoutesCount != 1 {
		return resp, fmt.Errorf("found %d Network static routes for %v", networkStaticRoutesCount, name)
	}
	firstRecord := (*listResult.NetworkRoutes)[0]
	networkRouteID := firstRecord.ID
	return client.GetNetworkStaticRoute(networkId, networkRouteID, &Request{})
}
