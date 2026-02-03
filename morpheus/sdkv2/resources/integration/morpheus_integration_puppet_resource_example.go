// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package integration

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
)

//go:generate ../../../../../bin/render -out examples/resources/morpheus_integration_puppet/resource.tf morpheus_integration_puppet_resource.tf.tmpl Name "tfexample puppet integration" Enabled true PuppetMasterHostname peserver01.morpheusdata.com AllowImmediateExecution true PuppetMasterSshUsername root PuppetMasterSshPassword password123

func RenderIntegrationPuppetConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":                    "Example",
		"Enabled":                 "true",
		"PuppetMasterHostname":    "peserver01.morpheusdata.com",
		"AllowImmediateExecution": "true",
		"PuppetMasterSshUsername": "root",
		"PuppetMasterSshPassword": "password123",
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
	templatePath := filepath.Join(dir, "morpheus_integration_puppet_resource.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		args...,
	)
}
