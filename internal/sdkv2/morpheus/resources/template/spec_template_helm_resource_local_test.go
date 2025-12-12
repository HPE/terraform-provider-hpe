// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package template_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/testhelpers"
)

func RenderSpecTemplateHelmLocalConfig(
	t *testing.T,
	name string,
	overrides map[string]string,
) (string, error) {
	t.Helper()

	defaultSpecContent := "apiVersion: v1\nkind: Service\nmetadata:\nname: {{ template \"fullname\" . }}\n" +
		"labels:\n    chart: \"{{ .Chart.Name }}-{{ .Chart.Version | replace \"+\" \"_\" }}\"\nspec:\n" +
		"type: {{ .Values.service.type }}\nports:\n- port: {{ .Values.service.externalPort }}\n" +
		"    targetPort: {{ .Values.service.internalPort }}\n    protocol: TCP\n" +
		"    name: {{ .Values.service.name }}\nselector:\n    app: {{ template \"fullname\" . }}"

	defaults := map[string]string{
		"Name":        name,
		"SourceType":  "local",
		"SpecContent": defaultSpecContent,
	}

	for key, value := range overrides {
		defaults[key] = value
	}

	return testhelpers.RenderExample(
		t,
		"morpheus_spec_template_helm_resource_local.tf.tmpl",
		"Name", defaults["Name"],
		"SourceType", defaults["SourceType"],
		"SpecContent", defaults["SpecContent"],
	)
}

func TestAccMorpheusSpecTemplateHelmLocalExampleOk(t *testing.T) {
	t.Parallel()

	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())

	resourceConfig, err := RenderSpecTemplateHelmLocalConfig(t, name, map[string]string{})
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_spec_template_helm.tfexample_helm_spec_template_local",
			"name",
			name,
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_spec_template_helm.tfexample_helm_spec_template_local",
			"source_type",
			"local",
		),

		resource.TestCheckResourceAttrSet(
			"hpe_morpheus_spec_template_helm.tfexample_helm_spec_template_local",
			"spec_content",
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
