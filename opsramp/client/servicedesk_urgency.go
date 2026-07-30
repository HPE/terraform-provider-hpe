// (C) Copyright 2026 Hewlett Packard Enterprise Development LP
package client

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

type ServiceDeskUrgency struct {
	Id          string `json:"uniqueId,omitempty"`
	Name        string `json:"name"`
	Description string `json:"description"`
	State       bool   `json:"state"`
}

type ServiceDeskUrgencyDelete struct {
	Ids []string `json:"Ids"`
}

type SearchUrgencyResponse struct {
	Results []ServiceDeskUrgency `json:"results"`
}

// CreateServiceDeskUrgency - Create new ServiceDeskUrgency
func (c *OpsRampClient) CreateServiceDeskUrgency(resource ServiceDeskUrgency) (*ServiceDeskUrgency, error) {
	// Convert Request Data/Body to JSON
	rb, err := json.Marshal(resource)
	if err != nil {
		return nil, err
	}

	// Prepare the URL, Method and Payload fo the Client
	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/incidents/urgencies", c.BaseUrl, c.TenantId)
	method := "POST"

	// Create a new Request
	body, err := c.NewJsonRequest(method, apiUrl, rb)
	if err != nil {
		return nil, err
	}

	// Preparing Response Body to return and convert it to Golang Map Object
	var responseBody ServiceDeskUrgency
	err = json.Unmarshal(body, &responseBody)
	if err != nil {
		return nil, err
	}

	// Return ID of the record created
	return &responseBody, nil
}

// GetServiceDeskUrgency - Get the resource with ID of the resource
func (c *OpsRampClient) GetServiceDeskUrgency(id string) (*ServiceDeskUrgency, error) {
	// Prepare config for API request
	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/serviceDesk/config/urgencies?pageNo=1&pageSize=1&isDescendingOrder=false&sortName=name&queryString=id:%s", c.BaseUrl, c.TenantId, id)
	method := "GET"

	// Create a new Request
	body, err := c.NewJsonRequest(method, apiUrl, nil)
	if err != nil {
		// If error is not related to "No Details Found", return the error
		return nil, err
	}

	// Preparing Response Body to return and convert it Record Struct
	responseBody := SearchUrgencyResponse{}
	err = json.Unmarshal(body, &responseBody)
	if err != nil {
		return nil, err
	}
	if len(responseBody.Results) == 0 {
		// If no results are found, return nil
		return nil, nil
	}

	return &responseBody.Results[0], err
}

// UpdateServiceDeskUrgency - Update a resource using Az Client
func (c *OpsRampClient) UpdateServiceDeskUrgency(id string, updateRecord ServiceDeskUrgency) (*ServiceDeskUrgency, error) {
	// Convert data into JSON
	rb, err := json.Marshal(updateRecord)
	if err != nil {
		return nil, err
	}

	// Prepare config for API request
	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/incidents/urgencies/%s", c.BaseUrl, c.TenantId, id)
	method := "POST"

	// Create a new Request
	body, err := c.NewJsonRequest(method, apiUrl, rb)
	if err != nil {
		return nil, err
	}
	var urgency ServiceDeskUrgency

	err = json.Unmarshal(body, &urgency)
	if err != nil {
		return nil, err
	}

	// Return Body and error
	return &urgency, err
}

func (c *OpsRampClient) DeleteServiceDeskUrgency(id string) error {
	// Prepare config for API Request
	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/serviceDesk/urgencies", c.BaseUrl, c.TenantId)
	method := "DELETE"

	ids := ServiceDeskUrgencyDelete{Ids: []string{id}}

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

func (c *OpsRampClient) FindServiceDeskUrgencyByName(tenantId string, name string) (*ServiceDeskUrgency, error) {
	apiURL := fmt.Sprintf("%s/api/v2/tenants/%s/serviceDesk/config/urgencies?pageNo=1&pageSize=100&isDescendingOrder=false&sortName=name&queryString=name:%s", c.BaseUrl, tenantId, url.QueryEscape(name))
	body, err := c.NewJsonRequest("GET", apiURL, nil)
	if err != nil {
		return nil, err
	}

	var response SearchUrgencyResponse
	if err = json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	for i := range response.Results {
		if strings.EqualFold(strings.TrimSpace(response.Results[i].Name), strings.TrimSpace(name)) {
			return &response.Results[i], nil
		}
	}

	return nil, fmt.Errorf("service desk urgency with name %q not found", name)
}
