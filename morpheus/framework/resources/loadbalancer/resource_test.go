// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package loadbalancer_test

import (
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/loadbalancer"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

func TestMain(m *testing.M) {
	code := testhelpers.TestMain(m)
	testhelpers.WriteMergedResults()
	os.Exit(code)
}

func TestAccMorpheusLoadBalancerResourceHAProxyExampleOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.NetworkLoadBalancer, capabilities.NetworkLoadBalancerHAProxy)

	t.Parallel()

	if testing.Short() {
		t.Skip("Skipping slow acceptance test in short mode")
	}

	name := acctest.RandomWithPrefix(t.Name())
	// This resource only allows name to be 32 characters maximum.
	name = name[0:16] + name[len(name)-16:]

	config, err := loadbalancer.RenderLoadBalancerHAProxyConfig(t, map[string]string{
		"Name": name,
	})
	if err != nil {
		t.Fatalf("failed to render config: %s", err)
	}

	resourceName := "hpe_morpheus_load_balancer.haproxy"
	config = testhelpers.ProviderBlock() + config

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "visibility", "public"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"config",
					"config_haproxy",
					"cloud_id",
					// group_id and network_server_id are WriteOnly create
					// inputs: the framework never persists them to state (Read
					// sets them null) and the load-balancer GET never returns
					// them, so they cannot diff on import and need no ignore.
					// The WriteOnly conversion is the actual fix for
					// MORPH-10303.
					"write_only_revision",
				},
			},
		},
	})
}

func TestAccMorpheusLoadBalancerResourceHAProxyGenericExampleOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.NetworkLoadBalancer, capabilities.NetworkLoadBalancerHAProxy)

	t.Parallel()

	if testing.Short() {
		t.Skip("Skipping slow acceptance test in short mode")
	}

	name := acctest.RandomWithPrefix(t.Name())
	// This resource only allows name to be 32 characters maximum.
	name = name[0:16] + name[len(name)-16:]

	config, err := loadbalancer.RenderLoadBalancerHAProxyGenericConfig(t, map[string]string{
		"Name": name,
	})
	if err != nil {
		t.Fatalf("failed to render config: %s", err)
	}

	resourceName := "hpe_morpheus_load_balancer.haproxy_generic"
	config = testhelpers.ProviderBlock() + config

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "type_code", "haproxyContainer"),
					resource.TestCheckResourceAttr(resourceName, "visibility", "public"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"config",
					"config_haproxy",
					"cloud_id",
					"type_code",
					"permissions.groups",
					// group_id/network_server_id are WriteOnly (never in state)
					// and not returned by GET, so they cannot diff on import.
					// See MORPH-10303.
					"write_only_revision",
				},
			},
		},
	})
}

// Test validation: permissions.all conflicts with permissions.groups
func TestAccMorpheusLoadBalancerResourceValidationPermissionsConflict(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.NetworkLoadBalancer)

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()
	name := acctest.RandomWithPrefix(t.Name())
	name = name[0:16] + name[len(name)-16:]

	resourceConfig := `
resource "hpe_morpheus_load_balancer" "test" {
  name = "` + name + `"

  config = {}

  permissions = {
    all    = true
    groups = [1, 2]
  }
}
`

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config:      providerConfig + resourceConfig,
				ExpectError: regexp.MustCompile(`(?i)attribute "permissions.groups" cannot be specified when`),
			},
		},
	})
}

