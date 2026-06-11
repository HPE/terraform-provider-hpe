// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkrouterbgpneighbor_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/HPE/terraform-provider-hpe/morpheus"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/networkrouterbgpneighbor"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
)

func TestMain(m *testing.M) {
	code := testhelpers.TestMain(m)

	testhelpers.WriteMergedResults()

	os.Exit(code)
}

// existingTier0RouterID is a pre-provisioned, fully-realized NSX-T tier-0 gateway
// (BGP enabled, with an associated edge cluster and local AS) on integration 5.
// BGP neighbors attach to the tier-0's locale-services, which are only populated
// in Morpheus after a sync of a realized gateway; creating a tier-0 per test
// races that sync, so we reference this existing gateway.
const existingTier0RouterID = "28"

func TestAccMorpheusNetworkRouterBgpNeighborResourceExampleOk(t *testing.T) {
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
	ipAddress := "192.168.10." + acctest.RandStringFromCharSet(2, "123456789")

	config, err := networkrouterbgpneighbor.RenderBgpNeighborConfig(t, map[string]string{
		"Description": name,
		"IpAddress":   ipAddress,
		"RouterId":    existingTier0RouterID,
	})
	if err != nil {
		t.Fatalf("failed to render config: %s", err)
	}

	resourceName := "hpe_morpheus_network_router_bgp_neighbor.example"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "ip_address", ipAddress),
					resource.TestCheckResourceAttr(resourceName, "description", name),
					resource.TestCheckResourceAttr(resourceName, "remote_as", "65001"),
					resource.TestCheckResourceAttr(resourceName, "weight", "60"),
					resource.TestCheckResourceAttr(resourceName, "keep_alive", "60"),
					resource.TestCheckResourceAttr(resourceName, "hold_down", "180"),
				),
			},
		},
	})
}

func TestAccMorpheusNetworkRouterBgpNeighborResourceCreate(t *testing.T) {
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
	ipAddress := "192.168.1." + acctest.RandStringFromCharSet(2, "123456789")

	configText := providerConfig + `
resource "hpe_morpheus_network_router_bgp_neighbor" "test" {
  router_id   = 28
  ip_address  = "` + ipAddress + `"
  description = "` + name + `"
  remote_as   = "65001"
  weight      = 60
  keep_alive  = 60
  hold_down   = 180
}
`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config: configText,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"hpe_morpheus_network_router_bgp_neighbor.test", "id",
					),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network_router_bgp_neighbor.test",
						"ip_address", ipAddress,
					),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network_router_bgp_neighbor.test",
						"remote_as", "65001",
					),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network_router_bgp_neighbor.test",
						"weight", "60",
					),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network_router_bgp_neighbor.test",
						"keep_alive", "60",
					),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network_router_bgp_neighbor.test",
						"hold_down", "180",
					),
				),
			},
		},
	})
}

func TestAccMorpheusNetworkRouterBgpNeighborResourceCreateAllAttrs(t *testing.T) {
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
	ipAddress := "10.0.0." + acctest.RandStringFromCharSet(2, "123456789")

	configText := providerConfig + `
resource "hpe_morpheus_network_router_bgp_neighbor" "all_attrs" {
  router_id            = 28
  ip_address           = "` + ipAddress + `"
  description          = "` + name + `"
  remote_as            = "65002"
  weight               = 100
  keep_alive           = 30
  hold_down            = 90
  bfd_enabled          = true
  bfd_interval         = 500
  bfd_multiple         = 3
  allow_as_in          = true
  hop_limit            = 2
  restart_mode         = "HELPER_ONLY"
  route_filtering_type = "IPV4"
  route_filtering_in   = "filter-in"
  route_filtering_out  = "filter-out"
}
`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config: configText,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"hpe_morpheus_network_router_bgp_neighbor.all_attrs", "id",
					),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network_router_bgp_neighbor.all_attrs",
						"ip_address", ipAddress,
					),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network_router_bgp_neighbor.all_attrs",
						"description", name,
					),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network_router_bgp_neighbor.all_attrs",
						"remote_as", "65002",
					),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network_router_bgp_neighbor.all_attrs",
						"weight", "100",
					),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network_router_bgp_neighbor.all_attrs",
						"keep_alive", "30",
					),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network_router_bgp_neighbor.all_attrs",
						"hold_down", "90",
					),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network_router_bgp_neighbor.all_attrs",
						"bfd_enabled", "true",
					),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network_router_bgp_neighbor.all_attrs",
						"bfd_interval", "500",
					),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network_router_bgp_neighbor.all_attrs",
						"bfd_multiple", "3",
					),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network_router_bgp_neighbor.all_attrs",
						"allow_as_in", "true",
					),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network_router_bgp_neighbor.all_attrs",
						"hop_limit", "2",
					),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network_router_bgp_neighbor.all_attrs",
						"restart_mode", "HELPER_ONLY",
					),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network_router_bgp_neighbor.all_attrs",
						"route_filtering_type", "IPV4",
					),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network_router_bgp_neighbor.all_attrs",
						"route_filtering_in", "filter-in",
					),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network_router_bgp_neighbor.all_attrs",
						"route_filtering_out", "filter-out",
					),
				),
			},
		},
	})
}

