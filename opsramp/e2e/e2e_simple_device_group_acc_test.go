// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package e2e_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/opsramp/acctest"
)

// TestAccE2ESimpleDeviceGroup exercises the simple-device-group e2e scenario:
// resources + root device group + child groups with resources, queries, and mixed.
func TestAccE2ESimpleDeviceGroup(t *testing.T) {
	acctest.SkipIfNotClient(t)
	clientOverride := acctest.OptionalClientOverride(t)

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
					Config: testAccE2ESimpleDeviceGroupConfig(res1, res2, res3, rootGroup, childRes, childQuery, childMixed, clientOverride),
					Check: resource.ComposeAggregateTestCheckFunc(
						testAccEnsureResourceExists(t, "hpe_opsramp_resource.resource1"),
						testAccEnsureResourceExists(t, "hpe_opsramp_resource.resource2"),
						testAccEnsureResourceExists(t, "hpe_opsramp_resource.resource3"),
						testAccEnsureDeviceGroupExists(t, "hpe_opsramp_device_group.device_group_root"),
						testAccEnsureDeviceGroupExists(t, "hpe_opsramp_device_group.device_group_resources"),
						testAccEnsureDeviceGroupExists(t, "hpe_opsramp_device_group.device_group_query"),
						testAccEnsureDeviceGroupExists(t, "hpe_opsramp_device_group.device_group_mixed"),
						resource.TestCheckResourceAttr("hpe_opsramp_device_group.device_group_root", "name", rootGroup),
						resource.TestCheckResourceAttr("hpe_opsramp_device_group.device_group_resources", "name", childRes),
						resource.TestCheckResourceAttr("hpe_opsramp_device_group.device_group_query", "name", childQuery),
						resource.TestCheckResourceAttr("hpe_opsramp_device_group.device_group_mixed", "name", childMixed),
					),
				},
			},
		})
	})
}

func testAccE2ESimpleDeviceGroupConfig(
	res1, res2, res3, rootGroup, childRes, childQuery,
	childMixed, clientOverride string,
) string {
	clientAttr := acctest.ClientAttrHCL(clientOverride)

	return fmt.Sprintf(
		`
%s
resource "hpe_opsramp_resource" "resource1" {
	resource_name = "%s"
	resource_type = "Linux"
	%s
}

resource "hpe_opsramp_resource" "resource2" {
	resource_name = "%s"
	resource_type = "Linux"
	%s
}

resource "hpe_opsramp_resource" "resource3" {
	resource_name = "%s"
	resource_type = "Linux"
	%s
}

resource "hpe_opsramp_device_group" "device_group_root" {
	name      = "%s"
	resources = []
	%s
}

resource "hpe_opsramp_device_group" "device_group_resources" {
	parent_id = hpe_opsramp_device_group.device_group_root.id
	name      = "%s"
	resources = [hpe_opsramp_resource.resource1.uuid]
	%s
}

resource "hpe_opsramp_device_group" "device_group_query" {
	parent_id    = hpe_opsramp_device_group.device_group_root.id
	name         = "%s"
	search_query = format("resourceType = \"Linux\" AND uuid = \"%%s\"", hpe_opsramp_resource.resource2.uuid)
	%s
}

resource "hpe_opsramp_device_group" "device_group_mixed" {
	parent_id    = hpe_opsramp_device_group.device_group_root.id
	name         = "%s"
	search_query = format("resourceType = \"Linux\" AND uuid = \"%%s\"", hpe_opsramp_resource.resource2.uuid)
	resources    = [hpe_opsramp_resource.resource3.uuid]
	%s
}
`,
		acctest.ProviderConfigHCL(),
		res1, clientAttr,
		res2, clientAttr,
		res3, clientAttr,
		rootGroup, clientAttr,
		childRes, clientAttr,
		childQuery, clientAttr,
		childMixed, clientAttr,
	)
}
