package integration_test

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

func TestAccIntegrationResource_basic(t *testing.T) {
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
				Config: providerConfig + testAccIntegrationConfig(rName, "docker.registry", true),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("hpe_morpheus_integration.test", "id"),
					resource.TestCheckResourceAttr("hpe_morpheus_integration.test", "name", rName),
					resource.TestCheckResourceAttr("hpe_morpheus_integration.test", "type", "docker.registry"),
					resource.TestCheckResourceAttr("hpe_morpheus_integration.test", "enabled", "true"),
				),
			},
			// ImportState
			{
				ResourceName:      "hpe_morpheus_integration.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update enabled
			{
				Config: providerConfig + testAccIntegrationConfig(rName, "docker.registry", false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hpe_morpheus_integration.test", "enabled", "false"),
				),
			},
		},
	})
}

func testAccIntegrationConfig(name, intType string, enabled bool) string {
	return fmt.Sprintf(`
resource "hpe_morpheus_integration" "test" {
  name    = %q
  type    = %q
  enabled = %t
}
`, name, intType, enabled)
}
