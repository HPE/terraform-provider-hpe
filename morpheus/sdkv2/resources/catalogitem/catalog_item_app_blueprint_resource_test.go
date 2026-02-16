// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package catalogitem_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus"
	sdkv2morpheus "github.com/HPE/terraform-provider-hpe/morpheus/sdkv2"
	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/resources/catalogitem"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
)

func TestAccMorpheusCatalogItemAppBlueprintExampleOk(t *testing.T) {
	t.Parallel()

	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Skip("Skipping due to missing infrastructure in test environment")

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())

	resourceConfig, err := catalogitem.RenderCatalogItemAppBlueprintConfig(t, map[string]string{
		"Name": name,
	})
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_catalog_item_app_blueprint.example",
			"app_spec",
			"file(\"${path.module}/appSpec.yaml\")",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_catalog_item_app_blueprint.example",
			"blueprint_id",
			"5",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_catalog_item_app_blueprint.example",
			"content",
			"file(\"${path.module}/catalog-data.md\")",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_catalog_item_app_blueprint.example",
			"dark_logo_image_name",
			"tfexampledark.png",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_catalog_item_app_blueprint.example",
			"dark_logo_image_path",
			"tfexampledark.png",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_catalog_item_app_blueprint.example",
			"description",
			"terraform example app blueprint catalog item",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_catalog_item_app_blueprint.example",
			"enabled",
			"true",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_catalog_item_app_blueprint.example",
			"featured",
			"true",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_catalog_item_app_blueprint.example",
			"labels",
			"[\"aws\", \"demo\", \"testing\"]",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_catalog_item_app_blueprint.example",
			"logo_image_name",
			"tfexample.png",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_catalog_item_app_blueprint.example",
			"logo_image_path",
			"tfexample.png",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_catalog_item_app_blueprint.example",
			"name",
			"tfexample_app_blueprint_catalog",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_catalog_item_app_blueprint.example",
			"option_type_ids",
			"[2056, 2006, 2058]",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_catalog_item_app_blueprint.example",
			"visibility",
			"public",
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
