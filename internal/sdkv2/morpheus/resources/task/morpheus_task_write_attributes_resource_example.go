// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package task

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/testhelpers"
)

//go:generate ../../../../../bin/render -out examples/resources/morpheus_task_write_attributes/resource.tf morpheus_task_write_attributes_resource.tf.tmpl Name tfexample_write_attributes Code tfexample_write_attributes Label1 demo Label2 terraform Attributes {"demo":"test"} Retryable true RetryCount 1 RetryDelaySeconds 10 AllowCustomConfig true

// RenderTaskWriteAttributesConfig generates a Terraform configuration for testing
// the task_write_attributes resource. It accepts overrides to customize field values.
func RenderTaskWriteAttributesConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":              "Example",
		"Code":              "example",
		"Label1":            "demo",
		"Label2":            "terraform",
		"Attributes":        `{"demo":"test"}`,
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
	templatePath := filepath.Join(dir, "morpheus_task_write_attributes_resource.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}
