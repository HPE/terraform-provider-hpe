// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkserver_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/HPE/terraform-provider-hpe/morpheus/framework/resources/networkserver"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

// envSdnControllerAddress names a spare CN2 controller address the SDN network
// server create test can point a new network server at. Creating a second
// network server against the SAME controller as the fixture one can be rejected
// or trigger cluster on-boarding, so the create test only runs when a spare
// controller is supplied; otherwise the SDN coverage is an import-and-read of
// the fixture network server.
const envSdnControllerAddress = "TF_ACC_SDN_CONTROLLER_ADDRESS"

func TestMain(m *testing.M) {
	code := testhelpers.TestMain(m)
	testhelpers.WriteMergedResults()
	os.Exit(code)
}

// TestAccMorpheusNetworkServerResourceJuniperSdnImportOk imports the fixture
// SDN network server and verifies its type and controller address. It makes no
// changes to the appliance, so it is safe to run against the shared fixture.
func TestAccMorpheusNetworkServerResourceJuniperSdnImportOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.JuniperSDN)
	f := testhelpers.RequireSdnFixture(t)

	t.Parallel()

	resourceName := "hpe_morpheus_network_server.juniper_sdn"

	// A minimal declaration whose only purpose is to receive the imported
	// state; the controller address is immutable and read back from the API.
	config := testhelpers.ProviderBlock() + `
resource "hpe_morpheus_network_server" "juniper_sdn" {
  name      = "imported-sdn"
  type_code = "morpheus-junipersdn-network"
  cloud_id  = ` + f.CloudID + `

  config_juniper_sdn = {
    controller_address = "placeholder"
  }
}
`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				// Import only: the fixture network server is read, not created.
				Config:        config,
				ResourceName:  resourceName,
				ImportState:   true,
				ImportStateId: f.NetworkServerID,
				ImportStateCheck: func(states []*terraform.InstanceState) error {
					if len(states) == 0 {
						return fmt.Errorf("imported network server %s produced no state", f.NetworkServerID)
					}

					is := states[0]

					if got := is.Attributes["type_code"]; got != "morpheus-junipersdn-network" {
						return fmt.Errorf(
							"type_code = %q after import, want morpheus-junipersdn-network", got)
					}

					if is.Attributes["config_juniper_sdn.controller_address"] == "" {
						return fmt.Errorf(
							"config_juniper_sdn.controller_address is empty after import, want it set")
					}

					return nil
				},
			},
		},
	})
}

// TestAccMorpheusNetworkServerResourceJuniperSdnCreateOk creates an SDN network
// server against a spare controller. It only runs when TF_ACC_SDN_CONTROLLER_ADDRESS
// names a controller other than the fixture one, since pointing a second
// network server at the fixture controller may be rejected or trigger cluster
// on-boarding.
func TestAccMorpheusNetworkServerResourceJuniperSdnCreateOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.JuniperSDN)
	f := testhelpers.RequireSdnFixture(t)

	controller := os.Getenv(envSdnControllerAddress)
	if controller == "" {
		t.Skip("spare SDN controller not configured; set: " + envSdnControllerAddress)
	}

	t.Parallel()

	name := acctest.RandomWithPrefix(t.Name())
	resourceName := "hpe_morpheus_network_server.juniper_sdn"

	config, err := networkserver.RenderNetworkServerJuniperSdnConfig(t, map[string]string{
		"Name":              name,
		"CloudId":           f.CloudID,
		"ControllerAddress": controller,
	})
	if err != nil {
		t.Fatal(err)
	}

	config = testhelpers.ProviderBlock() + config

	checks := resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr(resourceName, "name", name),
		resource.TestCheckResourceAttr(resourceName, "type_code", "morpheus-junipersdn-network"),
		resource.TestCheckResourceAttr(resourceName, "config_juniper_sdn.controller_address", controller),
		resource.TestCheckResourceAttrSet(resourceName, "id"),
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
				ImportStateVerify: true,
				// write_only_revision is the provider's own marker: create and
				// import each assign a fresh Unix-second timestamp, so the two
				// differ whenever the steps straddle a second boundary.
				ImportStateVerifyIgnore: []string{
					"service_password_wo", "service_password_wo_version", "write_only_revision",
				},
			},
		},
	})
}
