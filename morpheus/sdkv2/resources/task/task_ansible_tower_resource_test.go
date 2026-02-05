// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package task_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus"
	sdkv2morpheus "github.com/HPE/terraform-provider-hpe/morpheus/sdkv2"
	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/resources/task"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
)

func TestAccMorpheusTaskAnsibleTowerExampleOk(t *testing.T) {
	t.Parallel()

	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	// t.Skip("Skipping due to API error")
	// t.Skip("Skipping due to missing infrastructure in test environment")
	// t.Skip("Skipping due to missing resource implementation")
	// t.Skip("Skipping due to mismatch between Morpheus API and Terraform schema")

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())

	resourceConfig, err := task.RenderTaskAnsibleTowerConfig(t, map[string]string{
		"Name": name,
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

		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_ansible_tower.example",
			"ansible_tower_integration_id",
			"1",
		),

		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_ansible_tower.example",
			"ansible_tower_inventory_id",
			"5",
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

		resource.TestCheckResourceAttr(
			"hpe_morpheus_task_ansible_tower.example",
			"job_template_id",
			"3",
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
			"public",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), sdkv2morpheus.Provider()),
		Steps: []resource.TestStep{
			// Apply
			{
				Config:             providerConfig + resourceConfig,
				ExpectNonEmptyPlan: false,
				Check:              checkFn,
			},
			// Plan after apply
			{
				Config:             providerConfig + resourceConfig,
				ExpectNonEmptyPlan: false,
				PlanOnly:           true,
			},
		},
	})
}
