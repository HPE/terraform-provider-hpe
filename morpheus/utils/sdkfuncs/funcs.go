// (C) Copyright 2025-2026 Hewlett Packard Enterprise Development LP

package sdkfuncs

import "github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"

var (
	// Cloud
	NewAwsCloudConfig    = sdk.NewAddCloudsRequestZoneConfigAnyOf
	NewAzureCloudConfig  = sdk.NewAddCloudsRequestZoneConfigAnyOf1
	NewHvmCloudConfig    = sdk.NewAddCloudsRequestZoneConfigAnyOf2
	NewVmwareCloudConfig = sdk.NewAddCloudsRequestZoneConfigAnyOf3

	// Cluster
	NewHvmClusterServerConfig = sdk.NewAddClusterRequestClusterServerConfigAnyOfOneOf6
	// nolint:lll
	NewHvmClusterServerConfigAsAnyOf = sdk.AddClusterRequestClusterServerConfigAnyOfOneOf6AsAddClusterRequestClusterServerConfigAnyOf
)
