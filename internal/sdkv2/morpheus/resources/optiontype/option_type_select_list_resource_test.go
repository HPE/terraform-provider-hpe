// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package optiontype_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus"
	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/testhelpers"
	sdkv2morpheus "github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus"
	"github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus/resources/optiontype"
)

func TestAccMorpheusOptionTypeSelectListExampleOk(t *testing.T) {
	t.Parallel()

	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Skip("Skipping due to mismatch between Morpheus API and Terraform schema")

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())

	resourceConfig, err := optiontype.RenderOptionTypeSelectListConfig(t, map[string]string{
		"Name": name,
	})
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_option_type_select_list.example",
			"default_value",
			"testing",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_option_type_select_list.example",
			"dependent_field",
			"dependent_example",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_option_type_select_list.example",
			"description",
			"Terraform select list option type example",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_option_type_select_list.example",
			"display_value_on_details",
			"true",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_option_type_select_list.example",
			"editable",
			"true",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_option_type_select_list.example",
			"export_meta",
			"true",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_option_type_select_list.example",
			"field_label",
			"numbers",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_option_type_select_list.example",
			"field_name",
			"tfSelectExample",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_option_type_select_list.example",
			"help_block",
			"fiwefw",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_option_type_select_list.example",
			"labels",
			"[\"demo\", \"terraform\"]",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_option_type_select_list.example",
			"name",
			"tf_example_select_list_option_type",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_option_type_select_list.example",
			"option_list_id",
			"3",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_option_type_select_list.example",
			"require_field",
			"require_example",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_option_type_select_list.example",
			"required",
			"true",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_option_type_select_list.example",
			"show_on_edit",
			"true",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_option_type_select_list.example",
			"visibility_field",
			"visibility_example",
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
