package storage_volume_test

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

func TestAccMorpheusStorageVolumeResourceBasic(t *testing.T) {
	if capabilities.Missing(t, capabilities.Alletra) {
		t.Log("Skipping test due to missing capabilities")

		return
	}
	t.Skip("Skipping: requires pre-existing storage server infrastructure")

	providerConfig := testhelpers.ProviderBlock()

	rName := acctest.RandomWithPrefix(t.Name())
	rNameUpdated := acctest.RandomWithPrefix(t.Name() + "-updated")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			// Create
			{
				Config: providerConfig + testAccStorageVolumeConfig(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("hpe_morpheus_storage_volume.test", "id"),
					resource.TestCheckResourceAttr("hpe_morpheus_storage_volume.test", "name", rName),
				),
			},
			// ImportState
			{
				ResourceName:      "hpe_morpheus_storage_volume.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update name
			{
				Config: providerConfig + testAccStorageVolumeConfig(rNameUpdated),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hpe_morpheus_storage_volume.test", "name", rNameUpdated),
				),
			},
		},
	})
}

func testAccStorageVolumeConfig(name string) string {
	return fmt.Sprintf(`
resource "hpe_morpheus_storage_volume" "test" {
  name = %q
}
`, name)
}
