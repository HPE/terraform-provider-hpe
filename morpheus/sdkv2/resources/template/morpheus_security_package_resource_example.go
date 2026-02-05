// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package template

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
)

//go:generate../../../../bin/render -out examples/resources/morpheus_security_package/resource.tf morpheus_security_package_resource.tf.tmpl Name "tf_example_security_package" Description "Terraform security package example" Labels "[\"demo\", \"terraform\"]" Enabled true Url "https://github.com/ComplianceAsCode/content/releases/download/v0.1.59/scap-security-guide-0.1.59.zip"

func RenderSecurityPackageConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":        "Example",
		"Description": "Terraform security package example",
		"Labels":      "[\"demo\", \"terraform\"]",
		"Enabled":     "true",
		"Url": "https://github.com/ComplianceAsCode/content/releases/download/" +
			"v0.1.59/scap-security-guide-0.1.59.zip",
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
	templatePath := filepath.Join(dir, "morpheus_security_package_resource.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}
