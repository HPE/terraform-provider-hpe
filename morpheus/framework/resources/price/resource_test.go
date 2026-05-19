package price_test

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

func TestAccMorpheusPriceBasic(t *testing.T) {
	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	testSystem := systemoverride.GetPreferred(t, "feature")
	providerConfig := testhelpers.ProviderBlockForServer(testSystem)

	rName := acctest.RandomWithPrefix("tf-acc-price")
	rCode := acctest.RandomWithPrefix("tf-acc-price")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + testAccPriceConfig(rName, rCode, 10.0),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("hpe_morpheus_price.test", "id"),
					resource.TestCheckResourceAttr("hpe_morpheus_price.test", "name", rName),
					resource.TestCheckResourceAttr("hpe_morpheus_price.test", "code", rCode),
					resource.TestCheckResourceAttr("hpe_morpheus_price.test", "cost", "10"),
				),
			},
			{
				ResourceName:      "hpe_morpheus_price.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccMorpheusPriceUpdate(t *testing.T) {
	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	testSystem := systemoverride.GetPreferred(t, "feature")
	providerConfig := testhelpers.ProviderBlockForServer(testSystem)

	rName := acctest.RandomWithPrefix("tf-acc-price")
	rCode := acctest.RandomWithPrefix("tf-acc-price")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + testAccPriceConfig(rName, rCode, 10.0),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hpe_morpheus_price.test", "cost", "10"),
				),
			},
			{
				Config: providerConfig + testAccPriceConfig(rName, rCode, 20.0),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hpe_morpheus_price.test", "cost", "20"),
				),
			},
		},
	})
}

func testAccPriceConfig(name, code string, cost float64) string {
	return fmt.Sprintf(`
resource "hpe_morpheus_price" "test" {
  name       = %q
  code       = %q
  price_type = "fixed"
  price_unit = "month"
  cost       = %g
}
`, name, code, cost)
}
