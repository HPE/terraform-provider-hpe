// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package blueprint_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus/resources/blueprint"
)

func TestAccMorpheusNodeTypeExampleOk(t *testing.T) {
	t.Parallel()

	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	resourceConfig, err := blueprint.RenderNodeTypeConfig(t,
		map[string]string{
			"Name": "tf_example_node_type",
		},
	)
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_node_type.tf_example_node",
			"category",
			"tfexample",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_node_type.tf_example_node",
			"labels.#",
			"3",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_node_type.tf_example_node",
			"name",
			"tf_example_node_type",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_node_type.tf_example_node",
			"service_port.0.name",
			"web",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_node_type.tf_example_node",
			"service_port.0.port",
			"8080",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_node_type.tf_example_node",
			"service_port.0.protocol",
			"HTTP",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_node_type.tf_example_node",
			"service_port.1.name",
			"secureweb",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_node_type.tf_example_node",
			"service_port.1.port",
			"8443",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_node_type.tf_example_node",
			"service_port.1.protocol",
			"HTTPS",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_node_type.tf_example_node",
			"short_name",
			"tfexamplenodetype",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_node_type.tf_example_node",
			"technology",
			"vmware",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_node_type.tf_example_node",
			"version",
			"2.0",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_node_type.tf_example_node",
			"virtual_image_id",
			"10",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Plan
			{
				Config:             providerConfig + resourceConfig,
				ExpectNonEmptyPlan: true,
				Check:              checkFn,
				PlanOnly:           true,
			},
			// Apply
			{
				Config: providerConfig + resourceConfig,
				Check:  checkFn,
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
