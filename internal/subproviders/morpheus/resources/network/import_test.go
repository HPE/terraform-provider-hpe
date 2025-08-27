// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

//go:build experimental

package network_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/testhelpers"
)

func TestAccMorpheusNetworkImport(t *testing.T) {
	defer testhelpers.RecordResult(t)
	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	// Generate unique name for this test run
	uniqueName := acctest.RandomWithPrefix(t.Name())

	providerConfig := testhelpers.ProviderBlock()

	// nolint: gosec
	resourceCfg := `
variable "name" {
  description = "Network name"
  type        = string
  default     = "TestAccMorpheusNetworkImport"
}

variable "description" {
  description = "Network description"
  type        = string
  default     = "mclaren 1 updated"
}

variable "cloud_id" {
  description = "Cloud (zone) id"
  type        = number
  default     = 4617
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
  default     = 35
}

variable "cidr" {
  description = "CIDR Network"
  type        = string
  default     = "10.0.0.0/8"
}

variable "visibility" {
  description = "Network visibility"
  type        = string
  default     = "private"
}

variable "config_resource_group_id" {
  description = "Resource Group ID for network config"
  type        = string
  default     = "morph-qa"
}

variable "config_subnet_name" {
  description = "Subnet name for network config"
  type        = string
  default     = "mclaren-subnet-3"
}

variable "config_subnet_cidr" {
  description = "Subnet CIDR for network config"
  type        = string
  default     = "10.0.1.0/24"
}

resource "hpe_morpheus_network" "net1" {
	name = var.name
	description = var.description
	cloud_id = var.cloud_id
	pool_id = var.pool_id
	group_id = var.group_id
	type_id = var.type_id
	config = {
		"resourceGroupId" = var.config_resource_group_id
		"subnetName" = var.config_subnet_name
		"subnetCidr" = var.config_subnet_cidr
	}
	active = true
	dhcp_server = false
	search_domains = null
	appliance_url_proxy_bypass = true
	no_proxy = null
	resource_permissions = {
		all = true
	}
	tenant_ids = [1,2]
	visibility = var.visibility
	cidr = var.cidr
}
`

	resourceCfgRemove := `
# This allows us to create a resource using the provider
# and then import it in a separate resource.Test.
#
# A regular 'Config:' test step (with import block) can be used for the
# import test (rather than an 'ImportState:' style test step)
#
# The state is preserved after import and available to
# subsequent tests.
#
# The 'removed' block means the resource is removed from state
# (and terraform control) but not deleted.
#
# This avoids both triggering two deletes for the same resource,
# and the dreaded "resource is already under terraform control"
# error when running the follow on import test.
removed {
from = hpe_morpheus_network.net1

lifecycle {
destroy = false
}
}
`

	baseChecks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_network.net1",
			"name",
			uniqueName,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_network.net1",
			"type_id",
			"35",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_network.net1",
			"group_id",
			"1",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_network.net1",
			"cidr",
			"10.0.0.0/8",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_network.net1",
			"description",
			"mclaren 1 updated",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_network.net1",
			"cloud_id",
			"4617",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_network.net1",
			"pool_id",
			"1",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_network.net1",
			"active",
			"true",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_network.net1",
			"dhcp_server",
			"false",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_network.net1",
			"appliance_url_proxy_bypass",
			"true",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_network.net1",
			"tenant_ids.#",
			"2",
		),
		resource.TestCheckTypeSetElemAttr(
			"hpe_morpheus_network.net1",
			"tenant_ids.*",
			"1",
		),
		resource.TestCheckTypeSetElemAttr(
			"hpe_morpheus_network.net1",
			"tenant_ids.*",
			"2",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_network.net1",
			"visibility",
			"private",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_network.net1",
			"resource_permissions.all",
			"true",
		),
		// Note: Removed config checks as it's not returned by API and causes state drift
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(baseChecks...)

	var cachedID string

	// This is a new TestCase - we know for sure
	// we inherit no state from the TestCase above
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + resourceCfg,
				ConfigVariables: config.Variables{
					"name": config.StringVariable(uniqueName),
					// All other values use defaults
				},
				Check: resource.ComposeTestCheckFunc(
					func(s *terraform.State) error {
						// Cache ID for use later
						rs := s.RootModule().Resources["hpe_morpheus_network.net1"]
						if rs == nil {
							return fmt.Errorf("resource not found")
						}
						cachedID = rs.Primary.ID

						return nil
					},
					checkFn,
				),
				PlanOnly: false,
			},
			{
				// remove resource from terraform state (without deleting it)
				Config:   providerConfig + resourceCfgRemove,
				PlanOnly: false,
			},
		},
	})

	importCfg := resourceCfg + `
import {
  to = hpe_morpheus_network.net1
  id = ` + cachedID + `
}
`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + importCfg,
				ConfigVariables: config.Variables{
					"name": config.StringVariable(uniqueName),
					// All other values use defaults
				},
				PlanOnly: false,
				Check: resource.ComposeTestCheckFunc(
					checkFn,
				),
			},
			{
				// check that a plan after import detects no changes
				Config: providerConfig + resourceCfg,
				ConfigVariables: config.Variables{
					"name": config.StringVariable(uniqueName),
					// All other values use defaults
				},
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}
