// SPDX-FileCopyrightText: Copyright Hewlett Packard Enterprise Development LP
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"encoding/json"
	"fmt"
)

// CreateAlertEscalationPolicy creates a new alert escalation policy
func (c *OpsRampClient) CreateAlertEscalationPolicy(tenantId string, policy AlertEscalationPolicy) (*AlertEscalationPolicy, error) {
	rb, err := json.Marshal(policy)
	if err != nil {
		return nil, err
	}

	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/escalations", c.BaseUrl, tenantId)

	body, err := c.NewJsonRequest("POST", apiUrl, rb)
	if err != nil {
		return nil, err
	}

	var responseBody AlertEscalationPolicy
	err = json.Unmarshal([]byte(body), &responseBody)
	if err != nil {
		return nil, err
	}

	return &responseBody, nil
}

// GetAlertEscalationPolicy retrieves a specific alert escalation policy by ID
func (c *OpsRampClient) GetAlertEscalationPolicy(tenantId string, policyId string) (*AlertEscalationPolicy, error) {
	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/escalations/%s", c.BaseUrl, tenantId, policyId)

	body, err := c.NewJsonRequest("GET", apiUrl, nil)
	if err != nil {
		return nil, err
	}

	var responseBody AlertEscalationPolicy
	err = json.Unmarshal([]byte(body), &responseBody)
	if err != nil {
		return nil, err
	}

	return &responseBody, nil
}

// UpdateAlertEscalationPolicy updates an existing alert escalation policy
func (c *OpsRampClient) UpdateAlertEscalationPolicy(tenantId string, policyId string, policy AlertEscalationPolicy) (*AlertEscalationPolicy, error) {
	rb, err := json.Marshal(policy)
	if err != nil {
		return nil, err
	}

	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/escalations/%s", c.BaseUrl, tenantId, policyId)

	body, err := c.NewJsonRequest("POST", apiUrl, rb)
	if err != nil {
		return nil, err
	}

	var responseBody AlertEscalationPolicy
	err = json.Unmarshal([]byte(body), &responseBody)
	if err != nil {
		return nil, err
	}

	return &responseBody, nil
}

// DeleteAlertEscalationPolicy deletes an alert escalation policy
func (c *OpsRampClient) DeleteAlertEscalationPolicy(tenantId string, policyId string) error {
	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/escalations/%s", c.BaseUrl, tenantId, policyId)

	_, err := c.NewJsonRequest("DELETE", apiUrl, nil)
	return err
}
