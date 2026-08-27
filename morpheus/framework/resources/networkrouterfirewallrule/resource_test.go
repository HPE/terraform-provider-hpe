package networkrouterfirewallrule_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/networkrouterfirewallrule"
	firewallrulegroupresource "github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/networkrouterfirewallrulegroup"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

// firewallRuleFixture renders a per-test NSX-T Tier-1 gateway router labelled
// hpe_morpheus_network_router.fw_tier1, plus a firewall rule group on it
// labelled hpe_morpheus_network_router_firewall_rule_group.example.
//
// parent_id is Required and RequiresReplace: NSX-T rejects a rule without a
// parent group. Gateway firewall rule groups attach to the router's gateway
// policy, which does not require an edge cluster or tier-0 connection -- a
// simple dhcpLocal Tier-1 is sufficient. This mirrors the fixture used by the
// firewall rule group resource/data source tests.
//
// QA constants: group_id 3, network_integration_id 5.
func firewallRuleFixture(t *testing.T, name string) string {
	t.Helper()

	routerConfig := `
resource "hpe_morpheus_network_router" "fw_tier1" {
  name                   = "` + name + `-tier1"
  group_id               = 3
  network_integration_id = 5

  config_nsxt_gateway_tier1 = {
    ip_management_type = "dhcpLocal"
  }
}
`

	groupConfig, err := firewallrulegroupresource.RenderNetworkRouterFirewallRuleGroupConfig(
		t, map[string]string{
			"RouterId": "hpe_morpheus_network_router.fw_tier1.id",
			"Name":     name + "-group",
		})
	if err != nil {
		t.Fatal(err)
	}

	return routerConfig + groupConfig
}

func TestMain(m *testing.M) {
	code := testhelpers.TestMain(m)
	testhelpers.WriteMergedResults()
	os.Exit(code)
}

func TestAccMorpheusNetworkRouterFirewallRuleResourceExampleOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.NetworkRouter, capabilities.NetworkFirewall)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()
	name := acctest.RandomWithPrefix(t.Name())
	resourceName := "hpe_morpheus_network_router_firewall_rule.example"

	routerConfig := firewallRuleFixture(t, name)

	resourceConfig, err := networkrouterfirewallrule.RenderNetworkRouterFirewallRuleConfig(t, map[string]string{
		"RouterId": "hpe_morpheus_network_router.fw_tier1.id",
		"Name":     name,
	})
	if err != nil {
		t.Fatal(err)
	}

	checks := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttrPair(resourceName, "router_id", "hpe_morpheus_network_router.fw_tier1", "id"),
		resource.TestCheckResourceAttr(resourceName, "name", name),
		resource.TestCheckResourceAttr(resourceName, "policy", "accept"),
		resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
		resource.TestCheckResourceAttrSet(resourceName, "id"),
		// description is Optional and never returned by the API. It must settle
		// as null when unset: when it was Optional+Computed the apply failed
		// with "Provider returned invalid result object after apply".
		resource.TestCheckNoResourceAttr(resourceName, "description"),
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + routerConfig + resourceConfig,
				Check:  checks,
			},
			{
				Config:             providerConfig + routerConfig + resourceConfig,
				ExpectNonEmptyPlan: false,
				PlanOnly:           true,
			},
			{
				ImportState:       true,
				ImportStateVerify: true,
				ResourceName:      "hpe_morpheus_network_router_firewall_rule.example",
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources["hpe_morpheus_network_router_firewall_rule.example"]
					if !ok {
						return "", fmt.Errorf("resource not found")
					}

					return rs.Primary.Attributes["router_id"] + "." + rs.Primary.Attributes["id"], nil
				},
			},
		},
	})
}

