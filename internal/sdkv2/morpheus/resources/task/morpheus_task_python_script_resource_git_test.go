// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package task_test

import (
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus/resources/task"

)

func TestAccMorpheusTaskPythonScriptGitExampleOk(t *testing.T) {
	t.Parallel()

	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())

	resourceConfig, err := task.RenderTaskPythonScriptGitConfig(t, name, nil)
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
			strings.ToLower(name),
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
			"repository",
		),
		resource.TestCheckResourceAttr(
			resourceName,
			"result_type",
			"json",
		),
		resource.TestCheckResourceAttr(
			resourceName,
			"script_path",
			"example.py",
		),
		resource.TestCheckResourceAttr(
			resourceName,
			"version_ref",
			"master",
		),
		resource.TestCheckResourceAttr(
			resourceName,
			"repository_id",
			"0",
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
