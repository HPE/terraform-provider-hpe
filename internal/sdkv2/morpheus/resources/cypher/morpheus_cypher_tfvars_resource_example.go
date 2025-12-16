// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package cypher

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"
	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/testhelpers"
)

//go:generate sh -c "go run ../../../../../cmd/render -out examples/resources/morpheus_cypher_tfvars/resource.tf hpe_morpheus_cypher_tfvars_resource.tf.tmpl Key securetfvars Value 'account=12345\npassword=supersecure' Ttl 86400"

func RenderCypherTfvarsConfig(
	t *testing.T,
	name string,
	overrides map[string]string,
) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Key":   name,
		"Ttl":   "86400",
		"Value": "account=12345\npassword=supersecure",
	}

	for k, v := range overrides {
		defaults[k] = v
	}

	// Get the directory where this source file is located
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("unable to get current file path")
	}
	dir := filepath.Dir(filename)
	templatePath := filepath.Join(dir, "hpe_morpheus_cypher_tfvars_resource.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		"Key", defaults["Key"],
		"Ttl", defaults["Ttl"],
		"Value", defaults["Value"],
	)
}

// RenderMorpheusCypherSecretConfig generates a Terraform configuration
// for the hpe_morpheus_cypher_secret resource from the template file.
func RenderMorpheusCypherSecretConfig(
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

	// Build arguments for RenderExample
	var args []string
	for key, value := range defaults {
		args = append(args, key, value)
	}

	return testhelpers.RenderExample(t, "morpheus_cypher_secret_resource_tf.tmpl", args...)
}

