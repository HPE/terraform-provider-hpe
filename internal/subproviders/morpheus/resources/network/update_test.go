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

// TestAccMorpheusNetworkResourceUpdateOk tests updating a network resource
// with comprehensive validation of all updateable fields
func TestAccMorpheusNetworkResourceUpdateOk(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}
	defer testhelpers.RecordResult(t)

	// Generate unique name for this test run
	uniqueName := acctest.RandomWithPrefix(t.Name())

	providerConfig := testhelpers.ProviderBlock()

	// Base configuration with variables for all parameters
	baseConfigText := providerConfig + `
variable "name" {
  description = "Network name"
  type        = string
  default     = "TestAccMorpheusNetworkResourceUpdateOk"
}

variable "description" {
  description = "Network description"
  type        = string
  default     = "Initial network description"
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
  default     = "10.50.0.0/16"
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

variable "config_resource_group_id" {
  description = "Resource Group ID for network config"
  type        = string
  default     = "initial-resource-group"
}

variable "config_subnet_name" {
  description = "Subnet name for network config"
  type        = string
  default     = "initial-subnet"
}

variable "config_subnet_cidr" {
  description = "Subnet CIDR for network config"
  type        = string
  default     = "10.50.1.0/24"
}

variable "resource_permissions_all" {
  description = "Whether all groups have access"
  type        = bool
  default     = true
}

variable "resource_permissions_group_ids" {
  description = "Specific group IDs with access"
  type        = list(number)
  default     = null
}

variable "tenant_ids" {
  description = "List of tenant IDs"
  type        = list(number)
  default     = [1, 2]
}

resource "hpe_morpheus_network" "foo" {
	name                         = var.name
	description                  = var.description
	cloud_id                     = var.cloud_id
	pool_id                      = var.pool_id
	group_id                     = var.group_id
	type_id                      = var.type_id
	cidr                         = var.cidr
	visibility                   = var.visibility
	active                       = var.active
	dhcp_server                  = var.dhcp_server
	appliance_url_proxy_bypass   = var.appliance_url_proxy_bypass
	config = {
		"resourceGroupId" = var.config_resource_group_id
		"subnetName"      = var.config_subnet_name
		"subnetCidr"      = var.config_subnet_cidr
	}
	resource_permissions = {
		all = var.resource_permissions_all
		group_ids = var.resource_permissions_all ? null : var.resource_permissions_group_ids
	}
	tenant_ids = var.tenant_ids
}`

	// Common base configuration variables for initial state
	baseConfigVars := config.Variables{
		"name": config.StringVariable(uniqueName),
		// All other values use defaults from variable declarations
	}

	baseChecks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(
			"hpe_morpheus_network.foo",
			"name",
			uniqueName,
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_network.foo",
			"description",
			"Initial network description",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_network.foo",
			"cloud_id",
			"4617",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_network.foo",
			"pool_id",
			"1",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_network.foo",
			"group_id",
			"1",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_network.foo",
			"type_id",
			"35",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_network.foo",
			"cidr",
			"10.50.0.0/16",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_network.foo",
			"visibility",
			"private",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_network.foo",
			"active",
			"true",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_network.foo",
			"dhcp_server",
			"false",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_network.foo",
			"appliance_url_proxy_bypass",
			"true",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_network.foo",
			"resource_permissions.all",
			"true",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_network.foo",
			"tenant_ids.#",
			"2",
		),
		resource.TestCheckTypeSetElemAttr(
			"hpe_morpheus_network.foo",
			"tenant_ids.*",
			"1",
		),
		resource.TestCheckTypeSetElemAttr(
			"hpe_morpheus_network.foo",
			"tenant_ids.*",
			"2",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_network.foo",
			"config.resourceGroupId",
			"initial-resource-group",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_network.foo",
			"config.subnetName",
			"initial-subnet",
		),
		resource.TestCheckResourceAttr(
			"hpe_morpheus_network.foo",
			"config.subnetCidr",
			"10.50.1.0/24",
		),
	}

	checkFn := resource.ComposeAggregateTestCheckFunc(baseChecks...)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{{ // Step 1: Create resource
			Config:          baseConfigText,
			ConfigVariables: baseConfigVars,
			Check:           checkFn,
			PlanOnly:        false,
		}, { // Step 2: Plan only - verify no changes
			Config:             baseConfigText,
			ConfigVariables:    baseConfigVars,
			Check:              checkFn,
			ExpectNonEmptyPlan: false,
			PlanOnly:           true,
		}, { // Step 3: Plan only - test name change detection
			Config: baseConfigText,
			ConfigVariables: config.Variables{
				"name": config.StringVariable(uniqueName + "New"), // CHANGED
			},
			ExpectNonEmptyPlan: true,
			PlanOnly:           true,
		}, { // Step 4: Plan only - test description change detection
			Config: baseConfigText,
			ConfigVariables: config.Variables{
				"name":        config.StringVariable(uniqueName),
				"description": config.StringVariable("Changed network description"),
			},
			ExpectNonEmptyPlan: true,
			PlanOnly:           true,
		}, { // Step 5: Plan only - test pool_id change detection
			Config: baseConfigText,
			ConfigVariables: config.Variables{
				"name":    config.StringVariable(uniqueName),
				"pool_id": config.IntegerVariable(2), // CHANGED
			},
			ExpectNonEmptyPlan: true,
			PlanOnly:           true,
		}, { // Step 6: Plan only - test cidr change detection
			Config: baseConfigText,
			ConfigVariables: config.Variables{
				"name": config.StringVariable(uniqueName),
				"cidr": config.StringVariable("10.51.0.0/16"), // CHANGED
			},
			ExpectNonEmptyPlan: true,
			PlanOnly:           true,
		}, { // Step 7: Plan only - test visibility change detection
			Config: baseConfigText,
			ConfigVariables: config.Variables{
				"name":       config.StringVariable(uniqueName),
				"visibility": config.StringVariable("public"), // CHANGED
			},
			ExpectNonEmptyPlan: true,
			PlanOnly:           true,
		}, { // Step 8: Plan only - test active change detection
			Config: baseConfigText,
			ConfigVariables: config.Variables{
				"name":   config.StringVariable(uniqueName),
				"active": config.BoolVariable(false), // CHANGED
			},
			ExpectNonEmptyPlan: true,
			PlanOnly:           true,
		}, { // Step 9: Plan only - test dhcp_server change detection
			Config: baseConfigText,
			ConfigVariables: config.Variables{
				"name":        config.StringVariable(uniqueName),
				"dhcp_server": config.BoolVariable(true), // CHANGED
			},
			ExpectNonEmptyPlan: true,
			PlanOnly:           true,
		}, { // Step 10: Plan only - test appliance_url_proxy_bypass change detection
			Config: baseConfigText,
			ConfigVariables: config.Variables{
				"name":                       config.StringVariable(uniqueName),
				"appliance_url_proxy_bypass": config.BoolVariable(false), // CHANGED
			},
			ExpectNonEmptyPlan: true,
			PlanOnly:           true,
		}, { // Step 11: Plan only - test config change detection
			Config: baseConfigText,
			ConfigVariables: config.Variables{
				"name":                     config.StringVariable(uniqueName),
				"config_resource_group_id": config.StringVariable("changed-resource-group"),
			},
			ExpectNonEmptyPlan: true,
			PlanOnly:           true,
		}, { // Step 12: Plan only - test resource_permissions change detection
			Config: baseConfigText,
			ConfigVariables: config.Variables{
				"name":                     config.StringVariable(uniqueName),
				"resource_permissions_all": config.BoolVariable(false),
				"resource_permissions_group_ids": config.ListVariable(
					config.IntegerVariable(1), config.IntegerVariable(2),
				),
			},
			ExpectNonEmptyPlan: true,
			PlanOnly:           true,
		}, { // Step 13: Plan only - test tenant_ids change detection
			Config: baseConfigText,
			ConfigVariables: config.Variables{
				"name": config.StringVariable(uniqueName),
				"tenant_ids": config.ListVariable(
					config.IntegerVariable(1), config.IntegerVariable(3),
				),
			},
			ExpectNonEmptyPlan: true,
			PlanOnly:           true,
		}, { // Step 14: Apply comprehensive changes - modify many fields
			Config: baseConfigText,
			ConfigVariables: config.Variables{
				// Keep original name
				"name":        config.StringVariable(uniqueName),
				"description": config.StringVariable("Comprehensive update test"),
				"pool_id":     config.IntegerVariable(2),
				// Keep original CIDR
				"cidr":                       config.StringVariable("10.50.0.0/16"),
				"visibility":                 config.StringVariable("public"),
				"active":                     config.BoolVariable(false),
				"dhcp_server":                config.BoolVariable(true),
				"appliance_url_proxy_bypass": config.BoolVariable(false),
				"config_resource_group_id":   config.StringVariable("updated-resource-group"),
				"config_subnet_name":         config.StringVariable("updated-subnet"),
				"config_subnet_cidr":         config.StringVariable("10.99.1.0/24"),
				"resource_permissions_all":   config.BoolVariable(false),
				"resource_permissions_group_ids": config.ListVariable(
					config.IntegerVariable(1), config.IntegerVariable(3),
				),
				"tenant_ids": config.ListVariable(
					config.IntegerVariable(1), config.IntegerVariable(2), config.IntegerVariable(3),
				),
			},
			Check: resource.ComposeTestCheckFunc(
				// Verify all updated fields
				resource.TestCheckResourceAttr(
					"hpe_morpheus_network.foo", "name", uniqueName), // Name unchanged
				resource.TestCheckResourceAttr(
					"hpe_morpheus_network.foo", "description", "Comprehensive update test"),
				resource.TestCheckResourceAttr(
					"hpe_morpheus_network.foo", "pool_id", "2"),
				resource.TestCheckResourceAttr(
					"hpe_morpheus_network.foo", "cidr", "10.50.0.0/16"), // CIDR unchanged
				resource.TestCheckResourceAttr(
					"hpe_morpheus_network.foo", "visibility", "public"),
				resource.TestCheckResourceAttr(
					"hpe_morpheus_network.foo", "active", "false"),
				resource.TestCheckResourceAttr(
					"hpe_morpheus_network.foo", "dhcp_server", "true"),
				resource.TestCheckResourceAttr(
					"hpe_morpheus_network.foo", "appliance_url_proxy_bypass", "false"),
				resource.TestCheckResourceAttr(
					"hpe_morpheus_network.foo", "config.resourceGroupId", "updated-resource-group"),
				resource.TestCheckResourceAttr(
					"hpe_morpheus_network.foo", "config.subnetName", "updated-subnet"),
				resource.TestCheckResourceAttr(
					"hpe_morpheus_network.foo", "config.subnetCidr", "10.99.1.0/24"),
				resource.TestCheckResourceAttr(
					"hpe_morpheus_network.foo", "resource_permissions.all", "false"),
				resource.TestCheckResourceAttr(
					"hpe_morpheus_network.foo", "resource_permissions.group_ids.#", "2"),
				resource.TestCheckTypeSetElemAttr(
					"hpe_morpheus_network.foo", "resource_permissions.group_ids.*", "1"),
				resource.TestCheckTypeSetElemAttr(
					"hpe_morpheus_network.foo", "resource_permissions.group_ids.*", "3"),
				resource.TestCheckResourceAttr(
					"hpe_morpheus_network.foo", "tenant_ids.#", "3"),
				resource.TestCheckTypeSetElemAttr(
					"hpe_morpheus_network.foo", "tenant_ids.*", "1"),
				resource.TestCheckTypeSetElemAttr(
					"hpe_morpheus_network.foo", "tenant_ids.*", "2"),
				resource.TestCheckTypeSetElemAttr(
					"hpe_morpheus_network.foo", "tenant_ids.*", "3"),
				// Verify fields that shouldn't change
				resource.TestCheckResourceAttr(
					"hpe_morpheus_network.foo", "cloud_id", "4617"),
				resource.TestCheckResourceAttr(
					"hpe_morpheus_network.foo", "group_id", "1"),
				resource.TestCheckResourceAttr(
					"hpe_morpheus_network.foo", "type_id", "35"),
				resource.TestCheckResourceAttrSet(
					"hpe_morpheus_network.foo", "id"),
			),
			PlanOnly: false,
		},
		},
	})
}

