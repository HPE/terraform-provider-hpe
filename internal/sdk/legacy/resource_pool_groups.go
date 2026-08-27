package morpheus

import (
	"fmt"
)

var (
	// ResourcePoolGroupsPath is the API endpoint for resource pool groups
	ResourcePoolGroupsPath = "/api/resource-pools/groups"
)

// ResourcePoolGroup structures for use in request and response payloads
type ResourcePoolGroup struct {
	ID          int64   `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Visibility  string  `json:"visibility"`
	Mode        string  `json:"mode"`
	Pools       []int64 `json:"pools"`
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
	} `json:"resourcePermission"`
}

type ListResourcePoolGroupsResult struct {
	ResourcePoolGroups *[]ResourcePoolGroup `json:"resourcePoolGroups"`
	Meta               *MetaResult          `json:"meta"`
}

type GetResourcePoolGroupResult struct {
	ResourcePoolGroup *ResourcePoolGroup `json:"resourcePoolGroup"`
}

type CreateResourcePoolGroupResult struct {
	Success           bool               `json:"success"`
	Message           string             `json:"msg"`
	Errors            map[string]string  `json:"errors"`
	ResourcePoolGroup *ResourcePoolGroup `json:"resourcePoolGroup"`
}

type UpdateResourcePoolGroupResult struct {
	CreateResourcePoolGroupResult
}

type DeleteResourcePoolGroupResult struct {
	DeleteResult
}

func (client *Client) ListResourcePoolGroups(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        ResourcePoolGroupsPath,
		QueryParams: req.QueryParams,
		Result:      &ListResourcePoolGroupsResult{},
	})
}

func (client *Client) GetResourcePoolGroup(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        fmt.Sprintf("%s/%d", ResourcePoolGroupsPath, id),
		QueryParams: req.QueryParams,
		Result:      &GetResourcePoolGroupResult{},
	})
}

func (client *Client) CreateResourcePoolGroup(req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "POST",
		Path:        ResourcePoolGroupsPath,
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &CreateResourcePoolGroupResult{},
	})
}

func (client *Client) UpdateResourcePoolGroup(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "PUT",
		Path:        fmt.Sprintf("%s/%d", ResourcePoolGroupsPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &UpdateResourcePoolGroupResult{},
	})
}

func (client *Client) DeleteResourcePoolGroup(id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "DELETE",
		Path:        fmt.Sprintf("%s/%d", ResourcePoolGroupsPath, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &DeleteResourcePoolGroupResult{},
	})
}

// FindResourcePoolGroupByName gets an existing resource pool group by name
func (client *Client) FindResourcePoolGroupByName(name string) (*Response, error) {
	// Find by name, then get by ID
	resp, err := client.ListResourcePoolGroups(&Request{
		QueryParams: map[string]string{
			"name": name,
		},
	})
	if err != nil {
		return resp, err
	}
	listResult := resp.Result.(*ListResourcePoolGroupsResult)
	resourcePoolGroupCount := len(*listResult.ResourcePoolGroups)
	if resourcePoolGroupCount != 1 {
		return resp, fmt.Errorf("found %d resource pool groups for %v", resourcePoolGroupCount, name)
	}
	firstRecord := (*listResult.ResourcePoolGroups)[0]
	resourcePoolGroupID := firstRecord.ID
	return client.GetResourcePoolGroup(resourcePoolGroupID, &Request{})
}
