// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package optiontype

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/testhelpers"
)

//go:generate go run ../../../../../cmd/render -out examples/resources/morpheus_option_type_textarea/resource.tf morpheus_option_type_textarea_resource.tf.tmpl Name tf_example_textarea_option_type Description "Terraform text area option type example" Labels ["demo","terraform"] FieldName textareaExample ExportMeta true DependentField dependent_example VisibilityField visibility_example RequireField require_example ShowOnEdit true Editable true DisplayValueOnDetails true FieldLabel "Text Area Example" Rows 5 Placeholder "example text" DefaultValue example HelpBlock "Terraform text area option type example" Required true VerifyPattern "a\\\\D{4}"

func RenderOptionTypeTextareaConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":                  "Example",
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

	var args []string
	for key, value := range defaults {
		args = append(args, key, value)
	}

	// Get the directory where this source file is located
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("unable to get current file path")
	}
	dir := filepath.Dir(filename)
	templatePath := filepath.Join(dir, "morpheus_option_type_textarea_resource.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}
