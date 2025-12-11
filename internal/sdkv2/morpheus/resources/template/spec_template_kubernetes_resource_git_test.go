// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package template_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/testhelpers"
)

// RenderSpecTemplateKubernetesResourceGitConfig renders the configuration for the Git-based
// Kubernetes spec template resource. Pass overrides as a map to customize field values.
func RenderSpecTemplateKubernetesResourceGitConfig(
	t *testing.T,
	overrides map[string]string,
) string {
	t.Helper()

	defaults := map[string]string{
		"Name":         acctest.RandomWithPrefix(t.Name()),
		"SourceType":   "repository",
		"RepositoryId": "2",
		"VersionRef":   "main",
		"SpecPath":     "./spec.yaml",
	}

	for key, value := range overrides {
		defaults[key] = value
	}

	resourceConfig, err := testhelpers.RenderExample(
		t,
		"morpheus_spec_template_kubernetes_resource_git.tf.tmpl",
		"Name", defaults["Name"],
		"SourceType", defaults["SourceType"],
		"RepositoryId", defaults["RepositoryId"],
		"VersionRef", defaults["VersionRef"],
		"SpecPath", defaults["SpecPath"],
	)
	if err != nil {
		t.Fatal(err)
	}

	return resourceConfig
}

func TestAccMorpheusSpecTemplateKubernetesResourceGitExampleOk(t *testing.T) {
	t.Parallel()

	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())

	resourceConfig := RenderSpecTemplateKubernetesResourceGitConfig(t, map[string]string{
		"Name": name,
	})

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_spec_template_kubernetes.tfexample_kubernetes_spec_template_git",
			"name",
			name,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_spec_template_kubernetes.tfexample_kubernetes_spec_template_git",
			"source_type",
			"repository",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_spec_template_kubernetes.tfexample_kubernetes_spec_template_git",
			"repository_id",
			"2",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_spec_template_kubernetes.tfexample_kubernetes_spec_template_git",
			"version_ref",
			"main",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_spec_template_kubernetes.tfexample_kubernetes_spec_template_git",
			"spec_path",
			"./spec.yaml",
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
