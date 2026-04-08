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

func TestAccMorpheusFormNumberOk(t *testing.T) {
	t.Parallel()

	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()
	name := acctest.RandomWithPrefix(t.Name())

	resourceConfig, err := form.RenderNumberConfig(t, map[string]string{"Name": name})
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
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "option_type1_type", "number"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "option_type1_code", "number-input"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "option_type1_default_value", "4"),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_form.example",
						"option_type1_description",
						"Terraform number example",
					),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "option_type1_display_value_on_details", "true"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "option_type1_exclude_from_search", "true"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "option_type1_export_meta", "true"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "option_type1_field_label", "number input"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "option_type1_field_name", "numberInput"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "option_type1_help_block", "Is this working now"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "option_type1_hidden", "true"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "option_type1_locked", "true"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "option_type1_max_value", "44"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "option_type1_min_value", "3"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "option_type1_name", "tf number input example"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "option_type1_placeholder", "Testing 123"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "option_type1_required", "true"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "option_type1_step", "2"),
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
