// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkconfig

import (
	"fmt"
	"sort"
	"strings"
)

// BlockField maps one attribute of a typed block to the key Morpheus reads
// it under in the network's config map.
type BlockField struct {
	// Attr is the attribute name inside the block, e.g. "vlan_trunks".
	Attr string
	// Key is the config-map key, i.e. the option type's field name, e.g.
	// "hvmVlanTrunks".
	Key string
}

// Block declares one typed, write-only configuration block of the network
// resource: the network types it serves and how its attributes map to the
// config-context option types of those types. The blocks are hand-written —
// the attribute names are the provider's — and a unit test checks each one
// against the generated catalogue: every key must be a config-context field
// of every type the block serves, and every such field must be either mapped
// or listed in Omit.
type Block struct {
	// Attr is the block's attribute name, e.g. "config_ovs_port_group".
	Attr string
	// Types are the codes of the network types the block serves.
	Types []string
	// Fields are the block's attributes and their config-map keys.
	Fields []BlockField
	// Omit lists config-context fields of the served types that deliberately
	// have no attribute: values the server populates itself, such as the NSX-T
	// manager's version, or that the UI keeps hidden.
	Omit []string
}

// Blocks are the typed configuration blocks, in the order they are documented.
// Omit lists the served types' config-context fields that the server fills in
// itself: the NSX manager's version and project scope are hidden UI fields.
var Blocks = []Block{
	{
		Attr:  "config_ovs_port_group",
		Types: []string{"ovsPortGroup"},
		Fields: []BlockField{
			{Attr: "vlan_trunks", Key: "hvmVlanTrunks"},
			{Attr: "trunk_only", Key: "hvmTrunkOnly"},
		},
	},
	{
		Attr:  "config_nsxt",
		Types: []string{"nsxtLogicalSwitch"},
		Fields: []BlockField{
			{Attr: "connected_gateway", Key: "connectedGateway"},
			{Attr: "vlan_ids", Key: "vlanIDs"},
			{Attr: "subnet_ip_management_type", Key: "subnetIpManagementType"},
			{Attr: "subnet_ip_server_id", Key: "subnetIpServerId"},
			{Attr: "subnet_dhcp_server_address", Key: "subnetDhcpServerAddress"},
			{Attr: "dhcp_range", Key: "dhcpRange"},
			{Attr: "subnet_dhcp_lease_time", Key: "subnetDhcpLeaseTime"},
			{Attr: "subnet_dhcp_server_address_ipv6", Key: "subnetDhcpServerAddressIPv6"},
			{Attr: "dhcp_range_ipv6", Key: "dhcpRangeIPv6"},
			{Attr: "dhcp_excluded_ranges_ipv6", Key: "dhcpExcludedRangesIPv6"},
			{Attr: "subnet_dhcp_lease_time_ipv6", Key: "subnetDhcpLeaseTimeIPv6"},
			{Attr: "dhcp_domain_names_ipv6", Key: "dhcpDomainNamesIPv6"},
			{Attr: "dhcp_preferred_time_ipv6", Key: "dhcpPreferredTimeIPv6"},
			{Attr: "dhcp_sntp_servers_ipv6", Key: "dhcpSntpServersIPv6"},
		},
		Omit: []string{"systemVersion", "project"},
	},
	{
		Attr:  "config_nsx_cloud",
		Types: []string{"nsxCloudLogicalSwitch"},
		Fields: []BlockField{
			{Attr: "connected_gateway", Key: "connectedGateway"},
			{Attr: "vlan_ids", Key: "vlanIDs"},
			{Attr: "subnet_ip_management_type", Key: "subnetIpManagementType"},
			{Attr: "subnet_ip_server_id", Key: "subnetIpServerId"},
			{Attr: "subnet_dhcp_server_address", Key: "subnetDhcpServerAddress"},
			{Attr: "dhcp_range", Key: "dhcpRange"},
			{Attr: "subnet_dhcp_lease_time", Key: "subnetDhcpLeaseTime"},
		},
		Omit: []string{"systemVersion", "project"},
	},
	{
		Attr:  "config_hvm_data",
		Types: []string{"hvmPortGroupDataNetwork"},
		Fields: []BlockField{
			{Attr: "cluster_network_layout", Key: "hvmClusterNetworkLayout"},
			{Attr: "host_network_interface", Key: "hvmHostNetworkInterface"},
			{Attr: "virtual_switch_segment_id", Key: "virtualSwitchSegmentId"},
			{Attr: "mtu", Key: "hvmMtu"},
		},
	},
	{
		Attr:  "config_hvm_overlay",
		Types: []string{"hvmPortGroupOverlayNetwork"},
		Fields: []BlockField{
			{Attr: "overlay_network_interface", Key: "hvmOverlayNetworkInterface"},
			{Attr: "overlay_port", Key: "hvmOverlayPort"},
			{Attr: "overlay_id", Key: "hvmOverlayId"},
		},
	},
	{
		Attr:  "config_hvm_private",
		Types: []string{"hvmPortGroupPrivateNetwork"},
		Fields: []BlockField{
			{Attr: "cluster_network_layout", Key: "hvmClusterNetworkLayout"},
			{Attr: "host_network_interface", Key: "hvmHostNetworkInterface"},
			{Attr: "virtual_switch_id", Key: "virtualSwitchId"},
		},
	},
	{
		Attr:  "config_hvm_sriov",
		Types: []string{"hvmPortGroupSriovNetwork"},
		Fields: []BlockField{
			{Attr: "physical_function", Key: "hvmPhysicalFunction"},
		},
	},
	{
		Attr:  "config_azure",
		Types: []string{"azureNetwork"},
		Fields: []BlockField{
			{Attr: "resource_group_id", Key: "resourceGroupId"},
			{Attr: "subnet_name", Key: "subnetName"},
			{Attr: "subnet_cidr", Key: "subnetCidr"},
		},
	},
	{
		Attr:  "config_gcp",
		Types: []string{"googleNetwork"},
		Fields: []BlockField{
			{Attr: "auto_create", Key: "autoCreate"},
			{Attr: "mtu", Key: "mtu"},
		},
	},
	{
		Attr:  "config_openstack",
		Types: []string{"openstackPrivate", "huaweiPrivate", "opentelekomPrivate"},
		Fields: []BlockField{
			{Attr: "dhcp_range", Key: "dhcpRange"},
			{Attr: "host_routes", Key: "hostRoutes"},
		},
	},
}

