// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package task_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/testhelpers"
)

// RenderTaskJavascriptConfig generates a terraform configuration string
// for task javascript resource.
// It accepts a name and a map to override default field values.
func RenderTaskJavascriptConfig(
	t *testing.T,
	name string,
	overrides map[string]string,
) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":              name,
		"Code":              name,
		"Labels":            `["demo","terraform"]`,
		"ScriptContent":     `console.log("testing")`,
		"Retryable":         "true",
		"RetryCount":        "1",
		"RetryDelaySeconds": "10",
		"AllowCustomConfig": "true",
	}

	for key, value := range overrides {
		defaults[key] = value
	}

	return testhelpers.RenderExample(
		t,
		"morpheus_task_javascript_resource.tf.tmpl",
		"Name", defaults["Name"],
		"Code", defaults["Code"],
		"Labels", defaults["Labels"],
		"ScriptContent", defaults["ScriptContent"],
		"Retryable", defaults["Retryable"],
		"RetryCount", defaults["RetryCount"],
		"RetryDelaySeconds", defaults["RetryDelaySeconds"],
		"AllowCustomConfig", defaults["AllowCustomConfig"],
	)
}

func TestAccMorpheusTaskJavascriptExampleOk(t *testing.T) {
	t.Parallel()

	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())

	resourceConfig, err := RenderTaskJavascriptConfig(t, name, map[string]string{})
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_javascript.tfexample_javascript",
			"name",
			name,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_javascript.tfexample_javascript",
			"code",
			name,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_javascript.tfexample_javascript",
			"labels.0",
			"demo",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_javascript.tfexample_javascript",
			"labels.1",
			"terraform",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_javascript.tfexample_javascript",
			"script_content",
			`console.log("testing")`,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_javascript.tfexample_javascript",
			"retryable",
			"true",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_javascript.tfexample_javascript",
			"retry_count",
			"1",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_javascript.tfexample_javascript",
			"retry_delay_seconds",
			"10",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_javascript.tfexample_javascript",
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
