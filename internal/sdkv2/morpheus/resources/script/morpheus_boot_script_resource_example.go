// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package script

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"
	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/testhelpers"
)

//go:generate go run ../../../../../cmd/render -out examples/resources/morpheus_boot_script/resource.tf morpheus_boot_script_resource.tf.tmpl Name "TF Example Boot Script" Content "ls"

// RenderPreseedScriptConfig renders a Terraform configuration
// for preseed_script resource.
// It accepts a name and a map of overrides to customize the default field values.
func RenderPreseedScriptConfig(
	t *testing.T,
	name string,
	overrides map[string]string,
) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":    name,
		"Content": "ls",
	}

	for key, value := range overrides {
		defaults[key] = value
	}

	// Get the directory where this source file is located
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("unable to get current file path")
	}
	dir := filepath.Dir(filename)
	templatePath := filepath.Join(dir, "morpheus_preseed_script_resource.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		"Name", defaults["Name"],
		"Content", defaults["Content"],
	)
}

// RenderBootScriptConfig renders a Terraform configuration for boot_script resource.
// It accepts a name and a map of overrides to customize the default field values.
func RenderBootScriptConfig(
	t *testing.T,
	name string,
	overrides map[string]string,
) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":    name,
		"Content": "ls",
	}

	for key, value := range overrides {
		defaults[key] = value
	}

	// Get the directory where this source file is located
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("unable to get current file path")
	}
	dir := filepath.Dir(filename)
	templatePath := filepath.Join(dir, "morpheus_boot_script_resource.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		"Name", defaults["Name"],
		"Content", defaults["Content"],
	)
}

