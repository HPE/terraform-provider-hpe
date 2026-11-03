// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package subnet_test

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

// Removing description, labels or resource_permission_group_ids clears it on
// the subnet (MORPH-17990). All three are Optional without Computed: the
// apply used to fail with an inconsistent result because Update never sent
// the removal.
func TestAccMorpheusSubnetRemovedAttributesAreCleared(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.Subnet)

	// NOTE: not t.Parallel() — subnet tests share an Azure VNet, which
	// serialises per-VNet operations; this test uses its own free /24.
	name := acctest.RandomWithPrefix(t.Name())
	const rn = "hpe_morpheus_subnet.test"

	config := func(with ...string) string {
		set := map[string]string{
			"description":                    `description                    = "to be cleared"`,
			"labels":                         `labels                         = ["terraform", "acctest"]`,
			"resource_permission_groups_all": `resource_permission_groups_all = false`,
			"resource_permission_group_ids":  `resource_permission_group_ids  = [1]`,
		}

		lines := make([]string, 0, len(with))
		for _, a := range with {
			lines = append(lines, "  "+set[a])
		}

		// A VNet of the QA Azure cloud without a delete lock, found by name;
		// 7 is the Azure subnet type.
		return testhelpers.ProviderBlock() + fmt.Sprintf(`
data "hpe_morpheus_network" "vnet" {
  name     = "vnet-westus-1"
  cloud_id = 4617
}

resource "hpe_morpheus_subnet" "test" {
  name                           = %q
  type_id                        = 7
  network_id                     = data.hpe_morpheus_network.vnet.id
  visibility                     = "private"
%s

  config = {
    subnetName = %q
    subnetCidr = "172.16.251.0/24"
  }
}
`, name, strings.Join(lines, "\n"), name)
	}

	all := []string{"description", "labels", "resource_permission_groups_all", "resource_permission_group_ids"}

	groupsCleared := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr(rn, "resource_permission_groups_all", "false"),
		resource.TestCheckNoResourceAttr(rn, "resource_permission_group_ids.#"),
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: config(all...),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(rn, "description", "to be cleared"),
					resource.TestCheckResourceAttr(rn, "labels.#", "2"),
					resource.TestCheckResourceAttr(rn, "resource_permission_groups_all", "false"),
					resource.TestCheckResourceAttr(rn, "resource_permission_group_ids.#", "1"),
				),
			},
			{
				// The group list removed while the "all" flag stays configured:
				// the permissions object is sent with the flag as planned.
				Config: config("resource_permission_groups_all"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(rn, plancheck.ResourceActionUpdate),
						plancheck.ExpectKnownValue(rn, tfjsonpath.New("description"), knownvalue.Null()),
						plancheck.ExpectKnownValue(rn, tfjsonpath.New("labels"), knownvalue.Null()),
						plancheck.ExpectKnownValue(rn, tfjsonpath.New("resource_permission_group_ids"), knownvalue.Null()),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckNoResourceAttr(rn, "description"),
					resource.TestCheckNoResourceAttr(rn, "labels.#"),
					groupsCleared,
				),
			},
			{
				Config:   config("resource_permission_groups_all"),
				PlanOnly: true,
			},
			{
				// The groups back, to remove them again the other way.
				Config: config("resource_permission_groups_all", "resource_permission_group_ids"),
				Check:  resource.TestCheckResourceAttr(rn, "resource_permission_group_ids.#", "1"),
			},
			{
				// The group list and the "all" flag removed together: the flag,
				// Optional+Computed, is unknown in the plan, and the permissions
				// object is sent with the flag as it is in state so the empty
				// list reaches the API.
				Config: config(),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(rn, plancheck.ResourceActionUpdate),
						plancheck.ExpectUnknownValue(rn, tfjsonpath.New("resource_permission_groups_all")),
						plancheck.ExpectKnownValue(rn, tfjsonpath.New("resource_permission_group_ids"), knownvalue.Null()),
					},
				},
				Check: groupsCleared,
			},
			{
				Config:   config(),
				PlanOnly: true,
			},
		},
	})
}