// TestMorpheusLoadBalancerWriteOnlyConfigParsing verifies that write-only attributes
// (group_id, network_server_id) are accepted in config and produce a valid plan.
// The framework nullifies write-only values in the plan, so Create must source them
// from req.Config — this test confirms the schema accepts the attributes correctly.
func TestAccMorpheusLoadBalancerWriteOnlyConfigParsing(t *testing.T) {
	defer testhelpers.RecordResult(t)
	t.Parallel()

	if testing.Short() {
		t.Skip("Skipping acceptance test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()
	resourceConfig := `
resource "hpe_morpheus_load_balancer" "test" {
  name              = "test-wo-attrs"
  group_id          = 1
  network_server_id = 42
  config_nsxt = {
    admin_state   = true
    log_level     = "INFO"
    size          = "SMALL"
    tier1_gateway = "/infra/tier-1s/test"
  }
}
`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config:             providerConfig + resourceConfig,
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

// TestMorpheusLoadBalancerWriteOnlyVersionCompanions verifies that the version
// companion attributes (group_id_version, network_server_id_version) are accepted
// in config alongside the write-only attributes they track.
func TestAccMorpheusLoadBalancerWriteOnlyVersionCompanions(t *testing.T) {
	defer testhelpers.RecordResult(t)
	t.Parallel()

	if testing.Short() {
		t.Skip("Skipping acceptance test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()
	resourceConfig := `
resource "hpe_morpheus_load_balancer" "test" {
  name                      = "test-wo-companions"
  group_id                  = 1
  group_id_version          = 1
  network_server_id         = 42
  network_server_id_version = 1
  config_nsxt = {
    admin_state   = true
    log_level     = "INFO"
    size          = "SMALL"
    tier1_gateway = "/infra/tier-1s/test"
  }
}
`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config:             providerConfig + resourceConfig,
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

// TestAccMorpheusLoadBalancerResourceJuniperSdnExampleOk creates an SDN
// (juniper-sdn-lb) load balancer on a juniperSdn virtual network. The SDN LB
// configuration is immutable after create, so this is a create-and-import test.
func TestAccMorpheusLoadBalancerResourceJuniperSdnExampleOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.JuniperSDN)
	// Held: deleting any load balancer on Morpheus 9.1.x wipes every other load
	// balancer's pool and virtual-server records (MORPH-18174), so this test
	// must not run on a shared appliance until that fix lands.
	if os.Getenv(testhelpers.EnvSdnLoadBalancer) == "" {
		t.Skip("SDN load balancer test held for MORPH-18174; set " + testhelpers.EnvSdnLoadBalancer + "=1 to run")
	}
	f := testhelpers.RequireSdnFixture(t)

	t.Parallel()

	if testing.Short() {
		t.Skip("Skipping slow acceptance test in short mode")
	}

	name := acctest.RandomWithPrefix(t.Name())
	name = name[0:16] + name[len(name)-16:]

	// The SDN virtual network the LB's virtual_network_uid points at; created
	// per-test so parallel runs do not collide.
	network := `
resource "hpe_morpheus_network" "sdn" {
  name              = "` + name + `-net"
  cloud_id          = ` + f.CloudID + `
  group_id          = ` + f.GroupID + `
  type_code         = "juniperSdn"
  network_server_id = ` + f.NetworkServerID + `
  zone_pool_id      = ` + f.ClusterPoolID + `
  cidr              = "10.92.0.0/24"
  gateway           = "10.92.0.1"
  visibility        = "private"
  tenant_ids        = [1]
}
`

	lb, err := loadbalancer.RenderLoadBalancerJuniperSdnConfig(t, map[string]string{
		"Name":              name,
		"CloudId":           f.CloudID,
		"NetworkServerId":   f.NetworkServerID,
		"ResourcePoolId":    f.ClusterPoolID,
		"VirtualNetworkUid": "hpe_morpheus_network.sdn.unique_id",
		"VirtualIp":         "10.92.0.100",
		"SourcePort":        "80",
		"DestinationPort":   "8080",
		"Protocol":          "TCP",
	})
	if err != nil {
		t.Fatalf("failed to render config: %s", err)
	}

	resourceName := "hpe_morpheus_load_balancer.sdn"
	config := testhelpers.ProviderBlock() + network + lb

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "type_code", "juniper-sdn-lb"),
					resource.TestCheckResourceAttr(resourceName, "config_juniper_sdn.virtual_ip", "10.92.0.100"),
					resource.TestCheckResourceAttr(resourceName, "config_juniper_sdn.source_port", "80"),
					resource.TestCheckResourceAttr(resourceName, "config_juniper_sdn.destination_port", "8080"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"cloud_id",
					"network_server_id",
					"write_only_revision",
				},
			},
		},
	})
}

// TestAccMorpheusLoadBalancerResourceJuniperSdnTypeCodeOmittedOk creates the
// same SDN load balancer with type_code left out: the config_juniper_sdn block
// implies the type, so type_code is computed from the appliance, the next
// plan is empty, and the import matches (MORPH-20774).
func TestAccMorpheusLoadBalancerResourceJuniperSdnTypeCodeOmittedOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.JuniperSDN)
	// Held: deleting any load balancer on Morpheus 9.1.x wipes every other load
	// balancer's pool and virtual-server records (MORPH-18174), so this test
	// must not run on a shared appliance until that fix lands.
	if os.Getenv(testhelpers.EnvSdnLoadBalancer) == "" {
		t.Skip("SDN load balancer test held for MORPH-18174; set " + testhelpers.EnvSdnLoadBalancer + "=1 to run")
	}
	f := testhelpers.RequireSdnFixture(t)

	t.Parallel()

	if testing.Short() {
		t.Skip("Skipping slow acceptance test in short mode")
	}

	name := acctest.RandomWithPrefix(t.Name())
	name = name[0:16] + name[len(name)-16:]

	// The SDN virtual network the LB's virtual_network_uid points at; created
	// per-test so parallel runs do not collide.
	network := `
resource "hpe_morpheus_network" "sdn" {
  name              = "` + name + `-net"
  cloud_id          = ` + f.CloudID + `
  group_id          = ` + f.GroupID + `
  type_code         = "juniperSdn"
  network_server_id = ` + f.NetworkServerID + `
  zone_pool_id      = ` + f.ClusterPoolID + `
  cidr              = "10.92.0.0/24"
  gateway           = "10.92.0.1"
  visibility        = "private"
  tenant_ids        = [1]
}
`

	lb, err := loadbalancer.RenderLoadBalancerJuniperSdnConfig(t, map[string]string{
		"Name":              name,
		"CloudId":           f.CloudID,
		"NetworkServerId":   f.NetworkServerID,
		"ResourcePoolId":    f.ClusterPoolID,
		"VirtualNetworkUid": "hpe_morpheus_network.sdn.unique_id",
		"VirtualIp":         "10.92.0.100",
		"SourcePort":        "80",
		"DestinationPort":   "8080",
		"Protocol":          "TCP",
	})
	if err != nil {
		t.Fatalf("failed to render config: %s", err)
	}

	// The example sets type_code; this variant drops the line.
	lb = strings.Replace(lb, "\n  type_code         = \"juniper-sdn-lb\"", "", 1)
	if strings.Contains(lb, "type_code") {
		t.Fatal("type_code still present in the rendered configuration")
	}

	resourceName := "hpe_morpheus_load_balancer.sdn"
	config := testhelpers.ProviderBlock() + network + lb

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "type_code", "juniper-sdn-lb"),
					resource.TestCheckResourceAttr(resourceName, "config_juniper_sdn.virtual_ip", "10.92.0.100"),
					resource.TestCheckResourceAttr(resourceName, "config_juniper_sdn.source_port", "80"),
					resource.TestCheckResourceAttr(resourceName, "config_juniper_sdn.destination_port", "8080"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			{
				// Nothing to plan: the computed type_code is kept from state.
				Config:   config,
				PlanOnly: true,
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"cloud_id",
					"network_server_id",
					"write_only_revision",
				},
			},
		},
	})
}
