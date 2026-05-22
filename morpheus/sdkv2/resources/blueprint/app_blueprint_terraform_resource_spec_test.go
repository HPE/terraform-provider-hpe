// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package blueprint_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus"
	sdkv2morpheus "github.com/HPE/terraform-provider-hpe/morpheus/sdkv2"
	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/resources/blueprint"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
)

func TestAccMorpheusAppBlueprintTerraformSpecExampleOk(t *testing.T) {
	if capabilities.Missing(t, capabilities.All) {
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

	resourceConfig, err := blueprint.RenderAppBlueprintTerraformSpecConfig(t, map[string]string{
		"Name": name,
	})
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_app_blueprint_terraform.example",
			"category",
			"terraformdemo",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_app_blueprint_terraform.example",
			"description",
			"testing terraform",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_app_blueprint_terraform.example",
			"name",
			name,
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_app_blueprint_terraform.example",
			"source_type",
			"spec",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_app_blueprint_terraform.example",
			"terraform_options",
			"-var foo=bar",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_app_blueprint_terraform.example",
			"terraform_version",
			"1.1.1",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_app_blueprint_terraform.example",
			"tfvar_secret",
			"tfvars/rdsdemo-secrets",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), sdkv2morpheus.Provider()),
		Steps: []resource.TestStep{
			// Apply
			{
				Config:             providerConfig + resourceConfig,
				ExpectNonEmptyPlan: false,
				Check:              checkFn,
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
