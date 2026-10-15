// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package testhelpers

import (
	"os"
	"testing"
)

// The HVM fixture: the objects the HVM/KVM instance tests provision against,
// named by environment variables and resolved by Terraform data sources inside
// each test configuration. Nothing here is an id. Ids are allocated per
// appliance, so a configuration that carries one only ever works on the
// appliance it was read from; a name resolves wherever the object exists.
//
// There are no defaults, for the same reason the affinity-group variables have
// none: a guessed name turns a test that should skip into one that fails
// against the API. Each accessor skips the test, naming the variable and what
// it should point at, when its variable is unset.
const (
	// EnvHVMCloudName names the Morpheus private cloud the HVM cluster is in.
	EnvHVMCloudName = "TF_VAR_testacc_morpheus_hvm_cloud_name"

	// EnvHVMClusterName names the HVM cluster to provision into. The cluster's
	// resource pool (hpe_morpheus_cluster: permissions.resource_pool.id) is the
	// instance's config_hvm.resource_pool_id, and its networks and datastores
	// are where hvm_network_name and hvm_datastore_name are looked up.
	EnvHVMClusterName = "TF_VAR_testacc_morpheus_hvm_cluster_name"

	// EnvHVMNetworkName names a network of that cluster that a test VM can be
	// attached to, as shown by GET /api/networks?clusterId=<cluster>.
	EnvHVMNetworkName = "TF_VAR_testacc_morpheus_hvm_network_name"

	// EnvHVMDatastoreName names a datastore of that cluster for the test VM's
	// volumes, as shown by GET /api/clusters/<cluster>/datastores. Optional in
	// the instance tests: unset, the volumes use datastore_auto_selection.
	EnvHVMDatastoreName = "TF_VAR_testacc_morpheus_hvm_datastore_name"

	// EnvHVMLayoutName and EnvHVMLayoutVersion name the instance type layout
	// to provision with; the layout resolves its instance type too. Layout
	// names repeat once per guest version, so the version is what pins one
	// down.
	EnvHVMLayoutName    = "TF_VAR_testacc_morpheus_hvm_layout_name"
	EnvHVMLayoutVersion = "TF_VAR_testacc_morpheus_hvm_layout_version"

	// EnvHVMPlanName names the kvm service plan for the test VMs; something
	// small, such as "1 CPU, 1GB Memory".
	EnvHVMPlanName = "TF_VAR_testacc_morpheus_hvm_plan_name"

	// EnvHVMImageName names the virtual image to provision from, when the
	// layout does not carry one. Optional in the instance tests: unset, the
	// layout's own image is used.
	EnvHVMImageName = "TF_VAR_testacc_morpheus_hvm_image_name"

	// EnvHVMGrowingImageName names a virtual image whose raw size is a little
	// over its minimum disk, so a volume requested at the minimum is grown by
	// the platform (a 10 GB request on a 10.002 GB image becomes 11 GB). The
	// volume actual_size test needs exactly that image; most appliances have
	// none, and the test skips without it.
	EnvHVMGrowingImageName = "TF_VAR_testacc_morpheus_hvm_growing_image_name"

	// EnvHVMNoAgent, when "true", provisions the test VMs without the Morpheus
	// agent. Default: with the agent. The agent is how the appliance learns a
	// VM's address and applies environment variables on a Linux guest;
	// without it Morpheus falls back to SSH from the appliance into the
	// guest, which needs a route the appliance may not have -- on the lab,
	// a no-agent VM never reports an address and every environment-variable
	// write hangs until the appliance's proxy gives up (120 s, answered as a
	// bare 405). Set it only where the appliance can SSH to the test VMs.
	EnvHVMNoAgent = "TF_VAR_testacc_morpheus_hvm_no_agent"

	// EnvHVMOtherClusterName names a second HVM cluster in the same cloud, for
	// the test that provisions into one cluster's pool with a network of the
	// other's and expects Morpheus to refuse.
	EnvHVMOtherClusterName = "TF_VAR_testacc_morpheus_hvm_other_cluster_name"

	// EnvHVMPoollessNetworkName names a network of the HVM cloud that belongs
	// to no resource pool (a docker bridge, say), for the test that expects
	// the provider to list the networks Morpheus would have accepted instead.
	EnvHVMPoollessNetworkName = "TF_VAR_testacc_morpheus_hvm_poolless_network_name"
)

