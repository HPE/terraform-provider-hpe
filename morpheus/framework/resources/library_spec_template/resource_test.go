package library_spec_template_test

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
	code := testhelpers.TestMain(m)
	testhelpers.WriteMergedResults()
	os.Exit(code)
}

func TestAccMorpheusLibrarySpecTemplateBasic(t *testing.T) {
	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	testSystem := systemoverride.GetPreferred(t, "feature")
	providerConfig := testhelpers.ProviderBlockForServer(testSystem)

	rName := acctest.RandomWithPrefix("tf-acc-spectempl")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + testAccLibrarySpecTemplateConfig(rName, "{}"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("hpe_morpheus_library_spec_template.test", "id"),
					resource.TestCheckResourceAttr("hpe_morpheus_library_spec_template.test", "name", rName),
					resource.TestCheckResourceAttr("hpe_morpheus_library_spec_template.test", "type", "terraform"),
					resource.TestCheckResourceAttr("hpe_morpheus_library_spec_template.test", "content", "{}"),
				),
			},
			{
				ResourceName:      "hpe_morpheus_library_spec_template.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccMorpheusLibrarySpecTemplateUpdate(t *testing.T) {
	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	testSystem := systemoverride.GetPreferred(t, "feature")
	providerConfig := testhelpers.ProviderBlockForServer(testSystem)

	rName := acctest.RandomWithPrefix("tf-acc-spectempl")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + testAccLibrarySpecTemplateConfig(rName, "{}"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hpe_morpheus_library_spec_template.test", "content", "{}"),
				),
			},
			{
				Config: providerConfig + testAccLibrarySpecTemplateConfig(rName, `{"updated": true}`),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hpe_morpheus_library_spec_template.test", "content", `{"updated": true}`),
				),
			},
		},
	})
}

func testAccLibrarySpecTemplateConfig(name, content string) string {
	return fmt.Sprintf(`
resource "hpe_morpheus_library_spec_template" "test" {
  name    = %q
  type    = "terraform"
  source  = "local"
  content = %q
}
`, name, content)
}
