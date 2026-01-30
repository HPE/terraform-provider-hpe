// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package template

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/testhelpers"
)

//go:generate ../../../../../bin/render -out examples/resources/morpheus_cluster_package/resource.tf cluster_package_resource.tf.tmpl Name tf_example_cluster_package Code tf-example-cluster-package Description "Terraform example cluster package" PackageVersion 1.2.3 Type apps PackageType example Enabled true RepeatInstall true SpecTemplateIds [1,2]

// RenderClusterPackageConfig generates a test configuration for cluster package resource.
// It accepts a name and a map of field overrides to customize the default values.
func RenderClusterPackageConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":            "Example",
		"Code":            "tf-example-cluster-package",
		"Description":     "Terraform example cluster package",
		"PackageVersion":  "1.2.3",
		"Type":            "apps",
		"PackageType":     "example",
		"Enabled":         "true",
		"RepeatInstall":   "true",
		"SpecTemplateIds": "[1,2]",
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
	templatePath := filepath.Join(dir, "cluster_package_resource.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}
