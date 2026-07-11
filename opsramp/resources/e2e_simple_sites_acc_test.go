// SPDX-FileCopyrightText: Copyright Hewlett Packard Enterprise Development LP
// SPDX-License-Identifier: Apache-2.0

package resources_test

import (
	"fmt"
	"testing"

	"github.com/HPE/terraform-provider-hpe/opsramp/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccE2ESimpleSites exercises the simple-sites e2e scenario:
// root site with child sites using resources, queries, and mixed approaches.
func TestAccE2ESimpleSites(t *testing.T) {
	t.Run("full hierarchy", func(t *testing.T) {
		res1 := acctest.RandomName("site-res1")
		res2 := acctest.RandomName("site-res2")
		res3 := acctest.RandomName("site-res3")
		rootSite := acctest.RandomName("site-root")
		childValencia := acctest.RandomName("site-vlc")
		childMadrid := acctest.RandomName("site-mad")
		childBarcelona := acctest.RandomName("site-bcn")

		resource.ParallelTest(t, resource.TestCase{
			PreCheck:                 acctest.PreCheck(t),
			ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories(),
			CheckDestroy:             testAccCheckSiteDestroy(t),
			Steps: []resource.TestStep{
				{
					Config: testAccE2ESimpleSitesConfig(res1, res2, res3, rootSite, childValencia, childMadrid, childBarcelona),
					Check: resource.ComposeAggregateTestCheckFunc(
						testAccEnsureSiteExists(t, "hpe_opsramp_site.site_root"),
						testAccEnsureSiteExists(t, "hpe_opsramp_site.site_valencia"),
						testAccEnsureSiteExists(t, "hpe_opsramp_site.site_madrid"),
						testAccEnsureSiteExists(t, "hpe_opsramp_site.site_barcelona"),
						resource.TestCheckResourceAttr("hpe_opsramp_site.site_root", "name", rootSite),
						resource.TestCheckResourceAttr("hpe_opsramp_site.site_valencia", "name", childValencia),
						resource.TestCheckResourceAttr("hpe_opsramp_site.site_valencia", "country", "Spain"),
						resource.TestCheckResourceAttr("hpe_opsramp_site.site_madrid", "name", childMadrid),
						resource.TestCheckResourceAttr("hpe_opsramp_site.site_barcelona", "name", childBarcelona),
					),
				},
			},
		})
	})
}

func testAccE2ESimpleSitesConfig(res1, res2, res3, rootSite, childValencia, childMadrid, childBarcelona string) string {
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

resource "hpe_opsramp_resource" "resource3" {
	resource_name = "%s"
	resource_type = "Linux"
}

resource "hpe_opsramp_site" "site_root" {
	name    = "%s"
	country = "Spain"
}

resource "hpe_opsramp_site" "site_valencia" {
	parent_id = hpe_opsramp_site.site_root.id
	name      = "%s"
	address   = "Av. del General Avilés, 35-37, Benicalap"
	country   = "Spain"
	zip       = "46035"
	state     = "Comunitat Valenciana"
	city      = "València"
	resources = [hpe_opsramp_resource.resource1.uuid]
}

resource "hpe_opsramp_site" "site_madrid" {
	parent_id    = hpe_opsramp_site.site_root.id
	name         = "%s"
	address      = "Calle Vicente Aleixandre, 1"
	country      = "Spain"
	zip          = "28232"
	state        = "Madrid"
	city         = "Las Rozas de Madrid"
	search_query = format("uuid = \"%%s\"", hpe_opsramp_resource.resource2.uuid)
}

resource "hpe_opsramp_site" "site_barcelona" {
	parent_id    = hpe_opsramp_site.site_root.id
	name         = "%s"
	address      = "Carrer de Tànger, 66"
	country      = "Spain"
	zip          = "08018"
	state        = "Barcelona"
	city         = "Sant Martí"
	search_query = format("uuid = \"%%s\"", hpe_opsramp_resource.resource2.uuid)
	resources    = [hpe_opsramp_resource.resource3.uuid]
}
`, acctest.ProviderConfigHCL(), res1, res2, res3, rootSite, childValencia, childMadrid, childBarcelona)
}
