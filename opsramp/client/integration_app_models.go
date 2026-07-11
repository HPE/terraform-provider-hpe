// SPDX-FileCopyrightText: Copyright Hewlett Packard Enterprise Development LP
// SPDX-License-Identifier: Apache-2.0

package client

// InstallIntegrationV3Request represents the request body to install an SDK APP via v3.
type InstallIntegrationV3Request struct {
	App                       string            `json:"app"`
	Version                   string            `json:"version"`
	Profile                   *InstallV3Profile `json:"profile"`
	MultiAppsDiscoveryEnabled bool              `json:"multiAppsDiscoveryEnabled"`
}

// InstallV3Profile is the gateway/profile reference for v3 installs.
type InstallV3Profile struct {
	UuId string `json:"uuId"`
}

// IntegrationResponse represents the API response for an installed integration
type IntegrationResponseV3 struct {
	ID            string                    `json:"id"`
	DisplayName   string                    `json:"displayName"`
	Description   string                    `json:"description,omitempty"`
	Integration   *IntegrationDetails       `json:"integration,omitempty"`
	InboundConfig *IntegrationInboundConfig `json:"inboundConfig,omitempty"`
	Category      string                    `json:"category,omitempty"`
	InstalledBy   string                    `json:"installedBy,omitempty"`
	InstalledTime string                    `json:"installedTime,omitempty"`
	Status        string                    `json:"status,omitempty"`
	State         string                    `json:"state,omitempty"`
	IPAddress     string                    `json:"ipAddress,omitempty"`
	CredentialSet string                    `json:"credentialSet,omitempty"`
	ConfigFiles   []ConfigFile              `json:"configFiles,omitempty"`
	AlertSource   *AlertSource              `json:"alertSource,omitempty"`
	Version       string                    `json:"version,omitempty"`

	// Configuration integrations
	ProviderProps             map[string]any     `json:"providerProps,omitempty"`
	ConfigDetails             map[string]any     `json:"configDetails,omitempty"`
	MultiAppsDiscoveryEnabled bool               `json:"multiAppsDiscoveryEnabled,omitempty"`
	DiscoveryProfiles         []DiscoveryProfile `json:"discoveryProfiles,omitempty"`
	AdditionalProperties      map[string]any     `json:"additionalProperties,omitempty"`
}

// DeleteIntegrationRequestV3 represents the request body to uninstall an integration
type DeleteIntegrationRequestV3 struct {
	UninstallReason string `json:"uninstallReason,omitempty"`
}
