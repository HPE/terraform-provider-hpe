// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package task_test

import (
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/testhelpers"
)

func RenderMorpheusTaskPythonScriptConfig(
	t *testing.T, name string,
	overrides map[string]string,
) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":               name,
		"Code":               strings.ToLower(name),
		"Labels":             "[\"demo\", \"terraform\"]",
		"SourceType":         "local",
		"ScriptContent":      "print('morpheus')\\nprint('python')",
		"CommandArguments":   "example",
		"AdditionalPackages": "pyyaml",
		"PythonBinary":       "/usr/bin/python3",
		"Retryable":          "true",
		"RetryCount":         "1",
		"RetryDelaySeconds":  "10",
		"AllowCustomConfig":  "true",
	}

	for key, value := range overrides {
		defaults[key] = value
	}

	return testhelpers.RenderExample(
		t,
		"morpheus_task_python_script_resource.tf.tmpl",
		"Name", defaults["Name"],
		"Code", defaults["Code"],
		"Labels", defaults["Labels"],
		"SourceType", defaults["SourceType"],
		"ScriptContent", defaults["ScriptContent"],
		"CommandArguments", defaults["CommandArguments"],
		"AdditionalPackages", defaults["AdditionalPackages"],
		"PythonBinary", defaults["PythonBinary"],
		"Retryable", defaults["Retryable"],
		"RetryCount", defaults["RetryCount"],
		"RetryDelaySeconds", defaults["RetryDelaySeconds"],
		"AllowCustomConfig", defaults["AllowCustomConfig"],
	)
}

func TestAccMorpheusTaskPythonScriptExampleOk(t *testing.T) {
	t.Parallel()

	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())

	resourceConfig, err := RenderMorpheusTaskPythonScriptConfig(t, name, nil)
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
			"local",
		),
		resource.TestCheckResourceAttr(
			resourceName,
			"script_content",
			"print('morpheus')\\nprint('python')\\n",
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
			// Skipping for now
			// // Apply
			// {
			// 	Config: providerConfig + resourceConfig,
			// 	Check:  checkFn,
			// },
			// // Plan after apply
			// {
			// 	Config:             providerConfig + resourceConfig,
			// 	ExpectNonEmptyPlan: false,
			// 	PlanOnly:           true,
			// },
		},
	})
}
