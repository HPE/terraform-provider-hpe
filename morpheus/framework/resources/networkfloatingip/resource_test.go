// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkfloatingip_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

func TestMain(m *testing.M) {
	code := testhelpers.TestMain(m)
	testhelpers.WriteMergedResults()
	os.Exit(code)
}

// TestAccMorpheusNetworkFloatingIpResourceOk allocates and releases a floating
// IP end to end. It skips without the SDN fixture: the SDN plugin, its
// controller and a floating IP pool exist only on a lab appliance.
func TestAccMorpheusNetworkFloatingIpResourceOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.JuniperSDN)
	f := testhelpers.RequireSdnFixture(t)
	testhelpers.LockFloatingIpPool(t, f.FloatingIpPoolID)

	t.Parallel()

	resourceName := "hpe_morpheus_network_floating_ip.example"

	config := testhelpers.ProviderBlock() + fmt.Sprintf(`
resource "hpe_morpheus_network_floating_ip" "example" {
  network_server_id   = %s
  floating_ip_pool_id = %s
}
`, f.NetworkServerID, f.FloatingIpPoolID)

	checks := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr(resourceName, "network_server_id", f.NetworkServerID),
		resource.TestCheckResourceAttr(resourceName, "floating_ip_pool_id", f.FloatingIpPoolID),
		resource.TestCheckResourceAttrSet(resourceName, "id"),
		resource.TestCheckResourceAttrSet(resourceName, "ip_address"),
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check:  checks,
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: importIDFunc(resourceName, f.NetworkServerID, f.FloatingIpPoolID),
				ImportStateVerify: true,
			},
		},
	})
}

// importIDFunc builds the long-form import ID "<id>:<network_server_id>:
// <floating_ip_pool_id>" so the imported state matches the inputs the API does
// not return.
func importIDFunc(resourceName, networkServerID, poolID string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("resource %s not found in state", resourceName)
		}

		return fmt.Sprintf("%s:%s:%s", rs.Primary.ID, networkServerID, poolID), nil
	}
}
