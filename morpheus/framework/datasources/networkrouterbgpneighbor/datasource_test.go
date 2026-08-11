// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkrouterbgpneighbor_test

import (
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/networkrouterbgpneighbor"
	bgpresource "github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/networkrouterbgpneighbor"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/nsxt"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
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

// neighborFixture renders a BGP neighbor on the pre-provisioned tier-0 gateway
// (resolved by name via nsxt.Tier0Config), labelled
// hpe_morpheus_network_router_bgp_neighbor.example.
//
// The returned config does not include nsxt.Tier0Config; callers emit it.
func neighborFixture(t *testing.T, name, ipAddress string) string {
	t.Helper()

	cfg, err := bgpresource.RenderBgpNeighborConfig(t, map[string]string{
		"RouterId":        nsxt.Tier0RouterIDRef,
		"IpAddress":       ipAddress,
		"Description":     name,
		"SourceAddresses": nsxt.BgpSourceAddressValue(),
	})
	if err != nil {
		t.Fatal(err)
	}

	return cfg
}

func TestAccMorpheusFindNetworkRouterBgpNeighborByIpAddress(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.NetworkRouter)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()
	name := acctest.RandomWithPrefix(t.Name())
	ipAddress := "192.168.50." + acctest.RandStringFromCharSet(2, "123456789")

	// Look up the neighbor created in the same config by IP. depends_on defers
	// the data source read until the neighbor exists.
	dataSourceConfig := `
data "hpe_morpheus_network_router_bgp_neighbor" "example" {
  ip_address = "` + ipAddress + `"
  router_id  = ` + nsxt.Tier0RouterIDRef + `
  depends_on = [hpe_morpheus_network_router_bgp_neighbor.example]
}
`

	checkFn := resource.ComposeAggregateTestCheckFunc(bgpNeighborChecks()...)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + nsxt.Tier0Config() +
					neighborFixture(t, name, ipAddress) + dataSourceConfig,
				Check: checkFn,
			},
		},
	})
}

func TestAccMorpheusFindNetworkRouterBgpNeighborById(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.NetworkRouter)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()
	name := acctest.RandomWithPrefix(t.Name())
	ipAddress := "192.168.51." + acctest.RandStringFromCharSet(2, "123456789")

	// id and router_id reference the created resources, deferring the read.
	dataSourceConfig, err := networkrouterbgpneighbor.RenderBgpNeighborByIdConfig(t, map[string]string{
		"Id":       "hpe_morpheus_network_router_bgp_neighbor.example.id",
		"RouterId": nsxt.Tier0RouterIDRef,
	})
	if err != nil {
		t.Fatal(err)
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(bgpNeighborChecks()...)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + nsxt.Tier0Config() +
					neighborFixture(t, name, ipAddress) + dataSourceConfig,
				Check: checkFn,
			},
		},
	})
}

func TestAccMorpheusFindNetworkRouterBgpNeighborNotFound(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.NetworkRouter)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()

	// Search a real (existing) router for a neighbor IP that does not exist.
	dataSourceConfig := `
data "hpe_morpheus_network_router_bgp_neighbor" "example" {
  ip_address = "0.0.0.0"
  router_id  = ` + nsxt.Tier0RouterIDRef + `
}
`

	expected := regexp.MustCompile(`no network router BGP neighbor found`)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config:      providerConfig + nsxt.Tier0Config() + dataSourceConfig,
				ExpectError: expected,
			},
		},
	})
}

func TestAccMorpheusFindNetworkRouterBgpNeighborNoSearchAttrs(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.NetworkRouter)

	t.Parallel()

	config := providerConfigOffline + `
      data "hpe_morpheus_network_router_bgp_neighbor" "test" {
        router_id = 1
      }`

	expected := networkrouterbgpneighbor.ErrorNoValidSearchTerms

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

func bgpNeighborChecks() []resource.TestCheckFunc {
	ds := "data.hpe_morpheus_network_router_bgp_neighbor.example"

	return []resource.TestCheckFunc{
		resource.TestCheckResourceAttrSet(ds, "id"),
		resource.TestCheckResourceAttrSet(ds, "router_id"),
		resource.TestCheckResourceAttrSet(ds, "ip_address"),
		resource.TestCheckResourceAttrSet(ds, "remote_as"),
	}
}
