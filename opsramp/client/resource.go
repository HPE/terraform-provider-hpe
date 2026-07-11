// SPDX-FileCopyrightText: Copyright Hewlett Packard Enterprise Development LP
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"encoding/json"
	"fmt"
)

type OpsQLSearchRequest struct {
	ObjectType     string   `json:"objectType"`
	Fields         []string `json:"fields"`
	FilterCriteria string   `json:"filterCriteria"`
}

func (c *OpsRampClient) GetResourceTypes(tenantId string) ([]string, error) {
	// Prepare the URL, Method and Payload fo the Client
	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/allResourceTypes/minimal", c.BaseUrl, tenantId)
	method := "GET"

	// Create a new Request
	body, err := c.NewJsonRequest(method, apiUrl, nil)
	if err != nil {
		return nil, err
	}

	// Preparing Response Body to return and convert it to Golang Map Object
	var responseBody []ResourceType
	err = json.Unmarshal([]byte(body), &responseBody)
	if err != nil {
		return nil, err
	}

	var response []string
	for _, resourceType := range responseBody {
		response = append(response, resourceType.Name)
	}

	// Return ID of the record created
	return response, nil
}

// CreateResource - Create new Resource
func (c *OpsRampClient) CreateResource(tenantId string, resource CreateResource) (*ResourceCreated, error) {

	// Convert Request Data/Body to JSON
	rb, err := json.Marshal(resource)
	if err != nil {
		return nil, err
	}

	// Prepare the URL, Method and Payload fo the Client
	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/resources", c.BaseUrl, tenantId)
	method := "POST"

	// Create a new Request
	body, err := c.NewJsonRequest(method, apiUrl, rb)
	if err != nil {
		return nil, err
	}

	// Preparing Response Body to return and convert it to Golang Map Object
	var responseBody ResourceCreated
	err = json.Unmarshal([]byte(body), &responseBody)
	if err != nil {
		return nil, err
	}

	// Return ID of the record created
	return &responseBody, nil

}

// GetResource - Get the resource with ID of the resource
func (c *OpsRampClient) GetResource(tenantId string, uuid string) (*GetResource, error) {

	// Prepare config for API request
	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/resources/%s", c.BaseUrl, tenantId, uuid)
	method := "GET"

	// Prepare the request

	// Create a new Request
	body, err := c.NewJsonRequest(method, apiUrl, nil)
	if err != nil {
		return nil, err
	}

	// Preparing Response Body to return and convert it Record Struct
	responseBody := GetResource{}
	err = json.Unmarshal(body, &responseBody)
	if err != nil {
		return nil, err
	}

	return &responseBody, err
}

// UpdateResource - Update a resource using Az Client
func (c *OpsRampClient) UpdateResource(tenantId string, uuid string, updateRecord UpdateResource) (interface{}, error) {

	// Convert data into JSON
	rb, err := json.Marshal(updateRecord)
	if err != nil {
		return nil, err
	}

	// Prepare config for API request
	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/resources/%s", c.BaseUrl, tenantId, uuid)
	method := "POST"

	// Create a new Request
	body, err := c.NewJsonRequest(method, apiUrl, rb)
	if err != nil {
		return nil, err
	}

	// Return Body and error
	return body, err
}

func (c *OpsRampClient) DeleteResource(tenantId string, uuid string) (interface{}, error) {

	// Prepare config for API Request
	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/resources/%s", c.BaseUrl, tenantId, uuid)
	method := "DELETE"

	// Create a new Request
	body, err := c.NewJsonRequest(method, apiUrl, nil)
	if err != nil {
		return nil, err
	}

	// Return ResponseBody and error
	return body, err
}

func (c *OpsRampClient) FindResourceByName(name string) (interface{}, error) {

	// Prepare config for API Request
	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/resources/search?queryString=resourceName:%s", c.BaseUrl, c.TenantId, name)
	method := "GET"

	// Create a new Request
	body, err := c.NewJsonRequest(method, apiUrl, nil)
	if err != nil {
		return nil, err
	}

	// Parse response and extract UUID
	var resources SearchResources
	err = json.Unmarshal(body, &resources)
	if err == nil {
		return nil, fmt.Errorf("resource with name '%s' not found", name)
	}

	// Return first matching resource UUID
	return resources.Results[0].Uuid, nil
}

// OpsQLSearch - Query resources using OpsRamp Query Language (OpsQL)
// Returns true if at least one resource matches the query, the ID of the first matching resource (if any), the count of matching resources, and an error if the API call fails.
func (c *OpsRampClient) QueryResources(filterCriteria string) (bool, string, int, error) {
	updateRecord := OpsQLSearchRequest{
		ObjectType:     "resource",
		Fields:         []string{"id", "name", "resourceType"},
		FilterCriteria: filterCriteria,
	}

	// Convert data into JSON
	rb, err := json.Marshal(updateRecord)
	if err != nil {
		return false, "", -1, err
	}

	// Prepare config for API Request
	apiUrl := fmt.Sprintf("%s/opsql/api/v3/tenants/%s/queries", c.BaseUrl, c.TenantId)
	method := "POST"

	// Create a new Request
	body, err := c.NewJsonRequest(method, apiUrl, rb)
	if err != nil {
		return false, "", 0, err
	}

	// Parse response
	var resources SearchResources
	err = json.Unmarshal(body, &resources)
	if err != nil {
		return false, "", 0, err
	}

	count := len(resources.Results)
	exists := count > 0
	var firstId string
	if exists {
		firstId = resources.Results[0].Uuid
	}

	return exists, firstId, count, nil
}
