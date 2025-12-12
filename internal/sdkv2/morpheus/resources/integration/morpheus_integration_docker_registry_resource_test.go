// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package integration_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/testhelpers"
)

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

	return testhelpers.RenderExample(
		t,
		"morpheus_integration_docker_registry_resource.tf.tmpl",
		"Name", defaults["Name"],
		"Enabled", defaults["Enabled"],
		"Url", defaults["Url"],
		"Username", defaults["Username"],
		"Password", defaults["Password"],
	)
}

func TestAccMorpheusIntegrationDockerRegistryExampleOk(t *testing.T) {
	t.Parallel()

	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())

	resourceConfig, err := RenderMorpheusIntegrationDockerRegistryConfig(t, map[string]string{
		"Name": name,
	})
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_integration_docker_registry.tf_example_docker_registry_integration",
			"name",
			name,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_integration_docker_registry.tf_example_docker_registry_integration",
			"enabled",
			"true",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_integration_docker_registry.tf_example_docker_registry_integration",
			"url",
			"https://index.docker.io/v1/",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_integration_docker_registry.tf_example_docker_registry_integration",
			"username",
			"admin",
		),
		resource.TestCheckResourceAttrSet(
			"hpe_morpheus_integration_docker_registry.tf_example_docker_registry_integration",
			"password",
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
