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

// --- env var helpers (inlined because testhelpers/config.go is outside edit scope) ---

const (
	envAffinityGroupCloudAGID   = "TF_VAR_testacc_morpheus_affinity_group_cloud_ag_id"
	envAffinityGroupClusterAGID = "TF_VAR_testacc_morpheus_affinity_group_cluster_ag_id"
)

// affinityGroupCloudAGID returns a cloud affinity group to provision into.
//
// The group must already contain at least one compute server. vSphere DRS
// rules require two virtual machines, so provisioning the first machine into an
// empty group builds a one-VM rule that vCenter rejects; Morpheus reports
// "Error running vm" at Power On and fails the provision. An empty group
// therefore fails this test for reasons that have nothing to do with the
// provider.
func affinityGroupCloudAGID(t *testing.T) string {
	t.Helper()

	v := os.Getenv(envAffinityGroupCloudAGID)
	if v == "" {
		t.Skip(envAffinityGroupCloudAGID +
			" not set; skipping test requiring a cloud affinity group " +
			"(which must already have at least one member) for VMware instance placement")
	}

	return v
}

func affinityGroupClusterAGID(t *testing.T) string {
	t.Helper()

	v := os.Getenv(envAffinityGroupClusterAGID)
	if v == "" {
		t.Skip(envAffinityGroupClusterAGID +
			" not set; skipping test requiring a cluster affinity group for HVM instance placement")
	}

	return v
}

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
	agID := affinityGroupCloudAGID(t)
	poolID := testhelpers.AffinityPoolID(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()
	name := acctest.RandomWithPrefix(t.Name())

	resourceConfig := fmt.Sprintf(`
data "hpe_morpheus_service_plan" "vmware_1cpu" {
  name                = "1 CPU, 1GB Memory"
  provision_type_code = "vmware"
}

data "hpe_morpheus_instance_type_layout" "vmware" {
  name    = "VMware VM"
  version = "22.04"
}

resource "hpe_morpheus_instance" "affinity_vmware" {
  name             = "%[1]s"
  cloud_id         = %[2]s
  layout_id        = data.hpe_morpheus_instance_type_layout.vmware.id
  instance_type_id = 9
  group_id         = 28
  plan_id          = data.hpe_morpheus_service_plan.vmware_1cpu.id

  instance_context = "dev"

  network_interfaces = [
    {
      network_id = 86657
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

  config_vmware = {
    resource_pool_id      = "pool-%[3]s"
    nested_virtualization = "off"
    no_agent              = true
    create_user           = false
    vmware_folder_id      = "group-v79"
    affinity_group_id     = %[4]s
  }

  timeouts = {
    create = "1h"
    delete = "20m"
  }
}

# Read the affinity group AFTER the instance is created to check membership.
data "hpe_morpheus_cloud_affinity_group" "check" {
  cloud_id = %[2]s
  id       = %[4]s

  depends_on = [hpe_morpheus_instance.affinity_vmware]
}
`, name, cloudID, poolID, agID)

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
func TestAccMorpheusInstanceResourceAffinityGroupHVM(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.HVM, capabilities.AffinityGroup)

	clusterID := testhelpers.AffinityClusterID(t)
	agID := affinityGroupClusterAGID(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()
	name := acctest.RandomWithPrefix(t.Name())

	resourceConfig := fmt.Sprintf(`
data "hpe_morpheus_cloud" "hvm_cloud" {
  name = "QA HVM"
}

resource "hpe_morpheus_instance" "affinity_hvm" {
  name             = "%[1]s"
  cloud_id         = data.hpe_morpheus_cloud.hvm_cloud.id
  layout_id        = 5385
  instance_type_id = 9
  group_id         = 1
  plan_id          = 176

  instance_context = "dev"

  network_interfaces = [
    {
      network_id = 235699
    }
  ]

  volumes = [
    {
      root_volume              = true
      name                     = "root"
      size                     = 80
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
    resource_pool_id  = "pool-153047"
    no_agent          = true
    create_user       = false
    affinity_group_id = %[2]s
  }

  timeouts = {
    create = "1h"
    delete = "20m"
  }
}

# Read the cluster affinity group AFTER the instance is created.
data "hpe_morpheus_cluster_affinity_group" "check" {
  cluster_id = %[3]s
  id         = %[2]s

  depends_on = [hpe_morpheus_instance.affinity_hvm]
}
`, name, agID, clusterID)

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

	_ = testhelpers.AffinityClusterID(t) // ensures env is set
	agID := affinityGroupClusterAGID(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()
	name := acctest.RandomWithPrefix(t.Name())

	// affinityGroupLine is either the attribute assignment or nothing at all,
	// so the second step drops the attribute entirely.
	makeConfig := func(affinityGroupLine string) string {
		return fmt.Sprintf(`
data "hpe_morpheus_cloud" "hvm_cloud" {
  name = "QA HVM"
}

resource "hpe_morpheus_instance" "replace_test" {
  name             = "%[1]s"
  cloud_id         = data.hpe_morpheus_cloud.hvm_cloud.id
  layout_id        = 5385
  instance_type_id = 9
  group_id         = 1
  plan_id          = 176

  instance_context = "dev"

  network_interfaces = [
    {
      network_id = 235699
    }
  ]

  volumes = [
    {
      root_volume              = true
      name                     = "root"
      size                     = 80
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
    resource_pool_id  = "pool-153047"
    no_agent          = true
    create_user       = false
%[2]s
  }

  timeouts = {
    create = "1h"
    delete = "20m"
  }
}
`, name, affinityGroupLine)
	}

	resourceName := "hpe_morpheus_instance.replace_test"

	withGroup := fmt.Sprintf("    affinity_group_id = %s", agID)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + makeConfig(withGroup),
			},
			{
				// Dropping affinity_group_id must replace the instance, not
				// update it in place.
				Config: providerConfig + makeConfig(""),
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
