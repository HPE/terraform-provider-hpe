package network_router_route_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/HPE/terraform-provider-hpe/morpheus"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/systemoverride"
)

func TestMain(m *testing.M) {
	systemoverride.ParseFlags()
	code := testhelpers.TestMain(m)
	testhelpers.WriteMergedResults()
	os.Exit(code)
}

func TestAccMorpheusNetworkRouterRouteBasic(t *testing.T) {
	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	routerID := os.Getenv("TF_ACC_MORPHEUS_ROUTER_ID")
	if routerID == "" {
		t.Skip("TF_ACC_MORPHEUS_ROUTER_ID not set")
	}

	t.Parallel()

	testSystem := systemoverride.GetPreferred(t, "zodiac")
	providerConfig := testhelpers.ProviderBlockForServer(testSystem)

	name := acctest.RandomWithPrefix("tf-acc-rtr-route")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + fmt.Sprintf(`
resource "hpe_morpheus_network_router_route" "test" {
  router_id   = %s
  name        = %q
  source      = "10.0.0.0/24"
  destination = "10.1.0.0/24"
}
`, routerID, name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hpe_morpheus_network_router_route.test", "name", name),
					resource.TestCheckResourceAttr("hpe_morpheus_network_router_route.test", "source", "10.0.0.0/24"),
					resource.TestCheckResourceAttr("hpe_morpheus_network_router_route.test", "destination", "10.1.0.0/24"),
					resource.TestCheckResourceAttrSet("hpe_morpheus_network_router_route.test", "id"),
				),
			},
			{
				ImportState:       true,
				ImportStateVerify: true,
				ResourceName:      "hpe_morpheus_network_router_route.test",
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources["hpe_morpheus_network_router_route.test"]
					if !ok {
						return "", fmt.Errorf("resource not found")
					}

					return rs.Primary.Attributes["router_id"] + "/" + rs.Primary.ID, nil
				},
			},
		},
	})
}

func TestAccMorpheusNetworkRouterRouteUpdate(t *testing.T) {
	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	routerID := os.Getenv("TF_ACC_MORPHEUS_ROUTER_ID")
	if routerID == "" {
		t.Skip("TF_ACC_MORPHEUS_ROUTER_ID not set")
	}

	t.Parallel()

	testSystem := systemoverride.GetPreferred(t, "zodiac")
	providerConfig := testhelpers.ProviderBlockForServer(testSystem)

	name := acctest.RandomWithPrefix("tf-acc-rtr-route")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + fmt.Sprintf(`
resource "hpe_morpheus_network_router_route" "test" {
  router_id   = %s
  name        = %q
  source      = "10.0.0.0/24"
  destination = "10.1.0.0/24"
  description = "initial"
}
`, routerID, name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hpe_morpheus_network_router_route.test", "description", "initial"),
				),
			},
			{
				Config: providerConfig + fmt.Sprintf(`
resource "hpe_morpheus_network_router_route" "test" {
  router_id   = %s
  name        = %q
  source      = "10.0.0.0/24"
  destination = "10.1.0.0/24"
  description = "updated"
}
`, routerID, name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hpe_morpheus_network_router_route.test", "description", "updated"),
				),
			},
		},
	})
}
