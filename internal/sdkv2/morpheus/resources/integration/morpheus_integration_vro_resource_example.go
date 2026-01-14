// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package integration

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/testhelpers"
)

//go:generate go run ../../../../../cmd/render -out examples/resources/morpheus_integration_vro/resource.tf morpheus_integration_vro_resource.tf.tmpl Name "tfexample vro" Enabled true Url https://myvro/vco/api Username my-vro-username Password my-vro-password AuthType basic Tenant vsphere.local AuthId "1"

// RenderIntegrationVroConfig generates a Terraform configuration for the VRO integration
// resource. It accepts an optional map of field overrides to customize the default values.
func RenderIntegrationVroConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	// Default field values
	defaults := map[string]string{
		"Name":     "Example",
		"Enabled":  "true",
		"Url":      "https://myvro/vco/api",
		"Username": "my-vro-username",
		"Password": "my-vro-password",
		"AuthType": "basic",
		"Tenant":   "vsphere.local",
		"AuthId":   "1",
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
	templatePath := filepath.Join(dir, "morpheus_integration_vro_resource.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}
