// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package catalogitem

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/testhelpers"
)

//go:generate go run ../../../../../cmd/render -out examples/resources/catalog_item_workflow/catalog_item_workflow_resource.tf catalog_item_workflow_resource.tf.tmpl ContextType "appliance" Content "file(\"${path.module}/catalog-data.md\")" DarkLogoImageName "wordpressbak.png" DarkLogoImagePath "wordpressbak.png" Description "Example Terraform workflow catalog item" Enabled "true" Featured "true" Labels "[\"aws\", \"demo\"]" LogoImageName "wordpress.png" LogoImagePath "wordpress.png" Name "tfexample_workflow_catalog_item" OptionTypeIds "[2056, 2006]" Visibility "public" WorkflowId "1"

// RenderCatalogItemWorkflowConfig generates a Terraform configuration for the tenant resource.
// It accepts optional overrides for field values. Default values are used if not overridden.

func RenderCatalogItemWorkflowConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"ContextType":       "appliance",
		"Content":           "file(\"${path.module}/catalog-data.md\")",
		"DarkLogoImageName": "wordpressbak.png",
		"DarkLogoImagePath": "wordpressbak.png",
		"Description":       "Example Terraform workflow catalog item",
		"Enabled":           "true",
		"Featured":          "true",
		"Labels":            "[\"aws\", \"demo\"]",
		"LogoImageName":     "wordpress.png",
		"LogoImagePath":     "wordpress.png",
		"Name":              "tfexample_workflow_catalog_item",
		"OptionTypeIds":     "[2056, 2006]",
		"Visibility":        "public",
		"WorkflowId":        "1",
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
	templatePath := filepath.Join(dir, "catalog_item_workflow_resource.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}
