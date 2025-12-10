// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package template_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/testhelpers"
)

func RenderSpecTemplateHelmGitConfig(
	t *testing.T,
	overrides map[string]string,
) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":         acctest.RandomWithPrefix(t.Name()),
		"RepositoryId": "2",
		"SourceType":   "repository",
		"SpecPath":     "./spec.yaml",
		"VersionRef":   "main",
	}

	for key, value := range overrides {
		defaults[key] = value
	}

	return testhelpers.RenderExample(
		t,
		"spec_template_helm_resource_git.tf.tmpl",
		"Name", defaults["Name"],
		"RepositoryId", defaults["RepositoryId"],
		"SourceType", defaults["SourceType"],
		"SpecPath", defaults["SpecPath"],
		"VersionRef", defaults["VersionRef"],
	)
}

func TestAccMorpheusSpecTemplateHelmGitExampleOk(t *testing.T) {
	t.Parallel()

	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())

	resourceConfig, err := RenderSpecTemplateHelmGitConfig(t, map[string]string{
		"Name": name,
	})
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_spec_template_helm.tfexample_helm_spec_template_git",
			"name",
			name,
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_spec_template_helm.tfexample_helm_spec_template_git",
			"repository_id",
			"2",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_spec_template_helm.tfexample_helm_spec_template_git",
			"source_type",
			"repository",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_spec_template_helm.tfexample_helm_spec_template_git",
			"spec_path",
			"./spec.yaml",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_spec_template_helm.tfexample_helm_spec_template_git",
			"version_ref",
			"main",
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
