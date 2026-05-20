// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkfirewallrulegroup_test

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/HPE/terraform-provider-hpe/morpheus"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/networkfirewallrulegroup"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/systemoverride"
)

func TestMain(m *testing.M) {
	systemoverride.ParseFlags()
	code := m.Run()
	testhelpers.WriteMergedResults()
	os.Exit(code)
}

func TestAccMorpheusNetworkFirewallRuleGroupExampleOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	testSystem := systemoverride.GetPreferred(t, "zodiac")
	providerConfig := testhelpers.ProviderBlockForServer(testSystem)

	name := acctest.RandomWithPrefix(t.Name())

	resourceConfig, err := networkfirewallrulegroup.RenderNetworkFirewallRuleGroupConfig(
		t,
		map[string]string{
			"Name": name,
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_network_firewall_rule_group.example",
			"name",
			name,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_network_firewall_rule_group.example",
			"description",
			"An example firewall rule group",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_network_firewall_rule_group.example",
			"priority",
			"100",
		),
		resource.TestCheckResourceAttrSet(
			"hpe_morpheus_network_firewall_rule_group.example",
			"id",
		),
		resource.TestCheckResourceAttrSet(
			"hpe_morpheus_network_firewall_rule_group.example",
			"network_integration_id",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_network_firewall_rule_group.example",
			"external_type",
			"SecurityPolicy",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(
			t, morpheus.New(), nil,
		),
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
					rs, ok := s.RootModule().
						Resources["hpe_morpheus_network_firewall_rule_group.example"]
					if !ok {
						return "", fmt.Errorf("resource not found")
					}

					return rs.Primary.Attributes["network_integration_id"] +
						"." + rs.Primary.Attributes["id"] +
						"." + rs.Primary.Attributes["external_type"], nil
				},
				ResourceName: "hpe_morpheus_network_firewall_rule_group.example",
				Check:        checkFn,
			},
		},
	})
}

func TestAccMorpheusNetworkFirewallRuleGroupUpdateOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	testSystem := systemoverride.GetPreferred(t, "feature")
	providerConfig := testhelpers.ProviderBlockForServer(testSystem)

	name := acctest.RandomWithPrefix(t.Name())

	createConfig, err := networkfirewallrulegroup.RenderNetworkFirewallRuleGroupConfig(
		t,
		map[string]string{
			"Name":        name,
			"Description": "Initial description",
			"Priority":    "100",
			"GroupLayer":  "Application",
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	updatedName := name + "-updated"

	updateConfig, err := networkfirewallrulegroup.RenderNetworkFirewallRuleGroupConfig(
		t,
		map[string]string{
			"Name":        updatedName,
			"Description": "Updated description",
			"Priority":    "200",
			"GroupLayer":  "Application",
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	createChecks := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr(
			"hpe_morpheus_network_firewall_rule_group.example",
			"name",
			name,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_network_firewall_rule_group.example",
			"description",
			"Initial description",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_network_firewall_rule_group.example",
			"priority",
			"100",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_network_firewall_rule_group.example",
			"group_layer",
			"Application",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_network_firewall_rule_group.example",
			"external_type",
			"SecurityPolicy",
		),
	)

	updateChecks := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr(
			"hpe_morpheus_network_firewall_rule_group.example",
			"name",
			updatedName,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_network_firewall_rule_group.example",
			"description",
			"Updated description",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_network_firewall_rule_group.example",
			"priority",
			"200",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_network_firewall_rule_group.example",
			"group_layer",
			"Application",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_network_firewall_rule_group.example",
			"external_type",
			"SecurityPolicy",
		),
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(
			t, morpheus.New(), nil,
		),
		Steps: []resource.TestStep{
			{
				Config:   providerConfig + createConfig,
				Check:    createChecks,
				PlanOnly: false,
			},
			{
				Config:             providerConfig + createConfig,
				Check:              createChecks,
				ExpectNonEmptyPlan: false,
				PlanOnly:           true,
			},
			{
				Config:   providerConfig + updateConfig,
				Check:    updateChecks,
				PlanOnly: false,
			},
			{
				Config:             providerConfig + updateConfig,
				Check:              updateChecks,
				ExpectNonEmptyPlan: false,
				PlanOnly:           true,
			},
		},
	})
}

