// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkfirewallrulegroup_test

import (
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/networkfirewallrulegroup"
	fwrulegroupresource "github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/networkfirewallrulegroup"
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

func TestAccMorpheusNetworkFirewallRuleGroupByIdOk(t *testing.T) {
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

	groupConfig, err := fwrulegroupresource.RenderNetworkFirewallRuleGroupConfig(t, map[string]string{
		"Name": name,
	})
	if err != nil {
		t.Fatal(err)
	}

	// Read back the freshly created group by its ID. The id reference creates an
	// implicit dependency so the group exists before the data source is read,
	// keeping the test independent of any pre-existing group.
	dataSourceConfig := `
data "hpe_morpheus_network_firewall_rule_group" "example" {
  network_integration_id = hpe_morpheus_network_firewall_rule_group.example.network_integration_id
  id                     = hpe_morpheus_network_firewall_rule_group.example.id
}
`

	ds := "data.hpe_morpheus_network_firewall_rule_group.example"

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttrSet(ds, "id"),
		resource.TestCheckResourceAttr(ds, "name", name),
		resource.TestCheckResourceAttrSet(ds, "network_integration_id"),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + groupConfig + dataSourceConfig,
				Check:  checkFn,
			},
		},
	})
}

func TestAccMorpheusNetworkFirewallRuleGroupByNameOk(t *testing.T) {
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

	groupConfig, err := fwrulegroupresource.RenderNetworkFirewallRuleGroupConfig(t, map[string]string{
		"Name": name,
	})
	if err != nil {
		t.Fatal(err)
	}

	// Read back the freshly created group by its name. Referencing the resource's
	// name attribute guarantees an exact match and an implicit dependency.
	dataSourceConfig := `
data "hpe_morpheus_network_firewall_rule_group" "example" {
  network_integration_id = hpe_morpheus_network_firewall_rule_group.example.network_integration_id
  name                   = hpe_morpheus_network_firewall_rule_group.example.name
}
`

	ds := "data.hpe_morpheus_network_firewall_rule_group.example"

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttrSet(ds, "id"),
		resource.TestCheckResourceAttr(ds, "name", name),
		resource.TestCheckResourceAttrSet(ds, "network_integration_id"),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + groupConfig + dataSourceConfig,
				Check:  checkFn,
			},
		},
	})
}

func TestAccMorpheusNetworkFirewallRuleGroupNotFound(t *testing.T) {
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

	dataSourceConfig, err := networkfirewallrulegroup.RenderDataSourceByNameConfig(t, map[string]string{
		"Name": "______nonexistent______",
	})
	if err != nil {
		t.Fatal(err)
	}

	expected := regexp.MustCompile(networkfirewallrulegroup.ErrorNotFound)

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

func TestAccMorpheusNetworkFirewallRuleGroupNoSearchAttrs(t *testing.T) {
	if capabilities.Missing(t, capabilities.NetworkFirewall) {
		t.Log("Skipping test due to missing capabilities")

		return
	}
	defer testhelpers.RecordResult(t)

	t.Parallel()

	config := providerConfigOffline + `
      data "hpe_morpheus_network_firewall_rule_group" "test" {
        network_integration_id = 1
      }`

	expected := regexp.MustCompile(`At least one attribute out of \[id,name\] must be specified`)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: expected,
			},
		},
	})
}

func TestAccMorpheusNetworkFirewallRuleGroupBothSearchAttrs(t *testing.T) {
	if capabilities.Missing(t, capabilities.NetworkFirewall) {
		t.Log("Skipping test due to missing capabilities")

		return
	}
	defer testhelpers.RecordResult(t)

	t.Parallel()

	config := providerConfigOffline + `
      data "hpe_morpheus_network_firewall_rule_group" "test" {
        network_integration_id = 1
        id                     = 1
        name                   = "______"
      }`

	expected := regexp.MustCompile(`cannot be specified when`)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: expected,
			},
		},
	})
}
