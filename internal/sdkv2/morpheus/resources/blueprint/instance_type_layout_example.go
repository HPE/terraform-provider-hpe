// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package blueprint

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/testhelpers"
)

//go:generate go run ../../../../../cmd/render -out examples/resources/morpheus_instance_type_layout/resource.tf instance_type_layout_resource.tf.tmpl InstanceTypeId "data.hpe_morpheus_instance_type.example.id" Labels "[\"demo\", \"layout\", \"terraform\"]" Name "todo_app_frontend" Technology "vmware" Version "1.0"

// RenderInstanceTypeLayoutConfig generates a Terraform configuration for the tenant resource.
// It accepts optional overrides for field values. Default values are used if not overridden.

func RenderInstanceTypeLayoutConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"InstanceTypeId": "data.hpe_morpheus_instance_type.example.id",
		"Labels":         "[\"demo\", \"layout\", \"terraform\"]",
		"Name":           "todo_app_frontend",
		"Technology":     "vmware",
		"Version":        "1.0",
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
	templatePath := filepath.Join(dir, "instance_type_layout_resource.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}
