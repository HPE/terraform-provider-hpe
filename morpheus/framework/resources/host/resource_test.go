// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package host_test

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

func TestMain(m *testing.M) {
	code := testhelpers.TestMain(m)
	testhelpers.WriteMergedResults()
	os.Exit(code)
}

// bmaasFixture reads the Omega bare-metal fixture IDs from the environment.
func bmaasFixture(t *testing.T) (cloudID, groupID, resourcePoolID string) {
	t.Helper()

	cloudID = os.Getenv("TF_VAR_testacc_morpheus_bmaas_cloud_id")
	groupID = os.Getenv("TF_VAR_testacc_morpheus_bmaas_group_id")
	resourcePoolID = os.Getenv("TF_VAR_testacc_morpheus_bmaas_resource_pool_id")

	if cloudID == "" || groupID == "" || resourcePoolID == "" {
		t.Skip("TF_VAR_testacc_morpheus_bmaas_{cloud_id,group_id,resource_pool_id} must all be set")
	}

	return cloudID, groupID, resourcePoolID
}

func omegaHostConfig(name, cloudID, groupID, poolID, description string) string {
	return omegaHostConfigWith(name, cloudID, groupID, poolID, fmt.Sprintf(`
  description         = %q
  labels              = ["tfacc", "omega"]
  ssh_username        = "tfacc"
  visibility          = "private"
`, description))
}

// omegaHostConfigWith is omegaHostConfig with the optional, in-place-updatable
// attributes supplied by the caller (possibly none).
func omegaHostConfigWith(name, cloudID, groupID, poolID, optional string) string {
	return omegaHostConfigWithPassword(name, cloudID, groupID, poolID, optional, "x")
}

// omegaHostConfigWithPassword is omegaHostConfigWith with the console password
// supplied, so a test can change the write-only value.
func omegaHostConfigWithPassword(name, cloudID, groupID, poolID, optional, consolePassword string) string {
	return fmt.Sprintf(`
resource "hpe_morpheus_host" "test" {
  name                = %q
  cloud_id            = %s
  group_id            = %s
  compute_server_type = "omega.baremetal.stub-server"
%s
  config_omega_bmaas = {
    console_host        = "10.0.0.1"
    console_username    = "admin"
    console_password_wo = %q
    num_nics            = 1
    resource_pool_id    = %q
  }
}
`, name, cloudID, groupID, optional, consolePassword, poolID)
}

// TestAccMorpheusHostResourceOmega creates an Omega stub host, reads it, updates
// description/labels in place, imports it, and destroys it (delete flags unset →
// Morpheus default).
func TestAccMorpheusHostResourceOmega(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.BareMetalStub)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	cloudID, groupID, poolID := bmaasFixture(t)
	name := acctest.RandomWithPrefix("TestAccMorpheusHost")

	providerConfig := testhelpers.ProviderBlock()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + omegaHostConfig(name, cloudID, groupID, poolID, "initial"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hpe_morpheus_host.test", "name", name),
					resource.TestCheckResourceAttr("hpe_morpheus_host.test", "compute_server_type", "omega.baremetal.stub-server"),
					resource.TestCheckResourceAttr("hpe_morpheus_host.test", "description", "initial"),
					resource.TestCheckResourceAttrSet("hpe_morpheus_host.test", "id"),
					resource.TestCheckResourceAttrSet("hpe_morpheus_host.test", "status"),
					resource.TestCheckResourceAttrSet("hpe_morpheus_host.test", "power_state"),
				),
			},
			{
				Config: providerConfig + omegaHostConfig(name, cloudID, groupID, poolID, "updated"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hpe_morpheus_host.test", "description", "updated"),
					resource.TestCheckResourceAttr("hpe_morpheus_host.test", "ssh_username", "tfacc"),
					resource.TestCheckResourceAttr("hpe_morpheus_host.test", "visibility", "private"),
				),
			},
			{
				// visibility is updated in place (the update handler binds it
				// even though the documented body omits it).
				Config: providerConfig + omegaHostConfigWith(name, cloudID, groupID, poolID, `
  description         = "updated"
  labels              = ["tfacc", "omega"]
  ssh_username        = "tfacc"
  visibility          = "public"
`),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("hpe_morpheus_host.test", plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.TestCheckResourceAttr("hpe_morpheus_host.test", "visibility", "public"),
			},
			{
				// Removing description, labels and ssh_username clears them on
				// the host; visibility falls back to its default.
				Config: providerConfig + omegaHostConfigWith(name, cloudID, groupID, poolID, ""),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("hpe_morpheus_host.test", plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckNoResourceAttr("hpe_morpheus_host.test", "description"),
					resource.TestCheckNoResourceAttr("hpe_morpheus_host.test", "ssh_username"),
					resource.TestCheckResourceAttr("hpe_morpheus_host.test", "labels.#", "0"),
					resource.TestCheckResourceAttr("hpe_morpheus_host.test", "visibility", "private"),
				),
			},
			{
				// The variant password is sent only when the host is added, so a
				// change to it replaces the host (the digest tracker marks only
				// the variant block in use -- the other, null, block must not be
				// named or Terraform rejects the plan).
				Config: providerConfig + omegaHostConfigWithPassword(name, cloudID, groupID, poolID, "", "y"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("hpe_morpheus_host.test", plancheck.ResourceActionReplace),
					},
				},
				Check: resource.TestCheckResourceAttr("hpe_morpheus_host.test", "name", name),
			},
			{
				ResourceName:      "hpe_morpheus_host.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"config_omega_bmaas",
					"delete_options",
					"timeouts",
					"write_only_revision",
				},
			},
		},
	})
}

