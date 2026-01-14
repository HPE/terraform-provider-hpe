// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package template_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus/resources/template"
)

func TestAccMorpheusScriptTemplateExampleOk(t *testing.T) {
	t.Parallel()

	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())

	resourceConfig, err := template.RenderScriptTemplateConfig(t, map[string]string{
		"Name": name,
	})
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_script_template.tfexample_script_template",
			"name",
			name,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_script_template.tfexample_script_template",
			"labels.#",
			"3",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_script_template.tfexample_script_template",
			"labels.0",
			"demo",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_script_template.tfexample_script_template",
			"labels.1",
			"template",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_script_template.tfexample_script_template",
			"labels.2",
			"terraform",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_script_template.tfexample_script_template",
			"script_type",
			"bash",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_script_template.tfexample_script_template",
			"script_phase",
			"provision",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_script_template.tfexample_script_template",
			"script_content",
			"echo \"testing\"",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_script_template.tfexample_script_template",
			"run_as_user",
			"root",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_script_template.tfexample_script_template",
			"sudo",
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
