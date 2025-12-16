// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package optiontype_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-mux/tf5to6server"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/testhelpers"
	sdkv2morpheus "github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus"
	"github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus/resources/optiontype"

)

func newProviderWithError() (tfprotov6.ProviderServer, error) {
	return tf5to6server.UpgradeServer(context.Background(), sdkv2morpheus.Provider().GRPCProvider)
}

var testAccProtoV6ProviderFactories = map[string]func() (
	tfprotov6.ProviderServer, error,
){
	"hpe": newProviderWithError,
}

func TestAccMorpheusOptionTypeTextareaExampleOk(t *testing.T) {
	t.Parallel()

	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	resourceConfig := optiontype.RenderMorpheusOptionTypeTextareaConfig(t, nil)

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
