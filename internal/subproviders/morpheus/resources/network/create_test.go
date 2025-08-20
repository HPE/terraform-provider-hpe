// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

//go:build experimental

package network_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/internal/subproviders/morpheus/testhelpers"
)

func TestAccMorpheusNetworkResourceCreateRequiredAttrsOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	// Generate unique name for this test run
	uniqueName := acctest.RandomWithPrefix(t.Name())

	// Build the configuration with variables and defaults for required fields only
	providerConfig := testhelpers.ProviderBlock()
	configText := providerConfig + `
variable "name" {
  description = "Network name"
  type        = string
  default     = "terraform-network-minimal"
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
  default     = "example-resource-group"
}

variable "config_subnet_name" {
  description = "Subnet name for network config"
  type        = string
  default     = "example-subnet"
}

variable "config_subnet_cidr" {
  description = "Subnet CIDR for network config"
  type        = string
  default     = "10.0.1.0/24"
}

variable "cidr" {
  description = "CIDR Network"
  type        = string
  default     = "10.0.0.0/8"
}

resource "hpe_morpheus_network" "foo" {
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

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: configText,
				ConfigVariables: config.Variables{
					"name": config.StringVariable(uniqueName),
					// All other values use defaults
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.foo", "name", uniqueName),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.foo", "cloud_id", "4617"),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.foo", "group_id", "1"),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.foo", "type_id", "35"),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.foo", "cidr", "10.0.0.0/8"),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.foo", "config.resourceGroupId", "example-resource-group"),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.foo", "config.subnetName", "example-subnet"),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.foo", "config.subnetCidr", "10.0.1.0/24"),
					// Check that the resource was created with an ID
					resource.TestCheckResourceAttrSet(
						"hpe_morpheus_network.foo", "id"),
				),
			},
		},
	})
}

// TestAccMorpheusNetworkResourceCreateAllAttrsOk tests creating a network resource
// with all available fields populated and validates that each field is set correctly
func TestAccMorpheusNetworkResourceCreateAllAttrsOk(t *testing.T) {
	defer testhelpers.RecordResult(t)

	// Generate unique name for this test run
	uniqueName := acctest.RandomWithPrefix(t.Name())

	// Build the configuration with all available fields
	providerConfig := testhelpers.ProviderBlock()
	configText := providerConfig + `
variable "name" {
  description = "Network name"
  type        = string
  default     = "terraform-network-all-attrs"
}

variable "description" {
  description = "Network description"
  type        = string
  default     = "Network with all attributes set"
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
  default     = "10.100.0.0/16"
}

variable "visibility" {
  description = "Network visibility"
  type        = string
  default     = "public"
}

variable "active" {
  description = "Whether network is active"
  type        = bool
  default     = true
}

variable "dhcp_server" {
  description = "Whether DHCP server is enabled"
  type        = bool
  default     = true
}

variable "appliance_url_proxy_bypass" {
  description = "Whether to bypass proxy for appliance URL"
  type        = bool
  default     = false
}

variable "config_resource_group_id" {
  description = "Resource Group ID for network config"
  type        = string
  default     = "all-attrs-resource-group"
}

variable "config_subnet_name" {
  description = "Subnet name for network config"
  type        = string
  default     = "all-attrs-subnet"
}

variable "config_subnet_cidr" {
  description = "Subnet CIDR for network config"
  type        = string
  default     = "10.100.1.0/24"
}

variable "config_location" {
  description = "Location for network config"
  type        = string
  default     = "eastus"
}

variable "config_additional_field" {
  description = "Additional config field"
  type        = string
  default     = "test-value"
}

