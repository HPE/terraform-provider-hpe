// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkrouterroute_test

import (
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/networkrouterroute"
	routeresource "github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/network_router_route"
	networkrouterresource "github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/networkrouter"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
)

func TestMain(m *testing.M) {
	code := m.Run()
	testhelpers.WriteMergedResults()
	os.Exit(code)
}

const providerConfigOffline = `
provider "hpe" {
  morpheus {
    url          = ""
    username     = ""
    password     = ""
  }
}
`

// routerFixture renders a self-contained NSX-T network router that NAT/route
// resources are proven to create on the QA appliance (group 3, NSX-T
// integration 5, tier-1 gateway type). The router is labelled
// hpe_morpheus_network_router.example.
func routerFixture(t *testing.T, name string) string {
	t.Helper()

	cfg, err := networkrouterresource.RenderNetworkRouterGenericConfig(t, map[string]string{
		"Name":                 name + "-router",
		"TypeId":               "9",
		"GroupId":              "3",
		"NetworkIntegrationId": "5",
	})
	if err != nil {
		t.Fatal(err)
	}

	return cfg
}

// routeFixture renders a route on the router fixture, labelled
// hpe_morpheus_network_router_route.example.
func routeFixture(t *testing.T, name string) string {
	t.Helper()

	cfg, err := routeresource.RenderNetworkRouterRouteConfig(t, map[string]string{
		"RouterId": "hpe_morpheus_network_router.example.id",
		"Name":     name,
	})
	if err != nil {
		t.Fatal(err)
	}

	return cfg
}

func TestAccMorpheusFindNetworkRouterRouteByName(t *testing.T) {
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

	// Look up the route created in the same config by name. depends_on defers
	// the data source read until the route exists (router_id only creates a
	// dependency on the router, not the route).
	dataSourceConfig := `
data "hpe_morpheus_network_router_route" "example" {
  name       = "` + name + `"
  router_id  = hpe_morpheus_network_router.example.id
  depends_on = [hpe_morpheus_network_router_route.example]
}
`

	checkFn := resource.ComposeAggregateTestCheckFunc(routeChecks()...)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + routerFixture(t, name) + routeFixture(t, name) + dataSourceConfig,
				Check:  checkFn,
			},
		},
	})
}

func TestAccMorpheusFindNetworkRouterRouteById(t *testing.T) {
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

	// Both id and router_id reference the created resources, so the data
	// source read is deferred until they exist.
	dataSourceConfig, err := networkrouterroute.RenderRouteByIdConfig(t, map[string]string{
		"Id":       "hpe_morpheus_network_router_route.example.id",
		"RouterId": "hpe_morpheus_network_router.example.id",
	})
	if err != nil {
		t.Fatal(err)
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(routeChecks()...)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + routerFixture(t, name) + routeFixture(t, name) + dataSourceConfig,
				Check:  checkFn,
			},
		},
	})
}

func TestAccMorpheusFindNetworkRouterRouteNotFound(t *testing.T) {
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

	// Search a real (created) router for a route name that does not exist.
	// router_id references the router so the read is deferred until it exists.
	dataSourceConfig := `
data "hpe_morpheus_network_router_route" "example" {
  name      = "nonexistent-route-name-that-should-not-exist"
  router_id = hpe_morpheus_network_router.example.id
}
`

	expected := regexp.MustCompile(`no network router route found`)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config:      providerConfig + routerFixture(t, name) + dataSourceConfig,
				ExpectError: expected,
			},
		},
	})
}

func TestAccMorpheusFindNetworkRouterRouteNoSearchAttrs(t *testing.T) {
	if capabilities.Missing(t, capabilities.NetworkRouter) {
		t.Log("Skipping test due to missing capabilities")

		return
	}
	defer testhelpers.RecordResult(t)

	t.Parallel()

	config := providerConfigOffline + `
      data "hpe_morpheus_network_router_route" "test" {
        router_id = 1
      }`

	expected := networkrouterroute.ErrorNoValidSearchTerms

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile(expected),
			},
		},
	})
}

func routeChecks() []resource.TestCheckFunc {
	ds := "data.hpe_morpheus_network_router_route.example"

	return []resource.TestCheckFunc{
		resource.TestCheckResourceAttrSet(ds, "id"),
		resource.TestCheckResourceAttrSet(ds, "router_id"),
		resource.TestCheckResourceAttrSet(ds, "name"),
		// 'code' is not populated for a freshly-created route (only synced/seeded
		// routes have it), so it is not asserted here.
		resource.TestCheckResourceAttrSet(ds, "route_type"),
		resource.TestCheckResourceAttrSet(ds, "source_type"),
		resource.TestCheckResourceAttrSet(ds, "external_id"),
		resource.TestCheckResourceAttrSet(ds, "provider_id"),
		resource.TestCheckResourceAttrSet(ds, "default_route"),
		resource.TestCheckResourceAttrSet(ds, "enabled"),
	}
}
