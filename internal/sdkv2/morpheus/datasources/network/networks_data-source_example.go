// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package network

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/testhelpers"
)

//go:generate sh -c "go run ../../../../../cmd/render -out examples/data-sources/morpheus_networks/data-source.tf networks_data-source.tf.tmpl CloudId '3' Name '\"name\"' SortAscending 'true' Values '[\"Test*\"]'"

// RenderNetworksConfig generates a Terraform configuration for the networks resource.
// It accepts optional overrides for field values. Default values are used if not overridden.
func RenderNetworksConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"CloudId":       "3",
		"Name":          "\"name\"",
		"SortAscending": "true",
		"Values":        "[\"Test*\"]",
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
	templatePath := filepath.Join(dir, "networks_data-source.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}
