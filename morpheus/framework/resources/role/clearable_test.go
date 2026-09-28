// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package role_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"

	sdkv2morpheus "github.com/HPE/terraform-provider-hpe/morpheus/sdkv2"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

// Removing description or landing_url clears it on the role (MORPH-17990).
// Both are Optional without Computed: Terraform plans the null and Update
// sends it.
func TestAccMorpheusRoleRemovedAttributesAreCleared(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	name := acctest.RandomWithPrefix(t.Name())
	const rn = "hpe_morpheus_role.test"

	config := func(extra string) string {
		return testhelpers.ProviderBlock() + fmt.Sprintf(`
resource "hpe_morpheus_role" "test" {
  name      = %q
  role_type = "user"
%s
}
`, name, extra)
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), sdkv2morpheus.Provider()),
		Steps: []resource.TestStep{
			{
				Config: config(`  description = "to be cleared"
  landing_url = "https://example.test/landing"`),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(rn, "description", "to be cleared"),
					resource.TestCheckResourceAttr(rn, "landing_url", "https://example.test/landing"),
				),
			},
			{
				Config: config(""),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(rn, plancheck.ResourceActionUpdate),
						plancheck.ExpectKnownValue(rn, tfjsonpath.New("description"), knownvalue.Null()),
						plancheck.ExpectKnownValue(rn, tfjsonpath.New("landing_url"), knownvalue.Null()),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckNoResourceAttr(rn, "description"),
					resource.TestCheckNoResourceAttr(rn, "landing_url"),
				),
			},
			{
				Config:   config(""),
				PlanOnly: true,
			},
		},
	})
}
