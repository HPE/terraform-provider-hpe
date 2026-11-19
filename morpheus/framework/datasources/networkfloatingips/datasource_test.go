// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkfloatingips_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

const dataSourceName = "data.hpe_morpheus_network_floating_ips.example"

func TestMain(m *testing.M) {
	code := testhelpers.TestMain(m)
	testhelpers.WriteMergedResults()
	os.Exit(code)
}

// TestAccMorpheusListNetworkFloatingIps allocates a floating IP via the
// resource, then lists floating IPs through the data source. It skips without
// the SDN fixture, since the list is only meaningful in an SDN-enabled
// environment where floating IPs exist.
func TestAccMorpheusListNetworkFloatingIps(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.JuniperSDN)
	f := testhelpers.RequireSdnFixture(t)
	testhelpers.LockFloatingIpPool(t, f.FloatingIpPoolID)

	t.Parallel()

	config := testhelpers.ProviderBlock() + fmt.Sprintf(`
resource "hpe_morpheus_network_floating_ip" "example" {
  network_server_id   = %s
  floating_ip_pool_id = %s
}

data "hpe_morpheus_network_floating_ips" "example" {
  depends_on = [hpe_morpheus_network_floating_ip.example]
}
`, f.NetworkServerID, f.FloatingIpPoolID)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "network_floating_ips.#"),
					resource.TestCheckTypeSetElemAttrPair(
						dataSourceName, "network_floating_ips.*.id",
						"hpe_morpheus_network_floating_ip.example", "id",
					),
				),
			},
		},
	})
}
