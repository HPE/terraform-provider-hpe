// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package template_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/resources/template"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/systemoverride"
)

func TestAccMorpheusSpecTemplateTerraformResourceLocalExampleOk(t *testing.T) {
	t.Parallel()

	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	testSystem := systemoverride.GetPreferred(t, "feature")
	providerConfig := testhelpers.ProviderBlock(testSystem)

	name := acctest.RandomWithPrefix(t.Name())

	specContent := `resource "aws_instance" "instance_1" {
  ami           = "ami-0b91a410940e82c54"
  instance_type = "t2.micro"
}
`

	resourceConfig, err := template.RenderSpecTemplateTerraformLocalConfig(t, map[string]string{
		"Name":        name,
		"SpecContent": specContent,
	})
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_spec_template_terraform.tfexample_terraform_spec_terraform_local",
			"name",
			name,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_spec_template_terraform.tfexample_terraform_spec_terraform_local",
			"source_type",
			"local",
		),
		// TODO: Get the DiffSuppressFunc working
		// resource.TestCheckResourceAttr(
		// 	"hpe_morpheus_spec_template_terraform.tfexample_terraform_spec_terraform_local",
		// 	"spec_content",
		// 	specContent,
		// ),
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
