// (C) Copyright 2025-2026 Hewlett Packard Enterprise Development LP

package blueprint_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	sdkv2morpheus "github.com/HPE/terraform-provider-hpe/morpheus/sdkv2"
	dsblueprint "github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/datasources/blueprint"
	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/resources/blueprint"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

func TestAccMorpheusDataSourceNodeTypeExampleOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	t.Parallel()

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())

	var dependenciesConfig string

	// Unique per-run labels -- see MORPH-16400 and the resource example test. Both
	// tests render this resource with t.Parallel(); identical static labels race
	// Morpheus's find-or-create of the shared org-labels on concurrent node type
	// creates, intermittently failing with "labels[N].name must be unique" or a 500.
	labelSuffix := acctest.RandString(8)
	labels := fmt.Sprintf(`["demo-%s", "node-%s", "tf-%s"]`, labelSuffix, labelSuffix, labelSuffix)

	if currentDependency, err := blueprint.RenderNodeTypeConfig(t, map[string]string{
		"Name":   name,
		"Code":   strings.ToLower(name),
		"Labels": labels,
	}); err != nil {
		t.Fatal(err)
	} else {
		dependenciesConfig += currentDependency
	}

	datasourceConfig, err := dsblueprint.RenderNodeTypeConfig(t, map[string]string{
		"Name": "resource.hpe_morpheus_node_type.example.name",
	})
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_node_type.example",
			"name",
			name,
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), sdkv2morpheus.Provider()),
		Steps: []resource.TestStep{
			{
				Config:             providerConfig + dependenciesConfig + datasourceConfig,
				ExpectNonEmptyPlan: false,
				Check:              checkFn,
			},
		},
	})
}
