// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkrouternat_test

import (
	"fmt"
	"strings"
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

// Removing an optional attribute clears it on the NAT rule (MORPH-17990): the
// description and the source and destination matches, all Optional+Computed
// and cleared once an apply has recorded them. The rule is an SNAT rule on a
// tier-1 gateway the test builds on the NSX-T fixture; translated_network is
// what SNAT translates to and translated_ports applies to DNAT only, so
// neither is removed here.
func TestAccMorpheusNetworkRouterNatRemovedAttributesAreCleared(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.NetworkRouter)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	name := acctest.RandomWithPrefix(t.Name())
	const rn = "hpe_morpheus_network_router_nat.test"

	routerConfig := nsxtTier1RouterConfig(t, name)

	config := func(with ...string) string {
		set := map[string]string{
			"description":         `description         = "to be cleared"`,
			"source_network":      `source_network      = "10.0.0.0/24"`,
			"destination_network": `destination_network = "10.1.0.0/24"`,
		}

		lines := make([]string, 0, len(with))
		for _, a := range with {
			lines = append(lines, "  "+set[a])
		}

		return testhelpers.ProviderBlock() + routerConfig + fmt.Sprintf(`
resource "hpe_morpheus_network_router_nat" "test" {
  router_id          = hpe_morpheus_network_router.nat_tier1.id
  name               = %q
  action             = "SNAT"
  translated_network = "192.168.1.1"
%s
}
`, name, strings.Join(lines, "\n"))
	}

	all := []string{"description", "source_network", "destination_network"}

	nullChecks := []plancheck.PlanCheck{plancheck.ExpectResourceAction(rn, plancheck.ResourceActionUpdate)}
	var noAttrs []resource.TestCheckFunc
	for _, a := range all {
		nullChecks = append(nullChecks, plancheck.ExpectKnownValue(rn, tfjsonpath.New(a), knownvalue.Null()))
		noAttrs = append(noAttrs, resource.TestCheckNoResourceAttr(rn, a))
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: config(all...),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(rn, "description", "to be cleared"),
					resource.TestCheckResourceAttr(rn, "source_network", "10.0.0.0/24"),
					resource.TestCheckResourceAttr(rn, "destination_network", "10.1.0.0/24"),
				),
			},
			{
				Config:           config(),
				ConfigPlanChecks: resource.ConfigPlanChecks{PreApply: nullChecks},
				Check:            resource.ComposeAggregateTestCheckFunc(noAttrs...),
			},
			{
				Config:   config(),
				PlanOnly: true,
			},
		},
	})
}
