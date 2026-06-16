// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package instancepowerstate_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	morpheus "github.com/HPE/terraform-provider-hpe/morpheus"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
)

func TestMain(m *testing.M) {
	code := testhelpers.TestMain(m)
	testhelpers.WriteMergedResults()
	os.Exit(code)
}

// TestAccMorpheusInstancePowerStateResourceRunning tests transitioning an
// instance to the running state. Requires TF_VAR_instance_id to reference
// an existing instance.
func TestAccMorpheusInstancePowerStateResourceRunning(t *testing.T) {
	if capabilities.Missing(t, capabilities.All) {
		t.Log("Skipping test due to missing capabilities")

		return
	}
	defer testhelpers.RecordResult(t)
	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	instanceID := os.Getenv("TF_VAR_instance_id")
	if instanceID == "" {
		t.Skip("TF_VAR_instance_id not set; skipping acceptance test")
	}

	providerConfig := testhelpers.ProviderBlock()

	resourceConfig := fmt.Sprintf(`
resource "hpe_morpheus_instance_power_state" "test" {
  instance_id   = %s
  desired_state = "running"
}
`, instanceID)

	resourceConfigStopped := fmt.Sprintf(`
resource "hpe_morpheus_instance_power_state" "test" {
  instance_id   = %s
  desired_state = "stopped"
}
`, instanceID)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + resourceConfig,
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
					resource.TestCheckResourceAttr(
						"hpe_morpheus_instance_power_state.test",
						"instance_id",
						instanceID,
					),
				),
			},
			// Update: running → stopped
			{
				Config: providerConfig + resourceConfigStopped,
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
				Config: providerConfig + resourceConfig,
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