// HVMCloudName returns the HVM cloud's name, skipping the test when
// EnvHVMCloudName is unset.
func HVMCloudName(t *testing.T) string {
	t.Helper()

	return envOrSkip(t, EnvHVMCloudName, "the Morpheus private cloud of an HVM cluster")
}

// HVMClusterName returns the HVM cluster's name, skipping the test when
// EnvHVMClusterName is unset.
func HVMClusterName(t *testing.T) string {
	t.Helper()

	return envOrSkip(t, EnvHVMClusterName, "an HVM cluster to provision into")
}

// HVMNetworkName returns the name of a network of the HVM cluster, skipping the
// test when EnvHVMNetworkName is unset.
func HVMNetworkName(t *testing.T) string {
	t.Helper()

	return envOrSkip(t, EnvHVMNetworkName, "a network of the HVM cluster")
}

// HVMDatastoreName returns the name of a datastore of the HVM cluster, skipping
// the test when EnvHVMDatastoreName is unset. Tests that can do without one
// should read EnvHVMDatastoreName themselves and fall back to automatic
// datastore selection rather than skip.
func HVMDatastoreName(t *testing.T) string {
	t.Helper()

	return envOrSkip(t, EnvHVMDatastoreName, "a datastore of the HVM cluster")
}

// HVMLayout returns the name and version of the layout to provision with,
// skipping the test when either of EnvHVMLayoutName and EnvHVMLayoutVersion is
// unset.
func HVMLayout(t *testing.T) (name, version string) {
	t.Helper()

	name = envOrSkip(t, EnvHVMLayoutName, "an HVM/KVM instance type layout")
	version = envOrSkip(t, EnvHVMLayoutVersion, "the version of that layout")

	return name, version
}

// HVMPlanName returns the name of the kvm service plan for the test VMs,
// skipping the test when EnvHVMPlanName is unset.
func HVMPlanName(t *testing.T) string {
	t.Helper()

	return envOrSkip(t, EnvHVMPlanName, "a kvm service plan for small test VMs")
}

// HVMOtherClusterName returns the name of a second HVM cluster in the HVM
// cloud, skipping the test when EnvHVMOtherClusterName is unset.
func HVMOtherClusterName(t *testing.T) string {
	t.Helper()

	return envOrSkip(t, EnvHVMOtherClusterName, "a second HVM cluster in the same cloud")
}

// HVMPoollessNetworkName returns the name of a network of the HVM cloud that
// belongs to no resource pool, skipping the test when
// EnvHVMPoollessNetworkName is unset.
func HVMPoollessNetworkName(t *testing.T) string {
	t.Helper()

	return envOrSkip(t, EnvHVMPoollessNetworkName, "a network of the HVM cloud that belongs to no resource pool")
}

// HVMNoAgent reports whether the test VMs are provisioned without the
// Morpheus agent (EnvHVMNoAgent == "true"); the default is with it.
func HVMNoAgent() bool {
	return os.Getenv(EnvHVMNoAgent) == "true"
}

// HVMGrowingImageName returns the name of the image the platform grows a
// minimum-size volume for, skipping the test when EnvHVMGrowingImageName is
// unset.
func HVMGrowingImageName(t *testing.T) string {
	t.Helper()

	return envOrSkip(t, EnvHVMGrowingImageName,
		"a virtual image whose raw size is just over its minimum disk (10 GB requested -> 11 GB)")
}
