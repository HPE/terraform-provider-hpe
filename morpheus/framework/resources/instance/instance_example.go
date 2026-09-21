// (C) Copyright 2025-2026 Hewlett Packard Enterprise Development LP

package instance

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/group"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
)

//go:generate ../../../../bin/render example.tf.tmpl Name "TestInstance" InstanceType "9" ResourcePool "pool-700" CloudName "MyCloud" NetworkId "755" LayoutId "644" DatastoreId "555" InstanceContext "dev" MultipleTags "true" GroupName "\"TestGroup\""
//go:generate ../../../../bin/render example_twonetworks.tf.tmpl Name "TestInstance" InstanceType "9" ResourcePool "pool-62299" GroupName "\"TestGroup\""
//go:generate ../../../../bin/render example_timeouts.tf.tmpl Name "TestInstance" InstanceType "9" ResourcePool "pool-62299" GroupName "\"TestGroup\""
//go:generate ../../../../bin/render example_evars.tf.tmpl Name "TestInstance" InstanceType "9" ResourcePool "pool-62299"
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
		// MORPH-16458: GroupName is the value rendered into the config's
		// hpe_morpheus_group data source lookup. Tests replace it below with a
		// reference to a per-run group resource, so no fixture code appears in
		// the published example (which renders a plain quoted literal).
		"GroupName": `"TestGroup"`,
		// MaxMemory renders a service_plan_options block, for resize coverage.
		"MaxMemory": "",
		// SingleVolume omits the non-root volume, for volume add/remove coverage.
		"SingleVolume": "",
		// Evars and Ports take a pre-rendered HCL fragment (the list elements),
		// so a test can express add/modify/remove scenarios exactly.
		"Evars": "",
		"Ports": "",
	}

	for key, value := range overrides {
		defaults[key] = value
	}

	// MORPH-16458: tests create a per-run group instead of relying on the
	// master tenant's seeded group (which a sub-tenant cannot use). The group's
	// name defaults to the instance name — tests pass a unique, randomized
	// Name, so the group is likewise unique per run (group names must be
	// unique per account, so a fixed name would collide on the second run).
	// The go:generate directives render the template directly, so the
	// published examples carry only the data-source lookup.
	groupName, ok := overrides["GroupName"]
	if !ok {
		groupName = defaults["Name"]
	}

	groupConfig, err := group.RenderGroupConfig(t, map[string]string{
		"Name":     groupName,
		"Code":     strings.ToLower(groupName),
		"GroupId1": "data.hpe_morpheus_cloud.vme_cloud.id",
	})
	if err != nil {
		return "", err
	}

	// The instance config's data source reads the group by referencing the
	// resource's name, so Terraform's dependency graph orders the read after
	// the group is created within the same apply.
	defaults["GroupName"] = "hpe_morpheus_group.example.name"

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

	instanceConfig, err := testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
	if err != nil {
		return "", err
	}

	return groupConfig + instanceConfig, nil
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

// RenderInstanceHVMWrongPoolConfig renders a test fixture (not an example) in
// which the requested network belongs to a different resource pool than
// config_hvm.resource_pool_id, so that Morpheus rejects the create during
// validation with "Invalid network" before anything is provisioned. The
// defaults target the HVM cluster on the shared acceptance appliance: network 1
// belongs to cluster "Duck" (pool 1); pool 1702 is a different cluster.
func RenderInstanceHVMWrongPoolConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":           `"tfacc-wrong-pool"`,
		"CloudId":        "1",
		"GroupId":        "265",
		"InstanceTypeId": "4",
		"LayoutId":       "1193",
		"PlanId":         "19",
		"NetworkId":      "1",
		"ResourcePoolId": `"pool-1702"`,
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

	templatePath := filepath.Join(filepath.Dir(filename), "example_hvm_wrong_pool.tf.tmpl")

	return testhelpers.RenderExample(t, templatePath, args...)
}
