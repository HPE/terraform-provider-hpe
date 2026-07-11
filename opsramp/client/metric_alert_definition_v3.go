// SPDX-FileCopyrightText: Copyright Hewlett Packard Enterprise Development LP
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"encoding/json"
	"fmt"
)

// --- Metric Alert Definition CRUD ---

// CreateMetricAlertDefinition creates a new metric alert definition.
func (c *OpsRampClient) CreateMetricAlertDefinition(tenantId string, request MetricAlertDefinitionRequest) (*MetricAlertDefinitionResponse, error) {
	rb, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	apiUrl := fmt.Sprintf("%s/alertdefinitions/api/v3/tenants/%s/metric-alerts", c.BaseUrl, tenantId)

	body, err := c.NewJsonRequest("POST", apiUrl, rb)
	if err != nil {
		return nil, err
	}

	var response MetricAlertDefinitionResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// UpdateMetricAlertDefinition updates an existing metric alert definition.
func (c *OpsRampClient) UpdateMetricAlertDefinition(tenantId string, id string, request MetricAlertDefinitionRequest) (*MetricAlertDefinitionResponse, error) {
	rb, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	apiUrl := fmt.Sprintf("%s/alertdefinitions/api/v3/tenants/%s/metric-alerts/%s", c.BaseUrl, tenantId, id)

	body, err := c.NewJsonRequest("PUT", apiUrl, rb)
	if err != nil {
		return nil, err
	}

	var response MetricAlertDefinitionResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// DeleteMetricAlertDefinition deletes a metric alert definition.
func (c *OpsRampClient) DeleteMetricAlertDefinition(tenantId string, id string) error {
	apiUrl := fmt.Sprintf("%s/alertdefinitions/api/v3/tenants/%s/metric-alerts/%s", c.BaseUrl, tenantId, id)

	_, err := c.NewJsonRequest("DELETE", apiUrl, nil)
	return err
}
