// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkfirewallrule_test

import (
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/networkfirewallrule"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/systemoverride"
)

func TestMain(m *testing.M) {
	systemoverride.ParseFlags()
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

func TestAccMorpheusFindNetworkFirewallRuleByName(t *testing.T) {
	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	testSystem := systemoverride.GetPreferred(t, "zodiac")
	providerConfig := testhelpers.ProviderBlockForServer(testSystem)

	dataSourceConfig, err := networkfirewallrule.RenderNetworkFirewallRuleByNameConfig(t, nil)
	if err != nil {
		t.Fatal(err)
	}

	checks := networkFirewallRuleChecks()

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + dataSourceConfig,
				Check:  checkFn,
			},
		},
	})
}

func TestAccMorpheusFindNetworkFirewallRuleById(t *testing.T) {
	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	testSystem := systemoverride.GetPreferred(t, "zodiac")
	providerConfig := testhelpers.ProviderBlockForServer(testSystem)

	dataSourceConfig, err := networkfirewallrule.RenderNetworkFirewallRuleByIdConfig(t, nil)
	if err != nil {
		t.Fatal(err)
	}

	checks := networkFirewallRuleChecks()

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + dataSourceConfig,
				Check:  checkFn,
			},
		},
	})
}

func TestAccMorpheusFindNetworkFirewallRuleNotFound(t *testing.T) {
	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	testSystem := systemoverride.GetPreferred(t, "zodiac")
	providerConfig := testhelpers.ProviderBlockForServer(testSystem)

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
