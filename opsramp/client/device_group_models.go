// SPDX-FileCopyrightText: Copyright Hewlett Packard Enterprise Development LP
// SPDX-License-Identifier: Apache-2.0

package client

type FilterCriteria struct {
	Id          int    `json:"id,omitempty"`
	SearchQuery string `json:"searchQuery,omitempty"`
}

type DeviceGroupAPI struct {
	Id             string          `json:"id,omitempty"`
	Name           string          `json:"name"`
	EntityType     string          `json:"entityType,omitempty"`
	FilterCriteria *FilterCriteria `json:"filterCriteria,omitempty"`
	Parent         *Parent         `json:"parent,omitempty"`
	Description    string          `json:"description,omitempty"`
}

// DeviceGroupChild represents a device/resource entry in a device group's child list.
type DeviceGroupChild struct {
	Id   string `json:"id"`
	Name string `json:"name,omitempty"`
	Type string `json:"type"`
}
