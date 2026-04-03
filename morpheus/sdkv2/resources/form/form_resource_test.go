// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package form_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus"
	sdkv2morpheus "github.com/HPE/terraform-provider-hpe/morpheus/sdkv2"
	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/resources/form"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
)

func TestMain(m *testing.M) {
	code := m.Run()

	testhelpers.WriteMergedResults()

	os.Exit(code)
}

func TestAccMorpheusFormExampleOk(t *testing.T) {
	t.Parallel()

	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	// t.Skip("Skipping due to error in resource code")

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())

	resourceConfig, err := form.RenderFormConfig(t, map[string]string{
		"Name": name,
	})
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"code",
			"demo",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"description",
			"demo",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"field_group1_collapsed_by_default",
			"true",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"field_group1_collapsible",
			"true",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"field_group1_description",
			"testin",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"field_group1_name",
			"fg1",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"field_group1_option_type_code",
			"test-input",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"field_group1_option_type_default_value",
			"Demo123",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"field_group1_option_type_description",
			"Terraform text input example",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"field_group1_option_type_display_value_on_details",
			"true",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"field_group1_option_type_exclude_from_search",
			"true",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"field_group1_option_type_export_meta",
			"true",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"field_group1_option_type_field_label",
			"Testin",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"field_group1_option_type_field_name",
			"test",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"field_group1_option_type_help_block",
			"Is this working now",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"field_group1_option_type_hidden",
			"false",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"field_group1_option_type_locked",
			"true",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"field_group1_option_type_name",
			"tf field group 1 text input example",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"field_group1_option_type_placeholder",
			"Testing 123",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"field_group1_option_type_required",
			"true",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"field_group1_option_type_type",
			"text",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"field_group2_collapsed_by_default",
			"true",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"field_group2_collapsible",
			"true",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"field_group2_description",
			"testin",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"field_group2_name",
			"fg2",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"field_group2_option_type_code",
			"test-input",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"field_group2_option_type_default_value",
			"Demo123",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"field_group2_option_type_description",
			"Terraform text input example",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"field_group2_option_type_display_value_on_details",
			"true",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"field_group2_option_type_exclude_from_search",
			"true",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"field_group2_option_type_export_meta",
			"true",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"field_group2_option_type_field_label",
			"Testin",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"field_group2_option_type_field_name",
			"test",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"field_group2_option_type_help_block",
			"Is this working now",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"field_group2_option_type_hidden",
			"false",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"field_group2_option_type_locked",
			"true",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"field_group2_option_type_name",
			"tf field group 2 text input example",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"field_group2_option_type_placeholder",
			"Testing 123",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"field_group2_option_type_required",
			"true",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"field_group2_option_type_type",
			"text",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"labels.#",
			"2",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"name",
			name,
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"option_type1_code",
			"select-input",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"option_type1_default_value",
			"test123",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"option_type1_description",
			"Terraform select example",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"option_type1_display_value_on_details",
			"true",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"option_type1_exclude_from_search",
			"true",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"option_type1_export_meta",
			"true",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"option_type1_field_label",
			"Select Test",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"option_type1_field_name",
			"selectTest",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"option_type1_help_block",
			"Select an option",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"option_type1_hidden",
			"true",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"option_type1_locked",
			"true",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"option_type1_name",
			"tf example select",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"option_type1_option_list_id",
			"1",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"option_type1_placeholder",
			"Testing 123",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"option_type1_required",
			"true",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"option_type1_type",
			"select",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"option_type2_code",
			"radio-input",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"option_type2_default_value",
			"Demo123",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"option_type2_description",
			"Terraform radio example",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"option_type2_display_value_on_details",
			"true",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"option_type2_exclude_from_search",
			"true",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"option_type2_export_meta",
			"true",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"option_type2_field_label",
			"Radio Test",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"option_type2_field_name",
			"radioTest",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"option_type2_help_block",
			"Select an option",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"option_type2_hidden",
			"true",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"option_type2_locked",
			"true",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"option_type2_name",
			"tf radio example",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"option_type2_option_list_id",
			"1",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"option_type2_placeholder",
			"Testing 123",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"option_type2_required",
			"true",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"option_type2_type",
			"radio",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"option_type3_code",
			"test-input",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"option_type3_default_value",
			"Demo123",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"option_type3_description",
			"Terraform text example",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"option_type3_display_value_on_details",
			"true",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"option_type3_exclude_from_search",
			"true",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"option_type3_export_meta",
			"true",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"option_type3_field_label",
			"Testin",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"option_type3_field_name",
			"test",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"option_type3_help_block",
			"Is this working now",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"option_type3_hidden",
			"true",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"option_type3_locked",
			"true",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"option_type3_name",
			"tf text example",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"option_type3_placeholder",
			"Testing 123",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"option_type3_required",
			"true",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"option_type3_type",
			"text",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"option_type4_code",
			"checkbox-input",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"option_type4_default_checked",
			"true",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"option_type4_description",
			"Terraform checkbox example",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"option_type4_display_value_on_details",
			"true",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"option_type4_exclude_from_search",
			"true",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"option_type4_export_meta",
			"true",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"option_type4_field_label",
			"checkbox input",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"option_type4_field_name",
			"checkboxInput",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"option_type4_help_block",
			"Is this working now",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"option_type4_hidden",
			"true",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"option_type4_locked",
			"true",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"option_type4_name",
			"tf checkbox example",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"option_type4_placeholder",
			"Testing 123",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"option_type4_required",
			"true",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"option_type4_type",
			"checkbox",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"option_type5_code",
			"hidden-input",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"option_type5_default_value",
			"test",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"option_type5_description",
			"Terraform hidden input example",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"option_type5_display_value_on_details",
			"true",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"option_type5_exclude_from_search",
			"true",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"option_type5_export_meta",
			"true",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"option_type5_field_label",
			"hidden input",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"option_type5_field_name",
			"hiddenInput",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"option_type5_help_block",
			"Is this working now",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"option_type5_hidden",
			"true",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"option_type5_locked",
			"true",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"option_type5_name",
			"tf hidden input example",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"option_type5_placeholder",
			"Testing 123",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"option_type5_required",
			"true",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"option_type5_type",
			"hidden",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"option_type6_code",
			"number-input",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"option_type6_default_value",
			"4",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"option_type6_description",
			"Terraform number example",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"option_type6_display_value_on_details",
			"true",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"option_type6_exclude_from_search",
			"true",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"option_type6_export_meta",
			"true",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"option_type6_field_label",
			"number input",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"option_type6_field_name",
			"numberInput",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"option_type6_help_block",
			"Is this working now",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"option_type6_hidden",
			"true",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"option_type6_locked",
			"true",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"option_type6_max_value",
			"44",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"option_type6_min_value",
			"3",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"option_type6_name",
			"tf number input example",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"option_type6_placeholder",
			"Testing 123",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"option_type6_required",
			"true",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"option_type6_step",
			"2",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"option_type6_type",
			"number",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"option_type7_code",
			"network-manager-input",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"option_type7_default_value",
			"test123",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"option_type7_description",
			"Terraform network manager example",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"option_type7_display_value_on_details",
			"true",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"option_type7_enable_ip_mode_selection",
			"true",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"option_type7_exclude_from_search",
			"true",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"option_type7_export_meta",
			"true",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"option_type7_field_label",
			"network input",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"option_type7_field_name",
			"networkInput",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"option_type7_help_block",
			"Select a network",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"option_type7_hidden",
			"false",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"option_type7_locked",
			"true",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"option_type7_name",
			"tf network manager example",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"option_type7_placeholder",
			"Select network",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"option_type7_required",
			"true",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"option_type7_show_network_type_selection",
			"true",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"option_type7_type",
			"networkManager",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"option_type8_code",
			"cloud-input",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"option_type8_default_value",
			"test123",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"option_type8_description",
			"Terraform cloud example",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"option_type8_display_value_on_details",
			"true",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"option_type8_exclude_from_search",
			"true",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"option_type8_export_meta",
			"true",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"option_type8_field_label",
			"cloud input",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"option_type8_field_name",
			"cloudInput",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"option_type8_filter_from_resource",
			"true",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"option_type8_help_block",
			"Select a cloud",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"option_type8_hidden",
			"false",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"option_type8_locked",
			"true",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"option_type8_name",
			"tf cloud example",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"option_type8_placeholder",
			"Select cloud",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"option_type8_required",
			"true",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_form.example",
			"option_type8_type",
			"cloud",
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
