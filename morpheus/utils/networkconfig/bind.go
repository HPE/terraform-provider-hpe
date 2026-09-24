// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkconfig

import "sort"

// How Morpheus applies a network-context field on update is decided by its
// network service, not only by the field's option type. When a network is
// updated, the service binds a fixed list of properties directly: some on
// every type, some only when the type's matching flag is on, one only on
// create. An option type whose field is on that list is applied (or not)
// according to the list, whatever its own editable flag says; only the option
// types not on the list are filtered by editable. The three sets below mirror
// that list.

// alwaysBound are bound directly on every update, for every type.
var alwaysBound = map[string]bool{
	"description": true, "displayName": true, "labelString": true, "labelArray": true,
	"zonePool": true, "zonePool.id": true,
	"networkDomain": true, "networkDomain.id": true,
	"networkProxy": true, "networkProxy.id": true,
	"noProxy": true, "assignPublicIp": true, "allowStaticOverride": true, "scanNetwork": true,
	"active": true, "visibility": true, "searchDomains": true, "applianceUrlProxyBypass": true,
	"guestConsoleJumpHost": true, "guestConsoleJumpPort": true, "guestConsoleJumpUsername": true,
	"guestConsoleJumpPassword": true, "guestConsoleJumpKeypair": true,
}

// flagBound are bound on update only when the type's flag is on; otherwise a
// new value is dropped without error. (The pool fields follow canAssignPool
// but are excluded on create as well when it is off, so they are never
// immutable: they are simply not available on such types.)
var flagBound = map[string]func(TypeFlags) bool{
	"name":             func(f TypeFlags) bool { return f.Name },
	"cidr":             func(f TypeFlags) bool { return f.Cidr },
	"cidrIPv6":         func(f TypeFlags) bool { return f.Cidr },
	"ipv4Enabled":      func(f TypeFlags) bool { return f.Cidr },
	"ipv6Enabled":      func(f TypeFlags) bool { return f.Cidr },
	"gateway":          func(f TypeFlags) bool { return f.Gateway },
	"gatewayIPv6":      func(f TypeFlags) bool { return f.Gateway },
	"vlanId":           func(f TypeFlags) bool { return f.VlanID },
	"dhcpServer":       func(f TypeFlags) bool { return f.DhcpServer },
	"dhcpServerIPv6":   func(f TypeFlags) bool { return f.DhcpServer },
	"dnsPrimary":       func(f TypeFlags) bool { return f.DNS },
	"dnsSecondary":     func(f TypeFlags) bool { return f.DNS },
	"dnsPrimaryIPv6":   func(f TypeFlags) bool { return f.DNS },
	"dnsSecondaryIPv6": func(f TypeFlags) bool { return f.DNS },
}

// createOnly are bound on create and never on update.
var createOnly = map[string]bool{
	"availabilityZone": true,
}

// attributeFor maps a network-context field name to the resource attribute
// that carries it. Fields with no attribute yet are listed in
// unsupportedNetworkFields so the coverage test notices a new one.
var attributeFor = map[string]string{
	"name":                "name",
	"cidr":                "cidr",
	"cidrIPv6":            "cidr_ipv6",
	"ipv4Enabled":         "ipv4enabled",
	"ipv6Enabled":         "ipv6enabled",
	"gateway":             "gateway",
	"gatewayIPv6":         "gateway_ipv6",
	"vlanId":              "vlan_id",
	"dhcpServer":          "dhcp_server",
	"dhcpServerIPv6":      "dhcp_server_ipv6",
	"dnsPrimary":          "dns_primary",
	"dnsSecondary":        "dns_secondary",
	"dnsPrimaryIPv6":      "dns_primary_ipv6",
	"dnsSecondaryIPv6":    "dns_secondary_ipv6",
	"active":              "active",
	"allowStaticOverride": "allow_static_override",
	"assignPublicIp":      "assign_public_ip",
	"networkDomain.id":    "network_domain_id",
	"zonePool.id":         "zone_pool_id",
	"switchId":            "switch_id",
	// Accepted by the API but never returned; carried as write-only scalars.
	"scopeId":          "scope_id",
	"vxlanId":          "vxlan_id",
	"parentNetwork.id": "parent_network_id",
	"availabilityZone": "availability_zone",
}

// unsupportedNetworkFields are network-context option types the resource has
// no attribute for, with the reason.
var unsupportedNetworkFields = map[string]string{
	"bridgeName":       "KVM/OLVM bridge name: no cloud to exercise it against yet",
	"networkPool.pool": "OpenStack pool binding by name; pool_id carries the pool by id",
	"streamType":       "single PowerVC type; no appliance to exercise it against",
	"tenantName":       "single OpenStack type; no appliance to exercise it against",
}

// NetworkAttribute returns the resource attribute that carries the
// network-context field named field, if the resource has one.
func NetworkAttribute(field string) (string, bool) {
	attr, ok := attributeFor[field]

	return attr, ok
}

// ImmutableAttributes returns, sorted, the top-level attributes whose change
// requires replacing a network of type t, because Morpheus would accept the
// update and not apply it: the flag-bound fields whose flag is off, the
// create-only fields, and the type's network-context option types that are
// neither on the bind list nor editable. Only fields the resource has an
// attribute for are returned.
func ImmutableAttributes(t TypeSpec) []string {
	set := map[string]bool{}

	for field, on := range flagBound {
		if !on(t.Flags) {
			if attr, ok := attributeFor[field]; ok {
				set[attr] = true
			}
		}
	}

	for _, f := range t.NetworkFields() {
		if alwaysBound[f.Name] {
			continue
		}

		if _, ok := flagBound[f.Name]; ok {
			continue
		}

		if createOnly[f.Name] || !f.Editable {
			if attr, ok := attributeFor[f.Name]; ok {
				set[attr] = true
			}
		}
	}

	out := make([]string, 0, len(set))
	for attr := range set {
		out = append(out, attr)
	}

	sort.Strings(out)

	return out
}
