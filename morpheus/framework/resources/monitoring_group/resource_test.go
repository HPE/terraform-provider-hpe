package monitoring_group_test

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

func TestAccMorpheusMonitoringGroupBasic(t *testing.T) {
	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	testSystem := systemoverride.GetPreferred(t, "feature")
	providerConfig := testhelpers.ProviderBlockForServer(testSystem)

	rName := acctest.RandomWithPrefix("tf-acc-mongroup")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + testAccMonitoringGroupConfig(rName, ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("hpe_morpheus_monitoring_group.test", "id"),
					resource.TestCheckResourceAttr("hpe_morpheus_monitoring_group.test", "name", rName),
				),
			},
			{
				ResourceName:            "hpe_morpheus_monitoring_group.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"active"},
			},
		},
	})
}

func TestAccMorpheusMonitoringGroupUpdate(t *testing.T) {
	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	testSystem := systemoverride.GetPreferred(t, "feature")
	providerConfig := testhelpers.ProviderBlockForServer(testSystem)

	rName := acctest.RandomWithPrefix("tf-acc-mongroup")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + testAccMonitoringGroupConfig(rName, ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("hpe_morpheus_monitoring_group.test", "id"),
				),
			},
			{
				Config: providerConfig + testAccMonitoringGroupConfig(rName, "updated description"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hpe_morpheus_monitoring_group.test", "description", "updated description"),
				),
			},
		},
	})
}

func testAccMonitoringGroupConfig(name, description string) string {
	desc := ""
	if description != "" {
		desc = fmt.Sprintf(`description = %q`, description)
	}

	return fmt.Sprintf(`
resource "hpe_morpheus_monitoring_group" "test" {
  name     = %q
  severity = "critical"
  %s
}
`, name, desc)
}
