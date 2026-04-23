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
)

func TestAccMorpheusAppBlueprintKubernetesYamlExampleOk(t *testing.T) {
	t.Parallel()

	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	// t.Skip("Skipping due to missing infrastructure in test environment")

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())

	resourceConfig, err := blueprint.RenderAppBlueprintKubernetesYamlConfig(t, map[string]string{
		"Name": name,
	})
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_app_blueprint_kubernetes.tfexample_kubernetes_app_blueprint_yaml",
			"blueprint_content",
			"---\napiVersion: apps/v1\nkind: Deployment\nmetadata:\n name: nginx-deployment\n"+
				" labels:\n app: nginx\nspec:\n replicas: 3\n selector:\n matchLabels:\n app: nginx\n"+
				" template:\n metadata:\n labels:\n app: nginx\n spec:\n containers:\n - name: nginx\n"+
				" image: nginx:1.14.2\n ports:\n - containerPort: 80",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_app_blueprint_kubernetes.tfexample_kubernetes_app_blueprint_yaml",
			"category",
			"k8s",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_app_blueprint_kubernetes.tfexample_kubernetes_app_blueprint_yaml",
			"description",
			"tf example kubernetes app blueprint",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_app_blueprint_kubernetes.tfexample_kubernetes_app_blueprint_yaml",
			"name",
			name,
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_app_blueprint_kubernetes.tfexample_kubernetes_app_blueprint_yaml",
			"source_type",
			"yaml",
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
