// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package cloud

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
)

//go:generate ../../../../../../bin/render -out examples/resources/morpheus_cloud/example.tf example.tf.tmpl Name "TestCloud" TenantId "1" GroupId "1" Code "aCode" Label "aLabel" ApplianceUrl "https://somewhere.com"
//go:generate ../../../../../../bin/render -out examples/resources/morpheus_cloud/example_generic.tf example_generic.tf.tmpl Name "TestCloud" TenantId "1" GroupId "1" Code "aCode" Label "aLabel" ApplianceUrl "https://somewhere.com"

func RenderCloudConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":         "TestCloud",
		"TenantId":     "1",
		"GroupId":      "1",
		"Code":         "testcloud",
		"Label":        "aLabel",
		"ApplianceUrl": "https://somewhere.com",
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
	templatePath := filepath.Join(dir, "example.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}

func RenderCloudGenericConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":         "TestCloud",
		"TenantId":     "1",
		"GroupId":      "1",
		"Code":         "testcloud",
		"Label":        "aLabel",
		"ApplianceUrl": "https://somewhere.com",
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
	templatePath := filepath.Join(dir, "example_generic.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}
