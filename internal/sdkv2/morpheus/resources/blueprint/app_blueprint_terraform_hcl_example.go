// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package blueprint

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/testhelpers"
)

//go:generate go run ../../../../../cmd/render -out examples/resources/app_blueprint_terraform/app_blueprint_terraform_resource_hcl.tf app_blueprint_terraform_resource_hcl.tf.tmpl BlueprintContent "..." Category "terraformdemo" Description "testing terraform" Name "tfappbluedemo" SourceType "hcl" TerraformOptions "-var 'foo=bar'" TerraformVersion "1.1.1" TfvarSecret "tfvars/rdsdemo-secrets"

// RenderAppBlueprintTerraformHclConfig generates a Terraform configuration for the tenant resource.
// It accepts optional overrides for field values. Default values are used if not overridden.

func RenderAppBlueprintTerraformHclConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"BlueprintContent": `variable "master_username" {
  type = string
}

variable "master_password" {
  type      = string
  sensitive = true
}

variable "engine_version" {
  type = string
}

variable "instance_class" {
  type = string
}

resource "local_file" "foo" {
    content  = "foo!"
    filename = "${path.module}/foo.bar"
}`,
		"Category":         "terraformdemo",
		"Description":      "testing terraform",
		"Name":             "tfappbluedemo",
		"SourceType":       "hcl",
		"TerraformOptions": "-var 'foo=bar'",
		"TerraformVersion": "1.1.1",
		"TfvarSecret":      "tfvars/rdsdemo-secrets",
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
	templatePath := filepath.Join(dir, "app_blueprint_terraform_resource_hcl.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}
