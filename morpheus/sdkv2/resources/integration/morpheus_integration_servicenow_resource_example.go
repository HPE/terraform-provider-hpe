// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package integration

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
)

//go:generate../../../../bin/render -out examples/resources/morpheus_integration_servicenow/resource.tf morpheus_integration_servicenow_resource.tf.tmpl Name "terraform servicenow integration" Enabled true Url "https://servicenowprod.service-now.com" Username "my-snow-username" Password "my-snow-password" DefaultCmdbBusinessClass "demo"

func RenderIntegrationServicenowConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":                     "Example",
		"Enabled":                  "true",
		"Url":                      "https://servicenowprod.service-now.com",
		"Username":                 "my-snow-username",
		"Password":                 "my-snow-password",
		"DefaultCmdbBusinessClass": "demo",
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
	templatePath := filepath.Join(dir, "morpheus_integration_servicenow_resource.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}
