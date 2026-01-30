// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package optiontype

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/testhelpers"
)

//go:generate sh -c " ../../../../../bin/render -out examples/resources/morpheus_option_type_text/resource.tf option_type_text_resource.tf.tmpl DefaultValue 'testing' DependentField 'dependent_example' Description 'Terraform text option type example' DisplayValueOnDetails 'true' Editable 'true' ExportMeta 'true' FieldLabel 'numbers' FieldName 'test1' HelpBlock 'fiwefw' Labels '[\"demo\", \"terraform\"]' Name 'tf_example_text_option_type' Placeholder 'fewf' RequireField 'require_example' Required 'true' ShowOnEdit 'true' VerifyPattern 'a\\\\D{4}' VisibilityField 'visibility_example'"

// RenderOptionTypeTextConfig generates a Terraform configuration for the tenant resource.
// It accepts optional overrides for field values. Default values are used if not overridden.
func RenderOptionTypeTextConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"DefaultValue":          "testing",
		"DependentField":        "dependent_example",
		"Description":           "Terraform text option type example",
		"DisplayValueOnDetails": "true",
		"Editable":              "true",
		"ExportMeta":            "true",
		"FieldLabel":            "numbers",
		"FieldName":             "test1",
		"HelpBlock":             "fiwefw",
		"Labels":                "[\"demo\", \"terraform\"]",
		"Name":                  "tf_example_text_option_type",
		"Placeholder":           "fewf",
		"RequireField":          "require_example",
		"Required":              "true",
		"ShowOnEdit":            "true",
		"VerifyPattern":         "a\\\\D{4}",
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
	templatePath := filepath.Join(dir, "option_type_text_resource.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}
