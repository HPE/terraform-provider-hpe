// SPDX-FileCopyrightText: Copyright Hewlett Packard Enterprise Development LP
// SPDX-License-Identifier: Apache-2.0

package client

func TranslateGetToCreateServicemapModels(readModel ReadServicemap) CreateServicemap {
	var services []CreateServicemap
	for _, i := range readModel.Services {
		services = append(services, TranslateGetToCreateServicemapModels(i))
	}
	var nodeType string
	switch readModel.NodeType {
	case "SERVICE":
		nodeType = "Service"
	case "Google":
		nodeType = "Resource"
	}

	return CreateServicemap{
		Name: readModel.Name,
		Type: nodeType,
		Id:   readModel.Id,

		// For Children
		AvailabilityThreshold:      readModel.AvailabilityThreshold,
		ServiceAvailabilityMonitor: readModel.ServiceAvailabilityMonitor,
		Services:                   services,
		Resources:                  readModel.ResourceIds,
		FilterCriteria:             readModel.FilterCriteria,
	}
}

type Parent struct {
	Id string `json:"id"`
}

type ServiceAvailabilityMonitor struct {
	AlertType string   `json:"alertType"`
	Metrics   []string `json:"metrics"`
	MatchType string   `json:"matchType,omitempty"`
}

type AvailabilityThreshold struct {
	ThresholdType  string `json:"thresholdType"`
	ThresholdLimit int    `json:"thresholdLimit"`
}

type CreateServicemap struct {
	// Mandatory
	Name   string  `json:"name"`
	Type   string  `json:"type"`
	Parent *Parent `json:"parent,omitempty"`

	// For update
	Id string `json:"id,omitempty"`

	// Optional fields
	SearchQuery string `json:"searchQuery,omitempty"`

	// For Children
	AvailabilityThreshold      AvailabilityThreshold       `json:"availabilityThreshold"`
	ServiceAvailabilityMonitor *ServiceAvailabilityMonitor `json:"serviceAvailabilityMonitor,omitempty"`

	Services    []CreateServicemap `json:"services,omitempty"`
	LinkedNodes []string           `json:"linkedNodes,omitempty"`
	Resources   []string           `json:"resources,omitempty"`

	FilterCriteria *FilterCriteria `json:"filterCriteria,omitempty"`
}

type ReadServicemap struct {
	// Mandatory
	Name     string `json:"name"`
	NodeType string `json:"nodeType"`

	// For update
	Id string `json:"id,omitempty"`

	// For Children
	AvailabilityThreshold      AvailabilityThreshold       `json:"availabilityThreshold"`
	ServiceAvailabilityMonitor *ServiceAvailabilityMonitor `json:"serviceAvailabilityMonitor,omitempty"`
	Services                   []ReadServicemap            `json:"services,omitempty"`
	ResourceIds                []string                    `json:"resourceIds,omitempty"`

	FilterCriteria *FilterCriteria `json:"filterCriteria,omitempty"`
}
