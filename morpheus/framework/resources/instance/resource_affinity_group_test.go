// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package instance_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

// checkServerInAffinityGroup asserts that at least one value in the instance's
// compute_servers set appears in the affinity group data source's servers set.
// This is the authoritative membership check — the affinity group is the source
// of truth for which servers are placed in it.
//
// Note this proves the payload was accepted and membership recorded. It does
// not prove which hypervisor host the instance landed on: the provider does not
// expose a compute server's parent host, so host placement cannot be asserted
// here. See exampleHCL/ for the manual host comparison.
func checkServerInAffinityGroup(instanceAddr, agDataSourceAddr string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		instRS, ok := s.RootModule().Resources[instanceAddr]
		if !ok {
			return fmt.Errorf("resource %s not found in state", instanceAddr)
		}

		agRS, ok := s.RootModule().Resources[agDataSourceAddr]
		if !ok {
			return fmt.Errorf("data source %s not found in state", agDataSourceAddr)
		}

		// Collect compute_servers from the instance (set stored as *.# + *.<hash>)
		computeServers := collectSetValues(instRS.Primary.Attributes, "compute_servers")
		if len(computeServers) == 0 {
			return fmt.Errorf("%s has no compute_servers in state", instanceAddr)
		}

		// Collect servers from the affinity group data source
		agServers := collectSetValues(agRS.Primary.Attributes, "servers")
		if len(agServers) == 0 {
			return fmt.Errorf(
				"%s has no servers in state — the affinity group is empty, "+
					"meaning the instance was NOT placed into it", agDataSourceAddr)
		}

		// Check that at least one compute server is in the affinity group
		for _, cs := range computeServers {
			for _, ags := range agServers {
				if cs == ags {
					return nil
				}
			}
		}

		return fmt.Errorf(
			"none of the instance's compute_servers %v appear in the affinity group's servers %v",
			computeServers, agServers)
	}
}

// collectSetValues extracts all values from a Terraform set attribute stored in
// the flat state map (e.g. "attr.#" = "2", "attr.12345" = "val").
func collectSetValues(attrs map[string]string, prefix string) []string {
	var vals []string

	for k, v := range attrs {
		if k == prefix+".#" || k == prefix+".%" {
			continue
		}
		if len(k) > len(prefix)+1 && k[:len(prefix)+1] == prefix+"." {
			vals = append(vals, v)
		}
	}

	return vals
}

// TestAccMorpheusInstanceResourceAffinityGroupVMware provisions a VMware instance
// with config_vmware.affinity_group_id and asserts the instance's compute server
// appears in the affinity group's server membership list. Asserting on the group's
// servers attribute is the only reliable proof of placement — the instance echoes
// the attribute back unchanged regardless of whether placement succeeded.
//
// The instance body mirrors example_vmware.tf.tmpl, the provider's proven VMware
// example, so a failure points at affinity group handling rather than unrelated
// provisioning settings.
func TestAccMorpheusInstanceResourceAffinityGroupVMware(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.VMware, capabilities.AffinityGroup)

	cloudID := testhelpers.AffinityCloudID(t)
	poolID := testhelpers.AffinityPoolID(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()
	name := acctest.RandomWithPrefix(t.Name())

	// The test provisions its own seed instance and creates its own affinity
	// group, rather than using a group maintained outside the test.
	//
	// It has to. On VMware a DRS rule requires two virtual machines, so
	// provisioning the first machine into an empty group builds a one-machine
	// rule that vCenter rejects and the instance fails at power on. Depending on
	// a shared group also makes the test fail confusingly if someone deletes the
	// server that happened to be seeding it.
	resourceConfig := fmt.Sprintf(`
data "hpe_morpheus_service_plan" "vmware_1cpu" {
  name                = "1 CPU, 1GB Memory"
  provision_type_code = "vmware"
}

data "hpe_morpheus_instance_type_layout" "vmware" {
  name    = "VMware VM"
  version = "22.04"
}

locals {
  vmware_volumes = [
    {
      root_volume              = true
      name                     = "root"
      size                     = 10
      storage_type_id          = 1
      datastore_auto_selection = "auto"
    }
  ]
  vmware_tags = [
    { name = "sweepable", value = "true" },
    { name = "managed_by", value = "terraform" },
  ]
  vmware_config = {
    resource_pool_id      = "pool-%[3]s"
    nested_virtualization = "off"
    no_agent              = true
    create_user           = false
    vmware_folder_id      = "group-v79"
  }
}

resource "hpe_morpheus_instance" "affinity_vmware_seed" {
  name             = "%[1]s-seed"
  cloud_id         = %[2]s
  layout_id        = data.hpe_morpheus_instance_type_layout.vmware.id
  instance_type_id = 9
  group_id         = 28
  plan_id          = data.hpe_morpheus_service_plan.vmware_1cpu.id
  instance_context = "dev"

  network_interfaces = [{ network_id = 86657 }]
  volumes            = local.vmware_volumes
  tags               = local.vmware_tags
  config_vmware      = local.vmware_config

  timeouts = {
    create = "1h"
    delete = "20m"
  }
}

resource "hpe_morpheus_cloud_affinity_group" "check_group" {
  cloud_id      = %[2]s
  pool_id       = %[3]s
  name          = "%[1]s"
  affinity_type = "KEEP_SEPARATE"
  active        = true

  servers = hpe_morpheus_instance.affinity_vmware_seed.compute_servers

  # Morpheus adds the instance under test to this group itself, at provision
  # time, in response to affinity_group_id. Without this the next apply would
  # treat that member as drift and remove it.
  lifecycle {
    ignore_changes = [servers]
  }
}

resource "hpe_morpheus_instance" "affinity_vmware" {
  name             = "%[1]s"
  cloud_id         = %[2]s
  layout_id        = data.hpe_morpheus_instance_type_layout.vmware.id
  instance_type_id = 9
  group_id         = 28
  plan_id          = data.hpe_morpheus_service_plan.vmware_1cpu.id
  instance_context = "dev"

  network_interfaces = [{ network_id = 86657 }]
  volumes            = local.vmware_volumes
  tags               = local.vmware_tags

  config_vmware = merge(local.vmware_config, {
    affinity_group_id = hpe_morpheus_cloud_affinity_group.check_group.id
  })

  timeouts = {
    create = "1h"
    delete = "20m"
  }
}

# Read the affinity group AFTER the instance is created to check membership.
data "hpe_morpheus_cloud_affinity_group" "check" {
  cloud_id = %[2]s
  id       = hpe_morpheus_cloud_affinity_group.check_group.id

  depends_on = [hpe_morpheus_instance.affinity_vmware]
}
`, name, cloudID, poolID)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + resourceConfig,
				Check: checkServerInAffinityGroup(
					"hpe_morpheus_instance.affinity_vmware",
					"data.hpe_morpheus_cloud_affinity_group.check",
				),
			},
		},
	})
}

