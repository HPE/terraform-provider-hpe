// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkrouter_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/networkrouter"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

const testResourceName = "hpe_morpheus_network_router.example"

func TestMain(m *testing.M) {
	code := testhelpers.TestMain(m)
	testhelpers.WriteMergedResults()
	os.Exit(code)
}

func TestAccMorpheusNetworkRouterResourceGenericExampleOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.NetworkRouter)

	t.Parallel()

	if testing.Short() {
		t.Skip("Skipping slow acceptance test in short mode")
	}

	name := acctest.RandomWithPrefix(t.Name())

	config, err := networkrouter.RenderNetworkRouterGenericConfig(t, map[string]string{
		"Name":                 name,
		"TypeId":               "9", // tier 1 gateway
		"GroupId":              "3",
		"NetworkIntegrationId": "5",
	})
	if err != nil {
		t.Fatalf("failed to render config: %s", err)
	}

	resourceName := testResourceName
	config = testhelpers.ProviderBlock() + config

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "type_id", "9"), // tier 1
					resource.TestCheckResourceAttr(resourceName, "group_id", "3"),
					resource.TestCheckResourceAttr(resourceName, "network_integration_id", "5"),
					resource.TestCheckResourceAttr(resourceName, "config.ipManagementType", "dhcpLocal"),
				),
			},
		},
	})
}

func TestAccMorpheusNetworkRouterResourceNSXTGatewayTier0ExampleOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.NSXT)

	t.Parallel()

	if testing.Short() {
		t.Skip("Skipping slow acceptance test in short mode")
	}

	name := acctest.RandomWithPrefix(t.Name())

	config, err := networkrouter.RenderNetworkRouterNSXTGatewayTier0Config(t, map[string]string{
		"Name":                 name,
		"GroupId":              "3",
		"NetworkIntegrationId": "5",
	})
	if err != nil {
		t.Fatalf("failed to render config: %s", err)
	}

	resourceName := testResourceName
	config = testhelpers.ProviderBlock() + config

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "type_id", "8"), // tier 0
					resource.TestCheckResourceAttr(resourceName, "group_id", "3"),
					resource.TestCheckResourceAttr(resourceName, "network_integration_id", "5"),
					resource.TestCheckResourceAttr(resourceName, "config_nsxt_gateway_tier0.ha_mode", "ACTIVE_ACTIVE"),
					resource.TestCheckResourceAttr(resourceName, "config_nsxt_gateway_tier0.restart_mode", "HELPER_ONLY"),
				),
			},
		},
	})
}

func TestAccMorpheusNetworkRouterResourceNSXTGatewayTier1ExampleOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.NSXT)

	t.Parallel()

	if testing.Short() {
		t.Skip("Skipping slow acceptance test in short mode")
	}

	name := acctest.RandomWithPrefix(t.Name())

	config, err := networkrouter.RenderNetworkRouterNSXTGatewayTier1Config(t, map[string]string{
		"Name":                 name,
		"GroupId":              "3",
		"NetworkIntegrationId": "5",
	})
	if err != nil {
		t.Fatalf("failed to render config: %s", err)
	}

	resourceName := testResourceName
	config = testhelpers.ProviderBlock() + config

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "type_id", "9"), // tier 1
					resource.TestCheckResourceAttr(resourceName, "group_id", "3"),
					resource.TestCheckResourceAttr(resourceName, "network_integration_id", "5"),
					resource.TestCheckResourceAttr(resourceName, "config_nsxt_gateway_tier1.ip_management_type", "dhcpLocal"),
				),
			},
		},
	})
}

// sdnNetworkPrereq renders a juniperSdn virtual network the SDN router joins.
// The SDN router needs at least one SDN virtual network on the same cluster
// resource pool; it is created per-test so parallel runs do not collide. It
// is named from the full test name, not the router's short name, so the
// network sweeper (which matches the suite-wide prefix only) can reclaim it.
func sdnNetworkPrereq(name string, f testhelpers.SdnFixture) string {
	return `
resource "hpe_morpheus_network" "sdn" {
  name              = "` + name + `-net"
  cloud_id          = ` + f.CloudID + `
  group_id          = ` + f.GroupID + `
  type_code         = "juniperSdn"
  network_server_id = ` + f.NetworkServerID + `
  zone_pool_id      = ` + f.ClusterPoolID + `
  cidr              = "10.91.0.0/24"
  gateway           = "10.91.0.1"
  visibility        = "private"
  tenant_ids        = [1]
}
`
}

