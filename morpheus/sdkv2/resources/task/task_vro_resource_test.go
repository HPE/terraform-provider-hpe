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
	"github.com/HPE/terraform-provider-hpe/utils/adapter"
)

func TestAccMorpheusTaskVroExampleOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.VRO)

	t.Parallel()

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())

	expectedBody := "{\n \"parameters\": [\n {\n \"name\": \"vmName\",\n" +
		" \"type\": \"string\",\n \"value\": {\n \"string\": {\n" +
		" \"value\": \"<%=instance.hostname%>\"\n }\n }\n }\n ]\n}"

	resourceConfig, err := task.RenderTaskVroConfig(t, map[string]string{
		"Name": name,
	})
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_vro.example",
			"body",
			expectedBody,
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_vro.example",
			"code",
			"tfexample-vro-task",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_vro.example",
			"execute_target",
			"local",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_vro.example",
			"labels.#",
			"2",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_vro.example",
			"name",
			name,
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_vro.example",
			"retryable",
			"false",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_vro.example",
			"vro_integration_id",
			"1",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_vro.example",
			"vro_workflow_value",
			"1",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewAdaptedMorpheus(), sdkv2morpheus.Provider()),
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
