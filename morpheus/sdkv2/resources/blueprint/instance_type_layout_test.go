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

func TestAccMorpheusInstanceTypeLayoutExampleOk(t *testing.T) {
	t.Parallel()

	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	testSystem := systemoverride.GetPreferred(t, "feature")
	providerConfig := testhelpers.ProviderBlock(testSystem)

	name := acctest.RandomWithPrefix(t.Name())

	resourceConfig, err := blueprint.RenderInstanceTypeLayoutConfig(t, map[string]string{
		"Name":           name,
		"InstanceTypeId": "data.hpe_morpheus_instance_type.example.id",
	})
	if err != nil {
		t.Fatal(err)
	}

	// use one of the system instance types
	resourceConfig += `
	data "hpe_morpheus_instance_type" "example" {
	  name = "KVM"
	}
	`

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttrSet(
			"hpe_morpheus_instance_type_layout.example",
			"instance_type_id",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_instance_type_layout.example",
			"labels.0",
			"demo",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_instance_type_layout.example",
			"labels.1",
			"layout",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_instance_type_layout.example",
			"labels.2",
			"terraform",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_instance_type_layout.example",
			"name",
			name,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_instance_type_layout.example",
			"technology",
			"vmware",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_instance_type_layout.example",
			"version",
			"1.0",
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
