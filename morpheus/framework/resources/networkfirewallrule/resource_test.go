// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkfirewallrule_test

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/HPE/terraform-provider-hpe/morpheus"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/networkfirewallrule"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/networkfirewallrulegroup"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/clientfactory"
)

func TestMain(m *testing.M) {
	code := testhelpers.TestMain(m)
	testhelpers.WriteMergedResults()
	os.Exit(code)
}

// firewallRuleConfigWithGroup renders a self-contained firewall rule group plus a
// firewall rule that references it via rule_group_id. The Morpheus API resolves
// the rule's security group from rule.ruleGroup.id (see NetworkFirewallRulesController
// .save -> ruleParams.securityGroupId), so a real, freshly created group is required;
// hard-coding a group ID fails with "no security group specified" on environments
// where that ID does not exist.
func firewallRuleConfigWithGroup(
	t *testing.T,
	groupName string,
	ruleOverrides map[string]string,
) string {
	t.Helper()

	groupConfig, err := networkfirewallrulegroup.RenderNetworkFirewallRuleGroupConfig(
		t,
		map[string]string{
			"Name": groupName,
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	if ruleOverrides == nil {
		ruleOverrides = map[string]string{}
	}
	ruleOverrides["RuleGroupId"] = "hpe_morpheus_network_firewall_rule_group.example.id"

	ruleConfig, err := networkfirewallrule.RenderNetworkFirewallRuleConfig(t, ruleOverrides)
	if err != nil {
		t.Fatal(err)
	}

	return groupConfig + ruleConfig
}

func checkDestroy(t *testing.T) resource.TestCheckFunc {
	t.Helper()

	return func(s *terraform.State) error {
		ctx := context.Background()
		client := clientfactory.NewAPIClient(
			ctx,
			os.Getenv("TF_VAR_testacc_morpheus_url"),
			os.Getenv("TF_VAR_testacc_morpheus_username"),
			os.Getenv("TF_VAR_testacc_morpheus_password"),
			os.Getenv("TF_VAR_testacc_morpheus_tenant_subdomain"),
			os.Getenv("TF_VAR_testacc_morpheus_access_token"),
			clientfactory.WithInsecureTLS(),
		)

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "hpe_morpheus_network_firewall_rule" {
				continue
			}

			id, err := strconv.ParseInt(rs.Primary.ID, 10, 64)
			if err != nil {
				return fmt.Errorf("invalid ID %q: %w", rs.Primary.ID, err)
			}

			serverId, err := strconv.ParseInt(rs.Primary.Attributes["network_integration_id"], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid network_integration_id %q: %w", rs.Primary.Attributes["network_integration_id"], err)
			}

			_, httpResp, _ := client.NetworksAPI.
				GetNetworkFirewallRule(ctx, id, serverId).Execute()
			if httpResp != nil && httpResp.StatusCode != http.StatusNotFound {
				return fmt.Errorf("network firewall rule %d still exists", id)
			}
		}

		return nil
	}
}

func TestAccMorpheusNetworkFirewallRuleResourceExampleOk(t *testing.T) {
	if capabilities.Missing(t, capabilities.NetworkFirewall) {
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

	resourceConfig := firewallRuleConfigWithGroup(
		t,
		name+"-group",
		map[string]string{
			"Name":        name,
			"Description": "test description",
			"Priority":    "10",
		},
	)

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_network_firewall_rule.example",
			"name",
			name,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_network_firewall_rule.example",
			"direction",
			"Ingress",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_network_firewall_rule.example",
			"policy",
			"Accept",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_network_firewall_rule.example",
			"enabled",
			"true",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_network_firewall_rule.example",
			"description",
			"test description",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_network_firewall_rule.example",
			"priority",
			"10",
		),
		resource.TestCheckResourceAttrSet(
			"hpe_morpheus_network_firewall_rule.example",
			"id",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_network_firewall_rule.example",
			"network_integration_id",
			"5",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		CheckDestroy:             checkDestroy(t),
		Steps: []resource.TestStep{
			{
				Config:             providerConfig + resourceConfig,
				ExpectNonEmptyPlan: false,
				Check:              checkFn,
				PlanOnly:           false,
			},
			{
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources["hpe_morpheus_network_firewall_rule.example"]
					if !ok {
						return "", fmt.Errorf("resource not found")
					}

					return rs.Primary.Attributes["network_integration_id"] + "." + rs.Primary.Attributes["id"], nil
				},
				ResourceName: "hpe_morpheus_network_firewall_rule.example",
			},
		},
	})
}

