// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package virtualswitch

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Traffic types, grouped as Morpheus's VirtualSwitchService groups them.
const (
	trafficVM            = "vm"
	trafficManagement    = "management"
	trafficDataNfs       = "data-nfs"
	trafficDataIscsi     = "data-iscsi"
	trafficLiveMigration = "live-migration"
	trafficSDN           = "sdn"
)

const maxSegmentsPerDataType = 8

// dataTrafficTypes are the segment types Morpheus limits to 8 per switch, with
// at most one untagged and no duplicate VLAN.
var dataTrafficTypes = map[string]bool{
	trafficDataNfs:   true,
	trafficDataIscsi: true,
	trafficSDN:       true,
}

// ValidateConfig rejects, at plan time and on known values only, the
// configurations Morpheus's VirtualSwitchService.validateSegmentsAndUplinks
// would reject at the API, naming the offending attribute instead of surfacing
// a 400 at apply. Rules that need appliance state (name uniqueness, VLAN vs
// existing networks, switches-per-cluster, cross-switch iSCSI VLAN, one
// live-migration per cluster, layout version) stay API-time.
func (r *Resource) ValidateConfig(
	ctx context.Context,
	req resource.ValidateConfigRequest,
	resp *resource.ValidateConfigResponse,
) {
	var config VirtualSwitchModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	segs, d := segmentInputs(ctx, config)
	if resp.Diagnostics.Append(d...); resp.Diagnostics.HasError() {
		return
	}

	ups, d := uplinkInputs(ctx, config)
	if resp.Diagnostics.Append(d...); resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(validateSegmentsAndUplinks(config.BaseType, config.BondMode, segs, ups)...)
}

// validateSegmentsAndUplinks is the condition table, separated from the
// framework plumbing so it can be unit-tested directly.
func validateSegmentsAndUplinks(
	baseType types.String,
	bondMode types.String,
	segs []segmentInput,
	ups []uplinkInput,
) diag.Diagnostics {
	var diags diag.Diagnostics

	present := map[string]bool{}
	for _, s := range segs {
		present[s.trafficType] = true
	}

	hasDataIscsi := present[trafficDataIscsi]
	hasSDN := present[trafficSDN]
	hasLiveMigration := present[trafficLiveMigration]

	otherDataBeside := func(self string) bool {
		for t := range dataTrafficTypes {
			if t != self && present[t] {
				return true
			}
		}

		return false
	}

	// iscsiSingleNic: a data-iscsi segment forbids a bond and any nic2.
	if hasDataIscsi {
		if !bondMode.IsNull() && !bondMode.IsUnknown() {
			diags.AddAttributeError(path.Root("bond_mode"),
				"Bonding not allowed with data-iscsi",
				"bond_mode must be unset when a data-iscsi segment is present: iSCSI traffic "+
					"requires a single, unbonded NIC (iscsiSingleNic).")
		}

		for _, u := range ups {
			if u.nic2 != nil {
				diags.AddAttributeError(path.Root("uplinks"),
					"Second NIC not allowed with data-iscsi",
					fmt.Sprintf("uplink for server_id %d sets nic2, but a data-iscsi segment "+
						"requires single-NIC uplinks (iscsiSingleNic). Remove nic2.", u.serverID))

				break
			}
		}
	}

	// iscsiExclusive: data-iscsi may not share the switch with another data type.
	if hasDataIscsi && otherDataBeside(trafficDataIscsi) {
		diags.AddAttributeError(path.Root("segments"),
			"data-iscsi must be exclusive",
			"a data-iscsi segment cannot share a virtual switch with another data traffic "+
				"type (data-nfs or sdn) (iscsiExclusive).")
	}

	// sdnExclusive: sdn may not share the switch with another data type.
	if hasSDN && otherDataBeside(trafficSDN) {
		diags.AddAttributeError(path.Root("segments"),
			"sdn must be exclusive",
			"an sdn segment cannot share a virtual switch with another data traffic type "+
				"(data-nfs or data-iscsi) (sdnExclusive).")
	}

	// liveMigrationOnData: live-migration requires base_type != iscsi and no
	// data-iscsi segment.
	if hasLiveMigration {
		if !baseType.IsNull() && !baseType.IsUnknown() && baseType.ValueString() == "iscsi" {
			diags.AddAttributeError(path.Root("segments"),
				"live-migration not allowed on an iscsi switch",
				"a live-migration segment requires base_type other than \"iscsi\" "+
					"(liveMigrationOnData).")
		}

		if hasDataIscsi {
			diags.AddAttributeError(path.Root("segments"),
				"live-migration not allowed with data-iscsi",
				"a live-migration segment cannot coexist with a data-iscsi segment "+
					"(liveMigrationOnData).")
		}
	}

	diags.Append(validatePerDataType(segs)...)
	diags.Append(validateSingletons(segs)...)
	diags.Append(validateNetmaskRequired(segs)...)
	diags.Append(validateManagementAddress(segs)...)

	return diags
}

