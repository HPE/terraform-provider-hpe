// SPDX-FileCopyrightText: Copyright Hewlett Packard Enterprise Development LP
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"encoding/json"
	"fmt"
)

// CreateUser creates a new user
func (c *OpsRampClient) CreateUser(tenantId string, userData CreateUser) (*UserResponse, error) {
	// Convert Request Data/Body to JSON
	rb, err := json.Marshal(userData)
	if err != nil {
		return nil, err
	}

	// Prepare the URL, Method and Payload for the Client
	// Use provided tenantId to support both partner and client-level users
	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/users", c.BaseUrl, tenantId)
	method := "POST"

	// Create a new Request
	body, err := c.NewJsonRequest(method, apiUrl, rb)
	if err != nil {
		return nil, err
	}

	// Preparing Response Body to return and convert it to Golang Map Object
	var responseBody UserResponse
	err = json.Unmarshal([]byte(body), &responseBody)
	if err != nil {
		return nil, err
	}

	return &responseBody, nil
}

// GetUser retrieves a specific user by ID
func (c *OpsRampClient) GetUser(tenantId string, userId string) (*UserResponse, error) {
	// Prepare the URL, Method for the Client
	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/users/%s", c.BaseUrl, tenantId, userId)
	method := "GET"

	// Create a new Request
	body, err := c.NewJsonRequest(method, apiUrl, nil)
	if err != nil {
		return nil, err
	}

	// Preparing Response Body
	var responseBody UserResponse
	err = json.Unmarshal([]byte(body), &responseBody)
	if err != nil {
		return nil, err
	}

	return &responseBody, nil
}

// UpdateUser updates an existing user
func (c *OpsRampClient) UpdateUser(tenantId string, userId string, userData UpdateUser) (*UserResponse, error) {
	// Convert Request Data/Body to JSON
	rb, err := json.Marshal(userData)
	if err != nil {
		return nil, err
	}

	// Prepare the URL, Method and Payload for the Client
	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/users/%s", c.BaseUrl, tenantId, userId)
	method := "POST"

	// Create a new Request
	body, err := c.NewJsonRequest(method, apiUrl, rb)
	if err != nil {
		return nil, err
	}

	// Preparing Response Body
	var responseBody UserResponse

	if len(body) == 0 {
		// If the response body is empty, look for the user resource
		response, err := c.GetUser(tenantId, userId)
		if err != nil {
			return nil, fmt.Errorf("user updated but failed to retrieve updated user: %w", err)
		}

		responseBody = *response
	} else {
		err = json.Unmarshal([]byte(body), &responseBody)
		if err != nil {
			return nil, err
		}
	}

	return &responseBody, nil
}

// DeleteUser deletes an existing user
func (c *OpsRampClient) DeleteUser(tenantId string, userId string) error {

	data := DeleteUserRequest{
		TerminateReason: "Deleted via API",
		MaskType:        "FULL",
	}

	// Convert Request Data/Body to JSON
	rb, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// Prepare the URL, Method and Payload for the Client
	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/users/%s/TERMINATE", c.BaseUrl, tenantId, userId)
	method := "POST"

	// Create a new Request
	_, err = c.NewJsonRequest(method, apiUrl, rb)
	if err != nil {
		return err
	}

	return nil
}

// SearchUsers searches for users
func (c *OpsRampClient) SearchUsers(tenantId string, queryString string) (*UserSearchResponse, error) {
	// Prepare the URL, Method for the Client
	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/users/search", c.BaseUrl, tenantId)
	if queryString != "" {
		apiUrl = fmt.Sprintf("%s?queryString=%s", apiUrl, queryString)
	}
	method := "GET"

	// Create a new Request
	body, err := c.NewJsonRequest(method, apiUrl, nil)
	if err != nil {
		return nil, err
	}

	// Preparing Response Body
	var responseBody UserSearchResponse
	err = json.Unmarshal([]byte(body), &responseBody)
	if err != nil {
		return nil, err
	}

	return &responseBody, nil
}

// ChangeUserPassword changes a user's password
func (c *OpsRampClient) ChangeUserPassword(tenantId string, userId string, passwordData ChangePasswordRequest) error {
	// Convert Request Data/Body to JSON
	rb, err := json.Marshal(passwordData)
	if err != nil {
		return err
	}

	// Prepare the URL, Method and Payload for the Client
	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/users/%s/changePassword", c.BaseUrl, tenantId, userId)
	method := "POST"

	// Create a new Request
	_, err = c.NewJsonRequest(method, apiUrl, rb)
	return err
}

// FindUserByLoginName searches for a user by login name
func (c *OpsRampClient) FindUserByLoginName(tenantId string, loginName string) (*UserResponse, error) {
	response, err := c.SearchUsers(tenantId, fmt.Sprintf("loginName:%s", loginName))
	if err != nil {
		return nil, err
	}

	for _, user := range response.Results {
		if user.LoginName == loginName {
			return &user, nil
		}
	}

	return nil, fmt.Errorf("user with login name '%s' not found", loginName)
}

// GetUsersMinimal retrieves a minimal list of users
func (c *OpsRampClient) GetUsersMinimal(tenantId string) ([]UserMinimalResponse, error) {
	// Prepare the URL, Method for the Client
	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/users/minimal", c.BaseUrl, tenantId)
	method := "GET"

	// Create a new Request
	body, err := c.NewJsonRequest(method, apiUrl, nil)
	if err != nil {
		return nil, err
	}

	// Preparing Response Body
	var responseBody []UserMinimalResponse
	err = json.Unmarshal([]byte(body), &responseBody)
	if err != nil {
		return nil, err
	}

	return responseBody, nil
}
