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
	overrides map[string]string,
) (string, error) {
	t.Helper()

	defaultSpecContent := `apiVersion: v1
kind: Service
metadata:
name: {{ template "fullname" . }}
labels:
    chart: "{{ .Chart.Name }}-{{ .Chart.Version | replace "+" "_" }}"
spec:
type: {{ .Values.service.type }}
ports:
- port: {{ .Values.service.externalPort }}
    targetPort: {{ .Values.service.internalPort }}
    protocol: TCP
    name: {{ .Values.service.name }}
selector:
    app: {{ template "fullname" . }}`

	defaults := map[string]string{
		"Name":        acctest.RandomWithPrefix(t.Name()),
		"SourceType":  "local",
		"SpecContent": defaultSpecContent,
	}

	for key, value := range overrides {
		defaults[key] = value
	}

	return testhelpers.RenderExample(
		t,
		"spec_template_helm_resource_local.tf.tmpl",
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

	specContent := `apiVersion: v1
kind: Service
metadata:
name: {{ template "fullname" . }}
labels:
    chart: "{{ .Chart.Name }}-{{ .Chart.Version | replace "+" "_" }}"
spec:
type: {{ .Values.service.type }}
ports:
- port: {{ .Values.service.externalPort }}
    targetPort: {{ .Values.service.internalPort }}
    protocol: TCP
    name: {{ .Values.service.name }}
selector:
    app: {{ template "fullname" . }}`

	resourceConfig, err := RenderSpecTemplateHelmLocalConfig(t, map[string]string{
		"Name":        name,
		"SpecContent": specContent,
	})
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

		resource.TestCheckResourceAttr(
			"hpe_morpheus_spec_template_helm.tfexample_helm_spec_template_local",
			"spec_content",
			specContent,
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
