// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package blueprint_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus"
	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/testhelpers"
	sdkv2morpheus "github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus"
	"github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus/resources/blueprint"
)

func TestAccMorpheusAppBlueprintTerraformExampleOk(t *testing.T) {
	t.Parallel()

	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Skip("Skipping test as it requires additional external services to be configured")

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
			"hpe_morpheus_app_blueprint_terraform.tfapp_blueprint_specs",
			"name",
			name,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_app_blueprint_terraform.tfapp_blueprint_specs",
			"description",
			"testing terraform",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_app_blueprint_terraform.tfapp_blueprint_specs",
			"category",
			"terraformdemo",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_app_blueprint_terraform.tfapp_blueprint_specs",
			"source_type",
			"spec",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_app_blueprint_terraform.tfapp_blueprint_specs",
			"terraform_version",
			"1.1.1",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_app_blueprint_terraform.tfapp_blueprint_specs",
			"terraform_options",
			"-var 'foo=bar'",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_app_blueprint_terraform.tfapp_blueprint_specs",
			"tfvar_secret",
			"tfvars/rdsdemo-secrets",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(
			t,
			morpheus.New(),
			sdkv2morpheus.Provider(),
		),
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
