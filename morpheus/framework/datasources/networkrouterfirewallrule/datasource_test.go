// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkrouterfirewallrule_test

import (
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/networkrouterfirewallrule"
	networkrouterresource "github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/networkrouter"
	firewallruleresource "github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/networkrouterfirewallrule"
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

// routerFixture renders a self-contained NSX-T network router.
func routerFixture(t *testing.T, name string) string {
	t.Helper()

	cfg, err := networkrouterresource.RenderNetworkRouterGenericConfig(t, map[string]string{
		"Name":                 name + "-router",
		"TypeId":               "9",
		"GroupId":              "3",
		"NetworkIntegrationId": "5",
	})
	if err != nil {
		t.Fatal(err)
	}

	return cfg
}

func TestAccMorpheusFindNetworkRouterFirewallRuleByName(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.NetworkRouter, capabilities.NetworkFirewall)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	parentID := os.Getenv("TF_ACC_FIREWALL_RULE_GROUP_ID")
	if parentID == "" {
		t.Skip("TF_ACC_FIREWALL_RULE_GROUP_ID not set; skipping test requiring a known NSX-T firewall rule group")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()
	name := acctest.RandomWithPrefix(t.Name())

	routerConfig := routerFixture(t, name)

	resourceConfig, err := firewallruleresource.RenderNetworkRouterFirewallRuleConfig(t, map[string]string{
		"RouterId": "hpe_morpheus_network_router.example.id",
		"ParentId": parentID,
		"Name":     name,
	})
	if err != nil {
		t.Fatal(err)
	}

	dataSourceConfig := `
data "hpe_morpheus_network_router_firewall_rule" "example" {
  name       = "` + name + `"
  router_id  = hpe_morpheus_network_router.example.id
  depends_on = [hpe_morpheus_network_router_firewall_rule.example]
}
`

	checkFn := resource.ComposeAggregateTestCheckFunc(firewallRuleChecks()...)

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

func TestAccMorpheusFindNetworkRouterFirewallRuleById(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.NetworkRouter, capabilities.NetworkFirewall)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	parentID := os.Getenv("TF_ACC_FIREWALL_RULE_GROUP_ID")
	if parentID == "" {
		t.Skip("TF_ACC_FIREWALL_RULE_GROUP_ID not set; skipping test requiring a known NSX-T firewall rule group")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()
	name := acctest.RandomWithPrefix(t.Name())

	routerConfig := routerFixture(t, name)

	resourceConfig, err := firewallruleresource.RenderNetworkRouterFirewallRuleConfig(t, map[string]string{
		"RouterId": "hpe_morpheus_network_router.example.id",
		"ParentId": parentID,
		"Name":     name,
	})
	if err != nil {
		t.Fatal(err)
	}

	dataSourceConfig, err := networkrouterfirewallrule.RenderFirewallRuleByIdConfig(t, map[string]string{
		"Id":       "hpe_morpheus_network_router_firewall_rule.example.id",
		"RouterId": "hpe_morpheus_network_router.example.id",
	})
	if err != nil {
		t.Fatal(err)
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(firewallRuleChecks()...)

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

func TestAccMorpheusFindNetworkRouterFirewallRuleNotFound(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.NetworkRouter, capabilities.NetworkFirewall)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()
	name := acctest.RandomWithPrefix(t.Name())

	routerConfig := routerFixture(t, name)

	dataSourceConfig := `
data "hpe_morpheus_network_router_firewall_rule" "example" {
  name      = "nonexistent-firewall-rule-name-that-should-not-exist"
  router_id = hpe_morpheus_network_router.example.id
}
`

	expected := regexp.MustCompile(`no network router firewall rule found`)

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

func TestAccMorpheusFindNetworkRouterFirewallRuleNoSearchAttrs(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.NetworkRouter, capabilities.NetworkFirewall)

	t.Parallel()

	config := providerConfigOffline + `
      data "hpe_morpheus_network_router_firewall_rule" "test" {
        router_id = 1
      }`

	expected := networkrouterfirewallrule.ErrorNoValidSearchTerms

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

func firewallRuleChecks() []resource.TestCheckFunc {
	ds := "data.hpe_morpheus_network_router_firewall_rule.example"

	return []resource.TestCheckFunc{
		resource.TestCheckResourceAttrSet(ds, "id"),
		resource.TestCheckResourceAttrSet(ds, "router_id"),
		resource.TestCheckResourceAttrSet(ds, "name"),
		resource.TestCheckResourceAttrSet(ds, "enabled"),
		resource.TestCheckResourceAttrSet(ds, "policy"),
	}
}
