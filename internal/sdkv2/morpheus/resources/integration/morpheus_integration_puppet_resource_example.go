// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package integration

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/testhelpers"
)

//go:generate go run ../../../../../cmd/render -out examples/resources/morpheus_integration_puppet/resource.tf morpheus_integration_puppet_resource.tf.tmpl Name "tfexample puppet integration" Enabled true PuppetMasterHostname peserver01.morpheusdata.com AllowImmediateExecution true PuppetMasterSshUsername root PuppetMasterSshPassword password123

func RenderIntegrationPuppetConfig(t *testing.T, name string, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":                    name,
		"Enabled":                 "true",
		"PuppetMasterHostname":    "peserver01.morpheusdata.com",
		"AllowImmediateExecution": "true",
		"PuppetMasterSshUsername": "root",
		"PuppetMasterSshPassword": "password123",
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
	templatePath := filepath.Join(dir, "morpheus_integration_puppet_resource.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		"Name", defaults["Name"],
		"Enabled", defaults["Enabled"],
		"PuppetMasterHostname", defaults["PuppetMasterHostname"],
		"AllowImmediateExecution", defaults["AllowImmediateExecution"],
		"PuppetMasterSshUsername", defaults["PuppetMasterSshUsername"],
		"PuppetMasterSshPassword", defaults["PuppetMasterSshPassword"],
	)
}
