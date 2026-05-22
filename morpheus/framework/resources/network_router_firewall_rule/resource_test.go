package network_router_firewall_rule_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/HPE/terraform-provider-hpe/morpheus"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
)

func TestMain(m *testing.M) {
	code := m.Run()
	testhelpers.WriteMergedResults()
	os.Exit(code)
}

func TestAccMorpheusNetworkRouterFirewallRuleBasic(t *testing.T) {
	if capabilities.Missing(t, capabilities.All, capabilities.NetworkRouter, capabilities.NetworkFirewall) {
		t.Log("Skipping test due to missing capabilities")

		return
	}
	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	routerID := os.Getenv("TF_ACC_MORPHEUS_ROUTER_ID")
	if routerID == "" {
		t.Skip("TF_ACC_MORPHEUS_ROUTER_ID not set")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + fmt.Sprintf(`
resource "hpe_morpheus_network_router_firewall_rule" "test" {
  router_id = %s
  name      = %q
  policy    = "accept"
}
`, routerID, name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hpe_morpheus_network_router_firewall_rule.test", "name", name),
					resource.TestCheckResourceAttr("hpe_morpheus_network_router_firewall_rule.test", "policy", "accept"),
					resource.TestCheckResourceAttrSet("hpe_morpheus_network_router_firewall_rule.test", "id"),
				),
			},
			{
				ImportState:       true,
				ImportStateVerify: true,
				ResourceName:      "hpe_morpheus_network_router_firewall_rule.test",
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources["hpe_morpheus_network_router_firewall_rule.test"]
					if !ok {
						return "", fmt.Errorf("resource not found")
					}

					return rs.Primary.Attributes["router_id"] + "/" + rs.Primary.ID, nil
				},
			},
		},
	})
}

func TestAccMorpheusNetworkRouterFirewallRuleUpdate(t *testing.T) {
	if capabilities.Missing(t, capabilities.All, capabilities.NetworkRouter, capabilities.NetworkFirewall) {
		t.Log("Skipping test due to missing capabilities")

		return
	}
	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	routerID := os.Getenv("TF_ACC_MORPHEUS_ROUTER_ID")
	if routerID == "" {
		t.Skip("TF_ACC_MORPHEUS_ROUTER_ID not set")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + fmt.Sprintf(`
resource "hpe_morpheus_network_router_firewall_rule" "test" {
  router_id = %s
  name      = %q
  policy    = "accept"
  enabled   = true
}
`, routerID, name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hpe_morpheus_network_router_firewall_rule.test", "policy", "accept"),
					resource.TestCheckResourceAttr("hpe_morpheus_network_router_firewall_rule.test", "enabled", "true"),
				),
			},
			{
				Config: providerConfig + fmt.Sprintf(`
resource "hpe_morpheus_network_router_firewall_rule" "test" {
  router_id = %s
  name      = %q
  policy    = "deny"
  enabled   = false
}
`, routerID, name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hpe_morpheus_network_router_firewall_rule.test", "policy", "deny"),
					resource.TestCheckResourceAttr("hpe_morpheus_network_router_firewall_rule.test", "enabled", "false"),
				),
			},
		},
	})
}