func TestAccMorpheusNetworkRouterFirewallRuleResourceUpdateOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.NetworkRouter, capabilities.NetworkFirewall)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()
	name := acctest.RandomWithPrefix(t.Name())
	resourceName := "hpe_morpheus_network_router_firewall_rule.example"

	routerConfig := firewallRuleFixture(t, name)

	createConfig, err := networkrouterfirewallrule.RenderNetworkRouterFirewallRuleConfig(t, map[string]string{
		"RouterId": "hpe_morpheus_network_router.fw_tier1.id",
		"Name":     name,
	})
	if err != nil {
		t.Fatal(err)
	}

	updateConfig := `
resource "hpe_morpheus_network_router_firewall_rule" "example" {
  router_id = hpe_morpheus_network_router.fw_tier1.id
  parent_id = hpe_morpheus_network_router_firewall_rule_group.example.external_id
  name      = "` + name + `"
  policy    = "deny"
  enabled   = false
}
`

	createChecks := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttrPair(resourceName, "router_id", "hpe_morpheus_network_router.fw_tier1", "id"),
		resource.TestCheckResourceAttr(resourceName, "name", name),
		resource.TestCheckResourceAttr(resourceName, "policy", "accept"),
		resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
	)

	updateChecks := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttrPair(resourceName, "router_id", "hpe_morpheus_network_router.fw_tier1", "id"),
		resource.TestCheckResourceAttr(resourceName, "name", name),
		resource.TestCheckResourceAttr(resourceName, "policy", "deny"),
		resource.TestCheckResourceAttr(resourceName, "enabled", "false"),
	)

	checkInPlaceUpdate := resource.ConfigPlanChecks{
		PreApply: []plancheck.PlanCheck{
			plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate),
		},
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{Config: providerConfig + routerConfig + createConfig, Check: createChecks},
			{Config: providerConfig + routerConfig + updateConfig, Check: updateChecks, ConfigPlanChecks: checkInPlaceUpdate},
			{Config: providerConfig + routerConfig + updateConfig, ExpectNonEmptyPlan: false, PlanOnly: true},
		},
	})
}

// TestAccMorpheusNetworkRouterFirewallRuleResourceDescriptionOk exercises the
// full lifecycle of the description attribute: absent on create, added, then
// removed again.
//
// description is accepted by the API on create and update but is never returned
// by the GET, so every assertion here is necessarily state-only. In particular
// the final step cannot prove the appliance actually cleared the description --
// only that Terraform's state and plan are consistent. Confirming the appliance
// side needs a manual check in the Morpheus UI.
func TestAccMorpheusNetworkRouterFirewallRuleResourceDescriptionOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.NetworkRouter, capabilities.NetworkFirewall)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()
	name := acctest.RandomWithPrefix(t.Name())
	resourceName := "hpe_morpheus_network_router_firewall_rule.example"

	routerConfig := firewallRuleFixture(t, name)

	ruleConfig := func(description string) string {
		descriptionLine := ""
		if description != "" {
			descriptionLine = `  description = "` + description + `"` + "\n"
		}

		return `
resource "hpe_morpheus_network_router_firewall_rule" "example" {
  router_id = hpe_morpheus_network_router.fw_tier1.id
  parent_id = hpe_morpheus_network_router_firewall_rule_group.example.external_id
  name      = "` + name + `"
` + descriptionLine + `}
`
	}

	checkInPlaceUpdate := resource.ConfigPlanChecks{
		PreApply: []plancheck.PlanCheck{
			plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate),
		},
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			// Created without a description: the value must be known (null)
			// after apply, not unknown.
			{
				Config: providerConfig + routerConfig + ruleConfig(""),
				Check:  resource.TestCheckNoResourceAttr(resourceName, "description"),
			},
			// Adding a description is an in-place update.
			{
				Config:           providerConfig + routerConfig + ruleConfig("managed by terraform"),
				Check:            resource.TestCheckResourceAttr(resourceName, "description", "managed by terraform"),
				ConfigPlanChecks: checkInPlaceUpdate,
			},
			// Removing it again clears it on the appliance and returns state to
			// null, leaving no residual diff.
			{
				Config:           providerConfig + routerConfig + ruleConfig(""),
				Check:            resource.TestCheckNoResourceAttr(resourceName, "description"),
				ConfigPlanChecks: checkInPlaceUpdate,
			},
			{
				Config:             providerConfig + routerConfig + ruleConfig(""),
				ExpectNonEmptyPlan: false,
				PlanOnly:           true,
			},
		},
	})
}
