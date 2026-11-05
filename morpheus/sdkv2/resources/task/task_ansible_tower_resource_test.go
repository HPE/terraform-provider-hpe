// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package task_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	sdkv2morpheus "github.com/HPE/terraform-provider-hpe/morpheus/sdkv2"
	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/resources/task"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

// TestAccMorpheusTaskAnsibleTowerExampleOk creates a Tower task against the
// Ansible Tower integration the environment names (testhelpers/tower_fixture.go):
// the integration, the inventory and the job template are all resolved by name
// through data sources, and the task takes their ids from those.
func TestAccMorpheusTaskAnsibleTowerExampleOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.Ansible, capabilities.AnsibleTower)

	integrationName := testhelpers.AnsibleTowerIntegrationName(t)
	inventoryName := testhelpers.AnsibleTowerInventoryName(t)
	jobTemplateName := testhelpers.AnsibleTowerJobTemplateName(t)

	t.Parallel()

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()
	visibility := testhelpers.TenantVisibility(t)

	name := acctest.RandomWithPrefix(t.Name())

	dependenciesConfig := testhelpers.AnsibleTowerIntegrationBlock(integrationName) + fmt.Sprintf(`
data "hpe_morpheus_ansible_tower_inventory" "example" {
  name                         = %s
  ansible_tower_integration_id = %s
}

data "hpe_morpheus_ansible_tower_job_template" "example" {
  name                         = %s
  ansible_tower_integration_id = %s
}
`, testhelpers.HCLQuote(inventoryName), testhelpers.AnsibleTowerIntegrationIDRef,
		testhelpers.HCLQuote(jobTemplateName), testhelpers.AnsibleTowerIntegrationIDRef)

	resourceConfig, err := task.RenderTaskAnsibleTowerConfig(t, map[string]string{
		"Name":                      name,
		"Visibility":                visibility,
		"AnsibleTowerIntegrationId": testhelpers.AnsibleTowerIntegrationIDRef,
		"AnsibleTowerInventoryId":   "data.hpe_morpheus_ansible_tower_inventory.example.id",
		"JobTemplateId":             "data.hpe_morpheus_ansible_tower_job_template.example.id",
	})
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_ansible_tower.example",
			"allow_custom_config",
			"true",
		),

		resource.TestCheckResourceAttrPair(
			"hpe_morpheus_task_ansible_tower.example", "ansible_tower_integration_id",
			"data.hpe_morpheus_integration.ansible_tower", "id",
		),

		resource.TestCheckResourceAttrPair(
			"hpe_morpheus_task_ansible_tower.example", "ansible_tower_inventory_id",
			"data.hpe_morpheus_ansible_tower_inventory.example", "id",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_ansible_tower.example",
			"code",
			"tfexample-ansible-tower-task",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_ansible_tower.example",
			"execute_mode",
			"executeAll",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_ansible_tower.example",
			"execute_target",
			"local",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_ansible_tower.example",
			"group",
			"demo",
		),

		resource.TestCheckResourceAttrPair(
			"hpe_morpheus_task_ansible_tower.example", "job_template_id",
			"data.hpe_morpheus_ansible_tower_job_template.example", "id",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_ansible_tower.example",
			"labels.#",
			"2",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_ansible_tower.example",
			"name",
			name,
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_ansible_tower.example",
			"retry_count",
			"5",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_ansible_tower.example",
			"retry_delay_seconds",
			"10",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_ansible_tower.example",
			"retryable",
			"true",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_ansible_tower.example",
			"scm_override",
			"main",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_ansible_tower.example",
			"visibility",
			visibility,
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), sdkv2morpheus.Provider()),
		Steps: []resource.TestStep{
			// Apply
			{
				Config:             providerConfig + dependenciesConfig + resourceConfig,
				ExpectNonEmptyPlan: false,
				Check:              checkFn,
			},
			// Plan after apply
			{
				Config:             providerConfig + dependenciesConfig + resourceConfig,
				ExpectNonEmptyPlan: false,
				PlanOnly:           true,
			},
		},
	})
}
