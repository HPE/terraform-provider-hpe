// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package task_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	sdkv2morpheus "github.com/HPE/terraform-provider-hpe/morpheus/sdkv2"
	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/resources/task"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

func TestAccMorpheusTaskNestedWorkflowExampleOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	t.Parallel()

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	// t.Skip("Skipping due to API error")
	// t.Skip("Skipping due to missing infrastructure in test environment")
	// t.Skip("Skipping due to missing resource implementation")
	// t.Skip("Skipping due to mismatch between Morpheus API and Terraform schema")

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())

	resourceConfig, err := task.RenderTaskNestedWorkflowConfig(t, map[string]string{
		"Name": name,
		// Environment-specific fixture: an operational workflow that exists in
		// the acceptance test environment. Kept out of the example defaults so
		// it doesn't leak into user-facing docs.
		"OperationalWorkflowId":   "797",
		"OperationalWorkflowName": "full-operational-workflow-1",
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
			"797",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_nested_workflow.example",
			"operational_workflow_name",
			"full-operational-workflow-1",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), sdkv2morpheus.Provider()),
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
