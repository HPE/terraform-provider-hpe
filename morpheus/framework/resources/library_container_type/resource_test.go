package library_container_type_test

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
	code := m.Run()
	testhelpers.WriteMergedResults()
	os.Exit(code)
}

func TestAccMorpheusLibraryContainerTypeBasic(t *testing.T) {
	if capabilities.Missing(t, capabilities.All, capabilities.Docker) {
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
	rShort := acctest.RandomWithPrefix("tfaccnt")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + testAccLibraryContainerTypeConfig(rName, rShort, "1.0", ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("hpe_morpheus_library_container_type.test", "id"),
					resource.TestCheckResourceAttr("hpe_morpheus_library_container_type.test", "name", rName),
					resource.TestCheckResourceAttr("hpe_morpheus_library_container_type.test", "short_name", rShort),
					resource.TestCheckResourceAttr("hpe_morpheus_library_container_type.test", "container_version", "1.0"),
					resource.TestCheckResourceAttr("hpe_morpheus_library_container_type.test", "provision_type_code", "docker"),
				),
			},
			{
				Config: providerConfig + testAccLibraryContainerTypeConfig(rName, rShort, "2.0", "updated description"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hpe_morpheus_library_container_type.test", "container_version", "2.0"),
					resource.TestCheckResourceAttr("hpe_morpheus_library_container_type.test", "description", "updated description"),
				),
			},
		},
	})
}

func testAccLibraryContainerTypeConfig(name, shortName, version, description string) string {
	desc := ""
	if description != "" {
		desc = fmt.Sprintf(`description = %q`, description)
	}

	return fmt.Sprintf(`
resource "hpe_morpheus_library_container_type" "test" {
  name                = %q
  short_name          = %q
  container_version   = %q
  provision_type_code = "docker"
  %s
}
`, name, shortName, version, desc)
}
