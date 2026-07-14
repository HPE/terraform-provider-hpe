package networkrouterfirewallrulegroup_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/networkrouter"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/networkrouterfirewallrulegroup"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

func TestMain(m *testing.M) {
	code := testhelpers.TestMain(m)
	testhelpers.WriteMergedResults()
	os.Exit(code)
}

func TestAccMorpheusNetworkRouterFirewallRuleGroupResourceExampleOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.NetworkRouter, capabilities.NetworkFirewall)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()
	name := acctest.RandomWithPrefix(t.Name())
	resourceName := "hpe_morpheus_network_router_firewall_rule_group.example"

	routerConfig, err := networkrouter.RenderNetworkRouterGenericConfig(t, map[string]string{
		"Name":                 name + "-router",
		"TypeId":               "9",
		"GroupId":              "3",
		"NetworkIntegrationId": "5",
	})
	if err != nil {
		t.Fatal(err)
	}

	resourceConfig, err := networkrouterfirewallrulegroup.RenderNetworkRouterFirewallRuleGroupConfig(t, map[string]string{
		"RouterId":     "hpe_morpheus_network_router.example.id",
		"Name":         name,
		"ExternalType": "SecurityPolicy",
	})
	if err != nil {
		t.Fatal(err)
	}

	checks := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttrPair(resourceName, "router_id", "hpe_morpheus_network_router.example", "id"),
		resource.TestCheckResourceAttr(resourceName, "name", name),
		resource.TestCheckResourceAttr(resourceName, "external_type", "SecurityPolicy"),
		resource.TestCheckResourceAttrSet(resourceName, "id"),
		resource.TestCheckResourceAttrSet(resourceName, "external_id"),
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + routerConfig + resourceConfig,
				Check:  checks,
			},
			{
				// Verify no-op plan after apply.
				Config:             providerConfig + routerConfig + resourceConfig,
				ExpectNonEmptyPlan: false,
				PlanOnly:           true,
			},
		},
	})
}

func TestAccMorpheusNetworkRouterFirewallRuleGroupResourceUpdateOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.NetworkRouter, capabilities.NetworkFirewall)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()
	name := acctest.RandomWithPrefix(t.Name())
	resourceName := "hpe_morpheus_network_router_firewall_rule_group.example"

	routerConfig, err := networkrouter.RenderNetworkRouterGenericConfig(t, map[string]string{
		"Name":                 name + "-router",
		"TypeId":               "9",
		"GroupId":              "3",
		"NetworkIntegrationId": "5",
	})
	if err != nil {
		t.Fatal(err)
	}

	createConfig, err := networkrouterfirewallrulegroup.RenderNetworkRouterFirewallRuleGroupConfig(t, map[string]string{
		"RouterId":     "hpe_morpheus_network_router.example.id",
		"Name":         name,
		"ExternalType": "SecurityPolicy",
	})
	if err != nil {
		t.Fatal(err)
	}

	updatedName := name + "-updated"
	updateConfig := fmt.Sprintf(`
resource "hpe_morpheus_network_router_firewall_rule_group" "example" {
  router_id     = hpe_morpheus_network_router.example.id
  name          = %q
  external_type = "SecurityPolicy"
  priority      = 10
  visibility    = "private"
}
`, updatedName)

	createChecks := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttrPair(resourceName, "router_id", "hpe_morpheus_network_router.example", "id"),
		resource.TestCheckResourceAttr(resourceName, "name", name),
		resource.TestCheckResourceAttr(resourceName, "external_type", "SecurityPolicy"),
	)

	updateChecks := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttrPair(resourceName, "router_id", "hpe_morpheus_network_router.example", "id"),
		resource.TestCheckResourceAttr(resourceName, "name", updatedName),
		resource.TestCheckResourceAttr(resourceName, "priority", "10"),
		resource.TestCheckResourceAttr(resourceName, "visibility", "private"),
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

func TestAccMorpheusNetworkRouterFirewallRuleGroupResourceImportOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.NetworkRouter, capabilities.NetworkFirewall)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()
	name := acctest.RandomWithPrefix(t.Name())
	resourceName := "hpe_morpheus_network_router_firewall_rule_group.example"

	routerConfig, err := networkrouter.RenderNetworkRouterGenericConfig(t, map[string]string{
		"Name":                 name + "-router",
		"TypeId":               "9",
		"GroupId":              "3",
		"NetworkIntegrationId": "5",
	})
	if err != nil {
		t.Fatal(err)
	}

	resourceConfig, err := networkrouterfirewallrulegroup.RenderNetworkRouterFirewallRuleGroupConfig(t, map[string]string{
		"RouterId":     "hpe_morpheus_network_router.example.id",
		"Name":         name,
		"ExternalType": "SecurityPolicy",
	})
	if err != nil {
		t.Fatal(err)
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + routerConfig + resourceConfig,
			},
			{
				ImportState:       true,
				ImportStateVerify: false, // write-only fields (external_type, visibility, tenant_ids) will be null after import
				ResourceName:      resourceName,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[resourceName]
					if !ok {
						return "", fmt.Errorf("resource not found")
					}

					return rs.Primary.Attributes["router_id"] + "." + rs.Primary.Attributes["id"], nil
				},
			},
		},
	})
}
