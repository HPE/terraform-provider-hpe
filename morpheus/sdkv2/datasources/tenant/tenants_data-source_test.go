// (C) Copyright 2025-2026 Hewlett Packard Enterprise Development LP

package tenant_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus"
	sdkv2morpheus "github.com/HPE/terraform-provider-hpe/morpheus/sdkv2"
	dstenant "github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/datasources/tenant"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
)

func TestAccMorpheusDataSourceTenantsExampleOk(t *testing.T) {
	t.Parallel()

	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	var dependenciesConfig string

	// We're always guaranteed to find at least one tenant (the Master Tenant),
	// so we don't need to create anything for this test.
	datasourceConfig, err := dstenant.RenderTenantsConfig(t, map[string]string{
		"Values": "\".*\"", // find all tenants
	})
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		// just check that we get at least one ID with data source.
		resource.TestCheckResourceAttrSet(
			"data.hpe_morpheus_tenants.example",
			"ids.#",
		),
		// ...and make sure that these have been set correctly.
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_tenants.example",
			"filter.0.name",
			"name",
		),

		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_tenants.example",
			"sort_ascending",
			"true",
		),

		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_tenants.example",
			"filter.0.values.0",
			".*",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), sdkv2morpheus.Provider()),
		Steps: []resource.TestStep{
			{
				Config:             providerConfig + dependenciesConfig + datasourceConfig,
				ExpectNonEmptyPlan: false,
				Check:              checkFn,
			},
		},
	})
}
