// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package client

// AdapterType represents a monitoring adapter type (used as metric group "applicationType"
// values for APPLICATION-adapter metric groups, e.g. RSE, HAProxy).
type AdapterType struct {
	Id          int    `json:"id"`
	Version     int    `json:"version,omitempty"`
	Name        string `json:"name"`
	DisplayName string `json:"displayName,omitempty"`
	GroupName   string `json:"groupName,omitempty"`
	Category    string `json:"category,omitempty"`
}
