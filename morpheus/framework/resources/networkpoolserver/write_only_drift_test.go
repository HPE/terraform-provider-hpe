// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkpoolserver_test

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

// TestAccMorpheusNetworkPoolServerResourceWriteOnlyDriftDetected: the service
// password is rotated out of band (a direct API PUT, as the Morpheus UI would
// do). The next plan sees the appliance's stored-secret hash no longer matches
// the baseline recorded at the last send and plans an update that re-sends the
// configured value; the plan is then clean. Needs a reachable Infoblox (the API
// validates the service URL on create).
func TestAccMorpheusNetworkPoolServerResourceWriteOnlyDriftDetected(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.NetworkPool)
	testhelpers.SkipUnlessAcc(t)

	// Morpheus 9.1.2 validates the IPAM credentials on create and update (a
	// reachable host with a wrong password is refused: "Error connecting to
	// infoblox"), so an out-of-band rotation cannot be staged against an
	// ordinary live integration. Set TF_ACC_IPAM_ROTATABLE=1 to re-enable this
	// test against an IPAM whose credentials can be rotated. Skip first, before
	// the fixture helper, so the body is not executed.
	if os.Getenv("TF_ACC_IPAM_ROTATABLE") == "" {
		t.Skip("set TF_ACC_IPAM_ROTATABLE=1 against an IPAM whose credentials can be rotated; MORPH-19005 follow-up")
	}

	// Morpheus 9.1.2 verifies the service URL on create, so a reachable IPAM
	// is required; without the fixture there is no appliance to create against.
	f := testhelpers.RequireIpamFixture(t)

	name := acctest.RandomWithPrefix(t.Name())
	addr := "hpe_morpheus_network_pool_server.wo"

	config := testhelpers.ProviderBlock() + f.WithAppID(`
resource "hpe_morpheus_network_pool_server" "wo" {
  name                = "`+name+`"
  type_code           = "`+f.TypeCode+`"
  enabled             = false
  service_url         = "`+f.URL+`"
  service_username    = "`+f.Username+`"
  service_password_wo = "`+f.Password+`"
  ignore_ssl          = true
  network_filter      = "192.168.0.0/16"
}
`)

	var id int64

	capture := func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[addr]
		if !ok {
			return fmt.Errorf("%s not in state", addr)
		}

		var err error
		id, err = strconv.ParseInt(rs.Primary.ID, 10, 64)

		return err
	}

	rotate := func() {
		if id == 0 {
			t.Fatal("capture did not run: network pool server id is 0")
		}

		ctx := context.Background()

		client, err := testhelpers.NewClientForServer(ctx, "")
		if err != nil {
			t.Fatalf("client: %s", err)
		}

		// The update request is a oneOf union typed per pool server type;
		// pick the fixture's variant so the rotation is a well-formed update
		// for that type, as the UI would send.
		rotated := "Rotated456!"
		body := sdk.UpdateNetworkPoolServerRequest{
			NetworkPoolServer: rotateVariant(t, f.TypeCode, rotated),
		}

		// A rotation that is accepted but not stored (the hash is unchanged)
		// would surface later as a confusing plan-check failure; compare the
		// hash around the PUT so the failure names the cause.
		before := poolServerPasswordHash(t, client, id)

		// With TF_ACC_IPAM_ROTATABLE set the target IPAM is expected to accept
		// the rotation; a rejected PUT is now a real failure (the test can no
		// longer silently pass by skipping here).
		if _, httpResp, err := client.NetworksAPI.
			UpdateNetworkPoolServer(ctx, id).
			UpdateNetworkPoolServerRequest(body).Execute(); err != nil {
			t.Fatalf("out-of-band rotate rejected by appliance: %s (status %v)", err, httpResp)
		}

		if after := poolServerPasswordHash(t, client, id); after == before {
			t.Fatalf("rotate PUT did not change servicePasswordHash (%q); drift cannot be detected", before)
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

// rotateVariant builds the typed update variant for the pool server type with
// only the service password set.
func rotateVariant(t *testing.T, typeCode, password string) *sdk.UpdateNetworkPoolServerRequestNetworkPoolServer {
	t.Helper()

	pw := *sdk.NewNullableString(&password)

	switch typeCode {
	case "phpipam":
		return &sdk.UpdateNetworkPoolServerRequestNetworkPoolServer{
			PhpIPAMNetworkPoolServerUpdate: &sdk.PhpIPAMNetworkPoolServerUpdate{ServicePassword: pw},
		}
	case "bluecat":
		return &sdk.UpdateNetworkPoolServerRequestNetworkPoolServer{
			BluecatNetworkPoolServerUpdate: &sdk.BluecatNetworkPoolServerUpdate{ServicePassword: pw},
		}
	case "solarwindsipam":
		return &sdk.UpdateNetworkPoolServerRequestNetworkPoolServer{
			SolarWindsNetworkPoolServerUpdate: &sdk.SolarWindsNetworkPoolServerUpdate{ServicePassword: pw},
		}
	case "infoblox":
		return &sdk.UpdateNetworkPoolServerRequestNetworkPoolServer{
			InfobloxNetworkPoolServerUpdate: &sdk.InfobloxNetworkPoolServerUpdate{ServicePassword: pw},
		}
	default:
		t.Skipf("no typed update variant for pool server type %q; rotation not staged", typeCode)

		return nil
	}
}

// poolServerPasswordHash reads the appliance's servicePasswordHash for the
// pool server ("" when not rendered).
func poolServerPasswordHash(t *testing.T, client *sdk.APIClient, id int64) string {
	t.Helper()

	resp, _, err := client.NetworksAPI.GetNetworkPoolServer(context.Background(), id).Execute()
	if err != nil || resp == nil || resp.NetworkPoolServer == nil || resp.NetworkPoolServer.ServicePasswordHash == nil {
		return ""
	}

	return *resp.NetworkPoolServer.ServicePasswordHash
}
