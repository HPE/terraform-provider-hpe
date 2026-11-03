// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package client

import (
	"encoding/json"
	"fmt"
	"net/url"
)

// SearchTemplates searches G2/legacy templates by name using the v2 search endpoint.
// Pass an empty name to list all templates visible to the tenant. scope filters by
// template scope (e.g. "Global", "Service Provider", "Partner", "Client"); pass an
// empty string to skip scope filtering.
func (c *OpsRampClient) SearchTemplates(tenantId string, name string, scope string) ([]TemplateSearchResult, error) {
	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/monitoring/templates/search", c.BaseUrl, tenantId)

	params := url.Values{}

	if name != "" || scope != "" {
		queryString := ""
		if name != "" {
			queryString += "name:" + name
		}
		if scope != "" {
			if name != "" {
				queryString += "+"
			}
			queryString += "scope:" + scope
		}
		params.Set("queryString", queryString)
	}

	if len(params) > 0 {
		apiUrl += "?" + params.Encode()
	}

	body, err := c.NewJsonRequest("GET", apiUrl, nil)
	if err != nil {
		return nil, err
	}

	var response TemplateSearchResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return response.Results, nil
}

// FindTemplateByName searches for a template by exact name match within the given scope.
func (c *OpsRampClient) FindTemplateByName(tenantId string, name string, scope string) (*TemplateSearchResult, error) {
	results, err := c.SearchTemplates(tenantId, name, scope)
	if err != nil {
		return nil, err
	}

	for _, r := range results {
		if r.Name == name {
			return &r, nil
		}
	}

	return nil, fmt.Errorf("no template found with name %q", name)
}
