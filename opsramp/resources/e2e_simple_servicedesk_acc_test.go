// SPDX-FileCopyrightText: Copyright Hewlett Packard Enterprise Development LP
// SPDX-License-Identifier: Apache-2.0

package resources_test

import (
	"fmt"
	"testing"

	"github.com/HPE/terraform-provider-opsramp/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccE2ESimpleServicedesk exercises the simple-servicedesk e2e scenario:
// multiple categories, business impacts, and urgencies.
func TestAccE2ESimpleServicedesk(t *testing.T) {
	t.Run("multiple servicedesk entries", func(t *testing.T) {
		cat1 := acctest.RandomName("e2e-cat1")
		cat2 := acctest.RandomName("e2e-cat2")
		cat3 := acctest.RandomName("e2e-cat3")
		impact1 := acctest.RandomName("e2e-impact1")
		impact2 := acctest.RandomName("e2e-impact2")
		impact3 := acctest.RandomName("e2e-impact3")
		urgency1 := acctest.RandomName("e2e-urgency1")
		urgency2 := acctest.RandomName("e2e-urgency2")
		urgency3 := acctest.RandomName("e2e-urgency3")

		resource.ParallelTest(t, resource.TestCase{
			PreCheck:                 acctest.PreCheck(t),
			ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories(),
			Steps: []resource.TestStep{
				{
					Config: testAccE2ESimpleServicedeskConfig(cat1, cat2, cat3, impact1, impact2, impact3, urgency1, urgency2, urgency3),
					Check: resource.ComposeAggregateTestCheckFunc(
						testAccEnsureServicedeskCategoryExists(t, "opsramp_servicedesk_category.category1"),
						testAccEnsureServicedeskCategoryExists(t, "opsramp_servicedesk_category.category2"),
						testAccEnsureServicedeskCategoryExists(t, "opsramp_servicedesk_category.category3"),
						testAccEnsureServicedeskBusinessImpactExists(t, "opsramp_servicedesk_business_impact.impact1"),
						testAccEnsureServicedeskBusinessImpactExists(t, "opsramp_servicedesk_business_impact.impact2"),
						testAccEnsureServicedeskBusinessImpactExists(t, "opsramp_servicedesk_business_impact.impact3"),
						testAccEnsureServicedeskUrgencyExists(t, "opsramp_servicedesk_urgency.urgency1"),
						testAccEnsureServicedeskUrgencyExists(t, "opsramp_servicedesk_urgency.urgency2"),
						testAccEnsureServicedeskUrgencyExists(t, "opsramp_servicedesk_urgency.urgency3"),
						resource.TestCheckResourceAttr("opsramp_servicedesk_category.category1", "name", cat1),
						resource.TestCheckResourceAttr("opsramp_servicedesk_business_impact.impact1", "name", impact1),
						resource.TestCheckResourceAttr("opsramp_servicedesk_urgency.urgency1", "name", urgency1),
					),
				},
			},
		})
	})
}

func testAccE2ESimpleServicedeskConfig(cat1, cat2, cat3, impact1, impact2, impact3, urgency1, urgency2, urgency3 string) string {
	return fmt.Sprintf(`
%s
resource "opsramp_servicedesk_category" "category1" {
	name        = "%s"
	description = "Category1 Description"
	ticket_type = "serviceRequests"
}

resource "opsramp_servicedesk_category" "category2" {
	name        = "%s"
	description = "Category2 Description"
	ticket_type = "incidents"
}

resource "opsramp_servicedesk_category" "category3" {
	name        = "%s"
	description = "Category3 Description"
	ticket_type = "problems"
}

resource "opsramp_servicedesk_business_impact" "impact1" {
	name        = "%s"
	description = "Business Impact1 Description"
}

resource "opsramp_servicedesk_business_impact" "impact2" {
	name        = "%s"
	description = "Business Impact2 Description"
}

resource "opsramp_servicedesk_business_impact" "impact3" {
	name        = "%s"
	description = "Business Impact3 Description"
}

resource "opsramp_servicedesk_urgency" "urgency1" {
	name        = "%s"
	description = "Urgency1 Description"
}

resource "opsramp_servicedesk_urgency" "urgency2" {
	name        = "%s"
	description = "Urgency2 Description"
}

resource "opsramp_servicedesk_urgency" "urgency3" {
	name        = "%s"
	description = "Urgency3 Description"
}
`, acctest.ProviderConfigHCL(), cat1, cat2, cat3, impact1, impact2, impact3, urgency1, urgency2, urgency3)
}
