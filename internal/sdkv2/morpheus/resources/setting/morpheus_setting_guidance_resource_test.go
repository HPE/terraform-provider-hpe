// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package setting_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/testhelpers"
)

func RenderMorpheusSettingGuidanceConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"PowerSettingsAverageCpu":         "75",
		"PowerSettingsMaximumCpu":         "500",
		"PowerSettingsNetworkThreshold":   "2000",
		"CpuUpsizeAverageCpu":             "50",
		"CpuUpsizeMaximumCpu":             "99",
		"MemoryUpsizeMinimumFreeMemory":   "10",
		"MemoryDownsizeAverageFreeMemory": "60",
		"MemoryDownsizeMaximumFreeMemory": "30",
	}

	for key, value := range overrides {
		defaults[key] = value
	}

	return testhelpers.RenderExample(t, "morpheus_setting_guidance_resource.tf.tmpl",
		"PowerSettingsAverageCpu", defaults["PowerSettingsAverageCpu"],
		"PowerSettingsMaximumCpu", defaults["PowerSettingsMaximumCpu"],
		"PowerSettingsNetworkThreshold", defaults["PowerSettingsNetworkThreshold"],
		"CpuUpsizeAverageCpu", defaults["CpuUpsizeAverageCpu"],
		"CpuUpsizeMaximumCpu", defaults["CpuUpsizeMaximumCpu"],
		"MemoryUpsizeMinimumFreeMemory", defaults["MemoryUpsizeMinimumFreeMemory"],
		"MemoryDownsizeAverageFreeMemory", defaults["MemoryDownsizeAverageFreeMemory"],
		"MemoryDownsizeMaximumFreeMemory", defaults["MemoryDownsizeMaximumFreeMemory"],
	)
}

func TestAccMorpheusSettingGuidanceExampleOk(t *testing.T) {
	t.Parallel()

	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())
	_ = name // name is reserved for future use

	resourceConfig, err := RenderMorpheusSettingGuidanceConfig(t, nil)
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
