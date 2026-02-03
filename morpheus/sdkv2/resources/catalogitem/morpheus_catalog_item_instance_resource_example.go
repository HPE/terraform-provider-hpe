// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package catalogitem

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
)

//go:generate../../../../bin/render -out examples/resources/morpheus_catalog_item_instance/resource.tf morpheus_catalog_item_instance_resource.tf.tmpl Name tfexample_instance_catalog Description "terraform example instance catalog item" ImagePath tfexample.png ImageName tfexample.png Enabled true Featured true Content "{\"name\":\"test\"}" Config "{\"name\":\"test\"}" Visibility private

// RenderCatalogItemInstanceConfig generates a Terraform configuration for catalog item instance resource.
// It accepts a name and a map of field overrides to customize the default values.
func RenderCatalogItemInstanceConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":        "Example",
		"Config":      "{\"name\":\"test\"}",
		"Content":     "{\"name\":\"test\"}",
		"Description": "terraform example instance catalog item",
		"Enabled":     "true",
		"Featured":    "true",
		"ImageName":   "tfexample.png",
		"ImagePath":   "tfexample.png",
		"Visibility":  "private",
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
	templatePath := filepath.Join(dir, "morpheus_catalog_item_instance_resource.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}
