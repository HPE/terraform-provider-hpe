// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package e2e_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/opsramp/acctest"
)

// TestAccE2ESimpleResources exercises the simple-resources e2e scenario:
// creating multiple resources with different identification methods.
func TestAccE2ESimpleResources(t *testing.T) {
	acctest.SkipIfNotClient(t)
	clientOverride := acctest.OptionalClientOverride(t)

	t.Run("multiple resources", func(t *testing.T) {
		res1Name := acctest.RandomName("e2e-res1")
		res2Name := acctest.RandomName("e2e-res2")
		res2Host := acctest.RandomName("e2e-host")

		resource.ParallelTest(t, resource.TestCase{
			PreCheck:                 acctest.PreCheck(t),
			ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories(),
			Steps: []resource.TestStep{
				{
					Config: testAccE2ESimpleResourcesConfig(res1Name, res2Name, res2Host, clientOverride),
					Check: resource.ComposeAggregateTestCheckFunc(
						testAccEnsureResourceExists(t, "hpe_opsramp_resource.resource1"),
						testAccEnsureResourceExists(t, "hpe_opsramp_resource.resource2"),
						resource.TestCheckResourceAttr("hpe_opsramp_resource.resource1", "resource_name", res1Name),
						resource.TestCheckResourceAttr("hpe_opsramp_resource.resource2", "hostname", res2Host),
					),
				},
			},
		})
	})
}

func testAccE2ESimpleResourcesConfig(res1Name string, res2Name string, res2Host string, clientOverride string) string {
	clientAttr := acctest.ClientAttrHCL(clientOverride)

	return fmt.Sprintf(`
%s
resource "hpe_opsramp_resource" "resource1" {
	alias_name    = "%s"
	resource_name = "%s"
	resource_type = "Other"
	%s
}

resource "hpe_opsramp_resource" "resource2" {
	alias_name    = "%s"
	hostname      = "%s"
	resource_type = "Other"
	%s
}
`, acctest.ProviderConfigHCL(), res1Name, res1Name, clientAttr, res2Name, res2Host, clientAttr)
}
