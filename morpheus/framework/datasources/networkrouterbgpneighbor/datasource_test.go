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

// existingTier0RouterID is a pre-provisioned, fully-realized NSX-T tier-0 gateway
// (BGP enabled, with an associated edge cluster and local AS). BGP neighbors
// attach to the tier-0's locale-services, which are only populated in Morpheus
// after a sync of a realized gateway; the provider cannot create such a fixture
// per test, so these tests require an existing gateway supplied via
// TF_ACC_BGP_ROUTER_ID and skip when it is unset.
var existingTier0RouterID = os.Getenv("TF_ACC_BGP_ROUTER_ID")

// bgpNeighborSourceAddress is a valid IP on the tier-0's interface, required for
// EBGP multihop neighbors (see resource test for details). Override with
// TF_ACC_BGP_SOURCE_ADDRESS.
var bgpNeighborSourceAddress = bgpSourceAddressOrDefault()

func bgpSourceAddressOrDefault() string {
	if v := os.Getenv("TF_ACC_BGP_SOURCE_ADDRESS"); v != "" {
		return v
	}

	return "10.100.10.1"
}

// skipUnlessBGPRouter skips the test unless a pre-provisioned BGP-enabled NSX-T
// tier-0 gateway id is supplied via TF_ACC_BGP_ROUTER_ID.
func skipUnlessBGPRouter(t *testing.T) {
	t.Helper()

	if existingTier0RouterID == "" {
		t.Skip("TF_ACC_BGP_ROUTER_ID not set; skipping test requiring a pre-provisioned BGP-enabled NSX-T tier-0 gateway")
	}
}

// neighborFixture renders a BGP neighbor on the existing tier-0 router, labelled
// hpe_morpheus_network_router_bgp_neighbor.example.
func neighborFixture(t *testing.T, name, ipAddress string) string {
	t.Helper()

	cfg, err := bgpresource.RenderBgpNeighborConfig(t, map[string]string{
		"RouterId":        existingTier0RouterID,
		"IpAddress":       ipAddress,
		"Description":     name,
		"SourceAddresses": bgpNeighborSourceAddress,
	})
	if err != nil {
		t.Fatal(err)
	}

	return cfg
}

func TestAccMorpheusFindNetworkRouterBgpNeighborByIpAddress(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.NetworkRouter)

	skipUnlessBGPRouter(t)

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
  router_id  = ` + existingTier0RouterID + `
  depends_on = [hpe_morpheus_network_router_bgp_neighbor.example]
}
`

	checkFn := resource.ComposeAggregateTestCheckFunc(bgpNeighborChecks()...)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + neighborFixture(t, name, ipAddress) + dataSourceConfig,
				Check:  checkFn,
			},
		},
	})
}

func TestAccMorpheusFindNetworkRouterBgpNeighborById(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.NetworkRouter)

	skipUnlessBGPRouter(t)

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
		"RouterId": existingTier0RouterID,
	})
	if err != nil {
		t.Fatal(err)
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(bgpNeighborChecks()...)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + neighborFixture(t, name, ipAddress) + dataSourceConfig,
				Check:  checkFn,
			},
		},
	})
}

func TestAccMorpheusFindNetworkRouterBgpNeighborNotFound(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.NetworkRouter)

	skipUnlessBGPRouter(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()

	// Search a real (existing) router for a neighbor IP that does not exist.
	dataSourceConfig := `
data "hpe_morpheus_network_router_bgp_neighbor" "example" {
  ip_address = "0.0.0.0"
  router_id  = ` + existingTier0RouterID + `
}
`

	expected := regexp.MustCompile(`no network router BGP neighbor found`)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config:      providerConfig + dataSourceConfig,
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
