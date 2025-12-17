// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package license

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/testhelpers"
)

//go:generate go run ../../../../../cmd/render -out examples/resources/morpheus_license/resource.tf morpheus_license_resource.tf.tmpl Key 22324FEF3WMCDMMSWE

// RenderLicenseConfig generates a Terraform configuration for the morpheus_license resource.
// It accepts optional overrides for field values. Default values are used if not overridden.
func RenderLicenseConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Key": "22324FEF3WMCDMMSWE",
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
	templatePath := filepath.Join(dir, "morpheus_license_resource.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		"Key",
		defaults["Key"],
	)
}
