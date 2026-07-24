// (C) Copyright 2024-2025 Hewlett Packard Enterprise Development LP

package models

// Broker structs go here

// SubscriptionDetailsResponse is the response for Subscription Details from the broker
type SubscriptionDetailsResponse struct {
	ServiceInstanceID string `json:"ServiceInstanceID"`
	URL               string `json:"URL"`
}

// MorpheusTokenResponse is the response for Morpheus Token from the broker
type MorpheusTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Expires      int64  `json:"expires"`
	ExpiresIn    int64  `json:"expires_in"`
}

// TFMorpheusDetails is what we return to terraform
type TFMorpheusDetails struct {
	// ID is the ServiceInstanceID, added here for use by the provider when storing the data
	ID          string `json:"id"`
	AccessToken string `json:"access_token"`
	// ValidTill Unix timestamp of when the access_token expires in seconds
	ValidTill int64  `json:"valid_till"`
	URL       string `json:"URL"`
}

type CMPDetails struct {
	ServiceInstanceID string `json:"ServiceInstanceID"`
	TenantID          string `json:"TenantID"`
	TenantName        string `json:"TenantName"`
	LocationName      string `json:"LocationName"`
	URL               string `json:"URL"`
	TokenDetails      struct {
		AccessToken  string `json:"access_token"`
		ExpiresIn    int64  `json:"expires_in"`
		RefreshToken string `json:"refresh_token"`
		TokenType    string `json:"token_type"`
		Expires      int64  `json:"expires"`
		Success      bool   `json:"success"`
	} `json:"TokenDetails"`
}
