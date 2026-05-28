// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package template_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/resources/template"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
)

func TestAccMorpheusSpecTemplateTerraformResourceGitExampleOk(t *testing.T) {
	if capabilities.Missing(t, capabilities.AWS) {
		t.Log("Skipping test due to missing capabilities")

		return
	}
	t.Parallel()

	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())

	resourceConfig, err := template.RenderSpecTemplateTerraformGitConfig(t, map[string]string{
		"Name": name,
	})
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_spec_template_terraform.tfexample_terraform_spec_terraform_git",
			"name",
			name,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_spec_template_terraform.tfexample_terraform_spec_terraform_git",
			"source_type",
			"repository",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_spec_template_terraform.tfexample_terraform_spec_terraform_git",
			"repository_id",
			"2",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_spec_template_terraform.tfexample_terraform_spec_terraform_git",
			"version_ref",
			"main",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_spec_template_terraform.tfexample_terraform_spec_terraform_git",
			"spec_path",
			"Instance Types/Terraform/CloudResource/aws/vpc.tf",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:             providerConfig + resourceConfig,
				ExpectNonEmptyPlan: true,
				Check:              checkFn,
				PlanOnly:           true,
			},
			{
				Config: providerConfig + resourceConfig,
				Check:  checkFn,
			},
			{
				Config:             providerConfig + resourceConfig,
				ExpectNonEmptyPlan: false,
				PlanOnly:           true,
			},
		},
	})
}
