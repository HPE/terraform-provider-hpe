// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package workflow

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/testhelpers"
)

//go:generate ../../../../../bin/render -out examples/resources/morpheus_workflow_provisioning/resource.tf workflow_provisioning_resource.tf.tmpl Name "tf_example_provisioning_workflow" Description "Terraform provisioning workflow example" Labels "[\"demo\", \"terraform\"]" Platform "all" Visibility "private"

// RenderWorkflowProvisioningConfig generates a Terraform configuration for the workflow provisioning resource.
// It accepts optional overrides for field values. Default values are used if not overridden.
func RenderWorkflowProvisioningConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":        "tf_example_provisioning_workflow",
		"Description": "Terraform provisioning workflow example",
		"Labels":      "[\"demo\", \"terraform\"]",
		"Platform":    "all",
		"Visibility":  "private",
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
	templatePath := filepath.Join(dir, "workflow_provisioning_resource.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}
