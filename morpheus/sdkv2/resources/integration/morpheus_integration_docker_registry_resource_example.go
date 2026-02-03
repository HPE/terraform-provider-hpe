// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package integration

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
)

//go:generate../../../../bin/render -out examples/resources/morpheus_integration_docker_registry/resource.tf morpheus_integration_docker_registry_resource.tf.tmpl Name tfexampledockerregistry Enabled true Url https://index.docker.io/v1/ Username admin Password password123

// RenderIntegrationDockerRegistryConfig renders the Docker Registry integration
// resource configuration with the provided field overrides. Default values are used for any
// fields not specified.
func RenderIntegrationDockerRegistryConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":     "Example",
		"Enabled":  "true",
		"Url":      "https://index.docker.io/v1/",
		"Username": "admin",
		"Password": "password123",
	}

	// Apply overrides
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
	templatePath := filepath.Join(dir, "morpheus_integration_docker_registry_resource.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}
