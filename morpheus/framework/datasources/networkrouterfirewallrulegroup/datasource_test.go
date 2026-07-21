// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkrouterfirewallrulegroup_test

import (
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/networkrouterfirewallrulegroup"
	firewallrulegroupresource "github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/networkrouterfirewallrulegroup"
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
    url      = ""
    username = ""
    password = ""
  }
}
`

// nsxtTier1RouterConfig renders a self-contained NSX-T Tier-1 gateway router.
// QA constants: group_id 3, network_integration_id 5.
func nsxtTier1RouterConfig(name string) string {
	return `
resource "hpe_morpheus_network_router" "fw_tier1" {
  name                   = "` + name + `-tier1"
  group_id               = 3
  network_integration_id = 5

  config_nsxt_gateway_tier1 = {
    ip_management_type = "dhcpLocal"
  }
}
`
}

func TestAccMorpheusFindNetworkRouterFirewallRuleGroupByName(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.NetworkRouter, capabilities.NetworkFirewall)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()
	name := acctest.RandomWithPrefix(t.Name())

	routerConfig := nsxtTier1RouterConfig(name)

	resourceConfig, err := firewallrulegroupresource.RenderNetworkRouterFirewallRuleGroupConfig(t, map[string]string{
		"RouterId": "hpe_morpheus_network_router.fw_tier1.id",
		"Name":     name,
	})
	if err != nil {
		t.Fatal(err)
	}

	dataSourceConfig := `
data "hpe_morpheus_network_router_firewall_rule_group" "example" {
  name      = "` + name + `"
  router_id = hpe_morpheus_network_router.fw_tier1.id
  depends_on = [hpe_morpheus_network_router_firewall_rule_group.example]
}
`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + routerConfig + resourceConfig + dataSourceConfig,
				Check:  resource.ComposeAggregateTestCheckFunc(ruleGroupDataSourceChecks()...),
			},
		},
	})
}

func TestAccMorpheusFindNetworkRouterFirewallRuleGroupById(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.NetworkRouter, capabilities.NetworkFirewall)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()
	name := acctest.RandomWithPrefix(t.Name())

	routerConfig := nsxtTier1RouterConfig(name)

	resourceConfig, err := firewallrulegroupresource.RenderNetworkRouterFirewallRuleGroupConfig(t, map[string]string{
		"RouterId": "hpe_morpheus_network_router.fw_tier1.id",
		"Name":     name,
	})
	if err != nil {
		t.Fatal(err)
	}

	dataSourceConfig, err := networkrouterfirewallrulegroup.RenderFirewallRuleGroupByIdConfig(t, map[string]string{
		"Id":       "hpe_morpheus_network_router_firewall_rule_group.example.id",
		"RouterId": "hpe_morpheus_network_router.fw_tier1.id",
	})
	if err != nil {
		t.Fatal(err)
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + routerConfig + resourceConfig + dataSourceConfig,
				Check:  resource.ComposeAggregateTestCheckFunc(ruleGroupDataSourceChecks()...),
			},
		},
	})
}

func TestAccMorpheusFindNetworkRouterFirewallRuleGroupNotFound(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.NetworkRouter, capabilities.NetworkFirewall)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()
	name := acctest.RandomWithPrefix(t.Name())

	routerConfig := nsxtTier1RouterConfig(name)

	dataSourceConfig := `
data "hpe_morpheus_network_router_firewall_rule_group" "example" {
  name      = "nonexistent-firewall-rule-group-name-that-should-not-exist"
  router_id = hpe_morpheus_network_router.fw_tier1.id
}
`

	expected := regexp.MustCompile(networkrouterfirewallrulegroup.ErrorNoNetworkRouterFirewallRuleGroupFound)

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

func TestAccMorpheusFindNetworkRouterFirewallRuleGroupNoSearchAttrs(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.NetworkRouter, capabilities.NetworkFirewall)

	t.Parallel()

	config := providerConfigOffline + `
data "hpe_morpheus_network_router_firewall_rule_group" "test" {
  router_id = 1
}`

	expected := networkrouterfirewallrulegroup.ErrorNoValidSearchTerms

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile(regexp.QuoteMeta(expected)),
			},
		},
	})
}

func ruleGroupDataSourceChecks() []resource.TestCheckFunc {
	ds := "data.hpe_morpheus_network_router_firewall_rule_group.example"

	return []resource.TestCheckFunc{
		resource.TestCheckResourceAttrSet(ds, "id"),
		resource.TestCheckResourceAttrSet(ds, "router_id"),
		resource.TestCheckResourceAttrSet(ds, "name"),
		resource.TestCheckResourceAttrSet(ds, "external_id"),
		resource.TestCheckResourceAttrSet(ds, "status"),
	}
}