// validateManagementAddress enforces that a management segment carries an IP
// address (ip_address or a non-empty host_ips). The host's hvmcli refuses
// "virtswitch add-segment --traffic-type management" without --ip: the create
// proceeds on every host, the management segment is then refused, and Morpheus
// rolls the whole switch back with an opaque "Virtual switch create on <host>:
// error". The accompanying netmask requirement is covered by
// validateNetmaskRequired. Only management is enforced here; live-migration may
// carry the same host requirement, but that is unverified against a real
// appliance, so it is left to the API for now.
func validateManagementAddress(segs []segmentInput) diag.Diagnostics {
	var diags diag.Diagnostics

	for _, s := range segs {
		if s.trafficType != trafficManagement || s.addressUnknown {
			continue
		}

		hasIP := s.ipAddress != nil || len(s.hostIps) > 0
		if !hasIP {
			diags.AddAttributeError(path.Root("segments"),
				"Management segment needs an address",
				fmt.Sprintf("segment %q is a management segment but sets no IP address; "+
					"set ip_address or a non-empty host_ips (and netmask). The host requires "+
					"an IP for a management segment, and Morpheus would otherwise roll back "+
					"the create.", s.key))
		}
	}

	return diags
}

// validatePerDataType enforces, per data traffic type: at most 8 segments, at
// most one untagged, and no duplicate VLAN.
func validatePerDataType(segs []segmentInput) diag.Diagnostics {
	var diags diag.Diagnostics

	counts := map[string]int{}
	untagged := map[string]int{}
	vlans := map[string]map[int32]bool{}

	for _, s := range segs {
		if s.trafficTypeUnknown || !dataTrafficTypes[s.trafficType] {
			continue
		}

		counts[s.trafficType]++

		if s.vlanUnknown {
			// Not yet known: neither untagged nor a VLAN to check for duplicates.
			continue
		}

		if s.vlanID == nil {
			untagged[s.trafficType]++
		} else {
			if vlans[s.trafficType] == nil {
				vlans[s.trafficType] = map[int32]bool{}
			}

			if vlans[s.trafficType][*s.vlanID] {
				diags.AddAttributeError(path.Root("segments"),
					"Duplicate VLAN",
					fmt.Sprintf("two %s segments share VLAN %d; each %s segment needs a "+
						"distinct VLAN (duplicateVlan).", s.trafficType, *s.vlanID, s.trafficType))
			}

			vlans[s.trafficType][*s.vlanID] = true
		}
	}

	for t, c := range counts {
		if c > maxSegmentsPerDataType {
			diags.AddAttributeError(path.Root("segments"),
				"Too many segments",
				fmt.Sprintf("a virtual switch allows at most %d %s segments, but %d are "+
					"configured (maxSegments).", maxSegmentsPerDataType, t, c))
		}
	}

	for t, c := range untagged {
		if c > 1 {
			diags.AddAttributeError(path.Root("segments"),
				"Too many untagged segments",
				fmt.Sprintf("at most one %s segment may be untagged (omit vlan_id), but %d "+
					"are (vlanRequired).", t, c))
		}
	}

	return diags
}

// validateSingletons enforces at most one live-migration and one management
// segment.
func validateSingletons(segs []segmentInput) diag.Diagnostics {
	var diags diag.Diagnostics

	counts := map[string]int{}
	for _, s := range segs {
		if s.trafficTypeUnknown {
			continue
		}

		counts[s.trafficType]++
	}

	if counts[trafficLiveMigration] > 1 {
		diags.AddAttributeError(path.Root("segments"),
			"Too many live-migration segments",
			"a virtual switch allows at most one live-migration segment.")
	}

	if counts[trafficManagement] > 1 {
		diags.AddAttributeError(path.Root("segments"),
			"Too many management segments",
			"a virtual switch allows at most one management segment.")
	}

	return diags
}

// validateNetmaskRequired enforces that a segment with an ip_address or any
// host_ips also sets netmask.
func validateNetmaskRequired(segs []segmentInput) diag.Diagnostics {
	var diags diag.Diagnostics

	for _, s := range segs {
		if s.netmaskUnknown {
			continue
		}

		hasIP := s.ipAddress != nil || len(s.hostIps) > 0
		if hasIP && s.netmask == nil {
			diags.AddAttributeError(path.Root("segments"),
				"Missing netmask",
				fmt.Sprintf("segment %q sets an IP address (ip_address or host_ips) but no "+
					"netmask; a netmask is required whenever an IP is assigned.", s.key))
		}
	}

	return diags
}
