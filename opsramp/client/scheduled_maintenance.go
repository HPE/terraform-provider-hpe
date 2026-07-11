// SPDX-FileCopyrightText: Copyright Hewlett Packard Enterprise Development LP
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"encoding/json"
	"fmt"
)

// CreateScheduledMaintenance creates a new scheduled maintenance window.
func (c *OpsRampClient) CreateScheduledMaintenance(tenantId string, request ScheduledMaintenanceRequest) (*ScheduledMaintenanceCreateResponse, error) {
	rb, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/scheduleMaintenances", c.BaseUrl, tenantId)

	body, err := c.NewJsonRequest("POST", apiUrl, rb)
	if err != nil {
		return nil, err
	}

	var response ScheduledMaintenanceCreateResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// GetScheduledMaintenance retrieves a scheduled maintenance window by ID.
func (c *OpsRampClient) GetScheduledMaintenance(tenantId string, smId string) (*ScheduledMaintenanceResponse, error) {
	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/scheduleMaintenances/%s", c.BaseUrl, tenantId, smId)

	body, err := c.NewJsonRequest("GET", apiUrl, nil)
	if err != nil {
		return nil, err
	}

	var response ScheduledMaintenanceResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// UpdateScheduledMaintenance updates an existing scheduled maintenance window.
func (c *OpsRampClient) UpdateScheduledMaintenance(tenantId string, smId string, request ScheduledMaintenanceRequest) (*ScheduledMaintenanceCreateResponse, error) {
	rb, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/scheduleMaintenances/%s", c.BaseUrl, tenantId, smId)

	body, err := c.NewJsonRequest("POST", apiUrl, rb)
	if err != nil {
		return nil, err
	}

	var response ScheduledMaintenanceCreateResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// DeleteScheduledMaintenance deletes a scheduled maintenance window.
func (c *OpsRampClient) DeleteScheduledMaintenance(tenantId string, smId string) error {
	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/scheduleMaintenances/%s", c.BaseUrl, tenantId, smId)

	_, err := c.NewJsonRequest("DELETE", apiUrl, nil)
	return err
}

// AddScheduledMaintenanceResources adds devices, device groups, and locations to a maintenance window.
func (c *OpsRampClient) AddScheduledMaintenanceResources(tenantId string, smId string, request ScheduledMaintenanceResourcesRequest) error {
	rb, err := json.Marshal(request)
	if err != nil {
		return err
	}

	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/scheduleMaintenances/%s/resources", c.BaseUrl, tenantId, smId)

	_, err = c.NewJsonRequest("POST", apiUrl, rb)
	return err
}

// RemoveScheduledMaintenanceResources removes devices, device groups, and locations from a maintenance window.
func (c *OpsRampClient) RemoveScheduledMaintenanceResources(tenantId string, smId string, request ScheduledMaintenanceResourcesRequest) error {
	rb, err := json.Marshal(request)
	if err != nil {
		return err
	}

	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/scheduleMaintenances/%s/resources", c.BaseUrl, tenantId, smId)

	_, err = c.NewJsonRequest("DELETE", apiUrl, rb)
	return err
}

// GetScheduledMaintenanceResourcesByType returns assigned resources of the given type.
// resourceType must be one of: "resources", "deviceGroups", "sites".
func (c *OpsRampClient) GetScheduledMaintenanceResourcesByType(tenantId string, smId string, resourceType string) (*ScheduledMaintenanceResourceTypeResponse, error) {
	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/scheduleMaintenances/%s/resources/%s", c.BaseUrl, tenantId, smId, resourceType)

	body, err := c.NewJsonRequest("GET", apiUrl, nil)
	if err != nil {
		return nil, err
	}

	var response ScheduledMaintenanceResourceTypeResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// ScheduledMaintenanceAction executes an action (resume or suspend) on a maintenance window.
func (c *OpsRampClient) ScheduledMaintenanceAction(tenantId string, smId string, action string) error {
	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/scheduleMaintenances/%s/%s", c.BaseUrl, tenantId, smId, action)

	_, err := c.NewJsonRequest("POST", apiUrl, nil)
	return err
}
