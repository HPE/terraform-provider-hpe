// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package network_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/internal/provider"
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus"
	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/testhelpers"
)

func TestMain(m *testing.M) {
	code := m.Run()
	testhelpers.WriteMergedResults()
	os.Exit(code)
}

func newProviderWithError() (tfprotov6.ProviderServer, error) {
	providerInstance := provider.New("test", morpheus.New())()

	return providerserver.NewProtocol6WithError(providerInstance)()
}

var testAccProtoV6ProviderFactories = map[string]func() (
	tfprotov6.ProviderServer, error,
){
	"hpe": newProviderWithError,
}

func TestNetworkDataSourceExample(t *testing.T) {
	defer testhelpers.RecordResult(t)
	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	// Use the standard provider config from testhelpers (no gock mocking)
	providerConfig := testhelpers.ProviderBlock()

	// Generate unique name for this test run
	uniqueName := acctest.RandomWithPrefix(t.Name())

	// Network resource configuration (using host network type for simplicity)
	networkResourceConfig := `
variable "name" {
  description = "Network name"
  type        = string
  default     = "TestNetworkDataSourceExample"
}

variable "description" {
  description = "Network description"
  type        = string
  default     = "A test network for datasource acceptance testing"
}

variable "display_name" {
  description = "Network display name"
  type        = string
  default     = "Test Network Display Name"
}

variable "cloud_id" {
  description = "Cloud (zone) id"
  type        = number
  default     = 17
}

variable "pool_id" {
  description = "Network pool id"
  type        = number
  default     = 1
}

variable "group_id" {
  description = "Group (site) id"
  type        = number
  default     = 1
}

variable "type_id" {
  description = "Network type id"
  type        = number
  default     = 1
}

variable "cidr" {
  description = "CIDR Network"
  type        = string
  default     = "10.0.0.0/24"
}

variable "visibility" {
  description = "Network visibility"
  type        = string
  default     = "private"
}

variable "active" {
  description = "Whether network is active"
  type        = bool
  default     = true
}

variable "dhcp_server" {
  description = "Whether DHCP server is enabled"
  type        = bool
  default     = false
}

variable "appliance_url_proxy_bypass" {
  description = "Whether to bypass proxy for appliance URL"
  type        = bool
  default     = true
}

resource "hpe_morpheus_network" "test" {
  name                       = var.name
  description                = var.description
  display_name               = var.display_name
  cloud_id                   = var.cloud_id
  pool_id                    = var.pool_id
  group_id                   = var.group_id
  type_id                    = var.type_id
  config                     = {}
  active                     = var.active
  dhcp_server                = var.dhcp_server
  appliance_url_proxy_bypass = var.appliance_url_proxy_bypass
  tenant_ids                 = [1]
  visibility                 = var.visibility
  cidr                       = var.cidr
}

data "hpe_morpheus_network" "example" {
  name = hpe_morpheus_network.test.name
}
`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + networkResourceConfig,
				ConfigVariables: config.Variables{
					"name": config.StringVariable(uniqueName),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					// Check that datasource reads the created network correctly
					resource.TestCheckResourceAttr(
						"data.hpe_morpheus_network.example",
						"name",
						uniqueName,
					),
					resource.TestCheckResourceAttrPair(
						"data.hpe_morpheus_network.example",
						"id",
						"hpe_morpheus_network.test",
						"id",
					),
					resource.TestCheckResourceAttr(
						"data.hpe_morpheus_network.example",
						"description",
						"A test network for datasource acceptance testing",
					),
					resource.TestCheckResourceAttr(
						"data.hpe_morpheus_network.example",
						"display_name",
						"Test Network Display Name",
					),
					resource.TestCheckResourceAttr(
						"data.hpe_morpheus_network.example",
						"cidr",
						"10.0.0.0/24",
					),
					resource.TestCheckResourceAttr(
						"data.hpe_morpheus_network.example",
						"visibility",
						"private",
					),
					resource.TestCheckResourceAttr(
						"data.hpe_morpheus_network.example",
						"active",
						"true",
					),
					// Note: cloud_id, group_id, and type_id are not available in the datasource schema
					// Only checking the attributes that are actually exposed by the datasource
				),
			},
		},
	})
}
