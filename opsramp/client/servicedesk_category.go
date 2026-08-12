// (C) Copyright 2026 Hewlett Packard Enterprise Development LP
package client

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

type ServiceDeskCategory struct {
	Id          string `json:"uniqueId,omitempty"`
	Name        string `json:"name"`
	Description string `json:"description"`
	TicketType  string `json:"ticketType"`
}
type ServiceDeskCategoryDelete struct {
	Ids []string `json:"Ids"`
}

// CreateServiceDeskCategory - Create new ServiceDeskCategory
func (c *OpsRampClient) CreateServiceDeskCategory(tenantId string, resource ServiceDeskCategory) (*ServiceDeskCategory, error) {
	// Convert Request Data/Body to JSON
	rb, err := json.Marshal(resource)
	if err != nil {
		return nil, err
	}

	// Prepare the URL, Method and Payload fo the Client
	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/serviceDesk/config/categories", c.BaseUrl, tenantId)
	method := "POST"

	// Create a new Request
	body, err := c.NewJsonRequest(method, apiUrl, rb)
	if err != nil {
		return nil, err
	}

	// Preparing Response Body to return and convert it to Golang Map Object
	var responseBody []ServiceDeskCategory
	err = json.Unmarshal(body, &responseBody)
	if err != nil {
		return nil, err
	}

	// Return ID of the record created
	return &responseBody[0], nil
}

// GetServiceDeskCategory - Get the resource with ID of the resource
func (c *OpsRampClient) GetServiceDeskCategory(tenantId string, id string) (*ServiceDeskCategory, error) {
	// Prepare config for API request
	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/serviceDesk/config/category/%s", c.BaseUrl, tenantId, id)
	method := "GET"

	// Prepare the request

	// Create a new Request
	body, err := c.NewJsonRequest(method, apiUrl, nil)
	if err != nil {

		if !strings.Contains(err.Error(), "No Details Found") {
			// Not exists
			return nil, nil
		}

		// If error is not related to "No Details Found", return the error
		return nil, err
	}

	// Preparing Response Body to return and convert it Record Struct
	responseBody := ServiceDeskCategory{}
	err = json.Unmarshal(body, &responseBody)
	if err != nil {
		return nil, err
	}

	return &responseBody, nil
}

// UpdateServiceDeskCategory - Update a resource using Az Client
func (c *OpsRampClient) UpdateServiceDeskCategory(
	tenantId string,
	id string,
	updateRecord ServiceDeskCategory,
) (*ServiceDeskCategory, error) {
	// Convert data into JSON
	rb, err := json.Marshal(updateRecord)
	if err != nil {
		return nil, err
	}

	// Prepare config for API request
	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/serviceDesk/config/category/%s", c.BaseUrl, tenantId, id)
	method := "PUT"

	// Create a new Request
	body, err := c.NewJsonRequest(method, apiUrl, rb)
	if err != nil {
		return nil, err
	}

	var responseBody ServiceDeskCategory

	// Preparing Response Body to return and convert it to Golang Struct
	err = json.Unmarshal(body, &responseBody)
	if err != nil {
		return nil, err
	}

	// Return Body and error
	return &responseBody, err
}

func (c *OpsRampClient) DeleteServiceDeskCategory(tenantId string, id string) error {
	// Prepare config for API Request
	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/serviceDesk/categories", c.BaseUrl, tenantId)
	method := "DELETE"

	req := ServiceDeskCategoryDelete{Ids: []string{id}}

	// Convert Request Data/Body to JSON
	rb, err := json.Marshal(req)
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

func normalizeServiceDeskCategoryTicketType(ticketType string) string {
	switch strings.TrimSpace(ticketType) {
	case "incidents", "Incident":
		return "Incident"
	case "problems", "Problem":
		return "Problem"
	case "serviceRequests", "Service Request":
		return "Service Request"
	default:
		return strings.TrimSpace(ticketType)
	}
}

type searchServiceDeskCategoryResponse struct {
	Results []ServiceDeskCategory `json:"results"`
}

func (c *OpsRampClient) FindServiceDeskCategoryByName(
	tenantId string,
	name string,
	ticketType string,
) (*ServiceDeskCategory, error) {
	apiURL := fmt.Sprintf(
		"%s/api/v2/tenants/%s/serviceDesk/config/categories?pageNo=1&pageSize=100&isDescendingOrder=false&sortName=name&queryString=name:%s",
		c.BaseUrl,
		tenantId,
		url.QueryEscape(name),
	)
	body, err := c.NewJsonRequest("GET", apiURL, nil)
	if err != nil {
		return nil, err
	}

	var response searchServiceDeskCategoryResponse
	if err = json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	normalizedTicketType := normalizeServiceDeskCategoryTicketType(ticketType)
	for i := range response.Results {
		if !strings.EqualFold(strings.TrimSpace(response.Results[i].Name), strings.TrimSpace(name)) {
			continue
		}
		if normalizedTicketType != "" &&
			normalizeServiceDeskCategoryTicketType(response.Results[i].TicketType) != normalizedTicketType {
			continue
		}

		return &response.Results[i], nil
	}

	if normalizedTicketType != "" {
		return nil, fmt.Errorf("service desk category with name %q and ticket type %q not found", name, ticketType)
	}

	return nil, fmt.Errorf("service desk category with name %q not found", name)
}
