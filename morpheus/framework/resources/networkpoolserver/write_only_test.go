// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkpoolserver_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

// TestAccMorpheusNetworkPoolServerResourceWriteOnlyOk: changes to
// service_password_wo are detected by digest (MORPH-18005); the deprecated
// marker still works; both import forms leave a clean plan. Needs a reachable
// Infoblox (the API validates the service URL on create).
func TestAccMorpheusNetworkPoolServerResourceWriteOnlyOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.NetworkPool)
	testhelpers.SkipUnlessAcc(t)

	// Morpheus 9.1.2 validates the IPAM credentials on create and update (a
	// reachable host with a wrong password is refused: "Error connecting to
	// infoblox"), so neither a change-detection matrix that sends three
	// different passwords nor an out-of-band rotation to a wrong value can be
	// staged against an ordinary live integration. The password handling is
	// covered by the unit tests and by the same wodigest mechanism exercised
	// live on resources whose secret the appliance stores without validating
	// (hpe_morpheus_image). Set TF_ACC_IPAM_ROTATABLE=1 to re-enable this test
	// against an IPAM whose credentials can be rotated. Skip first, before the
	// fixture helper, so the body is not executed.
	if os.Getenv("TF_ACC_IPAM_ROTATABLE") == "" {
		t.Skip("set TF_ACC_IPAM_ROTATABLE=1 against an IPAM whose credentials can be rotated; MORPH-19005 follow-up")
	}

	// Morpheus 9.1.2 verifies the service URL on create, so a reachable IPAM
	// is required; without the fixture there is no appliance to create against.
	f := testhelpers.RequireIpamFixture(t)

	name := acctest.RandomWithPrefix(t.Name())

	c := testhelpers.WriteOnlyCase{
		Address:   "hpe_morpheus_network_pool_server.wo",
		Attribute: "service_password_wo",
		Marker:    "service_password_wo_version",
		APIPath:   "/api/networks/pool-servers",
		APIKey:    "networkPoolServer.servicePassword",
		Config: func(value, extra string, unrelated bool) string {
			filter := "192.168.0.0/16"
			if unrelated {
				filter = "192.168.1.0/24"
			}

			return testhelpers.ProviderBlock() + f.WithAppID(`
resource "hpe_morpheus_network_pool_server" "wo" {
  name                = "`+name+`"
  type_code           = "`+f.TypeCode+`"
  enabled             = false
  service_url         = "`+f.URL+`"
  service_username    = "`+f.Username+`"
  service_password_wo = "`+value+`"`+extra+`
  ignore_ssl          = true
  network_filter      = "`+filter+`"
}
`)
		},
	}

	rec := testhelpers.RecordRequests(t)

	steps := testhelpers.WriteOnlyMatrixSteps(c, rec, "Secret123!", "Secret456!", "Secret789!")
	steps = append(steps, testhelpers.WriteOnlyImportSteps(c, "Secret789!", "service_username")...)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps:                    steps,
	})
}
