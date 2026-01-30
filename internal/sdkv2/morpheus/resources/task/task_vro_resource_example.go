// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package task

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/testhelpers"
)

//go:generate sh -c " ../../../../../bin/render -out examples/resources/morpheus_task_vro/resource.tf task_vro_resource.tf.tmpl Body '<<EOF\n{\n \"parameters\": [\n {\n \"name\": \"vmName\",\n \"type\": \"string\",\n \"value\": {\n \"string\": {\n \"value\": \"<%=instance.hostname%>\"\n }\n }\n }\n ]\n}\nEOF' Code 'tfexample-vro-task' ExecuteTarget 'local' Labels '[\"demo\", \"terraform\"]' Name 'tfexample vro-task' Retryable 'false' VroIntegrationId '1' VroWorkflowValue '1'"

// RenderTaskVroConfig generates a Terraform configuration for the tenant resource.
// It accepts optional overrides for field values. Default values are used if not overridden.
func RenderTaskVroConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	bodyDefault := "<<EOF\n{\n \"parameters\": [\n {\n \"name\": \"vmName\",\n" +
		" \"type\": \"string\",\n \"value\": {\n \"string\": {\n" +
		" \"value\": \"<%=instance.hostname%>\"\n }\n }\n }\n ]\n}\nEOF"

	defaults := map[string]string{
		"Body":             bodyDefault,
		"Code":             "tfexample-vro-task",
		"ExecuteTarget":    "local",
		"Labels":           "[\"demo\", \"terraform\"]",
		"Name":             "tfexample vro-task",
		"Retryable":        "false",
		"VroIntegrationId": "1",
		"VroWorkflowValue": "1",
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
	templatePath := filepath.Join(dir, "task_vro_resource.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}
