// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package optionlist_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/testhelpers"
)

func RenderMorpheusOptionListApiConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":        acctest.RandomWithPrefix(t.Name()),
		"Description": "Terraform Morpheus API option list example",
		"Visibility":  "private",
		"OptionList":  "instances",
		"TranslationScript": "var i=0;\nresults = [];\nfor(i; i<data.length; i++) {" +
			"\n  results.push({name: data[i].name, value: data[i].name});\n}",
	}

	for key, value := range overrides {
		defaults[key] = value
	}

	return testhelpers.RenderExample(
		t,
		"morpheus_option_list_api_resource.tf.tmpl",
		"Name", defaults["Name"],
		"Description", defaults["Description"],
		"Visibility", defaults["Visibility"],
		"OptionList", defaults["OptionList"],
		"TranslationScript", defaults["TranslationScript"],
	)
}

func TestAccMorpheusOptionListApiExampleOk(t *testing.T) {
	t.Parallel()

	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())

	resourceConfig, err := RenderMorpheusOptionListApiConfig(t, map[string]string{
		"Name": name,
	})
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_option_list_api.tf_example_api_option_list",
			"name",
			name,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_option_list_api.tf_example_api_option_list",
			"description",
			"Terraform Morpheus API option list example",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_option_list_api.tf_example_api_option_list",
			"visibility",
			"private",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_option_list_api.tf_example_api_option_list",
			"option_list",
			"instances",
		),
		resource.TestCheckResourceAttrSet(
			"hpe_morpheus_option_list_api.tf_example_api_option_list",
			"translation_script",
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
			// Apply - Note: translation_script formatting causes persistent drift
			{
				Config:             providerConfig + resourceConfig,
				ExpectNonEmptyPlan: true,
				Check:              checkFn,
			},
		},
	})
}
