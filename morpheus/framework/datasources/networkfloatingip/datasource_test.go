// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkfloatingip_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

func TestMain(m *testing.M) {
	code := testhelpers.TestMain(m)
	testhelpers.WriteMergedResults()
	os.Exit(code)
}

// TestAccMorpheusFindNetworkFloatingIpById allocates a floating IP via the
// resource, then reads it back through the data source by id and by address
// (the paged, exact-match lookup). It skips without the SDN fixture, since
// floating IPs only exist in an SDN-enabled environment.
func TestAccMorpheusFindNetworkFloatingIpById(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.JuniperSDN)
	f := testhelpers.RequireSdnFixture(t)
	testhelpers.LockFloatingIpPool(t, f.FloatingIpPoolID)

	t.Parallel()

	dataSourceName := "data.hpe_morpheus_network_floating_ip.example"

	config := testhelpers.ProviderBlock() + fmt.Sprintf(`
resource "hpe_morpheus_network_floating_ip" "example" {
  network_server_id   = %s
  floating_ip_pool_id = %s
}

data "hpe_morpheus_network_floating_ip" "example" {
  id = hpe_morpheus_network_floating_ip.example.id
}

data "hpe_morpheus_network_floating_ip" "by_ip" {
  ip_address = hpe_morpheus_network_floating_ip.example.ip_address
}
`, f.NetworkServerID, f.FloatingIpPoolID)

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttrPair(
			dataSourceName, "id",
			"hpe_morpheus_network_floating_ip.example", "id"),
		resource.TestCheckResourceAttrSet(dataSourceName, "ip_address"),
		resource.TestCheckResourceAttrPair(
			"data.hpe_morpheus_network_floating_ip.by_ip", "id",
			"hpe_morpheus_network_floating_ip.example", "id"),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check:  resource.ComposeAggregateTestCheckFunc(checks...),
			},
		},
	})
}
