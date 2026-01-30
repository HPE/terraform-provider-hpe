// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package integration

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/testhelpers"
)

//go:generate ../../../../../bin/render -out examples/resources/morpheus_integration_chef/resource.tf morpheus_integration_chef_resource.tf.tmpl Name "tfexample chef integration" Enabled true Url "https://chef.morpheusdata.com" Version "15.9.38" WindowsVersion "15.9.38" WindowsMsiInstallUrl "https://packages.chef.io" Organization "morpheus" Username "admin" PrivateKey "EXAMPLEPRIVATEKEY" OrganizationValidatorKey "EXAMPLEPRIVATEKEY"

// RenderIntegrationChefConfig generates a Terraform configuration for the Chef integration resource
// using default values that can be overridden via the overrides map.
func RenderIntegrationChefConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":                     "Example",
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
	templatePath := filepath.Join(dir, "morpheus_integration_chef_resource.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}
