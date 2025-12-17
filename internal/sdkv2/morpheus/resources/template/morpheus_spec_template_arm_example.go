// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package template

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/testhelpers"
)

//go:generate go run ../../../../../cmd/render -out examples/resources/morpheus_spec_template_arm/resource_git.tf morpheus_spec_template_arm_resource_git.tf.tmpl Name tf-arm-spec-example-git SourceType repository RepositoryId 2 VersionRef main SpecPath ./test.json
//go:generate go run ../../../../../cmd/render -out examples/resources/morpheus_spec_template_arm/resource_local.tf morpheus_spec_template_arm_resource_local.tf.tmpl Name tf-arm-spec-example-local SourceType local
//go:generate go run ../../../../../cmd/render -out examples/resources/morpheus_spec_template_arm/resource_url.tf morpheus_spec_template_arm_resource_url.tf.tmpl Name tf-arm-spec-example-url SourceType url SpecPath http://example.com/spec.json

func RenderSpecTemplateArmLocalConfig(
	t *testing.T,
	name string,
	overrides map[string]string,
) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":       name,
		"SourceType": "local",
	}

	for key, value := range overrides {
		defaults[key] = value
	}

	// Get the directory where this source file is located
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("unable to get current file path")
	}
	dir := filepath.Dir(filename)
	templatePath := filepath.Join(dir, "morpheus_spec_template_arm_resource_local.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		"Name", defaults["Name"],
		"SourceType", defaults["SourceType"],
	)
}

func RenderSpecTemplateArmUrlConfig(
	t *testing.T,
	name string,
	overrides map[string]string,
) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":       name,
		"SourceType": "url",
		"SpecPath":   "http://example.com/spec.json",
	}

	for key, value := range overrides {
		defaults[key] = value
	}

	// Get the directory where this source file is located
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("unable to get current file path")
	}
	dir := filepath.Dir(filename)
	templatePath := filepath.Join(dir, "morpheus_spec_template_arm_resource_url.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		"Name", defaults["Name"],
		"SourceType", defaults["SourceType"],
		"SpecPath", defaults["SpecPath"],
	)
}

func RenderSpecTemplateArmGitConfig(
	t *testing.T,
	name string,
	overrides map[string]string,
) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":         name,
		"SourceType":   "repository",
		"RepositoryId": "2",
		"VersionRef":   "main",
		"SpecPath":     "./test.json",
	}

	for key, value := range overrides {
		defaults[key] = value
	}

	// Get the directory where this source file is located
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("unable to get current file path")
	}
	dir := filepath.Dir(filename)
	templatePath := filepath.Join(dir, "morpheus_spec_template_arm_resource_git.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		"Name", defaults["Name"],
		"SourceType", defaults["SourceType"],
		"RepositoryId", defaults["RepositoryId"],
		"VersionRef", defaults["VersionRef"],
		"SpecPath", defaults["SpecPath"],
	)
}
