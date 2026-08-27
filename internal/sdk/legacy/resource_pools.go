package morpheus

import (
	"fmt"
)

// ResourcePool structures for use in request and response payloads
type ResourcePool struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Zone        struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	} `json:"zone"`
	Active      bool        `json:"active"`
	Type        string      `json:"type"`
	ExternalId  string      `json:"externalId"`
	ReadOnly    bool        `json:"readOnly"`
	DefaultPool bool        `json:"defaultPool"`
	RegionCode  interface{} `json:"regionCode"`
	IacId       interface{} `json:"iacId"`
	Status      string      `json:"status"`
	Inventory   bool        `json:"inventory"`
	Visibility  string      `json:"visibility"`
	Config      struct {
		CidrBlock string `json:"cidrBlock"`
		Tenancy   string `json:"tenancy"`
	} `json:"config"`
	Tenants []struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	} `json:"tenants"`
	ResourcePermission struct {
		All      string        `json:"all"`
		Sites    []interface{} `json:"sites"`
		Plans    []interface{} `json:"plans"`
		AllPlans string        `json:"allPlans"`
	} `json:"resourcePermission"`
}

type ListResourcePoolsResult struct {
	ResourcePools *[]ResourcePool `json:"resourcePools"`
	Meta          *MetaResult     `json:"meta"`
}

type GetResourcePoolResult struct {
	ResourcePool *ResourcePool `json:"resourcePool"`
}

type CreateResourcePoolResult struct {
	Success      bool              `json:"success"`
	Message      string            `json:"msg"`
	Errors       map[string]string `json:"errors"`
	ResourcePool *ResourcePool     `json:"resourcePool"`
}

type UpdateResourcePoolResult struct {
	CreateResourcePoolResult
}

type DeleteResourcePoolResult struct {
	DeleteResult
}

type ResourcePoolPayload struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Active      bool   `json:"active"`
	Visibility  string `json:"visibility"`
}

type CreateResourcePoolPayload struct {
	ResourcePoolPayload *ResourcePoolPayload `json:"resourcePool"`
}

type UpdateResourcePoolPayload struct {
	CreateResourcePoolPayload
}

func (client *Client) ListResourcePools(cloudId int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        fmt.Sprintf("/api/zones/%d/resource-pools", cloudId),
		QueryParams: req.QueryParams,
		Result:      &ListResourcePoolsResult{},
	})
}

func (client *Client) GetResourcePool(cloudId int64, id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "GET",
		Path:        fmt.Sprintf("/api/zones/%d/resource-pools/%d", cloudId, id),
		QueryParams: req.QueryParams,
		Result:      &GetResourcePoolResult{},
	})
}

func (client *Client) CreateResourcePool(cloudId int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "POST",
		Path:        fmt.Sprintf("/api/zones/%d/resource-pools", cloudId),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &CreateResourcePoolResult{},
	})
}

func (client *Client) UpdateResourcePool(cloudId int64, id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "PUT",
		Path:        fmt.Sprintf("/api/zones/%d/resource-pools/%d", cloudId, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &UpdateResourcePoolResult{},
	})
}

func (client *Client) DeleteResourcePool(cloudId int64, id int64, req *Request) (*Response, error) {
	return client.Execute(&Request{
		Method:      "DELETE",
		Path:        fmt.Sprintf("/api/zones/%d/resource-pools/%d", cloudId, id),
		QueryParams: req.QueryParams,
		Body:        req.Body,
		Result:      &DeleteResourcePoolResult{},
	})
}

// FindResourcePoolByName gets an existing resource pool by name
func (client *Client) FindResourcePoolByName(cloudId int64, name string) (*Response, error) {
	// Find by name, then get by ID
	resp, err := client.ListResourcePools(cloudId, &Request{
		QueryParams: map[string]string{
			"name": name,
		},
	})
	if err != nil {
		return resp, err
	}
	listResult := resp.Result.(*ListResourcePoolsResult)
	// Go through listResult and find any entries with the matching name. Check the length of this new list
	// if the list has more than one entry return an error
	var matchingPools []ResourcePool
	for _, pool := range *listResult.ResourcePools {
		if pool.Name == name {
			matchingPools = append(matchingPools, pool)
		}
	}
	resourcePoolCount := len(matchingPools)
	if resourcePoolCount != 1 {
		return resp, fmt.Errorf("found %d resourcePools for %v", resourcePoolCount, name)
	}
	firstRecord := matchingPools[0]
	resourcePoolID := firstRecord.ID
	return client.GetResourcePool(cloudId, resourcePoolID, &Request{})
}
