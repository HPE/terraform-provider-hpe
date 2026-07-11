// SPDX-FileCopyrightText: Copyright Hewlett Packard Enterprise Development LP
// SPDX-License-Identifier: Apache-2.0

package client

// IntegrationConfigResponse represents a single config object under an installed integration.
type IntegrationConfigResponse struct {
	ID      string         `json:"id"`
	Name    string         `json:"name"`
	Config  map[string]any `json:"config"`
	State   string         `json:"state"`
	InfoMap *InfoMap       `json:"infoMap,omitempty"`
}

// IntegrationConfigListResponse represents the paginated list of configs.
type IntegrationConfigListResponse struct {
	Results      []IntegrationConfigResponse `json:"results"`
	TotalResults int                         `json:"totalResults"`
}

// IntegrationConfigRequest represents the request body to create/update a config.
type IntegrationConfigRequest struct {
	Name         string         `json:"name"`
	Config       map[string]any `json:"config"`
	AllResources bool           `json:"allResources"`
	ScheduleNone bool           `json:"scheduleNone"` // For backward compatibility with older API versions
	Schedule     *Schedule      `json:"schedule,omitempty"`
	InfoMap      InfoMap        `json:"infoMap"`
}

type Schedule struct {
	PatternType string `json:"patternType"`
	Pattern     int64  `json:"pattern"`
	StartTime   string `json:"startTime"`
}

type InfoMap struct {
	KubernetesProxy *ProxyRef `json:"kubernetesProxy,omitempty"`
}

type ProxyRef struct {
	Uuid string `json:"uuid"`
	Name string `json:"name,omitempty"`
}
