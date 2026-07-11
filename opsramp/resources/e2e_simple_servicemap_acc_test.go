// SPDX-FileCopyrightText: Copyright Hewlett Packard Enterprise Development LP
// SPDX-License-Identifier: Apache-2.0

package resources_test

import (
	"fmt"
	"testing"

	"github.com/HPE/terraform-provider-hpe/opsramp/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccE2ESimpleServicemap exercises the simple-servicemap e2e scenario:
// root servicemap, child services, resource-based children, and links.
func TestAccE2ESimpleServicemap(t *testing.T) {
	t.Run("full hierarchy with links", func(t *testing.T) {
		res1 := acctest.RandomName("sm-res1")
		res2 := acctest.RandomName("sm-res2")
		rootName := acctest.RandomName("sm-root")
		child1 := acctest.RandomName("sm-child1")
		child2 := acctest.RandomName("sm-child2")
		child21 := acctest.RandomName("sm-child21")
		child22 := acctest.RandomName("sm-child22")
		linkedRoot := acctest.RandomName("sm-linked")

		resource.ParallelTest(t, resource.TestCase{
			PreCheck:                 acctest.PreCheck(t),
			ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories(),
			CheckDestroy:             testAccCheckServicemapDestroy(t),
			Steps: []resource.TestStep{
				{
					Config: testAccE2ESimpleServicemapConfig(res1, res2, rootName, child1, child2, child21, child22, linkedRoot),
					Check: resource.ComposeAggregateTestCheckFunc(
						testAccEnsureServicemapExists(t, "hpe_opsramp_servicemap.servicemap_root"),
						testAccEnsureServicemapExists(t, "hpe_opsramp_servicemap.servicemap_child1"),
						testAccEnsureServicemapExists(t, "hpe_opsramp_servicemap.servicemap_child2"),
						testAccEnsureServicemapExists(t, "hpe_opsramp_servicemap.servicemap_child21"),
						testAccEnsureServicemapExists(t, "hpe_opsramp_servicemap.servicemap_child22"),
						testAccEnsureServicemapExists(t, "hpe_opsramp_servicemap.servicemap_linked_root"),
						resource.TestCheckResourceAttrPair(
							"hpe_opsramp_servicemap_link.servicemap_link", "parent",
							"hpe_opsramp_servicemap.servicemap_root", "id",
						),
						resource.TestCheckResourceAttrPair(
							"hpe_opsramp_servicemap_link.servicemap_link", "link",
							"hpe_opsramp_servicemap.servicemap_linked_root", "id",
						),
						resource.TestCheckResourceAttr("hpe_opsramp_servicemap.servicemap_root", "name", rootName),
						resource.TestCheckResourceAttr("hpe_opsramp_servicemap.servicemap_linked_root", "name", linkedRoot),
					),
				},
			},
		})
	})
}

func testAccE2ESimpleServicemapConfig(res1, res2, rootName, child1, child2, child21, child22, linkedRoot string) string {
	return fmt.Sprintf(`
%s
resource "hpe_opsramp_resource" "resource1" {
	resource_name = "%s"
	resource_type = "Linux"
}

resource "hpe_opsramp_resource" "resource2" {
	resource_name = "%s"
	resource_type = "Linux"
}

resource "hpe_opsramp_servicemap" "servicemap_root" {
	name = "%s"
	type = "Service"
}

resource "hpe_opsramp_servicemap" "servicemap_child1" {
	name   = "%s"
	type   = "Service"
	parent = hpe_opsramp_servicemap.servicemap_root.id
}

resource "hpe_opsramp_servicemap" "servicemap_child2" {
	name   = "%s"
	type   = "Service"
	parent = hpe_opsramp_servicemap.servicemap_root.id
}

resource "hpe_opsramp_servicemap" "servicemap_child21" {
	name      = "%s"
	type      = "Resource"
	parent    = hpe_opsramp_servicemap.servicemap_child2.id
	resources = [hpe_opsramp_resource.resource1.uuid]
}

resource "hpe_opsramp_servicemap" "servicemap_child22" {
	name         = "%s"
	type         = "Resource"
	parent       = hpe_opsramp_servicemap.servicemap_child2.id
	search_query = "resourceType = \"Server\" AND name CONTAINS \"Test\""
}

resource "hpe_opsramp_servicemap" "servicemap_linked_root" {
	name = "%s"
	type = "Service"
}

resource "hpe_opsramp_servicemap_link" "servicemap_link" {
	parent = hpe_opsramp_servicemap.servicemap_root.id
	link   = hpe_opsramp_servicemap.servicemap_linked_root.id
}
`, acctest.ProviderConfigHCL(), res1, res2, rootName, child1, child2, child21, child22, linkedRoot)
}
