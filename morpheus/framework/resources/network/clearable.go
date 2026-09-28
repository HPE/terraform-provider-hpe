// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package network

import "github.com/HPE/terraform-provider-hpe/morpheus/utils/clearable"

// clearableAttributes are the optional attributes a user can remove from the
// configuration to clear on the appliance. ModifyPlan plans the clear, Update
// sends it in the shape each field takes on the wire (a nullable field set to
// null, a plain string sent empty — Morpheus stores both as null — a list sent
// empty), and Create and Update record which of them the configuration set so
// the next plan can tell a removal from a value the server computed.
//
// Whether Morpheus applies a clear to gateway, dns_* or switch_id depends on
// the network type, as for any other change to them: on a type whose flag is
// off the clear is planned as a replacement, by the same rule ModifyPlan
// applies to a changed value.
//
// A clear the server verified as taking effect (a PUT with null, or "", read
// back as null) is Clear. Two are Release — planned unknown, and whatever
// comes back is the new state: connected_gateway is read back from the config
// map, which the server may or may not hold an empty key in; tenant_ids sent
// empty revokes every tenant's access but the owner's, whose permission row
// also carries the network's group access and is re-created on every update,
// so it reads back as the owner alone.
var clearableAttributes = clearable.Attributes{
	// Nullable in the update request: sent as null.
	{Name: "description", Mode: clearable.Clear},
	{Name: "gateway_ipv6", Mode: clearable.Clear},
	{Name: "netmask_ipv6", Mode: clearable.Clear},
	{Name: "dns_primary_ipv6", Mode: clearable.Clear},
	{Name: "dns_secondary_ipv6", Mode: clearable.Clear},
	{Name: "no_proxy", Mode: clearable.Clear},
	{Name: "pool_id", Mode: clearable.Clear},
	{Name: "pool_ipv6_id", Mode: clearable.Clear},
	// Plain strings in the update request: sent empty.
	{Name: "display_name", Mode: clearable.Clear},
	{Name: "gateway", Mode: clearable.Clear},
	{Name: "dns_primary", Mode: clearable.Clear},
	{Name: "dns_secondary", Mode: clearable.Clear},
	{Name: "search_domains", Mode: clearable.Clear},
	{Name: "switch_id", Mode: clearable.Clear},
	// Optional, not computed: the plan already carries the null; Update has
	// to send it (an empty string) or the apply fails.
	{Name: "cidr", Mode: clearable.Clear},
	{Name: "cidr_ipv6", Mode: clearable.Clear},
	// Nullable integer; and associations, sent as {"id": null}.
	{Name: "vlan_id", Mode: clearable.Clear},
	{Name: "network_domain_id", Mode: clearable.Clear},
	{Name: "network_proxy_id", Mode: clearable.Clear},
	{Name: "zone_pool_id", Mode: clearable.Clear},
	// Lists: sent empty.
	{Name: "labels", Mode: clearable.Clear},
	{Name: "tenant_ids", Mode: clearable.Release},
	// Read from the config map; cleared by sending an empty value there.
	{Name: "connected_gateway", Mode: clearable.Release},
}

// poolIsCluster names the network types whose zone_pool_id is the HVM cluster
// the network lives on. The HVM network plugin resolves the cluster from the
// pool for every operation — create, update and delete alike — and fails
// without it ("Unable to identify cluster."), so a clear the API would accept
// breaks the update that follows and leaves the network undeletable.
// ModifyPlan refuses it on these types; on every other type the pool is an
// association Morpheus removes on request.
var poolIsCluster = map[string]bool{
	"ovsPortGroup":               true,
	"ovsBridge":                  true,
	"ovsOverlay":                 true,
	"hvmPortGroupDataNetwork":    true,
	"hvmPortGroupOverlayNetwork": true,
	"hvmPortGroupPrivateNetwork": true,
	"hvmPortGroupSriovNetwork":   true,
}

// notClearable records why each remaining optional attribute is not cleared
// when it is removed from the configuration. The unit test holds the two
// lists against the schema so a new attribute has to be classified.
var notClearable = clearable.NotClearable{
	// Booleans have no cleared state; a removed boolean keeps its value.
	"active":                     "boolean",
	"allow_static_override":      "boolean",
	"appliance_url_proxy_bypass": "boolean",
	"assign_public_ip":           "boolean",
	"dhcp_server":                "boolean",
	"dhcp_server_ipv6":           "boolean",
	"ipv4enabled":                "boolean",
	"ipv6enabled":                "boolean",
	// Accepted on create and never returned, or create-only: the state keeps
	// the configured value, and there is nothing on the server to clear.
	"network_server_id": "never returned by the API; kept from the configuration",
	"scope_id":          "never returned by the API; kept from the configuration",
	"vxlan_id":          "never returned by the API; kept from the configuration",
	"parent_network_id": "never returned by the API; kept from the configuration",
	"availability_zone": "create-only; a change replaces the network",
	// One of the pair is always set; visibility has a default.
	"type_id":    "one of type_id/type_code is always configured",
	"type_code":  "one of type_id/type_code is always configured",
	"visibility": "enum with a default, no empty value",
}
