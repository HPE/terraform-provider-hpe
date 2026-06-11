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
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/networkrouter"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
)

func TestMain(m *testing.M) {
	code := testhelpers.TestMain(m)
	testhelpers.WriteMergedResults()
	os.Exit(code)
}

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

	routerConfig, err := networkrouter.RenderNetworkRouterGenericConfig(t, map[string]string{
		"Name":                 name + "-router",
		"TypeId":               "9",
		"GroupId":              "3",
		"NetworkIntegrationId": "5",
	})
	if err != nil {
		t.Fatal(err)
	}

	resourceConfig, err := network_router_nat.RenderNetworkRouterNatConfig(t, map[string]string{
		"RouterId": "hpe_morpheus_network_router.example.id",
		"Name":     name,
	})
	if err != nil {
		t.Fatal(err)
	}

	checks := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttrPair(resourceName, "router_id", "hpe_morpheus_network_router.example", "id"),
		resource.TestCheckResourceAttr(resourceName, "name", name),
		resource.TestCheckResourceAttr(resourceName, "source_network", "10.0.0.0/24"),
		resource.TestCheckResourceAttr(resourceName, "description", "Example SNAT rule"),
		resource.TestCheckResourceAttrSet(resourceName, "id"),
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + routerConfig + resourceConfig,
				Check:  checks,
			},
			{
				Config:             providerConfig + routerConfig + resourceConfig,
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

	routerConfig, err := networkrouter.RenderNetworkRouterGenericConfig(t, map[string]string{
		"Name":                 name + "-router",
		"TypeId":               "9",
		"GroupId":              "3",
		"NetworkIntegrationId": "5",
	})
	if err != nil {
		t.Fatal(err)
	}

	createConfig, err := network_router_nat.RenderNetworkRouterNatConfig(t, map[string]string{
		"RouterId": "hpe_morpheus_network_router.example.id",
		"Name":     name,
	})
	if err != nil {
		t.Fatal(err)
	}

	updateConfig := `
resource "hpe_morpheus_network_router_nat" "example" {
  router_id          = hpe_morpheus_network_router.example.id
  name               = "` + name + `"
  action             = "SNAT"
  source_network     = "10.1.0.0/24"
  translated_network = "192.168.1.1"
  description        = "Updated SNAT rule"
}
`

	createChecks := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttrPair(resourceName, "router_id", "hpe_morpheus_network_router.example", "id"),
		resource.TestCheckResourceAttr(resourceName, "name", name),
		resource.TestCheckResourceAttr(resourceName, "source_network", "10.0.0.0/24"),
		resource.TestCheckResourceAttr(resourceName, "description", "Example SNAT rule"),
	)

	updateChecks := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttrPair(resourceName, "router_id", "hpe_morpheus_network_router.example", "id"),
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
			{Config: providerConfig + routerConfig + createConfig, Check: createChecks},
			{Config: providerConfig + routerConfig + updateConfig, Check: updateChecks, ConfigPlanChecks: checkInPlaceUpdate},
			{Config: providerConfig + routerConfig + updateConfig, ExpectNonEmptyPlan: false, PlanOnly: true},
		},
	})
}
