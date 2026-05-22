package security_group_rule_test

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

func TestAccMorpheusSecurityGroupRuleBasic(t *testing.T) {
	if capabilities.Missing(t, capabilities.All) {
		t.Log("Skipping test due to missing capabilities")

		return
	}
	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()

	sgName := acctest.RandomWithPrefix(t.Name())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + fmt.Sprintf(`
resource "hpe_morpheus_security_group" "test" {
  name    = %q
}

resource "hpe_morpheus_security_group_rule" "test" {
  security_group_id = hpe_morpheus_security_group.test.id
  protocol          = "tcp"
  rule_type         = "customRule"
  source            = "0.0.0.0/0"
  destination       = "10.0.0.0/8"
  port_range        = "80-443"
}
`, sgName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hpe_morpheus_security_group_rule.test", "protocol", "tcp"),
					resource.TestCheckResourceAttr("hpe_morpheus_security_group_rule.test", "rule_type", "customRule"),
					resource.TestCheckResourceAttr("hpe_morpheus_security_group_rule.test", "source", "0.0.0.0/0"),
					resource.TestCheckResourceAttrSet("hpe_morpheus_security_group_rule.test", "id"),
				),
			},
			{
				ImportState:       true,
				ImportStateVerify: true,
				ResourceName:      "hpe_morpheus_security_group_rule.test",
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources["hpe_morpheus_security_group_rule.test"]
					if !ok {
						return "", fmt.Errorf("resource not found")
					}

					return rs.Primary.Attributes["security_group_id"] + "/" + rs.Primary.ID, nil
				},
			},
		},
	})
}

func TestAccMorpheusSecurityGroupRuleUpdate(t *testing.T) {
	if capabilities.Missing(t, capabilities.All) {
		t.Log("Skipping test due to missing capabilities")

		return
	}
	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()

	sgName := acctest.RandomWithPrefix(t.Name())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + fmt.Sprintf(`
resource "hpe_morpheus_security_group" "test" {
  name    = %q
}

resource "hpe_morpheus_security_group_rule" "test" {
  security_group_id = hpe_morpheus_security_group.test.id
  protocol          = "tcp"
  rule_type         = "customRule"
  source            = "0.0.0.0/0"
  destination       = "10.0.0.0/8"
  port_range        = "80-443"
}
`, sgName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hpe_morpheus_security_group_rule.test", "protocol", "tcp"),
				),
			},
			{
				Config: providerConfig + fmt.Sprintf(`
resource "hpe_morpheus_security_group" "test" {
  name    = %q
}

resource "hpe_morpheus_security_group_rule" "test" {
  security_group_id = hpe_morpheus_security_group.test.id
  protocol          = "udp"
  rule_type         = "customRule"
  source            = "10.0.0.0/8"
  destination       = "192.168.0.0/16"
  port_range        = "53"
}
`, sgName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hpe_morpheus_security_group_rule.test", "protocol", "udp"),
					resource.TestCheckResourceAttr("hpe_morpheus_security_group_rule.test", "source", "10.0.0.0/8"),
				),
			},
		},
	})
}
