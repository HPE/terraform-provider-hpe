package networkrouternat_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/networkrouternat"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/nsxt"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

func TestMain(m *testing.M) {
	code := testhelpers.TestMain(m)
	testhelpers.WriteMergedResults()
	os.Exit(code)
}

// NSX-T NAT rules attach to a gateway's policy path, so they require a realized
// tier-1 that is connected to a tier-0 and has an edge cluster. nsxt.Tier1Config
// renders that tier-1 by resolving the pre-provisioned tier-0 (and its group,
// network integration and edge cluster) by name.

func TestAccMorpheusNetworkRouterNatResourceExampleOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.NetworkRouter)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()
	name := acctest.RandomWithPrefix(t.Name())
	resourceName := "hpe_morpheus_network_router_nat.example"

	routerConfig := nsxt.Tier1Config(name)

	resourceConfig, err := networkrouternat.RenderNetworkRouterNatConfig(t, map[string]string{
		"RouterId": nsxt.Tier1RouterIDRef,
		"Name":     name,
	})
	if err != nil {
		t.Fatal(err)
	}

	checks := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttrPair(resourceName, "router_id", nsxt.Tier1Resource, "id"),
		resource.TestCheckResourceAttr(resourceName, "name", name),
		resource.TestCheckResourceAttr(resourceName, "source_network", "10.0.0.0/24"),
		resource.TestCheckResourceAttr(resourceName, "description", "Example SNAT rule"),
		resource.TestCheckResourceAttrSet(resourceName, "id"),
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
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
				// action and firewall are create-only inputs that the NAT read
				// (GET) API does not return, so they cannot round-trip through
				// import (verified: only these two differ after import).
				ImportStateVerifyIgnore: []string{"action", "firewall"},
				ResourceName:            "hpe_morpheus_network_router_nat.example",
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
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.NetworkRouter)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()
	name := acctest.RandomWithPrefix(t.Name())
	resourceName := "hpe_morpheus_network_router_nat.example"

	routerConfig := nsxt.Tier1Config(name)

	createConfig, err := networkrouternat.RenderNetworkRouterNatConfig(t, map[string]string{
		"RouterId": nsxt.Tier1RouterIDRef,
		"Name":     name,
	})
	if err != nil {
		t.Fatal(err)
	}

	updateConfig := `
resource "hpe_morpheus_network_router_nat" "example" {
  router_id          = ` + nsxt.Tier1RouterIDRef + `
  name               = "` + name + `"
  action             = "SNAT"
  source_network     = "10.1.0.0/24"
  translated_network = "192.168.1.1"
  description        = "Updated SNAT rule"
  protocol           = "tcp"
}
`

	createChecks := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttrPair(resourceName, "router_id", nsxt.Tier1Resource, "id"),
		resource.TestCheckResourceAttr(resourceName, "name", name),
		resource.TestCheckResourceAttr(resourceName, "source_network", "10.0.0.0/24"),
		resource.TestCheckResourceAttr(resourceName, "description", "Example SNAT rule"),
	)

	updateChecks := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttrPair(resourceName, "router_id", nsxt.Tier1Resource, "id"),
		resource.TestCheckResourceAttr(resourceName, "name", name),
		resource.TestCheckResourceAttr(resourceName, "source_network", "10.1.0.0/24"),
		resource.TestCheckResourceAttr(resourceName, "description", "Updated SNAT rule"),
		// protocol is deprecated and dropped by the API; verify the configured
		// value still round-trips (regression guard for the inconsistent-result
		// -after-apply defect).
		resource.TestCheckResourceAttr(resourceName, "protocol", "tcp"),
	)

	checkInPlaceUpdate := resource.ConfigPlanChecks{
		PreApply: []plancheck.PlanCheck{
			plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate),
		},
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{Config: providerConfig + routerConfig + createConfig, Check: createChecks},
			{Config: providerConfig + routerConfig + updateConfig, Check: updateChecks, ConfigPlanChecks: checkInPlaceUpdate},
			{Config: providerConfig + routerConfig + updateConfig, ExpectNonEmptyPlan: false, PlanOnly: true},
		},
	})
}
