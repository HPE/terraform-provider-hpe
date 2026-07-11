// SPDX-FileCopyrightText: Copyright Hewlett Packard Enterprise Development LP
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"encoding/json"
	"fmt"
)

// CreateScheduleMaintenance creates a new scheduled maintenance window.
func (c *OpsRampClient) CreateScheduleMaintenance(tenantId string, request ScheduleMaintenanceRequest) (*ScheduleMaintenanceCreateResponse, error) {
	rb, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/scheduleMaintenances", c.BaseUrl, tenantId)

	body, err := c.NewJsonRequest("POST", apiUrl, rb)
	if err != nil {
		return nil, err
	}

	var response ScheduleMaintenanceCreateResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// GetScheduleMaintenance retrieves a scheduled maintenance window by ID.
func (c *OpsRampClient) GetScheduleMaintenance(tenantId string, smId string) (*ScheduleMaintenanceResponse, error) {
	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/scheduleMaintenances/%s", c.BaseUrl, tenantId, smId)

	body, err := c.NewJsonRequest("GET", apiUrl, nil)
	if err != nil {
		return nil, err
	}

	var response ScheduleMaintenanceResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// UpdateScheduleMaintenance updates an existing scheduled maintenance window.
func (c *OpsRampClient) UpdateScheduleMaintenance(tenantId string, smId string, request ScheduleMaintenanceRequest) (*ScheduleMaintenanceCreateResponse, error) {
	rb, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/scheduleMaintenances/%s", c.BaseUrl, tenantId, smId)

	body, err := c.NewJsonRequest("POST", apiUrl, rb)
	if err != nil {
		return nil, err
	}

	var response ScheduleMaintenanceCreateResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// DeleteScheduleMaintenance deletes a scheduled maintenance window.
func (c *OpsRampClient) DeleteScheduleMaintenance(tenantId string, smId string) error {
	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/scheduleMaintenances/%s", c.BaseUrl, tenantId, smId)

	_, err := c.NewJsonRequest("DELETE", apiUrl, nil)
	return err
}

// AddScheduleMaintenanceResources adds devices, device groups, and locations to a maintenance window.
func (c *OpsRampClient) AddScheduleMaintenanceResources(tenantId string, smId string, request ScheduleMaintenanceResourcesRequest) error {
	rb, err := json.Marshal(request)
	if err != nil {
		return err
	}

	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/scheduleMaintenances/%s/resources", c.BaseUrl, tenantId, smId)

	_, err = c.NewJsonRequest("POST", apiUrl, rb)
	return err
}

// RemoveScheduleMaintenanceResources removes devices, device groups, and locations from a maintenance window.
func (c *OpsRampClient) RemoveScheduleMaintenanceResources(tenantId string, smId string, request ScheduleMaintenanceResourcesRequest) error {
	rb, err := json.Marshal(request)
	if err != nil {
		return err
	}

	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/scheduleMaintenances/%s/resources", c.BaseUrl, tenantId, smId)

	_, err = c.NewJsonRequest("DELETE", apiUrl, rb)
	return err
}

// GetScheduleMaintenanceResourcesByType returns assigned resources of the given type.
// resourceType must be one of: "resources", "deviceGroups", "sites".
func (c *OpsRampClient) GetScheduleMaintenanceResourcesByType(tenantId string, smId string, resourceType string) (*ScheduleMaintenanceResourceTypeResponse, error) {
	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/scheduleMaintenances/%s/resources/%s", c.BaseUrl, tenantId, smId, resourceType)

	body, err := c.NewJsonRequest("GET", apiUrl, nil)
	if err != nil {
		return nil, err
	}

	var response ScheduleMaintenanceResourceTypeResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// ScheduleMaintenanceAction executes an action (resume or suspend) on a maintenance window.
func (c *OpsRampClient) ScheduleMaintenanceAction(tenantId string, smId string, action string) error {
	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/scheduleMaintenances/%s/%s", c.BaseUrl, tenantId, smId, action)

	_, err := c.NewJsonRequest("POST", apiUrl, nil)
	return err
}
