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

func TestAccMorpheusTaskPowershellScriptResourceGitExampleOk(t *testing.T) {
	t.Parallel()

	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())

	resourceConfig, err := task.RenderTaskPowershellScriptGitConfig(t, name, nil)
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_powershell_script."+name,
			"name",
			name,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_powershell_script."+name,
			"code",
			strings.ToLower(name),
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_powershell_script."+name,
			"source_type",
			"repository",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_powershell_script."+name,
			"result_type",
			"json",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_powershell_script."+name,
			"script_path",
			"example.ps",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_powershell_script."+name,
			"version_ref",
			"master",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_powershell_script."+name,
			"repository_id",
			"0",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_powershell_script."+name,
			"elevated_shell",
			"true",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_powershell_script."+name,
			"retryable",
			"true",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_powershell_script."+name,
			"retry_count",
			"1",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_powershell_script."+name,
			"retry_delay_seconds",
			"10",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_powershell_script."+name,
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
