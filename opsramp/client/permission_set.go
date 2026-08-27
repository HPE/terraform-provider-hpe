// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package client

import (
	"encoding/json"
	"fmt"
)

// CreatePermissionSet creates a new permission set
func (c *OpsRampClient) CreatePermissionSet(tenantId string, permSetData CreatePermissionSet) (*PermissionSetResponse, error) {
	// Convert Request Data/Body to JSON
	rb, err := json.Marshal(permSetData)
	if err != nil {
		return nil, err
	}

	// Prepare the URL, Method and Payload for the Client
	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/permissionSets", c.BaseUrl, tenantId)
	method := "POST"

	// Create a new Request
	body, err := c.NewJsonRequest(method, apiUrl, rb)
	if err != nil {
		return nil, err
	}

	// Preparing Response Body to return and convert it to Golang Map Object
	var responseBody PermissionSetResponse
	err = json.Unmarshal(body, &responseBody)
	if err != nil {
		return nil, err
	}

	return &responseBody, nil
}

// GetPermissionSet retrieves a specific permission set by ID
func (c *OpsRampClient) GetPermissionSet(tenantId string, permSetId string) (*PermissionSetResponse, error) {
	// Prepare the URL, Method for the Client
	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/permissionSets/%s", c.BaseUrl, tenantId, permSetId)
	method := "GET"

	// Create a new Request
	body, err := c.NewJsonRequest(method, apiUrl, nil)
	if err != nil {
		return nil, err
	}

	// Preparing Response Body
	var responseBody PermissionSetResponse
	err = json.Unmarshal(body, &responseBody)
	if err != nil {
		return nil, err
	}

	return &responseBody, nil
}

// UpdatePermissionSet updates an existing permission set
func (c *OpsRampClient) UpdatePermissionSet(
	tenantId string,
	permSetId string,
	permSetData UpdatePermissionSet,
) (*PermissionSetResponse, error) {
	// Convert Request Data/Body to JSON
	rb, err := json.Marshal(permSetData)
	if err != nil {
		return nil, err
	}

	// Prepare the URL, Method and Payload for the Client
	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/permissionSets/%s", c.BaseUrl, tenantId, permSetId)
	method := "PUT"

	// Create a new Request
	body, err := c.NewJsonRequest(method, apiUrl, rb)
	if err != nil {
		return nil, err
	}

	// Preparing Response Body
	var responseBody PermissionSetResponse
	err = json.Unmarshal(body, &responseBody)
	if err != nil {
		return nil, err
	}

	return &responseBody, nil
}

// DeletePermissionSet deletes a permission set
func (c *OpsRampClient) DeletePermissionSet(tenantId string, permSetId string) error {
	// Prepare the URL, Method for the Client
	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/permissionSets/%s", c.BaseUrl, tenantId, permSetId)
	method := "DELETE"

	// Create a new Request
	_, err := c.NewJsonRequest(method, apiUrl, nil)

	return err
}

// GetPermissionSets retrieves all permission sets
func (c *OpsRampClient) GetPermissionSets(tenantId string) ([]PermissionSetResponse, error) {
	// Prepare the URL, Method for the Client
	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/permissionSets", c.BaseUrl, tenantId)
	method := "GET"

	// Create a new Request
	body, err := c.NewJsonRequest(method, apiUrl, nil)
	if err != nil {
		return nil, err
	}

	// Preparing Response Body
	var responseBody []PermissionSetResponse
	err = json.Unmarshal(body, &responseBody)
	if err != nil {
		return nil, err
	}

	return responseBody, nil
}

// FindPermissionSetByName searches for a permission set by name
func (c *OpsRampClient) FindPermissionSetByName(tenantId string, name string) (*PermissionSetResponse, error) {
	permSets, err := c.GetPermissionSets(tenantId)
	if err != nil {
		return nil, err
	}

	for _, permSet := range permSets {
		if permSet.Name == name {
			return &permSet, nil
		}
	}

	return nil, fmt.Errorf("permission set with name '%s' not found", name)
}

// SearchPermissionSets searches permission sets using the paginated endpoint (returns numeric IDs)
func (c *OpsRampClient) SearchPermissionSets(tenantId string) (*PermissionSetListResponse, error) {
	// Prepare the URL with pagination parameters
	apiUrl := fmt.Sprintf(
		"%s/api/v2/tenants/%s/permissionSets?sortName=name&isDescendingOrder=false&pageNo=1&pageSize=100",
		c.BaseUrl,
		tenantId,
	)
	method := "GET"

	// Create a new Request
	body, err := c.NewJsonRequest(method, apiUrl, nil)
	if err != nil {
		return nil, err
	}

	// Preparing Response Body
	var responseBody PermissionSetListResponse
	err = json.Unmarshal(body, &responseBody)
	if err != nil {
		return nil, err
	}

	return &responseBody, nil
}

// GetPermissionSetIdByUniqueId looks up the numeric ID of a permission set by its unique ID
func (c *OpsRampClient) GetPermissionSetIdByUniqueId(tenantId string, uniqueId string) (int, error) {
	result, err := c.SearchPermissionSets(tenantId)
	if err != nil {
		return 0, err
	}

	for _, permSet := range result.Results {
		if permSet.UniqueId == uniqueId {
			return permSet.Id, nil
		}
	}

	return 0, fmt.Errorf("permission set with uniqueId '%s' not found", uniqueId)
}
