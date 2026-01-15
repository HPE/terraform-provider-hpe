// (C) Copyright 2025-2026 Hewlett Packard Enterprise Development LP

package network_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus"
	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/testhelpers"
	sdkv2morpheus "github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus"
	dsnetwork "github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus/datasources/network"
)

func TestAccMorpheusDataSourceNetworksExampleOk(t *testing.T) {
	t.Parallel()

	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Skip("Skipping due to missing infrastructure in test environment")

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())

	var dependenciesConfig string

	datasourceConfig, err := dsnetwork.RenderNetworksConfig(t, map[string]string{
		"Name": name,
	})
	if err != nil {
		t.Fatal(err)
	}

	t.Log(providerConfig + dependenciesConfig + datasourceConfig)

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_networks.example",
			"cloud_id",
			"3",
		),

		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_networks.example",
			"name",
			"name",
		),

		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_networks.example",
			"sort_ascending",
			"true",
		),

		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_networks.example",
			"values",
			"[\"Test*\"]",
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
