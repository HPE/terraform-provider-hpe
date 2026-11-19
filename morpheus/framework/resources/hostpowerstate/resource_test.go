// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package hostpowerstate_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

func TestMain(m *testing.M) {
	code := testhelpers.TestMain(m)
	testhelpers.WriteMergedResults()
	os.Exit(code)
}

// TestAccMorpheusHostPowerStateResource creates an Omega stub host and takes it
// off then on. Destroying the power-state resource performs no power action; the
// host is destroyed with it.
func TestAccMorpheusHostPowerStateResource(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.BareMetalStub)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	cloudID := os.Getenv("TF_VAR_testacc_morpheus_bmaas_cloud_id")
	groupID := os.Getenv("TF_VAR_testacc_morpheus_bmaas_group_id")
	poolID := os.Getenv("TF_VAR_testacc_morpheus_bmaas_resource_pool_id")
	if cloudID == "" || groupID == "" || poolID == "" {
		t.Skip("TF_VAR_testacc_morpheus_bmaas_{cloud_id,group_id,resource_pool_id} must all be set")
	}

	name := acctest.RandomWithPrefix("TestAccMorpheusHostPower")
	providerConfig := testhelpers.ProviderBlock()

	hostBlock := fmt.Sprintf(`
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
}
`, name, cloudID, groupID, poolID)

	powerBlock := func(state string) string {
		return fmt.Sprintf(`
resource "hpe_morpheus_host_power_state" "test" {
  host_id     = hpe_morpheus_host.test.id
  power_state = %q
}
`, state)
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + hostBlock + powerBlock("off"),
				Check: resource.TestCheckResourceAttr(
					"hpe_morpheus_host_power_state.test", "power_state", "off"),
			},
			{
				Config: providerConfig + hostBlock + powerBlock("on"),
				Check: resource.TestCheckResourceAttr(
					"hpe_morpheus_host_power_state.test", "power_state", "on"),
			},
		},
	})
}
