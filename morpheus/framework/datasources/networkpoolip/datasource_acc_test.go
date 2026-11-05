// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkpoolip_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
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

// TestAccMorpheusNetworkPoolIPLookups creates a reservation (hostname only) in
// the pool identified by TF_ACC_NETWORK_POOL_ID, then reads it back through the
// data source by id, by ip_address and by hostname in a single config to prove
// all three lookup paths resolve to the same reservation.
func TestAccMorpheusNetworkPoolIPLookups(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.NetworkPool)

	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	t.Parallel()

	poolID := os.Getenv("TF_ACC_NETWORK_POOL_ID")
	if poolID == "" {
		t.Skip("TF_ACC_NETWORK_POOL_ID not set; skipping network pool IP acc test")
	}

	hostname := acctest.RandomWithPrefix(t.Name())

	config := testhelpers.ProviderBlock() + fmt.Sprintf(`
resource "hpe_morpheus_network_pool_ip" "test" {
  network_pool_id = %[1]s
  hostname        = %[2]q
}

data "hpe_morpheus_network_pool_ip" "by_id" {
  network_pool_id = %[1]s
  id              = hpe_morpheus_network_pool_ip.test.id
}

data "hpe_morpheus_network_pool_ip" "by_ip" {
  network_pool_id = %[1]s
  ip_address      = hpe_morpheus_network_pool_ip.test.ip_address
}

# The hostname is a literal, so take it from the resource to make Terraform
# wait for the reservation to exist before reading it back.
data "hpe_morpheus_network_pool_ip" "by_hostname" {
  network_pool_id = %[1]s
  hostname        = hpe_morpheus_network_pool_ip.test.hostname
}
`, poolID, hostname)

	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttrPair(
			"data.hpe_morpheus_network_pool_ip.by_id", "id",
			"hpe_morpheus_network_pool_ip.test", "id",
		),
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_network_pool_ip.by_id", "hostname", hostname,
		),
		resource.TestCheckResourceAttrPair(
			"data.hpe_morpheus_network_pool_ip.by_ip", "id",
			"hpe_morpheus_network_pool_ip.test", "id",
		),
		resource.TestCheckResourceAttrPair(
			"data.hpe_morpheus_network_pool_ip.by_hostname", "id",
			"hpe_morpheus_network_pool_ip.test", "id",
		),
		resource.TestCheckResourceAttr(
			"data.hpe_morpheus_network_pool_ip.by_hostname", "hostname", hostname,
		),
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
