// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package client

import (
	"encoding/json"
	"fmt"
)

// CreateFirstResponsePolicy creates a new first response policy
func (c *OpsRampClient) CreateFirstResponsePolicy(tenantId string, policy FirstResponsePolicy) (*FirstResponsePolicy, error) {
	rb, err := json.Marshal(policy)
	if err != nil {
		return nil, err
	}

	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/policies/firstResponse", c.BaseUrl, tenantId)

	body, err := c.NewJsonRequest("POST", apiUrl, rb)
	if err != nil {
		return nil, err
	}

	var responseBody FirstResponsePolicy
	err = json.Unmarshal(body, &responseBody)
	if err != nil {
		return nil, err
	}

	return &responseBody, nil
}

// GetFirstResponsePolicy retrieves a specific first response policy by ID
func (c *OpsRampClient) GetFirstResponsePolicy(tenantId string, policyId string) (*FirstResponsePolicy, error) {
	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/policies/firstResponse/%s", c.BaseUrl, tenantId, policyId)

	body, err := c.NewJsonRequest("GET", apiUrl, nil)
	if err != nil {
		return nil, err
	}

	var responseBody FirstResponsePolicy
	err = json.Unmarshal(body, &responseBody)
	if err != nil {
		return nil, err
	}

	return &responseBody, nil
}

// UpdateFirstResponsePolicy updates an existing first response policy
func (c *OpsRampClient) UpdateFirstResponsePolicy(
	tenantId string,
	policyId string,
	policy FirstResponsePolicy,
) (*FirstResponsePolicy, error) {
	rb, err := json.Marshal(policy)
	if err != nil {
		return nil, err
	}

	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/policies/firstResponse/%s", c.BaseUrl, tenantId, policyId)

	body, err := c.NewJsonRequest("POST", apiUrl, rb)
	if err != nil {
		return nil, err
	}

	var responseBody FirstResponsePolicy
	err = json.Unmarshal(body, &responseBody)
	if err != nil {
		return nil, err
	}

	return &responseBody, nil
}

// DeleteFirstResponsePolicy deletes a first response policy
func (c *OpsRampClient) DeleteFirstResponsePolicy(tenantId string, policyId string) error {
	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/policies/firstResponse/%s", c.BaseUrl, tenantId, policyId)

	_, err := c.NewJsonRequest("DELETE", apiUrl, nil)

	return err
}

// SetFirstResponsePolicyStatus enables, disables, or sets a first response policy to observed mode
func (c *OpsRampClient) SetFirstResponsePolicyStatus(tenantId string, policyId string, status string) error {
	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/policies/firstResponse/%s/%s", c.BaseUrl, tenantId, policyId, status)

	_, err := c.NewJsonRequest("POST", apiUrl, nil)

	return err
}
