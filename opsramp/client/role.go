// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package client

import (
	"encoding/json"
	"fmt"
	"net/url"
)

// CreateRole creates a new role
func (c *OpsRampClient) CreateRole(tenantId string, roleData Role) (*Role, error) {
	// Convert Request Data/Body to JSON
	rb, err := json.Marshal(roleData)
	if err != nil {
		return nil, err
	}

	// Prepare the URL, Method and Payload for the Client
	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/roles", c.BaseUrl, tenantId)
	method := "POST"

	// Create a new Request
	body, err := c.NewJsonRequest(method, apiUrl, rb)
	if err != nil {
		return nil, err
	}

	// Preparing Response Body to return and convert it to Golang Map Object
	var responseBody Role
	err = json.Unmarshal(body, &responseBody)
	if err != nil {
		return nil, err
	}

	return &responseBody, nil
}

// GetRole retrieves a specific role by ID
func (c *OpsRampClient) GetRole(tenantId string, roleId string) (*Role, error) {
	// Prepare the URL, Method for the Client
	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/roles/%s", c.BaseUrl, tenantId, roleId)
	method := "GET"

	// Create a new Request
	body, err := c.NewJsonRequest(method, apiUrl, nil)
	if err != nil {
		return nil, err
	}

	// Preparing Response Body
	var responseBody Role
	err = json.Unmarshal(body, &responseBody)
	if err != nil {
		return nil, err
	}

	return &responseBody, nil
}

// UpdateRole updates an existing role
func (c *OpsRampClient) UpdateRole(tenantId string, roleId string, roleData Role) (*Role, error) {
	// Convert Request Data/Body to JSON
	rb, err := json.Marshal(roleData)
	if err != nil {
		return nil, err
	}

	// Prepare the URL, Method and Payload for the Client
	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/roles/%s", c.BaseUrl, tenantId, roleId)
	method := "POST"

	// Create a new Request
	body, err := c.NewJsonRequest(method, apiUrl, rb)
	if err != nil {
		return nil, err
	}

	// Preparing Response Body
	var responseBody Role
	err = json.Unmarshal(body, &responseBody)
	if err != nil {
		return nil, err
	}

	return &responseBody, nil
}

// DeleteRole deletes a role
func (c *OpsRampClient) DeleteRole(tenantId string, roleId string) error {
	// Prepare the URL, Method for the Client
	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/roles/%s", c.BaseUrl, tenantId, roleId)
	method := "DELETE"

	// Create a new Request
	_, err := c.NewJsonRequest(method, apiUrl, nil)

	return err
}

// SearchRoles searches for roles
func (c *OpsRampClient) SearchRoles(tenantId string, queryString string) (*RoleSearchResponse, error) {
	// Prepare the URL, Method for the Client
	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/roles/search", c.BaseUrl, tenantId)
	if queryString != "" {
		apiUrl = fmt.Sprintf("%s?queryString=%s", apiUrl, url.QueryEscape(queryString))
	}
	method := "GET"

	// Create a new Request
	body, err := c.NewJsonRequest(method, apiUrl, nil)
	if err != nil {
		return nil, err
	}

	// Preparing Response Body
	var responseBody RoleSearchResponse
	err = json.Unmarshal(body, &responseBody)
	if err != nil {
		return nil, err
	}

	return &responseBody, nil
}

// FindRoleByName searches for a role by name
func (c *OpsRampClient) FindRoleByName(tenantId string, name string) (*Role, error) {
	response, err := c.SearchRoles(tenantId, fmt.Sprintf("name:%s", name))
	if err != nil {
		return nil, err
	}

	for _, role := range response.Results {
		if role.Name == name {
			return &role, nil
		}
	}

	return nil, fmt.Errorf("role with name '%s' not found", name)
}
