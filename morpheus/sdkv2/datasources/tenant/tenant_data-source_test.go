// (C) Copyright 2025-2026 Hewlett Packard Enterprise Development LP

package tenant_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/role"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/tenant"
	sdkv2morpheus "github.com/HPE/terraform-provider-hpe/morpheus/sdkv2"
	dstenant "github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/datasources/tenant"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

func TestAccMorpheusDataSourceTenantExampleOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	t.Parallel()

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())
	subdomain := "tf" + acctest.RandStringFromCharSet(12, acctest.CharSetAlpha)

	// Create a tenant (and its base role) as a dependency so we can test
	// searching by name. The name of the master tenant is not guaranteed across
	// appliances, so we prefer to create one.
	roleConfig, err := role.RenderRoleTenantConfig(t, map[string]string{
		"Name": name,
	})
	if err != nil {
		t.Fatal(err)
	}

	tenantConfig, err := tenant.RenderTenantConfig(t, map[string]string{
		"Name":       name,
		"Subdomain":  subdomain,
		"BaseRoleId": "hpe_morpheus_role.example.id",
	})
	if err != nil {
		t.Fatal(err)
	}

	dependenciesConfig := roleConfig + tenantConfig

	datasourceConfig, err := dstenant.RenderTenantConfig(t, map[string]string{
		"Name": "resource.hpe_morpheus_tenant.example.name",
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
			"hpe_morpheus_tenant.example",
			"id",
		),
		resource.TestCheckResourceAttrPair(
			"data.hpe_morpheus_tenant.example",
			"account_number",
			"hpe_morpheus_tenant.example",
			"account_number",
		),
		resource.TestCheckResourceAttrPair(
			"data.hpe_morpheus_tenant.example",
			"account_name",
			"hpe_morpheus_tenant.example",
			"account_name",
		),
		resource.TestCheckResourceAttrPair(
			"data.hpe_morpheus_tenant.example",
			"customer_number",
			"hpe_morpheus_tenant.example",
			"customer_number",
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
