// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package client

import (
	"encoding/json"
	"fmt"
	"strings"
)

// CreateManagementProfile creates a new management profile.
func (c *OpsRampClient) CreateManagementProfile(tenantId string, request ManagementProfile) (*ManagementProfile, error) {
	rb, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/managementProfiles", c.BaseUrl, tenantId)

	body, err := c.NewJsonRequest("POST", apiUrl, rb)
	if err != nil {
		return nil, err
	}

	var response ManagementProfile
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// GetManagementProfile retrieves a management profile by ID.
func (c *OpsRampClient) GetManagementProfile(tenantId string, profileId int) (*ManagementProfile, error) {
	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/managementProfiles/%d", c.BaseUrl, tenantId, profileId)

	body, err := c.NewJsonRequest("GET", apiUrl, nil)
	if err != nil {
		return nil, err
	}

	var response ManagementProfile
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// UpdateManagementProfile updates an existing management profile.
func (c *OpsRampClient) UpdateManagementProfile(tenantId string, profileId int, request ManagementProfile) (*ManagementProfile, error) {
	rb, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/managementProfiles/%d", c.BaseUrl, tenantId, profileId)

	body, err := c.NewJsonRequest("POST", apiUrl, rb)
	if err != nil {
		return nil, err
	}

	var response ManagementProfile
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// DeleteManagementProfile deletes a management profile by ID.
func (c *OpsRampClient) DeleteManagementProfile(tenantId string, profileId int) error {
	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/managementProfiles/%d", c.BaseUrl, tenantId, profileId)

	_, err := c.NewJsonRequest("DELETE", apiUrl, nil)
	if err != nil {
		if strings.Contains(err.Error(), "404") || strings.Contains(err.Error(), "not found") {
			return nil
		}
		return err
	}

	return nil
}

// SearchManagementProfiles returns management profiles matching queryName.
// Pass an empty queryName to list all profiles.
func (c *OpsRampClient) SearchManagementProfiles(tenantId string, queryName string) ([]ManagementProfile, error) {
	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/managementProfiles/search", c.BaseUrl, tenantId)
	if queryName != "" {
		apiUrl += fmt.Sprintf("?name=%s", queryName)
	}

	body, err := c.NewJsonRequest("GET", apiUrl, nil)
	if err != nil {
		return nil, err
	}

	var response ManagementProfileSearchResponse
	if err := json.Unmarshal(body, &response); err != nil {
		// Fall back to a direct slice response
		var results []ManagementProfile
		if err2 := json.Unmarshal(body, &results); err2 != nil {
			return nil, fmt.Errorf("failed to parse search response: %v", err)
		}
		return results, nil
	}

	return response.Results, nil
}
