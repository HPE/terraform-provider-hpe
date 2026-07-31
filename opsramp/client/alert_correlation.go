// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package client

import (
	"encoding/json"
	"fmt"
)

// CreateAlertCorrelationPolicy creates a new alert correlation policy
func (c *OpsRampClient) CreateAlertCorrelationPolicy(
	tenantId string,
	policy AlertCorrelationPolicy,
) (*AlertCorrelationPolicy, error) {
	rb, err := json.Marshal(policy)
	if err != nil {
		return nil, err
	}

	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/policies/alertCorrelation", c.BaseUrl, tenantId)

	body, err := c.NewJsonRequest("POST", apiUrl, rb)
	if err != nil {
		return nil, err
	}

	var responseBody AlertCorrelationPolicy
	err = json.Unmarshal(body, &responseBody)
	if err != nil {
		return nil, err
	}

	return &responseBody, nil
}

// GetAlertCorrelationPolicy retrieves a specific alert correlation policy by ID
func (c *OpsRampClient) GetAlertCorrelationPolicy(tenantId string, policyId string) (*AlertCorrelationPolicy, error) {
	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/policies/alertCorrelation/%s", c.BaseUrl, tenantId, policyId)

	body, err := c.NewJsonRequest("GET", apiUrl, nil)
	if err != nil {
		return nil, err
	}

	var responseBody AlertCorrelationPolicy
	err = json.Unmarshal(body, &responseBody)
	if err != nil {
		return nil, err
	}

	return &responseBody, nil
}

// UpdateAlertCorrelationPolicy updates an existing alert correlation policy
func (c *OpsRampClient) UpdateAlertCorrelationPolicy(
	tenantId string,
	policyId string,
	policy AlertCorrelationPolicy,
) (*AlertCorrelationPolicy, error) {
	rb, err := json.Marshal(policy)
	if err != nil {
		return nil, err
	}

	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/policies/alertCorrelation/%s", c.BaseUrl, tenantId, policyId)

	body, err := c.NewJsonRequest("POST", apiUrl, rb)
	if err != nil {
		return nil, err
	}

	var responseBody AlertCorrelationPolicy
	err = json.Unmarshal(body, &responseBody)
	if err != nil {
		return nil, err
	}

	return &responseBody, nil
}

// DeleteAlertCorrelationPolicy deletes an alert correlation policy
func (c *OpsRampClient) DeleteAlertCorrelationPolicy(tenantId string, policyId string) error {
	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/policies/alertCorrelation/%s", c.BaseUrl, tenantId, policyId)

	_, err := c.NewJsonRequest("DELETE", apiUrl, nil)

	return err
}

// SetAlertCorrelationPolicyMode enables, disables, or sets a policy to observed mode
func (c *OpsRampClient) SetAlertCorrelationPolicyMode(tenantId string, policyId string, mode string) error {
	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/policies/alertCorrelation/%s/%s", c.BaseUrl, tenantId, policyId, mode)

	_, err := c.NewJsonRequest("POST", apiUrl, nil)

	return err
}
