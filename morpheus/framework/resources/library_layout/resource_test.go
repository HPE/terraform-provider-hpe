package library_layout_test

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

func TestAccMorpheusLibraryLayoutBasic(t *testing.T) {
	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	testSystem := systemoverride.GetPreferred(t, "feature")
	providerConfig := testhelpers.ProviderBlockForServer(testSystem)

	rName := acctest.RandomWithPrefix("tf-acc-layout")
	itName := acctest.RandomWithPrefix("tf-acc-insttype")
	itCode := acctest.RandomWithPrefix("tf-acc-insttype")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + testAccLibraryLayoutConfig(itName, itCode, rName, "1.0", ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("hpe_morpheus_library_layout.test", "id"),
					resource.TestCheckResourceAttr("hpe_morpheus_library_layout.test", "name", rName),
					resource.TestCheckResourceAttr("hpe_morpheus_library_layout.test", "instance_version", "1.0"),
					resource.TestCheckResourceAttr("hpe_morpheus_library_layout.test", "provision_type_code", "docker"),
				),
			},
			{
				Config: providerConfig + testAccLibraryLayoutConfig(itName, itCode, rName, "2.0", "updated description"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hpe_morpheus_library_layout.test", "instance_version", "2.0"),
					resource.TestCheckResourceAttr("hpe_morpheus_library_layout.test", "description", "updated description"),
				),
			},
		},
	})
}

func testAccLibraryLayoutConfig(itName, itCode, name, version, description string) string {
	desc := ""
	if description != "" {
		desc = fmt.Sprintf(`description = %q`, description)
	}
	return fmt.Sprintf(`
resource "hpe_morpheus_library_instance_type" "test" {
  name       = %q
  code       = %q
  category   = "os"
  visibility = "private"
}

resource "hpe_morpheus_library_layout" "test" {
  instance_type_id    = hpe_morpheus_library_instance_type.test.id
  name                = %q
  instance_version    = %q
  provision_type_code = "docker"
  %s
}
`, itName, itCode, name, version, desc)
}
