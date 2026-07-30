// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package client

import (
	"encoding/json"
	"fmt"
)

// CreateUserGroup creates a new user group
func (c *OpsRampClient) CreateUserGroup(tenantId string, groupData CreateUserGroup) (*UserGroupResponse, error) {
	// Convert Request Data/Body to JSON
	rb, err := json.Marshal(groupData)
	if err != nil {
		return nil, err
	}

	// Prepare the URL, Method and Payload for the Client
	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/userGroups", c.BaseUrl, tenantId)
	method := "POST"

	// Create a new Request
	body, err := c.NewJsonRequest(method, apiUrl, rb)
	if err != nil {
		return nil, err
	}

	// Preparing Response Body to return and convert it to Golang Map Object
	var responseBody UserGroupResponse
	err = json.Unmarshal(body, &responseBody)
	if err != nil {
		return nil, err
	}

	return &responseBody, nil
}

// GetUserGroup retrieves a specific user group by ID
func (c *OpsRampClient) GetUserGroup(tenantId string, groupId string) (*UserGroupResponse, error) {
	// Prepare the URL, Method for the Client
	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/userGroups/%s", c.BaseUrl, tenantId, groupId)
	method := "GET"

	// Create a new Request
	body, err := c.NewJsonRequest(method, apiUrl, nil)
	if err != nil {
		return nil, err
	}

	// Preparing Response Body
	var responseBody UserGroupResponse
	err = json.Unmarshal(body, &responseBody)
	if err != nil {
		return nil, err
	}

	return &responseBody, nil
}

// UpdateUserGroup updates an existing user group
func (c *OpsRampClient) UpdateUserGroup(tenantId string, groupId string, groupData UpdateUserGroup) (*UserGroupResponse, error) {
	// Convert Request Data/Body to JSON
	rb, err := json.Marshal(groupData)
	if err != nil {
		return nil, err
	}

	// Prepare the URL, Method and Payload for the Client
	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/userGroups/%s", c.BaseUrl, tenantId, groupId)
	method := "POST"

	// Create a new Request
	body, err := c.NewJsonRequest(method, apiUrl, rb)
	if err != nil {
		return nil, err
	}

	// Preparing Response Body
	var responseBody UserGroupResponse
	err = json.Unmarshal(body, &responseBody)
	if err != nil {
		return nil, err
	}

	return &responseBody, nil
}

// DeleteUserGroup deletes a user group
func (c *OpsRampClient) DeleteUserGroup(tenantId string, groupId string) error {
	// Prepare the URL, Method for the Client
	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/userGroups/%s", c.BaseUrl, tenantId, groupId)
	method := "DELETE"

	// Create a new Request
	_, err := c.NewJsonRequest(method, apiUrl, nil)

	return err
}

// GetUserGroups retrieves all user groups
func (c *OpsRampClient) GetUserGroups(tenantId string) ([]UserGroupResponse, error) {
	// Prepare the URL, Method for the Client
	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/userGroups", c.BaseUrl, tenantId)
	method := "GET"

	// Create a new Request
	body, err := c.NewJsonRequest(method, apiUrl, nil)
	if err != nil {
		return nil, err
	}

	// Preparing Response Body
	var responseBody []UserGroupResponse
	err = json.Unmarshal(body, &responseBody)
	if err != nil {
		return nil, err
	}

	return responseBody, nil
}

// FindUserGroupByName searches for a user group by name
func (c *OpsRampClient) FindUserGroupByName(tenantId string, name string) (*UserGroupResponse, error) {
	groups, err := c.GetUserGroups(tenantId)
	if err != nil {
		return nil, err
	}

	for _, group := range groups {
		if group.Name == name {
			return &group, nil
		}
	}

	return nil, fmt.Errorf("user group with name '%s' not found", name)
}

func (c *OpsRampClient) AddUserToUserGroup(tenantId string, userGroupId string, users []UserRef) error {
	// Convert Request Data/Body to JSON
	rb, err := json.Marshal(users)
	if err != nil {
		return err
	}

	// Prepare the URL, Method and Payload for the Client
	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/userGroups/%s/users", c.BaseUrl, tenantId, userGroupId)
	method := "POST"

	// Create a new Request
	_, err = c.NewJsonRequest(method, apiUrl, rb)
	if err != nil {
		return err
	}

	return nil
}

// GetUserGroupUsers retrieves users in a specific user group
func (c *OpsRampClient) GetUserGroupUsers(tenantId string, groupId string) ([]UserRef, error) {
	// Prepare the URL, Method for the Client
	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/userGroups/%s/users", c.BaseUrl, tenantId, groupId)
	method := "GET"

	// Create a new Request
	body, err := c.NewJsonRequest(method, apiUrl, nil)
	if err != nil {
		return nil, err
	}

	// Preparing Response Body
	var responseBody GetUsersFromUserGroupResponse
	err = json.Unmarshal(body, &responseBody)
	if err != nil {
		return nil, err
	}

	return responseBody.Results, nil
}

// RemoveUsersFromUserGroup removes users from a user group
func (c *OpsRampClient) RemoveUsersFromUserGroup(tenantId string, groupId string, users []UserRef) error {
	// Convert Request Data/Body to JSON
	rb, err := json.Marshal(users)
	if err != nil {
		return err
	}

	// Prepare the URL, Method and Payload for the Client
	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/userGroups/%s/users", c.BaseUrl, tenantId, groupId)
	method := "DELETE"

	// Create a new Request
	_, err = c.NewJsonRequest(method, apiUrl, rb)

	return err
}
