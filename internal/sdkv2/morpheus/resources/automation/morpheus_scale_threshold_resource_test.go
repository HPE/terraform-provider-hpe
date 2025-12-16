// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package automation_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus/resources/automation"

)

func TestAccMorpheusScaleThresholdExampleOk(t *testing.T) {
	t.Parallel()

	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())

	resourceConfig, err := automation.RenderScaleThresholdConfig(
		t,
		name,
		map[string]string{},
	)
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_scale_threshold.tf_example_scale_threshold",
			"name",
			name,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_scale_threshold.tf_example_scale_threshold",
			"auto_upscale",
			"true",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_scale_threshold.tf_example_scale_threshold",
			"auto_downscale",
			"true",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_scale_threshold.tf_example_scale_threshold",
			"min_count",
			"1",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_scale_threshold.tf_example_scale_threshold",
			"max_count",
			"3",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_scale_threshold.tf_example_scale_threshold",
			"enable_cpu_threshold",
			"true",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_scale_threshold.tf_example_scale_threshold",
			"min_cpu_percentage",
			"30",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_scale_threshold.tf_example_scale_threshold",
			"max_cpu_percentage",
			"75",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_scale_threshold.tf_example_scale_threshold",
			"enable_memory_threshold",
			"true",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_scale_threshold.tf_example_scale_threshold",
			"min_memory_percentage",
			"20",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_scale_threshold.tf_example_scale_threshold",
			"max_memory_percentage",
			"60",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_scale_threshold.tf_example_scale_threshold",
			"enable_disk_threshold",
			"true",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_scale_threshold.tf_example_scale_threshold",
			"min_disk_percentage",
			"25",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_scale_threshold.tf_example_scale_threshold",
			"max_disk_percentage",
			"80",
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
