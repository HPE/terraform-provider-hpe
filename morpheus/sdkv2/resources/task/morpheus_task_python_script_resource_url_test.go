// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package task_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/resources/task"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
)

func TestAccMorpheusTaskPythonScriptUrlExampleOk(t *testing.T) {
	if capabilities.Missing(t, capabilities.All) {
		t.Log("Skipping test due to missing capabilities")

		return
	}
	t.Parallel()

	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())

	resourceConfig, err := task.RenderTaskPythonScriptUrlConfig(t, map[string]string{
		"Name": name,
		"Code": name,
	})
	if err != nil {
		t.Fatal(err)
	}

	resourceName := "hpe_morpheus_task_python_script." + name

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			resourceName,
			"name",
			name,
		),
		resource.TestCheckResourceAttr(
			resourceName,
			"code",
			name,
		),
		resource.TestCheckResourceAttr(
			resourceName,
			"labels.0",
			"demo",
		),
		resource.TestCheckResourceAttr(
			resourceName,
			"labels.1",
			"terraform",
		),
		resource.TestCheckResourceAttr(
			resourceName,
			"source_type",
			"url",
		),
		resource.TestCheckResourceAttr(
			resourceName,
			"result_type",
			"json",
		),
		resource.TestCheckResourceAttr(
			resourceName,
			"script_path",
			"https://example.com/example.py",
		),
		resource.TestCheckResourceAttr(
			resourceName,
			"command_arguments",
			"example",
		),
		resource.TestCheckResourceAttr(
			resourceName,
			"additional_packages",
			"pyyaml",
		),
		resource.TestCheckResourceAttr(
			resourceName,
			"python_binary",
			"/usr/bin/python3",
		),
		resource.TestCheckResourceAttr(
			resourceName,
			"retryable",
			"true",
		),
		resource.TestCheckResourceAttr(
			resourceName,
			"retry_count",
			"1",
		),
		resource.TestCheckResourceAttr(
			resourceName,
			"retry_delay_seconds",
			"10",
		),
		resource.TestCheckResourceAttr(
			resourceName,
			"allow_custom_config",
			"true",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Plan
			{
				Config:             providerConfig + resourceConfig,
				ExpectNonEmptyPlan: true,
				Check:              checkFn,
				PlanOnly:           true,
			},
			// Apply
			{
				Config: providerConfig + resourceConfig,
				Check:  checkFn,
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
