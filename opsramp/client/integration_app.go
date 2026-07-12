// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package client

import (
	"encoding/json"
	"fmt"
)

// InstallIntegrationV3 installs an SDK APP integration using the v3 endpoint.
func (c *OpsRampClient) InstallIntegrationV3(tenantId string, request InstallIntegrationV3Request) (*IntegrationResponse, error) {
	rb, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	apiUrl := fmt.Sprintf("%s/api/v3/tenants/%s/apps/install", c.BaseUrl, tenantId)

	body, err := c.NewJsonRequest("POST", apiUrl, rb)
	if err != nil {
		return nil, err
	}

	var response IntegrationResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// GetIntegrationV3 retrieves an installed SDK APP integration using the v3 endpoint.
func (c *OpsRampClient) GetIntegrationV3(tenantId string, integrationId string) (*IntegrationResponseV3, error) {
	apiUrl := fmt.Sprintf("%s/api/v3/tenants/%s/apps/install/%s", c.BaseUrl, tenantId, integrationId)

	body, err := c.NewJsonRequest("GET", apiUrl, nil)
	if err != nil {
		return nil, err
	}

	var response IntegrationResponseV3
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// DeleteIntegrationV3 deletes an installed SDK APP integration using the v3 endpoint.
func (c *OpsRampClient) DeleteIntegrationV3(tenantId string, integrationId string, reason string) error {
	req := DeleteIntegrationRequestV3{UninstallReason: reason}
	rb, err := json.Marshal(req)
	if err != nil {
		return err
	}

	apiUrl := fmt.Sprintf("%s/api/v3/tenants/%s/apps/install/%s", c.BaseUrl, tenantId, integrationId)

	_, err = c.NewJsonRequest("DELETE", apiUrl, rb)
	return err
}
