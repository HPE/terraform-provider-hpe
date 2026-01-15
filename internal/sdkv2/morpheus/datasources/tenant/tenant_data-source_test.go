// (C) Copyright 2025-2026 Hewlett Packard Enterprise Development LP

package tenant_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus"
	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/resources/role"
	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/testhelpers"
	sdkv2morpheus "github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus"
	dstenant "github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus/datasources/tenant"
	"github.com/HPE/terraform-provider-hpe/internal/sdkv2/morpheus/resources/tenant"
)

func TestAccMorpheusDataSourceTenantExampleOk(t *testing.T) {
	t.Parallel()

	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())

	var dependenciesConfig string

	// Need a role to create a tenant.
	// Although they're guarantted to exist,
	// we want to avoid using the admin roles provided by the Morpheus appliance.
	if currentDependency, err := role.RenderRoleTenantConfig(t, map[string]string{
		"Name": name,
	}); err != nil {
		t.Fatal(err)
	} else {
		dependenciesConfig += currentDependency
	}

	// Create a tenant as a dependency so we can test searching by name.
	// The name of the master tenant is not guaranteed to be the same across
	// different Morpheus appliances, so we prefer to create one for testing.
	if currentDependency, err := tenant.RenderTenantConfig(t, map[string]string{
		"Name":       name,
		"BaseRoleId": "resource.hpe_morpheus_role.example.id",
	}); err != nil {
		t.Fatal(err)
	} else {
		dependenciesConfig += currentDependency
	}

	datasourceConfig, err := dstenant.RenderTenantConfig(t, map[string]string{
		"Name": "resource.hpe_morpheus_tenant.tf_example_tenant.name",
	})
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_tenant.example",
			"name",
			name,
		),
		// check if the tenant data source has read
		// the same values as the resource has in state.
		resource.TestCheckResourceAttrPair(
			"data.hpe_morpheus_tenant.example",
			"id",
			"hpe_morpheus_tenant.tf_example_tenant",
			"id",
		),
		resource.TestCheckResourceAttrPair(
			"data.hpe_morpheus_tenant.example",
			"account_number",
			"hpe_morpheus_tenant.tf_example_tenant",
			"account_number",
		),
		resource.TestCheckResourceAttrPair(
			"data.hpe_morpheus_tenant.example",
			"account_name",
			"hpe_morpheus_tenant.tf_example_tenant",
			"account_name",
		),
		resource.TestCheckResourceAttrPair(
			"data.hpe_morpheus_tenant.example",
			"customer_number",
			"hpe_morpheus_tenant.tf_example_tenant",
			"customer_number",
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
