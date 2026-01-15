// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package plan_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus"
	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/testhelpers"
	sdkv2morpheus "github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus"
	"github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus/resources/plan"
)

func TestAccMorpheusPriceSetExampleOk(t *testing.T) {
	t.Parallel()

	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())

	dependencyConfig, err := plan.RenderPriceConfig(t, map[string]string{
		"Name": name,
		"Code": name,
	})
	if err != nil {
		t.Fatal(err)
	}

	resourceConfig, err := plan.RenderPriceSetConfig(t, map[string]string{
		"Name":     name,
		"Code":     name,
		"PriceIds": "[hpe_morpheus_price.example.id]",
	})
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_price_set.example",
			"name",
			name,
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_price_set.example",
			"code",
			name,
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_price_set.example",
			"region_code",
			"us-west-2",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_price_set.example",
			"price_unit",
			"minute",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_price_set.example",
			"type",
			"fixed",
		),

		resource.TestCheckResourceAttrSet(
			"hpe_morpheus_price_set.example",
			"price_ids.0",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), sdkv2morpheus.Provider()),
		Steps: []resource.TestStep{
			// Apply
			{
				Config:             providerConfig + dependencyConfig + resourceConfig,
				ExpectNonEmptyPlan: false,
				Check:              checkFn,
			},
			// Plan after apply
			{
				Config:             providerConfig + dependencyConfig + resourceConfig,
				ExpectNonEmptyPlan: false,
				PlanOnly:           true,
			},
		},
	})
}
