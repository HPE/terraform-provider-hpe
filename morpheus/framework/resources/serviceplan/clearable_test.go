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

// Removing the description or the disk count limit clears it on the service
// plan (MORPH-17990). The description is Optional-only: Terraform plans the
// null, and Update sends it. max_disks is Optional+Computed: ModifyPlan plans
// the clear once an apply has recorded it as configured, and Update sends
// null -- never 0, which Morpheus stores as a limit that rejects every
// instance on the plan (MORPH-18341); a cleared limit reads back unset, like
// one never configured.
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
				Config: config("  description = \"to be cleared\"\n  max_disks   = 2"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(rn, "description", "to be cleared"),
					resource.TestCheckResourceAttr(rn, "max_disks", "2"),
				),
			},
			{
				Config: config(""),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(rn, plancheck.ResourceActionUpdate),
						plancheck.ExpectKnownValue(rn, tfjsonpath.New("description"), knownvalue.Null()),
						plancheck.ExpectKnownValue(rn, tfjsonpath.New("max_disks"), knownvalue.Null()),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckNoResourceAttr(rn, "description"),
					resource.TestCheckNoResourceAttr(rn, "max_disks"),
				),
			},
			{
				Config:   config(""),
				PlanOnly: true,
			},
		},
	})
}
