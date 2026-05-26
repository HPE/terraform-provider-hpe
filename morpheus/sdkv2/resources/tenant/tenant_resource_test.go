// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package tenant_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/role"
	sdkv2morpheus "github.com/HPE/terraform-provider-hpe/morpheus/sdkv2"
	"github.com/HPE/terraform-provider-hpe/morpheus/sdkv2/resources/tenant"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
)

func TestMain(m *testing.M) {
	code := m.Run()

	testhelpers.WriteMergedResults()

	os.Exit(code)
}

func TestAccMorpheusTenantExampleOk(t *testing.T) {
	if capabilities.Missing(t, capabilities.All) {
		t.Log("Skipping test due to missing capabilities")

		return
	}
	t.Parallel()

	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())

	dependencyResourceConfig, err := role.RenderRoleTenantConfig(t, map[string]string{
		"Name": name,
	})
	if err != nil {
		t.Fatal(err)
	}

	resourceConfig, err := tenant.RenderTenantConfig(t, map[string]string{
		"Name":       name,
		"BaseRoleId": "hpe_morpheus_role.example.id",
	})
	if err != nil {
		t.Fatal(err)
	}

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_tenant.tf_example_tenant",
			"name",
			name,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_tenant.tf_example_tenant",
			"description",
			"Terraform example tenant",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_tenant.tf_example_tenant",
			"enabled",
			"true",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_tenant.tf_example_tenant",
			"subdomain",
			"tfexample",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_tenant.tf_example_tenant",
			"currency",
			"USD",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_tenant.tf_example_tenant",
			"account_number",
			"12345",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_tenant.tf_example_tenant",
			"account_name",
			"tenant 12345",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_tenant.tf_example_tenant",
			"customer_number",
			"12345",
		),
		resource.TestCheckResourceAttrSet(
			"hpe_morpheus_tenant.tf_example_tenant",
			"base_role_id",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(checks...)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), sdkv2morpheus.Provider()),
		Steps: []resource.TestStep{
			// Apply
			{
				Config:             providerConfig + dependencyResourceConfig + resourceConfig,
				ExpectNonEmptyPlan: false,
				Check:              checkFn,
			},
			// Plan after apply
			{
				Config:             providerConfig + dependencyResourceConfig + resourceConfig,
				ExpectNonEmptyPlan: false,
				PlanOnly:           true,
			},
		},
	})
}
