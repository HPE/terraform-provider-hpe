// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

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

func TestAccMorpheusOptionTypePasswordExampleOk(t *testing.T) {
	t.Parallel()

	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())

	resourceConfig, err := optiontype.RenderOptionTypePasswordConfig(t, map[string]string{"Name": name})
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_option_type_password.tf_example_password_option_type",
			"name",
			name,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_option_type_password.tf_example_password_option_type",
			"description",
			"Terraform password option type example",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_option_type_password.tf_example_password_option_type",
			"field_name",
			"tfPasswordExample",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_option_type_password.tf_example_password_option_type",
			"export_meta",
			"true",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_option_type_password.tf_example_password_option_type",
			"dependent_field",
			"dependent_example",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_option_type_password.tf_example_password_option_type",
			"visibility_field",
			"visibility_example",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_option_type_password.tf_example_password_option_type",
			"require_field",
			"require_example",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_option_type_password.tf_example_password_option_type",
			"show_on_edit",
			"true",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_option_type_password.tf_example_password_option_type",
			"editable",
			"true",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_option_type_password.tf_example_password_option_type",
			"display_value_on_details",
			"true",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_option_type_password.tf_example_password_option_type",
			"field_label",
			"numbers",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_option_type_password.tf_example_password_option_type",
			"placeholder",
			"fewf",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_option_type_password.tf_example_password_option_type",
			"default_value",
			"testing",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_option_type_password.tf_example_password_option_type",
			"help_block",
			"fiwefw",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_option_type_password.tf_example_password_option_type",
			"required",
			"true",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_option_type_password.tf_example_password_option_type",
			"verify_pattern",
			"a\\D{4}",
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
