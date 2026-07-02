// (C) Copyright 2025-2026 Hewlett Packard Enterprise Development LP

package network_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	sdkv2morpheus "github.com/HPE/terraform-provider-hpe/morpheus/sdkv2"
	dsnetwork "github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/datasources/network"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

func TestAccMorpheusDataSourceNetworksExampleOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.Network)

	t.Parallel()

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	var dependenciesConfig string

	// Use template defaults (CloudId=3, Name="name", SortAscending=true,
	// Values=["Test*"]) which the checks below assert against. Do not override
	// Name with a random value; the template renders it unquoted and a random
	// string produces an invalid HCL reference.
	datasourceConfig, err := dsnetwork.RenderNetworksConfig(t, map[string]string{})
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
			"sort_ascending",
			"true",
		),

		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_networks.example",
			"filter.#",
			"1",
		),

		resource.TestCheckTypeSetElemNestedAttrs(
			"data.hpe_morpheus_networks.example",
			"filter.*",
			map[string]string{
				"name":     "name",
				"values.0": "Test*",
				"values.#": "1",
			},
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
