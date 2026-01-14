// (C) Copyright 2025-2026 Hewlett Packard Enterprise Development LP

package cloud_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus"
	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/testhelpers"
	sdkv2morpheus "github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus"
	dscloud "github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus/datasources/cloud"
)

func TestAccMorpheusDataSourceCloudByNameExampleOk(t *testing.T) {
	t.Parallel()

	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	var dependenciesConfig string

	if currentDependency, err := dscloud.RenderCloudsConfig(t, map[string]string{}); err != nil {
		t.Fatal(err)
	} else {
		dependenciesConfig += currentDependency
	}

	if currentDependency, err := dscloud.RenderCloudByIdConfig(t, map[string]string{
		"Id": "data.hpe_morpheus_clouds.example.ids[0]",
	}); err != nil {
		t.Fatal(err)
	} else {
		dependenciesConfig += currentDependency
	}

	datasourceConfig, err := dscloud.RenderCloudByNameConfig(t, map[string]string{
		"Name": "data.hpe_morpheus_cloud.example.name",
	})
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttrPair(
			"data.hpe_morpheus_cloud.example_by_name",
			"id",
			"data.hpe_morpheus_cloud.example",
			"id",
		),
	}

	t.Log(providerConfig + dependenciesConfig + datasourceConfig)

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), sdkv2morpheus.Provider()),
		Steps: []resource.TestStep{
			// Apply
			{
				Config:             providerConfig + dependenciesConfig + datasourceConfig,
				ExpectNonEmptyPlan: false,
				Check:              checkFn,
			},
			// Plan after apply
			{
				Config:             providerConfig + dependenciesConfig + datasourceConfig,
				ExpectNonEmptyPlan: false,
				PlanOnly:           true,
			},
		},
	})
}
