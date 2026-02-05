// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package task

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
)

//go:generate sh -c "../../../../bin/render -out examples/resources/morpheus_task_library_script/resource.tf task_library_script_resource.tf.tmpl AllowCustomConfig 'true' Code 'tf-example-library-script-task' ExecuteTarget 'resource' Labels '[\"demo\", \"library\", \"terraform\"]' Name 'Example Terraform Library Script Task' RetryCount '1' RetryDelaySeconds '10' Retryable 'true' ScriptTemplate 'My script template' ScriptTemplateId '1'"

// RenderTaskLibraryScriptConfig generates a Terraform configuration for the tenant resource.
// It accepts optional overrides for field values. Default values are used if not overridden.
func RenderTaskLibraryScriptConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"AllowCustomConfig": "true",
		"Code":              "tf-example-library-script-task",
		"ExecuteTarget":     "resource",
		"Labels":            "[\"demo\", \"library\", \"terraform\"]",
		"Name":              "Example Terraform Library Script Task",
		"RetryCount":        "1",
		"RetryDelaySeconds": "10",
		"Retryable":         "true",
		"ScriptTemplate":    "My script template",
		"ScriptTemplateId":  "1",
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
	templatePath := filepath.Join(dir, "task_library_script_resource.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}
