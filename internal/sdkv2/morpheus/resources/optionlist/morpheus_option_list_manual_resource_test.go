// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package optionlist_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/testhelpers"
)

// RenderMorpheusOptionListManualConfig generates a Terraform configuration for the
// morpheus_option_list_manual resource. It accepts an optional map of field overrides to
// customize the default values. Supported override keys: "Name", "Description", "Dataset", "RealTime"
func RenderMorpheusOptionListManualConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":        acctest.RandomWithPrefix(t.Name()),
		"Description": "Terraform manual option list example",
		"Dataset": "[{\"name\": \"Level 1\",\"value\":\"level1\"},\n " +
			"{\"name\": \"Level 2\",\"value\":\"level2\"},\n " +
			"{\"name\": \"Level 3\",\"value\":\"level3\"}\n]",
		"RealTime": "true",
	}

	// Apply overrides
	for key, value := range overrides {
		defaults[key] = value
	}

	//nolint: lll
	return testhelpers.RenderExample(t, "morpheus_option_list_manual_resource.tf.tmpl",
		"Name", defaults["Name"],
		"Description", defaults["Description"],
		"Dataset", defaults["Dataset"],
		"RealTime", defaults["RealTime"],
	)
}

func TestAccMorpheusOptionListManualExampleOk(t *testing.T) {
	t.Parallel()

	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())

	resourceConfig, err := RenderMorpheusOptionListManualConfig(t, map[string]string{"Name": name})
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_option_list_manual.tf_example_manual_option_list",
			"name",
			name,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_option_list_manual.tf_example_manual_option_list",
			"description",
			"Terraform manual option list example",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_option_list_manual.tf_example_manual_option_list",
			"real_time",
			"true",
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
