// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package resources_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/HPE/terraform-provider-hpe/opsramp/acctest"
)

func TestAccServicemapLinkResource(t *testing.T) {
	acctest.SkipIfNotClient(t)
	clientOverride := acctest.RequireClientScope(t)

	t.Run("create_and_import", func(t *testing.T) {
		rootName := acctest.RandomName("sm-link-root")
		linkedName := acctest.RandomName("sm-link-target")

		resource.ParallelTest(t, resource.TestCase{
			PreCheck:                 acctest.PreCheck(t),
			ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories(),
			Steps: []resource.TestStep{
				{
					Config: testAccServicemapLinkConfig(rootName, linkedName, clientOverride),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttrSet("hpe_opsramp_servicemap_link.test_link", "parent"),
						resource.TestCheckResourceAttrSet("hpe_opsramp_servicemap_link.test_link", "link"),
						resource.TestCheckResourceAttrPair(
							"hpe_opsramp_servicemap_link.test_link", "parent",
							"hpe_opsramp_servicemap.test_link_root", "id",
						),
						resource.TestCheckResourceAttrPair(
							"hpe_opsramp_servicemap_link.test_link", "link",
							"hpe_opsramp_servicemap.test_link_target", "id",
						),
					),
				},
				// ImportState testing — import ID is <parent_id>:<link_id>
				{
					ResourceName: "hpe_opsramp_servicemap_link.test_link",
					ImportState:  true,
					ImportStateIdFunc: testAccServicemapLinkImportStateIdFunc(
						"hpe_opsramp_servicemap_link.test_link",
						clientOverride,
					),
					ImportStateVerify:                    true,
					ImportStateVerifyIdentifierAttribute: "parent",
				},
			},
		})
	})
}

func testAccServicemapLinkConfig(rootName string, linkedName string, clientOverride string) string {
	clientAttr := acctest.ClientAttrHCL(clientOverride)

	return fmt.Sprintf(`
%s
resource "hpe_opsramp_servicemap" "test_link_root" {
	name = "%s"
	type = "Service"
	%s
}

resource "hpe_opsramp_servicemap" "test_link_target" {
	name = "%s"
	type = "Service"
	%s
}

resource "hpe_opsramp_servicemap_link" "test_link" {
	parent = hpe_opsramp_servicemap.test_link_root.id
	link   = hpe_opsramp_servicemap.test_link_target.id
	%s
}
`, acctest.ProviderConfigHCL(), rootName, clientAttr, linkedName, clientAttr, clientAttr)
}

func testAccServicemapLinkImportStateIdFunc(resourceName string, clientOverride string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Not found: %s", resourceName)
		}

		parentID := rs.Primary.Attributes["parent"]
		linkID := rs.Primary.Attributes["link"]

		id := parentID + ":" + linkID
		if clientOverride != "" {
			return clientOverride + ":" + id, nil
		}

		return id, nil
	}
}
