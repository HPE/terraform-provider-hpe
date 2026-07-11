// SPDX-FileCopyrightText: Copyright Hewlett Packard Enterprise Development LP
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"encoding/json"
	"fmt"
)

// CreateClient creates a new client (sub-tenant)
func (c *OpsRampClient) CreateClient(clientData CreateClient) (*ClientMinimal, error) {
	// Convert Request Data/Body to JSON
	rb, err := json.Marshal(clientData)
	if err != nil {
		return nil, err
	}

	// Prepare the URL, Method and Payload for the Client
	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/clients", c.BaseUrl, c.TenantId)
	method := "POST"

	// Create a new Request
	body, err := c.NewJsonRequest(method, apiUrl, rb)
	if err != nil {
		return nil, err
	}

	// Preparing Response Body to return and convert it to Golang Map Object
	var responseBody ClientMinimal
	err = json.Unmarshal([]byte(body), &responseBody)
	if err != nil {
		return nil, err
	}

	return &responseBody, nil
}

// GetClients retrieves all clients (minimal info)
func (c *OpsRampClient) GetClients() ([]ClientMinimal, error) {
	// Prepare the URL, Method for the Client
	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/clients/minimal", c.BaseUrl, c.TenantId)
	method := "GET"

	// Create a new Request
	body, err := c.NewJsonRequest(method, apiUrl, nil)
	if err != nil {
		return nil, err
	}

	// Preparing Response Body to return and convert it to Golang slice
	var responseBody []ClientMinimal
	err = json.Unmarshal([]byte(body), &responseBody)
	if err != nil {
		return nil, err
	}

	return responseBody, nil
}

// GetClient retrieves a specific client by ID
func (c *OpsRampClient) GetClient(clientId string) (*ClientResponse, error) {
	// Prepare the URL, Method for the Client
	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/clients/%s", c.BaseUrl, c.TenantId, clientId)
	method := "GET"

	// Create a new Request
	body, err := c.NewJsonRequest(method, apiUrl, nil)
	if err != nil {
		return nil, err
	}

	// Preparing Response Body
	var responseBody ClientResponse
	err = json.Unmarshal([]byte(body), &responseBody)
	if err != nil {
		return nil, err
	}

	return &responseBody, nil
}

// UpdateClient updates an existing client
func (c *OpsRampClient) UpdateClient(clientId string, clientData CreateClient) (*ClientMinimal, error) {
	// Convert Request Data/Body to JSON
	rb, err := json.Marshal(clientData)
	if err != nil {
		return nil, err
	}

	// Prepare the URL, Method and Payload for the Client
	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/clients/%s", c.BaseUrl, c.TenantId, clientId)
	method := "POST"

	// Create a new Request
	body, err := c.NewJsonRequest(method, apiUrl, rb)
	if err != nil {
		return nil, err
	}

	// Preparing Response Body
	var responseBody ClientMinimal
	err = json.Unmarshal([]byte(body), &responseBody)
	if err != nil {
		return nil, err
	}

	return &responseBody, nil
}

// DeleteClient deletes a client
func (c *OpsRampClient) DeleteClient(clientId string) error {
	// Prepare the URL, Method for the Client
	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/clients/%s/terminate", c.BaseUrl, c.TenantId, clientId)
	method := "POST"

	// Create a new Request
	_, err := c.NewJsonRequest(method, apiUrl, nil)
	return err
}

// FindClientByName searches for a client by name
func (c *OpsRampClient) FindClientByName(name string) (*ClientMinimal, error) {
	clients, err := c.GetClients()
	if err != nil {
		return nil, err
	}

	for _, client := range clients {
		if client.Name == name {
			return &client, nil
		}
	}

	return nil, fmt.Errorf("client with name '%s' not found", name)
}
