// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package tenant_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

// Removing an optional attribute clears it on the tenant (MORPH-17990):
// description (Optional without Computed) and the Optional+Computed billing
// fields and subdomain, whose removal used to be silently ignored.
func TestAccMorpheusTenantRemovedAttributesAreCleared(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All, capabilities.MasterTenant)

	name := acctest.RandomWithPrefix(t.Name())
	// Subdomains are globally unique and limited to 63 characters.
	subdomain := "clr-" + acctest.RandString(12)
	const rn = "hpe_morpheus_tenant.test"

	config := func(with ...string) string {
		set := map[string]string{
			"description":     `description     = "to be cleared"`,
			"subdomain":       `subdomain       = "` + strings.ToLower(subdomain) + `"`,
			"account_number":  `account_number  = "ACC-1"`,
			"account_name":    `account_name    = "Account One"`,
			"customer_number": `customer_number = "CUST-1"`,
		}

		lines := make([]string, 0, len(with))
		for _, a := range with {
			lines = append(lines, "  "+set[a])
		}

		return testhelpers.ProviderBlock() + fmt.Sprintf(`
resource "hpe_morpheus_tenant" "test" {
  name = %q
%s
}
`, name, strings.Join(lines, "\n"))
	}

	all := []string{"description", "subdomain", "account_number", "account_name", "customer_number"}

	nullChecks := []plancheck.PlanCheck{plancheck.ExpectResourceAction(rn, plancheck.ResourceActionUpdate)}
	for _, a := range all {
		nullChecks = append(nullChecks, plancheck.ExpectKnownValue(rn, tfjsonpath.New(a), knownvalue.Null()))
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: config(all...),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(rn, "description", "to be cleared"),
					resource.TestCheckResourceAttr(rn, "subdomain", strings.ToLower(subdomain)),
					resource.TestCheckResourceAttr(rn, "account_number", "ACC-1"),
					resource.TestCheckResourceAttr(rn, "account_name", "Account One"),
					resource.TestCheckResourceAttr(rn, "customer_number", "CUST-1"),
				),
			},
			{
				Config:           config(),
				ConfigPlanChecks: resource.ConfigPlanChecks{PreApply: nullChecks},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckNoResourceAttr(rn, "description"),
					resource.TestCheckNoResourceAttr(rn, "subdomain"),
					resource.TestCheckNoResourceAttr(rn, "account_number"),
					resource.TestCheckNoResourceAttr(rn, "account_name"),
					resource.TestCheckNoResourceAttr(rn, "customer_number"),
				),
			},
			{
				Config:   config(),
				PlanOnly: true,
			},
		},
	})
}
