package backup_test

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

func TestAccMorpheusBackupBasic(t *testing.T) {
	defer testhelpers.RecordResult(t)
	t.Skip("requires instance infrastructure to create a backup")

	t.Parallel()

	testSystem := systemoverride.GetPreferred(t, "feature")
	providerConfig := testhelpers.ProviderBlockForServer(testSystem)

	rName := acctest.RandomWithPrefix("tf-acc-backup")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + testAccBackupConfig(rName, true),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("hpe_morpheus_backup.test", "id"),
					resource.TestCheckResourceAttr("hpe_morpheus_backup.test", "name", rName),
					resource.TestCheckResourceAttr("hpe_morpheus_backup.test", "enabled", "true"),
				),
			},
			{
				ResourceName:      "hpe_morpheus_backup.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccMorpheusBackupUpdate(t *testing.T) {
	defer testhelpers.RecordResult(t)
	t.Skip("requires instance infrastructure to create a backup")

	t.Parallel()

	testSystem := systemoverride.GetPreferred(t, "feature")
	providerConfig := testhelpers.ProviderBlockForServer(testSystem)

	rName := acctest.RandomWithPrefix("tf-acc-backup")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + testAccBackupConfig(rName, true),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hpe_morpheus_backup.test", "enabled", "true"),
				),
			},
			{
				Config: providerConfig + testAccBackupConfig(rName, false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hpe_morpheus_backup.test", "enabled", "false"),
				),
			},
		},
	})
}

func testAccBackupConfig(name string, enabled bool) string {
	return fmt.Sprintf(`
resource "hpe_morpheus_backup" "test" {
  name    = %q
  enabled = %t
}
`, name, enabled)
}
