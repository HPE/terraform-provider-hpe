// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkrouterroute_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/HPE/terraform-provider-hpe/morpheus"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/networkrouter"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/networkrouterroute"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/systemoverride"
)

func TestMain(m *testing.M) {
	systemoverride.ParseFlags()

	code := testhelpers.TestMain(m)

	testhelpers.WriteMergedResults()

	os.Exit(code)
}

func TestAccMorpheusNetworkRouterRouteExampleOk(t *testing.T) {
	t.Parallel()

	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())

	routerConfig, err := networkrouter.RenderNetworkRouterNSXTGatewayTier0Config(t, map[string]string{
		"Name":                 name,
		"GroupId":              "3",
		"NetworkIntegrationId": "5",
	})
	if err != nil {
		t.Fatalf("failed to render router config: %s", err)
	}

	routeConfig, err := networkrouterroute.RenderRouteConfig(t, map[string]string{
		"RouterId":    "hpe_morpheus_network_router.example.id",
		"Name":        name,
		"Description": name,
		"Network":     "10.0.0.0/24",
		"NextHop":     "10.0.0.1",
	})
	if err != nil {
		t.Fatalf("failed to render route config: %s", err)
	}

	resourceName := "hpe_morpheus_network_router_route.example"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + routerConfig + routeConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "router_id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", name),
					resource.TestCheckResourceAttr(resourceName, "network", "10.0.0.0/24"),
					resource.TestCheckResourceAttr(resourceName, "next_hop", "10.0.0.1"),
					resource.TestCheckResourceAttr(resourceName, "mtu", "1500"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "default_route", "false"),
					resource.TestCheckResourceAttrSet(resourceName, "route_type"),
					resource.TestCheckResourceAttrSet(resourceName, "source_type"),
					resource.TestCheckResourceAttrSet(resourceName, "external_id"),
					resource.TestCheckResourceAttrSet(resourceName, "provider_id"),
				),
			},
			{
				Config:             providerConfig + routerConfig + routeConfig,
				ExpectNonEmptyPlan: false,
				PlanOnly:           true,
			},
			{
				ResourceName: resourceName,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources["hpe_morpheus_network_router_route.example"]
					if !ok {
						return "", fmt.Errorf("resource not found")
					}

					return rs.Primary.Attributes["router_id"] + "." + rs.Primary.Attributes["id"], nil
				},
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
