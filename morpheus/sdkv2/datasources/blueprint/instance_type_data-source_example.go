// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package blueprint

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
)

//go:generate sh -c "../../../../bin/render -out examples/data-sources/morpheus_instance_type/data-source.tf instance_type_data-source.tf.tmpl Name '\"apache\"'"

// RenderInstanceTypeConfig generates a Terraform configuration for the tenant resource.
// It accepts optional overrides for field values. Default values are used if not overridden.
func RenderInstanceTypeConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name": "\"apache\"",
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
	templatePath := filepath.Join(dir, "instance_type_data-source.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}
