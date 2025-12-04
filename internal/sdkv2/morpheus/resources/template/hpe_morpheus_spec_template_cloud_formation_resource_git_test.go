// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package template_test

import (
	"context"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-mux/tf5to6server"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/testhelpers"
	sdkv2morpheus "github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus"
)

func TestMain(m *testing.M) {
	code := m.Run()

	testhelpers.WriteMergedResults()

	os.Exit(code)
}

func newProviderWithError() (tfprotov6.ProviderServer, error) {
	return tf5to6server.UpgradeServer(context.Background(), sdkv2morpheus.Provider().GRPCProvider)
}

var testAccProtoV6ProviderFactories = map[string]func() (
	tfprotov6.ProviderServer, error,
){
	"hpe": newProviderWithError,
}

func TestAccMorpheusSpecTemplateCloudFormationResourceGitExampleOk(t *testing.T) {
	t.Parallel()

	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())

	resourceConfig, err := testhelpers.RenderExample(t, "hpe_morpheus_spec_template_cloud_formation_resource_git.tf.tmpl",
		"Name", name,
		"SourceType", "repository",
		"RepositoryId", "2",
		"VersionRef", "main",
		"SpecPath", "./spec.yaml",
		"CapabilityIam", "true",
		"CapabilityNamedIam", "true",
		"CapabilityAutoExpand", "true",
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