// BlockFor returns the block whose attribute is named attr.
func BlockFor(attr string) (Block, bool) {
	for _, b := range Blocks {
		if b.Attr == attr {
			return b, true
		}
	}

	return Block{}, false
}

// BlockForType returns the block that serves the network type with the given
// code, if any. Types without a block take the free-form config map.
func BlockForType(code string) (Block, bool) {
	for _, b := range Blocks {
		for _, t := range b.Types {
			if t == code {
				return b, true
			}
		}
	}

	return Block{}, false
}

// Serves reports whether the block serves the network type with the given code.
func (b Block) Serves(code string) bool {
	for _, t := range b.Types {
		if t == code {
			return true
		}
	}

	return false
}

// Key returns the config-map key of the block attribute named attr.
func (b Block) Key(attr string) (string, bool) {
	for _, f := range b.Fields {
		if f.Attr == attr {
			return f.Key, true
		}
	}

	return "", false
}

// TypeNames lists the block's types as "code (Name)" for messages, in order.
func (b Block) TypeNames() string {
	names := make([]string, 0, len(b.Types))

	for _, code := range b.Types {
		if t, ok := Lookup(code); ok && t.Name != "" {
			names = append(names, fmt.Sprintf("%s (%s)", code, t.Name))
		} else {
			names = append(names, code)
		}
	}

	return strings.Join(names, ", ")
}

// Undeclared returns the block's attributes that are set in values but whose
// key the network type t does not declare as a config-context field. Morpheus
// applies only the option types a type declares, so such a value is accepted
// and ignored: a block attribute that a newer release or plugin added is not
// available on an appliance whose type does not have it yet.
func (b Block) Undeclared(t TypeSpec, values map[string]any) []BlockField {
	var out []BlockField

	for _, f := range b.Fields {
		if _, set := values[f.Attr]; !set {
			continue
		}

		if spec, ok := t.Field(f.Key); !ok || spec.Context != ContextConfig {
			out = append(out, f)
		}
	}

	return out
}

