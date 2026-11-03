// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package client

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// monitoringVersionV2Headers requests the New UI (V2) payload shape (metricGroups)
// for the G2-only template endpoints.
var monitoringVersionV2Headers = map[string]string{"Monitoring-Version": "V2"}

// CreateTemplate creates a new G2 monitoring template.
func (c *OpsRampClient) CreateTemplate(tenantId string, request TemplateRequest) (*TemplateCreateResponse, error) {
	rb, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	apiUrl := fmt.Sprintf("%s/monitoring/api/v3/tenants/%s/templates", c.BaseUrl, tenantId)

	body, err := c.NewJsonRequestWithHeaders("POST", apiUrl, rb, monitoringVersionV2Headers)
	if err != nil {
		return nil, err
	}

	var created TemplateCreateResponse
	if err := json.Unmarshal(body, &created); err != nil {
		return nil, err
	}

	// Return the created template response as-is so Terraform can fetch
	// additional details using GetTemplate if needed.
	return &created, nil
}

// GetTemplate retrieves a G2 monitoring template by templateId.
func (c *OpsRampClient) GetTemplate(tenantId string, templateId string) (*Template, error) {
	apiUrl := fmt.Sprintf("%s/monitoring/api/v3/tenants/%s/templates/%s", c.BaseUrl, tenantId, templateId)

	body, err := c.NewJsonRequestWithHeaders("GET", apiUrl, nil, monitoringVersionV2Headers)
	if err != nil {
		return nil, err
	}

	var response Template
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// UpdateTemplate updates an existing G2 monitoring template.
func (c *OpsRampClient) UpdateTemplate(tenantId string, id string, request TemplateRequest) (*Template, error) {
	rb, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	apiUrl := fmt.Sprintf("%s/monitoring/api/v3/tenants/%s/templates/%s", c.BaseUrl, tenantId, id)

	_, err = c.NewJsonRequestWithHeaders("PUT", apiUrl, rb, monitoringVersionV2Headers)
	if err != nil {
		return nil, err
	}

	return c.GetTemplate(tenantId, id)
}

// isRetryableTemplateDeleteError reports whether the delete failed due to an
// optimistic locking conflict (code "0006"), which clears up on retry.
func isRetryableTemplateDeleteError(err error) bool {
	if err == nil {
		return false
	}

	errText := strings.ToLower(err.Error())

	return strings.Contains(errText, `"code":"0006"`)
}

// DeleteTemplate deletes a G2 monitoring template by ID.
// Templates assigned to resources cannot be deleted; unassign before destroying.
func (c *OpsRampClient) DeleteTemplate(tenantId string, id string) error {
	return c.deleteTemplate(tenantId, id, 3)
}

func (c *OpsRampClient) deleteTemplate(tenantId string, id string, retries int) error {
	apiUrl := fmt.Sprintf("%s/monitoring/api/v3/tenants/%s/templates/%s", c.BaseUrl, tenantId, id)

	_, err := c.NewJsonRequestWithHeaders("DELETE", apiUrl, nil, monitoringVersionV2Headers)
	if err == nil {
		return nil
	}

	if retries <= 0 || !isRetryableTemplateDeleteError(err) {
		return err
	}

	time.Sleep(2 * time.Second)

	return c.deleteTemplate(tenantId, id, retries-1)
}
