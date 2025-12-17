// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package integration

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/testhelpers"
)

//go:generate go run ../../../../../cmd/render -out examples/resources/morpheus_integration_vro/resource.tf morpheus_integration_vro_resource.tf.tmpl Name "tfexample vro" Enabled true Url https://myvro/vco/api Username my-vro-username Password my-vro-password AuthType basic Tenant vsphere.local

// RenderIntegrationVroConfig generates a Terraform configuration for the VRO integration
// resource. It accepts an optional map of field overrides to customize the default values.
func RenderIntegrationVroConfig(t *testing.T, name string, overrides map[string]string) (string, error) {
	t.Helper()

	// Default field values
	defaults := map[string]string{
		"Name":     name,
		"Enabled":  "true",
		"Url":      "https://myvro/vco/api",
		"Username": "my-vro-username",
		"Password": "my-vro-password",
		"AuthType": "basic",
		"Tenant":   "vsphere.local",
	}

	// Apply overrides to defaults
	for key, value := range overrides {
		defaults[key] = value
	}

	// Get the directory where this source file is located
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("unable to get current file path")
	}
	dir := filepath.Dir(filename)
	templatePath := filepath.Join(dir, "morpheus_integration_vro_resource.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		"Name", defaults["Name"],
		"Enabled", defaults["Enabled"],
		"Url", defaults["Url"],
		"Username", defaults["Username"],
		"Password", defaults["Password"],
		"AuthType", defaults["AuthType"],
		"Tenant", defaults["Tenant"],
	)
}
