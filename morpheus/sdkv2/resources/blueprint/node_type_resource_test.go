// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package blueprint_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	sdkv2morpheus "github.com/HPE/terraform-provider-hpe/morpheus/sdkv2"
	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/resources/blueprint"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

func TestAccMorpheusNodeTypeExampleOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.VMware)

	t.Parallel()

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	// t.Skip("Skipping due to missing infrastructure in test environment")

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())

	// Unique per-run labels. This test and the data-source example test both render
	// this resource with t.Parallel(); with the example's static labels, concurrent
	// node type creates race Morpheus's find-or-create of the shared org-labels and
	// intermittently fail with "labels[N].name must be unique" (or a 500). Deriving
	// the labels from a per-run token keeps them distinct across tests and re-runs.
	// See MORPH-16400.
	labelSuffix := acctest.RandString(8)
	labels := fmt.Sprintf(`["demo-%s", "node-%s", "tf-%s"]`, labelSuffix, labelSuffix, labelSuffix)

	resourceConfig, err := blueprint.RenderNodeTypeConfig(t, map[string]string{
		"Name":   name,
		"Labels": labels,
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
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), sdkv2morpheus.Provider()),
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
