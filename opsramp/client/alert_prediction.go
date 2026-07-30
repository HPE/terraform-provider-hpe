// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package client

import (
	"encoding/json"
	"fmt"
)

// CreateAlertPredictionPolicy creates a new alert prediction policy
func (c *OpsRampClient) CreateAlertPredictionPolicy(tenantId string, policy AlertPredictionPolicy) (*AlertPredictionPolicy, error) {
	rb, err := json.Marshal(policy)
	if err != nil {
		return nil, err
	}

	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/policies/alertprediction", c.BaseUrl, tenantId)

	body, err := c.NewJsonRequest("POST", apiUrl, rb)
	if err != nil {
		return nil, err
	}

	var responseBody AlertPredictionPolicy
	err = json.Unmarshal(body, &responseBody)
	if err != nil {
		return nil, err
	}

	return &responseBody, nil
}

// GetAlertPredictionPolicy retrieves a specific alert prediction policy by ID
func (c *OpsRampClient) GetAlertPredictionPolicy(tenantId string, policyId string) (*AlertPredictionPolicy, error) {
	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/policies/alertprediction/%s", c.BaseUrl, tenantId, policyId)

	body, err := c.NewJsonRequest("GET", apiUrl, nil)
	if err != nil {
		return nil, err
	}

	var responseBody AlertPredictionPolicy
	err = json.Unmarshal(body, &responseBody)
	if err != nil {
		return nil, err
	}

	return &responseBody, nil
}

// UpdateAlertPredictionPolicy updates an existing alert prediction policy
func (c *OpsRampClient) UpdateAlertPredictionPolicy(
	tenantId string,
	policyId string,
	policy AlertPredictionPolicy,
) (*AlertPredictionPolicy, error) {
	rb, err := json.Marshal(policy)
	if err != nil {
		return nil, err
	}

	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/policies/alertprediction/%s", c.BaseUrl, tenantId, policyId)

	body, err := c.NewJsonRequest("POST", apiUrl, rb)
	if err != nil {
		return nil, err
	}

	var responseBody AlertPredictionPolicy
	err = json.Unmarshal(body, &responseBody)
	if err != nil {
		return nil, err
	}

	return &responseBody, nil
}

// DeleteAlertPredictionPolicy deletes an alert prediction policy
func (c *OpsRampClient) DeleteAlertPredictionPolicy(tenantId string, policyId string) error {
	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/policies/alertprediction/%s", c.BaseUrl, tenantId, policyId)

	_, err := c.NewJsonRequest("DELETE", apiUrl, nil)

	return err
}
