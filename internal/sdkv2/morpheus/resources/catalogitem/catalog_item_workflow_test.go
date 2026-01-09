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
	"github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus/resources/workflow"
)

func TestAccCatalogItemWorkflowExampleOk(t *testing.T) {
	t.Parallel()

	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())

	dependencyConfig, err := workflow.RenderWorkflowOperationalConfig(t, map[string]string{
		"Name": name,
	})

	resourceConfig, err := catalogitem.RenderCatalogItemWorkflowConfig(t, map[string]string{
		"Name":       name,
		"WorkflowId": "hpe_morpheus_workflow_operational.example.id",
	})
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_catalog_item_workflow.example",
			"name",
			name,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_catalog_item_workflow.example",
			"description",
			"Example Terraform workflow catalog item",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_catalog_item_workflow.example",
			"logo_image_path",
			"tfexample.png",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_catalog_item_workflow.example",
			"logo_image_name",
			"tfexample.png",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_catalog_item_workflow.example",
			"dark_logo_image_path",
			"tfexample.png",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_catalog_item_workflow.example",
			"dark_logo_image_name",
			"tfexample.png",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_catalog_item_workflow.example",
			"enabled",
			"true",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_catalog_item_workflow.example",
			"featured",
			"true",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_catalog_item_workflow.example",
			"labels.#",
			"2",
		),
		resource.TestCheckResourceAttrSet(
			"hpe_morpheus_catalog_item_workflow.example",
			"workflow_id",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_catalog_item_workflow.example",
			"context_type",
			"appliance",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_catalog_item_workflow.example",
			"content",
			"Example catalog content",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_catalog_item_workflow.example",
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
