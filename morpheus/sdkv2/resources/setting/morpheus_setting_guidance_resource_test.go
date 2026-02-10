// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package setting_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus"
	sdkv2morpheus "github.com/HPE/terraform-provider-hpe/morpheus/sdkv2"
	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/resources/setting"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/systemoverride"
)

func TestAccMorpheusSettingGuidanceExampleOk(t *testing.T) {
	t.Parallel()

	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	testSystem := systemoverride.GetPreferred(t, "feature")
	providerConfig := testhelpers.ProviderBlock(testSystem)

	resourceConfig, err := setting.RenderSettingGuidanceConfig(t, nil)
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_setting_guidance.tf_example_guidance_setting",
			"power_settings_average_cpu",
			"75",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_setting_guidance.tf_example_guidance_setting",
			"power_settings_maximum_cpu",
			"500",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_setting_guidance.tf_example_guidance_setting",
			"power_settings_network_threshold",
			"2000",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_setting_guidance.tf_example_guidance_setting",
			"cpu_upsize_average_cpu",
			"50",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_setting_guidance.tf_example_guidance_setting",
			"cpu_upsize_maximum_cpu",
			"99",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_setting_guidance.tf_example_guidance_setting",
			"memory_upsize_minimum_free_memory",
			"10",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_setting_guidance.tf_example_guidance_setting",
			"memory_downsize_average_free_memory",
			"60",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_setting_guidance.tf_example_guidance_setting",
			"memory_downsize_maximum_free_memory",
			"30",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), sdkv2morpheus.Provider()),
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
