// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package optiontype

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
)

//go:generate ../../../../../bin/render -out examples/resources/morpheus_option_type_hidden/resource.tf morpheus_option_type_hidden_resource.tf.tmpl Name tf_example_hidden_option_type Description "Terraform hidden option type example" Labels ["demo","terraform"] FieldName hidden_example ExportMeta true DependentField dependent_example VisibilityField visibility_example RequireField require_example ShowOnEdit true Editable true DisplayValueOnDetails true DefaultValue example

// RenderOptionTypeHiddenConfig generates a Terraform configuration for the
// morpheus_option_type_hidden resource. It accepts an optional map of field overrides to customize
// the default values. Supported override keys: "Name", "Description", "Labels", "FieldName",
// "ExportMeta", "DependentField", "VisibilityField", "RequireField", "ShowOnEdit", "Editable",
// "DisplayValueOnDetails", "DefaultValue"
func RenderOptionTypeHiddenConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":                  "Example",
		"Description":           "Terraform hidden option type example",
		"Labels":                `["demo","terraform"]`,
		"FieldName":             "hidden_example",
		"ExportMeta":            "true",
		"DependentField":        "dependent_example",
		"VisibilityField":       "visibility_example",
		"RequireField":          "require_example",
		"ShowOnEdit":            "true",
		"Editable":              "true",
		"DisplayValueOnDetails": "true",
		"DefaultValue":          "example",
	}

	// Apply overrides
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
	templatePath := filepath.Join(dir, "morpheus_option_type_hidden_resource.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}
