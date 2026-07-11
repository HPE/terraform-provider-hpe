// SPDX-FileCopyrightText: Copyright Hewlett Packard Enterprise Development LP
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"encoding/json"
	"fmt"
)

// InstallIntegration installs a new integration of the given application type.
// The application parameter is the integration type ID (e.g. "CUSTOM", "CUSTOM-EVENT", "NEWRELIC", "VMWARE").
func (c *OpsRampClient) InstallIntegration(tenantId string, application string, request InstallIntegrationRequest) (*IntegrationResponse, error) {
	rb, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/integrations/install/%s", c.BaseUrl, tenantId, application)

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

// GetIntegration retrieves an installed integration by its ID.
func (c *OpsRampClient) GetIntegration(tenantId string, integrationId string) (*IntegrationResponse, error) {
	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/integrations/installed/%s", c.BaseUrl, tenantId, integrationId)

	body, err := c.NewJsonRequest("GET", apiUrl, nil)
	if err != nil {
		return nil, err
	}

	var response IntegrationResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// DeleteIntegration uninstalls an integration.
func (c *OpsRampClient) DeleteIntegration(tenantId string, integrationId string, reason string) error {
	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/integrations/installed/%s", c.BaseUrl, tenantId, integrationId)

	req := DeleteIntegrationRequest{UninstallReason: reason}
	rb, err := json.Marshal(req)
	if err != nil {
		return err
	}

	_, err = c.NewJsonRequest("DELETE", apiUrl, rb)
	return err
}

// GetInboundAuthentication retrieves the inbound authentication configuration for an integration.
func (c *OpsRampClient) GetInboundAuthentication(tenantId string, integrationId string) (*InboundAuthResponse, error) {
	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/integrations/installed/%s/inbound/authentication", c.BaseUrl, tenantId, integrationId)

	body, err := c.NewJsonRequest("GET", apiUrl, nil)
	if err != nil {
		return nil, err
	}

	var response InboundAuthResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// SetInboundAuthentication configures the inbound authentication for an integration.
func (c *OpsRampClient) SetInboundAuthentication(tenantId string, integrationId string, request SetInboundAuthRequest) (*InboundAuthResponse, error) {
	rb, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/integrations/installed/%s/inbound/authentication", c.BaseUrl, tenantId, integrationId)

	body, err := c.NewJsonRequest("POST", apiUrl, rb)
	if err != nil {
		return nil, err
	}

	var response InboundAuthResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// SetMappingAttributes configures the inbound attribute mapping for an integration.
func (c *OpsRampClient) SetMappingAttributes(tenantId string, integrationId string, request MappingAttributesRequest) error {
	rb, err := json.Marshal(request)
	if err != nil {
		return err
	}

	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/integrations/installed/%s/mappingAttr", c.BaseUrl, tenantId, integrationId)

	_, err = c.NewJsonRequest("POST", apiUrl, rb)
	return err
}

// DeleteMappingAttribute deletes a single installed inbound mapping attribute by its unique ID.
func (c *OpsRampClient) DeleteMappingAttribute(tenantId string, integrationId string, mappingId string) error {
	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/integrations/installed/%s/inbound/mappingAttr/%s", c.BaseUrl, tenantId, integrationId, mappingId)
	_, err := c.NewJsonRequest("DELETE", apiUrl, nil)
	return err
}

// DeleteOutboundMappingAttribute deletes a single installed outbound mapping attribute by its unique ID.
func (c *OpsRampClient) DeleteOutboundMappingAttribute(tenantId string, integrationId string, mappingId string) error {
	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/integrations/installed/%s/outbound/mappingAttr/%s", c.BaseUrl, tenantId, integrationId, mappingId)
	_, err := c.NewJsonRequest("DELETE", apiUrl, nil)
	return err
}

// GetInstalledMappingAttributes retrieves the current inbound attribute mappings for an integration.
func (c *OpsRampClient) GetInstalledMappingAttributes(tenantId string, integrationId string) ([]InstalledMappingResult, error) {
	var all []InstalledMappingResult
	pageNo := 1
	pageSize := 100
	for {
		apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/integrations/installed/%s/inbound/installedIntgMappings?pageNo=%d&pageSize=%d", c.BaseUrl, tenantId, integrationId, pageNo, pageSize)
		body, err := c.NewJsonRequest("GET", apiUrl, nil)
		if err != nil {
			return nil, err
		}
		var page InstalledMappingResponse
		if err := json.Unmarshal(body, &page); err != nil {
			return nil, err
		}
		all = append(all, page.Results...)
		if !page.NextPage {
			break
		}
		pageNo++
	}
	return all, nil
}

// SetEnableDropAlerts enables or disables dropping alerts for an integration.
func (c *OpsRampClient) SetEnableDropAlerts(tenantId string, integrationId string, enable bool) error {
	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/integrations/installed/%s/enableDropAlerts?enableDropAlerts=%t", c.BaseUrl, tenantId, integrationId, enable)

	_, err := c.NewJsonRequest("GET", apiUrl, nil)
	return err
}

// AssignProcessDefinition assigns or unassigns a process definition to an integration.
func (c *OpsRampClient) AssignProcessDefinition(tenantId string, integrationId string, processId string, assign bool) error {
	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/integrations/installed/%s/processDefinitons/%s/assignOrUnassign?assign=%t",
		c.BaseUrl, tenantId, integrationId, processId, assign)

	_, err := c.NewJsonRequest("GET", apiUrl, nil)
	return err
}

// SetNotifier configures the outbound notifier for an integration.
func (c *OpsRampClient) SetNotifier(tenantId string, integrationId string, request NotifierRequest) error {
	rb, err := json.Marshal(request)
	if err != nil {
		return err
	}

	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/integrations/installed/%s/notifier", c.BaseUrl, tenantId, integrationId)

	_, err = c.NewJsonRequest("POST", apiUrl, rb)
	return err
}

// SetAdditionalProperties sets additional properties on an integration.
func (c *OpsRampClient) SetAdditionalProperties(tenantId string, integrationId string, props map[string]string) error {
	rb, err := json.Marshal(props)
	if err != nil {
		return err
	}

	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/integrations/installed/%s/additionalProps", c.BaseUrl, tenantId, integrationId)

	_, err = c.NewJsonRequest("POST", apiUrl, rb)
	return err
}

// RegenerateIntegrationToken regenerates the webhook token for an integration.
func (c *OpsRampClient) RegenerateIntegrationToken(tenantId string, integrationId string) (*InboundAuthResponse, error) {
	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/integrations/installed/%s/inbound/authentication/regenerateSecretOrToken", c.BaseUrl, tenantId, integrationId)

	body, err := c.NewJsonRequest("POST", apiUrl, nil)
	if err != nil {
		return nil, err
	}

	var response InboundAuthResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// GetAvailableAlertSources retrieves the list of available alert sources for a given integration type.
// Used to validate alert_source_id for CUSTOM-EVENT integrations.
func (c *OpsRampClient) GetAvailableAlertSources(tenantId string, application string) ([]AlertSource, error) {
	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/cfg/alertSource/available/custIntg/%s", c.BaseUrl, tenantId, application)

	body, err := c.NewJsonRequest("GET", apiUrl, nil)
	if err != nil {
		return nil, err
	}

	var response []AlertSource
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return response, nil
}

// GetInboundEntityProperties retrieves the available properties for mapping on an installed integration.
// Entity is typically "ALERT".
func (c *OpsRampClient) GetInboundEntityProperties(tenantId string, integrationId string, entity string) ([]EntityProperty, error) {
	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/integrations/installed/%s/inbound/entities/%s/properties", c.BaseUrl, tenantId, integrationId, entity)

	body, err := c.NewJsonRequest("GET", apiUrl, nil)
	if err != nil {
		return nil, err
	}

	var response []EntityProperty
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return response, nil
}

// GetOutboundEntityProperties retrieves the available properties for mapping on an installed integration.
// Entity is typically "ALERT".
func (c *OpsRampClient) GetOutboundEntityProperties(tenantId string, integrationId string, entity string) ([]EntityProperty, error) {
	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/integrations/installed/%s/outbound/entities/%s/properties", c.BaseUrl, tenantId, integrationId, entity)

	body, err := c.NewJsonRequest("GET", apiUrl, nil)
	if err != nil {
		return nil, err
	}

	var response []EntityProperty
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return response, nil
}

func (c *OpsRampClient) GetIntegrationAvailableRoles(tenantId string, integrationId string) ([]RoleClientRef, error) {
	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/integrations/installed/%s/availableRoles", c.BaseUrl, tenantId, integrationId)

	body, err := c.NewJsonRequest("GET", apiUrl, nil)
	if err != nil {
		return nil, err
	}

	var response []RoleClientRef
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return response, nil
}

// UpdateIntegration updates the base properties (displayName, description) of an installed integration.
func (c *OpsRampClient) UpdateIntegration(tenantId string, integrationId string, request InstallIntegrationRequest) (*IntegrationResponse, error) {
	rb, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/integrations/installed/%s", c.BaseUrl, tenantId, integrationId)

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

// SetWebhookHandshake sets or deletes webhook handshake properties on an integration.
func (c *OpsRampClient) SetWebhookHandshake(tenantId string, integrationId string, props map[string]any) error {
	payload := map[string]any{
		"webhookHandshakeProps": props,
	}
	rb, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/integrations/installed/%s/additionalProps?webhookHandshake=true", c.BaseUrl, tenantId, integrationId)

	_, err = c.NewJsonRequest("POST", apiUrl, rb)
	return err
}

// DeleteWebhookHandshake removes webhook handshake properties from an integration.
func (c *OpsRampClient) DeleteWebhookHandshake(tenantId string, integrationId string, props map[string]any) error {
	payload := map[string]any{
		"webhookHandshakeProps": props,
	}
	rb, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/integrations/installed/%s/additionalProps?webhookHandshake=true", c.BaseUrl, tenantId, integrationId)

	_, err = c.NewJsonRequest("DELETE", apiUrl, rb)
	return err
}

// SearchAvailableIntegrations retrieves the list of available integrations for a tenant.
func (c *OpsRampClient) SearchAvailableIntegrations(tenantId string, name string) ([]AvailableIntegration, error) {
	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/integrations/available/search?queryString=name:%s&pageSize=500", c.BaseUrl, tenantId, name)

	body, err := c.NewJsonRequest("GET", apiUrl, nil)
	if err != nil {
		return nil, err
	}

	var response AvailableIntegrationSearchResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return response.Results, nil
}

// GetAvailableIntegration retrieves a single available integration by its application ID.
func (c *OpsRampClient) GetAvailableIntegration(tenantId string, applicationId string) (*AvailableIntegration, error) {
	results, err := c.SearchAvailableIntegrations(tenantId, applicationId)
	if err != nil {
		return nil, err
	}

	if len(results) > 0 {
		return &results[0], nil
	}

	return nil, fmt.Errorf("available integration '%s' not found", applicationId)
}

// CreateIntegrationEvent creates an outbound event on an installed integration.
func (c *OpsRampClient) CreateIntegrationEvent(tenantId string, integrationId string, req IntegrationEventRequest) (*IntegrationEventResponse, error) {
	rb, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/integrations/installed/%s/event", c.BaseUrl, tenantId, integrationId)
	body, err := c.NewJsonRequest("POST", apiUrl, rb)
	if err != nil {
		return nil, err
	}
	var resp IntegrationEventResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetIntegrationEvent retrieves an outbound event by its ID.
func (c *OpsRampClient) GetIntegrationEvent(tenantId string, integrationId string, eventId string) (*IntegrationEventResponse, error) {
	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/integrations/installed/%s/event/%s", c.BaseUrl, tenantId, integrationId, eventId)
	body, err := c.NewJsonRequest("GET", apiUrl, nil)
	if err != nil {
		return nil, err
	}
	var resp IntegrationEventResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// UpdateIntegrationEvent updates an existing outbound event.
func (c *OpsRampClient) UpdateIntegrationEvent(tenantId string, integrationId string, eventId string, req IntegrationEventRequest) (*IntegrationEventResponse, error) {
	rb, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/integrations/installed/%s/event/%s", c.BaseUrl, tenantId, integrationId, eventId)
	body, err := c.NewJsonRequest("POST", apiUrl, rb)
	if err != nil {
		return nil, err
	}
	var resp IntegrationEventResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// SetIntegrationEventActive activates or deactivates an outbound integration event.
// The API endpoint is POST .../event/{eventId}/activate or .../event/{eventId}/deactivate.
func (c *OpsRampClient) SetIntegrationEventActive(tenantId, integrationId, eventId string, active bool) error {
	action := "activate"
	if !active {
		action = "deactivate"
	}
	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/integrations/installed/%s/event/%s/%s", c.BaseUrl, tenantId, integrationId, eventId, action)
	_, err := c.NewJsonRequest("POST", apiUrl, nil)
	return err
}

// DeleteIntegrationEvent deletes an outbound event.
func (c *OpsRampClient) DeleteIntegrationEvent(tenantId string, integrationId string, eventId string) error {
	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/integrations/installed/%s/event/%s", c.BaseUrl, tenantId, integrationId, eventId)
	_, err := c.NewJsonRequest("DELETE", apiUrl, nil)
	return err
}

// GetInstalledOutboundMappingAttributes retrieves outbound mapping attributes for an integration.
func (c *OpsRampClient) GetInstalledOutboundMappingAttributes(tenantId string, integrationId string) ([]InstalledMappingResult, error) {
	var all []InstalledMappingResult
	pageNo := 1
	pageSize := 100
	for {
		apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/integrations/installed/%s/outbound/installedIntgMappings?pageNo=%d&pageSize=%d", c.BaseUrl, tenantId, integrationId, pageNo, pageSize)
		body, err := c.NewJsonRequest("GET", apiUrl, nil)
		if err != nil {
			return nil, err
		}
		var page InstalledMappingResponse
		if err := json.Unmarshal(body, &page); err != nil {
			return nil, err
		}
		all = append(all, page.Results...)
		if !page.NextPage {
			break
		}
		pageNo++
	}
	return all, nil
}
