// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package compute_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus"
	sdkv2morpheus "github.com/HPE/terraform-provider-hpe/morpheus/sdkv2"
	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/resources/compute"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
)

func TestMain(m *testing.M) {
	code := m.Run()

	testhelpers.WriteMergedResults()

	os.Exit(code)
}

func TestAccResourcePoolGroupExampleOk(t *testing.T) {
	t.Parallel()

	defer testhelpers.RecordResult(t)

	t.Skip("Skipping due to missing infrastructure in test environment")

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())

	resourceConfig, err := compute.RenderResourcePoolGroupConfig(t, map[string]string{
		"Name": name,
	})
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_resource_pool_group.example",
			"name",
			name,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_resource_pool_group.example",
			"description",
			"TFExample Resource Pool Group",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_resource_pool_group.example",
			"mode",
			"roundRobin",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_resource_pool_group.example",
			"resource_pool_ids.0",
			"1",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_resource_pool_group.example",
			"resource_pool_ids.1",
			"2",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_resource_pool_group.example",
			"resource_pool_ids.2",
			"3",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_resource_pool_group.example",
			"all_group_access",
			"true",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_resource_pool_group.example",
			"group_access.0.group_id",
			"2",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_resource_pool_group.example",
			"group_access.0.default",
			"true",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_resource_pool_group.example",
			"visibility",
			"public",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_resource_pool_group.example",
			"tenant_ids.0",
			"1",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_resource_pool_group.example",
			"tenant_ids.1",
			"2",
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
