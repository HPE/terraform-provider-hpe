package security_group_rule_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/HPE/terraform-provider-hpe/morpheus"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/security_group_rule"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
)

func TestMain(m *testing.M) {
	code := testhelpers.TestMain(m)
	testhelpers.WriteMergedResults()
	os.Exit(code)
}

func TestAccMorpheusSecurityGroupRuleResourceExampleOk(t *testing.T) {
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
	name := acctest.RandomWithPrefix(t.Name())

	securityGroupConfig := `
resource "hpe_morpheus_security_group" "test" {
  name = "` + name + `"
}
`

	resourceConfig, err := security_group_rule.RenderSecurityGroupRuleConfig(t, map[string]string{
		"SecurityGroupId": "hpe_morpheus_security_group.test.id",
		"Name":            name + "-rule",
	})
	if err != nil {
		t.Fatal(err)
	}

	checks := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttrPair(
			"hpe_morpheus_security_group_rule.example",
			"security_group_id",
			"hpe_morpheus_security_group.test",
			"id",
		),
		resource.TestCheckResourceAttr("hpe_morpheus_security_group_rule.example", "name", name+"-rule"),
		resource.TestCheckResourceAttr("hpe_morpheus_security_group_rule.example", "protocol", "tcp"),
		resource.TestCheckResourceAttr("hpe_morpheus_security_group_rule.example", "rule_type", "customRule"),
		resource.TestCheckResourceAttr("hpe_morpheus_security_group_rule.example", "direction", "ingress"),
		resource.TestCheckResourceAttr("hpe_morpheus_security_group_rule.example", "port_range", "443"),
		resource.TestCheckResourceAttr("hpe_morpheus_security_group_rule.example", "source", "0.0.0.0/0"),
		resource.TestCheckResourceAttr("hpe_morpheus_security_group_rule.example", "destination", "0.0.0.0/0"),
		resource.TestCheckResourceAttr("hpe_morpheus_security_group_rule.example", "policy", "accept"),
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + securityGroupConfig + resourceConfig,
				Check:  checks,
			},
			{
				Config:             providerConfig + securityGroupConfig + resourceConfig,
				ExpectNonEmptyPlan: false,
				PlanOnly:           true,
			},
			{
				ImportState:       true,
				ImportStateVerify: true,
				ResourceName:      "hpe_morpheus_security_group_rule.example",
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources["hpe_morpheus_security_group_rule.example"]
					if !ok {
						return "", fmt.Errorf("resource not found")
					}

					return rs.Primary.Attributes["security_group_id"] + "." + rs.Primary.Attributes["id"], nil
				},
			},
		},
	})
}

func TestAccMorpheusSecurityGroupRuleResourceUpdateOk(t *testing.T) {
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
	name := acctest.RandomWithPrefix(t.Name())

	securityGroupConfig := `
resource "hpe_morpheus_security_group" "test" {
  name = "` + name + `"
}
`

	createConfig, err := security_group_rule.RenderSecurityGroupRuleConfig(t, map[string]string{
		"SecurityGroupId": "hpe_morpheus_security_group.test.id",
		"Name":            name + "-rule",
	})
	if err != nil {
		t.Fatal(err)
	}

	updateConfig, err := security_group_rule.RenderSecurityGroupRuleConfig(t, map[string]string{
		"SecurityGroupId": "hpe_morpheus_security_group.test.id",
		"Name":            name + "-rule",
		"Protocol":        "udp",
		"PortRange":       "53",
		"Source":          "10.0.0.0/8",
		"Destination":     "192.168.0.0/16",
	})
	if err != nil {
		t.Fatal(err)
	}

	resourceName := "hpe_morpheus_security_group_rule.example"
	createChecks := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttrPair(resourceName, "security_group_id", "hpe_morpheus_security_group.test", "id"),
		resource.TestCheckResourceAttr(resourceName, "name", name+"-rule"),
		resource.TestCheckResourceAttr(resourceName, "protocol", "tcp"),
		resource.TestCheckResourceAttr(resourceName, "rule_type", "customRule"),
		resource.TestCheckResourceAttr(resourceName, "direction", "ingress"),
		resource.TestCheckResourceAttr(resourceName, "port_range", "443"),
		resource.TestCheckResourceAttr(resourceName, "source", "0.0.0.0/0"),
		resource.TestCheckResourceAttr(resourceName, "destination", "0.0.0.0/0"),
		resource.TestCheckResourceAttr(resourceName, "policy", "accept"),
	)
	updateChecks := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttrPair(resourceName, "security_group_id", "hpe_morpheus_security_group.test", "id"),
		resource.TestCheckResourceAttr(resourceName, "name", name+"-rule"),
		resource.TestCheckResourceAttr(resourceName, "protocol", "udp"),
		resource.TestCheckResourceAttr(resourceName, "rule_type", "customRule"),
		resource.TestCheckResourceAttr(resourceName, "direction", "ingress"),
		resource.TestCheckResourceAttr(resourceName, "port_range", "53"),
		resource.TestCheckResourceAttr(resourceName, "source", "10.0.0.0/8"),
		resource.TestCheckResourceAttr(resourceName, "destination", "192.168.0.0/16"),
		resource.TestCheckResourceAttr(resourceName, "policy", "accept"),
	)

	checkInPlaceUpdate := resource.ConfigPlanChecks{
		PreApply: []plancheck.PlanCheck{
			plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate),
		},
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + securityGroupConfig + createConfig,
				Check:  createChecks,
			},
			{
				Config:           providerConfig + securityGroupConfig + updateConfig,
				Check:            updateChecks,
				ConfigPlanChecks: checkInPlaceUpdate,
			},
			{
				Config:             providerConfig + securityGroupConfig + updateConfig,
				ExpectNonEmptyPlan: false,
				PlanOnly:           true,
			},
		},
	})
}
