// (C) Copyright 2025-2026 Hewlett Packard Enterprise Development LP

package sdkfuncs

import "github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"

var (
	// Cloud
	NewHvmCloudConfig    = sdk.NewAddCloudsRequestZoneConfigAnyOf2
	NewVmwareCloudConfig = sdk.NewAddCloudsRequestZoneConfigAnyOf3
)
