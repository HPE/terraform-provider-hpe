// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package virtualswitch

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// markComputedUnknownForUpdate rewrites the plan's uplinks set and segments map
// so every server-side, volatile computed sub-attribute is unknown.
//
// WHY: uplinks is a set and segments entries carry computed fields. Terraform
// compares set elements as whole values and requires each planned element to
// correlate with an element in the apply result. The framework's default plan
// copies computed fields forward from prior state, but an async switch mutation
// can change them (e.g. an uplink's sync_message flips from "Interface
// configuration applied successfully." to "" after a rename). The planned
// element then fails to correlate with the actual element, producing
// "Provider produced inconsistent result after apply". Marking the volatile
// fields unknown tells Terraform these values are not known until apply, so no
// correlation is attempted on them. Unknown values inside set/map elements are
// permitted; only the whole element must not be needed for plan-time
// correlation, and the user-configured identity fields (server_id, nic1, nic2,
// traffic_type, vlan_id, ...) are left exactly as planned so correlation still
// works on them.
func markComputedUnknownForUpdate(
	ctx context.Context,
	plan *VirtualSwitchModel,
	keptSegIDs map[string]int64,
) diag.Diagnostics {
	var diags diag.Diagnostics

	ups, d := unknownComputedUplinks(ctx, plan.Uplinks)
	if diags.Append(d...); diags.HasError() {
		return diags
	}

	plan.Uplinks = ups

	segs, d := unknownComputedSegments(ctx, plan.Segments, keptSegIDs)
	if diags.Append(d...); diags.HasError() {
		return diags
	}

	plan.Segments = segs

	return diags
}

// unknownComputedUplinks rebuilds the uplinks set with id, server_name,
// sync_message and sync_status unknown, preserving the configured fields.
func unknownComputedUplinks(ctx context.Context, set types.Set) (types.Set, diag.Diagnostics) {
	var diags diag.Diagnostics

	upType := UplinksValue{}.Type(ctx)

	if set.IsNull() || set.IsUnknown() {
		return set, diags
	}

	var vals []UplinksValue
	if d := set.ElementsAs(ctx, &vals, false); d.HasError() {
		return set, d
	}

	elems := make([]attr.Value, 0, len(vals))
	for _, u := range vals {
		u.Id = types.Int64Unknown()
		u.ServerName = types.StringUnknown()
		u.SyncMessage = types.StringUnknown()
		u.SyncStatus = types.StringUnknown()
		u.state = attr.ValueStateKnown
		elems = append(elems, u)
	}

	out, d := types.SetValue(upType, elems)
	diags.Append(d...)

	return out, diags
}

// unknownComputedSegments rebuilds the segments map with each entry's computed
// fields unknown, preserving the configured fields and the map keys. A segment
// whose key is in keptIDs already exists and keeps that id: ids are stable, and
// a network built on the segment references it — unknown, it would plan that
// network for replacement on every edit of the switch.
func unknownComputedSegments(ctx context.Context, m types.Map, keptIDs map[string]int64) (types.Map, diag.Diagnostics) {
	var diags diag.Diagnostics

	segType := SegmentsValue{}.Type(ctx)

	if m.IsNull() || m.IsUnknown() {
		return m, diags
	}

	vals := make(map[string]SegmentsValue, len(m.Elements()))
	if d := m.ElementsAs(ctx, &vals, false); d.HasError() {
		return m, d
	}

	elems := make(map[string]attr.Value, len(vals))
	for key, s := range vals {
		if id, ok := keptIDs[key]; ok {
			s.Id = types.Int64Value(id)
		} else {
			s.Id = types.Int64Unknown()
		}

		s.InterfaceName = types.StringUnknown()
		s.OvsBridge = types.StringUnknown()
		s.OvsLibvirtNetwork = types.StringUnknown()
		s.OvsUuid = types.StringUnknown()
		s.Status = types.StringUnknown()
		s.StatusMessage = types.StringUnknown()
		// Morpheus' cluster refresh backfills a segment's ipAddress from its
		// first host entry when none was configured. UseStateForUnknown covers
		// segments that already exist; a segment added on this update has no
		// prior state, so its null ip_address must be unknown here or the read
		// back after apply is reported as an inconsistent result.
		if s.IpAddress.IsNull() && !s.HostIps.IsNull() && !s.HostIps.IsUnknown() && len(s.HostIps.Elements()) > 0 {
			s.IpAddress = types.StringUnknown()
		}
		// Likewise, a segment added through the composite edit comes back with
		// mtu set to the switch's mtu, whereas a segment created with the switch
		// keeps a null mtu. An unset mtu is therefore unknown on update.
		if s.Mtu.IsNull() {
			s.Mtu = types.Int64Unknown()
		}
		s.state = attr.ValueStateKnown
		elems[key] = s
	}

	out, d := types.MapValue(segType, elems)
	diags.Append(d...)

	return out, diags
}