func TestAccMorpheusNetworkRouterBgpNeighborResourceUpdate(t *testing.T) {
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
	ipAddress := "172.16.0." + acctest.RandStringFromCharSet(2, "123456789")

	createConfig := providerConfig + `
resource "hpe_morpheus_network_router_bgp_neighbor" "update_test" {
  router_id   = 28
  ip_address  = "` + ipAddress + `"
  description = "` + name + `"
  remote_as   = "65001"
  weight      = 60
  keep_alive  = 60
  hold_down   = 180
}
`

	updateConfig := providerConfig + `
resource "hpe_morpheus_network_router_bgp_neighbor" "update_test" {
  router_id    = 28
  ip_address   = "` + ipAddress + `"
  description  = "` + name + ` updated"
  remote_as    = "65002"
  weight       = 100
  keep_alive   = 30
  hold_down    = 90
  bfd_enabled  = true
  bfd_interval = 500
  bfd_multiple = 3
  allow_as_in  = true
  hop_limit    = 3
  restart_mode = "GRACEFUL_RESTART"
}
`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config: createConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"hpe_morpheus_network_router_bgp_neighbor.update_test", "id",
					),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network_router_bgp_neighbor.update_test",
						"remote_as", "65001",
					),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network_router_bgp_neighbor.update_test",
						"weight", "60",
					),
				),
			},
			{
				Config: updateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network_router_bgp_neighbor.update_test",
						"description", name+" updated",
					),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network_router_bgp_neighbor.update_test",
						"remote_as", "65002",
					),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network_router_bgp_neighbor.update_test",
						"weight", "100",
					),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network_router_bgp_neighbor.update_test",
						"keep_alive", "30",
					),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network_router_bgp_neighbor.update_test",
						"hold_down", "90",
					),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network_router_bgp_neighbor.update_test",
						"bfd_enabled", "true",
					),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network_router_bgp_neighbor.update_test",
						"bfd_interval", "500",
					),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network_router_bgp_neighbor.update_test",
						"bfd_multiple", "3",
					),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network_router_bgp_neighbor.update_test",
						"allow_as_in", "true",
					),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network_router_bgp_neighbor.update_test",
						"hop_limit", "3",
					),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network_router_bgp_neighbor.update_test",
						"restart_mode", "GRACEFUL_RESTART",
					),
				),
			},
		},
	})
}

func TestAccMorpheusNetworkRouterBgpNeighborResourceImport(t *testing.T) {
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
	ipAddress := "10.1.0." + acctest.RandStringFromCharSet(2, "123456789")

	configText := providerConfig + `
resource "hpe_morpheus_network_router_bgp_neighbor" "import_test" {
  router_id   = 28
  ip_address  = "` + ipAddress + `"
  description = "` + name + `"
  remote_as   = "65001"
}
`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config: configText,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"hpe_morpheus_network_router_bgp_neighbor.import_test", "id",
					),
				),
			},
			{
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password_wo", "password_wo_version"},
				ResourceName:            "hpe_morpheus_network_router_bgp_neighbor.import_test",
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources["hpe_morpheus_network_router_bgp_neighbor.import_test"]
					if !ok {
						return "", fmt.Errorf("resource not found")
					}

					return rs.Primary.Attributes["router_id"] + "." + rs.Primary.Attributes["id"], nil
				},
			},
		},
	})
}

func TestAccMorpheusNetworkRouterBgpNeighborResourceWithNsxtConfig(t *testing.T) {
	if capabilities.Missing(t, capabilities.NSXT) {
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
	ipAddress := "10.2.0." + acctest.RandStringFromCharSet(2, "123456789")

	configText := providerConfig + `
resource "hpe_morpheus_network_router_bgp_neighbor" "nsxt_test" {
  router_id   = 28
  ip_address  = "` + ipAddress + `"
  description = "` + name + `"
  remote_as   = "65010"

  config_nsxt = {
    source_addresses = ["10.0.0.1", "10.0.0.2"]
  }
}
`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config: configText,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"hpe_morpheus_network_router_bgp_neighbor.nsxt_test", "id",
					),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network_router_bgp_neighbor.nsxt_test",
						"ip_address", ipAddress,
					),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network_router_bgp_neighbor.nsxt_test",
						"remote_as", "65010",
					),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network_router_bgp_neighbor.nsxt_test",
						"config_nsxt.source_addresses.#", "2",
					),
				),
			},
		},
	})
}

// TestAccMorpheusNetworkRouterBgpNeighborResourceWithNsxvConfig requires a
// pre-existing NSX-V network router: the provider cannot create NSX-V routers
// (only NSX-T tier-0/tier-1 gateways), so this test cannot be made fully
// self-contained. Supply a real NSX-V router id via TF_VAR_nsxv_router_id.
//
// QA verify: set nsxv_router_id to a valid NSX-V router on the target system.
func TestAccMorpheusNetworkRouterBgpNeighborResourceWithNsxvConfig(t *testing.T) {
	if capabilities.Missing(t, capabilities.NSXV) {
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
	ipAddress := "10.3.0." + acctest.RandStringFromCharSet(2, "123456789")

	configText := providerConfig + `
variable "nsxv_router_id" {
  description = "ID of a pre-existing NSX-V network router (provider cannot create NSX-V routers)"
  type        = number
}

resource "hpe_morpheus_network_router_bgp_neighbor" "nsxv_test" {
  router_id   = var.nsxv_router_id
  ip_address  = "` + ipAddress + `"
  description = "` + name + `"
  remote_as   = "65020"

  config_nsxv = {
    router_id = "10.0.0.1"
    interface = "vNic_0"
  }
}
`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config: configText,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"hpe_morpheus_network_router_bgp_neighbor.nsxv_test", "id",
					),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network_router_bgp_neighbor.nsxv_test",
						"ip_address", ipAddress,
					),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network_router_bgp_neighbor.nsxv_test",
						"remote_as", "65020",
					),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network_router_bgp_neighbor.nsxv_test",
						"config_nsxv.router_id", "10.0.0.1",
					),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network_router_bgp_neighbor.nsxv_test",
						"config_nsxv.interface", "vNic_0",
					),
				),
			},
		},
	})
}
