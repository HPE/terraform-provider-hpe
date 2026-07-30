// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package client

import (
	"encoding/json"
	"fmt"
)

// CreateIntegrationConfig creates a new config under an installed integration.
func (c *OpsRampClient) CreateIntegrationConfig(tenantId string, integrationId string, request IntegrationConfigRequest) (*IntegrationConfigResponse, error) {
	rb, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	apiUrl := fmt.Sprintf("%s/api/v3/tenants/%s/apps/install/%s/config", c.BaseUrl, tenantId, integrationId)

	body, err := c.NewJsonRequest("POST", apiUrl, rb)
	if err != nil {
		return nil, err
	}

	var response IntegrationConfigResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// GetIntegrationConfig retrieves a specific config by its ID.
func (c *OpsRampClient) GetIntegrationConfig(tenantId string, integrationId string, configId string) (*IntegrationConfigResponse, error) {
	apiUrl := fmt.Sprintf("%s/api/v3/tenants/%s/apps/install/%s/config/%s", c.BaseUrl, tenantId, integrationId, configId)

	body, err := c.NewJsonRequest("GET", apiUrl, nil)
	if err != nil {
		return nil, err
	}

	var response IntegrationConfigResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// UpdateIntegrationConfig updates an existing config.
func (c *OpsRampClient) UpdateIntegrationConfig(tenantId string, integrationId string, configId string, request IntegrationConfigRequest) (*IntegrationConfigResponse, error) {
	rb, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	apiUrl := fmt.Sprintf("%s/api/v3/tenants/%s/apps/install/%s/config/%s", c.BaseUrl, tenantId, integrationId, configId)

	body, err := c.NewJsonRequest("POST", apiUrl, rb)
	if err != nil {
		return nil, err
	}

	var response IntegrationConfigResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// DeleteIntegrationConfig deletes a config by its ID.
func (c *OpsRampClient) DeleteIntegrationConfig(tenantId string, integrationId string, configId string) error {
	apiUrl := fmt.Sprintf("%s/api/v3/tenants/%s/apps/install/%s/config/%s?keepAgentInstalledResources=false", c.BaseUrl, tenantId, integrationId, configId)

	_, err := c.NewJsonRequest("DELETE", apiUrl, nil)

	return err
}
