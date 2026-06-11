package network_router_nat_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/HPE/terraform-provider-hpe/morpheus"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/network_router_nat"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
)

func TestMain(m *testing.M) {
	code := testhelpers.TestMain(m)
	testhelpers.WriteMergedResults()
	os.Exit(code)
}

// existingTier1RouterID is a pre-provisioned, fully-realized NSX-T tier-1 gateway
// (connected to a tier-0, with an associated edge cluster) on integration 5.
// NSX-T NAT rules attach to the gateway's policy path, so they require a realized
// tier-1; creating one per test races NSX-T realization and a standalone tier-1
// (not connected to a tier-0) is rejected, so we reference this existing gateway.
const existingTier1RouterID = "27"

func TestAccMorpheusNetworkRouterNatResourceExampleOk(t *testing.T) {
	if capabilities.Missing(t, capabilities.NetworkRouter) {
		t.Log("Skipping test due to missing capabilities")

		return
	}
	defer testhelpers.RecordResult(t)
	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()
	name := acctest.RandomWithPrefix(t.Name())
	resourceName := "hpe_morpheus_network_router_nat.example"

	resourceConfig, err := network_router_nat.RenderNetworkRouterNatConfig(t, map[string]string{
		"RouterId": existingTier1RouterID,
		"Name":     name,
	})
	if err != nil {
		t.Fatal(err)
	}

	checks := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr(resourceName, "router_id", existingTier1RouterID),
		resource.TestCheckResourceAttr(resourceName, "name", name),
		resource.TestCheckResourceAttr(resourceName, "source_network", "10.0.0.0/24"),
		resource.TestCheckResourceAttr(resourceName, "description", "Example SNAT rule"),
		resource.TestCheckResourceAttrSet(resourceName, "id"),
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + resourceConfig,
				Check:  checks,
			},
			{
				Config:             providerConfig + resourceConfig,
				ExpectNonEmptyPlan: false,
				PlanOnly:           true,
			},
			{
				ImportState:       true,
				ImportStateVerify: true,
				ResourceName:      "hpe_morpheus_network_router_nat.example",
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources["hpe_morpheus_network_router_nat.example"]
					if !ok {
						return "", fmt.Errorf("resource not found")
					}

					return rs.Primary.Attributes["router_id"] + "." + rs.Primary.Attributes["id"], nil
				},
			},
		},
	})
}

func TestAccMorpheusNetworkRouterNatResourceUpdateOk(t *testing.T) {
	if capabilities.Missing(t, capabilities.NetworkRouter) {
		t.Log("Skipping test due to missing capabilities")

		return
	}
	defer testhelpers.RecordResult(t)
	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()
	name := acctest.RandomWithPrefix(t.Name())
	resourceName := "hpe_morpheus_network_router_nat.example"

	createConfig, err := network_router_nat.RenderNetworkRouterNatConfig(t, map[string]string{
		"RouterId": existingTier1RouterID,
		"Name":     name,
	})
	if err != nil {
		t.Fatal(err)
	}

	updateConfig := `
resource "hpe_morpheus_network_router_nat" "example" {
  router_id          = ` + existingTier1RouterID + `
  name               = "` + name + `"
  action             = "SNAT"
  source_network     = "10.1.0.0/24"
  translated_network = "192.168.1.1"
  description        = "Updated SNAT rule"
}
`

	createChecks := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr(resourceName, "router_id", existingTier1RouterID),
		resource.TestCheckResourceAttr(resourceName, "name", name),
		resource.TestCheckResourceAttr(resourceName, "source_network", "10.0.0.0/24"),
		resource.TestCheckResourceAttr(resourceName, "description", "Example SNAT rule"),
	)

	updateChecks := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr(resourceName, "router_id", existingTier1RouterID),
		resource.TestCheckResourceAttr(resourceName, "name", name),
		resource.TestCheckResourceAttr(resourceName, "source_network", "10.1.0.0/24"),
		resource.TestCheckResourceAttr(resourceName, "description", "Updated SNAT rule"),
	)

	checkInPlaceUpdate := resource.ConfigPlanChecks{
		PreApply: []plancheck.PlanCheck{
			plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate),
		},
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{Config: providerConfig + createConfig, Check: createChecks},
			{Config: providerConfig + updateConfig, Check: updateChecks, ConfigPlanChecks: checkInPlaceUpdate},
			{Config: providerConfig + updateConfig, ExpectNonEmptyPlan: false, PlanOnly: true},
		},
	})
}
