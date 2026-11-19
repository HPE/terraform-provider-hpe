// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package networkconfig

// pluginTypes is the hand-maintained companion of the generated knownTypes
// table: network types that a Morpheus plugin registers, transcribed from the
// plugin's source because no reference appliance the catalogue is harvested
// from has the plugin installed. Each entry names the plugin, its version and
// the source it was read from, so the entry can be checked when the plugin
// changes. An appliance that has the plugin reports the type through
// GET /api/network-types like any other, and the live shape (FromAPI) takes
// precedence, exactly as for the generated entries.
//
// A code present in both tables is a bug: the plugin's type has made it into
// the harvested catalogue and the entry here must be removed. The catalogue
// test enforces that.
var pluginTypes = map[string]TypeSpec{
	// HPE Morpheus SDN network plugin for HPE VM Essentials (morpheus-juniper-sdn-plugin 1.1.0,
	// HpeJuniperSdnNetworkProvider.getNetworkTypes and buildNetworkOptionTypes). A virtual
	// network on a Juniper CN2 controller, scoped to an HVM cluster through zonePool.id; the
	// plugin computes the CIDR and gateway from the pool and manages DHCP itself, so only the
	// pool and name are taken. Nothing is editable after create (see the Flags comment on
	// gateway).
	"juniperSdn": {
		Code: "juniperSdn", Name: "HPE Morpheus SDN Virtual Network", Category: "", ExternalType: "",
		Creatable: true, Deletable: true, HasNetworkServer: true, HasCidr: true, CidrRequired: false,
		Overlay: false, HasStaticRoutes: false,
		// The plugin declares gatewayEditable but its update path refuses a
		// gateway change ("Gateway cannot be modified after network creation",
		// verified live against plugin 1.1.0), so the provider follows the
		// behaviour, not the declaration: a gateway change replaces the network.
		Flags: TypeFlags{
			Name: false, Cidr: false, Gateway: false, VlanID: false, DhcpServer: false, DNS: false,
			StaticOverride: false, NetworkDomain: false, CanAssignPool: true,
		},
		Fields: []FieldSpec{
			{
				Name: "zonePool.id", Context: ContextNetwork, Input: "select", Label: "Resource Pool",
				Required: true, Creatable: true, Editable: false, OptionSource: "clusterPoolsForNetworkType",
			},
		},
	},
}
