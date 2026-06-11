// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkrouterbgpneighbor_test

import (
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/networkrouterbgpneighbor"
	bgpresource "github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/networkrouterbgpneighbor"
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

// existingTier0RouterID is a pre-provisioned, fully-realized NSX-T tier-0 gateway
// (BGP enabled, with an associated edge cluster and local AS) on integration 5.
// BGP neighbors attach to the tier-0's locale-services, which are only populated
// in Morpheus after a sync of a realized gateway; creating a tier-0 per test
// races that sync, so we reference this existing gateway.
const existingTier0RouterID = "28"

// neighborFixture renders a BGP neighbor on the existing tier-0 router, labelled
// hpe_morpheus_network_router_bgp_neighbor.example.
func neighborFixture(t *testing.T, name, ipAddress string) string {
	t.Helper()

	cfg, err := bgpresource.RenderBgpNeighborConfig(t, map[string]string{
		"RouterId":    existingTier0RouterID,
		"IpAddress":   ipAddress,
		"Description": name,
	})
	if err != nil {
		t.Fatal(err)
	}

	return cfg
}

func TestAccMorpheusFindNetworkRouterBgpNeighborByIpAddress(t *testing.T) {
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
	ipAddress := "192.168.50." + acctest.RandStringFromCharSet(2, "123456789")

	// Look up the neighbor created in the same config by IP. depends_on defers
	// the data source read until the neighbor exists.
	dataSourceConfig := `
data "hpe_morpheus_network_router_bgp_neighbor" "example" {
  ip_address = "` + ipAddress + `"
  router_id  = 28
  depends_on = [hpe_morpheus_network_router_bgp_neighbor.example]
}
`

	checkFn := resource.ComposeAggregateTestCheckFunc(bgpNeighborChecks()...)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + neighborFixture(t, name, ipAddress) + dataSourceConfig,
				Check:  checkFn,
			},
		},
	})
}

func TestAccMorpheusFindNetworkRouterBgpNeighborById(t *testing.T) {
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
	ipAddress := "192.168.51." + acctest.RandStringFromCharSet(2, "123456789")

	// id and router_id reference the created resources, deferring the read.
	dataSourceConfig, err := networkrouterbgpneighbor.RenderBgpNeighborByIdConfig(t, map[string]string{
		"Id":       "hpe_morpheus_network_router_bgp_neighbor.example.id",
		"RouterId": "28",
	})
	if err != nil {
		t.Fatal(err)
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(bgpNeighborChecks()...)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + neighborFixture(t, name, ipAddress) + dataSourceConfig,
				Check:  checkFn,
			},
		},
	})
}

func TestAccMorpheusFindNetworkRouterBgpNeighborNotFound(t *testing.T) {
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

	// Search a real (existing) router for a neighbor IP that does not exist.
	dataSourceConfig := `
data "hpe_morpheus_network_router_bgp_neighbor" "example" {
  ip_address = "0.0.0.0"
  router_id  = 28
}
`

	expected := regexp.MustCompile(`no network router BGP neighbor found`)

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

func TestAccMorpheusFindNetworkRouterBgpNeighborNoSearchAttrs(t *testing.T) {
	if capabilities.Missing(t, capabilities.NetworkRouter) {
		t.Log("Skipping test due to missing capabilities")

		return
	}
	defer testhelpers.RecordResult(t)

	t.Parallel()

	config := providerConfigOffline + `
      data "hpe_morpheus_network_router_bgp_neighbor" "test" {
        router_id = 1
      }`

	expected := networkrouterbgpneighbor.ErrorNoValidSearchTerms

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

func bgpNeighborChecks() []resource.TestCheckFunc {
	ds := "data.hpe_morpheus_network_router_bgp_neighbor.example"

	return []resource.TestCheckFunc{
		resource.TestCheckResourceAttrSet(ds, "id"),
		resource.TestCheckResourceAttrSet(ds, "router_id"),
		resource.TestCheckResourceAttrSet(ds, "ip_address"),
		resource.TestCheckResourceAttrSet(ds, "remote_as"),
	}
}