resource "hpe_morpheus_network" "all_attrs" {
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
    "resourceGroupId"    = var.config_resource_group_id
    "subnetName"         = var.config_subnet_name
    "subnetCidr"         = var.config_subnet_cidr
    "location"           = var.config_location
    "additionalField"    = var.config_additional_field
  }
  resource_permissions = {
    all = true
  }
  tenant_ids = [1, 2, 3]
}
`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: configText,
				ConfigVariables: config.Variables{
					"name": config.StringVariable(uniqueName),
					// All other values use defaults
				},
				Check: resource.ComposeTestCheckFunc(
					// Check basic required fields
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.all_attrs", "name", uniqueName),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.all_attrs", "description", "Network with all attributes set"),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.all_attrs", "cloud_id", "4617"),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.all_attrs", "pool_id", "1"),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.all_attrs", "group_id", "1"),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.all_attrs", "type_id", "35"),

					// Check network configuration fields
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.all_attrs", "cidr", "10.100.0.0/16"),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.all_attrs", "visibility", "public"),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.all_attrs", "active", "true"),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.all_attrs", "dhcp_server", "true"),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.all_attrs", "appliance_url_proxy_bypass", "false"),

					// Check config object fields
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.all_attrs", "config.resourceGroupId", "all-attrs-resource-group"),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.all_attrs", "config.subnetName", "all-attrs-subnet"),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.all_attrs", "config.subnetCidr", "10.100.1.0/24"),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.all_attrs", "config.location", "eastus"),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.all_attrs", "config.additionalField", "test-value"),

					// Check resource permissions
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.all_attrs", "resource_permissions.all", "true"),

					// Check tenant_ids
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.all_attrs", "tenant_ids.#", "3"),
					resource.TestCheckTypeSetElemAttr(
						"hpe_morpheus_network.all_attrs", "tenant_ids.*", "1"),
					resource.TestCheckTypeSetElemAttr(
						"hpe_morpheus_network.all_attrs", "tenant_ids.*", "2"),
					resource.TestCheckTypeSetElemAttr(
						"hpe_morpheus_network.all_attrs", "tenant_ids.*", "3"),

					// Check that the resource was created with an ID
					resource.TestCheckResourceAttrSet(
						"hpe_morpheus_network.all_attrs", "id"),
				),
			},
		},
	})
}

func TestAccMorpheusNetworkResourceCreateResourcePermissionsAllFalse(t *testing.T) {
	// TODO: Write test when PCCP-3372 is fixed
}

func TestAccMorpheusNetworkResourceCreateResourcePermissionsWithGroupIds(t *testing.T) {
	// TODO: Write test when PCCP-4209 is fixed
}

// TestAccMorpheusNetworkHostConfig tests creating a host network resource
// with host-specific configuration and empty config object
func TestAccMorpheusNetworkResourceCreateHostConfig(t *testing.T) {
	defer testhelpers.RecordResult(t)

	// Generate unique name for this test run
	uniqueName := acctest.RandomWithPrefix(t.Name())

	// Build the configuration with variables and defaults for host network
	providerConfig := testhelpers.ProviderBlock()
	configText := providerConfig + `
variable "name" {
  description = "Network name"
  type        = string
  default     = "terraform-host-network"
}

