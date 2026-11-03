// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package client

// Metric represents a monitoring metric definition.
type Metric struct {
	MetricName         string  `json:"metricName"`
	MetricSource       string  `json:"metricSource,omitempty"`
	DisplayName        string  `json:"displayName,omitempty"`
	Description        string  `json:"description,omitempty"`
	Unit               string  `json:"unit,omitempty"`
	UnitLabel          string  `json:"unitLabel,omitempty"`
	GroupName          string  `json:"groupName,omitempty"`
	Factor             float64 `json:"factor,omitempty"`
	ConsolFunc         string  `json:"consolFunc,omitempty"`
	DisplayByComponent bool    `json:"displayByComponent"`
	DisplayOrder       int     `json:"displayOrder,omitempty"`
	TopMetric          bool    `json:"topMetric"`
	Status             string  `json:"status,omitempty"`
	Type               string  `json:"type,omitempty"`
	ResourceType       string  `json:"resourceType,omitempty"`
}

// MetricSearchResponse wraps the paged metrics search response.
type MetricSearchResponse struct {
	Results []Metric `json:"results"`
}
