// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package blueprint

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/testhelpers"
)

//go:generate go run ../../../../../cmd/render -out examples/resources/app_blueprint_arm/app_blueprint_arm_resource_git.tf app_blueprint_arm_resource_git.tf.tmpl Category "armtemplates" CloudInitEnabled "true" Description "example arm app blueprint" InstallAgent "true" IntegrationId "3" Name "example_app_arm_blueprint_git" OsType "linux" RepositoryId "1" SourceType "repository" VersionRef "main" WorkingPath "./test"

// RenderAppBlueprintArmConfig generates a Terraform configuration for the tenant resource.
// It accepts optional overrides for field values. Default values are used if not overridden.
func RenderAppBlueprintArmConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Category":         "armtemplates",
		"CloudInitEnabled": "true",
		"Description":      "example arm app blueprint",
		"InstallAgent":     "true",
		"IntegrationId":    "3",
		"Name":             "example_app_arm_blueprint_git",
		"OsType":           "linux",
		"RepositoryId":     "1",
		"SourceType":       "repository",
		"VersionRef":       "main",
		"WorkingPath":      "./test",
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
	templatePath := filepath.Join(dir, "app_blueprint_arm_resource_git.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}
