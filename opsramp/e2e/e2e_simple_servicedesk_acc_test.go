// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package e2e_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/opsramp/acctest"
)

// TestAccE2ESimpleServicedesk exercises the simple-servicedesk e2e scenario:
// multiple categories, business impacts, and urgencies.
func TestAccE2ESimpleServicedesk(t *testing.T) {
	clientOverride := acctest.OptionalClientOverride(t)

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
					Config: testAccE2ESimpleServicedeskConfig(
						cat1,
						cat2,
						cat3,
						impact1,
						impact2,
						impact3,
						urgency1,
						urgency2,
						urgency3,
						clientOverride,
					),
					Check: resource.ComposeAggregateTestCheckFunc(
						testAccEnsureServicedeskCategoryExists(t, "hpe_opsramp_servicedesk_category.category1"),
						testAccEnsureServicedeskCategoryExists(t, "hpe_opsramp_servicedesk_category.category2"),
						testAccEnsureServicedeskCategoryExists(t, "hpe_opsramp_servicedesk_category.category3"),
						testAccEnsureServicedeskBusinessImpactExists(t, "hpe_opsramp_servicedesk_business_impact.impact1"),
						testAccEnsureServicedeskBusinessImpactExists(t, "hpe_opsramp_servicedesk_business_impact.impact2"),
						testAccEnsureServicedeskBusinessImpactExists(t, "hpe_opsramp_servicedesk_business_impact.impact3"),
						testAccEnsureServicedeskUrgencyExists(t, "hpe_opsramp_servicedesk_urgency.urgency1"),
						testAccEnsureServicedeskUrgencyExists(t, "hpe_opsramp_servicedesk_urgency.urgency2"),
						testAccEnsureServicedeskUrgencyExists(t, "hpe_opsramp_servicedesk_urgency.urgency3"),
						resource.TestCheckResourceAttr("hpe_opsramp_servicedesk_category.category1", "name", cat1),
						resource.TestCheckResourceAttr("hpe_opsramp_servicedesk_business_impact.impact1", "name", impact1),
						resource.TestCheckResourceAttr("hpe_opsramp_servicedesk_urgency.urgency1", "name", urgency1),
					),
				},
			},
		})
	})
}

func testAccE2ESimpleServicedeskConfig(
	cat1, cat2, cat3, impact1, impact2, impact3,
	urgency1, urgency2, urgency3, clientOverride string,
) string {
	clientAttr := acctest.ClientAttrHCL(clientOverride)

	return fmt.Sprintf(
		`
%s
resource "hpe_opsramp_servicedesk_category" "category1" {
	name        = "%s"
	description = "Category1 Description"
	ticket_type = "serviceRequests"
	%s
}

resource "hpe_opsramp_servicedesk_category" "category2" {
	name        = "%s"
	description = "Category2 Description"
	ticket_type = "incidents"
	%s
}

resource "hpe_opsramp_servicedesk_category" "category3" {
	name        = "%s"
	description = "Category3 Description"
	ticket_type = "problems"
	%s
}

resource "hpe_opsramp_servicedesk_business_impact" "impact1" {
	name        = "%s"
	description = "Business Impact1 Description"
	%s
}

resource "hpe_opsramp_servicedesk_business_impact" "impact2" {
	name        = "%s"
	description = "Business Impact2 Description"
	%s
}

resource "hpe_opsramp_servicedesk_business_impact" "impact3" {
	name        = "%s"
	description = "Business Impact3 Description"
	%s
}

resource "hpe_opsramp_servicedesk_urgency" "urgency1" {
	name        = "%s"
	description = "Urgency1 Description"
	%s
}

resource "hpe_opsramp_servicedesk_urgency" "urgency2" {
	name        = "%s"
	description = "Urgency2 Description"
	%s
}

resource "hpe_opsramp_servicedesk_urgency" "urgency3" {
	name        = "%s"
	description = "Urgency3 Description"
	%s
}
`,
		acctest.ProviderConfigHCL(),
		cat1, clientAttr,
		cat2, clientAttr,
		cat3, clientAttr,
		impact1, clientAttr,
		impact2, clientAttr,
		impact3, clientAttr,
		urgency1, clientAttr,
		urgency2, clientAttr,
		urgency3, clientAttr,
	)
}
