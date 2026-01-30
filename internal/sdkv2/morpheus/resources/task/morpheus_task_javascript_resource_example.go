// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package task

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/testhelpers"
)

//go:generate ../../../../../bin/render -out examples/resources/morpheus_task_javascript/resource.tf morpheus_task_javascript_resource.tf.tmpl Name tfexample_javascript Code tfexample_javascript Labels ["demo","terraform"] ScriptContent console.log("testing") Retryable true RetryCount 1 RetryDelaySeconds 10 AllowCustomConfig true

// RenderTaskJavascriptConfig generates a terraform configuration string
// for task javascript resource.
// It accepts a name and a map to override default field values.
func RenderTaskJavascriptConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":              "Example",
		"Code":              "example",
		"Labels":            `["demo","terraform"]`,
		"ScriptContent":     `console.log("testing")`,
		"Retryable":         "true",
		"RetryCount":        "1",
		"RetryDelaySeconds": "10",
		"AllowCustomConfig": "true",
	}

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
	templatePath := filepath.Join(dir, "morpheus_task_javascript_resource.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}
