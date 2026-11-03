// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package network_test

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers"
	"github.com/HPE/terraform-provider-hpe/morpheus/testhelpers/capabilities"
	"github.com/HPE/terraform-provider-hpe/provider"
	"github.com/HPE/terraform-provider-hpe/provider/adapter"
)

func TestMain(m *testing.M) {
	code := m.Run()
	testhelpers.WriteMergedResults()
	os.Exit(code)
}

func newProviderWithError() (tfprotov6.ProviderServer, error) {
	providerInstance := provider.New("test", adapter.NewMorpheus())()

	return providerserver.NewProtocol6WithError(providerInstance)()
}

var testAccProtoV6ProviderFactories = map[string]func() (
	tfprotov6.ProviderServer, error,
){
	"hpe": newProviderWithError,
}

func TestNetworkDataSourceExample(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.NetworkServer)

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
  default     = 6446
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
  tenant_ids                 = [data.hpe_morpheus_whoami.current.tenant_id]
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
				Config: providerConfig + testhelpers.WhoamiBlock() + networkResourceConfig,
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

// A network two clouds see has a record in each; a lookup by name alone is
// ambiguous, and cloud_id resolves it. The fixture is a standing pair on the
// test appliance: the same Azure VNet seen by the resource-group-scoped cloud
// and by the subscription-wide cloud.
func TestAccNetworkDataSourceCloudIDSelectsAmongTwins(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.Azure)

	providerConfig := testhelpers.ProviderBlock()

	const (
		name        = "morph-qa-net"
		scopedCloud = "4617"
		wideCloud   = "4618"
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// By name alone: both records are found, and the error says which
				// cloud each belongs to.
				Config: providerConfig + `
data "hpe_morpheus_network" "ambiguous" {
  name = "` + name + `"
}
`,
				// Terraform re-wraps diagnostics, so every gap may be a newline.
				ExpectError: regexp.MustCompile(`multiple\s+networks\s+found\s+with\s+name\s+` + name +
					`[\s\S]*\(cloud\s+(` + scopedCloud + `|` + wideCloud + `)\)[\s\S]*\(cloud\s+(` + scopedCloud + `|` + wideCloud + `)\)`),
			},
			{
				// With cloud_id: the scoped cloud's record.
				Config: providerConfig + `
data "hpe_morpheus_network" "scoped" {
  name     = "` + name + `"
  cloud_id = ` + scopedCloud + `
}

data "hpe_morpheus_network" "wide" {
  name     = "` + name + `"
  cloud_id = ` + wideCloud + `
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.hpe_morpheus_network.scoped", "name", name),
					resource.TestCheckResourceAttr("data.hpe_morpheus_network.scoped", "cloud_id", scopedCloud),
					resource.TestCheckResourceAttr("data.hpe_morpheus_network.wide", "cloud_id", wideCloud),
					// Two records of one network.
					func(s *terraform.State) error {
						a := s.RootModule().Resources["data.hpe_morpheus_network.scoped"].Primary.ID
						b := s.RootModule().Resources["data.hpe_morpheus_network.wide"].Primary.ID
						if a == b {
							return fmt.Errorf("expected different records in the two clouds, both are %s", a)
						}

						return nil
					},
				),
			},
			{
				// By id with the wrong cloud: refused, naming the real one.
				Config: providerConfig + `
data "hpe_morpheus_network" "scoped" {
  name     = "` + name + `"
  cloud_id = ` + scopedCloud + `
}

data "hpe_morpheus_network" "wrong" {
  id       = data.hpe_morpheus_network.scoped.id
  cloud_id = ` + wideCloud + `
}
`,
				ExpectError: regexp.MustCompile(`belongs\s+to\s+cloud\s+` + scopedCloud + `,\s+not\s+cloud\s+` + wideCloud),
			},
		},
	})
}

// Each HVM cluster has networks of its own, so two clusters in one cloud can
// both have a network of one name -- the Management network Morpheus creates
// with every cluster is the standing case -- and name plus cloud_id cannot
// tell them apart. cluster_id limits the lookup to one cluster's networks.
// The fixture is the HVM cluster the instance tests use, resolved by name.
func TestAccNetworkDataSourceClusterIDScopesLookup(t *testing.T) {
	defer testhelpers.RecordResult(t)

	capabilities.MustHaveOrSkip(t, capabilities.HVM)

	clusterName := testhelpers.HVMClusterName(t)
	networkName := testhelpers.HVMNetworkName(t)

	providerConfig := testhelpers.ProviderBlock()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// By name within the cluster: found, and the filter is carried
				// through as configured.
				Config: providerConfig + `
data "hpe_morpheus_cluster" "hvm" {
  name = "` + clusterName + `"
}

data "hpe_morpheus_network" "scoped" {
  name       = "` + networkName + `"
  cluster_id = data.hpe_morpheus_cluster.hvm.id
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.hpe_morpheus_network.scoped", "name", networkName),
					resource.TestCheckResourceAttrPair(
						"data.hpe_morpheus_network.scoped", "cluster_id",
						"data.hpe_morpheus_cluster.hvm", "id",
					),
					resource.TestCheckResourceAttrPair(
						"data.hpe_morpheus_network.scoped", "cloud_id",
						"data.hpe_morpheus_cluster.hvm", "cloud.id",
					),
				),
			},
			{
				// By id with the cluster: the same network, accepted.
				Config: providerConfig + `
data "hpe_morpheus_cluster" "hvm" {
  name = "` + clusterName + `"
}

data "hpe_morpheus_network" "scoped" {
  name       = "` + networkName + `"
  cluster_id = data.hpe_morpheus_cluster.hvm.id
}

data "hpe_morpheus_network" "by_id" {
  id         = data.hpe_morpheus_network.scoped.id
  cluster_id = data.hpe_morpheus_cluster.hvm.id
}
`,
				Check: resource.TestCheckResourceAttrPair(
					"data.hpe_morpheus_network.by_id", "id",
					"data.hpe_morpheus_network.scoped", "id",
				),
			},
			{
				// A cluster that does not exist has no networks: not found, and
				// the error names the cluster.
				Config: providerConfig + `
data "hpe_morpheus_network" "elsewhere" {
  name       = "` + networkName + `"
  cluster_id = 2147483647
}
`,
				ExpectError: regexp.MustCompile(`not\s+found\s+in\s+cluster\s+2147483647`),
			},
			{
				// By id against a cluster that does not have it: refused, and
				// the error names both. The membership check lists the
				// cluster's networks of this network's name, so this is the
				// listing coming back empty.
				Config: providerConfig + `
data "hpe_morpheus_cluster" "hvm" {
  name = "` + clusterName + `"
}

data "hpe_morpheus_network" "scoped" {
  name       = "` + networkName + `"
  cluster_id = data.hpe_morpheus_cluster.hvm.id
}

data "hpe_morpheus_network" "misplaced" {
  id         = data.hpe_morpheus_network.scoped.id
  cluster_id = 2147483647
}
`,
				ExpectError: regexp.MustCompile(`is\s+not\s+one\s+of\s+cluster\s+2147483647's\s+networks`),
			},
		},
	})
}
