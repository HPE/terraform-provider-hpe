// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package instancepowerstate_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/instance"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

func TestMain(m *testing.M) {
	code := testhelpers.TestMain(m)
	testhelpers.WriteMergedResults()
	os.Exit(code)
}

// TestAccMorpheusInstancePowerStateResourceRunning takes an instance through
// running -> stopped -> running. The instance is the test's own: it is
// provisioned from the HVM fixture (testhelpers/hvm_fixture.go, every object
// named by environment variables), so no existing instance is stopped -- a
// borrowed one would be a gateway, a fixture or somebody's workload -- and
// the test skips cleanly where there is no HVM cluster to provision into.
// Destroying the power-state resource performs no power action; the instance
// is destroyed with it.
func TestAccMorpheusInstancePowerStateResourceRunning(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All, capabilities.HVM)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	// A short instance name: Morpheus derives the guest's hostname from it,
	// and a test name with a random suffix is longer than a hostname label.
	instanceConfig, err := instance.RenderInstanceHVMConfig(t, map[string]string{
		"Name": acctest.RandomWithPrefix("TestAccMorpheusPower"),
	})
	if err != nil {
		t.Fatal(err)
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()

	const instanceIDRef = "hpe_morpheus_instance.example.id"

	resourceConfig := fmt.Sprintf(`
resource "hpe_morpheus_instance_power_state" "test" {
  instance_id   = %s
  desired_state = "running"
}
`, instanceIDRef)

	resourceConfigStopped := fmt.Sprintf(`
resource "hpe_morpheus_instance_power_state" "test" {
  instance_id   = %s
  desired_state = "stopped"
}
`, instanceIDRef)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + instanceConfig + resourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						"hpe_morpheus_instance_power_state.test",
						"desired_state",
						"running",
					),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_instance_power_state.test",
						"current_state",
						"running",
					),
					resource.TestCheckResourceAttrPair(
						"hpe_morpheus_instance_power_state.test", "instance_id",
						"hpe_morpheus_instance.example", "id",
					),
				),
			},
			// Update: running → stopped
			{
				Config: providerConfig + instanceConfig + resourceConfigStopped,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						"hpe_morpheus_instance_power_state.test",
						"desired_state",
						"stopped",
					),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_instance_power_state.test",
						"current_state",
						"stopped",
					),
				),
			},
			// Re-start to clean state
			{
				Config: providerConfig + instanceConfig + resourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						"hpe_morpheus_instance_power_state.test",
						"current_state",
						"running",
					),
				),
			},
		},
	})
}
