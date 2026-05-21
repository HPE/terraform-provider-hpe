// (C) Copyright 2025-2026 Hewlett Packard Enterprise Development LP

package instance

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
)

//go:generate ../../../../bin/render example.tf.tmpl Name "TestInstance" InstanceType "9" ResourcePool "pool-700" CloudName "MyCloud" NetworkId "755" LayoutId "644" DatastoreId "555" InstanceContext "dev" MultipleTags "true"
//go:generate ../../../../bin/render example_twonetworks.tf.tmpl Name "TestInstance" InstanceType "9" ResourcePool "pool-62299"
//go:generate ../../../../bin/render example_timeouts.tf.tmpl Name "TestInstance" InstanceType "9" ResourcePool "pool-62299"
//go:generate ../../../../bin/render example_vmware.tf.tmpl Name "TestInstance" InstanceType "9" ResourcePool "pool-1"
//go:generate ../../../../bin/render example_vmware_sp_options.tf.tmpl Name "TestInstance" InstanceType "9" ResourcePool "pool-1"
//go:generate ../../../../bin/render example_metal.tf.tmpl Name "TestInstance" CloudName "aCloud" EnvironmentName "anEnvironment" GroupName "aGroup" InstanceTypeLayout "Single ILO Server" Role "aRole" PlanName "G3i"
//go:generate ../../../../bin/render example_aws.tf.tmpl Name "TestInstance" InstanceType "9" ResourcePool "pool-12284"
//go:generate ../../../../bin/render example_azure.tf.tmpl Name "TestInstance" InstanceType "9" ResourcePool "pool-12284" AzureRegion "eastus"

func RenderInstanceConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":            "TestInstance",
		"InstanceType":    "34",
		"ResourcePool":    "pool-1",
		"CloudName":       "hvm",
		"LayoutId":        "77",
		"NetworkId":       "1",
		"DatastoreId":     "1",
		"InstanceContext": "dev",
		"MultipleTags":    "",
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

func RenderInstanceAzureConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":         "TestInstance",
		"InstanceType": "9",
		"ResourcePool": "pool-12284",
		"AzureRegion":  "eastus",
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
	templatePath := filepath.Join(dir, "example_azure.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}
