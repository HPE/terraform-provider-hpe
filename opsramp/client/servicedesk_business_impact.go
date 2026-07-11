// SPDX-FileCopyrightText: Copyright Hewlett Packard Enterprise Development LP
// SPDX-License-Identifier: Apache-2.0
package client

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

type ServiceDeskBusinessImpact struct {
	Id          string `json:"uniqueId,omitempty"`
	Name        string `json:"name"`
	Description string `json:"description"`
	State       bool   `json:"state,omitempty"`
}

type ServiceDeskBusinessImpactDelete struct {
	Ids []string `json:"Ids"`
}

type SearchBusinessImpactResponse struct {
	Results []ServiceDeskBusinessImpact `json:"results"`
}

// CreateServiceDeskBusinessImpact - Create new ServiceDeskBusinessImpact
func (c *OpsRampClient) CreateServiceDeskBusinessImpact(resource ServiceDeskBusinessImpact) (*ServiceDeskBusinessImpact, error) {

	// Convert Request Data/Body to JSON
	rb, err := json.Marshal(resource)
	if err != nil {
		return nil, err
	}

	// Prepare the URL, Method and Payload fo the Client
	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/incidents/businessImpacts", c.BaseUrl, c.TenantId)
	method := "POST"

	// Create a new Request
	body, err := c.NewJsonRequest(method, apiUrl, rb)
	if err != nil {
		return nil, err
	}

	// Preparing Response Body to return and convert it to Golang Map Object
	var responseBody ServiceDeskBusinessImpact
	err = json.Unmarshal([]byte(body), &responseBody)
	if err != nil {
		return nil, err
	}

	// Return ID of the record created
	return &responseBody, nil

}

// GetServiceDeskBusinessImpact - Get the resource with ID of the resource
func (c *OpsRampClient) GetServiceDeskBusinessImpact(id string) (*ServiceDeskBusinessImpact, error) {

	// Prepare config for API request
	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/serviceDesk/config/businessImpacts?pageNo=1&pageSize=1&isDescendingOrder=false&sortName=name&queryString=id:%s", c.BaseUrl, c.TenantId, id)
	method := "GET"

	// Create a new Request
	body, err := c.NewJsonRequest(method, apiUrl, nil)
	if err != nil {
		// If error is not "No Details Found", return the error
		return nil, err
	}

	// Preparing Response Body to return and convert it Record Struct
	var responseBody SearchBusinessImpactResponse
	err = json.Unmarshal(body, &responseBody)
	if err != nil {
		return nil, err
	}

	// If no results are found, return nil
	if len(responseBody.Results) == 0 {
		return nil, nil
	}

	return &responseBody.Results[0], err
}

// UpdateServiceDeskBusinessImpact - Update a resource using Az Client
func (c *OpsRampClient) UpdateServiceDeskBusinessImpact(id string, updateRecord ServiceDeskBusinessImpact) (*ServiceDeskBusinessImpact, error) {

	// Convert data into JSON
	rb, err := json.Marshal(updateRecord)
	if err != nil {
		return nil, err
	}

	// Prepare config for API request
	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/incidents/businessImpacts/%s", c.BaseUrl, c.TenantId, id)
	method := "POST"

	// Create a new Request
	body, err := c.NewJsonRequest(method, apiUrl, rb)
	if err != nil {
		return nil, err
	}

	var businessImpact ServiceDeskBusinessImpact

	err = json.Unmarshal([]byte(body), &businessImpact)
	if err != nil {
		return nil, err
	}

	// Return Body and error
	return &businessImpact, err
}

func (c *OpsRampClient) DeleteServiceDeskBusinessImpact(id string) error {

	// Prepare config for API Request
	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/serviceDesk/businessImpacts", c.BaseUrl, c.TenantId)
	method := "DELETE"

	ids := ServiceDeskBusinessImpactDelete{Ids: []string{id}}

	// Convert data into JSON
	rb, err := json.Marshal(ids)
	if err != nil {
		return err
	}

	// Create a new Request
	_, err = c.NewJsonRequest(method, apiUrl, rb)
	if err != nil {
		return err
	}

	// Return ResponseBody and error
	return nil
}

func (c *OpsRampClient) FindServiceDeskBusinessImpactByName(tenantId string, name string) (*ServiceDeskBusinessImpact, error) {
	apiURL := fmt.Sprintf("%s/api/v2/tenants/%s/serviceDesk/config/businessImpacts?pageNo=1&pageSize=100&isDescendingOrder=false&sortName=name&queryString=name:%s", c.BaseUrl, tenantId, url.QueryEscape(name))
	body, err := c.NewJsonRequest("GET", apiURL, nil)
	if err != nil {
		return nil, err
	}

	var response SearchBusinessImpactResponse
	if err = json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	for i := range response.Results {
		if strings.EqualFold(strings.TrimSpace(response.Results[i].Name), strings.TrimSpace(name)) {
			return &response.Results[i], nil
		}
	}

	return nil, fmt.Errorf("service desk business impact with name %q not found", name)
}
