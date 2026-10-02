// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package client

import (
	"encoding/json"
	"fmt"
	"net/url"
)

// GetMetric retrieves a single metric definition by its exact metric name.
func (c *OpsRampClient) GetMetric(tenantId string, metricName string) (*Metric, error) {
	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/metrics/%s", c.BaseUrl, tenantId, url.PathEscape(metricName))

	body, err := c.NewJsonRequest("GET", apiUrl, nil)
	if err != nil {
		return nil, err
	}

	var response Metric
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// SearchMetrics searches metrics by name using the v2 metrics list endpoint. scope
// filters by metric scope (e.g. "global", "service provider", "partner", "client");
// pass an empty string to skip scope filtering.
func (c *OpsRampClient) SearchMetrics(tenantId string, metricName string, scope string) ([]Metric, error) {
	apiUrl := fmt.Sprintf("%s/api/v2/tenants/%s/metrics", c.BaseUrl, tenantId)

	params := url.Values{}

	if metricName != "" || scope != "" {
		queryString := ""
		if metricName != "" {
			queryString += "metricName:" + metricName
		}
		if scope != "" {
			if metricName != "" {
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

	var response MetricSearchResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return response.Results, nil
}

// FindMetricByName searches for a metric by exact name match within the given scope.
func (c *OpsRampClient) FindMetricByName(tenantId string, metricName string, scope string) (*Metric, error) {
	results, err := c.SearchMetrics(tenantId, metricName, scope)
	if err != nil {
		return nil, err
	}

	for _, m := range results {
		if m.MetricName == metricName {
			return &m, nil
		}
	}

	return nil, fmt.Errorf("no metric found with name %q", metricName)
}
