// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package catalogitem_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus"
	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/testhelpers"
	sdkv2morpheus "github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus"
	"github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus/resources/catalogitem"
)

func TestAccMorpheusCatalogItemAppBlueprintExampleOk(t *testing.T) {
	t.Parallel()

	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Skip("Skipping due to missing configured instrastructure")

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())

	resourceConfig, err := catalogitem.RenderCatalogItemAppBlueprintConfig(t, map[string]string{
		"Name":    name,
		"Content": "\"Example catalog content\"",
		"AppSpec": "\"name: example-app\\ndescription: Example app\\ntype: morpheus\"",
	})
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_catalog_item_app_blueprint.tf_example_app_blueprint_catalog_item",
			"AppSpec",
			"name: example-app\\ndescription: Example app\\ntype: morpheus",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_catalog_item_app_blueprint.tf_example_app_blueprint_catalog_item",
			"BlueprintId",
			"5",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_catalog_item_app_blueprint.tf_example_app_blueprint_catalog_item",
			"Content",
			"Example catalog content",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_catalog_item_app_blueprint.tf_example_app_blueprint_catalog_item",
			"DarkLogoImageName",
			"tfexampledark.png",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_catalog_item_app_blueprint.tf_example_app_blueprint_catalog_item",
			"DarkLogoImagePath",
			"tfexampledark.png",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_catalog_item_app_blueprint.tf_example_app_blueprint_catalog_item",
			"Description",
			"terraform example app blueprint catalog item",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_catalog_item_app_blueprint.tf_example_app_blueprint_catalog_item",
			"Enabled",
			"true",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_catalog_item_app_blueprint.tf_example_app_blueprint_catalog_item",
			"Featured",
			"true",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_catalog_item_app_blueprint.tf_example_app_blueprint_catalog_item",
			"Labels",
			"[\"aws\", \"demo\", \"testing\"]",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_catalog_item_app_blueprint.tf_example_app_blueprint_catalog_item",
			"LogoImageName",
			"tfexample.png",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_catalog_item_app_blueprint.tf_example_app_blueprint_catalog_item",
			"LogoImagePath",
			"tfexample.png",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_catalog_item_app_blueprint.tf_example_app_blueprint_catalog_item",
			"Name",
			"tfexample_app_blueprint_catalog",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_catalog_item_app_blueprint.tf_example_app_blueprint_catalog_item",
			"OptionTypeIds",
			"[2056, 2006, 2058]",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_catalog_item_app_blueprint.tf_example_app_blueprint_catalog_item",
			"Visibility",
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
