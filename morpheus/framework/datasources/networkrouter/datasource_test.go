// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkrouter_test

import (
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/networkrouter"
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
		"Name":                 name,
		"TypeId":               "9",
		"GroupId":              "3",
		"NetworkIntegrationId": "5",
	})
	if err != nil {
		t.Fatal(err)
	}

	return cfg
}

func TestAccMorpheusFindNetworkRouterByName(t *testing.T) {
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

	// Look up the router created in the same config by name. depends_on defers
	// the data source read until the router exists (the name is a literal and
	// otherwise creates no dependency).
	dataSourceConfig := `
data "hpe_morpheus_network_router" "example" {
  name       = "` + name + `"
  depends_on = [hpe_morpheus_network_router.example]
}
`

	checkFn := resource.ComposeAggregateTestCheckFunc(networkRouterChecks()...)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + routerFixture(t, name) + dataSourceConfig,
				Check:  checkFn,
			},
		},
	})
}

func TestAccMorpheusFindNetworkRouterById(t *testing.T) {
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

	// id references the created router, so the data source read is deferred
	// until it exists.
	dataSourceConfig, err := networkrouter.RenderNetworkRouterByIdConfig(t, map[string]string{
		"Id": "hpe_morpheus_network_router.example.id",
	})
	if err != nil {
		t.Fatal(err)
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(networkRouterChecks()...)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + routerFixture(t, name) + dataSourceConfig,
				Check:  checkFn,
			},
		},
	})
}

func TestAccMorpheusFindNetworkRouterNotFound(t *testing.T) {
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

	dataSourceConfig, err := networkrouter.RenderNetworkRouterByNameConfig(t,
		map[string]string{
			"Name": "______",
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	expected := regexp.MustCompile(`no network router found`)

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

func TestAccMorpheusFindNetworkRouterNoSearchAttrs(t *testing.T) {
	if capabilities.Missing(t, capabilities.NetworkRouter) {
		t.Log("Skipping test due to missing capabilities")

		return
	}
	defer testhelpers.RecordResult(t)

	t.Parallel()

	config := providerConfigOffline + `
      data "hpe_morpheus_network_router" "test" {
      }`

	expected := networkrouter.ErrorNoValidSearchTerms

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

func TestAccMorpheusFindNetworkRouterBothSearchAttrs(t *testing.T) {
	if capabilities.Missing(t, capabilities.NetworkRouter) {
		t.Log("Skipping test due to missing capabilities")

		return
	}
	defer testhelpers.RecordResult(t)

	t.Parallel()

	config := providerConfigOffline + `
      data "hpe_morpheus_network_router" "test" {
        id   = 1
        name = "______"
      }`

	expected := networkrouter.ErrorRunningPreApply

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

// networkRouterChecks asserts attributes that a freshly-created generic NSX-T
// router reliably returns. NOTE (QA verify): cloud.* and config.bridgeName were
// asserted against the previously hardcoded seeded router; a created tier-1
// gateway may not populate those, so they are intentionally not checked here.
func networkRouterChecks() []resource.TestCheckFunc {
	ds := "data.hpe_morpheus_network_router.example"

	return []resource.TestCheckFunc{
		resource.TestCheckResourceAttrSet(ds, "id"),
		resource.TestCheckResourceAttrSet(ds, "name"),
		resource.TestCheckResourceAttrSet(ds, "enabled"),
		resource.TestCheckResourceAttrSet(ds, "enable_bgp"),
		resource.TestCheckResourceAttrSet(ds, "group.id"),
		resource.TestCheckResourceAttrSet(ds, "group.name"),
		resource.TestCheckResourceAttrSet(ds, "permissions.visibility"),
	}
}
