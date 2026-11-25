// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package integration_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/resources/integration"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
)

func TestAccMorpheusIntegrationAnsibleExampleOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.Ansible)

	t.Parallel()

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	repositoryURL, repositoryBranch := gitRepositoryFromEnv(t)

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())

	resourceConfig, err := integration.RenderIntegrationAnsibleConfig(t, map[string]string{
		"Name":          name,
		"Url":           gitRepositoryURLExpr,
		"DefaultBranch": gitRepositoryBranchExpr,
	})
	if err != nil {
		t.Fatal(err)
	}
	config := providerConfig + gitRepositoryVariables + resourceConfig

	const addr = "hpe_morpheus_integration_ansible.tf_example_ansible_integration"
	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(addr, "name", name),
		resource.TestCheckResourceAttr(addr, "enabled", "true"),
		resource.TestCheckResourceAttr(addr, "url", repositoryURL),
		resource.TestCheckResourceAttr(addr, "default_branch", repositoryBranch),
		resource.TestCheckResourceAttr(addr, "playbooks_path", "/"),
		resource.TestCheckResourceAttr(addr, "roles_path", "/roles"),
		resource.TestCheckResourceAttr(addr, "group_variables_path", "/vars"),
		resource.TestCheckResourceAttr(addr, "host_variables_path", "/vars"),
		resource.TestCheckResourceAttr(addr, "enable_ansible_galaxy_install", "true"),
		resource.TestCheckResourceAttr(addr, "enable_verbose_logging", "true"),
		resource.TestCheckResourceAttr(addr, "enable_agent_command_bus", "true"),
		resource.TestCheckResourceAttr(addr, "enable_git_caching", "true"),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Plan
			{
				Config:             config,
				ExpectNonEmptyPlan: true,
				Check:              checkFn,
				PlanOnly:           true,
			},
			// Apply
			{
				Config: config,
				Check:  checkFn,
			},
			// Plan after apply
			{
				Config:             config,
				ExpectNonEmptyPlan: false,
				Check:              checkFn,
				PlanOnly:           true,
			},
		},
	})
}