func TestAccMorpheusNetworkFirewallRuleResourceUpdateOk(t *testing.T) {
	if capabilities.Missing(t, capabilities.NetworkFirewall) {
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

	createConfig := firewallRuleConfigWithGroup(
		t,
		name+"-group",
		map[string]string{
			"Name":        name,
			"Description": "initial description",
			"Priority":    "5",
		},
	)

	updatedName := name + "-updated"

	updateConfig := firewallRuleConfigWithGroup(
		t,
		name+"-group",
		map[string]string{
			"Name":        updatedName,
			"Direction":   "Egress",
			"Policy":      "Deny",
			"Enabled":     "false",
			"Description": "updated description",
			"Priority":    "20",
		},
	)

	baseChecks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_network_firewall_rule.example",
			"name",
			name,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_network_firewall_rule.example",
			"direction",
			"Ingress",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_network_firewall_rule.example",
			"policy",
			"Accept",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_network_firewall_rule.example",
			"enabled",
			"true",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_network_firewall_rule.example",
			"description",
			"initial description",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_network_firewall_rule.example",
			"priority",
			"5",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(baseChecks...)

	updateChecks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_network_firewall_rule.example",
			"name",
			updatedName,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_network_firewall_rule.example",
			"direction",
			"Egress",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_network_firewall_rule.example",
			"policy",
			"Deny",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_network_firewall_rule.example",
			"enabled",
			"false",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_network_firewall_rule.example",
			"description",
			"updated description",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_network_firewall_rule.example",
			"priority",
			"20",
		),
	}

	checkUpdateFn := resource.ComposeAggregateTestCheckFunc(updateChecks...)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		CheckDestroy:             checkDestroy(t),
		Steps: []resource.TestStep{
			{
				Config:   providerConfig + createConfig,
				Check:    checkFn,
				PlanOnly: false,
			},
			{
				Config:             providerConfig + createConfig,
				Check:              checkFn,
				ExpectNonEmptyPlan: false,
				PlanOnly:           true,
			},
			{
				Config:   providerConfig + updateConfig,
				Check:    checkUpdateFn,
				PlanOnly: false,
			},
			{
				Config:             providerConfig + updateConfig,
				Check:              checkUpdateFn,
				ExpectNonEmptyPlan: false,
				PlanOnly:           true,
			},
		},
	})
}

func TestAccMorpheusNetworkFirewallRuleResourceNestedAttributesOk(t *testing.T) {
	if capabilities.Missing(t, capabilities.NetworkFirewall) {
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

	resourceConfig := firewallRuleConfigWithGroup(
		t,
		name+"-group",
		map[string]string{
			"Name":           name,
			"DestinationIds": `"ANY"`,
		},
	)

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_network_firewall_rule.example",
			"name",
			name,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_network_firewall_rule.example",
			"destinations.id.#",
			"1",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		CheckDestroy:             checkDestroy(t),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + resourceConfig,
				Check:  checkFn,
			},
		},
	})
}

func TestAccMorpheusNetworkFirewallRuleResourceImportBadIDError(t *testing.T) {
	if capabilities.Missing(t, capabilities.NetworkFirewall) {
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

	resourceConfig := firewallRuleConfigWithGroup(
		t,
		name+"-group",
		map[string]string{
			"Name": name,
		},
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		CheckDestroy:             checkDestroy(t),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + resourceConfig,
			},
			{
				ImportState:   true,
				ImportStateId: "not-a-valid-id",
				ResourceName:  "hpe_morpheus_network_firewall_rule.example",
				ExpectError:   regexp.MustCompile(`expected format`),
			},
		},
	})
}

func TestAccMorpheusNetworkFirewallRuleResourceImportNonNumericIDError(t *testing.T) {
	if capabilities.Missing(t, capabilities.NetworkFirewall) {
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

	resourceConfig := firewallRuleConfigWithGroup(
		t,
		name+"-group",
		map[string]string{
			"Name": name,
		},
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		CheckDestroy:             checkDestroy(t),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + resourceConfig,
			},
			{
				ImportState:   true,
				ImportStateId: "abc.def",
				ResourceName:  "hpe_morpheus_network_firewall_rule.example",
				ExpectError:   regexp.MustCompile(`non-number`),
			},
		},
	})
}
