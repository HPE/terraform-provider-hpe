package job_test

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

func TestAccJobResource_basic(t *testing.T) {
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
				Config: providerConfig + testAccJobConfig(rName, true),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("hpe_morpheus_job.test", "id"),
					resource.TestCheckResourceAttr("hpe_morpheus_job.test", "name", rName),
					resource.TestCheckResourceAttr("hpe_morpheus_job.test", "enabled", "true"),
				),
			},
			// ImportState
			{
				ResourceName:      "hpe_morpheus_job.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update enabled
			{
				Config: providerConfig + testAccJobConfig(rName, false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hpe_morpheus_job.test", "enabled", "false"),
				),
			},
		},
	})
}

func testAccJobConfig(name string, enabled bool) string {
	return fmt.Sprintf(`
resource "hpe_morpheus_workflow" "test" {
  name       = "%s-wf"
  type       = "operation"
  visibility = "private"
}

resource "hpe_morpheus_job" "test" {
  name          = %q
  enabled       = %t
  schedule_mode = "manual"
  target_type   = "appliance"
  workflow_id   = hpe_morpheus_workflow.test.id
}
`, name, name, enabled)
}
