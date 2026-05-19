package library_instance_type_test

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

func TestAccMorpheusLibraryInstanceTypeBasic(t *testing.T) {
	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	testSystem := systemoverride.GetPreferred(t, "feature")
	providerConfig := testhelpers.ProviderBlockForServer(testSystem)

	rName := acctest.RandomWithPrefix("tf-acc-insttype")
	rCode := acctest.RandomWithPrefix("tf-acc-insttype")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + testAccLibraryInstanceTypeConfig(rName, rCode, ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("hpe_morpheus_library_instance_type.test", "id"),
					resource.TestCheckResourceAttr("hpe_morpheus_library_instance_type.test", "name", rName),
					resource.TestCheckResourceAttr("hpe_morpheus_library_instance_type.test", "code", rCode),
					resource.TestCheckResourceAttr("hpe_morpheus_library_instance_type.test", "visibility", "private"),
				),
			},
			{
				ResourceName:      "hpe_morpheus_library_instance_type.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccMorpheusLibraryInstanceTypeUpdate(t *testing.T) {
	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	testSystem := systemoverride.GetPreferred(t, "feature")
	providerConfig := testhelpers.ProviderBlockForServer(testSystem)

	rName := acctest.RandomWithPrefix("tf-acc-insttype")
	rCode := acctest.RandomWithPrefix("tf-acc-insttype")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + testAccLibraryInstanceTypeConfig(rName, rCode, ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("hpe_morpheus_library_instance_type.test", "id"),
				),
			},
			{
				Config: providerConfig + testAccLibraryInstanceTypeConfig(rName, rCode, "updated description"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hpe_morpheus_library_instance_type.test", "description", "updated description"),
				),
			},
		},
	})
}

func testAccLibraryInstanceTypeConfig(name, code, description string) string {
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
  %s
}
`, name, code, desc)
}
