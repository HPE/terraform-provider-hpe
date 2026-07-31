// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package client

// AlertPredictionPolicy represents an alert prediction policy
type AlertPredictionPolicy struct {
	Id                       string   `json:"id,omitempty"`
	Name                     string   `json:"name"`
	EnabledMode              string   `json:"enabledMode"`
	OrganizationMatchingType string   `json:"organizationMatchingType,omitempty"`
	IncludedClients          []string `json:"includedClients,omitempty"`
	SeasonalityTimeFrame     string   `json:"seasonalityTimeFrame"`
	GeneratePredictionAlert  bool     `json:"generatePredictionAlert"`
	FilterQuery              string   `json:"filterQuery"`
}
