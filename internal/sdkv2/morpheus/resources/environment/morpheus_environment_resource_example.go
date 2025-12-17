// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package environment

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/testhelpers"
)

//go:generate go run ../../../../../cmd/render -out examples/resources/morpheus_environment/resource.tf morpheus_environment_resource.tf.tmpl Active true Code tfexample Description "Terraform Example" Name tfexample

// RenderEnvironmentConfig renders the environment resource configuration with default values
// that can be overridden by providing a map of field name to value.
func RenderEnvironmentConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Active":      "true",
		"Code":        "tfexample",
		"Description": "Terraform Example",
		"Name":        acctest.RandomWithPrefix(t.Name()),
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
	templatePath := filepath.Join(dir, "morpheus_environment_resource.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		"Active", defaults["Active"],
		"Code", defaults["Code"],
		"Description", defaults["Description"],
		"Name", defaults["Name"],
	)
}

