// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package automation

//go:generate go run ../../../../../cmd/render -out examples/resources/morpheus_scale_threshold/resource.tf hpe_morpheus_scale_threshold_resource.tf.tmpl Name example_scale_threshold AutoUpscale true AutoDownscale true MinCount 1 MaxCount 3 EnableCpuThreshold true MinCpuPercentage 30.0 MaxCpuPercentage 75.0 EnableMemoryThreshold true MinMemoryPercentage 20.0 MaxMemoryPercentage 60.0 EnableDiskThreshold true MinDiskPercentage 25.0 MaxDiskPercentage 80.0
