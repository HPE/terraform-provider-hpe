// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkrouterbgpneighbor_test

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

// TestAccMorpheusNetworkRouterBgpNeighborResourceWriteOnlyDriftDetected: the
// BGP neighbor password is rotated out of band (a direct API PUT, as the
// Morpheus UI would do). The next plan sees the appliance's stored-secret hash
// no longer matches the baseline recorded at the last send and plans an update
// that re-sends the configured value; the plan is then clean. Gated like the
// other NSX-T BGP neighbor tests.
func TestAccMorpheusNetworkRouterBgpNeighborResourceWriteOnlyDriftDetected(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.NetworkRouter, capabilities.NSXT)
	testhelpers.SkipUnlessAcc(t)

	name := acctest.RandomWithPrefix(t.Name())
	ipAddress := "192.168.2." + acctest.RandStringFromCharSet(2, "123456789")
	addr := "hpe_morpheus_network_router_bgp_neighbor.wo"

	// The tier-0 comes from the shared NSX-T fixture so the test follows the
	// appliance like the other NSX-T tests; the BGP source address stays an
	// override because no fixture carries it.
	fixture := testhelpers.RequireNsxtFixture(t)
	routerID, sourceAddress := fixture.Tier0RouterID, bgpNeighborSourceAddress

	if v := os.Getenv("TF_ACC_BGP_SOURCE_ADDRESS"); v != "" {
		sourceAddress = v
	}

	config := testhelpers.ProviderBlock() + `
resource "hpe_morpheus_network_router_bgp_neighbor" "wo" {
  router_id   = ` + routerID + `
  ip_address  = "` + ipAddress + `"
  description = "` + name + `"
  remote_as   = "65001"
  weight      = 60
  keep_alive  = 60
  hold_down   = 180
  hop_limit   = ` + bgpNeighborHopLimit + `
  password_wo = "Secret123!"

  config_nsxt = {
    source_addresses = ["` + sourceAddress + `"]
  }
}
`

	// Captured from state by the create step, used by the next step's
	// PreConfig to rotate the secret out of band.
	var id, router int64

	capture := func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[addr]
		if !ok {
			return fmt.Errorf("%s not in state", addr)
		}

		var err error
		if router, err = strconv.ParseInt(rs.Primary.Attributes["router_id"], 10, 64); err != nil {
			return err
		}

		id, err = strconv.ParseInt(rs.Primary.ID, 10, 64)

		return err
	}

	// hashOf reads the appliance's stored-secret hash for the BGP neighbor.
	hashOf := func(ctx context.Context, client *sdk.APIClient) string {
		resp, httpResp, err := client.NetworksAPI.
			GetNetworkRouterBgpNeighbor(ctx, id, router).Execute()
		if err != nil {
			t.Fatalf("read BGP neighbor back: %s (status %v)", err, httpResp)
		}
		if resp.NetworkRouterBgpNeighbor == nil {
			t.Fatalf("read BGP neighbor back: nil NetworkRouterBgpNeighbor")
		}

		if v := resp.NetworkRouterBgpNeighbor.PasswordHash.Get(); v != nil {
			return *v
		}

		return ""
	}

	rotate := func() {
		if id == 0 {
			t.Fatal("capture did not run: BGP neighbor id is 0")
		}

		ctx := context.Background()

		client, err := testhelpers.NewClientForServer(ctx, "")
		if err != nil {
			t.Fatalf("client: %s", err)
		}

		before := hashOf(ctx, client)

		rotated := "Rotated456!"
		body := sdk.UpdateNetworkRouterBgpNeighborRequest{
			NetworkRouterBgpNeighbor: &sdk.UpdateNetworkRouterBgpNeighborRequestNetworkRouterBgpNeighbor{
				Password: &rotated,
			},
		}

		if _, httpResp, err := client.NetworksAPI.
			UpdateNetworkRouterBgpNeighbor(ctx, id, router).
			UpdateNetworkRouterBgpNeighborRequest(body).Execute(); err != nil {
			t.Fatalf("out-of-band rotate failed: %s (status %v)", err, httpResp)
		}

		// Assert the appliance actually re-stored the secret: if the stored
		// hash is unchanged there is no out-of-band drift to detect and the
		// test would pass even with drift detection removed -- a false pass.
		if after := hashOf(ctx, client); after == before {
			t.Fatalf("appliance did not re-store the BGP neighbor password: hash unchanged after rotate (%q); nothing to detect", after)
		}
	}

	update := resource.ConfigPlanChecks{PreApply: []plancheck.PlanCheck{
		plancheck.ExpectResourceAction(addr, plancheck.ResourceActionUpdate),
	}}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{Config: config, Check: capture},
			{
				PreConfig:        rotate,
				Config:           config,
				ConfigPlanChecks: update,
			},
			{Config: config, PlanOnly: true},
		},
	})
}
