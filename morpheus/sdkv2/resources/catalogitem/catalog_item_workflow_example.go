// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package catalogitem

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
)

//go:generate ../../../../../bin/render -out examples/resources/morpheus_catalog_item_workflow/resource.tf catalog_item_workflow_resource.tf.tmpl Name tfexample_workflow_catalog_item Description "Example Terraform workflow catalog item" LogoImagePath wordpress.png LogoImageName wordpress.png DarkLogoImagePath wordpressbak.png DarkLogoImageName wordpressbak.png Enabled true Featured true Labels "[\"terraform\",\"demo\"]" WorkflowId 1 ContextType appliance Content "\"Example catalog content\"" Visibility public

// RenderCatalogItemWorkflowConfig generates a Terraform configuration for catalog item workflow resource.
// It accepts a name and a map of field overrides to customize the default values.
func RenderCatalogItemWorkflowConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":              "tfexample_workflow_catalog_item",
		"Description":       "Example Terraform workflow catalog item",
		"LogoImagePath":     "tfexample.png",
		"LogoImageName":     "tfexample.png",
		"DarkLogoImagePath": "tfexample.png",
		"DarkLogoImageName": "tfexample.png",
		"Enabled":           "true",
		"Featured":          "true",
		"Labels":            "[\"terraform\",\"demo\"]",
		"WorkflowId":        "1",
		"ContextType":       "appliance",
		"Content":           "\"Example catalog content\"",
		"Visibility":        "public",
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
	templatePath := filepath.Join(dir, "catalog_item_workflow_resource.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}
