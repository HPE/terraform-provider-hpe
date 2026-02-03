// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package cypher

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
)

//go:generate ../../../../../bin/render -out examples/resources/morpheus_cypher_secret/resource.tf morpheus_cypher_secret_resource_tf.tmpl Key apipassword Value password123 Ttl 86400

// RenderCypherSecretConfig generates a Terraform configuration
// for the hpe_morpheus_cypher_secret resource from the template file.
func RenderCypherSecretConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	// Default field values
	defaults := map[string]string{
		"Key":   "Example",
		"Value": "password123",
		"Ttl":   "86400",
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
	templatePath := filepath.Join(dir, "morpheus_cypher_secret_resource_tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}
