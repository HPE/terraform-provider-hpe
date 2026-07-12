// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package client

// ManagementProfileResponse represents the API response for a management profile.
type ManagementProfile struct {
	Uuid        string `json:"uuId,omitempty"`
	Id          int    `json:"id,omitempty"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Type        string `json:"type"`
}

// ManagementProfileSearchResponse wraps the paged search response.
type ManagementProfileSearchResponse struct {
	Results []ManagementProfile `json:"results"`
}
