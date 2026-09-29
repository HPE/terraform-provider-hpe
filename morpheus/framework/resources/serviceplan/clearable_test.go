// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package serviceplan_test

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
)

// Removing the description clears it on the service plan (MORPH-17990). It
// is Optional-only: Terraform plans the null, and Update now sends it.
func TestAccMorpheusServicePlanRemovedAttributesAreCleared(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.All)

	name := acctest.RandomWithPrefix(t.Name())
	const rn = "hpe_morpheus_service_plan.test"

	config := func(extra string) string {
		return testhelpers.ProviderBlock() + fmt.Sprintf(`
resource "hpe_morpheus_service_plan" "test" {
  name                = %q
  code                = %q
  provision_type_code = "arm"
  max_memory          = 4294967296
  max_storage         = 536870912
  cores_per_socket    = 1
%s
}
`, name, strings.ToLower(name), extra)
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config(`  description = "to be cleared"`),
				Check:  resource.TestCheckResourceAttr(rn, "description", "to be cleared"),
			},
			{
				Config: config(""),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(rn, plancheck.ResourceActionUpdate),
						plancheck.ExpectKnownValue(rn, tfjsonpath.New("description"), knownvalue.Null()),
					},
				},
				Check: resource.TestCheckNoResourceAttr(rn, "description"),
			},
			{
				Config:   config(""),
				PlanOnly: true,
			},
		},
	})
}
