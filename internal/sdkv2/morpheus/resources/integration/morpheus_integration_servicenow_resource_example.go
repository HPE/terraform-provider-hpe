// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package integration

import (
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"

	"fmt"
	"path/filepath"
	"runtime"
	"testing"
	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/testhelpers"
)

//go:generate go run ../../../../../cmd/render -out examples/resources/morpheus_integration_servicenow/resource.tf morpheus_integration_servicenow_resource.tf.tmpl Name "terraform servicenow integration" Enabled true Url "https://servicenowprod.service-now.com" Username "my-snow-username" Password "my-snow-password" DefaultCmdbBusinessClass "demo"

// RenderMorpheusIntegrationVroConfig generates a Terraform configuration for the VRO integration
// resource. It accepts an optional map of field overrides to customize the default values.
func RenderMorpheusIntegrationVroConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	// Default field values
	defaults := map[string]string{
		"Name":     acctest.RandomWithPrefix(t.Name()),
		"Enabled":  "true",
		"Url":      "https://myvro/vco/api",
		"Username": "my-vro-username",
		"Password": "my-vro-password",
		"AuthType": "basic",
		"Tenant":   "vsphere.local",
	}

	// Apply overrides to defaults
	for key, value := range overrides {
		defaults[key] = value
	}

	// Get the directory where this source file is located
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("unable to get current file path")
	}
	dir := filepath.Dir(filename)
	templatePath := filepath.Join(dir, "morpheus_integration_vro_resource.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		"Name", defaults["Name"],
		"Enabled", defaults["Enabled"],
		"Url", defaults["Url"],
		"Username", defaults["Username"],
		"Password", defaults["Password"],
		"AuthType", defaults["AuthType"],
		"Tenant", defaults["Tenant"],
	)
}

func RenderIntegrationServicenowConfig(
	t *testing.T,
	name string,
	overrides map[string]string,
) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":                     name,
		"Enabled":                  "true",
		"Url":                      "https://servicenowprod.service-now.com",
		"Username":                 "my-snow-username",
		"Password":                 "my-snow-password",
		"DefaultCmdbBusinessClass": "demo",
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
	templatePath := filepath.Join(dir, "morpheus_integration_servicenow_resource.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		"Name", defaults["Name"],
		"Enabled", defaults["Enabled"],
		"Url", defaults["Url"],
		"Username", defaults["Username"],
		"Password", defaults["Password"],
		"DefaultCmdbBusinessClass", defaults["DefaultCmdbBusinessClass"],
	)
}

func RenderMorpheusIntegrationPuppetConfig(t *testing.T, name string, overrides map[string]string) (string, error) {
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

	args := []string{}
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
		templatePath, args...)
}

// RenderMorpheusIntegrationAnsibleTowerConfig renders the Ansible Tower integration resource configuration
// with default values that can be overridden via the overrides parameter.
func RenderMorpheusIntegrationAnsibleTowerConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":     "tf_test_ansible_tower",
		"Enabled":  "true",
		"Url":      "https://ansibletower01.morpheusdata.com",
		"Username": "admin",
		"Password": "password123",
	}

	// Override defaults with provided values
	for key, value := range overrides {
		defaults[key] = value
	}

	// Get the directory where this source file is located
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("unable to get current file path")
	}
	dir := filepath.Dir(filename)
	templatePath := filepath.Join(dir, "morpheus_integration_ansible_tower_resource.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		"Name", defaults["Name"],
		"Enabled", defaults["Enabled"],
		"Url", defaults["Url"],
		"Username", defaults["Username"],
		"Password", defaults["Password"],
	)
}

// RenderMorpheusIntegrationDockerRegistryConfig renders the Docker Registry integration
// resource configuration with the provided field overrides. Default values are used for any
// fields not specified.
func RenderMorpheusIntegrationDockerRegistryConfig(
	t *testing.T,
	overrides map[string]string,
) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":     acctest.RandomWithPrefix(t.Name()),
		"Enabled":  "true",
		"Url":      "https://index.docker.io/v1/",
		"Username": "admin",
		"Password": "password123",
	}

	// Apply overrides
	for key, value := range overrides {
		defaults[key] = value
	}

	// Get the directory where this source file is located
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("unable to get current file path")
	}
	dir := filepath.Dir(filename)
	templatePath := filepath.Join(dir, "morpheus_integration_docker_registry_resource.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath,
		"Name", defaults["Name"],
		"Enabled", defaults["Enabled"],
		"Url", defaults["Url"],
		"Username", defaults["Username"],
		"Password", defaults["Password"],
	)
}

// RenderIntegrationChefConfig generates a Terraform configuration for the Chef integration resource
// using default values that can be overridden via the overrides map.
func RenderIntegrationChefConfig(t *testing.T, overrides map[string]string) (string, error) {
	t.Helper()

	defaults := map[string]string{
		"Name":                     "test-chef-integration",
		"Enabled":                  "true",
		"Url":                      "https://chef.morpheusdata.com",
		"Version":                  "15.9.38",
		"WindowsVersion":           "15.9.38",
		"WindowsMsiInstallUrl":     "https://packages.chef.io",
		"Organization":             "morpheus",
		"Username":                 "admin",
		"PrivateKey":               "EXAMPLEPRIVATEKEY",
		"OrganizationValidatorKey": "EXAMPLEPRIVATEKEY",
	}

	for key, value := range overrides {
		defaults[key] = value
	}

	args := make([]string, 0, len(defaults)*2)
	for key, value := range defaults {
		args = append(args, key, value)
	}

	// Get the directory where this source file is located
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("unable to get current file path")
	}
	dir := filepath.Dir(filename)
	templatePath := filepath.Join(dir, "morpheus_integration_chef_resource.tf.tmpl")

	return testhelpers.RenderExample(
		t,
		templatePath, args...)
}

