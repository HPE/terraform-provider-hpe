// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package appliancesetup

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
)

//go:generate ../../../../bin/render -out examples/resources/morpheus_appliance_setup/example.tf example.tf.tmpl ResourceName example ApplianceName morpheus.example.com AccountName "My Company" Username admin Email admin@example.com FirstName Morpheus LastName Admin PasswordVariable admin_password LicenseVariable license_key CreateTimeout 40m

// RenderApplianceSetupConfig renders the example with the given overrides.
func RenderApplianceSetupConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"ResourceName":     "example",
		"ApplianceName":    "morpheus.example.com",
		"AccountName":      "My Company",
		"Username":         "admin",
		"Email":            "admin@example.com",
		"FirstName":        "Morpheus",
		"LastName":         "Admin",
		"PasswordVariable": "admin_password",
		"LicenseVariable":  "license_key",
		"CreateTimeout":    "40m",
	}

	for key, value := range overrides {
		defaults[key] = value
	}

	var args []string
	for key, value := range defaults {
		args = append(args, key, value)
	}

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("unable to get current file path")
	}

	templatePath := filepath.Join(filepath.Dir(filename), "example.tf.tmpl")

	return testhelpers.RenderExample(t, templatePath, args...)
}
