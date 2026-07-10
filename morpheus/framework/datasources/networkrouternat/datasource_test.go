// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkrouternat_test

import (
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/networkrouternat"
	natresource "github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/networkrouternat"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

func TestMain(m *testing.M) {
	code := testhelpers.TestMain(m)
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

// nsxtTier1RouterConfig renders an NSX-T tier-1 gateway required for NAT rules.
func nsxtTier1RouterConfig(name string) string {
	return `
data "hpe_morpheus_network_router" "nat_tier0" {
  id = 28
}

resource "hpe_morpheus_network_router" "nat_tier1" {
  name                   = "` + name + `-tier1"
  group_id               = 3
  network_integration_id = 5

  config_nsxt_gateway_tier1 = {
    ip_management_type = "dhcpLocal"
    edge_cluster       = "3de5f8d0-4f8a-433b-95ed-91020c948084"
    fail_over          = "NON_PREEMPTIVE"
    tier0_gateway      = data.hpe_morpheus_network_router.nat_tier0.provider_id
  }
}
`
}

func TestAccMorpheusFindNetworkRouterNatByName(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.NetworkRouter)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()
	name := acctest.RandomWithPrefix(t.Name())

	routerConfig := nsxtTier1RouterConfig(name)

	resourceConfig, err := natresource.RenderNetworkRouterNatConfig(t, map[string]string{
		"RouterId": "hpe_morpheus_network_router.nat_tier1.id",
		"Name":     name,
	})
	if err != nil {
		t.Fatal(err)
	}

	dataSourceConfig := `
data "hpe_morpheus_network_router_nat" "example" {
  name       = "` + name + `"
  router_id  = hpe_morpheus_network_router.nat_tier1.id
  depends_on = [hpe_morpheus_network_router_nat.example]
}
`

	checkFn := resource.ComposeAggregateTestCheckFunc(natChecks()...)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + routerConfig + resourceConfig + dataSourceConfig,
				Check:  checkFn,
			},
		},
	})
}

func TestAccMorpheusFindNetworkRouterNatById(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.NetworkRouter)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()
	name := acctest.RandomWithPrefix(t.Name())

	routerConfig := nsxtTier1RouterConfig(name)

	resourceConfig, err := natresource.RenderNetworkRouterNatConfig(t, map[string]string{
		"RouterId": "hpe_morpheus_network_router.nat_tier1.id",
		"Name":     name,
	})
	if err != nil {
		t.Fatal(err)
	}

	dataSourceConfig, err := networkrouternat.RenderNatByIdConfig(t, map[string]string{
		"Id":       "hpe_morpheus_network_router_nat.example.id",
		"RouterId": "hpe_morpheus_network_router.nat_tier1.id",
	})
	if err != nil {
		t.Fatal(err)
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(natChecks()...)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + routerConfig + resourceConfig + dataSourceConfig,
				Check:  checkFn,
			},
		},
	})
}

func TestAccMorpheusFindNetworkRouterNatNotFound(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.NetworkRouter)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()
	name := acctest.RandomWithPrefix(t.Name())

	routerConfig := nsxtTier1RouterConfig(name)

	dataSourceConfig := `
data "hpe_morpheus_network_router_nat" "example" {
  name      = "nonexistent-nat-name-that-should-not-exist"
  router_id = hpe_morpheus_network_router.nat_tier1.id
}
`

	expected := regexp.MustCompile(`no network router nat found`)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config:      providerConfig + routerConfig + dataSourceConfig,
				ExpectError: expected,
			},
		},
	})
}

func TestAccMorpheusFindNetworkRouterNatNoSearchAttrs(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.NetworkRouter)

	t.Parallel()

	config := providerConfigOffline + `
      data "hpe_morpheus_network_router_nat" "test" {
        router_id = 1
      }`

	expected := networkrouternat.ErrorNoValidSearchTerms

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile(expected),
			},
		},
	})
}

func natChecks() []resource.TestCheckFunc {
	ds := "data.hpe_morpheus_network_router_nat.example"

	return []resource.TestCheckFunc{
		resource.TestCheckResourceAttrSet(ds, "id"),
		resource.TestCheckResourceAttrSet(ds, "router_id"),
		resource.TestCheckResourceAttrSet(ds, "name"),
		resource.TestCheckResourceAttrSet(ds, "action"),
		resource.TestCheckResourceAttrSet(ds, "enabled"),
	}
}
