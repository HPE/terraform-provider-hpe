// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package integration

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/testhelpers"
)

//go:generate go run ../../../../../cmd/render -out examples/resources/morpheus_integration_docker_registry/resource.tf morpheus_integration_docker_registry_resource.tf.tmpl Name tfexampledockerregistry Enabled true Url https://index.docker.io/v1/ Username admin Password password123

// RenderMorpheusIntegrationDockerRegistryConfig renders the Docker Registry integration
// resource configuration with the provided field overrides. Default values are used for any
// fields not specified.
func RenderMorpheusIntegrationDockerRegistryConfig(
	t *testing.T,
	name string,
	overrides map[string]string,
) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":     name,
		"Enabled":  "true",
		"Url":      "https://index.docker.io/v1/",
		"Username": "admin",
		"Password": "password123",
	}

	// Apply overrides
	for key, value := range overrides {
		defaults[key] = value
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
		"Name", defaults["Name"],
		"Enabled", defaults["Enabled"],
		"Url", defaults["Url"],
		"Username", defaults["Username"],
		"Password", defaults["Password"],
	)
}
