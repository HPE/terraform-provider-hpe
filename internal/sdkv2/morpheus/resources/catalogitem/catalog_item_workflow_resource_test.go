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

func TestAccMorpheusCatalogItemWorkflowExampleOk(t *testing.T) {
	t.Parallel()

	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Skip("Skipping due to missing configured instrastructure")

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())

	resourceConfig, err := catalogitem.RenderCatalogItemWorkflowConfig(t, map[string]string{
		"Name":    name,
		"Content": "\"Example workflow catalog content\"",
	})
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_catalog_item_workflow.tfexample_workflow_catalog_item",
			"context_type",
			"appliance",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_catalog_item_workflow.tfexample_workflow_catalog_item",
			"content",
			"Example workflow catalog content",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_catalog_item_workflow.tfexample_workflow_catalog_item",
			"dark_logo_image_name",
			"wordpressbak.png",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_catalog_item_workflow.tfexample_workflow_catalog_item",
			"dark_logo_image_path",
			"wordpressbak.png",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_catalog_item_workflow.tfexample_workflow_catalog_item",
			"description",
			"Example Terraform workflow catalog item",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_catalog_item_workflow.tfexample_workflow_catalog_item",
			"enabled",
			"true",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_catalog_item_workflow.tfexample_workflow_catalog_item",
			"featured",
			"true",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_catalog_item_workflow.tfexample_workflow_catalog_item",
			"logo_image_name",
			"wordpress.png",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_catalog_item_workflow.tfexample_workflow_catalog_item",
			"logo_image_path",
			"wordpress.png",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_catalog_item_workflow.tfexample_workflow_catalog_item",
			"visibility",
			"public",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_catalog_item_workflow.tfexample_workflow_catalog_item",
			"workflow_id",
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
