// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkdomain

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
)

//go:generate ../../../../bin/render -out examples/resources/morpheus_network_domain/resource.tf example.tf.tmpl Name "tfexampledomain" Description "Terraform example network domain" PublicZone "false" Visibility "private" Active "true" DomainPasswordWo "Secret123!" DomainUsername "joiner"

// RenderNetworkDomainConfig renders the example template with the given
// overrides, for use as a fixture in acceptance tests.
func RenderNetworkDomainConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":             "tfexampledomain",
		"Description":      "Terraform example network domain",
		"PublicZone":       "false",
		"Visibility":       "private",
		"Active":           "true",
		"DomainPasswordWo": "Secret123!",
		"DomainUsername":   "joiner",
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

	dir := filepath.Dir(filename)
	templatePath := filepath.Join(dir, "example.tf.tmpl")

	return testhelpers.RenderExample(t, templatePath, args...)
}
