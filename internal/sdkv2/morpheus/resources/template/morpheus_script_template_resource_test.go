// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package template_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/testhelpers"
)

// RenderScriptTemplateConfig renders the template with provided overrides
func RenderScriptTemplateConfig(
	t *testing.T,
	name string,
	overrides map[string]string,
) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":          name,
		"Labels":        "[\"demo\", \"template\", \"terraform\"]",
		"ScriptType":    "bash",
		"ScriptPhase":   "provision",
		"ScriptContent": "echo \"testing\"",
		"RunAsUser":     "root",
		"Sudo":          "true",
	}

	for key, value := range overrides {
		defaults[key] = value
	}

	return testhelpers.RenderExample(
		t,
		"morpheus_script_template_resource.tf.tmpl",
		"Name", defaults["Name"],
		"Labels", defaults["Labels"],
		"ScriptType", defaults["ScriptType"],
		"ScriptPhase", defaults["ScriptPhase"],
		"ScriptContent", defaults["ScriptContent"],
		"RunAsUser", defaults["RunAsUser"],
		"Sudo", defaults["Sudo"],
	)
}

func TestAccMorpheusScriptTemplateExampleOk(t *testing.T) {
	t.Parallel()

	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())

	resourceConfig, err := RenderScriptTemplateConfig(
		t,
		name,
		map[string]string{},
	)
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