// TestAccMorpheusInstanceResourceAffinityGroupHVM provisions an HVM instance
// with config_hvm.affinity_group_id and asserts the instance's compute server
// appears in the affinity group's server membership list.
//
// NOTE on what this proves. HVM records membership unconditionally, so this
// test confirms the affinityGroup/affinityGroupId payload was sent and
// accepted. It does not prove host placement: that additionally requires the
// cluster to have dynamic placement enabled, a group with existing members, and
// a comparison against the host the scheduler would have chosen anyway.
// exampleHCL/hvm covers that end to end.
// affinityHVMConfig renders an HVM instance placed in a cluster affinity
// group, everything resolved by name: the cloud, cluster (and so the resource
// pool), layout (and so the instance type), network, plan and image from the
// HVM fixture (testhelpers/hvm_fixture.go), and the affinity group by name
// within that cluster (TF_VAR_testacc_morpheus_hvm_affinity_group_name).
// affinityGroupLine is the config_hvm line placing the instance, or empty for
// an instance in no group; withCheck adds a read of the group after the
// instance exists, for the membership assertion.
func affinityHVMConfig(t *testing.T, label, name, affinityGroupLine string, withCheck bool) string {
	t.Helper()

	cloudName := testhelpers.HVMCloudName(t)
	clusterName := testhelpers.HVMClusterName(t)
	networkName := testhelpers.HVMNetworkName(t)
	planName := testhelpers.HVMPlanName(t)
	layoutName, layoutVersion := testhelpers.HVMLayout(t)
	groupName := testhelpers.HVMAffinityGroupName(t)

	imageData, imageLine := "", ""
	if image := os.Getenv(testhelpers.EnvHVMImageName); image != "" {
		imageData = `
data "hpe_morpheus_image" "hvm" {
  name = "` + image + `"
}
`
		imageLine = "\n    image_id          = data.hpe_morpheus_image.hvm.id"
	}

	check := ""
	if withCheck {
		check = `
# Read the cluster affinity group AFTER the instance is created.
data "hpe_morpheus_cluster_affinity_group" "check" {
  cluster_id = data.hpe_morpheus_cluster.hvm.id
  id         = data.hpe_morpheus_cluster_affinity_group.target.id

  depends_on = [hpe_morpheus_instance.` + label + `]
}
`
	}

	return `
data "hpe_morpheus_cloud" "hvm_cloud" {
  name = "` + cloudName + `"
}

data "hpe_morpheus_cluster" "hvm" {
  name = "` + clusterName + `"
}

data "hpe_morpheus_instance_type_layout" "hvm" {
  name    = "` + layoutName + `"
  version = "` + layoutVersion + `"
}

data "hpe_morpheus_network" "hvm" {
  name       = "` + networkName + `"
  cluster_id = data.hpe_morpheus_cluster.hvm.id
}

data "hpe_morpheus_service_plan" "hvm" {
  name                = "` + planName + `"
  provision_type_code = "kvm"
}
` + imageData + `
# The affinity group to place the instance in, by name within the cluster.
data "hpe_morpheus_cluster_affinity_group" "target" {
  cluster_id = data.hpe_morpheus_cluster.hvm.id
  name       = "` + groupName + `"
}

# MORPH-16458: create a per-run group in the current tenant instead of pinning
# group_id = 1 (the master group), which a sub-tenant cannot use (403).
resource "hpe_morpheus_group" "test" {
  name      = "` + name + `-grp"
  cloud_ids = [data.hpe_morpheus_cloud.hvm_cloud.id]
}

resource "hpe_morpheus_instance" "` + label + `" {
  name             = "` + name + `"
  cloud_id         = data.hpe_morpheus_cloud.hvm_cloud.id
  layout_id        = data.hpe_morpheus_instance_type_layout.hvm.id
  instance_type_id = data.hpe_morpheus_instance_type_layout.hvm.instance_type.id
  group_id         = hpe_morpheus_group.test.id
  plan_id          = data.hpe_morpheus_service_plan.hvm.id

  instance_context = "dev"

  network_interfaces = [
    {
      network_id = data.hpe_morpheus_network.hvm.id
    }
  ]

  volumes = [
    {
      root_volume              = true
      name                     = "root"
      size                     = 10
      storage_type_id          = 1
      datastore_auto_selection = "auto"
    }
  ]

  tags = [
    {
      name  = "sweepable"
      value = "true"
    },
    {
      name  = "managed_by"
      value = "terraform"
    }
  ]

  config_hvm = {
    resource_pool_id  = "pool-${data.hpe_morpheus_cluster.hvm.permissions.resource_pool.id}"
    no_agent          = ` + fmt.Sprintf("%t", testhelpers.HVMNoAgent()) + `
    create_user       = false` + imageLine + `
` + affinityGroupLine + `
  }

  timeouts = {
    create = "1h"
    delete = "20m"
  }
}
` + check
}

