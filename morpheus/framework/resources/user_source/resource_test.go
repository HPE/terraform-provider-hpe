package user_source_test

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

func TestAccMorpheusUserSourceBasic(t *testing.T) {
	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Skip("requires LDAP infrastructure")

	t.Parallel()

	testSystem := systemoverride.GetPreferred(t, "zodiac")
	providerConfig := testhelpers.ProviderBlockForServer(testSystem)

	rName := acctest.RandomWithPrefix("tf-acc-test-usersource")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + testAccUserSourceConfig(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("hpe_morpheus_user_source.test", "id"),
					resource.TestCheckResourceAttr("hpe_morpheus_user_source.test", "name", rName),
					resource.TestCheckResourceAttr("hpe_morpheus_user_source.test", "type", "ldap"),
				),
			},
			{
				ResourceName:            "hpe_morpheus_user_source.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"account_id"},
			},
		},
	})
}

func testAccUserSourceConfig(name string) string {
	return fmt.Sprintf(`
resource "hpe_morpheus_user_source" "test" {
  name       = %q
  type       = "ldap"
  account_id = 1
}
`, name)
}
