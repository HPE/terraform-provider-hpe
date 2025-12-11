// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package automation_test

import (
	"context"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-mux/tf5to6server"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/testhelpers"
	sdkv2morpheus "github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus"
)

func TestMain(m *testing.M) {
	code := m.Run()

	testhelpers.WriteMergedResults()

	os.Exit(code)
}

func newProviderWithError() (tfprotov6.ProviderServer, error) {
	return tf5to6server.UpgradeServer(context.Background(), sdkv2morpheus.Provider().GRPCProvider)
}

var testAccProtoV6ProviderFactories = map[string]func() (
	tfprotov6.ProviderServer, error,
){
	"hpe": newProviderWithError,
}

func RenderMorpheusScaleThresholdConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":                  acctest.RandomWithPrefix(t.Name()),
		"AutoUpscale":           "true",
		"AutoDownscale":         "true",
		"MinCount":              "1",
		"MaxCount":              "3",
		"EnableCpuThreshold":    "true",
		"MinCpuPercentage":      "30.0",
		"MaxCpuPercentage":      "75.0",
		"EnableMemoryThreshold": "true",
		"MinMemoryPercentage":   "20.0",
		"MaxMemoryPercentage":   "60.0",
		"EnableDiskThreshold":   "true",
		"MinDiskPercentage":     "25.0",
		"MaxDiskPercentage":     "80.0",
	}

	for k, v := range overrides {
		defaults[k] = v
	}

	args := []string{}
	for k, v := range defaults {
		args = append(args, k, v)
	}

	return testhelpers.RenderExample(t, "morpheus_scale_threshold_resource.tf.tmpl", args...)
}

func TestAccMorpheusScaleThresholdExampleOk(t *testing.T) {
	t.Parallel()

	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())

	resourceConfig, err := RenderMorpheusScaleThresholdConfig(t, map[string]string{
		"Name": name,
	})
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