// TestAccMorpheusNetworkRouterResourceJuniperSdnExampleOk creates an SDN
// logical router joining a juniperSdn network, attaches a floating IP pool to
// it and removes the pool again — the one in-place edit that needs no second
// network — then imports it. The config_juniper_sdn block selects the
// juniperSdnRouter type.
func TestAccMorpheusNetworkRouterResourceJuniperSdnExampleOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.JuniperSDN)
	f := testhelpers.RequireSdnFixture(t)

	t.Parallel()

	// SDN router names must be under 63 characters (the plugin validates
	// the length), so a short prefix stands in for the full test name.
	name := acctest.RandomWithPrefix(networkrouter.SdnTestNamePrefix)

	router, err := networkrouter.RenderNetworkRouterJuniperSdnConfig(t, map[string]string{
		"Name":                 name,
		"GroupId":              f.GroupID,
		"NetworkIntegrationId": f.NetworkServerID,
		"ResourcePoolId":       f.ClusterPoolID,
		"VirtualNetworkUids":   "hpe_morpheus_network.sdn.unique_id",
	})
	if err != nil {
		t.Fatalf("failed to render config: %s", err)
	}

	resourceName := testResourceName
	config := testhelpers.ProviderBlock() + sdnNetworkPrereq(acctest.RandomWithPrefix(t.Name()), f) + router

	// The same router with a floating IP pool attached, to be removed again:
	// an omitted set must be sent as an empty list and read back as null.
	withPool := strings.Replace(config,
		"virtual_network_uids = [hpe_morpheus_network.sdn.unique_id]",
		"virtual_network_uids = [hpe_morpheus_network.sdn.unique_id]\n    floating_ip_pool_ids = ["+f.FloatingIpPoolID+"]", 1)
	if withPool == config {
		t.Fatal("test configuration did not take the floating IP pool")
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "config_juniper_sdn.resource_pool_id", f.ClusterPoolID),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckNoResourceAttr(resourceName, "config_juniper_sdn.floating_ip_pool_ids.#"),
				),
			},
			{
				Config: withPool,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "config_juniper_sdn.floating_ip_pool_ids.#", "1"),
				),
			},
			{
				// Removing the pools clears them on the appliance and leaves the
				// attribute null, with nothing further to plan.
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckNoResourceAttr(resourceName, "config_juniper_sdn.floating_ip_pool_ids.#"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				// write_only_revision is the provider's own marker (set on create for
				// every router); import assigns a fresh one.
				ImportStateVerifyIgnore: []string{"write_only_revision"},
			},
		},
	})
}

// TestAccMorpheusNetworkRouterResourceJuniperBgpExampleOk creates an SDN BGP
// router and changes its hold time in place, then rotates the write-only
// md5_key_wo (the digest moves write_only_revision, nothing else changes) and
// clears it with "". md5_key_wo is never in state; write_only_revision is the
// provider's own marker, so import cannot know it.
func TestAccMorpheusNetworkRouterResourceJuniperBgpExampleOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.JuniperSDN)
	f := testhelpers.RequireSdnFixture(t)

	t.Parallel()

	// SDN router names must be under 63 characters (the plugin validates
	// the length), so a short prefix stands in for the full test name.
	name := acctest.RandomWithPrefix(networkrouter.SdnTestNamePrefix)
	resourceName := testResourceName

	render := func(holdTime, md5Key string) string {
		cfg, err := networkrouter.RenderNetworkRouterJuniperBgpConfig(t, map[string]string{
			"Name":                 name,
			"GroupId":              f.GroupID,
			"NetworkIntegrationId": f.NetworkServerID,
			"Vendor":               "juniper",
			"Asn":                  "65000",
			"PeerAddress1":         "10.0.0.1",
			"AddressFamilies":      `"inet"`,
			"HoldTime":             holdTime,
		})
		if err != nil {
			t.Fatalf("failed to render config: %s", err)
		}

		// The example leaves md5_key_wo out; the test adds it after hold_time
		// (a configured "" clears the key on the appliance).
		cfg = strings.Replace(cfg, "    hold_time        = "+holdTime+"\n",
			"    hold_time        = "+holdTime+"\n    md5_key_wo       = \""+md5Key+"\"\n", 1)
		if !strings.Contains(cfg, "md5_key_wo") {
			t.Fatal("md5_key_wo not inserted into the rendered configuration")
		}

		return testhelpers.ProviderBlock() + cfg
	}

	var revision string

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: render("90", "s3cret-one"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "config_juniper_bgp.vendor", "juniper"),
					resource.TestCheckResourceAttr(resourceName, "config_juniper_bgp.asn", "65000"),
					resource.TestCheckResourceAttr(resourceName, "config_juniper_bgp.hold_time", "90"),
					// Write-only: never in state.
					resource.TestCheckNoResourceAttr(resourceName, "config_juniper_bgp.md5_key_wo"),
					resource.TestCheckResourceAttrWith(resourceName, "write_only_revision", func(v string) error {
						revision = v

						return nil
					}),
				),
			},
			{
				Config: render("120", "s3cret-one"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "config_juniper_bgp.hold_time", "120"),
					// An unchanged key does not move the revision.
					resource.TestCheckResourceAttrWith(resourceName, "write_only_revision", func(v string) error {
						if v != revision {
							return fmt.Errorf("write_only_revision moved to %s on an unrelated change (was %s)", v, revision)
						}

						return nil
					}),
				),
			},
			{
				// Rotating the key is a plan: the digest differs, the revision moves.
				Config: render("120", "s3cret-two"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrWith(resourceName, "write_only_revision", func(v string) error {
						if v == revision {
							return fmt.Errorf("write_only_revision did not move when md5_key_wo changed (%s)", v)
						}
						revision = v

						return nil
					}),
				),
			},
			{
				// "" clears the key on the appliance; the revision moves again.
				Config: render("120", ""),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrWith(resourceName, "write_only_revision", func(v string) error {
						if v == revision {
							return fmt.Errorf("write_only_revision did not move when md5_key_wo was cleared (%s)", v)
						}

						return nil
					}),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"write_only_revision"},
			},
		},
	})
}
