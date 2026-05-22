package storage_server_test

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

func TestAccStorageServerResource_basic(t *testing.T) {
	t.Skip("Skipping: requires external storage infrastructure (no standalone storage server types available)")

	testSystem := systemoverride.GetPreferred(t, "zodiac")
	providerConfig := testhelpers.ProviderBlockForServer(testSystem)

	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandStringFromCharSet(5, acctest.CharSetAlphaNum))

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			// Create
			{
				Config: providerConfig + testAccStorageServerConfig(rName, "local", ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("hpe_morpheus_storage_server.test", "id"),
					resource.TestCheckResourceAttr("hpe_morpheus_storage_server.test", "name", rName),
					resource.TestCheckResourceAttr("hpe_morpheus_storage_server.test", "type", "local"),
				),
			},
			// ImportState
			{
				ResourceName:      "hpe_morpheus_storage_server.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update description
			{
				Config: providerConfig + testAccStorageServerConfig(rName, "local", "updated description"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hpe_morpheus_storage_server.test", "description", "updated description"),
				),
			},
		},
	})
}

func testAccStorageServerConfig(name, serverType, description string) string {
	desc := ""
	if description != "" {
		desc = fmt.Sprintf(`  description = %q`, description)
	}

	return fmt.Sprintf(`
resource "hpe_morpheus_storage_server" "test" {
  name = %q
  type = %q
%s
}
`, name, serverType, desc)
}
