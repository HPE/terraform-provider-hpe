// SPDX-FileCopyrightText: Copyright Hewlett Packard Enterprise Development LP
// SPDX-License-Identifier: Apache-2.0

package resources_test

import (
	"fmt"
	"testing"

	"github.com/HPE/terraform-provider-opsramp/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccServicemapLinkResource(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		rootName := acctest.RandomName("sm-link-root")
		linkedName := acctest.RandomName("sm-link-target")

		resource.ParallelTest(t, resource.TestCase{
			PreCheck:                 acctest.PreCheck(t),
			ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories(),
			Steps: []resource.TestStep{
				{
					Config: testAccServicemapLinkConfig(rootName, linkedName),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttrSet("opsramp_servicemap_link.test_link", "parent"),
						resource.TestCheckResourceAttrSet("opsramp_servicemap_link.test_link", "link"),
						resource.TestCheckResourceAttrPair(
							"opsramp_servicemap_link.test_link", "parent",
							"opsramp_servicemap.test_link_root", "id",
						),
						resource.TestCheckResourceAttrPair(
							"opsramp_servicemap_link.test_link", "link",
							"opsramp_servicemap.test_link_target", "id",
						),
					),
				},
			},
		})
	})
}

func testAccServicemapLinkConfig(rootName string, linkedName string) string {
	return fmt.Sprintf(`
%s
resource "opsramp_servicemap" "test_link_root" {
	name = "%s"
	type = "Service"
}

resource "opsramp_servicemap" "test_link_target" {
	name = "%s"
	type = "Service"
}

resource "opsramp_servicemap_link" "test_link" {
	parent = opsramp_servicemap.test_link_root.id
	link   = opsramp_servicemap.test_link_target.id
}
`, acctest.ProviderConfigHCL(), rootName, linkedName)
}
