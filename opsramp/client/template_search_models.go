// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package client

// TemplateSearchResult represents a single template entry in a search response.
type TemplateSearchResult struct {
	Id            int    `json:"id"`
	UniqueId      string `json:"uniqueId"`
	Name          string `json:"name"`
	Description   string `json:"description,omitempty"`
	CollectorType string `json:"collectorType,omitempty"`
	Status        string `json:"status,omitempty"`
	Generation    int    `json:"generation,omitempty"`
	ResourceType  string `json:"resourceType,omitempty"`
	Tags          string `json:"tags,omitempty"`
	CreatedDate   string `json:"createdDate,omitempty"`
	UpdatedDate   string `json:"updatedDate,omitempty"`
}

// TemplateSearchResponse wraps the paged template search response.
type TemplateSearchResponse struct {
	Results      []TemplateSearchResult `json:"results"`
	TotalResults int                    `json:"totalResults"`
	PageNo       int                    `json:"pageNo"`
	PageSize     int                    `json:"pageSize"`
	TotalPages   int                    `json:"totalPages"`
}
