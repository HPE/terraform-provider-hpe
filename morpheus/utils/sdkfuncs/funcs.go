// (C) Copyright 2025-2026 Hewlett Packard Enterprise Development LP

package sdkfuncs

import "github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"

// Cloud

func NewAwsCloudConfig(endpoint string) *sdk.AddCloudsRequestZoneConfigAnyOf {
	cfg := &sdk.AddCloudsRequestZoneConfigAnyOf{}
	cfg.Endpoint = endpoint

	return cfg
}

func NewAzureCloudConfig() *sdk.AddCloudsRequestZoneConfigAnyOf1 {
	return &sdk.AddCloudsRequestZoneConfigAnyOf1{}
}

func NewHvmCloudConfig() *sdk.AddCloudsRequestZoneConfigAnyOf2 {
	return &sdk.AddCloudsRequestZoneConfigAnyOf2{}
}

func NewVmwareCloudConfig(apiUrl, apiVersion, datacenter string) *sdk.AddCloudsRequestZoneConfigAnyOf3 {
	cfg := &sdk.AddCloudsRequestZoneConfigAnyOf3{}
	cfg.ApiUrl = apiUrl
	cfg.ApiVersion = apiVersion
	cfg.Datacenter = datacenter

	return cfg
}

// Cluster

func NewHvmClusterServerConfig() *sdk.AddClusterRequestClusterServerConfigAnyOfOneOf6 {
	return &sdk.AddClusterRequestClusterServerConfigAnyOfOneOf6{}
}

func NewHvmClusterServerConfigAsAnyOf(
	cfg *sdk.AddClusterRequestClusterServerConfigAnyOfOneOf6,
) sdk.AddClusterRequestClusterServerConfigAnyOf {
	return sdk.AddClusterRequestClusterServerConfigAnyOf{
		AddClusterRequestClusterServerConfigAnyOfOneOf6: cfg,
	}
}
