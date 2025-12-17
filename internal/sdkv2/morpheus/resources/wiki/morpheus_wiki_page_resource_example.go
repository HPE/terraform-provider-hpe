// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package wiki

import (
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"

	"fmt"
	"path/filepath"
	"runtime"
	"testing"
	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/testhelpers"
)

//go:generate go run ../../../../../cmd/render -out examples/resources/morpheus_wiki_page/resource.tf morpheus_wiki_page_resource.tf.tmpl Name tfexample_wiki_page Category morpheus-terraform

// RenderWikiPageConfig generates a Terraform configuration for the wiki page resource.
// It accepts a map of field overrides to customize default values.
func RenderWikiPageConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	// Default field values
	defaults := map[string]string{
		"Name":     acctest.RandomWithPrefix(t.Name()),
		"Category": "morpheus-terraform",
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
	templatePath := filepath.Join(dir, "morpheus_wiki_page_resource.tf.tmpl")

	return testhelpers.RenderExample(t, templatePath,
		"Name", defaults["Name"],
		"Category", defaults["Category"],
	)
}