// TestAccMorpheusNetworkResourceUpdateNameChange tests that changing the name
// attribute forces resource replacement due to the RequiresReplace plan modifier
func TestAccMorpheusNetworkResourceUpdateNameChange(t *testing.T) {
	defer testhelpers.RecordResult(t)

	// Generate unique names for this test run
	initialName := acctest.RandomWithPrefix(t.Name() + "-initial")
	updatedName := acctest.RandomWithPrefix(t.Name() + "-updated")

	// Build the configuration with variables
	providerConfig := testhelpers.ProviderBlock()
	configText := providerConfig + `
variable "name" {
  description = "Network name"
  type        = string
  default     = "TestAccMorpheusNetworkResourceUpdateNameChange"
}

variable "cloud_id" {
  description = "Cloud (zone) id"
  type        = number
  default     = 4617
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

variable "config_resource_group_id" {
  description = "Resource Group ID for network config"
  type        = string
  default     = "name-change-resource-group"
}

variable "config_subnet_name" {
  description = "Subnet name for network config"
  type        = string
  default     = "name-change-subnet"
}

variable "config_subnet_cidr" {
  description = "Subnet CIDR for network config"
  type        = string
  default     = "10.0.3.0/24"
}

variable "cidr" {
  description = "CIDR Network"
  type        = string
  default     = "10.0.0.0/8"
}

resource "hpe_morpheus_network" "name_change_test" {
  name     = var.name
  cloud_id = var.cloud_id
  group_id = var.group_id
  type_id  = var.type_id
  cidr     = var.cidr
  config = {
    "resourceGroupId" = var.config_resource_group_id
    "subnetName"      = var.config_subnet_name
    "subnetCidr"      = var.config_subnet_cidr
  }
}
`

	var initialResourceId string

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Step 1: Create resource with initial name
				Config: configText,
				ConfigVariables: config.Variables{
					"name": config.StringVariable(initialName),
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.name_change_test", "name", initialName),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.name_change_test", "cloud_id", "4617"),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.name_change_test", "group_id", "1"),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.name_change_test", "type_id", "35"),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.name_change_test", "cidr", "10.0.0.0/8"),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.name_change_test", "config.resourceGroupId", "name-change-resource-group"),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.name_change_test", "config.subnetName", "name-change-subnet"),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.name_change_test", "config.subnetCidr", "10.0.3.0/24"),
					// Check that the resource was created with an ID and store it
					resource.TestCheckResourceAttrSet(
						"hpe_morpheus_network.name_change_test", "id"),
					// Store the initial ID for comparison
					func(s *terraform.State) error {
						resourceName := "hpe_morpheus_network.name_change_test"
						rs, ok := s.RootModule().Resources[resourceName]
						if !ok {
							return fmt.Errorf("Not found: %s", resourceName)
						}
						initialResourceId = rs.Primary.ID
						return nil
					},
				),
			},
			{
				// Step 2: Change the name and verify the resource is replaced (new ID)
				Config: configText,
				ConfigVariables: config.Variables{
					"name": config.StringVariable(updatedName),
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.name_change_test", "name", updatedName),
					// Verify other attributes remain the same
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.name_change_test", "cloud_id", "4617"),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.name_change_test", "group_id", "1"),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.name_change_test", "type_id", "35"),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.name_change_test", "cidr", "10.0.0.0/8"),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.name_change_test", "config.resourceGroupId", "name-change-resource-group"),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.name_change_test", "config.subnetName", "name-change-subnet"),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.name_change_test", "config.subnetCidr", "10.0.3.0/24"),
					// Check that the resource has a new ID (replacement occurred)
					resource.TestCheckResourceAttrSet(
						"hpe_morpheus_network.name_change_test", "id"),
					// Verify the ID changed (resource was replaced)
					func(s *terraform.State) error {
						rs, ok := s.RootModule().Resources["hpe_morpheus_network.name_change_test"]
						if !ok {
							return fmt.Errorf("Not found: hpe_morpheus_network.name_change_test")
						}
						currentResourceId := rs.Primary.ID
						if currentResourceId == initialResourceId {
							return fmt.Errorf("Expected resource ID to change due to name change (RequiresReplace), "+
								"but ID remained the same: %s", currentResourceId)
						}
						return nil
					},
				),
			},
		},
	})
}

