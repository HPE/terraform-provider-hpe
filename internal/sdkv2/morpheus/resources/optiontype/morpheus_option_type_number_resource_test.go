// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package optiontype_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/testhelpers"
)

// RenderMorpheusOptionTypeNumberConfig generates a Terraform configuration for the
// morpheus_option_type_number resource. It accepts an optional map of field overrides to
// customize the default values.
// Supported override keys: "Name", "Description", "Labels", "FieldName", "ExportMeta",
// "DependentField", "VisibilityField", "RequireField", "ShowOnEdit", "Editable",
// "DisplayValueOnDetails", "FieldLabel", "Placeholder", "DefaultValue", "HelpBlock", "Required"
func RenderMorpheusOptionTypeNumberConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":                  acctest.RandomWithPrefix(t.Name()),
		"Description":           "Terraform number option type example",
		"Labels":                "[\"demo\", \"terraform\"]",
		"FieldName":             "tfNumberExample",
		"ExportMeta":            "true",
		"DependentField":        "dependent_example",
		"VisibilityField":       "visibility_example",
		"RequireField":          "require_example",
		"ShowOnEdit":            "true",
		"Editable":              "true",
		"DisplayValueOnDetails": "true",
		"FieldLabel":            "Number Example",
		"Placeholder":           "12",
		"DefaultValue":          "1",
		"HelpBlock":             "Provide a number",
		"Required":              "true",
	}

	// Apply overrides
	for key, value := range overrides {
		defaults[key] = value
	}

	return testhelpers.RenderExample(t, "morpheus_option_type_number_resource.tf.tmpl",
		"Name", defaults["Name"],
		"Description", defaults["Description"],
		"Labels", defaults["Labels"],
		"FieldName", defaults["FieldName"],
		"ExportMeta", defaults["ExportMeta"],
		"DependentField", defaults["DependentField"],
		"VisibilityField", defaults["VisibilityField"],
		"RequireField", defaults["RequireField"],
		"ShowOnEdit", defaults["ShowOnEdit"],
		"Editable", defaults["Editable"],
		"DisplayValueOnDetails", defaults["DisplayValueOnDetails"],
		"FieldLabel", defaults["FieldLabel"],
		"Placeholder", defaults["Placeholder"],
		"DefaultValue", defaults["DefaultValue"],
		"HelpBlock", defaults["HelpBlock"],
		"Required", defaults["Required"],
	)
}

func TestAccMorpheusOptionTypeNumberExampleOk(t *testing.T) {
	t.Parallel()

	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())

	resourceConfig, err := RenderMorpheusOptionTypeNumberConfig(t, map[string]string{"Name": name})
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_option_type_number."+name,
			"name",
			name,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_option_type_number."+name,
			"description",
			"Terraform number option type example",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_option_type_number."+name,
			"field_name",
			"tfNumberExample",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_option_type_number."+name,
			"export_meta",
			"true",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_option_type_number."+name,
			"dependent_field",
			"dependent_example",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_option_type_number."+name,
			"visibility_field",
			"visibility_example",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_option_type_number."+name,
			"require_field",
			"require_example",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_option_type_number."+name,
			"show_on_edit",
			"true",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_option_type_number."+name,
			"editable",
			"true",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_option_type_number."+name,
			"display_value_on_details",
			"true",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_option_type_number."+name,
			"field_label",
			"Number Example",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_option_type_number."+name,
			"placeholder",
			"12",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_option_type_number."+name,
			"default_value",
			"1",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_option_type_number."+name,
			"help_block",
			"Provide a number",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_option_type_number."+name,
			"required",
			"true",
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
