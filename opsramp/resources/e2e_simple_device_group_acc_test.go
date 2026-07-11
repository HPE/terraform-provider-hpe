// SPDX-FileCopyrightText: Copyright Hewlett Packard Enterprise Development LP
// SPDX-License-Identifier: Apache-2.0

package resources_test

import (
	"fmt"
	"testing"

	"github.com/HPE/terraform-provider-hpe/opsramp/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccE2ESimpleDeviceGroup exercises the simple-device-group e2e scenario:
// resources + root device group + child groups with resources, queries, and mixed.
func TestAccE2ESimpleDeviceGroup(t *testing.T) {
	t.Run("full hierarchy", func(t *testing.T) {
		res1 := acctest.RandomName("dg-res1")
		res2 := acctest.RandomName("dg-res2")
		res3 := acctest.RandomName("dg-res3")
		rootGroup := acctest.RandomName("dg-root")
		childRes := acctest.RandomName("dg-child-res")
		childQuery := acctest.RandomName("dg-child-qry")
		childMixed := acctest.RandomName("dg-child-mix")

		resource.ParallelTest(t, resource.TestCase{
			PreCheck:                 acctest.PreCheck(t),
			ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories(),
			CheckDestroy:             testAccCheckDeviceGroupDestroy(t),
			Steps: []resource.TestStep{
				{
					Config: testAccE2ESimpleDeviceGroupConfig(res1, res2, res3, rootGroup, childRes, childQuery, childMixed),
					Check: resource.ComposeAggregateTestCheckFunc(
						testAccEnsureResourceExists(t, "opsramp_resource.resource1"),
						testAccEnsureResourceExists(t, "opsramp_resource.resource2"),
						testAccEnsureResourceExists(t, "opsramp_resource.resource3"),
						testAccEnsureDeviceGroupExists(t, "opsramp_device_group.device_group_root"),
						testAccEnsureDeviceGroupExists(t, "opsramp_device_group.device_group_resources"),
						testAccEnsureDeviceGroupExists(t, "opsramp_device_group.device_group_query"),
						testAccEnsureDeviceGroupExists(t, "opsramp_device_group.device_group_mixed"),
						resource.TestCheckResourceAttr("opsramp_device_group.device_group_root", "name", rootGroup),
						resource.TestCheckResourceAttr("opsramp_device_group.device_group_resources", "name", childRes),
						resource.TestCheckResourceAttr("opsramp_device_group.device_group_query", "name", childQuery),
						resource.TestCheckResourceAttr("opsramp_device_group.device_group_mixed", "name", childMixed),
					),
				},
			},
		})
	})
}

func testAccE2ESimpleDeviceGroupConfig(res1, res2, res3, rootGroup, childRes, childQuery, childMixed string) string {
	return fmt.Sprintf(`
%s
resource "opsramp_resource" "resource1" {
	resource_name = "%s"
	resource_type = "Linux"
}

resource "opsramp_resource" "resource2" {
	resource_name = "%s"
	resource_type = "Linux"
}

resource "opsramp_resource" "resource3" {
	resource_name = "%s"
	resource_type = "Linux"
}

resource "opsramp_device_group" "device_group_root" {
	name      = "%s"
	resources = []
}

resource "opsramp_device_group" "device_group_resources" {
	parent_id = opsramp_device_group.device_group_root.id
	name      = "%s"
	resources = [opsramp_resource.resource1.uuid]
}

resource "opsramp_device_group" "device_group_query" {
	parent_id    = opsramp_device_group.device_group_root.id
	name         = "%s"
	search_query = format("resourceType = \"Linux\" AND uuid = \"%%s\"", opsramp_resource.resource2.uuid)
}

resource "opsramp_device_group" "device_group_mixed" {
	parent_id    = opsramp_device_group.device_group_root.id
	name         = "%s"
	search_query = format("resourceType = \"Linux\" AND uuid = \"%%s\"", opsramp_resource.resource2.uuid)
	resources    = [opsramp_resource.resource3.uuid]
}
`, acctest.ProviderConfigHCL(), res1, res2, res3, rootGroup, childRes, childQuery, childMixed)
}
