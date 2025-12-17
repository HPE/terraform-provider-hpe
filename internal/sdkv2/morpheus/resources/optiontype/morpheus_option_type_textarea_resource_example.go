// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package optiontype

import (
	"testing"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/testhelpers"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

//go:generate go run ../../../../../cmd/render -out examples/resources/morpheus_option_type_textarea/resource.tf morpheus_option_type_textarea_resource.tf.tmpl Name tf_example_textarea_option_type Description "Terraform text area option type example" Labels ["demo","terraform"] FieldName textareaExample ExportMeta true DependentField dependent_example VisibilityField visibility_example RequireField require_example ShowOnEdit true Editable true DisplayValueOnDetails true FieldLabel "Text Area Example" Rows 5 Placeholder "example text" DefaultValue example HelpBlock "Terraform text area option type example" Required true VerifyPattern "a\\\\D{4}"

func RenderOptionTypeTextareaConfig(t *testing.T, overrides map[string]string) string {
	t.Helper()

	defaults := map[string]string{
		"Name":                  acctest.RandomWithPrefix(t.Name()),
		"Description":           "Terraform text area option type example",
		"Labels":                `["demo","terraform"]`,
		"FieldName":             "textareaExample",
		"ExportMeta":            "true",
		"DependentField":        "dependent_example",
		"VisibilityField":       "visibility_example",
		"RequireField":          "require_example",
		"ShowOnEdit":            "true",
		"Editable":              "true",
		"DisplayValueOnDetails": "true",
		"FieldLabel":            "Text Area Example",
		"Rows":                  "5",
		"Placeholder":           "example text",
		"DefaultValue":          "example",
		"HelpBlock":             "Terraform text area option type example",
		"Required":              "true",
		"VerifyPattern":         `a\\D{4}`,
	}

	for key, value := range overrides {
		defaults[key] = value
	}

	resourceConfig, err := testhelpers.RenderExample(t, "morpheus_option_type_textarea_resource.tf.tmpl",
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
		"Rows", defaults["Rows"],
		"Placeholder", defaults["Placeholder"],
		"DefaultValue", defaults["DefaultValue"],
		"HelpBlock", defaults["HelpBlock"],
		"Required", defaults["Required"],
		"VerifyPattern", defaults["VerifyPattern"],
	)
	if err != nil {
		t.Fatal(err)
	}

	return resourceConfig
}
