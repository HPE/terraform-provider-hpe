package network_pool_server_test

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

func TestAccMorpheusNetworkPoolServerBasic(t *testing.T) {
	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	testSystem := systemoverride.GetPreferred(t, "feature")
	providerConfig := testhelpers.ProviderBlockForServer(testSystem)

	name := acctest.RandomWithPrefix("tf-acc-pool-srv")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + fmt.Sprintf(`
resource "hpe_morpheus_network_pool_server" "test" {
  name             = %q
  type_id          = 4
  enabled          = true
  service_url      = "http://localhost:8080"
  service_username = "admin"
  service_password = "password123"
}
`, name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hpe_morpheus_network_pool_server.test", "name", name),
					resource.TestCheckResourceAttr("hpe_morpheus_network_pool_server.test", "type_id", "4"),
					resource.TestCheckResourceAttr("hpe_morpheus_network_pool_server.test", "enabled", "true"),
					resource.TestCheckResourceAttrSet("hpe_morpheus_network_pool_server.test", "id"),
				),
			},
			{
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"service_username", "service_password"},
				ResourceName:            "hpe_morpheus_network_pool_server.test",
			},
		},
	})
}

func TestAccMorpheusNetworkPoolServerUpdate(t *testing.T) {
	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	testSystem := systemoverride.GetPreferred(t, "feature")
	providerConfig := testhelpers.ProviderBlockForServer(testSystem)

	name := acctest.RandomWithPrefix("tf-acc-pool-srv")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + fmt.Sprintf(`
resource "hpe_morpheus_network_pool_server" "test" {
  name             = %q
  type_id          = 4
  enabled          = true
  service_url      = "http://localhost:8080"
  service_username = "admin"
  service_password = "password123"
}
`, name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hpe_morpheus_network_pool_server.test", "enabled", "true"),
				),
			},
			{
				Config: providerConfig + fmt.Sprintf(`
resource "hpe_morpheus_network_pool_server" "test" {
  name             = %q
  type_id          = 4
  enabled          = false
  service_url      = "http://localhost:8080"
  service_username = "admin"
  service_password = "password123"
}
`, name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hpe_morpheus_network_pool_server.test", "enabled", "false"),
				),
			},
		},
	})
}
