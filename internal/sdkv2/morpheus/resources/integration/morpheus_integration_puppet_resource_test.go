// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package integration_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/testhelpers"
)

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

	return testhelpers.RenderExample(t, "morpheus_integration_puppet_resource.tf.tmpl", args...)
}

func TestAccMorpheusIntegrationPuppetExampleOk(t *testing.T) {
	t.Skip("Skipping due to lack of available resources to test against")

	t.Parallel()

	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())

	resourceConfig, err := RenderMorpheusIntegrationPuppetConfig(t, name, nil)
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_integration_puppet.tf_example_puppet_integration",
			"name",
			name,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_integration_puppet.tf_example_puppet_integration",
			"enabled",
			"true",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_integration_puppet.tf_example_puppet_integration",
			"puppet_master_hostname",
			"peserver01.morpheusdata.com",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_integration_puppet.tf_example_puppet_integration",
			"allow_immediate_execution",
			"true",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_integration_puppet.tf_example_puppet_integration",
			"puppet_master_ssh_username",
			"root",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_integration_puppet.tf_example_puppet_integration",
			"puppet_master_ssh_password",
			"password123",
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
