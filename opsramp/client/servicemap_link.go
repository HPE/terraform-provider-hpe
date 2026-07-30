// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package client

import (
	"encoding/json"
	"fmt"
	"time"
)

// CreateServicemapLink - Create new Servicemap Link
// api: POST /api/v2/tenants/{clientId}/serviceGroups/link
func (c *OpsRampClient) CreateServicemapLink(tenantId string, servicemaplink CreateServicemapLink) (*CreateServicemapLink, error) {
	// Convert Request Data/Body to JSON
	payload, err := json.Marshal([]CreateServicemapLink{servicemaplink})
	if err != nil {
		return nil, err
	}

	// Prepare the URL, Method and Payload fo the Client
	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/serviceGroups/link", c.BaseUrl, tenantId)
	method := "POST"

	// Create a new Request
	_, err = c.NewJsonRequest(method, apiUrl, payload)
	if err != nil {
		return nil, err
	}

	time.Sleep(5 * time.Second)

	// Return ID of the record created
	return &servicemaplink, nil
}

// GetServicemap - Get existing Servicemap
func (c *OpsRampClient) GetServicemapLink(tenantId string, serviceMapLink CreateServicemapLink) (*CreateServicemapLink, error) {
	// Prepare the URL, Method and Payload fo the Client
	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/serviceGroups/%s/childs/search", c.BaseUrl, tenantId, serviceMapLink.Parent.Id)
	method := "GET"

	response, err := c.NewJsonRequest(method, apiUrl, nil)
	if err != nil {
		return nil, err
	}

	// Preparing Response Body to return and convert it to Golang Map Object
	var responseBody ReadServicemapLink
	err = json.Unmarshal(response, &responseBody)
	if err != nil {
		return nil, err
	}

	// Filter child service maps and find link
	var found *CreateServicemapLink
	for _, v := range responseBody.Results {
		if v.Id == serviceMapLink.Id {
			found = &CreateServicemapLink{
				Id: serviceMapLink.Id,
				Parent: &Parent{
					Id: serviceMapLink.Parent.Id,
				},
			}
		}
	}

	// Return ID of the record created
	return found, nil
}

// DeleteServicemapLink - DeleteServicemapLink
func (c *OpsRampClient) DeleteServicemapLink(tenantId string, serviceMapLink CreateServicemapLink) error {
	// Prepare the URL, Method and Payload fo the Client
	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/serviceGroups/unLink/%s/%s", c.BaseUrl, tenantId, serviceMapLink.Parent.Id, serviceMapLink.Id)
	method := "DELETE"

	// Create a new Request
	_, err := c.NewJsonRequest(method, apiUrl, nil)
	if err != nil {
		return err
	}

	return nil
}
