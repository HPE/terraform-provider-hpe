package subnet_test

import (
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

func TestAccMorpheusSubnetResourceBasic(t *testing.T) {
	if capabilities.Missing(t, capabilities.Subnet) {
		t.Log("Subnet tests require a configured cloud with networks - skipping in environment without infrastructure")

		return
	}
	defer testhelpers.RecordResult(t)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	providerConfig := testhelpers.ProviderBlock() //nolint:staticcheck // used after skip

	name := acctest.RandomWithPrefix(t.Name())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, morpheus.New(), nil),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
resource "hpe_morpheus_subnet" "test" {
	name       = "` + name + `"
  type_id    = 1
  visibility = "private"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hpe_morpheus_subnet.test", "name", name),
					resource.TestCheckResourceAttr("hpe_morpheus_subnet.test", "type_id", "1"),
					resource.TestCheckResourceAttr("hpe_morpheus_subnet.test", "visibility", "private"),
					resource.TestCheckResourceAttrSet("hpe_morpheus_subnet.test", "id"),
				),
			},
			{
				ImportState:       true,
				ImportStateVerify: true,
				ResourceName:      "hpe_morpheus_subnet.test",
			},
		},
	})
}
