// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package template

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/testhelpers"
)

//go:generate go run ../../../../../cmd/render -out examples/resources/morpheus_spec_template_terraform/resource_url.tf morpheus_spec_template_terraform_resource_url.tf.tmpl Name "tf-terraform-spec-example-url" SourceType "url" SpecPath "http://example.com/spec.tf"
//go:generate go run ../../../../../cmd/render -out examples/resources/morpheus_spec_template_terraform/resource_local.tf morpheus_spec_template_terraform_resource_local.tf.tmpl Name "tf-terraform-spec-example-local" SourceType "local" SpecContent "resource \"aws_instance\" \"instance_1\" {\n  ami           = \"ami-0b91a410940e82c54\"\n  instance_type = \"t2.micro\"\n}"
//go:generate go run ../../../../../cmd/render -out examples/resources/morpheus_spec_template_terraform/resource_git.tf morpheus_spec_template_terraform_resource_git.tf.tmpl Name "tf-terraform-spec-example-git" SourceType "repository" RepositoryId "2" VersionRef "main" SpecPath "Instance Types/Terraform/CloudResource/aws/vpc.tf"

// RenderSpecTemplateTerraformLocalConfig renders the Terraform config for
// spec_template_terraform_resource_local tests
func RenderSpecTemplateTerraformLocalConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":       "Example",
		"SourceType": "local",
		"SpecContent": `resource "aws_instance" "instance_1" {
  ami           = "ami-0b91a410940e82c54"
  instance_type = "t2.micro"
}`,
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
	templatePath := filepath.Join(dir, "morpheus_spec_template_terraform_resource_local.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}

// RenderSpecTemplateTerraformGitConfig renders the Terraform config for
// spec_template_terraform_resource_git tests
func RenderSpecTemplateTerraformGitConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":         "Example",
		"SourceType":   "repository",
		"RepositoryId": "2",
		"VersionRef":   "main",
		"SpecPath":     "Instance Types/Terraform/CloudResource/aws/vpc.tf",
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
	templatePath := filepath.Join(dir, "morpheus_spec_template_terraform_resource_git.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}

// RenderSpecTemplateTerraformUrlConfig renders the Terraform config for
// spec_template_terraform_resource_url tests
func RenderSpecTemplateTerraformUrlConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":       "Example",
		"SourceType": "url",
		"SpecPath":   "http://example.com/spec.tf",
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
	templatePath := filepath.Join(dir, "morpheus_spec_template_terraform_resource_url.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}
