package subnet_test

import (
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

func TestAccMorpheusSubnetBasic(t *testing.T) {
	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	testSystem := systemoverride.GetPreferred(t, "zodiac")
	providerConfig := testhelpers.ProviderBlockForServer(testSystem) //nolint:staticcheck // used after skip

	_ = acctest.RandomWithPrefix("tf-acc-subnet")

	// Subnet creation requires an existing network which requires a cloud (zone).
	// Skip if no clouds are configured in the test environment.
	t.Skip("Subnet tests require a configured cloud with networks - skipping in environment without infrastructure")

	_ = acctest.RandomWithPrefix("tf-acc-subnet")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
resource "hpe_morpheus_subnet" "test" {
  type_id    = 1
  visibility = "private"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hpe_morpheus_subnet.test", "type_id", "1"),
					resource.TestCheckResourceAttr("hpe_morpheus_subnet.test", "visibility", "private"),
					resource.TestCheckResourceAttrSet("hpe_morpheus_subnet.test", "id"),
				),
			},
			{
				ImportState:       true,
				ImportStateVerify: true,
				ResourceName:      "hpe_morpheus_subnet.test",
			},
		},
	})
}
