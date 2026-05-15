// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkrouterroute_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

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
					// TODO: verify - these may fail (may not be set with nsxt router routes)....
					resource.TestCheckResourceAttrSet(resourceName, "route_type"),
					resource.TestCheckResourceAttrSet(resourceName, "source_type"),
					resource.TestCheckResourceAttrSet(resourceName, "external_id"),
					resource.TestCheckResourceAttrSet(resourceName, "provider_id"),
				),
			},
		},
	})
}

func TestAccMorpheusNetworkRouterRouteUpdate(t *testing.T) {
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
		"Network":              "10.0.0.0/24",
		"NextHop":              "10.0.0.1",
	})
	if err != nil {
		t.Fatalf("failed to render router config: %s", err)
	}

	createRouteConfig, err := networkrouterroute.RenderRouteConfig(t, map[string]string{
		"RouterId":    "hpe_morpheus_network_router.example.id",
		"Name":        name,
		"Description": name,
		"Network":     "10.0.0.0/24",
		"NextHop":     "10.0.0.1",
		"Mtu":         "1500",
		"Enabled":     "true",
	})
	if err != nil {
		t.Fatalf("failed to render create route config: %s", err)
	}

	updateRouteConfig, err := networkrouterroute.RenderRouteConfig(t, map[string]string{
		"RouterId":    "hpe_morpheus_network_router.example.id",
		"Name":        name + "-updated",
		"Description": name + " updated",
		"Network":     "172.16.0.0/16",
		"NextHop":     "172.16.0.1",
		"Mtu":         "9000",
		"Enabled":     "false",
	})
	if err != nil {
		t.Fatalf("failed to render update route config: %s", err)
	}

	resourceName := "hpe_morpheus_network_router_route.example"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + routerConfig + createRouteConfig,
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
				),
			},
			{
				Config: providerConfig + routerConfig + updateRouteConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "router_id"),
					resource.TestCheckResourceAttr(resourceName, "name", name+"-updated"),
					resource.TestCheckResourceAttr(resourceName, "description", name+" updated"),
					resource.TestCheckResourceAttr(resourceName, "network", "172.16.0.0/16"),
					resource.TestCheckResourceAttr(resourceName, "next_hop", "172.16.0.1"),
					resource.TestCheckResourceAttr(resourceName, "mtu", "9000"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "default_route", "false"),
				),
			},
		},
	})
}
