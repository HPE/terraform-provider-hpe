// SPDX-FileCopyrightText: Copyright Hewlett Packard Enterprise Development LP
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"encoding/json"
	"fmt"
)

type TenancyContextResponse struct {
	Client ClientMinimal `json:"client"`
}

type ClientContextResponse struct {
	Id   uint   `json:"id"`
	UId  string `json:"uId"`
	Name string `json:"name"`
}

// CreateClient creates a new client (sub-tenant)
func (c *OpsRampClient) TenancyContext() (*ClientMinimal, error) {

	// Prepare the URL, Method and Payload for the Client
	apiUrl := fmt.Sprintf("%s/itop/tenancyContext", c.BaseUrl)
	method := "GET"

	// Create a new Request
	body, err := c.NewJsonRequest(method, apiUrl, nil)
	if err != nil {
		return nil, err
	}

	// Preparing Response Body to return and convert it to Golang Map Object
	var responseBody TenancyContextResponse
	err = json.Unmarshal([]byte(body), &responseBody)
	if err != nil {
		return nil, err
	}

	var clientContext ClientMinimal
	clientContext.Id = responseBody.Client.Id
	clientContext.UniqueId = responseBody.Client.UniqueId
	clientContext.Name = responseBody.Client.Name

	return &clientContext, nil
}

type MeResponse struct {
	ID               string `json:"id"`
	FirstName        string `json:"firstName"`
	LastName         string `json:"lastName"`
	Email            string `json:"email"`
	OrganizationName string `json:"organizationName"`
	OrgId            string `json:"orgId"`
	OrgType          string `json:"orgType"`
}

func (c *OpsRampClient) GetTenantInfo(tenantId string) (*ClientResponse, error) {
	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s", c.BaseUrl, tenantId)
	method := "GET"

	// Create a new Request
	body, err := c.NewJsonRequest(method, apiUrl, nil)
	if err != nil {
		return nil, err
	}

	// Preparing Response Body to return and convert it to Golang Map Object
	var context ClientResponse
	err = json.Unmarshal([]byte(body), &context)
	if err != nil {
		return nil, err
	}

	return &context, nil
}

func (c *OpsRampClient) GetClientInfo(clientId string) (*ClientResponse, error) {
	apiUrl := fmt.Sprintf("%s/api/v2/tenants/client/%s", c.BaseUrl, clientId)
	method := "GET"

	// Create a new Request
	body, err := c.NewJsonRequest(method, apiUrl, nil)
	if err != nil {
		return nil, err
	}

	// Preparing Response Body to return and convert it to Golang Map Object
	var context ClientResponse
	err = json.Unmarshal([]byte(body), &context)
	if err != nil {
		return nil, err
	}

	return &context, nil
}
