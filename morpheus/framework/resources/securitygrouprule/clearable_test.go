// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package securitygrouprule_test

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
func TestAccMorpheusSecurityGroupRuleRemovedAttributesAreCleared(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	name := acctest.RandomWithPrefix(t.Name())
	const rn = "hpe_morpheus_security_group_rule.test"

	config := func(with ...string) string {
		set := map[string]string{
			"name": `name              = "to be cleared"`,
		}

		lines := make([]string, 0, len(with))
		for _, a := range with {
			lines = append(lines, "  "+set[a])
		}

		return testhelpers.ProviderBlock() + fmt.Sprintf(`
resource "hpe_morpheus_security_group" "sg" {
  name = %q
}

resource "hpe_morpheus_security_group_rule" "test" {
  security_group_id = hpe_morpheus_security_group.sg.id
  protocol          = "tcp"
  rule_type         = "customRule"
  direction         = "ingress"
  port_range        = "8080"
  source            = "10.0.0.0/8"
  destination       = "10.1.0.0/16"
%s
}
`, name, strings.Join(lines, "\n"))
	}

	all := []string{"name"}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: config(all...),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(rn, "name"),
				),
			},
			{
				Config: config(),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(rn, plancheck.ResourceActionUpdate),
						plancheck.ExpectKnownValue(rn, tfjsonpath.New("name"), knownvalue.Null()),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckNoResourceAttr(rn, "name"),
				),
			},
			{
				Config:   config(),
				PlanOnly: true,
			},
		},
	})
}
