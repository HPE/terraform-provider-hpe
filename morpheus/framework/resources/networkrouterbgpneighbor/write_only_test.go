// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkrouterbgpneighbor_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

// TestAccMorpheusNetworkRouterBgpNeighborResourceWriteOnlyOk: changes to
// password_wo are detected by digest (MORPH-18005); the deprecated marker still
// works; both import forms -- on this resource's composite "router.id" id --
// leave a clean plan. Needs the NSX-T tier-0 router the other tests use.
func TestAccMorpheusNetworkRouterBgpNeighborResourceWriteOnlyOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.NetworkRouter, capabilities.NSXT)
	testhelpers.SkipUnlessAcc(t)

	name := acctest.RandomWithPrefix(t.Name())
	ipAddress := "192.168.2." + acctest.RandStringFromCharSet(2, "123456789")
	addr := "hpe_morpheus_network_router_bgp_neighbor.wo"

	// The tier-0 router and a source address on one of its interfaces; the
	// package constants are the defaults, overridable when that router is
	// gone from the appliance.
	routerID, sourceAddress := existingTier0RouterID, bgpNeighborSourceAddress
	if v := os.Getenv("TF_ACC_BGP_TIER0_ROUTER_ID"); v != "" {
		routerID = v
	}

	if v := os.Getenv("TF_ACC_BGP_SOURCE_ADDRESS"); v != "" {
		sourceAddress = v
	}

	c := testhelpers.WriteOnlyCase{
		Address:   addr,
		Attribute: "password_wo",
		Marker:    "password_wo_version",
		APIPath:   "/bgp-neighbors",
		APIKey:    "networkRouterBgpNeighbor.password",
		Config: func(value, extra string, unrelated bool) string {
			weight := "60"
			if unrelated {
				weight = "70"
			}

			return testhelpers.ProviderBlock() + `
resource "hpe_morpheus_network_router_bgp_neighbor" "wo" {
  router_id   = ` + routerID + `
  ip_address  = "` + ipAddress + `"
  description = "` + name + `"
  remote_as   = "65001"
  weight      = ` + weight + `
  keep_alive  = 60
  hold_down   = 180
  hop_limit   = ` + bgpNeighborHopLimit + `
  password_wo = "` + value + `"` + extra + `

  config_nsxt = {
    source_addresses = ["` + sourceAddress + `"]
  }
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

	steps := testhelpers.WriteOnlyMatrixSteps(c, rec, "Secret123!", "Secret456!", "Secret789!")
	// description is not returned by the API (as in the existing import test).
	steps = append(steps, testhelpers.WriteOnlyImportSteps(c, "Secret789!", "description")...)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps:                    steps,
	})
}
