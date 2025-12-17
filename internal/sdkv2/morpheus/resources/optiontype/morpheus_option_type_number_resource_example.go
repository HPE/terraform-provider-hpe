// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package optiontype

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/testhelpers"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

//go:generate go run ../../../../../cmd/render -out examples/resources/morpheus_option_type_number/resource.tf morpheus_option_type_number_resource.tf.tmpl Name tf_example_number_option_type Description "Terraform number option type example" Labels "[\"demo\", \"terraform\"]" FieldName tfNumberExample ExportMeta true DependentField dependent_example VisibilityField visibility_example RequireField require_example ShowOnEdit true Editable true DisplayValueOnDetails true FieldLabel "Number Example" Placeholder 12 DefaultValue 1 HelpBlock "Provide a number" Required true

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

	// Get the directory where this source file is located
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("unable to get current file path")
	}
	dir := filepath.Dir(filename)
	templatePath := filepath.Join(dir, "morpheus_option_type_number_resource.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
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
