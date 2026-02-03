// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package serviceplan

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
)

//go:generate ../../../../../../bin/render -out examples/resources/morpheus_service_plan/example.tf example.tf.tmpl Name "ExampleServicePlan" Code "exampleserviceplan" MaxMemory "4294967296" MaxStorage "536870912"  ProvisionTypeCode "arm" CustomMaxStorage "true" ConfigRangesMinStorage "268435456" ConfigRangesMaxStorage "536870912" SortOrder "10000" CoresPerSocket "1"

func RenderServicePlanConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":                   "ExampleServicePlan",
		"Code":                   "exampleserviceplan",
		"SortOrder":              "100",
		"MaxMemory":              "4294967296",
		"MaxStorage":             "536870912",
		"ProvisionTypeCode":      "arm",
		"CustomMaxStorage":       "true",
		"CoresPerSocket":         "1",
		"ConfigRangesMinStorage": "268435456",
		"ConfigRangesMaxStorage": "536870912",
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
	templatePath := filepath.Join(dir, "example.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}
