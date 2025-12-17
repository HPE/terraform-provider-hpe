// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package integration

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/testhelpers"
)

//go:generate go run ../../../../../cmd/render -out examples/resources/morpheus_integration_chef/resource.tf morpheus_integration_chef_resource.tf.tmpl Name "tfexample chef integration" Enabled true Url "https://chef.morpheusdata.com" Version "15.9.38" WindowsVersion "15.9.38" WindowsMsiInstallUrl "https://packages.chef.io" Organization "morpheus" Username "admin" PrivateKey "EXAMPLEPRIVATEKEY" OrganizationValidatorKey "EXAMPLEPRIVATEKEY"

// RenderIntegrationChefConfig generates a Terraform configuration for the Chef integration resource
// using default values that can be overridden via the overrides map.
func RenderIntegrationChefConfig(t *testing.T, name string, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":                     name,
		"Enabled":                  "true",
		"Url":                      "https://chef.morpheusdata.com",
		"Version":                  "15.9.38",
		"WindowsVersion":           "15.9.38",
		"WindowsMsiInstallUrl":     "https://packages.chef.io",
		"Organization":             "morpheus",
		"Username":                 "admin",
		"PrivateKey":               "EXAMPLEPRIVATEKEY",
		"OrganizationValidatorKey": "EXAMPLEPRIVATEKEY",
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
	templatePath := filepath.Join(dir, "morpheus_integration_chef_resource.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		"Name", defaults["Name"],
		"Enabled", defaults["Enabled"],
		"Url", defaults["Url"],
		"Version", defaults["Version"],
		"WindowsVersion", defaults["WindowsVersion"],
		"WindowsMsiInstallUrl", defaults["WindowsMsiInstallUrl"],
		"Organization", defaults["Organization"],
		"Username", defaults["Username"],
		"PrivateKey", defaults["PrivateKey"],
		"OrganizationValidatorKey", defaults["OrganizationValidatorKey"],
	)
}
