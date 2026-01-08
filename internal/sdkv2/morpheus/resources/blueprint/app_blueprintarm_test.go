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

func TestAccMorpheusAppBlueprintArmExampleOk(t *testing.T) {
	t.Parallel()

	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Skip("Skipping test as it requires additional external services to be configured")

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())

	resourceConfig, err := blueprint.RenderAppBlueprintArmConfig(t, map[string]string{
		"Name": name,
	})
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_app_blueprint_arm.tf_example_app_arm_blueprint_git",
			"name",
			name,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_app_blueprint_arm.tf_example_app_arm_blueprint_git",
			"description",
			"example arm app blueprint",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_app_blueprint_arm.tf_example_app_arm_blueprint_git",
			"category",
			"armtemplates",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_app_blueprint_arm.tf_example_app_arm_blueprint_git",
			"source_type",
			"repository",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_app_blueprint_arm.tf_example_app_arm_blueprint_git",
			"install_agent",
			"true",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_app_blueprint_arm.tf_example_app_arm_blueprint_git",
			"cloud_init_enabled",
			"true",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_app_blueprint_arm.tf_example_app_arm_blueprint_git",
			"os_type",
			"linux",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_app_blueprint_arm.tf_example_app_arm_blueprint_git",
			"working_path",
			"./test",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_app_blueprint_arm.tf_example_app_arm_blueprint_git",
			"integration_id",
			"3",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_app_blueprint_arm.tf_example_app_arm_blueprint_git",
			"repository_id",
			"1",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_app_blueprint_arm.tf_example_app_arm_blueprint_git",
			"version_ref",
			"main",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), sdkv2morpheus.Provider()),
		Steps: []resource.TestStep{
			{
				Config:             providerConfig + resourceConfig,
				ExpectNonEmptyPlan: false,
				Check:              checkFn,
			},
			{
				Config:             providerConfig + resourceConfig,
				ExpectNonEmptyPlan: false,
				PlanOnly:           true,
			},
		},
	})
}
