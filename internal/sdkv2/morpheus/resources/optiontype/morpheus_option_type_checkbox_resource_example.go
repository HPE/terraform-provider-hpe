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

//go:generate go run ../../../../../cmd/render -out examples/resources/morpheus_option_type_checkbox/resource.tf morpheus_option_type_checkbox_resource.tf.tmpl Name tfcheckboxexample Description "Terraform checkbox option type example" Labels "[\"demo\", \"terraform\"]" FieldName checkbox_example ExportMeta true DependentField dependent_example VisibilityField visibility_example RequireField require_example ShowOnEdit true Editable true DisplayValueOnDetails true FieldLabel "Checkbox Example" DefaultChecked true

// RenderOptionTypeCheckboxConfig generates a Terraform configuration for the morpheus_option_type_checkbox resource.
// It accepts an optional map of field overrides to customize the default values.
// Supported override keys: "Name", "Description", "Labels", "FieldName", "ExportMeta", "DependentField",
// "VisibilityField", "RequireField", "ShowOnEdit", "Editable", "DisplayValueOnDetails", "FieldLabel", "DefaultChecked"
func RenderOptionTypeCheckboxConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":                  acctest.RandomWithPrefix(t.Name()),
		"Description":           "Terraform checkbox option type example",
		"Labels":                "[\"demo\", \"terraform\"]",
		"FieldName":             "checkbox_example",
		"ExportMeta":            "true",
		"DependentField":        "dependent_example",
		"VisibilityField":       "visibility_example",
		"RequireField":          "require_example",
		"ShowOnEdit":            "true",
		"Editable":              "true",
		"DisplayValueOnDetails": "true",
		"FieldLabel":            "Checkbox Example",
		"DefaultChecked":        "true",
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
	templatePath := filepath.Join(dir, "morpheus_option_type_checkbox_resource.tf.tmpl")

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
		"DefaultChecked", defaults["DefaultChecked"],
	)
}
