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

func TestAccMorpheusAppBlueprintKubernetesSpecExampleOk(t *testing.T) {
	t.Parallel()

	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Skip("Skipping due to missing infrastructure in test environment")

	testSystem := systemoverride.GetPreferred(t, "feature")
	providerConfig := testhelpers.ProviderBlock(testSystem)

	name := acctest.RandomWithPrefix(t.Name())

	resourceConfig, err := blueprint.RenderAppBlueprintKubernetesSpecConfig(t, map[string]string{
		"Name": name,
	})
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_app_blueprint_kubernetes.example",
			"category",
			"k8s",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_app_blueprint_kubernetes.example",
			"description",
			"tf example kubernetes app blueprint",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_app_blueprint_kubernetes.example",
			"name",
			"tf-kubernetes-app-blueprint-example-spec",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_app_blueprint_kubernetes.example",
			"source_type",
			"spec",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_app_blueprint_kubernetes.example",
			"spec_template_ids",
			"[2, 3]",
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
