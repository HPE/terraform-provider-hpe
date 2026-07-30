// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package client

import (
	"encoding/json"
	"fmt"
)

// CreateDeviceGroup - Create new DeviceGroup (unofficial)
// api: POST v3 /api/v2/tenants/{clientId}/deviceGroups
func (c *OpsRampClient) CreateDeviceGroup(tenantId string, deviceGroup DeviceGroupAPI) (*DeviceGroupAPI, error) {
	// Transform object to array
	data := make([]DeviceGroupAPI, 1)
	data[0] = deviceGroup

	// Convert Request Data/Body to JSON
	rb, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	// Prepare the URL, Method and Payload fo the Client
	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/deviceGroups", c.BaseUrl, tenantId)
	method := "POST"

	// Create a new Request
	body, err := c.NewJsonRequest(method, apiUrl, rb)
	if err != nil {
		return nil, err
	}

	// Preparing Response Body to return and convert it to Golang Map Object
	var createdDeviceGroups []DeviceGroupAPI
	err = json.Unmarshal(body, &createdDeviceGroups)
	if err != nil {
		return nil, err
	}

	// Create child resources
	if len(createdDeviceGroups) < 1 {
		return nil, fmt.Errorf("device group not created")
	}

	// Return record created
	return &createdDeviceGroups[0], nil
}

// GetDeviceGroup - Get existing DeviceGroup
func (c *OpsRampClient) GetDeviceGroup(tenantId string, deviceGroupId string) (*DeviceGroupAPI, error) {
	// Prepare the URL, Method and Payload fo the Client
	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/deviceGroups/%s", c.BaseUrl, tenantId, deviceGroupId)
	method := "GET"

	// Create a new Request
	body, err := c.NewJsonRequest(method, apiUrl, nil)
	if err != nil {
		return nil, err
	}

	// Preparing Response Body to return and convert it to Golang Map Object
	var responseBody DeviceGroupAPI
	err = json.Unmarshal(body, &responseBody)
	if err != nil {
		return nil, err
	}

	// Return ID of the record created
	return &responseBody, nil
}

// GetDeviceGroupChilds - Get child resources assigned to a device group manually.
// api: GET /api/v2/tenants/{clientId}/deviceGroups/{resourceGroupId}/childs/search
func (c *OpsRampClient) GetDeviceGroupChilds(tenantId string, deviceGroupId string) ([]SearchResource, error) {
	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/deviceGroups/%s/childs/search?assignType=MANUAL&type=RESOURCE", c.BaseUrl, tenantId, deviceGroupId)
	method := "GET"

	body, err := c.NewJsonRequest(method, apiUrl, nil)
	if err != nil {
		return nil, err
	}

	var responseBody SearchResources
	err = json.Unmarshal(body, &responseBody)
	if err != nil {
		return nil, err
	}

	return responseBody.Results, nil
}

// DeleteDeviceGroup - DeleteDeviceGroup
func (c *OpsRampClient) DeleteDeviceGroup(tenantId string, deviceGroupId string) error {
	// Prepare the URL, Method and Payload fo the Client
	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/deviceGroups/%s", c.BaseUrl, tenantId, deviceGroupId)
	method := "DELETE"

	// Create a new Request
	_, err := c.NewJsonRequest(method, apiUrl, nil)
	if err != nil {
		return err
	}

	return nil
}

// AddDeviceGroupChilds - Add resources to a device group.
// api: POST /api/v2/tenants/{clientId}/deviceGroups/{resourceGroupId}/childs
func (c *OpsRampClient) AddDeviceGroupChilds(tenantId string, deviceGroupId string, ids []string) error {
	childs := make([]DeviceGroupChild, len(ids))
	for i, id := range ids {
		childs[i] = DeviceGroupChild{Id: id, Type: "DEVICE"}
	}

	rb, err := json.Marshal(childs)
	if err != nil {
		return err
	}

	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/deviceGroups/%s/childs", c.BaseUrl, tenantId, deviceGroupId)
	_, err = c.NewJsonRequest("POST", apiUrl, rb)

	return err
}

// RemoveDeviceGroupChilds - Remove resources from a device group.
// api: DELETE /api/v2/tenants/{clientId}/deviceGroups/{resourceGroupId}/childs
func (c *OpsRampClient) RemoveDeviceGroupChilds(tenantId string, deviceGroupId string, ids []string) error {
	childs := make([]DeviceGroupChild, len(ids))
	for i, id := range ids {
		childs[i] = DeviceGroupChild{Id: id, Type: "DEVICE"}
	}

	rb, err := json.Marshal(childs)
	if err != nil {
		return err
	}

	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/deviceGroups/%s/childs", c.BaseUrl, tenantId, deviceGroupId)
	_, err = c.NewJsonRequest("DELETE", apiUrl, rb)

	return err
}
