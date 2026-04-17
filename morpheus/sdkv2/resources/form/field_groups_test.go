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

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	defer testhelpers.RecordResult(t)
	providerConfig := testhelpers.ProviderBlock()
	name := acctest.RandomWithPrefix(t.Name())
	code := toCode(name)
	fg1OTCode := code + "-fg1-ot"
	fg1OTName := name + " fg1 option type"
	fg2OTCode := code + "-fg2-ot"
	fg2OTName := name + " fg2 option type"

	resourceConfig, err := form.RenderFieldGroupsConfig(t, map[string]string{
		"Name":                      name,
		"Code":                      code,
		"FieldGroup1OptionTypeCode": fg1OTCode,
		"FieldGroup1OptionTypeName": fg1OTName,
		"FieldGroup2OptionTypeCode": fg2OTCode,
		"FieldGroup2OptionTypeName": fg2OTName,
	})
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
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "code", code),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "description", "demo"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "field_group.0.collapsed_by_default", "true"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "field_group.0.collapsible", "true"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "field_group.0.description", "testin"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "field_group.0.name", "fg1"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "field_group.0.option_type.0.code", fg1OTCode),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_form.example",
						"field_group.0.option_type.0.default_value",
						"Demo123",
					),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_form.example",
						"field_group.0.option_type.0.description",
						"Terraform text input example",
					),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_form.example",
						"field_group.0.option_type.0.display_value_on_details",
						"true",
					),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_form.example",
						"field_group.0.option_type.0.exclude_from_search",
						"true",
					),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "field_group.0.option_type.0.export_meta", "true"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "field_group.0.option_type.0.field_label", "Testin"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "field_group.0.option_type.0.field_name", "test"),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_form.example",
						"field_group.0.option_type.0.help_block",
						"Help block example",
					),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "field_group.0.option_type.0.hidden", "false"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "field_group.0.option_type.0.locked", "true"),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_form.example",
						"field_group.0.option_type.0.name",
						fg1OTName,
					),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_form.example",
						"field_group.0.option_type.0.placeholder",
						"Testing 123",
					),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "field_group.0.option_type.0.required", "true"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "field_group.0.option_type.0.type", "text"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "field_group.1.collapsed_by_default", "true"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "field_group.1.collapsible", "true"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "field_group.1.description", "testin"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "field_group.1.name", "fg2"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "field_group.1.option_type.0.code", fg2OTCode),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_form.example",
						"field_group.1.option_type.0.default_value",
						"Demo123",
					),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_form.example",
						"field_group.1.option_type.0.description",
						"Terraform text input example",
					),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_form.example",
						"field_group.1.option_type.0.display_value_on_details",
						"true",
					),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_form.example",
						"field_group.1.option_type.0.exclude_from_search",
						"true",
					),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "field_group.1.option_type.0.export_meta", "true"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "field_group.1.option_type.0.field_label", "Testin"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "field_group.1.option_type.0.field_name", "test"),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_form.example",
						"field_group.1.option_type.0.help_block",
						"Help block example",
					),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "field_group.1.option_type.0.hidden", "false"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "field_group.1.option_type.0.locked", "true"),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_form.example",
						"field_group.1.option_type.0.name",
						fg2OTName,
					),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_form.example",
						"field_group.1.option_type.0.placeholder",
						"Testing 123",
					),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "field_group.1.option_type.0.required", "true"),
					resource.TestCheckResourceAttr("hpe_morpheus_form.example", "field_group.1.option_type.0.type", "text"),
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
