// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package securitygroup_test

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

// Removing an optional attribute clears it (MORPH-17990): the attributes here
// are Optional without Computed, so Terraform plans the null and the apply
// used to fail with an inconsistent result because Update never sent it.
func TestAccMorpheusSecurityGroupRemovedAttributesAreCleared(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	name := acctest.RandomWithPrefix(t.Name())
	const rn = "hpe_morpheus_security_group.test"

	config := func(with ...string) string {
		set := map[string]string{
			"description":                    `description                    = "to be cleared"`,
			"resource_permission_groups_all": `resource_permission_groups_all = false`,
			"resource_permission_group_ids":  `resource_permission_group_ids  = [1]`,
		}

		lines := make([]string, 0, len(with))
		for _, a := range with {
			lines = append(lines, "  "+set[a])
		}

		return testhelpers.ProviderBlock() + fmt.Sprintf(`
resource "hpe_morpheus_security_group" "test" {
  name = %q
%s
}
`, name, strings.Join(lines, "\n"))
	}

	all := []string{"description", "resource_permission_groups_all", "resource_permission_group_ids"}

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
					resource.TestCheckResourceAttr(rn, "resource_permission_groups_all", "false"),
					resource.TestCheckResourceAttr(rn, "resource_permission_group_ids.#", "1"),
				),
			},
			{
				// The group list removed while the "all" flag stays configured.
				Config: config("resource_permission_groups_all"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(rn, plancheck.ResourceActionUpdate),
						plancheck.ExpectKnownValue(rn, tfjsonpath.New("description"), knownvalue.Null()),
						plancheck.ExpectKnownValue(rn, tfjsonpath.New("resource_permission_group_ids"), knownvalue.Null()),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckNoResourceAttr(rn, "description"),
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
				// The group list and the "all" flag removed together: the flag
				// keeps its prior value in the plan (UseStateForUnknown), so the
				// permissions object is still sent, with the empty list.
				Config: config(),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(rn, plancheck.ResourceActionUpdate),
						plancheck.ExpectKnownValue(rn, tfjsonpath.New("resource_permission_groups_all"), knownvalue.Bool(false)),
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
