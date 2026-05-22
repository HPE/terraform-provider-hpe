package catalog_item_type_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/systemoverride"
)

func TestMain(m *testing.M) {
	systemoverride.ParseFlags()
	code := testhelpers.TestMain(m)
	testhelpers.WriteMergedResults()
	os.Exit(code)
}

func TestAccCatalogItemTypeResource_basic(t *testing.T) {
	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	testSystem := systemoverride.GetPreferred(t, "zodiac")
	providerConfig := testhelpers.ProviderBlockForServer(testSystem)

	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandStringFromCharSet(5, acctest.CharSetAlphaNum))

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			// Create
			{
				Config: providerConfig + testAccCatalogItemTypeConfig(rName, true, ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("hpe_morpheus_catalog_item_type.test", "id"),
					resource.TestCheckResourceAttr("hpe_morpheus_catalog_item_type.test", "name", rName),
					resource.TestCheckResourceAttr("hpe_morpheus_catalog_item_type.test", "enabled", "true"),
				),
			},
			// ImportState
			{
				ResourceName:      "hpe_morpheus_catalog_item_type.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update description
			{
				Config: providerConfig + testAccCatalogItemTypeConfig(rName, true, "updated description"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hpe_morpheus_catalog_item_type.test", "description", "updated description"),
				),
			},
		},
	})
}

func testAccCatalogItemTypeConfig(name string, enabled bool, description string) string {
	desc := ""
	if description != "" {
		desc = fmt.Sprintf(`  description = %q`, description)
	}

	return fmt.Sprintf(`
resource "hpe_morpheus_catalog_item_type" "test" {
  name    = %q
  enabled = %t
%s
}
`, name, enabled, desc)
}
