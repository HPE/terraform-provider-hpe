package provisioning_license_test

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

func TestAccMorpheusProvisioningLicenseResourceBasic(t *testing.T) {
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
				Config: providerConfig + testAccProvisioningLicenseConfig(rName, ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("hpe_morpheus_provisioning_license.test", "id"),
					resource.TestCheckResourceAttr("hpe_morpheus_provisioning_license.test", "name", rName),
					resource.TestCheckResourceAttr("hpe_morpheus_provisioning_license.test", "license_type", "win"),
				),
			},
			{
				ResourceName:            "hpe_morpheus_provisioning_license.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"license_key"},
			},
		},
	})
}

func TestAccMorpheusProvisioningLicenseResourceUpdate(t *testing.T) {
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
				Config: providerConfig + testAccProvisioningLicenseConfig(rName, ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("hpe_morpheus_provisioning_license.test", "id"),
				),
			},
			{
				Config: providerConfig + testAccProvisioningLicenseConfig(rName, "updated description"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hpe_morpheus_provisioning_license.test", "description", "updated description"),
				),
			},
		},
	})
}

func testAccProvisioningLicenseConfig(name, description string) string {
	desc := ""
	if description != "" {
		desc = fmt.Sprintf(`description = %q`, description)
	}

	return fmt.Sprintf(`
resource "hpe_morpheus_provisioning_license" "test" {
  name         = %q
  license_type = "win"
  license_key  = "XXXXX-XXXXX"
  %s
}
`, name, desc)
}
