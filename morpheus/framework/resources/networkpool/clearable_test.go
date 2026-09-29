// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkpool_test

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

// Removing the DNS domain, gateway, netmask or subnet address clears it on
// the pool (MORPH-17990); all are Optional+Computed and cleared once an apply
// has recorded them.
func TestAccMorpheusNetworkPoolRemovedAttributesAreCleared(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	name := acctest.RandomWithPrefix(t.Name())
	const rn = "hpe_morpheus_network_pool.test"

	config := func(with ...string) string {
		set := map[string]string{
			"dns_domain":     `dns_domain     = "clr.example.test"`,
			"gateway":        `gateway        = "10.77.1.1"`,
			"netmask":        `netmask        = "255.255.255.0"`,
			"subnet_address": `subnet_address = "10.77.1.0"`,
		}

		lines := make([]string, 0, len(with))
		for _, a := range with {
			lines = append(lines, "  "+set[a])
		}

		return testhelpers.ProviderBlock() + fmt.Sprintf(`
resource "hpe_morpheus_network_pool" "test" {
  name      = %q
  type_code = "morpheus"
%s

  ip_ranges = {
    starting_address = "10.77.1.10"
    ending_address   = "10.77.1.50"
  }
}
`, name, strings.Join(lines, "\n"))
	}

	all := []string{"dns_domain", "gateway", "netmask", "subnet_address"}

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
					resource.TestCheckResourceAttr(rn, "dns_domain", "clr.example.test"),
					resource.TestCheckResourceAttr(rn, "gateway", "10.77.1.1"),
					resource.TestCheckResourceAttr(rn, "netmask", "255.255.255.0"),
					resource.TestCheckResourceAttr(rn, "subnet_address", "10.77.1.0"),
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
