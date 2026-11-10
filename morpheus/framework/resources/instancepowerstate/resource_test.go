// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package instancepowerstate_test

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
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

// testCheckServersPowerState asserts, through the API, that every server
// behind the managed instance reports the given powerState. The instance
// status alone does not prove a power transition: Morpheus also reports an
// instance "stopped" after a service-only stop that leaves the VM running.
func testCheckServersPowerState(resourceName, want string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource %s not found in state", resourceName)
		}

		instanceID, err := strconv.ParseInt(rs.Primary.Attributes["instance_id"], 10, 64)
		if err != nil {
			return fmt.Errorf("%s: instance_id %q: %w", resourceName, rs.Primary.Attributes["instance_id"], err)
		}

		ctx := context.Background()

		client, err := testhelpers.NewClientForServer(ctx, "")
		if err != nil {
			return err
		}

		inst, hresp, err := client.InstancesAPI.GetInstance(ctx, instanceID).Execute()
		if err != nil {
			return fmt.Errorf("GET instance %d: %w (%v)", instanceID, err, hresp)
		}

		if inst.Instance == nil || len(inst.Instance.Servers) == 0 {
			return fmt.Errorf("instance %d: no servers to check", instanceID)
		}

		for _, serverID := range inst.Instance.Servers {
			host, hresp, err := client.HostsAPI.GetHost(ctx, sdk.GetHostIdParameter{Int64: &serverID}).Execute()
			if err != nil {
				return fmt.Errorf("GET server %d: %w (%v)", serverID, err, hresp)
			}

			if host.Server == nil || host.Server.PowerState == nil {
				return fmt.Errorf("server %d: no powerState in response", serverID)
			}

			if got := *host.Server.PowerState; got != want {
				return fmt.Errorf("server %d: powerState = %q, want %q", serverID, got, want)
			}
		}

		return nil
	}
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
					testCheckServersPowerState("hpe_morpheus_instance_power_state.test", "on"),
				),
			},
			// Update: running → stopped. The server must really be powered
			// off, not merely have had its service stopped.
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
					testCheckServersPowerState("hpe_morpheus_instance_power_state.test", "off"),
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
					testCheckServersPowerState("hpe_morpheus_instance_power_state.test", "on"),
				),
			},
		},
	})
}
