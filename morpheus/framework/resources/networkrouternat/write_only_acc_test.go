// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkrouternat_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

// TestAccMorpheusNetworkRouterNatResourceWriteOnlyOk: changes to the write-only
// service (one of action/firewall/service, which NSX-T updates in place) are
// detected by digest (MORPH-18005): an in-place update that re-sends the rule's
// configuration; the deprecated markers still work; both import forms -- on
// the composite "router.nat" id -- leave a clean plan. Needs the NSX-T fixture
// the other NAT tests use (a tier-1 router is created for it).
func TestAccMorpheusNetworkRouterNatResourceWriteOnlyOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.NetworkRouter)
	testhelpers.SkipUnlessAcc(t)

	fixture := testhelpers.RequireNsxtFixture(t)
	name := acctest.RandomWithPrefix("tfwo-nat")
	addr := "hpe_morpheus_network_router_nat.wo"

	router := `
data "hpe_morpheus_network_router" "nat_tier0" {
  id = ` + fixture.Tier0RouterID + `
}

resource "hpe_morpheus_network_router" "nat_tier1" {
  name                   = "` + name + `-tier1"
  group_id               = ` + fixture.GroupID + `
  network_integration_id = ` + fixture.NetworkServerID + `

  config_nsxt_gateway_tier1 = {
    ip_management_type = "dhcpLocal"
    edge_cluster       = "` + fixture.EdgeCluster + `"
    fail_over          = "NON_PREEMPTIVE"
    tier0_gateway      = data.hpe_morpheus_network_router.nat_tier0.provider_id
  }
}
`

	c := testhelpers.WriteOnlyCase{
		Address:   addr,
		Attribute: "service",
		Marker:    "service_version",
		APIPath:   "/nats",
		APIKey:    "networkRouterNAT.config.service",
		Config: func(value, extra string, unrelated bool) string {
			desc := "write-only matrix"
			if unrelated {
				desc = "write-only matrix, changed"
			}

			return testhelpers.ProviderBlock() + router + `
resource "hpe_morpheus_network_router_nat" "wo" {
  router_id          = hpe_morpheus_network_router.nat_tier1.id
  name               = "` + name + `"
  description        = "` + desc + `"
  action             = "SNAT"
  source_network     = "10.20.0.0/24"
  translated_network = "192.0.2.10"
  service            = "` + value + `"` + extra + `
}
`
		},
		ImportIDFunc: func(s *terraform.State) (string, error) {
			rs, ok := s.RootModule().Resources[addr]
			if !ok {
				return "", fmt.Errorf("%s not in state", addr)
			}

			return rs.Primary.Attributes["router_id"] + "." + rs.Primary.ID, nil
		},
	}

	rec := testhelpers.RecordRequests(t)

	steps := testhelpers.WriteOnlyMatrixSteps(c, rec, "HTTP", "HTTPS", "DNS")
	steps = append(steps, testhelpers.WriteOnlyImportSteps(c, "DNS", "action", "firewall", "protocol")...)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps:                    steps,
	})
}
