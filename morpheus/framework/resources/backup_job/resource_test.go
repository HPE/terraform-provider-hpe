package backup_job_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
)

func TestMain(m *testing.M) {
	code := testhelpers.TestMain(m)
	testhelpers.WriteMergedResults()
	os.Exit(code)
}

func TestAccMorpheusBackupJobResourceBasic(t *testing.T) {
	if capabilities.Missing(t, capabilities.All) {
		t.Log("Skipping test due to missing capabilities")

		return
	}
	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()

	rName := acctest.RandomWithPrefix(t.Name())
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + testAccBackupJobConfig(rName, true),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("hpe_morpheus_backup_job.test", "id"),
					resource.TestCheckResourceAttr("hpe_morpheus_backup_job.test", "name", rName),
					resource.TestCheckResourceAttr("hpe_morpheus_backup_job.test", "enabled", "true"),
				),
			},
			{
				ResourceName:            "hpe_morpheus_backup_job.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"enabled"},
			},
		},
	})
}

func TestAccMorpheusBackupJobResourceUpdate(t *testing.T) {
	if capabilities.Missing(t, capabilities.All) {
		t.Log("Skipping test due to missing capabilities")

		return
	}
	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock()

	rName := acctest.RandomWithPrefix(t.Name())
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + testAccBackupJobConfig(rName, true),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hpe_morpheus_backup_job.test", "enabled", "true"),
				),
			},
			{
				Config: providerConfig + testAccBackupJobConfig(rName, false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hpe_morpheus_backup_job.test", "enabled", "false"),
				),
			},
		},
	})
}

func testAccBackupJobConfig(name string, enabled bool) string {
	return fmt.Sprintf(`
resource "hpe_morpheus_backup_job" "test" {
  name    = %q
  enabled = %t
}
`, name, enabled)
}
