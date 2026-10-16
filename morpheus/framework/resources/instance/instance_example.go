// (C) Copyright 2025-2026 Hewlett Packard Enterprise Development LP

package instance

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"testing"

	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/group"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
)

//go:generate ../../../../bin/render example.tf.tmpl Name "TestInstance" CloudName "MyCloud" ClusterName "MyCluster" NetworkName "MyNetwork" DatastoreName "MyDatastore" LayoutName "Single HVM" LayoutVersion "1.0" PlanName "1 CPU, 1GB Memory" NoAgent "false" InstanceContext "dev" MultipleTags "true" GroupName "\"TestGroup\""
//go:generate ../../../../bin/render example_twonetworks.tf.tmpl Name "TestInstance" InstanceType "9" ResourcePool "pool-62299" GroupName "\"TestGroup\""
//go:generate ../../../../bin/render example_timeouts.tf.tmpl Name "TestInstance" InstanceType "9" ResourcePool "pool-62299" GroupName "\"TestGroup\""
//go:generate ../../../../bin/render example_evars.tf.tmpl Name "TestInstance" InstanceType "9" ResourcePool "pool-62299"
//go:generate ../../../../bin/render example_vmware.tf.tmpl Name "TestInstance" InstanceType "9" ResourcePool "pool-1"
//go:generate ../../../../bin/render example_vmware_sp_options.tf.tmpl Name "TestInstance" InstanceType "9" ResourcePool "pool-1"
//go:generate ../../../../bin/render example_hvm.tf.tmpl Name "TestInstance" CloudName "hvm" ClusterName "MyCluster" NetworkName "MyNetwork" PlanName "1 CPU, 512MB Memory" InstanceTypeLayout "Single KVM VM" LayoutVersion "11" ImageName "morpheus-central" NoAgent "false" GroupName "\"TestGroup\""
//go:generate ../../../../bin/render example_metal.tf.tmpl Name "TestInstance" CloudName "aCloud" EnvironmentName "anEnvironment" GroupName "aGroup" InstanceTypeLayout "Single ILO Server" Role "aRole" PlanName "G3i"
//go:generate ../../../../bin/render example_aws.tf.tmpl Name "TestInstance" InstanceType "9" ResourcePool "pool-12284"
//go:generate ../../../../bin/render example_azure.tf.tmpl Name "TestInstance" InstanceType "9" ResourcePool "pool-12284" AzureRegion "eastus"
//go:generate ../../../../bin/render example_azure_subnet.tf.tmpl Name "TestInstance" InstanceType "9" ResourcePool "pool-12284" AzureRegion "eastus" SubnetId "1"

// The data sources example.tf.tmpl declares, by address, for configurations
// that build on the rendered instance (the clone tests add a clone and a
// second instance beside it and reuse these lookups rather than declare their
// own). Keep in step with the template: a rename there would otherwise show
// up only as an acceptance failure in another package.
const (
	// ExampleCloudData is the cloud: data.hpe_morpheus_cloud.vme_cloud.
	ExampleCloudData = "data.hpe_morpheus_cloud.vme_cloud"
	// ExampleClusterData is the cluster: data.hpe_morpheus_cluster.vme.
	ExampleClusterData = "data.hpe_morpheus_cluster.vme"
	// ExampleLayoutData is the layout (and its instance_type):
	// data.hpe_morpheus_instance_type_layout.vme.
	ExampleLayoutData = "data.hpe_morpheus_instance_type_layout.vme"
	// ExampleNetworkData is the network, within the cluster:
	// data.hpe_morpheus_network.vme.
	ExampleNetworkData = "data.hpe_morpheus_network.vme"
	// ExampleDatastoreData is the datastore, within the cluster -- declared
	// only when the fixture names one (DatastoreName):
	// data.hpe_morpheus_datastore.vme.
	ExampleDatastoreData = "data.hpe_morpheus_datastore.vme"
	// ExamplePlanData is the service plan: data.hpe_morpheus_service_plan.vme_plan.
	ExamplePlanData = "data.hpe_morpheus_service_plan.vme_plan"
	// ExampleImageData is the image -- declared only when the fixture names one
	// (ImageName): data.hpe_morpheus_image.vme.
	ExampleImageData = "data.hpe_morpheus_image.vme"
	// ExampleGroupData is the group: data.hpe_morpheus_group.example. The
	// tests point it at a per-run group resource (GroupName).
	ExampleGroupData = "data.hpe_morpheus_group.example"
	// ExampleInstance is the rendered instance: hpe_morpheus_instance.example.
	ExampleInstance = "hpe_morpheus_instance.example"
)

