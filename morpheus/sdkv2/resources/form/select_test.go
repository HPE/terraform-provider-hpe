// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package form_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus"
	sdkv2morpheus "github.com/HPE/terraform-provider-hpe/morpheus/sdkv2"
	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/resources/form"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
)

func TestAccMorpheusFormSelectOk(t *testing.T) {
	t.Parallel()

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	defer testhelpers.RecordResult(t)
	providerConfig := testhelpers.ProviderBlock()
	name := acctest.RandomWithPrefix(t.Name())
	code := toCode(name)
	optTypeCode := code + "-ot"
	optTypeName := name + " option type"

	optionListConfig := fmt.Sprintf(`
resource "hpe_morpheus_option_list_manual" "example" {
  name    = %q
  dataset = "[{\"name\": \"Level 1\",\"value\":\"level1\"},{\"name\": \"Level 2\",\"value\":\"level2\"}]"
  real_time = true
}
`, name+" option list")

	resourceConfig, err := form.RenderSelectConfig(t, map[string]string{
		"Name":                   name,
		"Code":                   code,
		"OptionTypeCode":         optTypeCode,
		"OptionTypeName":         optTypeName,
		"OptionTypeOptionListId": "hpe_morpheus_option_list_manual.example.id",
	})
	if err != nil {
		t.Fatal(err)
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), sdkv2morpheus.Provider()),
		Steps: []resource.TestStep{
			{
				Config:             providerConfig + optionListConfig + resourceConfig,
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "code", code),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "description", "demo"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "option_type.0.type", "select"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "option_type.0.code", optTypeCode),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "option_type.0.default_value", "test123"),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_form.example",
						"option_type.0.description",
						"Terraform select example",
					),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "option_type.0.display_value_on_details", "true"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "option_type.0.exclude_from_search", "true"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "option_type.0.export_meta", "true"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "option_type.0.field_label", "Select Test"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "option_type.0.field_name", "selectTest"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "option_type.0.help_block", "Select an option"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "option_type.0.hidden", "true"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "option_type.0.locked", "true"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "option_type.0.name", optTypeName),
					resource.TestCheckResourceAttrPair(
						"hpe_morpheus_form.example", "option_type.0.option_list_id",
						"hpe_morpheus_option_list_manual.example", "id",
					),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "option_type.0.placeholder", "Testing 123"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "option_type.0.required", "true"),
				),
			},
			{
				Config:             providerConfig + optionListConfig + resourceConfig,
				ExpectNonEmptyPlan: false,
				PlanOnly:           true,
			},
		},
	})
}
