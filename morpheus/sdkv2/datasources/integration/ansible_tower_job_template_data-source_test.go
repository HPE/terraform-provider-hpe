// (C) Copyright 2025-2026 Hewlett Packard Enterprise Development LP

package integration_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	sdkv2morpheus "github.com/HPE/terraform-provider-hpe/morpheus/sdkv2"
	dsintegration "github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/datasources/integration"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

// TestAccMorpheusDataSourceAnsibleTowerJobTemplateExampleOk looks a job
// template up by name on the Ansible Tower integration the environment names
// (testhelpers/tower_fixture.go); the integration is resolved by name too.
func TestAccMorpheusDataSourceAnsibleTowerJobTemplateExampleOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.AnsibleTower)

	integrationName := testhelpers.AnsibleTowerIntegrationName(t)
	jobTemplateName := testhelpers.AnsibleTowerJobTemplateName(t)

	t.Parallel()

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	dependenciesConfig := testhelpers.AnsibleTowerIntegrationBlock(integrationName)

	datasourceConfig, err := dsintegration.RenderAnsibleTowerJobTemplateConfig(t, map[string]string{
		"Name":                      testhelpers.HCLQuote(jobTemplateName),
		"AnsibleTowerIntegrationId": testhelpers.AnsibleTowerIntegrationIDRef,
	})
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_ansible_tower_job_template.example",
			"name",
			jobTemplateName,
		),
		resource.TestCheckResourceAttrSet(
			"data.hpe_morpheus_ansible_tower_job_template.example",
			"id",
		),
		resource.TestCheckResourceAttrPair(
			"data.hpe_morpheus_ansible_tower_job_template.example", "ansible_tower_integration_id",
			"data.hpe_morpheus_integration.ansible_tower", "id",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), sdkv2morpheus.Provider()),
		Steps: []resource.TestStep{
			{
				Config:             providerConfig + dependenciesConfig + datasourceConfig,
				ExpectNonEmptyPlan: false,
				Check:              checkFn,
			},
		},
	})
}
