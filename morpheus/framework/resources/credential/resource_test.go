package credential_test

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
	code := m.Run()
	testhelpers.WriteMergedResults()
	os.Exit(code)
}

func TestAccCredentialResource_basic(t *testing.T) {
	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	testSystem := systemoverride.GetPreferred(t, "feature")
	providerConfig := testhelpers.ProviderBlockForServer(testSystem)

	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandStringFromCharSet(5, acctest.CharSetAlphaNum))

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			// Create
			{
				Config: providerConfig + testAccCredentialConfig(rName, "access-key-secret", ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("hpe_morpheus_credential.test", "id"),
					resource.TestCheckResourceAttr("hpe_morpheus_credential.test", "name", rName),
					resource.TestCheckResourceAttr("hpe_morpheus_credential.test", "type", "access-key-secret"),
				),
			},
			// ImportState (ignore sensitive fields)
			{
				ResourceName:            "hpe_morpheus_credential.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"username", "password"},
			},
			// Update description
			{
				Config: providerConfig + testAccCredentialConfig(rName, "access-key-secret", "updated description"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hpe_morpheus_credential.test", "description", "updated description"),
				),
			},
		},
	})
}

func testAccCredentialConfig(name, credType, description string) string {
	desc := ""
	if description != "" {
		desc = fmt.Sprintf(`  description = %q`, description)
	}

	return fmt.Sprintf(`
resource "hpe_morpheus_credential" "test" {
  name     = %q
  type     = %q
  username = "testaccesskey"
  password = "testsecretkey"
%s
}
`, name, credType, desc)
}