// RenderInstanceConfig renders example.tf.tmpl: an HVM/KVM instance whose
// cloud, cluster (and so resource pool), network, datastore, layout (and so
// instance type), plan and image are all looked up by name through data
// sources. The names come from the HVM fixture (testhelpers/hvm_fixture.go,
// TF_VAR_testacc_morpheus_hvm_*), so the test skips when they are unset; the
// datastore and image are optional there, and when unset the volumes use
// automatic datastore selection and the layout's own image.
func RenderInstanceConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	layoutName, layoutVersion := testhelpers.HVMLayout(t)

	defaults := map[string]string{
		"Name":            "TestInstance",
		"CloudName":       testhelpers.HVMCloudName(t),
		"ClusterName":     testhelpers.HVMClusterName(t),
		"NetworkName":     testhelpers.HVMNetworkName(t),
		"DatastoreName":   os.Getenv(testhelpers.EnvHVMDatastoreName),
		"LayoutName":      layoutName,
		"LayoutVersion":   layoutVersion,
		"PlanName":        testhelpers.HVMPlanName(t),
		"ImageName":       os.Getenv(testhelpers.EnvHVMImageName),
		"NoAgent":         strconv.FormatBool(testhelpers.HVMNoAgent()),
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
		// CreateTimeout adds a timeouts block with that create value. A test
		// may also override PlanName (a plan-change test needs a second plan).
		"CreateTimeout": "",
	}

	for key, value := range overrides {
		defaults[key] = value
	}

	// The go:generate directives render the template directly, so the
	// published examples carry only the data-source lookup; the tests get a
	// per-run group (see renderPerRunGroup).
	groupConfig, err := renderPerRunGroup(t, defaults["Name"], "data.hpe_morpheus_cloud.vme_cloud.id", overrides)
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

// RenderInstanceHVMConfig renders example_hvm.tf.tmpl, whose cloud, cluster,
// network, layout, plan and image are looked up by name; the names come from
// the HVM fixture (testhelpers/hvm_fixture.go), so the test skips when they
// are unset. The group is a per-run resource, as in RenderInstanceConfig.
//
// Unlike the VMware template this one always sets image_id: demonstrating the
// attribute is the reason the HVM example exists, so the image name is
// required here where RenderInstanceConfig treats it as optional.
func RenderInstanceHVMConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	layoutName, layoutVersion := testhelpers.HVMLayout(t)

	defaults := map[string]string{
		"Name":               "TestInstance",
		"CloudName":          testhelpers.HVMCloudName(t),
		"ClusterName":        testhelpers.HVMClusterName(t),
		"NetworkName":        testhelpers.HVMNetworkName(t),
		"PlanName":           testhelpers.HVMPlanName(t),
		"InstanceTypeLayout": layoutName,
		"LayoutVersion":      layoutVersion,
		"NoAgent":            strconv.FormatBool(testhelpers.HVMNoAgent()),
	}

	// The image comes from the fixture only when the test does not name one
	// of its own (the volume actual_size test brings a specific image), so
	// such a test is not skipped for want of the fixture's.
	if _, ok := overrides["ImageName"]; !ok {
		defaults["ImageName"] = envOrSkip(t, testhelpers.EnvHVMImageName, "a virtual image to provision from")
	}

	for key, value := range overrides {
		defaults[key] = value
	}

	groupConfig, err := renderPerRunGroup(t, defaults["Name"], "data.hpe_morpheus_cloud.hvm_cloud.id", overrides)
	if err != nil {
		return "", err
	}
	defaults["GroupName"] = "hpe_morpheus_group.example.name"

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

