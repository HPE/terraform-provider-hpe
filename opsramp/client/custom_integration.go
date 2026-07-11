// SPDX-FileCopyrightText: Copyright Hewlett Packard Enterprise Development LP
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"encoding/json"
	"fmt"
)

// CreateCustomIntegrationAPI creates a new custom integration to get API OAuth2 credentials
// tenantId can be either the provider's tenant ID or a client ID for client-level integrations
func (c *OpsRampClient) CreateCustomIntegration(tenantId string, integrationData CreateCustomIntegration) (*CustomIntegrationResponse, error) {
	// Convert Request Data/Body to JSON
	rb, err := json.Marshal(integrationData)
	if err != nil {
		return nil, err
	}

	// Prepare the URL, Method and Payload for the Client
	// API endpoint: /api/v2/tenants/{tenantId}/integrations/install/CUSTOM
	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/integrations/install/CUSTOM", c.BaseUrl, tenantId)
	method := "POST"

	// Create a new Request
	body, err := c.NewJsonRequest(method, apiUrl, rb)
	if err != nil {
		return nil, err
	}

	// Preparing Response Body to return and convert it to Golang Map Object
	var responseBody CustomIntegrationResponse
	err = json.Unmarshal([]byte(body), &responseBody)
	if err != nil {
		return nil, err
	}

	return &responseBody, nil
}

// GetCustomIntegration retrieves a specific custom integration by ID
func (c *OpsRampClient) GetCustomIntegration(tenantId string, integrationId string) (*CustomIntegrationResponse, error) {
	// Prepare the URL, Method for the Client
	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/integrations/installed/%s", c.BaseUrl, tenantId, integrationId)
	method := "GET"

	// Create a new Request
	body, err := c.NewJsonRequest(method, apiUrl, nil)
	if err != nil {
		return nil, err
	}

	// Preparing Response Body
	var responseBody CustomIntegrationResponse
	err = json.Unmarshal([]byte(body), &responseBody)
	if err != nil {
		return nil, err
	}

	return &responseBody, nil
}

// UpdateCustomIntegration updates an existing custom integration
func (c *OpsRampClient) UpdateCustomIntegration(tenantId string, integrationId string, integrationData CreateCustomIntegration) (*CustomIntegrationResponse, error) {
	// Convert Request Data/Body to JSON
	rb, err := json.Marshal(integrationData)
	if err != nil {
		return nil, err
	}

	// Prepare the URL, Method and Payload for the Client
	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/integrations/installed/%s", c.BaseUrl, tenantId, integrationId)
	method := "POST"

	// Create a new Request
	body, err := c.NewJsonRequest(method, apiUrl, rb)
	if err != nil {
		return nil, err
	}

	// Preparing Response Body
	var responseBody CustomIntegrationResponse
	err = json.Unmarshal([]byte(body), &responseBody)
	if err != nil {
		return nil, err
	}

	return &responseBody, nil
}

// DeleteCustomIntegration deletes a custom integration
func (c *OpsRampClient) DeleteCustomIntegration(tenantId string, integrationId string) error {
	// Prepare the URL, Method for the Client
	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/integrations/installed/%s", c.BaseUrl, tenantId, integrationId)
	method := "DELETE"

	deleteIntegrationData := DeleteCustomIntegration{
		UninstallReason: "Terraform - No longer needed",
	}

	rb, err := json.Marshal(deleteIntegrationData)
	if err != nil {
		return err
	}

	// Create a new Request
	_, err = c.NewJsonRequest(method, apiUrl, rb)
	return err
}

// ListCustomIntegrations lists all custom integrations for a tenant
func (c *OpsRampClient) ListCustomIntegrations(tenantId string) (*CustomIntegrationListResponse, error) {
	// Prepare the URL, Method for the Client
	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/integrations/installed/search?integrationType=CUSTOM", c.BaseUrl, tenantId)
	method := "GET"

	// Create a new Request
	body, err := c.NewJsonRequest(method, apiUrl, nil)
	if err != nil {
		return nil, err
	}

	// Preparing Response Body
	var responseBody CustomIntegrationListResponse
	err = json.Unmarshal([]byte(body), &responseBody)
	if err != nil {
		return nil, err
	}

	return &responseBody, nil
}

// RegenerateCustomIntegrationToken regenerates the OAuth2 token for a custom integration
// Returns the new AuthenticationConfig with fresh apiKeyPairs
func (c *OpsRampClient) RegenerateCustomIntegrationToken(tenantId string, integrationId string) (*AuthenticationConfig, error) {
	// Prepare the URL, Method for the Client
	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/integrations/installed/%s/inbound/authentication/regenerateSecretOrToken", c.BaseUrl, tenantId, integrationId)
	method := "POST"

	// Create a new Request (no body needed)
	body, err := c.NewJsonRequest(method, apiUrl, nil)
	if err != nil {
		return nil, err
	}

	// Preparing Response Body - returns AuthenticationConfig
	var responseBody AuthenticationConfig
	err = json.Unmarshal([]byte(body), &responseBody)
	if err != nil {
		return nil, err
	}

	return &responseBody, nil
}