variable "description" {
  description = "Network description"
  type        = string
  default     = "A test host network"
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
  default     = "10.0.0.0/8"
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

resource "hpe_morpheus_network" "foo" {
  name        = var.name
  description = var.description
  cloud_id    = var.cloud_id
  pool_id     = var.pool_id
  group_id    = var.group_id
  type_id     = var.type_id
  config = {}
  active                       = var.active
  dhcp_server                  = var.dhcp_server
  appliance_url_proxy_bypass   = var.appliance_url_proxy_bypass
  resource_permissions = {
    all = true
  }
  tenant_ids  = [1]
  visibility  = var.visibility
  cidr        = var.cidr
}
`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: configText,
				ConfigVariables: config.Variables{
					"name": config.StringVariable(uniqueName),
					// All other values use defaults
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.foo", "name", uniqueName),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.foo", "description", "A test host network"),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.foo", "cloud_id", "17"),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.foo", "pool_id", "1"),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.foo", "group_id", "1"),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.foo", "type_id", "1"),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.foo", "active", "true"),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.foo", "dhcp_server", "false"),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.foo", "appliance_url_proxy_bypass", "true"),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.foo", "visibility", "private"),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.foo", "cidr", "10.0.0.0/8"),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.foo", "resource_permissions.all", "true"),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.foo", "tenant_ids.#", "1"),
					resource.TestCheckTypeSetElemAttr(
						"hpe_morpheus_network.foo", "tenant_ids.*", "1"),
					// Check that the resource was created with an ID
					resource.TestCheckResourceAttrSet(
						"hpe_morpheus_network.foo", "id"),
				),
			},
		},
	})
}

// TestAccMorpheusNetworkAws tests creating an AWS subnet network
// resource with specific configuration including assignPublicIp and
// availabilityZone settings
func TestAccMorpheusNetworkResourceCreateAws(t *testing.T) {
	defer testhelpers.RecordResult(t)

	// Generate unique name for this test run
	uniqueName := acctest.RandomWithPrefix(t.Name())

	// Build the configuration with AWS-specific settings
	providerConfig := testhelpers.ProviderBlock()
	configText := providerConfig + `
variable "name" {
  description = "Network name"
  type        = string
  default     = "terraform-aws-test"
}

variable "description" {
  description = "Network description"
  type        = string
  default     = "AWS subnet"
}

variable "cloud_id" {
  description = "Cloud (zone) id"
  type        = number
  default     = 207
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
  default     = 36
}

variable "cidr" {
  description = "CIDR Network"
  type        = string
  default     = "10.200.99.0/24"
}

variable "zone_pool_id" {
  description = "Zone pool id"
  type        = number
  default     = 12329
}

variable "config_assign_public_ip" {
  description = "Assign public IP setting for network config"
  type        = bool
  default     = true
}

variable "config_availability_zone" {
  description = "Availability zone setting for network config"
  type        = string
  default     = "us-west-1a"
}

variable "active" {
  description = "Whether network is active"
  type        = bool
  default     = true
}

variable "dhcp_server" {
  description = "Whether DHCP server is enabled"
  type        = bool
  default     = true
}

variable "appliance_url_proxy_bypass" {
  description = "Whether to bypass proxy for appliance URL"
  type        = bool
  default     = true
}

variable "visibility" {
  description = "Network visibility"
  type        = string
  default     = "private"
}

resource "hpe_morpheus_network" "aws" {
  name                         = var.name
  description                  = var.description
  cloud_id                     = var.cloud_id
  pool_id                      = var.pool_id
  group_id                     = var.group_id
  type_id                      = var.type_id
  config = {
    assignPublicIp   = var.config_assign_public_ip
    availabilityZone = var.config_availability_zone
  }
  active                       = var.active
  dhcp_server                  = var.dhcp_server
  appliance_url_proxy_bypass   = var.appliance_url_proxy_bypass
  resource_permissions = {
    all = true
  }
  tenant_ids                   = [1]
  visibility                   = var.visibility
  cidr                         = var.cidr
  zone_pool_id                 = var.zone_pool_id

  lifecycle {
    ignore_changes = [ name, display_name, description ]
  }
}
`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: configText,
				ConfigVariables: config.Variables{
					"name": config.StringVariable(uniqueName),
					// All other values use defaults
				},
				Check: resource.ComposeTestCheckFunc(
					// Check basic required fields
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.aws", "name",
						uniqueName),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.aws", "description",
						"AWS subnet"),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.aws", "cloud_id",
						"207"),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.aws", "pool_id", "1"),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.aws", "group_id", "1"),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.aws", "type_id", "36"),

					// Check network configuration fields
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.aws", "active", "true"),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.aws", "dhcp_server",
						"true"),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.aws",
						"appliance_url_proxy_bypass", "true"),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.aws", "visibility",
						"private"),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.aws", "cidr",
						"10.200.99.0/24"),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.aws", "zone_pool_id",
						"12329"),

					// Check config object fields specific to AWS
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.aws",
						"config.assignPublicIp", "true"),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.aws",
						"config.availabilityZone", "us-west-1a"),

					// Check resource permissions
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.aws",
						"resource_permissions.all", "true"),

					// Check tenant_ids
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.aws", "tenant_ids.#",
						"1"),
					resource.TestCheckTypeSetElemAttr(
						"hpe_morpheus_network.aws", "tenant_ids.*",
						"1"),

					// Check that the resource was created with an ID
					resource.TestCheckResourceAttrSet(
						"hpe_morpheus_network.aws", "id"),
				),
			},
		},
	})
}
