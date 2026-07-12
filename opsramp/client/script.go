// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package client

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
)

// CreateScript creates a new RBA script inside a category
func (c *OpsRampClient) CreateScript(tenantId string, categoryId string, scriptData Script) (*Script, error) {
	script, err := json.Marshal(scriptData)
	if err != nil {
		return nil, err
	}

	apiUrl := fmt.Sprintf("%s/api/v3/tenants/%s/categories/%s/scripts", c.BaseUrl, tenantId, categoryId)

	body, err := c.NewJsonRequest("POST", apiUrl, script)
	if err != nil {
		return nil, err
	}

	var responseBody ScriptCreationResponse
	err = json.Unmarshal([]byte(body), &responseBody)
	if err != nil {
		return nil, err
	}

	if responseBody.StatusType != "CREATED" {
		return nil, fmt.Errorf("script creation failed with status: %s", responseBody.StatusType)
	}

	// Retrieve the created script to return full details
	createdScript, err := c.GetScript(tenantId, categoryId, responseBody.Entity)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve created script: %w", err)
	}
	return createdScript, nil
}

// GetScript retrieves a specific RBA script by ID
func (c *OpsRampClient) GetScript(tenantId string, categoryId string, scriptId string) (*Script, error) {
	apiUrl := fmt.Sprintf("%s/api/v3/tenants/%s/scripts/%s", c.BaseUrl, tenantId, scriptId)

	body, err := c.NewJsonRequest("GET", apiUrl, nil)
	if err != nil {
		return nil, err
	}

	var responseBody Script
	err = json.Unmarshal([]byte(body), &responseBody)
	if err != nil {
		return nil, err
	}

	// Transform base64 content
	scriptContent, err := base64.StdEncoding.DecodeString(responseBody.Attachment.File)
	if err != nil {
		return nil, err
	}
	responseBody.Attachment.File = string(scriptContent)

	return &responseBody, nil
}

// UpdateScript updates an existing RBA script
func (c *OpsRampClient) UpdateScript(tenantId string, categoryId string, scriptId string, scriptData Script) (*Script, error) {
	rb, err := json.Marshal(scriptData)
	if err != nil {
		return nil, err
	}

	apiUrl := fmt.Sprintf("%s/api/v3/tenants/%s/categories/%s/scripts/%s", c.BaseUrl, tenantId, categoryId, scriptId)

	body, err := c.NewJsonRequest("POST", apiUrl, rb)
	if err != nil {
		return nil, err
	}

	var responseBody ScriptCreationResponse
	err = json.Unmarshal([]byte(body), &responseBody)
	if err != nil {
		return nil, err
	}

	if responseBody.StatusType != "OK" {
		return nil, fmt.Errorf("script creation failed with status: %s", responseBody.StatusType)
	}

	// Retrieve the created script to return full details
	updatedScript, err := c.GetScript(tenantId, categoryId, scriptId)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve created script: %w", err)
	}
	return updatedScript, nil
}

// DeleteScript deletes an RBA script
func (c *OpsRampClient) DeleteScript(tenantId string, categoryId string, scriptId string) error {
	apiUrl := fmt.Sprintf("%s/api/v3/tenants/%s/categories/%s/scripts/%s", c.BaseUrl, tenantId, categoryId, scriptId)

	_, err := c.NewJsonRequest("DELETE", apiUrl, nil)
	return err
}
