// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package client

import (
	"encoding/json"
	"fmt"
)

// GetAdapterTypes retrieves all monitoring adapter types. This is a global endpoint
// (not tenant-scoped) that returns the current, API-authoritative list of valid
// "applicationType" values for APPLICATION-adapter metric groups.
func (c *OpsRampClient) GetAdapterTypes() ([]AdapterType, error) {
	apiUrl := fmt.Sprintf("%s/api/v2/monitoring/adapterTypes", c.BaseUrl)

	body, err := c.NewJsonRequest("GET", apiUrl, nil)
	if err != nil {
		return nil, err
	}

	var response []AdapterType
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return response, nil
}

// FindAdapterTypeByName looks up an adapter type by exact name match (case-sensitive,
// matching the API's "name" field, e.g. "RSE").
func (c *OpsRampClient) FindAdapterTypeByName(name string) (*AdapterType, error) {
	adapterTypes, err := c.GetAdapterTypes()
	if err != nil {
		return nil, err
	}

	for _, a := range adapterTypes {
		if a.Name == name {
			return &a, nil
		}
	}

	return nil, fmt.Errorf("no adapter type found with name %q", name)
}
