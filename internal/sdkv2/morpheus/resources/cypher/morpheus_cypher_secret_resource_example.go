// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package cypher

import (
	"testing"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/testhelpers"
)

//go:generate go run ../../../../../cmd/render -out examples/resources/morpheus_cypher_secret/resource.tf morpheus_cypher_secret_resource_tf.tmpl Key apipassword Value password123 Ttl 86400

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
