// SPDX-FileCopyrightText: Copyright Hewlett Packard Enterprise Development LP
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

func normalizeServicemapType(nodeType string) string {
	switch strings.ToLower(strings.TrimSpace(nodeType)) {
	case "service":
		return "service"
	case "resource":
		return "resource"
	default:
		return strings.ToLower(strings.TrimSpace(nodeType))
	}
}

func isRetryableServicemapCreateError(err error) bool {
	if err == nil {
		return false
	}

	errText := strings.ToLower(err.Error())
	return strings.Contains(errText, `"code":"0005"`)
}

func (c *OpsRampClient) findServicemapChildByName(tenantId string, parentId string, childName string, nodeType string) (*CreateServicemap, error) {
	parentMap, err := c.GetServicemap(tenantId, parentId)
	if err != nil {
		return nil, err
	}

	for _, service := range parentMap.Services {
		if service.Name == childName && normalizeServicemapType(service.Type) == normalizeServicemapType(nodeType) {
			serviceCopy := service
			return &serviceCopy, nil
		}
	}

	return nil, nil
}

func (c *OpsRampClient) createServicemapNode(tenantId string, servicemap CreateServicemap, retries int) (*CreateServicemap, error) {
	// Convert Request Data/Body to JSON
	rb, err := json.Marshal(servicemap)
	if err != nil {
		return nil, err
	}

	apiUrl := fmt.Sprintf("%s/api/v3/tenants/%s/serviceGroup", c.BaseUrl, tenantId)
	body, err := c.NewJsonRequest("POST", apiUrl, rb)
	if err == nil {
		var createdServicemap CreateServicemap
		unmarshalErr := json.Unmarshal([]byte(body), &createdServicemap)
		if unmarshalErr != nil {
			return nil, unmarshalErr
		}
		return &createdServicemap, nil
	}

	// Found error, check if it's a retryable one and if we have retries left
	// Child creation is recoverable: the API can create the node and still return 500
	// or a duplicate error while the parent relationship is settling.
	if servicemap.Parent != nil && servicemap.Parent.Id != "" {
		existing, lookupErr := c.findServicemapChildByName(tenantId, servicemap.Parent.Id, servicemap.Name, servicemap.Type)
		if lookupErr == nil && existing != nil {
			return existing, nil
		}
	}

	if retries <= 0 || !isRetryableServicemapCreateError(err) {
		return nil, err
	}

	time.Sleep(2 * time.Second)
	return c.createServicemapNode(tenantId, servicemap, retries-1)
}

// CreateServicemap - Create new Servicemap (unofficial)
// api: POST v3 /api/v3/tenants/{tenantId}/serviceGroup
func (c *OpsRampClient) CreateServicemap(tenantId string, servicemap CreateServicemap) (*CreateServicemap, error) {

	// Save and remove child services
	childs := servicemap.Services
	resources := servicemap.Resources
	servicemap.Services = make([]CreateServicemap, 0)
	servicemap.Resources = nil

	// Assert servicemap parameters
	if servicemap.Type == "service" && len(servicemap.Resources) > 0 {
		return nil, fmt.Errorf("servicemap node can't have 'type=Service' and 'resources' at the same time")
	}

	// Add required parameters in case type = "resource"

	if servicemap.ServiceAvailabilityMonitor == nil {
		servicemap.ServiceAvailabilityMonitor = &ServiceAvailabilityMonitor{
			AlertType: "availability",
			Metrics:   make([]string, 0),
		}
	}

	retries := 3
	createdServicemap, err := c.createServicemapNode(tenantId, servicemap, retries)
	if err != nil {
		return nil, err
	}

	// Create child resources
	if len(resources) > 0 {

		err = c.CreateServicemapResources(tenantId, createdServicemap.Id, resources)
		if err != nil {
			return nil, err
		}
	}

	// Recursive call to create children
	for _, child := range childs {

		// Adjust parent
		child.Parent = &Parent{
			Id: createdServicemap.Id,
		}

		createdChild, err := c.CreateServicemap(tenantId, child)
		if err != nil {
			return nil, err
		}

		createdServicemap.Services = append(createdServicemap.Services, *createdChild)
	}

	// Restore resources
	servicemap.Resources = resources

	// Return ID of the record created
	return createdServicemap, nil
}

func (c *OpsRampClient) CreateServicemapResources(tenantId string, parentId string, resources []string) error {
	var data []struct {
		Id string `json:"id"`
	}

	// Convert Request Data/Body to JSON
	for _, resource := range resources {
		res := struct {
			Id string `json:"id"`
		}{
			Id: resource,
		}
		data = append(data, res)
	}

	// Convert Request Data/Body to JSON
	rb, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// Prepare the URL, Method and Payload fo the Client
	apiUrl := fmt.Sprintf("%s/api/v3/tenants/%s/serviceGroups/%s/childs", c.BaseUrl, tenantId, parentId)
	method := "POST"

	// Create a new Request
	_, err = c.NewJsonRequest(method, apiUrl, rb)
	if err != nil {
		return err
	}

	time.Sleep(5 * time.Second)
	return nil
}

// GetServicemap - Get existing Servicemap
func (c *OpsRampClient) GetServicemap(tenantId string, servicemapId string) (*CreateServicemap, error) {

	// Prepare the URL, Method and Payload fo the Client
	apiUrl := fmt.Sprintf("%s/api/v3/tenants/%s/service-maps/%s", c.BaseUrl, tenantId, servicemapId)
	method := "GET"

	// Create a new Request
	body, err := c.NewJsonRequest(method, apiUrl, nil)
	if err != nil {
		return nil, err
	}

	// Preparing Response Body to return and convert it to Golang Map Object
	var responseBody ReadServicemap
	err = json.Unmarshal([]byte(body), &responseBody)
	if err != nil {
		return nil, err
	}

	translatedResponse := TranslateGetToCreateServicemapModels(responseBody)

	// Return ID of the record created
	return &translatedResponse, nil

}

// DeleteServicemap - DeleteServicemap
func (c *OpsRampClient) DeleteServicemap(tenantId string, servicemapId string) error {

	// Prepare the URL, Method and Payload fo the Client
	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/serviceGroups/%s", c.BaseUrl, tenantId, servicemapId)
	method := "DELETE"

	// Create a new Request
	_, err := c.NewJsonRequest(method, apiUrl, nil)
	if err != nil {
		return err
	}

	return nil
}
