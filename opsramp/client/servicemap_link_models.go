// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package client

type CreateServicemapLink struct {
	// Mandatory
	Id     string  `json:"id"`
	Parent *Parent `json:"parent,omitempty"`
}

type ReadServicemapLink struct {
	// Mandatory
	Results []ReadServicemap `json:"results"`
}
