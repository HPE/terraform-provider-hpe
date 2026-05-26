// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package optiontype_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus"
	sdkv2morpheus "github.com/HPE/terraform-provider-hpe/morpheus/sdkv2"
	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/resources/optiontype"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
)

func TestMain(m *testing.M) {
	code := m.Run()

	testhelpers.WriteMergedResults()

	os.Exit(code)
}

func TestAccMorpheusOptionTypeTextareaExampleOk(t *testing.T) {
	if capabilities.Missing(t, capabilities.All) {
		t.Log("Skipping test due to missing capabilities")

		return
	}
	t.Parallel()

	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())

	resourceConfig, err := optiontype.RenderOptionTypeTextareaConfig(t, map[string]string{
		"Name": name,
	})
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttrSet(
			"hpe_morpheus_option_type_textarea.tf_example_textarea_option_type",
			"name",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_option_type_textarea.tf_example_textarea_option_type",
			"description",
			"Terraform text area option type example",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_option_type_textarea.tf_example_textarea_option_type",
			"labels.0",
			"demo",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_option_type_textarea.tf_example_textarea_option_type",
			"labels.1",
			"terraform",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_option_type_textarea.tf_example_textarea_option_type",
			"field_name",
			"textareaExample",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_option_type_textarea.tf_example_textarea_option_type",
			"export_meta",
			"true",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_option_type_textarea.tf_example_textarea_option_type",
			"dependent_field",
			"dependent_example",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_option_type_textarea.tf_example_textarea_option_type",
			"visibility_field",
			"visibility_example",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_option_type_textarea.tf_example_textarea_option_type",
			"require_field",
			"require_example",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_option_type_textarea.tf_example_textarea_option_type",
			"show_on_edit",
			"true",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_option_type_textarea.tf_example_textarea_option_type",
			"editable",
			"true",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_option_type_textarea.tf_example_textarea_option_type",
			"display_value_on_details",
			"true",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_option_type_textarea.tf_example_textarea_option_type",
			"field_label",
			"Text Area Example",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_option_type_textarea.tf_example_textarea_option_type",
			"rows",
			"5",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_option_type_textarea.tf_example_textarea_option_type",
			"placeholder",
			"example text",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_option_type_textarea.tf_example_textarea_option_type",
			"default_value",
			"example",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_option_type_textarea.tf_example_textarea_option_type",
			"help_block",
			"Terraform text area option type example",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_option_type_textarea.tf_example_textarea_option_type",
			"required",
			"true",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_option_type_textarea.tf_example_textarea_option_type",
			"verify_pattern",
			`a\D{4}`,
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