// UndeclaredKeys returns, sorted, the keys of a free-form config map that the
// network type t does not declare as config-context fields.
func UndeclaredKeys(t TypeSpec, config map[string]any) []string {
	var out []string

	for k := range config {
		if spec, ok := t.Field(k); !ok || spec.Context != ContextConfig {
			out = append(out, k)
		}
	}

	sort.Strings(out)

	return out
}

// Route builds the config map Morpheus reads from the write-only inputs and
// returns it with its immutable subset.
//
// generic is the free-form "config" map (nil when unset); block holds the typed
// block's values by attribute name (nil when unset) and b is that block. The
// schema forbids setting both, so in practice one of them is nil; were both
// set, a typed value would override a generic one under the same key.
//
// A key is immutable when t, the network's type, declares it as a
// config-context option type that is not editable: Morpheus accepts a new
// value for it on update and ignores it, so a change must replace the network.
// With t nil — the type is not known at plan time — nothing is immutable: an
// update Morpheus ignores is visible in config_from_api, whereas a replacement
// on a guess would destroy a network for nothing.
func Route(t *TypeSpec, b *Block, generic, block map[string]any) (all, immutable map[string]any) {
	all = make(map[string]any, len(generic)+len(block))
	immutable = map[string]any{}

	for k, v := range generic {
		all[k] = v
	}

	if b != nil {
		for _, f := range b.Fields {
			if v, ok := block[f.Attr]; ok {
				all[f.Key] = v
			}
		}
	}

	if t == nil {
		return all, immutable
	}

	for k, v := range all {
		if f, ok := t.Field(k); ok && f.Context == ContextConfig && !f.Editable {
			immutable[k] = v
		}
	}

	return all, immutable
}

// Check validates the hand-written blocks against the generated catalogue and
// returns every inconsistency, so a re-harvest that changes a type's option
// types fails a unit test rather than silently misrouting a field.
func Check() []string {
	var problems []string

	seenAttr := map[string]bool{}
	seenType := map[string]string{}

	for _, b := range Blocks {
		if seenAttr[b.Attr] {
			problems = append(problems, fmt.Sprintf("%s: declared twice", b.Attr))
		}

		seenAttr[b.Attr] = true

		if len(b.Types) == 0 {
			problems = append(problems, fmt.Sprintf("%s: serves no network type", b.Attr))
		}

		for _, code := range b.Types {
			if other, dup := seenType[code]; dup {
				problems = append(problems, fmt.Sprintf("%s: type %q is already served by %s", b.Attr, code, other))
			}

			seenType[code] = b.Attr

			t, ok := Lookup(code)
			if !ok {
				problems = append(problems, fmt.Sprintf("%s: type %q is not in the catalogue", b.Attr, code))

				continue
			}

			problems = append(problems, checkBlockAgainstType(b, t)...)
		}
	}

	sort.Strings(problems)

	return problems
}

func checkBlockAgainstType(b Block, t TypeSpec) []string {
	var problems []string

	covered := map[string]bool{}

	for _, f := range b.Fields {
		covered[f.Key] = true

		spec, ok := t.Field(f.Key)

		switch {
		case !ok:
			problems = append(problems, fmt.Sprintf("%s.%s: %q is not an option type of %s", b.Attr, f.Attr, f.Key, t.Code))
		case spec.Context != ContextConfig:
			problems = append(problems, fmt.Sprintf("%s.%s: %q is a %s-context field of %s and belongs at the top level",
				b.Attr, f.Attr, f.Key, spec.Context, t.Code))
		}
	}

	for _, key := range b.Omit {
		covered[key] = true

		if _, ok := t.Field(key); !ok {
			problems = append(problems, fmt.Sprintf("%s: omitted %q is not an option type of %s", b.Attr, key, t.Code))
		}
	}

	for _, f := range t.ConfigFields() {
		if !covered[f.Name] {
			problems = append(problems, fmt.Sprintf("%s: config field %q of %s has no attribute and is not omitted",
				b.Attr, f.Name, t.Code))
		}
	}

	return problems
}
