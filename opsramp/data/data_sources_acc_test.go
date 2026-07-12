// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package data_test

import (
	"fmt"
	"testing"

	"github.com/HPE/terraform-provider-hpe/opsramp/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccTenantDataSource(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		resource.ParallelTest(t, resource.TestCase{
			PreCheck:                 acctest.PreCheck(t),
			ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories(),
			Steps: []resource.TestStep{
				{
					Config: testAccTenantDataSourceConfig(),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttrSet("data.hpe_opsramp_tenant.test", "uuid"),
						resource.TestCheckResourceAttrSet("data.hpe_opsramp_tenant.test", "name"),
					),
				},
			},
		})
	})
}

func TestAccRoleDataSource(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		resource.ParallelTest(t, resource.TestCase{
			PreCheck:                 acctest.PreCheck(t),
			ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories(),
			Steps: []resource.TestStep{
				{
					Config: testAccRoleDataSourceConfig(),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttrSet("data.hpe_opsramp_role.test", "id"),
						resource.TestCheckResourceAttr("data.hpe_opsramp_role.test", "name", "Client Administrator"),
					),
				},
			},
		})
	})
}

func TestAccResourceLookupDataSource(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		resourceName := acctest.RandomName("lookup-res")

		resource.ParallelTest(t, resource.TestCase{
			PreCheck:                 acctest.PreCheck(t),
			ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories(),
			Steps: []resource.TestStep{
				{
					Config: testAccResourceLookupDataSourceConfig(resourceName),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttrSet("data.hpe_opsramp_resource_lookup.test", "exists"),
					),
				},
			},
		})
	})
}

func TestAccServicedeskBusinessImpactDataSource(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		impactName := acctest.RandomName("ds-impact")

		resource.ParallelTest(t, resource.TestCase{
			PreCheck:                 acctest.PreCheck(t),
			ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories(),
			Steps: []resource.TestStep{
				{
					Config: testAccServicedeskBusinessImpactDataSourceConfig(impactName),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttrSet("data.hpe_opsramp_servicedesk_business_impact.test", "id"),
						resource.TestCheckResourceAttr("data.hpe_opsramp_servicedesk_business_impact.test", "name", impactName),
					),
				},
			},
		})
	})
}

func TestAccServicedeskCategoryDataSource(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		catName := acctest.RandomName("ds-cat")

		resource.ParallelTest(t, resource.TestCase{
			PreCheck:                 acctest.PreCheck(t),
			ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories(),
			Steps: []resource.TestStep{
				{
					Config: testAccServicedeskCategoryDataSourceConfig(catName),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttrSet("data.hpe_opsramp_servicedesk_category.test", "id"),
						resource.TestCheckResourceAttr("data.hpe_opsramp_servicedesk_category.test", "name", catName),
					),
				},
			},
		})
	})
}

func TestAccServicedeskUrgencyDataSource(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		urgencyName := acctest.RandomName("ds-urgency")

		resource.ParallelTest(t, resource.TestCase{
			PreCheck:                 acctest.PreCheck(t),
			ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories(),
			Steps: []resource.TestStep{
				{
					Config: testAccServicedeskUrgencyDataSourceConfig(urgencyName),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttrSet("data.hpe_opsramp_servicedesk_urgency.test", "id"),
						resource.TestCheckResourceAttr("data.hpe_opsramp_servicedesk_urgency.test", "name", urgencyName),
					),
				},
			},
		})
	})
}

func testAccTenantDataSourceConfig() string {
	return fmt.Sprintf(`
%s
data "hpe_opsramp_tenant" "test" {}
`, acctest.ProviderConfigHCL())
}

func testAccRoleDataSourceConfig() string {
	return fmt.Sprintf(`
%s
data "hpe_opsramp_role" "test" {
	name = "Client Administrator"
}
`, acctest.ProviderConfigHCL())
}

func testAccResourceLookupDataSourceConfig(resourceName string) string {
	return fmt.Sprintf(`
%s
resource "hpe_opsramp_resource" "lookup_test" {
	resource_name = "%s"
	resource_type = "Linux"
}

data "hpe_opsramp_resource_lookup" "test" {
	query = format("name = \"%%s\"", hpe_opsramp_resource.lookup_test.resource_name)
}
`, acctest.ProviderConfigHCL(), resourceName)
}

func testAccServicedeskBusinessImpactDataSourceConfig(name string) string {
	return fmt.Sprintf(`
%s
resource "hpe_opsramp_servicedesk_business_impact" "ds_test" {
	name        = "%s"
	description = "Data source test business impact"
}

data "hpe_opsramp_servicedesk_business_impact" "test" {
	name = hpe_opsramp_servicedesk_business_impact.ds_test.name
}
`, acctest.ProviderConfigHCL(), name)
}

func testAccServicedeskCategoryDataSourceConfig(name string) string {
	return fmt.Sprintf(`
%s
resource "hpe_opsramp_servicedesk_category" "ds_test" {
	name        = "%s"
	description = "Data source test category"
	ticket_type = "serviceRequests"
}

data "hpe_opsramp_servicedesk_category" "test" {
	name = hpe_opsramp_servicedesk_category.ds_test.name
}
`, acctest.ProviderConfigHCL(), name)
}

func testAccServicedeskUrgencyDataSourceConfig(name string) string {
	return fmt.Sprintf(`
%s
resource "hpe_opsramp_servicedesk_urgency" "ds_test" {
	name        = "%s"
	description = "Data source test urgency"
}

data "hpe_opsramp_servicedesk_urgency" "test" {
	name = hpe_opsramp_servicedesk_urgency.ds_test.name
}
`, acctest.ProviderConfigHCL(), name)
}
