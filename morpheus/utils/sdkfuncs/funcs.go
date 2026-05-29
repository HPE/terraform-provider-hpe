// (C) Copyright 2025-2026 Hewlett Packard Enterprise Development LP

package sdkfuncs

import "github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"

// Cloud

func NewAwsCloudConfig(endpoint string) *sdk.AddCloudsRequestZoneConfigAnyOf {
	cfg := sdk.NewAddCloudsRequestZoneConfigAnyOfWithDefaults()
	cfg.Endpoint = endpoint

	return cfg
}

func NewAzureCloudConfig() *sdk.AddCloudsRequestZoneConfigAnyOf1 {
	return sdk.NewAddCloudsRequestZoneConfigAnyOf1WithDefaults()
}

func NewHvmCloudConfig() *sdk.AddCloudsRequestZoneConfigAnyOf2 {
	return sdk.NewAddCloudsRequestZoneConfigAnyOf2WithDefaults()
}

func NewVmwareCloudConfig(apiUrl, apiVersion, datacenter string) *sdk.AddCloudsRequestZoneConfigAnyOf3 {
	cfg := sdk.NewAddCloudsRequestZoneConfigAnyOf3WithDefaults()
	cfg.ApiUrl = apiUrl
	cfg.ApiVersion = apiVersion
	cfg.Datacenter = datacenter

	return cfg
}

// Cluster

func NewHvmClusterServerConfig() *sdk.AddClusterRequestClusterServerConfigAnyOfOneOf6 {
	return sdk.NewAddClusterRequestClusterServerConfigAnyOfOneOf6WithDefaults()
}

func NewHvmClusterServerConfigAsAnyOf(
	cfg *sdk.AddClusterRequestClusterServerConfigAnyOfOneOf6,
) sdk.AddClusterRequestClusterServerConfigAnyOf {
	return sdk.AddClusterRequestClusterServerConfigAnyOf{
		AddClusterRequestClusterServerConfigAnyOfOneOf6: cfg,
	}
}
