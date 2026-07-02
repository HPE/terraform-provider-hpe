// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package optiontype_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	sdkv2morpheus "github.com/HPE/terraform-provider-hpe/morpheus/sdkv2"
	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/resources/optiontype"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

func TestAccMorpheusOptionTypeTypeaheadExampleOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	t.Parallel()

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())

	// Create the option list this option type depends on, rather than relying
	// on a hard-coded reference-environment option list ID.
	dependencyConfig := `
resource "hpe_morpheus_option_list" "dep" {
  name        = "` + name + `-list"
  description = "dependency option list for acceptance test"
  type        = "manual"
  visibility  = "public"
  real_time   = false
}
`

	resourceConfig, err := optiontype.RenderOptionTypeTypeaheadConfig(t, map[string]string{
		"Name":         name,
		"OptionListId": "hpe_morpheus_option_list.dep.id",
	})
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_option_type_typeahead.example",
			"allow_multiple_selections",
			"true",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_option_type_typeahead.example",
			"default_value",
			"testing",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_option_type_typeahead.example",
			"dependent_field",
			"dependent_example",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_option_type_typeahead.example",
			"description",
			"terraform example typeahead option type",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_option_type_typeahead.example",
			"display_value_on_details",
			"true",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_option_type_typeahead.example",
			"editable",
			"true",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_option_type_typeahead.example",
			"export_meta",
			"true",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_option_type_typeahead.example",
			"field_label",
			"numbers",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_option_type_typeahead.example",
			"field_name",
			"example",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_option_type_typeahead.example",
			"help_block",
			"terraform example typeahead",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_option_type_typeahead.example",
			"labels.#",
			"2",
		),
		resource.TestCheckTypeSetElemAttr(
			"hpe_morpheus_option_type_typeahead.example",
			"labels.*",
			"demo",
		),
		resource.TestCheckTypeSetElemAttr(
			"hpe_morpheus_option_type_typeahead.example",
			"labels.*",
			"terraform",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_option_type_typeahead.example",
			"name",
			name,
		),

		resource.TestCheckResourceAttrPair(
			"hpe_morpheus_option_type_typeahead.example",
			"option_list_id",
			"hpe_morpheus_option_list.dep",
			"id",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_option_type_typeahead.example",
			"placeholder",
			"enter text here",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_option_type_typeahead.example",
			"require_field",
			"require_example",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_option_type_typeahead.example",
			"required",
			"true",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_option_type_typeahead.example",
			"show_on_edit",
			"true",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_option_type_typeahead.example",
			"visibility_field",
			"visibility_example",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), sdkv2morpheus.Provider()),
		Steps: []resource.TestStep{
			// Apply
			{
				Config:             providerConfig + dependencyConfig + resourceConfig,
				ExpectNonEmptyPlan: false,
				Check:              checkFn,
			},
			// Plan after apply
			{
				Config:             providerConfig + dependencyConfig + resourceConfig,
				ExpectNonEmptyPlan: false,
				PlanOnly:           true,
			},
		},
	})
}
