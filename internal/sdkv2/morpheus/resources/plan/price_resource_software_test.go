// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package plan_test

import (
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus"
	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/testhelpers"
	sdkv2morpheus "github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus"
	"github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus/resources/plan"
)

func TestAccMorpheusPriceSoftwareExampleOk(t *testing.T) {
	t.Parallel()

	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())

	resourceConfig, err := plan.RenderPriceSoftwareConfig(t, map[string]string{
		"Name": name,
		"Code": strings.ToLower(name),
	})
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_price.example",
			"code",
			strings.ToLower(name),
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_price.example",
			"cost",
			"38",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_price.example",
			"currency",
			"USD",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_price.example",
			"incur_charges",
			"always",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_price.example",
			"name",
			name,
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_price.example",
			"price_type",
			"software",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_price.example",
			"price_unit",
			"minute",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_price.example",
			"software",
			"visio",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_price.example",
			"tenant_id",
			"1",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), sdkv2morpheus.Provider()),
		Steps: []resource.TestStep{
			// Apply
			{
				Config:             providerConfig + resourceConfig,
				ExpectNonEmptyPlan: false,
				Check:              checkFn,
			},
			// Plan after apply
			{
				Config:             providerConfig + resourceConfig,
				ExpectNonEmptyPlan: false,
				PlanOnly:           true,
			},
		},
	})
}
