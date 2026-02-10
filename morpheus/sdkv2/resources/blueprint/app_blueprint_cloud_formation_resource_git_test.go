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
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/systemoverride"
)

func TestAccMorpheusAppBlueprintCloudFormationGitExampleOk(t *testing.T) {
	t.Parallel()

	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Skip("Skipping due to missing infrastructure in test environment")

	testSystem := systemoverride.GetPreferred(t, "feature")
	providerConfig := testhelpers.ProviderBlock(testSystem)

	name := acctest.RandomWithPrefix(t.Name())

	resourceConfig, err := blueprint.RenderAppBlueprintCloudFormationGitConfig(t, map[string]string{
		"Name": name,
	})
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_app_blueprint_cloud_formation.example",
			"capability_auto_expand",
			"true",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_app_blueprint_cloud_formation.example",
			"capability_iam",
			"true",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_app_blueprint_cloud_formation.example",
			"capability_named_iam",
			"true",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_app_blueprint_cloud_formation.example",
			"category",
			"cloudformation",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_app_blueprint_cloud_formation.example",
			"cloud_init_enabled",
			"true",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_app_blueprint_cloud_formation.example",
			"description",
			"Example cloud formation app blueprint",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_app_blueprint_cloud_formation.example",
			"install_agent",
			"true",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_app_blueprint_cloud_formation.example",
			"integration_id",
			"3",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_app_blueprint_cloud_formation.example",
			"name",
			"example_cloud_formation_app_blueprint_git",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_app_blueprint_cloud_formation.example",
			"repository_id",
			"1",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_app_blueprint_cloud_formation.example",
			"source_type",
			"repository",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_app_blueprint_cloud_formation.example",
			"version_ref",
			"main",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_app_blueprint_cloud_formation.example",
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
