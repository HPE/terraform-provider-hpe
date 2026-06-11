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

// routerFixture renders a self-contained NSX-T tier-0 gateway router (group 3,
// NSX-T integration 5) with BGP enabled, labelled
// hpe_morpheus_network_router.example.
//
// QA verify: edge_cluster "qa-edge-cluster-01", local_as_num 65000, NSX-T
// integration 5 and group 3 are the QA appliance values.
func routerFixture(t *testing.T, name string) string {
	t.Helper()

	return `
resource "hpe_morpheus_network_router" "example" {
  name                   = "` + name + `-router"
  group_id               = 3
  network_integration_id = 5
  enable_bgp             = true

  config_nsxt_gateway_tier0 = {
    ha_mode      = "ACTIVE_ACTIVE"
    restart_mode = "HELPER_ONLY"
    edge_cluster = "qa-edge-cluster-01"
    fail_over    = "NON_PREEMPTIVE"
    local_as_num = "65000"
  }
}
`
}

// neighborFixture renders a BGP neighbor on the router fixture, labelled
// hpe_morpheus_network_router_bgp_neighbor.example.
func neighborFixture(t *testing.T, name, ipAddress string) string {
	t.Helper()

	cfg, err := bgpresource.RenderBgpNeighborConfig(t, map[string]string{
		"RouterId":    "hpe_morpheus_network_router.example.id",
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
  router_id  = hpe_morpheus_network_router.example.id
  depends_on = [hpe_morpheus_network_router_bgp_neighbor.example]
}
`

	checkFn := resource.ComposeAggregateTestCheckFunc(bgpNeighborChecks()...)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + routerFixture(t, name) + neighborFixture(t, name, ipAddress) + dataSourceConfig,
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
		"RouterId": "hpe_morpheus_network_router.example.id",
	})
	if err != nil {
		t.Fatal(err)
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(bgpNeighborChecks()...)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + routerFixture(t, name) + neighborFixture(t, name, ipAddress) + dataSourceConfig,
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
	name := acctest.RandomWithPrefix(t.Name())

	// Search a real (created) router for a neighbor IP that does not exist.
	dataSourceConfig := `
data "hpe_morpheus_network_router_bgp_neighbor" "example" {
  ip_address = "0.0.0.0"
  router_id  = hpe_morpheus_network_router.example.id
}
`

	expected := regexp.MustCompile(`no network router BGP neighbor found`)

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
