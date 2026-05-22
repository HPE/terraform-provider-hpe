// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkrouterroute_test

import (
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/networkrouterroute"
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

func TestAccMorpheusFindNetworkRouterRouteByName(t *testing.T) {
	if capabilities.Missing(t, capabilities.All, capabilities.NetworkRouter) {
		t.Log("Skipping test due to missing capabilities")

		return
	}
	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()

	dataSourceConfig, err := networkrouterroute.RenderRouteByNameConfig(t, map[string]string{
		"Name":     "mh-route",
		"RouterId": "19779",
	})
	if err != nil {
		t.Fatal(err)
	}

	checks := routeChecks()

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + dataSourceConfig,
				Check:  checkFn,
			},
		},
	})
}

func TestAccMorpheusFindNetworkRouterRouteById(t *testing.T) {
	if capabilities.Missing(t, capabilities.All, capabilities.NetworkRouter) {
		t.Log("Skipping test due to missing capabilities")

		return
	}
	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()

	dataSourceConfig, err := networkrouterroute.RenderRouteByIdConfig(t, map[string]string{
		"Id":       "33935",
		"RouterId": "19779",
	})
	if err != nil {
		t.Fatal(err)
	}

	checks := routeChecks()

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + dataSourceConfig,
				Check:  checkFn,
			},
		},
	})
}

func TestAccMorpheusFindNetworkRouterRouteNotFound(t *testing.T) {
	if capabilities.Missing(t, capabilities.All, capabilities.NetworkRouter) {
		t.Log("Skipping test due to missing capabilities")

		return
	}
	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()

	dataSourceConfig, err := networkrouterroute.RenderRouteByNameConfig(t,
		map[string]string{
			"Name":     "nonexistent-route-name-that-should-not-exist",
			"RouterId": "19779",
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	expected := regexp.MustCompile(`no network router route found`)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config:      providerConfig + dataSourceConfig,
				ExpectError: expected,
			},
		},
	})
}

func TestAccMorpheusFindNetworkRouterRouteNoSearchAttrs(t *testing.T) {
	if capabilities.Missing(t, capabilities.All, capabilities.NetworkRouter) {
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
		resource.TestCheckResourceAttrSet(ds, "code"),
		resource.TestCheckResourceAttrSet(ds, "route_type"),
		resource.TestCheckResourceAttrSet(ds, "source_type"),
		resource.TestCheckResourceAttrSet(ds, "external_id"),
		resource.TestCheckResourceAttrSet(ds, "provider_id"),
		resource.TestCheckResourceAttrSet(ds, "default_route"),
		resource.TestCheckResourceAttrSet(ds, "enabled"),
	}
}
