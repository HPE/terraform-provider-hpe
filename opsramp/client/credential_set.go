// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package client

import (
	"encoding/json"
	"fmt"
)

// CreateCredentialSet creates a new credential set
func (c *OpsRampClient) CreateCredentialSet(tenantId string, data CredentialSet) (*CredentialSet, error) {
	// The request body necessarily carries the credential set the practitioner supplied;
	// sending it is the purpose of the call.
	//nolint:gosec // G117: marshalling this secret is intentional
	rb, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/credentialSets", c.BaseUrl, tenantId)

	body, err := c.NewJsonRequest("POST", apiUrl, rb)
	if err != nil {
		return nil, err
	}

	var responseBody CredentialSet
	err = json.Unmarshal(body, &responseBody)
	if err != nil {
		return nil, err
	}

	return &responseBody, nil
}

// GetCredentialSet retrieves a credential set by ID
func (c *OpsRampClient) GetCredentialSet(tenantId string, credentialSetId string) (*CredentialSet, error) {
	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/credentialSets/%s", c.BaseUrl, tenantId, credentialSetId)

	body, err := c.NewJsonRequest("GET", apiUrl, nil)
	if err != nil {
		return nil, err
	}

	var responseBody CredentialSet
	err = json.Unmarshal(body, &responseBody)
	if err != nil {
		return nil, err
	}

	return &responseBody, nil
}

// UpdateCredentialSet updates an existing credential set
func (c *OpsRampClient) UpdateCredentialSet(tenantId string, credentialSetId string, data CredentialSet) (*CredentialSet, error) {
	// The request body necessarily carries the credential set the practitioner supplied;
	// sending it is the purpose of the call.
	//nolint:gosec // G117: marshalling this secret is intentional
	rb, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/credentialSets/%s", c.BaseUrl, tenantId, credentialSetId)

	body, err := c.NewJsonRequest("POST", apiUrl, rb)
	if err != nil {
		return nil, err
	}

	var responseBody CredentialSet
	err = json.Unmarshal(body, &responseBody)
	if err != nil {
		return nil, err
	}

	return &responseBody, nil
}

// DeleteCredentialSet deletes a credential set
func (c *OpsRampClient) DeleteCredentialSet(tenantId string, credentialSetId string) error {
	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/credentialSets/%s", c.BaseUrl, tenantId, credentialSetId)

	_, err := c.NewJsonRequest("DELETE", apiUrl, nil)

	return err
}
