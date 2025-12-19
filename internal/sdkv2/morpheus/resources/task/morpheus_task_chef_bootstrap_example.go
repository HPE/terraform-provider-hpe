// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package task

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/testhelpers"
)

//go:generate go run ../../../../../cmd/render -out examples/resources/morpheus_task_chef_bootstrap/resource.tf morpheus_task_chef_bootstrap_resource.tf.tmpl Name "terraform_example_chef" Code "terraform_example_chef" Labels "\"demo\", \"terraform\"" ChefServerId "9" Environment "dev" RunList "role[web]" DataBagKey "test123" DataBagKeyPath "/etc/chef/databag_secret" NodeName "demonode" NodeAttributes "\"test\":\"demo\"" Retryable "true" RetryCount "1" RetryDelaySeconds "10" AllowCustomConfig "true" Visibility "public"

// RenderChefBootstrapConfig renders the task chef bootstrap
// resource configuration with the provided overrides applied to default values.
func RenderChefBootstrapConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":              "Example",
		"Code":              "example",
		"Labels":            `"demo", "terraform"`,
		"ChefServerId":      "9",
		"Environment":       "dev",
		"RunList":           "role[web]",
		"DataBagKey":        "test123",
		"DataBagKeyPath":    "/etc/chef/databag_secret",
		"NodeName":          "demonode",
		"NodeAttributes":    `"test":"demo"`,
		"Retryable":         "true",
		"RetryCount":        "1",
		"RetryDelaySeconds": "10",
		"AllowCustomConfig": "true",
		"Visibility":        "public",
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
	templatePath := filepath.Join(dir, "morpheus_task_chef_bootstrap_resource.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}
