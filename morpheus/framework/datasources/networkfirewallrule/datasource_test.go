// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkfirewallrule_test

import (
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/networkfirewallrule"
	fwruleresource "github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/networkfirewallrule"
	fwrulegroupresource "github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/networkfirewallrulegroup"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
)

func TestMain(m *testing.M) {
	code := m.Run()
	testhelpers.WriteMergedResults()
	os.Exit(code)
}

// firewallRulePrereq renders a self-contained firewall rule group plus a firewall
// rule that references it. The data source tests read this freshly created rule
// instead of relying on a hard-coded rule ID that may not exist on the target
// environment.
func firewallRulePrereq(t *testing.T, name string) string {
	t.Helper()

	groupConfig, err := fwrulegroupresource.RenderNetworkFirewallRuleGroupConfig(t, map[string]string{
		"Name": name + "-group",
	})
	if err != nil {
		t.Fatal(err)
	}

	ruleConfig, err := fwruleresource.RenderNetworkFirewallRuleConfig(t, map[string]string{
		"Name":        name,
		"Priority":    "10",
		"Description": "data source acceptance test rule",
		"RuleGroupId": "hpe_morpheus_network_firewall_rule_group.example.id",
	})
	if err != nil {
		t.Fatal(err)
	}

	return groupConfig + ruleConfig
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

func TestAccMorpheusFindNetworkFirewallRuleByName(t *testing.T) {
	if capabilities.Missing(t, capabilities.NetworkFirewall) {
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
	prereq := firewallRulePrereq(t, name)

	// Read back the freshly created rule by name. Referencing the rule resource's
	// name attribute guarantees an exact match and an implicit dependency.
	dataSourceConfig := `
data "hpe_morpheus_network_firewall_rule" "example" {
  name                   = hpe_morpheus_network_firewall_rule.example.name
  network_integration_id = hpe_morpheus_network_firewall_rule.example.network_integration_id
}
`

	checks := networkFirewallRuleChecks()

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + prereq + dataSourceConfig,
				Check:  checkFn,
			},
		},
	})
}

func TestAccMorpheusFindNetworkFirewallRuleById(t *testing.T) {
	if capabilities.Missing(t, capabilities.NetworkFirewall) {
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
	prereq := firewallRulePrereq(t, name)

	// Read back the freshly created rule by ID. The id reference creates an
	// implicit dependency so the rule exists before the data source is read.
	dataSourceConfig := `
data "hpe_morpheus_network_firewall_rule" "example" {
  id                     = hpe_morpheus_network_firewall_rule.example.id
  network_integration_id = hpe_morpheus_network_firewall_rule.example.network_integration_id
}
`

	checks := networkFirewallRuleChecks()

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + prereq + dataSourceConfig,
				Check:  checkFn,
			},
		},
	})
}

func TestAccMorpheusFindNetworkFirewallRuleNotFound(t *testing.T) {
	if capabilities.Missing(t, capabilities.NetworkFirewall) {
		t.Log("Skipping test due to missing capabilities")

		return
	}
	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()

	dataSourceConfig, err := networkfirewallrule.RenderNetworkFirewallRuleByNameConfig(t,
		map[string]string{
			"Name": "______",
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	expected := regexp.MustCompile(`no network firewall rule found`)

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

func TestAccMorpheusFindNetworkFirewallRuleNoSearchAttrs(t *testing.T) {
	if capabilities.Missing(t, capabilities.NetworkFirewall) {
		t.Log("Skipping test due to missing capabilities")

		return
	}
	defer testhelpers.RecordResult(t)

	t.Parallel()

	config := providerConfigOffline + `
      data "hpe_morpheus_network_firewall_rule" "test" {
        network_integration_id = 1
      }`

	expected := `At least one attribute out of \[id,name\] must be specified`

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

func TestAccMorpheusFindNetworkFirewallRuleBothSearchAttrs(t *testing.T) {
	if capabilities.Missing(t, capabilities.NetworkFirewall) {
		t.Log("Skipping test due to missing capabilities")

		return
	}
	defer testhelpers.RecordResult(t)

	t.Parallel()

	config := providerConfigOffline + `
      data "hpe_morpheus_network_firewall_rule" "test" {
        id        = 1
        name      = "______"
        network_integration_id = 1
      }`

	expected := networkfirewallrule.ErrorRunningPreApply

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

func networkFirewallRuleChecks() []resource.TestCheckFunc {
	ds := "data.hpe_morpheus_network_firewall_rule.example"

	return []resource.TestCheckFunc{
		resource.TestCheckResourceAttrSet(ds, "id"),
		resource.TestCheckResourceAttrSet(ds, "name"),
		resource.TestCheckResourceAttrSet(ds, "network_integration_id"),
		resource.TestCheckResourceAttrSet(ds, "enabled"),
		resource.TestCheckResourceAttrSet(ds, "config.%"),
		resource.TestCheckResourceAttrSet(ds, "direction"),
		resource.TestCheckResourceAttrSet(ds, "policy"),
		resource.TestCheckResourceAttrSet(ds, "priority"),
		resource.TestCheckResourceAttrSet(ds, "source_type"),
		resource.TestCheckResourceAttrSet(ds, "destination_type"),
		resource.TestCheckResourceAttrSet(ds, "group_name"),
	}
}
