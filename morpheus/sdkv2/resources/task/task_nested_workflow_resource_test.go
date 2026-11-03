// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package task_test

import (
	"fmt"
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

	// MORPH-16489: the nestedWorkflow task API validates the ownership of the
	// referenced operational workflow, so a hardcoded master-owned workflow
	// (id 797) fails for a sub-tenant. Create a private operational workflow per
	// run and wire both id and name from it.
	workflowConfig := fmt.Sprintf(`
resource "hpe_morpheus_workflow_operational" "wf" {
  name                = %q
  description         = "MORPH-16489 per-run operational workflow"
  labels              = ["terraform", "acctest"]
  platform            = "all"
  visibility          = "private"
  allow_custom_config = true
}
`, name+"-wf")

	resourceConfig, err := task.RenderTaskNestedWorkflowConfig(t, map[string]string{
		"Name":                    name,
		"OperationalWorkflowId":   "hpe_morpheus_workflow_operational.wf.id",
		"OperationalWorkflowName": "hpe_morpheus_workflow_operational.wf.name",
	})
	if err != nil {
		t.Fatal(err)
	}
	resourceConfig = workflowConfig + resourceConfig

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

		resource.TestCheckResourceAttrPair(
			"hpe_morpheus_task_nested_workflow.example",
			"operational_workflow_id",
			"hpe_morpheus_workflow_operational.wf", "id",
		),

		resource.TestCheckResourceAttrPair(
			"hpe_morpheus_task_nested_workflow.example",
			"operational_workflow_name",
			"hpe_morpheus_workflow_operational.wf", "name",
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
