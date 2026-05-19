package scale_threshold_test

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

func TestAccScaleThresholdResource_basic(t *testing.T) {
	testSystem := systemoverride.GetPreferred(t, "feature")
	providerConfig := testhelpers.ProviderBlockForServer(testSystem)

	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandStringFromCharSet(5, acctest.CharSetAlphaNum))

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			// Create
			{
				Config: providerConfig + testAccScaleThresholdConfig(rName, true, 1, 3),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("hpe_morpheus_scale_threshold.test", "id"),
					resource.TestCheckResourceAttr("hpe_morpheus_scale_threshold.test", "name", rName),
					resource.TestCheckResourceAttr("hpe_morpheus_scale_threshold.test", "auto_upscale", "true"),
					resource.TestCheckResourceAttr("hpe_morpheus_scale_threshold.test", "min_count", "1"),
					resource.TestCheckResourceAttr("hpe_morpheus_scale_threshold.test", "max_count", "3"),
				),
			},
			// ImportState
			{
				ResourceName:      "hpe_morpheus_scale_threshold.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update max_count
			{
				Config: providerConfig + testAccScaleThresholdConfig(rName, true, 1, 5),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hpe_morpheus_scale_threshold.test", "max_count", "5"),
				),
			},
		},
	})
}

func testAccScaleThresholdConfig(name string, autoUpscale bool, minCount, maxCount int) string {
	return fmt.Sprintf(`
resource "hpe_morpheus_scale_threshold" "test" {
  name          = %q
  auto_upscale  = %t
  min_count     = %d
  max_count     = %d
}
`, name, autoUpscale, minCount, maxCount)
}
