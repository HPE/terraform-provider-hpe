// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package tenant_test

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/role"
	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/tenant"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/constants"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

func TestMain(m *testing.M) {
	code := testhelpers.TestMain(m)
	testhelpers.WriteMergedResults()
	os.Exit(code)
}

func TestAccMorpheusTenantExampleOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()
	name := acctest.RandomWithPrefix(t.Name())
	subdomain := name
	resourceName := "hpe_morpheus_tenant.example"

	// parent_id is computed to the authenticated user's tenant when omitted —
	// but only on appliances with the tenant hierarchy (8.1.0+). Older
	// appliances return no parent, so it is null there.
	parentIDCheck := resource.TestCheckNoResourceAttr(resourceName, "parent_id")
	if testhelpers.ApplianceVersionAtLeast(
		context.Background(), t, constants.TenantParentMinVersion,
	) {
		parentIDCheck = resource.TestCheckResourceAttrSet(resourceName, "parent_id")
	}

	// A tenant-scoped base role is required to create a tenant. Create one as a
	// dependency and reference its id, rather than depending on appliance-seeded
	// admin roles.
	roleConfig, err := role.RenderRoleTenantConfig(t, map[string]string{
		"Name": name,
	})
	if err != nil {
		t.Fatal(err)
	}

	config, err := tenant.RenderTenantConfig(t, map[string]string{
		"Name":       name,
		"Subdomain":  subdomain,
		"BaseRoleId": "hpe_morpheus_role.example.id",
	})
	if err != nil {
		t.Fatal(err)
	}

	updateConfig, err := tenant.RenderTenantConfig(t, map[string]string{
		"Name":        name,
		"Subdomain":   subdomain,
		"BaseRoleId":  "hpe_morpheus_role.example.id",
		"Description": "Updated tenant description",
		"AccountName": "tenant 99999",
	})
	if err != nil {
		t.Fatal(err)
	}

	// The same configuration with description removed entirely (an empty
	// override drops the line from the template). description is Optional and
	// not Computed precisely so that this clears the value: Update sends an
	// explicit null and Morpheus nullifies the field, rather than the previous
	// value being carried forward as it was with the SDKv2 resource.
	clearedConfig, err := tenant.RenderTenantConfig(t, map[string]string{
		"Name":        name,
		"Subdomain":   subdomain,
		"BaseRoleId":  "hpe_morpheus_role.example.id",
		"Description": "",
		"AccountName": "tenant 99999",
	})
	if err != nil {
		t.Fatal(err)
	}

	baseChecks := func(description, accountName string) resource.TestCheckFunc {
		descriptionCheck := resource.TestCheckResourceAttr(resourceName, "description", description)
		if description == "" {
			descriptionCheck = resource.TestCheckNoResourceAttr(resourceName, "description")
		}

		return resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttr(resourceName, "name", name),
			descriptionCheck,
			resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
			resource.TestCheckResourceAttr(resourceName, "subdomain", subdomain),
			resource.TestCheckResourceAttr(resourceName, "currency", "USD"),
			resource.TestCheckResourceAttr(resourceName, "account_number", "12345"),
			resource.TestCheckResourceAttr(resourceName, "account_name", accountName),
			resource.TestCheckResourceAttr(resourceName, "customer_number", "12345"),
			resource.TestCheckResourceAttr(resourceName, "master", "false"),
			resource.TestCheckResourceAttrSet(resourceName, "id"),
			resource.TestCheckResourceAttrSet(resourceName, "base_role_id"),
			parentIDCheck,
			resource.TestCheckResourceAttrSet(resourceName, "date_created"),
		)
	}

	checkInPlaceUpdate := resource.ConfigPlanChecks{
		PreApply: []plancheck.PlanCheck{
			plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate),
		},
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			// Apply
			{
				Config:             providerConfig + roleConfig + config,
				ExpectNonEmptyPlan: false,
				Check:              baseChecks("Terraform example tenant", "tenant 12345"),
			},
			// Plan after apply: no drift
			{
				Config:             providerConfig + roleConfig + config,
				ExpectNonEmptyPlan: false,
				PlanOnly:           true,
			},
			// Update in place
			{
				Config:           providerConfig + roleConfig + updateConfig,
				ConfigPlanChecks: checkInPlaceUpdate,
				Check:            baseChecks("Updated tenant description", "tenant 99999"),
			},
			// Clear the description by removing it from the configuration: an
			// in-place update after which the API reports no description.
			{
				Config:           providerConfig + roleConfig + clearedConfig,
				ConfigPlanChecks: checkInPlaceUpdate,
				Check:            baseChecks("", "tenant 99999"),
			},
			// Plan after clearing: no drift (the null is stable).
			{
				Config:             providerConfig + roleConfig + clearedConfig,
				ExpectNonEmptyPlan: false,
				PlanOnly:           true,
			},
			// Import
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[resourceName]
					if !ok {
						return "", fmt.Errorf("resource not found: %s", resourceName)
					}

					return rs.Primary.Attributes["id"], nil
				},
			},
		},
	})
}

// TestAccMorpheusTenantInvalidSubdomain asserts the plan-time subdomain
// validator rejects an all-numeric subdomain by naming the attribute, rather
// than surfacing a generic 400 at apply.
func TestAccMorpheusTenantInvalidSubdomain(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()

	// base_role_id is a literal so the plan does not depend on any other
	// resource; validation fails in ModifyPlan before any API call.
	config := `
resource "hpe_morpheus_tenant" "invalid" {
  name         = "tf-invalid-subdomain"
  base_role_id = 1
  subdomain    = "12345"
}
`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config:      providerConfig + config,
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`(?s)subdomain.*entirely numeric`),
			},
		},
	})
}

// TestAccMorpheusTenantParentIdRequiresNewerAppliance asserts the plan-time
// appliance version gate on parent_id: on appliances older than
// constants.TenantParentMinVersion (8.1.0, which introduced the tenant
// hierarchy) a configuration that nominates a parent is refused with a
// diagnostic naming the attribute and the required version, instead of being
// applied — where the appliance would ignore parentAccount and every
// subsequent plan would force a replacement. The gate is inert on 8.1.0+, so
// the test only runs against older appliances.
func TestAccMorpheusTenantParentIdRequiresNewerAppliance(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	if testhelpers.ApplianceVersionAtLeast(
		context.Background(), t, constants.TenantParentMinVersion,
	) {
		t.Skipf("appliance satisfies %q; the parent_id gate does not apply",
			constants.TenantParentMinVersion)
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()

	// Literal ids so the plan does not depend on any other resource; the gate
	// refuses in ModifyPlan before the parent or role is ever looked up.
	config := `
resource "hpe_morpheus_tenant" "gated" {
  name         = "tf-parent-id-gate"
  base_role_id = 1
  parent_id    = 1
}
`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config:      providerConfig + config,
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`(?s)parent_id.*require a Morpheus appliance version >= 8\.1\.0`),
			},
		},
	})
}
