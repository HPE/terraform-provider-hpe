// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package client

import (
	"encoding/json"
	"fmt"
	"net/url"
)

// CreateDeviceManagementPolicy creates a new device management (resource management) policy.
func (c *OpsRampClient) CreateDeviceManagementPolicy(
	tenantId string,
	policy DeviceManagementPolicy,
) (*DeviceManagementPolicy, error) {
	rb, err := json.Marshal([]DeviceManagementPolicy{policy})
	if err != nil {
		return nil, err
	}

	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/policies/management", c.BaseUrl, tenantId)

	body, err := c.NewJsonRequest("POST", apiUrl, rb)
	if err != nil {
		return nil, err
	}

	var response []DeviceManagementPolicy
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	if len(response) == 0 {
		return nil, fmt.Errorf("device management policy create response did not contain a policy")
	}

	return &response[0], nil
}

// UpdateDeviceManagementPolicy updates an existing device management policy. The policy's
// numeric Id must be set so the API can identify which policy to update.
func (c *OpsRampClient) UpdateDeviceManagementPolicy(
	tenantId string,
	policy DeviceManagementPolicy,
) (*DeviceManagementPolicy, error) {
	rb, err := json.Marshal([]DeviceManagementPolicy{policy})
	if err != nil {
		return nil, err
	}

	// check every policy has id
	if policy.Id == 0 {
		return nil, fmt.Errorf("device management policy must have a valid Id for update")
	}

	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/policies/management", c.BaseUrl, tenantId)

	body, err := c.NewJsonRequest("PUT", apiUrl, rb)
	if err != nil {
		return nil, err
	}

	var response []DeviceManagementPolicy
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	if len(response) == 0 {
		return nil, fmt.Errorf("device management policy update response did not contain a policy")
	}

	return &response[0], nil
}

// GetDeviceManagementPolicy retrieves a single device management policy by its uid.
func (c *OpsRampClient) GetDeviceManagementPolicy(tenantId string, uid string) (*DeviceManagementPolicy, error) {
	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/policies/management/%s", c.BaseUrl, tenantId, uid)

	body, err := c.NewJsonRequest("GET", apiUrl, nil)
	if err != nil {
		return nil, err
	}

	var response DeviceManagementPolicy
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// DeleteDeviceManagementPolicy removes a device management policy and unassigns its
// agent policies from resources.
func (c *OpsRampClient) DeleteDeviceManagementPolicy(tenantId string, uid string) error {
	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/policies/management/%s", c.BaseUrl, tenantId, uid)

	_, err := c.NewJsonRequest("DELETE", apiUrl, nil)

	return err
}

// RunDeviceManagementPolicy runs (saves and applies) a device management policy,
// assigning its actions to all currently matching resources.
func (c *OpsRampClient) RunDeviceManagementPolicy(tenantId string, uid string) error {
	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/policies/management/%s/action/run", c.BaseUrl, tenantId, uid)

	_, err := c.NewJsonRequest("GET", apiUrl, nil)

	return err
}

// SearchDeviceManagementPolicies searches device management policies by name using the
// v3 paginated search endpoint. Pass an empty name to list all policies visible to the tenant.
func (c *OpsRampClient) SearchDeviceManagementPolicies(tenantId string, name string) ([]DeviceManagementPolicy, error) {
	apiUrl := fmt.Sprintf("%s/api/v3/tenants/%s/resources/policies/search", c.BaseUrl, tenantId)

	params := url.Values{}
	params.Set("pageNo", "1")
	params.Set("pageSize", "100")

	if name != "" {
		params.Set("name", name)
		params.Set("orderBy", "name")
	}

	apiUrl += "?" + params.Encode()

	body, err := c.NewJsonRequest("GET", apiUrl, nil)
	if err != nil {
		return nil, err
	}

	var response DeviceManagementPolicySearchResult
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return response.Results, nil
}

// FindDeviceManagementPolicyByName searches for a device management policy by exact name match.
func (c *OpsRampClient) FindDeviceManagementPolicyByName(tenantId string, name string) (*DeviceManagementPolicy, error) {
	results, err := c.SearchDeviceManagementPolicies(tenantId, name)
	if err != nil {
		return nil, err
	}

	for _, r := range results {
		if r.Name == name {
			return &r, nil
		}
	}

	return nil, fmt.Errorf("no device management policy found with name %q", name)
}
