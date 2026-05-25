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

func TestAccMorpheusAppBlueprintArmGitExampleOk(t *testing.T) {
	if capabilities.Missing(t, capabilities.Git) {
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

	resourceConfig, err := blueprint.RenderAppBlueprintArmGitConfig(t, map[string]string{
		"Name": name,
	})
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_app_blueprint_arm.example",
			"category",
			"armtemplates",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_app_blueprint_arm.example",
			"cloud_init_enabled",
			"true",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_app_blueprint_arm.example",
			"description",
			"example arm app blueprint",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_app_blueprint_arm.example",
			"install_agent",
			"true",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_app_blueprint_arm.example",
			"integration_id",
			"3",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_app_blueprint_arm.example",
			"name",
			"example_app_arm_blueprint_git",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_app_blueprint_arm.example",
			"os_type",
			"linux",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_app_blueprint_arm.example",
			"repository_id",
			"1",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_app_blueprint_arm.example",
			"source_type",
			"repository",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_app_blueprint_arm.example",
			"version_ref",
			"main",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_app_blueprint_arm.example",
			"working_path",
			"./test",
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
