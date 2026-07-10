// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package role_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	sdkv2morpheus "github.com/HPE/terraform-provider-hpe/morpheus/sdkv2"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

// TestAccMorpheusRoleTenantCopiesOk verifies that the computed tenant_copies
// attribute is populated as a known list on a multitenant master role, and that
// each entry exposes tenant_id, role_id, and diverged.
//
// The role show API only returns tenantCopies on Morpheus 9.0.2 and later,
// so this test is version-gated via SkipUnlessApplianceVersionAtLeast: on older
// appliances it is skipped rather than asserting a field the appliance cannot
// return.
//
// Note: this asserts tenant_copies is a known (possibly empty) list. Verifying
// non-empty copies additionally requires a subtenant with the master role
// propagated into it; that fuller multitenant scenario is left as a follow-up.
func TestAccMorpheusRoleTenantCopiesOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	// tenant_copies is only populated by the role show API on Morpheus 9.0.2
	// and later; skip the assertion on older appliances.
	testhelpers.SkipUnlessApplianceVersionAtLeast(
		context.Background(), t, ">= 9.0.2",
	)

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()

	name := acctest.RandomWithPrefix(t.Name())

	resourceConfig := `
resource "hpe_morpheus_role" "tenant_copies" {
  name        = "` + name + `"
  role_type   = "user"
  multitenant = true
}
`
	resourceName := "hpe_morpheus_role.tenant_copies"

	checks := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr(resourceName, "name", name),
		resource.TestCheckResourceAttr(resourceName, "multitenant", "true"),
		// tenant_copies is a known, computed list on >= 9.0.2 (0+ entries).
		resource.TestCheckResourceAttrSet(resourceName, "tenant_copies.#"),
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(
			t, adapter.NewMorpheus(), sdkv2morpheus.Provider(),
		),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + resourceConfig,
				Check:  checks,
			},
		},
	})
}
