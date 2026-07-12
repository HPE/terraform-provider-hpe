// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package client

// Custom Integration models for API OAuth2 tokens

// CreateCustomIntegration represents the request to create a new custom integration
type CreateCustomIntegration struct {
	DisplayName   string        `json:"displayName"`
	Category      string        `json:"category"` // ej. Custom, Collaboration, Monitoring, SSO, Automation, ADAPTER_INTEGRATION
	InboundConfig InboundConfig `json:"inboundConfig"`
}

// InboundConfig represents the configuration for inbound integrations
type InboundConfig struct {
	// For API OAuth2 integrations, we can include additional fields if needed in the future
	Authentication AuthenticationConfig `json:"authentication,omitempty"`
}

type AuthenticationConfig struct {
	AuthType    string        `json:"authType"` // ej. WEBHOOK, OAUTH2, ALL
	ApiKeyPairs []ApiKeyPair  `json:"apiKeyPairs,omitempty"`
	Role        RoleClientRef `json:"role"`
}

type ApiKeyPair struct {
	Key    string `json:"key,omitempty"`
	Secret string `json:"secret,omitempty"`
}

// CustomIntegrationResponse represents the API response for a custom integration
type CustomIntegrationResponse struct {
	Id            string                    `json:"id"`
	DisplayName   string                    `json:"displayName"`
	Integration   *CustomIntegrationDetails `json:"integration,omitempty"`
	InboundConfig *InboundConfig            `json:"inboundConfig,omitempty"`
	Category      string                    `json:"category"`
	InstalledBy   string                    `json:"installedBy,omitempty"`
	InstalledTime string                    `json:"installedTime,omitempty"`
	Status        string                    `json:"status,omitempty"`
}

// CustomIntegrationDetails represents the integration details from the API
type CustomIntegrationDetails struct {
	Id                    string `json:"id,omitempty"`
	Name                  string `json:"name,omitempty"`
	MultipleInstallations bool   `json:"multipleInstallations,omitempty"`
	IntegrationId         int    `json:"integrationId,omitempty"`
	SubCategory           string `json:"subCategory,omitempty"`
}

// GetAPICredentials extracts the API key and secret from the response
// Returns empty strings if not found
func (r *CustomIntegrationResponse) GetAPICredentials() (key string, secret string) {
	if r.InboundConfig != nil &&
		len(r.InboundConfig.Authentication.ApiKeyPairs) > 0 {
		return r.InboundConfig.Authentication.ApiKeyPairs[0].Key,
			r.InboundConfig.Authentication.ApiKeyPairs[0].Secret
	}
	return "", ""
}

// CustomIntegrationListResponse represents the API response when listing integrations
type CustomIntegrationListResponse struct {
	Results         []CustomIntegrationResponse `json:"results"`
	TotalCount      int                         `json:"totalCount"`
	OrderBy         string                      `json:"orderBy"`
	PageNo          int                         `json:"pageNo"`
	PageSize        int                         `json:"pageSize"`
	TotalPages      int                         `json:"totalPages"`
	NextPage        bool                        `json:"nextPage"`
	PreviousPage    bool                        `json:"previousPageNo"`
	DescendingOrder bool                        `json:"descendingOrder"`
}

type DeleteCustomIntegration struct {
	UninstallReason             string `json:"uninstallReason,omitempty"`
	KeepAgentInstalledResources string `json:"keepAgentInstalledResources,omitempty"`
}
