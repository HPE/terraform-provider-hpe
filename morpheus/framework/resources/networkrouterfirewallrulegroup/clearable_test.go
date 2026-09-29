// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkrouterfirewallrulegroup_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

// Removing the description clears it on the rule group (MORPH-17990); it is
// Optional+Computed and cleared once an apply has recorded it. The group sits
// on a tier-1 gateway the test builds on the NSX-T fixture. tenant_ids is not
// covered: the API does not return it, so state keeps the planned value and a
// live check would prove nothing.
func TestAccMorpheusNetworkRouterFirewallRuleGroupRemovedAttributesAreCleared(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.NetworkRouter, capabilities.NetworkFirewall)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	name := acctest.RandomWithPrefix(t.Name())
	const rn = "hpe_morpheus_network_router_firewall_rule_group.test"

	routerConfig := nsxtTier1RouterConfig(t, name)

	config := func(extra string) string {
		return testhelpers.ProviderBlock() + routerConfig + fmt.Sprintf(`
resource "hpe_morpheus_network_router_firewall_rule_group" "test" {
  router_id   = hpe_morpheus_network_router.fw_tier1.id
  name        = %q
  group_layer = "LocalGatewayRules"
%s
}
`, name, extra)
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: config(`  description = "to be cleared"`),
				Check:  resource.TestCheckResourceAttr(rn, "description", "to be cleared"),
			},
			{
				Config: config(""),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(rn, plancheck.ResourceActionUpdate),
						plancheck.ExpectKnownValue(rn, tfjsonpath.New("description"), knownvalue.Null()),
					},
				},
				Check: resource.TestCheckNoResourceAttr(rn, "description"),
			},
			{
				Config:   config(""),
				PlanOnly: true,
			},
		},
	})
}
