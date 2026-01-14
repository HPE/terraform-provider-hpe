// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package catalogitem

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/testhelpers"
)

//go:generate sh -c "go run ../../../../../cmd/render -out examples/resources/morpheus_catalog_item_app_blueprint/resource.tf catalog_item_app_blueprint_resource.tf.tmpl AppSpec 'file(\"${path.module}/appSpec.yaml\")' BlueprintId '5' Content 'file(\"${path.module}/catalog-data.md\")' DarkLogoImageName 'tfexampledark.png' DarkLogoImagePath 'tfexampledark.png' Description 'terraform example app blueprint catalog item' Enabled 'true' Featured 'true' Labels '[\"aws\", \"demo\", \"testing\"]' LogoImageName 'tfexample.png' LogoImagePath 'tfexample.png' Name 'tfexample_app_blueprint_catalog' OptionTypeIds '[2056, 2006, 2058]' Visibility 'public'"

// RenderCatalogItemAppBlueprintConfig generates a Terraform configuration for the tenant resource.
// It accepts optional overrides for field values. Default values are used if not overridden.
func RenderCatalogItemAppBlueprintConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"AppSpec":           "file(\"${path.module}/appSpec.yaml\")",
		"BlueprintId":       "5",
		"Content":           "file(\"${path.module}/catalog-data.md\")",
		"DarkLogoImageName": "tfexampledark.png",
		"DarkLogoImagePath": "tfexampledark.png",
		"Description":       "terraform example app blueprint catalog item",
		"Enabled":           "true",
		"Featured":          "true",
		"Labels":            "[\"aws\", \"demo\", \"testing\"]",
		"LogoImageName":     "tfexample.png",
		"LogoImagePath":     "tfexample.png",
		"Name":              "tfexample_app_blueprint_catalog",
		"OptionTypeIds":     "[2056, 2006, 2058]",
		"Visibility":        "public",
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
	templatePath := filepath.Join(dir, "catalog_item_app_blueprint_resource.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}
