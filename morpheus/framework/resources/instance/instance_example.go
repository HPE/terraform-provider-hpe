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
//go:generate ../../../../bin/render example_hvm.tf.tmpl Name "TestInstance" CloudName "hvm" PlanName "1 CPU, 512MB Memory" InstanceTypeLayout "Single KVM VM" LayoutVersion "11" ImageName "morpheus-central" InstanceType "34" GroupId "1" NetworkId "1" ResourcePool "pool-1"
//go:generate ../../../../bin/render example_metal.tf.tmpl Name "TestInstance" CloudName "aCloud" EnvironmentName "anEnvironment" GroupName "aGroup" InstanceTypeLayout "Single ILO Server" Role "aRole" PlanName "G3i"
//go:generate ../../../../bin/render example_aws.tf.tmpl Name "TestInstance" InstanceType "9" ResourcePool "pool-12284"
//go:generate ../../../../bin/render example_azure.tf.tmpl Name "TestInstance" InstanceType "9" ResourcePool "pool-12284" AzureRegion "eastus"
//go:generate ../../../../bin/render example_azure_subnet.tf.tmpl Name "TestInstance" InstanceType "9" ResourcePool "pool-12284" AzureRegion "eastus" SubnetId "1"

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
		"UserGroup":       "",
		"StorageProfile":  "",
		// MaxMemory renders a service_plan_options block, for resize coverage.
		"MaxMemory": "",
		// SingleVolume omits the non-root volume, for volume add/remove coverage.
		"SingleVolume": "",
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

func RenderInstanceAzureSubnetConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":         "TestInstance",
		"InstanceType": "9",
		"ResourcePool": "pool-12284",
		"AzureRegion":  "eastus",
		"SubnetId":     "1",
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
	templatePath := filepath.Join(dir, "example_azure_subnet.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}

// RenderInstanceVMwareConfig renders example_vmware.tf.tmpl.
//
// ImageName is empty by default, which renders the example exactly as it is
// published: image_id commented out, to show that the layout's own image is
// used when the attribute is omitted. Supplying ImageName adds the
// hpe_morpheus_image data source and sets image_id from it.
//
// The image is looked up by name rather than id because ids are allocated per
// appliance, so a literal id only works on the appliance it was read from.
func RenderInstanceVMwareConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":         "TestInstance",
		"InstanceType": "9",
		"ResourcePool": "pool-1",
		"ImageName":    "",
	}

	for key, value := range overrides {
		defaults[key] = value
	}

	var args []string
	for key, value := range defaults {
		args = append(args, key, value)
	}

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("unable to get current file path")
	}
	dir := filepath.Dir(filename)
	templatePath := filepath.Join(dir, "example_vmware.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}

// RenderInstanceHVMConfig renders example_hvm.tf.tmpl. The defaults match the
// directive that renders that template, so the published example and the
// acceptance test exercise the same configuration.
//
// Unlike the VMware template this one always sets image_id: demonstrating the
// attribute is the reason the HVM example exists.
func RenderInstanceHVMConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":               "TestInstance",
		"CloudName":          "hvm",
		"PlanName":           "1 CPU, 512MB Memory",
		"InstanceTypeLayout": "Single KVM VM",
		// The layout name repeats once per guest version, so the version is
		// what actually pins it down. Version 11 is the Debian layout, which
		// pairs with instance type 34.
		"LayoutVersion": "11",
		"ImageName":     "morpheus-central",
		"InstanceType":  "34",
		"GroupId":       "1",
		"NetworkId":     "1",
		"ResourcePool":  "pool-1",
	}

	for key, value := range overrides {
		defaults[key] = value
	}

	var args []string
	for key, value := range defaults {
		args = append(args, key, value)
	}

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("unable to get current file path")
	}
	dir := filepath.Dir(filename)
	templatePath := filepath.Join(dir, "example_hvm.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}
