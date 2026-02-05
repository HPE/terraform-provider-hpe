// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package workflow

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
)

//go:generate../../../../bin/render -out examples/resources/morpheus_workflow_operational/resource.tf workflow_operational_resource.tf.tmpl Name tf_example_operational_workflow Description "Terraform operational workflow example" Labels "[\"demo\", \"terraform\"]" Platform all Visibility private AllowCustomConfig true

// RenderWorkflowOperationalConfig generates a Terraform configuration for the workflow operational resource.
// It accepts optional overrides for field values. Default values are used if not overridden.
func RenderWorkflowOperationalConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":              "tf_example_operational_workflow",
		"Description":       "Terraform operational workflow example",
		"Labels":            "[\"demo\", \"terraform\"]",
		"Platform":          "all",
		"Visibility":        "private",
		"AllowCustomConfig": "true",
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
	templatePath := filepath.Join(dir, "workflow_operational_resource.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}
