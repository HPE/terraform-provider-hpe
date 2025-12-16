// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package automation

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"
	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/testhelpers"
)

//go:generate go run ../../../../../cmd/render -out examples/resources/morpheus_scale_threshold/resource.tf morpheus_scale_threshold_resource.tf.tmpl Name example_scale_threshold AutoUpscale true AutoDownscale true MinCount 1 MaxCount 3 EnableCpuThreshold true MinCpuPercentage 30.0 MaxCpuPercentage 75.0 EnableMemoryThreshold true MinMemoryPercentage 20.0 MaxMemoryPercentage 60.0 EnableDiskThreshold true MinDiskPercentage 25.0 MaxDiskPercentage 80.0

func RenderScaleThresholdConfig(
	t *testing.T,
	name string,
	overrides map[string]string,
) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":                  name,
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

	for k, v := range overrides {
		defaults[k] = v
	}

	args := []string{}
	for k, v := range defaults {
		args = append(args, k, v)
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

func RenderExecuteScheduleConfig(
	t *testing.T,
	name string,
	overrides map[string]string,
) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":        name,
		"Description": "This schedule runs daily at 7 AM Mountain Time",
		"Enabled":     "false",
		"TimeZone":    "America/Denver",
		"Schedule":    "7 0 * * *",
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
	templatePath := filepath.Join(dir, "hpe_morpheus_execute_schedule_resource.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		"Name", defaults["Name"],
		"Description", defaults["Description"],
		"Enabled", defaults["Enabled"],
		"TimeZone", defaults["TimeZone"],
		"Schedule", defaults["Schedule"],
	)
}