// envOrSkip returns the named environment variable, skipping the test with a
// message naming the variable and what it should point at when it is unset.
func envOrSkip(t *testing.T, name, requirement string) string {
	t.Helper()

	value := os.Getenv(name)
	if value == "" {
		t.Skip(name + " not set; skipping test requiring " + requirement)
	}

	return value
}

// renderPerRunGroup renders the group resource every instance configuration
// provisions into. MORPH-16458: tests create a per-run group instead of
// relying on the master tenant's seeded group (which a sub-tenant cannot use).
// The group's name defaults to the instance name -- tests pass a unique,
// randomized Name, so the group is likewise unique per run (group names must
// be unique per account, so a fixed name would collide on the second run).
// The group is attached to the cloud the instance config resolves by name.
// cloudRef is the HCL reference to the cloud data source the template declares
// (example.tf.tmpl: data.hpe_morpheus_cloud.vme_cloud.id).
func renderPerRunGroup(t *testing.T, instanceName, cloudRef string, overrides map[string]string) (string, error) {
	t.Helper()

	groupName, ok := overrides["GroupName"]
	if !ok {
		groupName = instanceName
	}

	return group.RenderGroupConfig(t, map[string]string{
		"Name":     groupName,
		"Code":     strings.ToLower(groupName),
		"GroupId1": cloudRef,
	})
}

// RenderInstanceHVMWrongPoolConfig renders a test fixture (not an example) in
// which the requested network belongs to one cluster's resource pool while
// config_hvm.resource_pool_id is another cluster's, so that Morpheus rejects
// the create during validation with "Invalid network" before anything is
// provisioned. The network and its cluster come from the HVM fixture; the
// other cluster is EnvHVMOtherClusterName. The group is a per-run resource.
func RenderInstanceHVMWrongPoolConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	layoutName, layoutVersion := testhelpers.HVMLayout(t)

	defaults := map[string]string{
		"Name":          `"tfacc-wrong-pool"`,
		"CloudName":     testhelpers.HVMCloudName(t),
		"ClusterName":   testhelpers.HVMClusterName(t),
		"NetworkName":   testhelpers.HVMNetworkName(t),
		"LayoutName":    layoutName,
		"LayoutVersion": layoutVersion,
		"PlanName":      testhelpers.HVMPlanName(t),
		// NetworkPoolless looks the network up in the cloud rather than in
		// the cluster, for a network that belongs to no resource pool.
		"NetworkPoolless": "",
	}

	// The second cluster is only needed when the test does not say which
	// pool to target itself, so a test that targets the fixture cluster is
	// not skipped for want of a second one.
	if _, ok := overrides["OtherClusterName"]; !ok {
		defaults["OtherClusterName"] = testhelpers.HVMOtherClusterName(t)
	}

	for key, value := range overrides {
		defaults[key] = value
	}

	// The group resource needs the cloud data source the template declares
	// under its own name.
	groupName := strings.Trim(defaults["Name"], `"`)
	groupConfig, err := group.RenderGroupConfig(t, map[string]string{
		"Name":     groupName,
		"Code":     strings.ToLower(groupName),
		"GroupId1": "data.hpe_morpheus_cloud.wrong_pool.id",
	})
	if err != nil {
		return "", err
	}
	defaults["GroupName"] = "hpe_morpheus_group.example.name"

	var args []string
	for key, value := range defaults {
		args = append(args, key, value)
	}

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("unable to get current file path")
	}

	templatePath := filepath.Join(filepath.Dir(filename), "example_hvm_wrong_pool.tf.tmpl")

	instanceConfig, err := testhelpers.RenderExample(t, templatePath, args...)
	if err != nil {
		return "", err
	}

	return groupConfig + instanceConfig, nil
}
