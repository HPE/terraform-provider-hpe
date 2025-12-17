// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package optionlist

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/testhelpers"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

//go:generate go run ../../../../../cmd/render -out examples/resources/morpheus_option_list_manual/resource.tf morpheus_option_list_manual_resource.tf.tmpl Name tf_example_manual_option_list Description "Terraform manual option list example" Dataset "[{\"name\": \"Level 1\",\"value\":\"level1\"},\n {\"name\": \"Level 2\",\"value\":\"level2\"},\n {\"name\": \"Level 3\",\"value\":\"level3\"}\n]" RealTime true

// RenderOptionListManualConfig generates a Terraform configuration for the
// morpheus_option_list_manual resource. It accepts an optional map of field overrides to
// customize the default values. Supported override keys: "Name", "Description", "Dataset", "RealTime"
func RenderOptionListManualConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":        acctest.RandomWithPrefix(t.Name()),
		"Description": "Terraform manual option list example",
		"Dataset": "[{\"name\": \"Level 1\",\"value\":\"level1\"},\n " +
			"{\"name\": \"Level 2\",\"value\":\"level2\"},\n " +
			"{\"name\": \"Level 3\",\"value\":\"level3\"}\n]",
		"RealTime": "true",
	}

	// Apply overrides
	for key, value := range overrides {
		defaults[key] = value
	}

	//nolint: lll
	// Get the directory where this source file is located
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("unable to get current file path")
	}
	dir := filepath.Dir(filename)
	templatePath := filepath.Join(dir, "morpheus_option_list_manual_resource.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		"Name", defaults["Name"],
		"Description", defaults["Description"],
		"Dataset", defaults["Dataset"],
		"RealTime", defaults["RealTime"],
	)
}
