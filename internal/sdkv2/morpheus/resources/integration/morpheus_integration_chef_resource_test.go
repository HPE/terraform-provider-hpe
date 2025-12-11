// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package integration_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/testhelpers"
)

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

	return testhelpers.RenderExample(t, "morpheus_integration_chef_resource.tf.tmpl", args...)
}

func TestAccMorpheusIntegrationChefExampleOk(t *testing.T) {
	t.Skip("Skipping due to lack of available resources to test against")

	t.Parallel()

	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())

	resourceConfig, err := RenderIntegrationChefConfig(t, map[string]string{
		"Name": name,
	})
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_integration_chef.tf_example_chef_integration",
			"name",
			name,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_integration_chef.tf_example_chef_integration",
			"enabled",
			"true",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_integration_chef.tf_example_chef_integration",
			"url",
			"https://chef.morpheusdata.com",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_integration_chef.tf_example_chef_integration",
			"version",
			"15.9.38",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_integration_chef.tf_example_chef_integration",
			"windows_version",
			"15.9.38",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_integration_chef.tf_example_chef_integration",
			"windows_msi_install_url",
			"https://packages.chef.io",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_integration_chef.tf_example_chef_integration",
			"organization",
			"morpheus",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_integration_chef.tf_example_chef_integration",
			"username",
			"admin",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_integration_chef.tf_example_chef_integration",
			"private_key",
			"EXAMPLEPRIVATEKEY",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_integration_chef.tf_example_chef_integration",
			"organization_validator_key",
			"EXAMPLEPRIVATEKEY",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Plan
			{
				Config:             providerConfig + resourceConfig,
				ExpectNonEmptyPlan: true,
				Check:              checkFn,
				PlanOnly:           true,
			},
			// Apply
			{
				Config: providerConfig + resourceConfig,
				Check:  checkFn,
			},
			// Plan after apply
			{
				Config:             providerConfig + resourceConfig,
				ExpectNonEmptyPlan: false,
				Check:              checkFn,
				PlanOnly:           true,
			},
		},
	})
}
