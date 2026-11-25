// (C) Copyright 2025-2026 Hewlett Packard Enterprise Development LP

package integration_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	sdkv2morpheus "github.com/HPE/terraform-provider-hpe/morpheus/sdkv2"
	dsintegration "github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/datasources/integration"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

func TestMain(m *testing.M) {
	code := m.Run()

	testhelpers.WriteMergedResults()

	os.Exit(code)
}

// TestAccMorpheusDataSourceAnsibleTowerInventoryExampleOk looks an inventory
// up by name on the Ansible Tower integration the environment names
// (testhelpers/tower_fixture.go); the integration is resolved by name too, so
// no id is assumed.
func TestAccMorpheusDataSourceAnsibleTowerInventoryExampleOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.AnsibleTower)

	integrationName := testhelpers.AnsibleTowerIntegrationName(t)
	inventoryName := testhelpers.AnsibleTowerInventoryName(t)

	t.Parallel()

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	dependenciesConfig := testhelpers.AnsibleTowerIntegrationBlock(integrationName)

	datasourceConfig, err := dsintegration.RenderAnsibleTowerInventoryConfig(t, map[string]string{
		"Name":                      testhelpers.HCLQuote(inventoryName),
		"AnsibleTowerIntegrationId": testhelpers.AnsibleTowerIntegrationIDRef,
	})
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_ansible_tower_inventory.example",
			"name",
			inventoryName,
		),
		resource.TestCheckResourceAttrSet(
			"data.hpe_morpheus_ansible_tower_inventory.example",
			"id",
		),
		resource.TestCheckResourceAttrPair(
			"data.hpe_morpheus_ansible_tower_inventory.example", "ansible_tower_integration_id",
			"data.hpe_morpheus_integration.ansible_tower", "id",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), sdkv2morpheus.Provider()),
		Steps: []resource.TestStep{
			{
				Config:             providerConfig + dependenciesConfig + datasourceConfig,
				ExpectNonEmptyPlan: false,
				Check:              checkFn,
			},
		},
	})
}
