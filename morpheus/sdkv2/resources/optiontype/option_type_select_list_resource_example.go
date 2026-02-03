// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package optiontype

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
)

//go:generate sh -c " ../../../../../bin/render -out examples/resources/morpheus_option_type_select_list/resource.tf option_type_select_list_resource.tf.tmpl DefaultValue 'testing' DependentField 'dependent_example' Description 'Terraform select list option type example' DisplayValueOnDetails 'true' Editable 'true' ExportMeta 'true' FieldLabel 'numbers' FieldName 'tfSelectExample' HelpBlock 'fiwefw' Labels '[\"demo\", \"terraform\"]' Name 'tf_example_select_list_option_type' OptionListId '3' RequireField 'require_example' Required 'true' ShowOnEdit 'true' VisibilityField 'visibility_example'"

// RenderOptionTypeSelectListConfig generates a Terraform configuration for the tenant resource.
// It accepts optional overrides for field values. Default values are used if not overridden.
func RenderOptionTypeSelectListConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"DefaultValue":          "testing",
		"DependentField":        "dependent_example",
		"Description":           "Terraform select list option type example",
		"DisplayValueOnDetails": "true",
		"Editable":              "true",
		"ExportMeta":            "true",
		"FieldLabel":            "numbers",
		"FieldName":             "tfSelectExample",
		"HelpBlock":             "fiwefw",
		"Labels":                "[\"demo\", \"terraform\"]",
		"Name":                  "tf_example_select_list_option_type",
		"OptionListId":          "3",
		"RequireField":          "require_example",
		"Required":              "true",
		"ShowOnEdit":            "true",
		"VisibilityField":       "visibility_example",
	}

	// Apply overrides to defaults
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
	templatePath := filepath.Join(dir, "option_type_select_list_resource.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}
