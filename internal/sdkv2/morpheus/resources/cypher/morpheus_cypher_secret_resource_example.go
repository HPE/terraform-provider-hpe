// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package cypher

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/testhelpers"
)

//go:generate go run ../../../../../cmd/render -out examples/resources/morpheus_cypher_secret/resource.tf morpheus_cypher_secret_resource_tf.tmpl Key apipassword Value password123 Ttl 86400

// RenderCypherSecretConfig generates a Terraform configuration
// for the hpe_morpheus_cypher_secret resource from the template file.
func RenderCypherSecretConfig(
	t *testing.T,
	name string,
	overrides map[string]string,
) (string, error) {
	t.Helper()

	// Default field values
	defaults := map[string]string{
		"Key":   name,
		"Value": "password123",
		"Ttl":   "86400",
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
	templatePath := filepath.Join(dir, "morpheus_cypher_secret_resource_tf.tmpl")

	return testhelpers.RenderExample(t,
		templatePath,
		"Key", defaults["Key"],
		"Value", defaults["Value"],
		"Ttl", defaults["Ttl"],
	)
}
