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

func TestAccMorpheusNodeTypeExampleOk(t *testing.T) {
	if capabilities.Missing(t, capabilities.VMware) {
		t.Log("Skipping test due to missing capabilities")

		return
	}
	t.Parallel()

	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	// t.Skip("Skipping due to missing infrastructure in test environment")

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())

	resourceConfig, err := blueprint.RenderNodeTypeConfig(t, map[string]string{
		"Name": name,
	})
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_node_type.example",
			"category",
			"tfexample",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_node_type.example",
			"labels.#",
			"3",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_node_type.example",
			"name",
			name,
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_node_type.example",
			"short_name",
			"tfexamplenodetype",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_node_type.example",
			"technology",
			"vmware",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_node_type.example",
			"version",
			"2.0",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_node_type.example",
			"virtual_image_id",
			"10",
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
