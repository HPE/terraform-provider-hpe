// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package task_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/resources/task"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/systemoverride"
)

func TestAccMorpheusTaskShellScriptResourceGitExampleOk(t *testing.T) {
	t.Parallel()

	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	testSystem := systemoverride.GetPreferred(t, "feature")
	providerConfig := testhelpers.ProviderBlock(testSystem)

	name := acctest.RandomWithPrefix(t.Name())

	resourceConfig, err := task.RenderTaskShellScriptGitConfig(t, map[string]string{
		"Name": name,
		"Code": name,
	})
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_shell_script.tfexample_shell_git",
			"name",
			name,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_shell_script.tfexample_shell_git",
			"code",
			name,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_shell_script.tfexample_shell_git",
			"labels.0",
			"demo",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_shell_script.tfexample_shell_git",
			"labels.1",
			"terraform",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_shell_script.tfexample_shell_git",
			"source_type",
			"repository",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_shell_script.tfexample_shell_git",
			"result_type",
			"json",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_shell_script.tfexample_shell_git",
			"script_path",
			"example.sh",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_shell_script.tfexample_shell_git",
			"version_ref",
			"master",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_shell_script.tfexample_shell_git",
			"repository_id",
			"0",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_shell_script.tfexample_shell_git",
			"sudo",
			"true",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_shell_script.tfexample_shell_git",
			"retryable",
			"true",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_shell_script.tfexample_shell_git",
			"retry_count",
			"1",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_shell_script.tfexample_shell_git",
			"retry_delay_seconds",
			"10",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_shell_script.tfexample_shell_git",
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
