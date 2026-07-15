// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkrouterfirewallrulegroups_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/framework/datasources/networkrouterfirewallrulegroups"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

func TestMain(m *testing.M) {
	code := testhelpers.TestMain(m)
	testhelpers.WriteMergedResults()
	os.Exit(code)
}

// TestAccMorpheusNetworkRouterFirewallRuleGroupsBasic verifies the data source
// can be read against a live NSX-T router that has firewall rule groups.
// This test requires the NetworkFirewall capability because it targets a router
// with NSX-T distributed firewall rule groups.
func TestAccMorpheusNetworkRouterFirewallRuleGroupsBasic(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.NetworkFirewall)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	routerID := os.Getenv("TF_ACC_ROUTER_ID")
	if routerID == "" {
		t.Skip("TF_ACC_ROUTER_ID not set; skipping test requiring a known NSX-T router with firewall rule groups")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()

	dataSourceConfig, err := networkrouterfirewallrulegroups.RenderConfig(t, map[string]string{
		"RouterId": routerID,
	})
	if err != nil {
		t.Fatal(err)
	}

	config := providerConfig + dataSourceConfig

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttrSet(
			"data.hpe_morpheus_network_router_firewall_rule_groups.example", "router_id",
		),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check:  resource.ComposeAggregateTestCheckFunc(checks...),
			},
		},
	})
}

// TestAccMorpheusNetworkRouterFirewallRuleGroupsWithKnownRouter verifies the
// data source against a known router ID from the environment. Set
// TF_ACC_ROUTER_ID to the ID of an NSX-T router with firewall rule groups.
func TestAccMorpheusNetworkRouterFirewallRuleGroupsWithKnownRouter(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.NetworkFirewall)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	routerID := os.Getenv("TF_ACC_ROUTER_ID")
	if routerID == "" {
		t.Skip("TF_ACC_ROUTER_ID not set; skipping test requiring a known NSX-T router")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()

	dataSourceConfig, err := networkrouterfirewallrulegroups.RenderConfig(t, map[string]string{
		"RouterId": routerID,
	})
	if err != nil {
		t.Fatal(err)
	}

	config := providerConfig + dataSourceConfig

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttrSet(
			"data.hpe_morpheus_network_router_firewall_rule_groups.example", "router_id",
		),
		resource.TestCheckResourceAttrSet(
			"data.hpe_morpheus_network_router_firewall_rule_groups.example", "rule_groups.#",
		),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check:  resource.ComposeAggregateTestCheckFunc(checks...),
			},
		},
	})
}
