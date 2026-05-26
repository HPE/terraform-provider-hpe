// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package library_layout

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
)

//go:generate ../../../../bin/render -out examples/resources/hpe_morpheus_library_layout/example.tf example.tf.tmpl InstanceTypeId "1" Name "Example Layout" InstanceVersion "1.0" ProvisionTypeCode "docker"

func RenderLibraryLayoutConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"InstanceTypeId":    "1",
		"Name":              "Example Layout",
		"InstanceVersion":   "1.0",
		"ProvisionTypeCode": "docker",
	}

	for key, value := range overrides {
		defaults[key] = value
	}

	var args []string
	for key, value := range defaults {
		args = append(args, key, value)
	}

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("unable to get current file path")
	}
	dir := filepath.Dir(filename)
	templatePath := filepath.Join(dir, "example.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}
