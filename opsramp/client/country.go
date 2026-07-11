// SPDX-FileCopyrightText: Copyright Hewlett Packard Enterprise Development LP
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"encoding/json"
	"fmt"
	"strings"
)

// CountryResponse represents the API response for countries
type CountryResponse struct {
	Name string `json:"name"`
}

// GetCountries retrieves all available countries from the API
func (c *OpsRampClient) GetCountries() ([]CountryResponse, error) {
	// Prepare the URL, Method for the Client
	apiUrl := fmt.Sprintf("%s/api/v2/cfg/countries", c.BaseUrl)
	method := "GET"

	// Create a new Request
	body, err := c.NewJsonRequest(method, apiUrl, nil)
	if err != nil {
		return nil, err
	}

	// Preparing Response Body
	var responseBody []CountryResponse
	err = json.Unmarshal([]byte(body), &responseBody)
	if err != nil {
		return nil, err
	}

	return responseBody, nil
}

// ValidateCountry validates a country string against the available countries
// It returns the matching country name if found, or an error if not found
func (c *OpsRampClient) ValidateCountry(country string) (string, error) {
	countries, err := c.GetCountries()
	if err != nil {
		return "", fmt.Errorf("failed to fetch countries: %w", err)
	}

	// Search for the country by name (case-insensitive)
	for _, ct := range countries {
		if strings.EqualFold(ct.Name, country) {
			return ct.Name, nil
		}
	}

	// Build available countries list for error message (limit to first 20)
	var availableList []string
	for i, ct := range countries {
		if i >= 20 {
			availableList = append(availableList, "... and more")
			break
		}
		availableList = append(availableList, ct.Name)
	}

	return "", fmt.Errorf("country '%s' not found. Available countries include: %s",
		country, strings.Join(availableList, ", "))
}
