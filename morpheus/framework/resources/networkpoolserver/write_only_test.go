// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkpoolserver_test

import (
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

			return testhelpers.ProviderBlock() + `
resource "hpe_morpheus_network_pool_server" "wo" {
  name                = "` + name + `"
  type_code           = "infoblox"
  enabled             = false
  service_url         = "https://infoblox.example.com/wapi/v2.12"
  service_username    = "admin"
  service_password_wo = "` + value + `"` + extra + `
  ignore_ssl          = true
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
