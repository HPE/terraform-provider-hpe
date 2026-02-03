// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package template

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
)

//go:generate ../../../../../bin/render -out examples/resources/morpheus_spec_template_arm/resource_git.tf morpheus_spec_template_arm_resource_git.tf.tmpl Name tf-arm-spec-example-git SourceType repository RepositoryId 2 VersionRef main SpecPath ./test.json
//go:generate ../../../../../bin/render -out examples/resources/morpheus_spec_template_arm/resource_local.tf morpheus_spec_template_arm_resource_local.tf.tmpl Name tf-arm-spec-example-local SourceType local
//go:generate ../../../../../bin/render -out examples/resources/morpheus_spec_template_arm/resource_url.tf morpheus_spec_template_arm_resource_url.tf.tmpl Name tf-arm-spec-example-url SourceType url SpecPath http://example.com/spec.json

func RenderSpecTemplateArmLocalConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":       "Example",
		"SourceType": "local",
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
	templatePath := filepath.Join(dir, "morpheus_spec_template_arm_resource_local.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}

func RenderSpecTemplateArmUrlConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":       "Example",
		"SourceType": "url",
		"SpecPath":   "http://example.com/spec.json",
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
	templatePath := filepath.Join(dir, "morpheus_spec_template_arm_resource_url.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}

func RenderSpecTemplateArmGitConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":         "Example",
		"SourceType":   "repository",
		"RepositoryId": "2",
		"VersionRef":   "main",
		"SpecPath":     "./test.json",
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
	templatePath := filepath.Join(dir, "morpheus_spec_template_arm_resource_git.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}
