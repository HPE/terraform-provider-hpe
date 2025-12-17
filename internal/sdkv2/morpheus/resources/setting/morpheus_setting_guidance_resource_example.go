// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package setting

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"
	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/testhelpers"
)

//go:generate go run ../../../../../cmd/render -out examples/resources/morpheus_setting_guidance/resource.tf morpheus_setting_guidance_resource.tf.tmpl PowerSettingsAverageCpu 75 PowerSettingsMaximumCpu 500 PowerSettingsNetworkThreshold 2000 CpuUpsizeAverageCpu 50 CpuUpsizeMaximumCpu 99 MemoryUpsizeMinimumFreeMemory 10 MemoryDownsizeAverageFreeMemory 60 MemoryDownsizeMaximumFreeMemory 30

func RenderSettingGuidanceConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"PowerSettingsAverageCpu":         "75",
		"PowerSettingsMaximumCpu":         "500",
		"PowerSettingsNetworkThreshold":   "2000",
		"CpuUpsizeAverageCpu":             "50",
		"CpuUpsizeMaximumCpu":             "99",
		"MemoryUpsizeMinimumFreeMemory":   "10",
		"MemoryDownsizeAverageFreeMemory": "60",
		"MemoryDownsizeMaximumFreeMemory": "30",
	}

	for key, value := range overrides {
		defaults[key] = value
	}

	// Get the directory where this source file is located
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("unable to get current file path")
	}
	dir := filepath.Dir(filename)
	templatePath := filepath.Join(dir, "morpheus_setting_guidance_resource.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		"PowerSettingsAverageCpu", defaults["PowerSettingsAverageCpu"],
		"PowerSettingsMaximumCpu", defaults["PowerSettingsMaximumCpu"],
		"PowerSettingsNetworkThreshold", defaults["PowerSettingsNetworkThreshold"],
		"CpuUpsizeAverageCpu", defaults["CpuUpsizeAverageCpu"],
		"CpuUpsizeMaximumCpu", defaults["CpuUpsizeMaximumCpu"],
		"MemoryUpsizeMinimumFreeMemory", defaults["MemoryUpsizeMinimumFreeMemory"],
		"MemoryDownsizeAverageFreeMemory", defaults["MemoryDownsizeAverageFreeMemory"],
		"MemoryDownsizeMaximumFreeMemory", defaults["MemoryDownsizeMaximumFreeMemory"],
	)
}

