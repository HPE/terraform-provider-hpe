// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package task

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
)

//go:generate ../../../../../bin/render -out examples/resources/morpheus_task_restart/resource.tf morpheus_task_restart_resource.tf.tmpl Name tfexample_restart Code tfexample_restart Labels "[\"demo\", \"terraform\"]" Retryable true RetryCount 1 RetryDelaySeconds 10 AllowCustomConfig true

// RenderTaskRestartConfig renders the task restart resource configuration
// with the provided name and field overrides.
func RenderTaskRestartConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	// Default values
	defaults := map[string]string{
		"Name":              "Example",
		"Code":              "tfexample_restart",
		"Labels":            `["demo", "terraform"]`,
		"Retryable":         "true",
		"RetryCount":        "1",
		"RetryDelaySeconds": "10",
		"AllowCustomConfig": "true",
	}

	// Apply overrides
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
	templatePath := filepath.Join(dir, "morpheus_task_restart_resource.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}