const affinityGroupLine = "    affinity_group_id = data.hpe_morpheus_cluster_affinity_group.target.id"

func TestAccMorpheusInstanceResourceAffinityGroupHVM(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.HVM, capabilities.AffinityGroup)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()
	name := acctest.RandomWithPrefix(t.Name())

	resourceConfig := affinityHVMConfig(t, "affinity_hvm", name, affinityGroupLine, true)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + resourceConfig,
				Check: checkServerInAffinityGroup(
					"hpe_morpheus_instance.affinity_hvm",
					"data.hpe_morpheus_cluster_affinity_group.check",
				),
			},
		},
	})
}

// TestAccMorpheusInstanceResourceAffinityGroupRequiresReplace verifies that
// changing config_hvm.affinity_group_id triggers a destroy-before-create
// replacement, matching the RequiresReplace plan modifier on the attribute.
//
// The second step removes the attribute rather than pointing it at a second
// group. That exercises the same plan modifier, is a realistic thing for a user
// to do, and avoids requiring a second affinity group to exist purely for the
// test. The step applies for real: a plan-only step cannot carry a PreApply
// plan check.
func TestAccMorpheusInstanceResourceAffinityGroupRequiresReplace(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.HVM, capabilities.AffinityGroup)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()
	name := acctest.RandomWithPrefix(t.Name())

	resourceName := "hpe_morpheus_instance.replace_test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + affinityHVMConfig(t, "replace_test", name, affinityGroupLine, false),
			},
			{
				// Dropping affinity_group_id must replace the instance, not
				// update it in place.
				Config: providerConfig + affinityHVMConfig(t, "replace_test", name, "", false),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(
							resourceName, plancheck.ResourceActionDestroyBeforeCreate,
						),
					},
				},
			},
		},
	})
}

// TestAccMorpheusInstanceResourceAffinityGroupImport verifies that an instance
// created with config_hvm.affinity_group_id can be imported and that the
// imported state carries the same affinity_group_id, so a subsequent plan is
// clean.
//
// This is the case that makes reading affinity_group_id back from the API
// matter: without it, every imported instance would plan a replacement.
func TestAccMorpheusInstanceResourceAffinityGroupImport(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.HVM, capabilities.AffinityGroup)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()
	name := acctest.RandomWithPrefix(t.Name())

	resourceConfig := affinityHVMConfig(t, "import_test", name, affinityGroupLine, false)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + resourceConfig,
			},
			{
				// The import must round-trip affinity_group_id. The ignored
				// attributes are pre-existing import gaps on this resource and
				// are not what this test is about.
				ResourceName:      "hpe_morpheus_instance.import_test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"compute_servers",
					"config",
					"connection_info",
					"labels",
					"network_interfaces",
					"timeouts",
					"volumes",
				},
			},
		},
	})
}
