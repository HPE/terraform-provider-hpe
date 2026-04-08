// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package form_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus"
	sdkv2morpheus "github.com/HPE/terraform-provider-hpe/morpheus/sdkv2"
	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/resources/form"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
)

func TestAccMorpheusFormFieldGroupsOk(t *testing.T) {
	t.Parallel()

	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()
	name := acctest.RandomWithPrefix(t.Name())

	resourceConfig, err := form.RenderFieldGroupsConfig(t, map[string]string{"Name": name})
	if err != nil {
		t.Fatal(err)
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), sdkv2morpheus.Provider()),
		Steps: []resource.TestStep{
			{
				Config:             providerConfig + resourceConfig,
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "code", "demo"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "description", "demo"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "field_group1_collapsed_by_default", "true"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "field_group1_collapsible", "true"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "field_group1_description", "testin"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "field_group1_name", "fg1"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "field_group1_option_type_code", "test-input"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "field_group1_option_type_default_value", "Demo123"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "field_group1_option_type_description", "Terraform text input example"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "field_group1_option_type_display_value_on_details", "true"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "field_group1_option_type_exclude_from_search", "true"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "field_group1_option_type_export_meta", "true"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "field_group1_option_type_field_label", "Testin"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "field_group1_option_type_field_name", "test"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "field_group1_option_type_help_block", "Is this working now"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "field_group1_option_type_hidden", "false"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "field_group1_option_type_locked", "true"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "field_group1_option_type_name", "tf field group 1 text input example"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "field_group1_option_type_placeholder", "Testing 123"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "field_group1_option_type_required", "true"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "field_group1_option_type_type", "text"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "field_group2_collapsed_by_default", "true"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "field_group2_collapsible", "true"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "field_group2_description", "testin"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "field_group2_name", "fg2"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "field_group2_option_type_code", "test-input"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "field_group2_option_type_default_value", "Demo123"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "field_group2_option_type_description", "Terraform text input example"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "field_group2_option_type_display_value_on_details", "true"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "field_group2_option_type_exclude_from_search", "true"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "field_group2_option_type_export_meta", "true"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "field_group2_option_type_field_label", "Testin"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "field_group2_option_type_field_name", "test"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "field_group2_option_type_help_block", "Is this working now"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "field_group2_option_type_hidden", "false"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "field_group2_option_type_locked", "true"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "field_group2_option_type_name", "tf field group 2 text input example"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "field_group2_option_type_placeholder", "Testing 123"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "field_group2_option_type_required", "true"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "field_group2_option_type_type", "text"),
				),
			},
			{
				Config:             providerConfig + resourceConfig,
				ExpectNonEmptyPlan: false,
				PlanOnly:           true,
			},
		},
	})
}
