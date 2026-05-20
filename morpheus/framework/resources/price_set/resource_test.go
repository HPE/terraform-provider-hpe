package price_set_test

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

func TestAccMorpheusPriceSetBasic(t *testing.T) {
	t.Skip("resource registered in SDKv2 provider, not framework")
	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	testSystem := systemoverride.GetPreferred(t, "feature")
	providerConfig := testhelpers.ProviderBlockForServer(testSystem)

	rName := acctest.RandomWithPrefix("tf-acc-priceset")
	rCode := acctest.RandomWithPrefix("tf-acc-priceset")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + testAccPriceSetConfig(rName, rCode),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("hpe_morpheus_price_set.test", "id"),
					resource.TestCheckResourceAttr("hpe_morpheus_price_set.test", "name", rName),
					resource.TestCheckResourceAttr("hpe_morpheus_price_set.test", "code", rCode),
				),
			},
			{
				ResourceName:            "hpe_morpheus_price_set.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"code", "name", "price_unit", "type"},
			},
		},
	})
}

func TestAccMorpheusPriceSetUpdate(t *testing.T) {
	t.Skip("resource registered in SDKv2 provider, not framework")
	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	testSystem := systemoverride.GetPreferred(t, "feature")
	providerConfig := testhelpers.ProviderBlockForServer(testSystem)

	rName := acctest.RandomWithPrefix("tf-acc-priceset")
	rCode := acctest.RandomWithPrefix("tf-acc-priceset")
	updatedName := rName + "-updated"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + testAccPriceSetConfig(rName, rCode),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hpe_morpheus_price_set.test", "name", rName),
				),
			},
			{
				Config: providerConfig + testAccPriceSetConfig(updatedName, rCode),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hpe_morpheus_price_set.test", "name", updatedName),
				),
			},
		},
	})
}

func testAccPriceSetConfig(name, code string) string {
	return fmt.Sprintf(`
resource "hpe_morpheus_price_set" "test" {
  name       = %q
  code       = %q
  price_unit = "month"
  type       = "component"
}
`, name, code)
}
