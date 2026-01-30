// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package integration

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/testhelpers"
)

//go:generate ../../../../../bin/render -out examples/resources/morpheus_integration_ansible_tower/resource.tf morpheus_integration_ansible_tower_resource.tf.tmpl Name "tfexample ansible tower integration" Enabled true Url "https://ansibletower01.morpheusdata.com" Username admin Password password123

// RenderIntegrationAnsibleTowerConfig renders the Ansible Tower integration resource configuration
// with default values that can be overridden via the overrides parameter.
func RenderIntegrationAnsibleTowerConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":     "Example",
		"Enabled":  "true",
		"Url":      "https://ansibletower01.morpheusdata.com",
		"Username": "admin",
		"Password": "password123",
	}

	// Override defaults with provided values
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
	templatePath := filepath.Join(dir, "morpheus_integration_ansible_tower_resource.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}
