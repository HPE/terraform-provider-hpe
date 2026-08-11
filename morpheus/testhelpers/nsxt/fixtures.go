// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

// Package nsxt provides shared Terraform config fixtures for NSX-T acceptance
// tests.
//
// NSX-T tests need a realized tier-0 gateway to hang tier-1 gateways, BGP
// neighbors and NAT rules off. That tier-0 cannot be created by the provider
// (Morpheus only populates a gateway's locale-services after syncing a realized
// gateway, which a per-test create races), so it must be pre-provisioned on the
// target appliance.
//
// Previously each test package hardcoded the QA appliance's numeric ids for
// that tier-0 (28), its group (3), its network integration (5) and its edge
// cluster (a raw NSX-T UUID). Those ids are environment-specific and silently
// wrong anywhere else. The fixtures here look everything up by name instead:
// resolving the tier-0 yields its group and network integration for free, and
// the edge cluster is resolved through that same integration.
package nsxt

import (
	"fmt"
	"os"
)

// EnvTier0RouterName overrides the name of the pre-provisioned NSX-T tier-0
// gateway the fixtures resolve. Set it when the target appliance names its
// tier-0 something other than Tier0RouterName.
const EnvTier0RouterName = "TF_ACC_NSXT_TIER0_NAME"

// EnvEdgeClusterName overrides the name of the NSX-T edge cluster the fixtures
// resolve. Set it when the target appliance names its edge cluster something
// other than EdgeClusterName.
const EnvEdgeClusterName = "TF_ACC_NSXT_EDGE_CLUSTER_NAME"

// EnvBgpSourceAddress overrides the BGP neighbor source address. Set it
// alongside EnvTier0RouterName: the address must be an IP on an interface of
// whichever tier-0 the fixtures resolve.
const EnvBgpSourceAddress = "TF_ACC_NSXT_BGP_SOURCE_ADDRESS"

// Tier0RouterName is the name of the pre-provisioned, fully-realized NSX-T
// tier-0 gateway (BGP enabled, with an associated edge cluster and local AS)
// that the NSX-T fixtures build on.
//
// Known limitation: the network router data source resolves a name by listing
// routers, and the list endpoint is capped at 25 results (the OpenAPI spec does
// not expose the max/offset/phrase query parameters the API supports). If the
// target appliance accumulates enough routers that this tier-0 falls outside
// the first page, the lookup fails with "no network router found" rather than
// returning the wrong router. Sweeping leaked test routers keeps it in range.
const Tier0RouterName = "Terraform-NSX-T0"

// EdgeClusterName is the name of the NSX-T edge cluster on the same network
// integration as Tier0RouterName. Morpheus sets NetworkEdgeCluster.name from
// the NSX-T display_name, and exposes the UUID that gateway configs actually
// want as external_id, so the fixtures look the cluster up by this name and
// read external_id off it.
const EdgeClusterName = "qa-edge-cluster-01"

// BgpSourceAddress is an IP on an interface of Tier0RouterName. NSX-T requires
// a source address for EBGP multihop neighbors; without it the create fails
// with "BGP neighbor source address is mandatory for EBGP Multihop."
//
// This is paired with Tier0RouterName: point the fixtures at a different tier-0
// and this must move with it.
const BgpSourceAddress = "10.100.10.1"

// Terraform addresses emitted by the fixtures, for tests to reference.
const (
	// Tier0DataSource is the tier-0 gateway resolved by Tier0Config.
	Tier0DataSource = "data.hpe_morpheus_network_router.nsxt_tier0"

	// Tier0RouterIDRef resolves to the tier-0's Morpheus id.
	Tier0RouterIDRef = Tier0DataSource + ".id"

	// Tier1Resource is the per-test tier-1 gateway created by Tier1Config.
	Tier1Resource = "hpe_morpheus_network_router.nsxt_tier1"

	// Tier1RouterIDRef resolves to the per-test tier-1's Morpheus id.
	Tier1RouterIDRef = Tier1Resource + ".id"

	// Tier1ProviderIDRef resolves to the per-test tier-1's NSX-T policy path,
	// which is what a load balancer's tier1_gateway expects.
	Tier1ProviderIDRef = "data.hpe_morpheus_network_router.nsxt_tier1.provider_id"
)

// Tier0RouterNameValue returns the tier-0 gateway name to resolve, honouring
// EnvTier0RouterName.
func Tier0RouterNameValue() string {
	if v := os.Getenv(EnvTier0RouterName); v != "" {
		return v
	}

	return Tier0RouterName
}

// EdgeClusterNameValue returns the edge cluster name to resolve, honouring
// EnvEdgeClusterName.
func EdgeClusterNameValue() string {
	if v := os.Getenv(EnvEdgeClusterName); v != "" {
		return v
	}

	return EdgeClusterName
}

// BgpSourceAddressValue returns the BGP neighbor source address, honouring
// EnvBgpSourceAddress.
func BgpSourceAddressValue() string {
	if v := os.Getenv(EnvBgpSourceAddress); v != "" {
		return v
	}

	return BgpSourceAddress
}

// Tier0Config renders a data source resolving the pre-provisioned NSX-T tier-0
// gateway by name.
//
// Use Tier0RouterIDRef to reference its id. The data source also exposes
// group.id and network_integration.id, which Tier1Config feeds into the tier-1
// it creates, so no group or integration id ever has to be hardcoded.
func Tier0Config() string {
	return `
data "hpe_morpheus_network_router" "nsxt_tier0" {
  name = "` + Tier0RouterNameValue() + `"
}
`
}

// Tier1Config renders a per-test NSX-T tier-1 gateway named "<name>-tier1",
// connected to the tier-0 resolved by Tier0Config and placed on that tier-0's
// group, network integration and edge cluster.
//
// The returned config includes Tier0Config, so callers must not also emit it.
// It also includes a data source reading the created tier-1 back, exposing
// Tier1ProviderIDRef for consumers that need the gateway's NSX-T policy path.
//
// A tier-1 must be connected to a tier-0 and have an edge cluster before NSX-T
// will realize NAT rules or a load balancer service on it. NSX-T also permits
// only one load balancer service per tier-1, so each test creates its own
// rather than sharing one.
func Tier1Config(name string) string {
	return Tier0Config() + `
data "hpe_morpheus_network_edge_cluster" "nsxt_edge" {
  network_server_id = data.hpe_morpheus_network_router.nsxt_tier0.network_integration.id
  name              = "` + EdgeClusterNameValue() + `"
}

resource "hpe_morpheus_network_router" "nsxt_tier1" {
  name                   = "` + name + `-tier1"
  group_id               = data.hpe_morpheus_network_router.nsxt_tier0.group.id
  network_integration_id = data.hpe_morpheus_network_router.nsxt_tier0.network_integration.id

  config_nsxt_gateway_tier1 = {
    ip_management_type = "dhcpLocal"
    edge_cluster       = data.hpe_morpheus_network_edge_cluster.nsxt_edge.external_id
    fail_over          = "NON_PREEMPTIVE"
    tier0_gateway      = data.hpe_morpheus_network_router.nsxt_tier0.provider_id
  }
}

data "hpe_morpheus_network_router" "nsxt_tier1" {
  id = hpe_morpheus_network_router.nsxt_tier1.id
}
`
}

// BgpSourceAddressBlock renders the config_nsxt block a BGP neighbor needs to
// satisfy NSX-T's EBGP multihop source address requirement.
func BgpSourceAddressBlock() string {
	return fmt.Sprintf(`
  config_nsxt = {
    source_addresses = [%q]
  }
`, BgpSourceAddressValue())
}
