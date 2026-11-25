// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package integration

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
)

//go:generate ../../../../bin/render -out examples/resources/morpheus_integration_ansible/resource.tf morpheus_integration_ansible_resource.tf.tmpl Name "tfexample ansible integration" Enabled true Url "\"https://github.com/gomorpheus/morpheus-automation-examples.git\"" DefaultBranch "\"main\"" PlaybooksPath "/" RolesPath "/roles" GroupVariablesPath "/vars" HostVariablesPath "/vars" EnableAnsibleGalaxyInstall true EnableVerboseLogging true EnableAgentCommandBus true EnableGitCaching true

// RenderIntegrationAnsibleConfig generates a Terraform configuration for the
// Ansible integration resource using default values that can be overridden
// via the overrides map. The defaults point at a public repository over
// HTTPS, so no credentials or key pair are needed. Url and DefaultBranch —
// the values the acceptance test takes from the environment — are rendered
// as HCL expressions: their defaults carry their own quotes, and the test
// passes a Terraform variable reference instead of the value, so a
// repository URL or branch name is never spliced into HCL text.
func RenderIntegrationAnsibleConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":                       "Example",
		"Enabled":                    "true",
		"Url":                        `"https://github.com/gomorpheus/morpheus-automation-examples.git"`,
		"DefaultBranch":              `"main"`,
		"PlaybooksPath":              "/",
		"RolesPath":                  "/roles",
		"GroupVariablesPath":         "/vars",
		"HostVariablesPath":          "/vars",
		"EnableAnsibleGalaxyInstall": "true",
		"EnableVerboseLogging":       "true",
		"EnableAgentCommandBus":      "true",
		"EnableGitCaching":           "true",
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
	templatePath := filepath.Join(dir, "morpheus_integration_ansible_resource.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}
