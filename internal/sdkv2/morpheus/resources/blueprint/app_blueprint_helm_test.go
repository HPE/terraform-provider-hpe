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

func TestAccMorpheusAppBlueprintHelmExampleOk(t *testing.T) {
	t.Parallel()

	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())

	resourceConfig, err := blueprint.RenderAppBlueprintHelmConfig(t, map[string]string{
		"Name": name,
	})
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_app_blueprint_helm.tf_example_helm_app_blueprint",
			"name",
			name,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_app_blueprint_helm.tf_example_helm_app_blueprint",
			"description",
			"tf example helm app blueprint",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_app_blueprint_helm.tf_example_helm_app_blueprint",
			"category",
			"helm",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_app_blueprint_helm.tf_example_helm_app_blueprint",
			"integration_id",
			"3",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_app_blueprint_helm.tf_example_helm_app_blueprint",
			"repository_id",
			"1",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_app_blueprint_helm.tf_example_helm_app_blueprint",
			"version_ref",
			"main",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_app_blueprint_helm.tf_example_helm_app_blueprint",
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
