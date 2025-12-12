// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package task_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/testhelpers"
)

func RenderTaskShellScriptUrlConfig(
	t *testing.T,
	name string,
	overrides map[string]string,
) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":              name,
		"Code":              name,
		"Labels":            "[\"demo\", \"terraform\"]",
		"SourceType":        "url",
		"ResultType":        "json",
		"ScriptPath":        "https://example.com/example.sh",
		"Sudo":              "true",
		"Retryable":         "true",
		"RetryCount":        "1",
		"RetryDelaySeconds": "10",
		"AllowCustomConfig": "true",
	}

	for key, value := range overrides {
		defaults[key] = value
	}

	return testhelpers.RenderExample(t,
		"morpheus_task_shell_script_resource_url.tf.tmpl",
		"Name", defaults["Name"],
		"Code", defaults["Code"],
		"Labels", defaults["Labels"],
		"SourceType", defaults["SourceType"],
		"ResultType", defaults["ResultType"],
		"ScriptPath", defaults["ScriptPath"],
		"Sudo", defaults["Sudo"],
		"Retryable", defaults["Retryable"],
		"RetryCount", defaults["RetryCount"],
		"RetryDelaySeconds", defaults["RetryDelaySeconds"],
		"AllowCustomConfig", defaults["AllowCustomConfig"],
	)
}

func TestAccMorpheusTaskShellScriptResourceUrlExampleOk(t *testing.T) {
	t.Parallel()

	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())

	resourceConfig, err := RenderTaskShellScriptUrlConfig(t, name, map[string]string{})
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_shell_script.tfexample_shell_url",
			"name",
			name,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_shell_script.tfexample_shell_url",
			"code",
			name,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_shell_script.tfexample_shell_url",
			"labels.0",
			"demo",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_shell_script.tfexample_shell_url",
			"labels.1",
			"terraform",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_shell_script.tfexample_shell_url",
			"source_type",
			"url",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_shell_script.tfexample_shell_url",
			"result_type",
			"json",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_shell_script.tfexample_shell_url",
			"script_path",
			"https://example.com/example.sh",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_shell_script.tfexample_shell_url",
			"sudo",
			"true",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_shell_script.tfexample_shell_url",
			"retryable",
			"true",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_shell_script.tfexample_shell_url",
			"retry_count",
			"1",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_shell_script.tfexample_shell_url",
			"retry_delay_seconds",
			"10",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_shell_script.tfexample_shell_url",
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
