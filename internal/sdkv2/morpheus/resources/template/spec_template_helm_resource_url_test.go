// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package template_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/testhelpers"
)

func RenderSpecTemplateHelmUrlConfig(
	t *testing.T,
	overrides map[string]string,
) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":       acctest.RandomWithPrefix(t.Name()),
		"SourceType": "url",
		"SpecPath":   "http://example.com/chart.yaml",
	}

	for key, value := range overrides {
		defaults[key] = value
	}

	return testhelpers.RenderExample(
		t,
		"morpheus_spec_template_helm_resource_url.tf.tmpl",
		"Name", defaults["Name"],
		"SourceType", defaults["SourceType"],
		"SpecPath", defaults["SpecPath"],
	)
}

func TestAccMorpheusSpecTemplateHelmUrlExampleOk(t *testing.T) {
	t.Parallel()

	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())

	resourceConfig, err := RenderSpecTemplateHelmUrlConfig(t, map[string]string{
		"Name": name,
	})
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_spec_template_helm.tfexample_helm_spec_template_url",
			"name",
			name,
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_spec_template_helm.tfexample_helm_spec_template_url",
			"source_type",
			"url",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_spec_template_helm.tfexample_helm_spec_template_url",
			"spec_path",
			"http://example.com/chart.yaml",
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
