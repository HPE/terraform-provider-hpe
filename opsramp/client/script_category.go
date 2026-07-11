// SPDX-FileCopyrightText: Copyright Hewlett Packard Enterprise Development LP
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"encoding/json"
	"fmt"
)

// ListScriptCategories returns all Script categories (tree structure) for a tenant.
func (c *OpsRampClient) ListScriptCategories(tenantId string) ([]ScriptCategory, error) {
	apiUrl := fmt.Sprintf("%s/api/v3/tenants/%s/scripts-categories", c.BaseUrl, tenantId)

	body, err := c.NewJsonRequest("GET", apiUrl, nil)
	if err != nil {
		return nil, err
	}

	var responseBody []ScriptCategory
	err = json.Unmarshal([]byte(body), &responseBody)
	if err != nil {
		return nil, err
	}

	return responseBody, nil
}

// GetScriptCategory retrieves a specific RBA category by ID by searching the full list.
func (c *OpsRampClient) GetScriptCategory(tenantId string, categoryId string) (*ScriptCategory, error) {
	apiUrl := fmt.Sprintf("%s/api/v3/tenants/%s/scripts-category/%s", c.BaseUrl, tenantId, categoryId)

	body, err := c.NewJsonRequest("GET", apiUrl, nil)
	if err != nil {
		return nil, err
	}

	var responseBody ScriptCategory
	err = json.Unmarshal([]byte(body), &responseBody)
	if err != nil {
		return nil, err
	}

	return &responseBody, nil
}

// CreateScriptCategory creates a new RBA category.
func (c *OpsRampClient) CreateScriptCategory(tenantId string, categoryData ScriptCategory) (*ScriptCategory, error) {
	rb, err := json.Marshal(categoryData)
	if err != nil {
		return nil, err
	}

	apiUrl := fmt.Sprintf("%s/api/v3/tenants/%s/scripts-category", c.BaseUrl, tenantId)

	body, err := c.NewJsonRequest("POST", apiUrl, rb)
	if err != nil {
		return nil, err
	}

	var responseBody ScriptCategory
	err = json.Unmarshal([]byte(body), &responseBody)
	if err != nil {
		return nil, err
	}

	return &responseBody, nil
}

// UpdateScriptCategory updates an existing RBA category.
// The category ID must be set in the categoryData struct; the PUT URL has no ID path segment.
func (c *OpsRampClient) UpdateScriptCategory(tenantId string, categoryData ScriptCategory) (*ScriptCategory, error) {
	rb, err := json.Marshal(categoryData)
	if err != nil {
		return nil, err
	}

	apiUrl := fmt.Sprintf("%s/api/v3/tenants/%s/scripts-category", c.BaseUrl, tenantId)

	body, err := c.NewJsonRequest("PUT", apiUrl, rb)
	if err != nil {
		return nil, err
	}

	var responseBody ScriptCategory
	err = json.Unmarshal([]byte(body), &responseBody)
	if err != nil {
		return nil, err
	}

	return &responseBody, nil
}

// DeleteScriptCategory deletes an RBA category by ID.
func (c *OpsRampClient) DeleteScriptCategory(tenantId string, categoryId string) error {
	apiUrl := fmt.Sprintf("%s/api/v3/tenants/%s/scripts-category/%s", c.BaseUrl, tenantId, categoryId)

	_, err := c.NewJsonRequest("DELETE", apiUrl, nil)
	return err
}
