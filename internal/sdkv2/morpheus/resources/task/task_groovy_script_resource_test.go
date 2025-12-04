// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package task_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/testhelpers"
)

func TestAccMorpheusTaskGroovyScriptResourceExampleOk(t *testing.T) {
	t.Parallel()

	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())

	resourceConfig, err := testhelpers.RenderExample(t, "task_groovy_script_resource.tf.tmpl",
		"Name", name,
		"Code", "tfexample_groovy_local",
		"Labels", "[\"demo\", \"terraform\"]",
		"SourceType", "local",
		"ScriptContent", "println \"hello\"",
		"Retryable", "true",
		"RetryCount", "1",
		"RetryDelaySeconds", "10",
		"AllowCustomConfig", "true",
	)
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_groovy_script.tfexample_groovy_local",
			"name",
			name,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_groovy_script.tfexample_groovy_local",
			"code",
			"tfexample_groovy_local",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_groovy_script.tfexample_groovy_local",
			"labels.0",
			"demo",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_groovy_script.tfexample_groovy_local",
			"labels.1",
			"terraform",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_groovy_script.tfexample_groovy_local",
			"source_type",
			"local",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_groovy_script.tfexample_groovy_local",
			"script_content",
			"println \"hello\"\n",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_groovy_script.tfexample_groovy_local",
			"retryable",
			"true",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_groovy_script.tfexample_groovy_local",
			"retry_count",
			"1",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_groovy_script.tfexample_groovy_local",
			"retry_delay_seconds",
			"10",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_groovy_script.tfexample_groovy_local",
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
