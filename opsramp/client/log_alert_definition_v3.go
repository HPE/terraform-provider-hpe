// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package client

import (
	"encoding/json"
	"fmt"
)

// --- Log Alert Definition CRUD ---

// CreateLogAlertDefinition creates one or more log alert definitions.
func (c *OpsRampClient) CreateLogAlertDefinition(tenantId string, request LogAlertDefinitionRequest) (*LogAlertDefinitionResponse, error) {
	rb, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	apiUrl := fmt.Sprintf("%s/alertdefinitions/api/v3/tenants/%s/log-alerts", c.BaseUrl, tenantId)

	body, err := c.NewJsonRequest("POST", apiUrl, rb)
	if err != nil {
		return nil, err
	}

	var response LogAlertDefinitionResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// GetLogAlertDefinition retrieves a log alert definition by ID.
func (c *OpsRampClient) GetLogAlertDefinition(tenantId string, id string) (*LogAlertDefinition, error) {
	apiUrl := fmt.Sprintf("%s/alertdefinitions/api/v3/tenants/%s/log-alerts/%s", c.BaseUrl, tenantId, id)

	body, err := c.NewJsonRequest("GET", apiUrl, nil)
	if err != nil {
		return nil, err
	}

	var response LogAlertWrapper
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return &response.Alert, nil
}

// UpdateLogAlertDefinition updates a log alert definition.
func (c *OpsRampClient) UpdateLogAlertDefinition(tenantId string, id string, alert LogAlertDefinition) (*LogAlertDefinition, error) {
	wrapper := LogAlertWrapper{Alert: alert}
	rb, err := json.Marshal(wrapper)
	if err != nil {
		return nil, err
	}

	apiUrl := fmt.Sprintf("%s/alertdefinitions/api/v3/tenants/%s/log-alerts/%s", c.BaseUrl, tenantId, id)

	body, err := c.NewJsonRequest("PUT", apiUrl, rb)
	if err != nil {
		return nil, err
	}

	var response LogAlertWrapper
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return &response.Alert, nil
}

// DeleteLogAlertDefinition deletes a log alert definition by ID.
func (c *OpsRampClient) DeleteLogAlertDefinition(tenantId string, id string) error {
	apiUrl := fmt.Sprintf("%s/alertdefinitions/api/v3/tenants/%s/log-alerts/%s", c.BaseUrl, tenantId, id)

	_, err := c.NewJsonRequest("DELETE", apiUrl, nil)
	return err
}