// TestAccMorpheusHostResourceOmegaRemoveResources is the same create path but
// destroys with delete_options.remove_resources = true (the stub has no real
// metal to wipe).
func TestAccMorpheusHostResourceOmegaRemoveResources(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.BareMetalStub)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	cloudID, groupID, poolID := bmaasFixture(t)
	name := acctest.RandomWithPrefix("TestAccMorpheusHost")

	providerConfig := testhelpers.ProviderBlock()

	config := providerConfig + fmt.Sprintf(`
resource "hpe_morpheus_host" "test" {
  name                = %q
  cloud_id            = %s
  group_id            = %s
  compute_server_type = "omega.baremetal.stub-server"

  config_omega_bmaas = {
    console_host     = "10.0.0.1"
    console_username = "admin"
    num_nics         = 1
    resource_pool_id = %q
  }

  delete_options = {
    remove_resources = true
  }
}
`, name, cloudID, groupID, poolID)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.TestCheckResourceAttr(
					"hpe_morpheus_host.test", "delete_options.remove_resources", "true"),
			},
		},
	})
}

// TestAccMorpheusHostResourceNotCreatableType: an HVM hypervisor (mvm) cannot
// be added as a standalone host -- it is created by hpe_morpheus_cluster -- so
// the create is refused up front, naming the owning resource, instead of
// failing with the API's generic error.
func TestAccMorpheusHostResourceNotCreatableType(t *testing.T) {
	defer testhelpers.RecordResult(t)

	// The creatable check fails open, so if the appliance cannot serve the
	// type list this test would really POST an mvm host to the shared HVM
	// cloud. It therefore also needs the opt-in to host creation that the
	// Omega stub tests carry.
	capabilities.MustHaveOrSkip(t, capabilities.HVM, capabilities.BareMetalStub)
	cloudID := testhelpers.ComputeServerCloudID(t)
	groupID := os.Getenv("TF_VAR_testacc_morpheus_bmaas_group_id")
	if groupID == "" {
		t.Skip("TF_VAR_testacc_morpheus_bmaas_group_id must be set")
	}

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	config := testhelpers.ProviderBlock() + fmt.Sprintf(`
resource "hpe_morpheus_host" "hypervisor" {
  name                = "%s"
  cloud_id            = %s
  group_id            = %s
  compute_server_type = "mvm"
}
`, acctest.RandomWithPrefix(t.Name()), cloudID, groupID)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{{
			Config:      config,
			ExpectError: regexp.MustCompile(`(?s)Server type cannot be added as a host.*hpe_morpheus_cluster`),
		}},
	})
}

// TestAccMorpheusHostResourceOmegaVolumesAndNics creates an Omega stub host
// with a root volume (no id: Morpheus treats it as new), a network interface
// and user data; the first two are sent on the add and kept in state as
// configured (create-only, not read back), the user data is read back from
// the server's stored config, and the plan is empty afterwards. Verified
// by hand on the appliance: the server carries the 40 GiB root volume and one
// interface.
func TestAccMorpheusHostResourceOmegaVolumesAndNics(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.BareMetalStub)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	// Not parallel: the request recorder is process-wide.
	cloudID, groupID, poolID := bmaasFixture(t)
	name := acctest.RandomWithPrefix("TestAccMorpheusHost")
	rec := testhelpers.RecordRequests(t)

	cfg := testhelpers.ProviderBlock() + omegaHostConfigWith(name, cloudID, groupID, poolID, `
  user_data          = "#cloud-config\nhostname: tfacc\n"
  volumes            = [{ name = "root", root_volume = true, size = 40 }]
  network_interfaces = [{ ip_mode = "dhcp" }]
`)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				PreConfig: rec.Reset,
				Config:    cfg,
				Check: resource.ComposeAggregateTestCheckFunc(
					rec.ExpectSentKeys("/api/servers", map[string]bool{"volumes": true, "networkInterfaces": true}),
					resource.TestCheckResourceAttr("hpe_morpheus_host.test", "volumes.#", "1"),
					resource.TestCheckResourceAttr("hpe_morpheus_host.test", "volumes.0.name", "root"),
					resource.TestCheckResourceAttr("hpe_morpheus_host.test", "volumes.0.root_volume", "true"),
					resource.TestCheckResourceAttr("hpe_morpheus_host.test", "volumes.0.size", "40"),
					resource.TestCheckResourceAttr("hpe_morpheus_host.test", "network_interfaces.#", "1"),
					resource.TestCheckResourceAttr("hpe_morpheus_host.test", "network_interfaces.0.ip_mode", "dhcp"),
					// user_data is read back from the server's stored config.
					resource.TestCheckResourceAttr("hpe_morpheus_host.test", "user_data", "#cloud-config\nhostname: tfacc\n"),
				),
			},
			{Config: cfg, PlanOnly: true},
		},
	})
}
