// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package network_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
)

// Uses Azure
func TestAccMorpheusNetworkResourceCreateRequiredAttrsOk(t *testing.T) {
	if capabilities.Missing(t, capabilities.Network) {
		t.Log("Skipping test due to missing capabilities")

		return
	}
	defer testhelpers.RecordResult(t)
	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	// nolint: goconst
	providerConfig := testhelpers.ProviderBlock()

	// Generate unique name for this test run
	uniqueName := acctest.RandomWithPrefix(t.Name())

	// Build the configuration with variables and defaults for required fields only
	configText := providerConfig + `
variable "name" {
  description = "Network name"
  type        = string
  default     = "TestAccMorpheusNetworkResourceCreateRequiredAttrsOk"
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
  default     = "morph-qa"
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
  active   = true
  pool_id     = 6446
  tenant_ids = [1,2]
  name     = var.name
  cloud_id = var.cloud_id
  group_id = var.group_id
  type_id  = var.type_id
  cidr     = var.cidr
  labels   = ["terraform", "acctest", "hpe_morpheus_network", "sweepable"]
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
						"hpe_morpheus_network.foo", "name", uniqueName,
					),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.foo", "cloud_id", "4617",
					),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.foo", "group_id", "1",
					),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.foo", "type_id", "35",
					),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.foo", "cidr", "10.0.0.0/8",
					),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.foo", "config.resourceGroupId", "morph-qa",
					),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.foo", "config.subnetName", "example-subnet",
					),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.foo", "config.subnetCidr", "10.0.1.0/24",
					),
					// Check labels
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.foo", "labels.#", "4",
					),
					resource.TestCheckTypeSetElemAttr(
						"hpe_morpheus_network.foo", "labels.*", "terraform",
					),
					resource.TestCheckTypeSetElemAttr(
						"hpe_morpheus_network.foo", "labels.*", "acctest",
					),
					resource.TestCheckTypeSetElemAttr(
						"hpe_morpheus_network.foo", "labels.*", "hpe_morpheus_network",
					),
					resource.TestCheckTypeSetElemAttr(
						"hpe_morpheus_network.foo", "labels.*", "sweepable",
					),
					// Check resource permissions (computed-only)
					resource.TestCheckResourceAttrSet(
						"hpe_morpheus_network.foo", "resource_permissions.all",
					),
					// Check that the resource was created with an ID
					resource.TestCheckResourceAttrSet(
						"hpe_morpheus_network.foo", "id",
					),
				),
			},
		},
	})
}

// TestAccMorpheusNetworkResourceCreateAllAttrsOk tests creating a network resource
// with all available fields populated and validates that each field is set correctly
// Uses Azure
func TestAccMorpheusNetworkResourceCreateAllAttrsOk(t *testing.T) {
	if capabilities.Missing(t, capabilities.Azure, capabilities.NetworkDHCP) {
		t.Log("Skipping test due to missing capabilities")

		return
	}
	defer testhelpers.RecordResult(t)
	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	// nolint: goconst
	providerConfig := testhelpers.ProviderBlock()

	// Generate unique name for this test run
	uniqueName := acctest.RandomWithPrefix(t.Name())

	// Path to example configuration files
	examplePath := "../../../../examples/resources/hpe_morpheus_network/azure"

	// Read the resource.tf file from disk
	resourceContent, err := os.ReadFile(filepath.Join(examplePath, "resource.tf"))
	if err != nil {
		t.Fatalf("Failed to read resource.tf: %v", err)
	}

	// Combine provider config and resource file content
	configText := providerConfig + "\n" + string(resourceContent)

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
						"hpe_morpheus_network.all_attrs", "name", uniqueName,
					),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.all_attrs", "description", "Network with all attributes set",
					),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.all_attrs", "cloud_id", "4617",
					),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.all_attrs", "pool_id", "6446",
					),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.all_attrs", "group_id", "1",
					),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.all_attrs", "type_id", "35",
					),

					// Check network configuration fields
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.all_attrs", "cidr", "10.100.0.0/16",
					),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.all_attrs", "visibility", "public",
					),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.all_attrs", "active", "true",
					),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.all_attrs", "dhcp_server", "true",
					),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.all_attrs", "appliance_url_proxy_bypass", "false",
					),

					// Check config object fields
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.all_attrs", "config.resourceGroupId", "all-attrs-resource-group",
					),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.all_attrs", "config.subnetName", "all-attrs-subnet",
					),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.all_attrs", "config.subnetCidr", "10.100.1.0/24",
					),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.all_attrs", "config.location", "eastus",
					),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.all_attrs", "config.additionalField", "test-value",
					),

					// Check resource permissions (computed-only)
					resource.TestCheckResourceAttrSet(
						"hpe_morpheus_network.all_attrs", "resource_permissions.all",
					),

					// Check tenant_ids
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.all_attrs", "tenant_ids.#", "3",
					),
					resource.TestCheckTypeSetElemAttr(
						"hpe_morpheus_network.all_attrs", "tenant_ids.*", "1",
					),
					resource.TestCheckTypeSetElemAttr(
						"hpe_morpheus_network.all_attrs", "tenant_ids.*", "2",
					),
					resource.TestCheckTypeSetElemAttr(
						"hpe_morpheus_network.all_attrs", "tenant_ids.*", "3",
					),

					// Check that the resource was created with an ID
					resource.TestCheckResourceAttrSet(
						"hpe_morpheus_network.all_attrs", "id",
					),
				),
			},
		},
	})
}

// func TestAccMorpheusNetworkResourceCreateResourcePermissionsAllFalse(_ *testing.T) {
// 	// TODO: Write test when PCCP-3372 is fixed
// }

// func TestAccMorpheusNetworkResourceCreateResourcePermissionsWithGroupIds(_ *testing.T) {
// 	// TODO: Write test when PCCP-4209 is fixed
// }

// TestAccMorpheusNetworkHostConfig tests creating a host network resource
// with host-specific configuration and empty config object
func TestAccMorpheusNetworkResourceCreateHostConfig(t *testing.T) {
	if capabilities.Missing(t, capabilities.NetworkDHCP) {
		t.Log("Skipping test due to missing capabilities")

		return
	}
	defer testhelpers.RecordResult(t)
	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	// nolint: goconst
	providerConfig := testhelpers.ProviderBlock()

	// Generate unique name for this test run
	uniqueName := acctest.RandomWithPrefix(t.Name())

	// Path to example configuration files
	examplePath := "../../../../examples/resources/hpe_morpheus_network/host"

	// Read the resource.tf file from disk
	resourceContent, err := os.ReadFile(filepath.Join(examplePath, "resource.tf"))
	if err != nil {
		t.Fatalf("Failed to read resource.tf: %v", err)
	}

	// Combine provider config and resource file content
	configText := providerConfig + "\n" + string(resourceContent)

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
						"hpe_morpheus_network.host", "name", uniqueName,
					),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.host", "description", "A test host network",
					),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.host", "cloud_id", "17",
					),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.host", "pool_id", "6446",
					),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.host", "group_id", "1",
					),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.host", "type_id", "1",
					),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.host", "active", "true",
					),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.host", "dhcp_server", "false",
					),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.host", "appliance_url_proxy_bypass", "true",
					),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.host", "visibility", "private",
					),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.host", "cidr", "10.0.0.0/8",
					),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.host", "tenant_ids.#", "1",
					),
					resource.TestCheckTypeSetElemAttr(
						"hpe_morpheus_network.host", "tenant_ids.*", "1",
					),
					// Check resource permissions (computed-only)
					resource.TestCheckResourceAttrSet(
						"hpe_morpheus_network.host", "resource_permissions.all",
					),
					// Check that the resource was created with an ID
					resource.TestCheckResourceAttrSet(
						"hpe_morpheus_network.host", "id",
					),
				),
			},
		},
	})
}

// TestAccMorpheusNetworkResourceCreateAWSExample tests creating an AWS subnet network
// resource with specific configuration including assignPublicIp and
// availabilityZone settings using example files
func TestAccMorpheusNetworkResourceCreateAWSExample(t *testing.T) {
	if capabilities.Missing(t, capabilities.AWS, capabilities.NetworkDHCP) {
		t.Log("Skipping test due to missing capabilities")

		return
	}
	defer testhelpers.RecordResult(t)
	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	// nolint: goconst
	providerConfig := testhelpers.ProviderBlock()

	// Generate unique name for this test run
	uniqueName := acctest.RandomWithPrefix(t.Name())

	// Path to example configuration files
	examplePath := "../../../../examples/resources/hpe_morpheus_network"

	// Read the resource.tf file from disk
	resourceContent, err := os.ReadFile(filepath.Join(examplePath, "resource.tf"))
	if err != nil {
		t.Fatalf("Failed to read resource.tf: %v", err)
	}

	// Combine provider config and resource file content
	configText := providerConfig + "\n" + string(resourceContent)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: configText,
				ConfigVariables: config.Variables{
					"name": config.StringVariable(uniqueName),
					// All other values use defaults from resource.tf
				},
				Check: resource.ComposeTestCheckFunc(
					// Check basic required fields
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.aws", "name",
						uniqueName,
					),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.aws", "description",
						"AWS subnet",
					),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.aws", "cloud_id",
						"207",
					),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.aws", "pool_id", "6446",
					),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.aws", "group_id", "1",
					),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.aws", "type_id", "36",
					),

					// Check network configuration fields
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.aws", "active", "true",
					),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.aws", "dhcp_server",
						"true",
					),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.aws",
						"appliance_url_proxy_bypass", "true",
					),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.aws", "visibility",
						"private",
					),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.aws", "cidr",
						"10.200.99.0/24",
					),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.aws", "zone_pool_id",
						"12329",
					),

					// Check config object fields specific to AWS
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.aws",
						"config.assignPublicIp", "true",
					),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.aws",
						"config.availabilityZone", "us-west-1a",
					),

					// Check resource permissions (computed-only)
					resource.TestCheckResourceAttrSet(
						"hpe_morpheus_network.aws", "resource_permissions.all",
					),

					// Check tenant_ids
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.aws", "tenant_ids.#",
						"1",
					),
					resource.TestCheckTypeSetElemAttr(
						"hpe_morpheus_network.aws", "tenant_ids.*",
						"1",
					),

					// Check that the resource was created with an ID
					resource.TestCheckResourceAttrSet(
						"hpe_morpheus_network.aws", "id",
					),
				),
			},
		},
	})
}

// TestAccMorpheusNetworkResourceCreateGcp tests creating a GCP network
// resource with specific configuration including mtu and autoCreate settings
func TestAccMorpheusNetworkResourceCreateGcp(t *testing.T) {
	if capabilities.Missing(t, capabilities.GCP, capabilities.NetworkDHCP) {
		t.Log("Skipping test due to missing capabilities")

		return
	}
	defer testhelpers.RecordResult(t)
	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	// nolint: goconst
	providerConfig := testhelpers.ProviderBlock()

	// Generate unique name for this test run
	uniqueName := acctest.RandomWithPrefix(t.Name())

	// Path to example configuration files
	examplePath := "../../../../examples/resources/hpe_morpheus_network/gcp"

	// Read the resource.tf file from disk
	resourceContent, err := os.ReadFile(filepath.Join(examplePath, "resource.tf"))
	if err != nil {
		t.Fatalf("Failed to read resource.tf: %v", err)
	}

	// Combine provider config and resource file content
	configText := providerConfig + "\n" + string(resourceContent)

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
						"hpe_morpheus_network.gcp", "name", uniqueName,
					),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.gcp", "description", "GCP network",
					),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.gcp", "cloud_id", "6",
					),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.gcp", "pool_id", "6446",
					),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.gcp", "group_id", "8",
					),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.gcp", "type_id", "38",
					),

					// Check network configuration fields
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.gcp", "active", "true",
					),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.gcp", "dhcp_server", "false",
					),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.gcp", "appliance_url_proxy_bypass", "true",
					),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.gcp", "visibility", "private",
					),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.gcp", "cidr", "10.0.0.0/8",
					),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.gcp", "zone_pool_id", "85990",
					),

					// Check config object fields specific to GCP
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.gcp", "config.mtu", "1460",
					),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.gcp", "config.autoCreate", "true",
					),

					// Check resource permissions (computed-only)
					resource.TestCheckResourceAttrSet(
						"hpe_morpheus_network.gcp", "resource_permissions.all",
					),

					// Check tenant_ids
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.gcp", "tenant_ids.#", "1",
					),
					resource.TestCheckTypeSetElemAttr(
						"hpe_morpheus_network.gcp", "tenant_ids.*", "1",
					),

					// Check that the resource was created with an ID
					resource.TestCheckResourceAttrSet(
						"hpe_morpheus_network.gcp", "id",
					),
				),
			},
		},
	})
}

// TestAccMorpheusNetworkResourceCreateOVSPortGroup tests creating an OVS Port Group network
// for cloud ID 7714.
func TestAccMorpheusNetworkResourceCreateOVSPortGroup(t *testing.T) {
	if capabilities.Missing(t, capabilities.HVM) {
		t.Log("Skipping test due to missing capabilities")

		return
	}
	defer testhelpers.RecordResult(t)
	if testing.Short() {
		t.Skip("Skipping slow test in short mode")
	}

	// nolint: goconst
	providerConfig := testhelpers.ProviderBlock()

	// Generate unique name for this test run
	uniqueName := acctest.RandomWithPrefix(t.Name())

	// Path to example configuration files
	examplePath := "../../../../examples/resources/hpe_morpheus_network/ovs_port_group"

	// Read the resource.tf file from disk
	resourceContent, err := os.ReadFile(filepath.Join(examplePath, "resource.tf"))
	if err != nil {
		t.Fatalf("Failed to read resource.tf: %v", err)
	}

	// Combine provider config and resource file content
	configText := providerConfig + "\n" + string(resourceContent)

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
					// Standard checks
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.ovs_port_group", "name", uniqueName,
					),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.ovs_port_group", "description", "OVS Port Group network",
					),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.ovs_port_group", "cloud_id", "7714",
					),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.ovs_port_group", "pool_id", "3251",
					),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.ovs_port_group", "group_id", "1",
					),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.ovs_port_group", "type_id", "63",
					),
					// Check OVS-specific fields
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.ovs_port_group", "switch_id", "Compute",
					),
					// Additional standard checks
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.ovs_port_group", "active", "true",
					),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.ovs_port_group", "dhcp_server", "false",
					),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.ovs_port_group", "appliance_url_proxy_bypass", "true",
					),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.ovs_port_group", "visibility", "public",
					),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.ovs_port_group", "cidr", "10.32.148.0/22",
					),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.ovs_port_group", "zone_pool_id", "62299",
					),
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.ovs_port_group", "vlan_id", "43",
					),
					// Check permissions
					resource.TestCheckResourceAttrSet(
						"hpe_morpheus_network.ovs_port_group", "resource_permissions.all",
					),
					// Check tenant IDs
					resource.TestCheckResourceAttr(
						"hpe_morpheus_network.ovs_port_group", "tenant_ids.#", "1",
					),
					resource.TestCheckTypeSetElemAttr(
						"hpe_morpheus_network.ovs_port_group", "tenant_ids.*", "1",
					),
					// ID check
					resource.TestCheckResourceAttrSet(
						"hpe_morpheus_network.ovs_port_group", "id",
					),
				),
			},
		},
	})
}
