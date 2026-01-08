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

func TestAccMorpheusAppBlueprintCloudFormationExampleOk(t *testing.T) {
	t.Parallel()

	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())

	resourceConfig, err := blueprint.RenderAppBlueprintCloudFormationYamlConfig(t, map[string]string{
		"Name": name,
	})
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_app_blueprint_cloud_formation.tf_example_cloud_formation_app_blueprint_yaml",
			"name",
			name,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_app_blueprint_cloud_formation.tf_example_cloud_formation_app_blueprint_yaml",
			"description",
			"Example cloud formation app blueprint",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_app_blueprint_cloud_formation.tf_example_cloud_formation_app_blueprint_yaml",
			"category",
			"cloudformation",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_app_blueprint_cloud_formation.tf_example_cloud_formation_app_blueprint_yaml",
			"install_agent",
			"true",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_app_blueprint_cloud_formation.tf_example_cloud_formation_app_blueprint_yaml",
			"cloud_init_enabled",
			"true",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_app_blueprint_cloud_formation.tf_example_cloud_formation_app_blueprint_yaml",
			"capability_iam",
			"true",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_app_blueprint_cloud_formation.tf_example_cloud_formation_app_blueprint_yaml",
			"capability_named_iam",
			"true",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_app_blueprint_cloud_formation.tf_example_cloud_formation_app_blueprint_yaml",
			"capability_auto_expand",
			"true",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_app_blueprint_cloud_formation.tf_example_cloud_formation_app_blueprint_yaml",
			"source_type",
			"yaml",
		),
		resource.TestCheckResourceAttrSet(
			"hpe_morpheus_app_blueprint_cloud_formation.tf_example_cloud_formation_app_blueprint_yaml",
			"blueprint_content",
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