func TestAccMorpheusNetworkFirewallRuleGroupRequiresReplaceOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	testSystem := systemoverride.GetPreferred(t, "feature")
	providerConfig := testhelpers.ProviderBlockForServer(testSystem)

	name := acctest.RandomWithPrefix(t.Name())

	createConfig, err := networkfirewallrulegroup.RenderNetworkFirewallRuleGroupConfig(
		t,
		map[string]string{
			"Name":       name,
			"GroupLayer": "Application",
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	replaceConfig, err := networkfirewallrulegroup.RenderNetworkFirewallRuleGroupConfig(
		t,
		map[string]string{
			"Name":       name,
			"GroupLayer": "Ethernet",
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(
			t, morpheus.New(), nil,
		),
		Steps: []resource.TestStep{
			{
				Config:   providerConfig + createConfig,
				PlanOnly: false,
			},
			{
				Config:             providerConfig + replaceConfig,
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccMorpheusNetworkFirewallRuleGroupRequiresReplaceExternalTypeOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	testSystem := systemoverride.GetPreferred(t, "feature")
	providerConfig := testhelpers.ProviderBlockForServer(testSystem)

	name := acctest.RandomWithPrefix(t.Name())

	createConfig, err := networkfirewallrulegroup.RenderNetworkFirewallRuleGroupConfig(
		t,
		map[string]string{
			"Name":         name,
			"ExternalType": "SecurityPolicy",
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	replaceConfig, err := networkfirewallrulegroup.RenderNetworkFirewallRuleGroupConfig(
		t,
		map[string]string{
			"Name":         name,
			"ExternalType": "GatewayPolicy",
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(
			t, morpheus.New(), nil,
		),
		Steps: []resource.TestStep{
			{
				Config:   providerConfig + createConfig,
				PlanOnly: false,
			},
			{
				Config:             providerConfig + replaceConfig,
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccMorpheusNetworkFirewallRuleGroupRequiresReplaceNetworkIntegrationIdOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	testSystem := systemoverride.GetPreferred(t, "feature")
	providerConfig := testhelpers.ProviderBlockForServer(testSystem)

	name := acctest.RandomWithPrefix(t.Name())

	createConfig, err := networkfirewallrulegroup.RenderNetworkFirewallRuleGroupConfig(
		t,
		map[string]string{
			"Name": name,
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	replaceConfig, err := networkfirewallrulegroup.RenderNetworkFirewallRuleGroupConfig(
		t,
		map[string]string{
			"Name":                 name,
			"NetworkIntegrationId": "999",
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(
			t, morpheus.New(), nil,
		),
		Steps: []resource.TestStep{
			{
				Config:   providerConfig + createConfig,
				PlanOnly: false,
			},
			{
				Config:             providerConfig + replaceConfig,
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccMorpheusNetworkFirewallRuleGroupImportInvalidFormatErr(t *testing.T) {
	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	testSystem := systemoverride.GetPreferred(t, "feature")
	providerConfig := testhelpers.ProviderBlockForServer(testSystem)

	name := acctest.RandomWithPrefix(t.Name())

	resourceConfig, err := networkfirewallrulegroup.RenderNetworkFirewallRuleGroupConfig(
		t,
		map[string]string{
			"Name": name,
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	config := providerConfig + resourceConfig

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(
			t, morpheus.New(), nil,
		),
		Steps: []resource.TestStep{
			{
				Config: config,
			},
			{
				Config:        config,
				ImportState:   true,
				ImportStateId: "128:1",
				ResourceName:  "hpe_morpheus_network_firewall_rule_group.example",
				ExpectError:   regexp.MustCompile(`expected format`),
			},
		},
	})
}
