// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package setting

//go:generate go run ../../../../../cmd/render -out examples/resources/morpheus_setting_guidance/resource.tf morpheus_setting_guidance_resource.tf.tmpl PowerSettingsAverageCpu 75 PowerSettingsMaximumCpu 500 PowerSettingsNetworkThreshold 2000 CpuUpsizeAverageCpu 50 CpuUpsizeMaximumCpu 99 MemoryUpsizeMinimumFreeMemory 10 MemoryDownsizeAverageFreeMemory 60 MemoryDownsizeMaximumFreeMemory 30
