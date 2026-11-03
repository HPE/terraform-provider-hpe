// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package securitygrouprule_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

// TestAccMorpheusSecurityGroupRuleRequiredByCondition: a rule that leaves out
// an attribute Morpheus requires fails at plan, naming the attribute
// (MORPH-18279); the conditional cases Morpheus accepts still apply.
func TestAccMorpheusSecurityGroupRuleRequiredByCondition(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)
	testhelpers.SkipUnlessAcc(t)

	name := acctest.RandomWithPrefix(t.Name())

	rule := func(body string) string {
		return testhelpers.ProviderBlock() + `
resource "hpe_morpheus_security_group" "test" {
  name = "` + name + `"
}

resource "hpe_morpheus_security_group_rule" "rule" {
  security_group_id = hpe_morpheus_security_group.test.id
  name              = "` + name + `-rule"
  rule_type         = "customRule"
  direction         = "ingress"
  policy            = "accept"
` + body + `
}
`
	}

	addr := "hpe_morpheus_security_group_rule.rule"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			// The three omissions fail at plan, on the attribute.
			{
				Config:      rule(`  protocol = "tcp"` + "\n" + `  destination = "0.0.0.0/0"` + "\n" + `  port_range = "443"`),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`(?s)Missing source.*source_type is "cidr"`),
			},
			{
				Config:      rule(`  protocol = "tcp"` + "\n" + `  source = "0.0.0.0/0"` + "\n" + `  port_range = "443"`),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`(?s)Missing destination.*destination_type is "cidr"`),
			},
			{
				Config:      rule(`  protocol = "udp"` + "\n" + `  source = "0.0.0.0/0"` + "\n" + `  destination = "0.0.0.0/0"`),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`(?s)Missing port_range.*protocol "udp"`),
			},
			// What Morpheus accepts still applies: an ICMP rule has no ports.
			{
				Config: rule(`  protocol = "icmp"` + "\n" + `  source = "10.0.0.0/24"` + "\n" + `  destination = "0.0.0.0/0"`),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(addr, "protocol", "icmp"),
					resource.TestCheckNoResourceAttr(addr, "port_range"),
				),
			},
			{
				Config: rule(`  protocol = "tcp"` + "\n" + `  source = "10.0.0.0/24"` + "\n" +
					`  destination = "0.0.0.0/0"` + "\n" + `  port_range = "22"`),
				Check: resource.TestCheckResourceAttr(addr, "port_range", "22"),
			},
		},
	})
}
