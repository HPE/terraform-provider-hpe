// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkserver_test

import (
	"os"
	"testing"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

// TestAccMorpheusNetworkServerResourceWriteOnly: changes to
// service_password_wo are detected by digest; the deprecated marker still
// works; both import forms leave a clean plan.
//
// The write-only matrix creates an NSX-T network server, and Morpheus validates
// the NSX manager URL and credentials on create, so a placeholder URL is
// refused. The test therefore needs a reachable NSX manager whose credentials
// can be rotated three times, which no shared lab offers; it is gated on
// TF_ACC_NSXT_MANAGER_ROTATABLE naming that manager's URL and skips otherwise.
// The password handling itself is covered by the unit tests and by the same
// wodigest mechanism exercised live on other resources.
func TestAccMorpheusNetworkServerResourceWriteOnly(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.NSXT)
	testhelpers.SkipUnlessAcc(t)

	managerURL := os.Getenv("TF_ACC_NSXT_MANAGER_ROTATABLE")
	if managerURL == "" {
		t.Skip("TF_ACC_NSXT_MANAGER_ROTATABLE not set: creating an NSX-T network server needs a reachable " +
			"NSX manager with rotatable credentials (Morpheus validates the URL and credentials on create)")
	}

	cloudID := os.Getenv("TF_ACC_NSXT_CLOUD_ID")
	if cloudID == "" {
		t.Skip("TF_ACC_NSXT_CLOUD_ID not set: the cloud the NSX-T network server is attached to")
	}

	name := acctest.RandomWithPrefix(t.Name())

	c := testhelpers.WriteOnlyCase{
		Address:   "hpe_morpheus_network_server.wo",
		Attribute: "service_password_wo",
		Marker:    "service_password_wo_version",
		APIPath:   "/api/networks/servers",
		APIKey:    "networkServer.servicePassword",
		Config: func(value, extra string, unrelated bool) string {
			filter := "10.0.0.0/8"
			if unrelated {
				filter = "10.1.0.0/16"
			}

			return testhelpers.ProviderBlock() + `
resource "hpe_morpheus_network_server" "wo" {
  name                = "` + name + `"
  type_code           = "nsx-t"
  cloud_id            = ` + cloudID + `
  service_url         = "` + managerURL + `"
  service_username    = "admin"
  service_password_wo = "` + value + `"` + extra + `
  network_filter      = "` + filter + `"
}
`
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
