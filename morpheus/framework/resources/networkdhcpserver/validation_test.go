// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkdhcpserver_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

// TestAccMorpheusNetworkDhcpServerResourceValidationZeroLeaseTime verifies that
// a zero lease_time is rejected at plan time. Sending 0 to the API is silently
// coerced to the 86400 default, so we reject it up front and instruct callers
// to omit the field instead (MORPH-14701).
func TestAccMorpheusNetworkDhcpServerResourceValidationZeroLeaseTime(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.NSXT)

	providerConfig := testhelpers.ProviderBlock()

	config := providerConfig + `
resource "hpe_morpheus_network_dhcp_server" "test" {
	name                   = "TestAccMorpheusNetworkDhcpServerValidationZeroLeaseTime"
	network_integration_id = 1
	server_ip_address      = "192.168.1.1/24"
	lease_time             = 0
}
`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testhelpers.GetAccTestFactories(t, adapter.NewMorpheus(), nil),
		Steps: []resource.TestStep{
			{
				Config:      config,
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`must be at least 1`),
			},
		},
	})
}
