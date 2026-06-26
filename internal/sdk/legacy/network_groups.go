package morpheus

import (
	"fmt"
)

var (
	// NetworkGroupsPath is the API endpoint for network groups
	NetworkGroupsPath = "/api/networks/groups"
)

// NetworkGroup structures for use in request and response payloads
type NetworkGroup struct {
	ID          int64   `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Visibility  string  `json:"visibility"`
	Active      bool    `json:"active"`
	Networks    []int64 `json:"networks"`
	Subnets     []int64 `json:"subnets"`
	Tenants     []struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	}
	ResourcePermission struct {
		All   bool `json:"all"`
		Sites []struct {
			ID      int64  `json:"id"`
			Name    string `json:"name"`
			Default bool   `json:"default"`
		} `json:"sites"`
		AllPlans bool `json:"allPlans"`
	} `json:"resourcePermission"`
}

// ListNetworkGroupsResult structure parses the list network groups response payload
type ListNetworkGroupsResult struct {
	NetworkGroups *[]NetworkGroup `json:"networkGroups"`
	Meta          *MetaResult     `json:"meta"`
}

type GetNetworkGroupResult struct {
	NetworkGroup *NetworkGroup `json:"networkGroup"`
}

type CreateNetworkGroupResult struct {
	Success      bool              `json:"success"`
	Message      string            `json:"msg"`
	Errors       map[string]string `json:"errors"`
	NetworkGroup *NetworkGroup     `json:"networkGroup"`
	IsOwner      bool              `json:"isOwner"`
}

type UpdateNetworkGroupResult struct {
	CreateNetworkGroupResult
}

type DeleteNetworkGroupResult struct {
	DeleteResult
}

// ListNetworkGroups lists all network groups
func (client *Client) ListNetworkGroups(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        NetworkGroupsPath,
		QueryParams: req.QueryParams,
		Result:      &ListNetworkGroupsResult{},
	})
}

// GetNetworkGroup gets an existing network group
func (client *Client) GetNetworkGroup(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        fmt.Sprintf("%s/%d", NetworkGroupsPath, id),
		QueryParams: req.QueryParams,
		Result:      &GetNetworkGroupResult{},
	})
}

// CreateNetworkGroup creates a new network group
func (client *Client) CreateNetworkGroup(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "POST",
		Path:        NetworkGroupsPath,
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &CreateNetworkGroupResult{},
	})
}

// UpdateNetworkGroup updates an existing network group
func (client *Client) UpdateNetworkGroup(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "PUT",
		Path:        fmt.Sprintf("%s/%d", NetworkGroupsPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &UpdateNetworkGroupResult{},
	})
}

// DeleteNetworkGroup deletes an existing network group
func (client *Client) DeleteNetworkGroup(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "DELETE",
		Path:        fmt.Sprintf("%s/%d", NetworkGroupsPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &DeleteNetworkGroupResult{},
	})
}

// FindNetworkGroupByName gets an existing network group by name
func (client *Client) FindNetworkGroupByName(name string) (*Response, error) {
	// Find by name, then get by ID
	resp, err := client.ListNetworkGroups(&Request{
		QueryParams: map[string]string{
			"name": name,
		},
	})
	if err != nil {
		return resp, err
	}
	listResult := resp.Result.(*ListNetworkGroupsResult)
	networkGroupsCount := len(*listResult.NetworkGroups)
	if networkGroupsCount != 1 {
		return resp, fmt.Errorf("found %d Network Groups for %v", networkGroupsCount, name)
	}
	firstRecord := (*listResult.NetworkGroups)[0]
	networkGroupID := firstRecord.ID
	return client.GetNetworkGroup(networkGroupID, &Request{})
}
