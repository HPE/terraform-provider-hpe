// SPDX-FileCopyrightText: Copyright Hewlett Packard Enterprise Development LP
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"encoding/json"
	"fmt"
	"strings"
)

// GetTimezones retrieves all available timezones from the API
func (c *OpsRampClient) GetTimezones() ([]TimezoneResponse, error) {
	// Prepare the URL, Method for the Client
	apiUrl := fmt.Sprintf("%s/api/v2/cfg/timezones", c.BaseUrl)
	method := "GET"

	// Create a new Request
	body, err := c.NewJsonRequest(method, apiUrl, nil)
	if err != nil {
		return nil, err
	}

	// Preparing Response Body
	var responseBody []TimezoneResponse
	err = json.Unmarshal([]byte(body), &responseBody)
	if err != nil {
		return nil, err
	}

	return responseBody, nil
}

// ValidateTimezone validates a timezone string against the available timezones
// It returns the matching TimeZone object if found, or an error if not found
func (c *OpsRampClient) ValidateTimezone(timezone string) (*TimeZone, error) {
	timezones, err := c.GetTimezones()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch timezones: %w", err)
	}

	// Search for the timezone by code, name, or label (case-insensitive)
	timezoneLower := strings.ToLower(timezone)

	for _, tz := range timezones {
		// Match by exact code
		if strings.EqualFold(tz.Code, timezone) {
			return &TimeZone{
				Code:  tz.Code,
				Id:    tz.Id,
				Name:  tz.Name,
				Label: tz.Label,
			}, nil
		}
		// Match by name (e.g., "America/Los_Angeles")
		if strings.EqualFold(tz.Name, timezone) {
			return &TimeZone{
				Code:  tz.Code,
				Id:    tz.Id,
				Name:  tz.Name,
				Label: tz.Label,
			}, nil
		}
		// Match by label
		if strings.EqualFold(tz.Label, timezone) {
			return &TimeZone{
				Code:  tz.Code,
				Id:    tz.Id,
				Name:  tz.Name,
				Label: tz.Label,
			}, nil
		}
		// Partial match by name (case-insensitive)
		if strings.Contains(strings.ToLower(tz.Name), timezoneLower) {
			return &TimeZone{
				Code:  tz.Code,
				Id:    tz.Id,
				Name:  tz.Name,
				Label: tz.Label,
			}, nil
		}
	}

	// Build available timezones list for error message (limit to first 20)
	var availableList []string
	for i, tz := range timezones {
		if i >= 20 {
			availableList = append(availableList, "... and more")
			break
		}
		availableList = append(availableList, tz.Name)
	}

	return nil, fmt.Errorf("timezone '%s' not found. Available timezones include: %s",
		timezone, strings.Join(availableList, ", "))
}

// FindTimezoneByName finds a timezone by its exact name
func (c *OpsRampClient) FindTimezoneByName(name string) (*TimeZone, error) {
	timezones, err := c.GetTimezones()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch timezones: %w", err)
	}

	for _, tz := range timezones {
		if tz.Name == name {
			return &TimeZone{
				Code:  tz.Code,
				Id:    tz.Id,
				Name:  tz.Name,
				Label: tz.Label,
			}, nil
		}
	}

	return nil, fmt.Errorf("timezone with name '%s' not found", name)
}
