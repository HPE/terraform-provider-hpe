package vdi_pool_test

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

func TestAccMorpheusVdiPoolResourceBasic(t *testing.T) {
	if capabilities.Missing(t, capabilities.VDI) {
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
				Config: providerConfig + testAccVdiPoolConfig(rName, ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("hpe_morpheus_vdi_pool.test", "id"),
					resource.TestCheckResourceAttr("hpe_morpheus_vdi_pool.test", "name", rName),
					resource.TestCheckResourceAttr("hpe_morpheus_vdi_pool.test", "enabled", "true"),
				),
			},
			{
				ResourceName:      "hpe_morpheus_vdi_pool.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccMorpheusVdiPoolResourceUpdate(t *testing.T) {
	if capabilities.Missing(t, capabilities.VDI) {
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
				Config: providerConfig + testAccVdiPoolConfig(rName, ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("hpe_morpheus_vdi_pool.test", "id"),
				),
			},
			{
				Config: providerConfig + testAccVdiPoolConfig(rName, "updated description"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hpe_morpheus_vdi_pool.test", "description", "updated description"),
				),
			},
		},
	})
}

func testAccVdiPoolConfig(name, description string) string {
	desc := ""
	if description != "" {
		desc = fmt.Sprintf(`description = %q`, description)
	}

	return fmt.Sprintf(`
resource "hpe_morpheus_vdi_pool" "test" {
  name          = %q
  enabled       = true
  max_pool_size = 1
  %s
}
`, name, desc)
}
