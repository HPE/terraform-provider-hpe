// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package task_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/resources/task"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
)

func TestAccMorpheusTaskGroovyScriptGitExampleOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	t.Parallel()

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())

	// MORPH-16489: the groovy git task pins repository_id to a master-owned git
	// integration (id 1), so a sub-tenant stores 0 and the plan never
	// converges. Create a per-run git integration in the current tenant and
	// reference the discovered git repository id (not the integration id — the
	// task's file.repository.id is the SCM content-source id exposed via the
	// integration's computed repository_ids map). A public repository needs no
	// credentials.
	gitConfig := fmt.Sprintf(`
resource "hpe_morpheus_integration_git" "repo" {
  name           = %q
  url            = "https://github.com/gomorpheus/morpheus-cli.git"
  default_branch = "master"
}
`, name+"-repo")

	resourceConfig, err := task.RenderTaskGroovyScriptGitConfig(t, map[string]string{
		"Name":         name,
		"Code":         name,
		"RepositoryId": "one(values(hpe_morpheus_integration_git.repo.repository_ids))",
	})
	if err != nil {
		t.Fatal(err)
	}
	resourceConfig = gitConfig + resourceConfig

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_groovy_script.tfexample_groovy_git",
			"name",
			name,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_groovy_script.tfexample_groovy_git",
			"code",
			name,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_groovy_script.tfexample_groovy_git",
			"labels.0",
			"demo",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_groovy_script.tfexample_groovy_git",
			"labels.1",
			"terraform",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_groovy_script.tfexample_groovy_git",
			"source_type",
			"repository",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_groovy_script.tfexample_groovy_git",
			"result_type",
			"json",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_groovy_script.tfexample_groovy_git",
			"script_path",
			"example.groovy",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_groovy_script.tfexample_groovy_git",
			"version_ref",
			"master",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_groovy_script.tfexample_groovy_git",
			"retryable",
			"true",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_groovy_script.tfexample_groovy_git",
			"retry_count",
			"1",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_groovy_script.tfexample_groovy_git",
			"retry_delay_seconds",
			"10",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_groovy_script.tfexample_groovy_git",
			"allow_custom_config",
			"true",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// MORPH-16489: create the git integration first, on its own, so the
			// appliance scans the repository and populates repository_ids before
			// the task that references a discovered repository id is created. If
			// the integration and the task are created in the same apply, the
			// scan has not completed when the task is created and the task
			// stores repository_id = 0 (the drift the ticket describes).
			{
				Config:             providerConfig + gitConfig,
				ExpectNonEmptyPlan: false,
			},
			// Apply - repository_id references the per-run git integration's
			// discovered repository id, so the plan converges cleanly.
			{
				Config:             providerConfig + resourceConfig,
				Check:              checkFn,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}
