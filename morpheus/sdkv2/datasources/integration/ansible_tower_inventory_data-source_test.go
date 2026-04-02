// (C) Copyright 2025-2026 Hewlett Packard Enterprise Development LP

package integration_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus"
	sdkv2morpheus "github.com/HPE/terraform-provider-hpe/morpheus/sdkv2"
	dsintegration "github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/datasources/integration"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
)

func TestMain(m *testing.M) {
	code := m.Run()

	testhelpers.WriteMergedResults()

	os.Exit(code)
}

func TestAccMorpheusDataSourceAnsibleTowerInventoryExampleOk(t *testing.T) {
	t.Parallel()

	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Skip("Skipping due to missing infrastructure in test environment")

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())

	var dependenciesConfig string

	datasourceConfig, err := dsintegration.RenderAnsibleTowerInventoryConfig(t, map[string]string{
		"Name": name,
	})
	if err != nil {
		t.Fatal(err)
	}

	t.Log(providerConfig + dependenciesConfig + datasourceConfig)

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_ansible_tower_inventory.example",
			"name",
			"Demo Inventory",
		),
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_ansible_tower_inventory.example",
			"ansible_tower_integration_id",
			"1",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), sdkv2morpheus.Provider()),
		Steps: []resource.TestStep{
			{
				Config:             providerConfig + dependenciesConfig + datasourceConfig,
				ExpectNonEmptyPlan: false,
				Check:              checkFn,
			},
		},
	})
}