// TestAccMorpheusNetworkResourceUpdateCidrChange tests that changing the cidr
// or cidr_ipv6 attributes forces resource replacement due to RequiresReplace
func TestAccMorpheusNetworkResourceUpdateCidrChange(t *testing.T) {
	defer testhelpers.RecordResult(t)

	// Generate unique name for this test run
	uniqueName := acctest.RandomWithPrefix(t.Name())

	// Build the configuration with variables
	providerConfig := testhelpers.ProviderBlock()
	configText := providerConfig + `
variable "name" {
  description = "Network name"
  type        = string
  default     = "TestAccMorpheusNetworkResourceUpdateCidrChange"
}

variable "cloud_id" {
  description = "Cloud (zone) id"
  type        = number
  default     = 4617
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
  default     = "10.1.0.0/16"
}

variable "cidr_ipv6" {
  description = "IPv6 Network CIDR"
  type        = string
  default     = "2001:db8::/32"
}

variable "config_resource_group_id" {
  description = "Resource Group ID for network config"
  type        = string
  default     = "cidr-change-resource-group"
}

variable "config_subnet_name" {
  description = "Subnet name for network config"
  type        = string
  default     = "cidr-change-subnet"
}

variable "config_subnet_cidr" {
  description = "Subnet CIDR for network config"
  type        = string
  default     = "10.1.1.0/24"
}

resource "hpe_morpheus_network" "cidr_change_test" {
  name      = var.name
  cloud_id  = var.cloud_id
  group_id  = var.group_id
  type_id   = var.type_id
  cidr      = var.cidr
  cidr_ipv6 = var.cidr_ipv6
  config = {
    "resourceGroupId" = var.config_resource_group_id
    "subnetName"      = var.config_subnet_name
    "subnetCidr"      = var.config_subnet_cidr
  }
}
`

	var initialResourceId, afterCidrChangeId string

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Step 1: Create resource with initial CIDR values
				Config: configText,
				ConfigVariables: config.Variables{
					"name": config.StringVariable(uniqueName),
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.cidr_change_test", "name", uniqueName),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.cidr_change_test", "cloud_id", "4617"),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.cidr_change_test", "group_id", "1"),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.cidr_change_test", "type_id", "35"),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.cidr_change_test", "cidr", "10.1.0.0/16"),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.cidr_change_test", "cidr_ipv6", "2001:db8::/32"),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.cidr_change_test", "config.resourceGroupId", "cidr-change-resource-group"),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.cidr_change_test", "config.subnetName", "cidr-change-subnet"),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.cidr_change_test", "config.subnetCidr", "10.1.1.0/24"),
					// Check that the resource was created with an ID and store it
					resource.TestCheckResourceAttrSet(
						"hpe_morpheus_network.cidr_change_test", "id"),
					// Store the initial ID for comparison
					func(s *terraform.State) error {
						rs, ok := s.RootModule().Resources["hpe_morpheus_network.cidr_change_test"]
						if !ok {
							return fmt.Errorf("Not found: hpe_morpheus_network.cidr_change_test")
						}
						initialResourceId = rs.Primary.ID
						return nil
					},
				),
			},
			{
				// Step 2: Change the IPv4 CIDR and verify the resource is replaced (new ID)
				Config: configText,
				ConfigVariables: config.Variables{
					"name": config.StringVariable(uniqueName),
					"cidr": config.StringVariable("10.2.0.0/16"), // CHANGED
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.cidr_change_test", "name", uniqueName),
					// Verify the CIDR changed
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.cidr_change_test", "cidr", "10.2.0.0/16"),
					// Verify other attributes remain the same
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.cidr_change_test", "cloud_id", "4617"),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.cidr_change_test", "group_id", "1"),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.cidr_change_test", "type_id", "35"),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.cidr_change_test", "cidr_ipv6", "2001:db8::/32"), // unchanged
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.cidr_change_test", "config.resourceGroupId", "cidr-change-resource-group"),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.cidr_change_test", "config.subnetName", "cidr-change-subnet"),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.cidr_change_test", "config.subnetCidr", "10.1.1.0/24"),
					// Check that the resource has a new ID (replacement occurred)
					resource.TestCheckResourceAttrSet(
						"hpe_morpheus_network.cidr_change_test", "id"),
					// Verify the ID changed (resource was replaced)
					func(s *terraform.State) error {
						rs, ok := s.RootModule().Resources["hpe_morpheus_network.cidr_change_test"]
						if !ok {
							return fmt.Errorf("Not found: hpe_morpheus_network.cidr_change_test")
						}
						currentResourceId := rs.Primary.ID
						if currentResourceId == initialResourceId {
							return fmt.Errorf("Expected resource ID to change due to CIDR change (RequiresReplace), "+
								"but ID remained the same: %s", currentResourceId)
						}
						afterCidrChangeId = currentResourceId
						return nil
					},
				),
			},
			{
				// Step 3: Change the IPv6 CIDR and verify the resource is replaced again (new ID)
				Config: configText,
				ConfigVariables: config.Variables{
					"name":      config.StringVariable(uniqueName),
					"cidr":      config.StringVariable("10.2.0.0/16"),   // keep previous change
					"cidr_ipv6": config.StringVariable("2001:db9::/32"), // CHANGED
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.cidr_change_test", "name", uniqueName),
					// Verify the IPv6 CIDR changed
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.cidr_change_test", "cidr_ipv6", "2001:db9::/32"),
					// Verify other attributes remain the same
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.cidr_change_test", "cloud_id", "4617"),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.cidr_change_test", "group_id", "1"),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.cidr_change_test", "type_id", "35"),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.cidr_change_test", "cidr", "10.2.0.0/16"), // unchanged from step 2
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.cidr_change_test", "config.resourceGroupId", "cidr-change-resource-group"),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.cidr_change_test", "config.subnetName", "cidr-change-subnet"),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.cidr_change_test", "config.subnetCidr", "10.1.1.0/24"),
					// Check that the resource has a new ID (replacement occurred)
					resource.TestCheckResourceAttrSet(
						"hpe_morpheus_network.cidr_change_test", "id"),
					// Verify the ID changed again (resource was replaced)
					func(s *terraform.State) error {
						rs, ok := s.RootModule().Resources["hpe_morpheus_network.cidr_change_test"]
						if !ok {
							return fmt.Errorf("Not found: hpe_morpheus_network.cidr_change_test")
						}
						currentResourceId := rs.Primary.ID
						if currentResourceId == afterCidrChangeId {
							return fmt.Errorf("Expected resource ID to change due to IPv6 CIDR change (RequiresReplace), "+
								"but ID remained the same: %s", currentResourceId)
						}
						if currentResourceId == initialResourceId {
							return fmt.Errorf("Resource ID reverted to initial ID, "+
								"this should not happen: %s", currentResourceId)
						}
						return nil
					},
				),
			},
		},
	})
}
