// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package template_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/testhelpers"
)

func renderMorpheusSpecTemplateCloudFormationGitConfig(
	t *testing.T,
	name string,
	overrides map[string]string,
) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":                 name,
		"SourceType":           "repository",
		"RepositoryId":         "2",
		"VersionRef":           "main",
		"SpecPath":             "./spec.yaml",
		"CapabilityIam":        "true",
		"CapabilityNamedIam":   "true",
		"CapabilityAutoExpand": "true",
	}

	for key, value := range overrides {
		defaults[key] = value
	}

	args := []string{}
	for key, value := range defaults {
		args = append(args, key, value)
	}

	return testhelpers.RenderExample(
		t,
		"morpheus_spec_template_cloud_formation_resource_git.tf.tmpl",
		args...,
	)
}

func TestAccMorpheusSpecTemplateCloudFormationResourceGitExampleOk(t *testing.T) {
	t.Parallel()

	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())

	resourceConfig, err := renderMorpheusSpecTemplateCloudFormationGitConfig(
		t,
		name,
		map[string]string{},
	)
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_spec_template_cloud_formation.tfexample_cloud_formation_spec_template_git",
			"name",
			name,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_spec_template_cloud_formation.tfexample_cloud_formation_spec_template_git",
			"source_type",
			"repository",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_spec_template_cloud_formation.tfexample_cloud_formation_spec_template_git",
			"repository_id",
			"2",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_spec_template_cloud_formation.tfexample_cloud_formation_spec_template_git",
			"version_ref",
			"main",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_spec_template_cloud_formation.tfexample_cloud_formation_spec_template_git",
			"spec_path",
			"./spec.yaml",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_spec_template_cloud_formation.tfexample_cloud_formation_spec_template_git",
			"capability_iam",
			"true",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_spec_template_cloud_formation.tfexample_cloud_formation_spec_template_git",
			"capability_named_iam",
			"true",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_spec_template_cloud_formation.tfexample_cloud_formation_spec_template_git",
			"capability_auto_expand",
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
