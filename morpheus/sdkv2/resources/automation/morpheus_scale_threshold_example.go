// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package automation

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
)

//go:generate ../../../../../bin/render -out examples/resources/morpheus_scale_threshold/resource.tf morpheus_scale_threshold_resource.tf.tmpl Name example_scale_threshold AutoUpscale true AutoDownscale true MinCount 1 MaxCount 3 EnableCpuThreshold true MinCpuPercentage 30.0 MaxCpuPercentage 75.0 EnableMemoryThreshold true MinMemoryPercentage 20.0 MaxMemoryPercentage 60.0 EnableDiskThreshold true MinDiskPercentage 25.0 MaxDiskPercentage 80.0

func RenderScaleThresholdConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":                  "Example",
		"AutoUpscale":           "true",
		"AutoDownscale":         "true",
		"MinCount":              "1",
		"MaxCount":              "3",
		"EnableCpuThreshold":    "true",
		"MinCpuPercentage":      "30.0",
		"MaxCpuPercentage":      "75.0",
		"EnableMemoryThreshold": "true",
		"MinMemoryPercentage":   "20.0",
		"MaxMemoryPercentage":   "60.0",
		"EnableDiskThreshold":   "true",
		"MinDiskPercentage":     "25.0",
		"MaxDiskPercentage":     "80.0",
	}

	for key, value := range overrides {
		defaults[key] = value
	}

	var args []string
	for key, value := range defaults {
		args = append(args, key, value)
	}

	// Get the directory where this source file is located
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("unable to get current file path")
	}
	dir := filepath.Dir(filename)
	templatePath := filepath.Join(dir, "morpheus_scale_threshold_resource.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}
