// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package task_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus"
	sdkv2morpheus "github.com/HPE/terraform-provider-hpe/morpheus/sdkv2"
	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/resources/task"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/systemoverride"
)

func TestAccMorpheusTaskNestedWorkflowExampleOk(t *testing.T) {
	t.Parallel()

	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	// t.Skip("Skipping due to API error")
	// t.Skip("Skipping due to missing infrastructure in test environment")
	// t.Skip("Skipping due to missing resource implementation")
	// t.Skip("Skipping due to mismatch between Morpheus API and Terraform schema")

	testSystem := systemoverride.GetPreferred(t, "feature")
	providerConfig := testhelpers.ProviderBlock(testSystem)

	name := acctest.RandomWithPrefix(t.Name())

	resourceConfig, err := task.RenderTaskNestedWorkflowConfig(t, map[string]string{
		"Name": name,
	})
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_nested_workflow.example",
			"code",
			"tfexample_nested_workflow",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_nested_workflow.example",
			"labels.#",
			"2",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_nested_workflow.example",
			"name",
			name,
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_nested_workflow.example",
			"operational_workflow_id",
			"1",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_nested_workflow.example",
			"operational_workflow_name",
			"Example workflow",
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
