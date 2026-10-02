// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package data_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/opsramp/acctest"
)

func TestAccDeviceManagementPolicyDataSource(t *testing.T) {
	acctest.SkipIfNotClient(t)
	clientOverride := acctest.RequireClientScope(t)

	t.Run("happy path", func(t *testing.T) {
		name := acctest.RandomName("ds-device-mgmt-policy")

		resource.ParallelTest(t, resource.TestCase{
			PreCheck:                 acctest.PreCheck(t),
			ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories(),
			Steps: []resource.TestStep{
				{
					Config: testAccDeviceManagementPolicyDataSourceConfig(name, clientOverride),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttrSet("data.hpe_opsramp_device_management_policy.test", "id"),
						resource.TestCheckResourceAttrSet("data.hpe_opsramp_device_management_policy.test", "uid"),
						resource.TestCheckResourceAttr("data.hpe_opsramp_device_management_policy.test", "name", name),
						resource.TestCheckResourceAttrPair(
							"data.hpe_opsramp_device_management_policy.test", "id",
							"hpe_opsramp_device_management_policy.ds_test", "id",
						),
						resource.TestCheckResourceAttrPair(
							"data.hpe_opsramp_device_management_policy.test", "uid",
							"hpe_opsramp_device_management_policy.ds_test", "uid",
						),
					),
				},
			},
		})
	})
}

func testAccDeviceManagementPolicyDataSourceConfig(name string, clientOverride string) string {
	clientAttr := acctest.ClientAttrHCL(clientOverride)

	return fmt.Sprintf(`
%s
resource "hpe_opsramp_device_management_policy" "ds_test" {
	name = %q
	%s

	search_query = "name CONTAINS \"name\""

	actions = [
		{
			action = "Assign Availability Rule"
			items = [
				{
					item = "AND"
				}
			]
		}
	]
}

data "hpe_opsramp_device_management_policy" "test" {
	name = hpe_opsramp_device_management_policy.ds_test.name
	%s
}
`, acctest.ProviderConfigHCL(), name, clientAttr, clientAttr)
}
