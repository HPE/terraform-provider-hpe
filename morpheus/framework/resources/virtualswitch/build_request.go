// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package virtualswitch

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
)

// bondModeNone is the wire value that clears the switch's NIC bond.
const bondModeNone = "none"

// segmentInput is a normalised view of one configured segment, key included.
type segmentInput struct {
	key         string
	trafficType string
	vlanID      *int32
	ipAddress   *string
	netmask     *string
	gateway     *string
	mtu         *int32
	hostIps     []hostIPInput

	// Unknown-value markers for ValidateConfig, which runs before values
	// from other resources or data sources are known. A nil pointer above
	// cannot tell "unset" from "not yet known", and a check that read an
	// unknown as unset would misfire (a management segment whose host_ips
	// come from the available-NICs data source, a netmask from a variable,
	// two computed vlan_ids). The request builders never see unknowns.
	trafficTypeUnknown bool
	vlanUnknown        bool
	addressUnknown     bool // ip_address or host_ips (set or an element)
	netmaskUnknown     bool
	mtuUnknown         bool
}

type hostIPInput struct {
	serverID int64
	ip       string
}

type uplinkInput struct {
	serverID int64
	nic1     string
	nic2     *string
}

func int64ToInt32Ptr(v types.Int64) *int32 {
	if v.IsNull() || v.IsUnknown() {
		return nil
	}

	//nolint:gosec // schema validators bound vlan_id/mtu well within int32.
	i := int32(v.ValueInt64())

	return &i
}

func strToPtr(v types.String) *string {
	if v.IsNull() || v.IsUnknown() {
		return nil
	}

	s := v.ValueString()

	return &s
}

// segmentInputs extracts the configured/planned segments as a stable slice.
func segmentInputs(ctx context.Context, m VirtualSwitchModel) ([]segmentInput, diag.Diagnostics) {
	var diags diag.Diagnostics

	if m.Segments.IsNull() || m.Segments.IsUnknown() {
		return nil, diags
	}

	segsMap := make(map[string]SegmentsValue, len(m.Segments.Elements()))
	if d := m.Segments.ElementsAs(ctx, &segsMap, false); d.HasError() {
		diags.Append(d...)

		return nil, diags
	}

	out := make([]segmentInput, 0, len(segsMap))

	for key, seg := range segsMap {
		hostIps, d := hostIPInputs(ctx, seg.HostIps)
		if diags.Append(d...); diags.HasError() {
			return nil, diags
		}

		out = append(out, segmentInput{
			key:         key,
			trafficType: seg.TrafficType.ValueString(),
			vlanID:      int64ToInt32Ptr(seg.VlanId),
			ipAddress:   strToPtr(seg.IpAddress),
			netmask:     strToPtr(seg.Netmask),
			gateway:     strToPtr(seg.Gateway),
			mtu:         int64ToInt32Ptr(seg.Mtu),
			hostIps:     hostIps,

			trafficTypeUnknown: seg.TrafficType.IsUnknown(),
			vlanUnknown:        seg.VlanId.IsUnknown(),
			addressUnknown:     seg.IpAddress.IsUnknown() || hostIPsUnknown(seg.HostIps),
			netmaskUnknown:     seg.Netmask.IsUnknown(),
			mtuUnknown:         seg.Mtu.IsUnknown(),
		})
	}

	return out, diags
}

// hostIPsUnknown reports whether the host_ips set, or any element or element
// attribute in it, is not yet known.
func hostIPsUnknown(set types.Set) bool {
	if set.IsUnknown() {
		return true
	}

	for _, el := range set.Elements() {
		if el.IsUnknown() {
			return true
		}

		if obj, ok := el.(types.Object); ok {
			for _, v := range obj.Attributes() {
				if v.IsUnknown() {
					return true
				}
			}
		}
	}

	return false
}

func hostIPInputs(ctx context.Context, set types.Set) ([]hostIPInput, diag.Diagnostics) {
	var diags diag.Diagnostics

	if set.IsNull() || set.IsUnknown() {
		return nil, diags
	}

	vals := make([]HostIpsValue, 0, len(set.Elements()))
	if d := set.ElementsAs(ctx, &vals, false); d.HasError() {
		diags.Append(d...)

		return nil, diags
	}

	out := make([]hostIPInput, 0, len(vals))
	for _, v := range vals {
		out = append(out, hostIPInput{
			serverID: v.ServerId.ValueInt64(),
			ip:       v.Ip.ValueString(),
		})
	}

	return out, diags
}

func uplinkInputs(ctx context.Context, m VirtualSwitchModel) ([]uplinkInput, diag.Diagnostics) {
	ups, _, diags := uplinkInputsKnown(ctx, m)

	return ups, diags
}

// uplinkInputsKnown is uplinkInputs that also reports whether every uplink
// is known: false when the set itself, or any uplink's server_id, is not yet
// known (a deferred data source), in which case the plan cannot be judged.
func uplinkInputsKnown(ctx context.Context, m VirtualSwitchModel) ([]uplinkInput, bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	if m.Uplinks.IsNull() {
		return nil, true, diags
	}
	if m.Uplinks.IsUnknown() {
		return nil, false, diags
	}

	vals := make([]UplinksValue, 0, len(m.Uplinks.Elements()))
	if d := m.Uplinks.ElementsAs(ctx, &vals, false); d.HasError() {
		diags.Append(d...)

		return nil, false, diags
	}

	known := true
	out := make([]uplinkInput, 0, len(vals))
	for _, v := range vals {
		if v.IsUnknown() || v.ServerId.IsUnknown() {
			known = false
		}
		out = append(out, uplinkInput{
			serverID: v.ServerId.ValueInt64(),
			nic1:     v.Nic1.ValueString(),
			nic2:     strToPtr(v.Nic2),
		})
	}

	return out, known, diags
}

// buildCreateRequest assembles the POST body from the plan.
func buildCreateRequest(
	m VirtualSwitchModel,
	segs []segmentInput,
	ups []uplinkInput,
) sdk.AddVirtualSwitchRequest {
	vs := &sdk.AddVirtualSwitchRequestVirtualSwitch{
		Name: m.Name.ValueString(),
	}

	vs.Code = strToPtr(m.Code)
	vs.Description = strToPtr(m.Description)
	vs.BaseType = strToPtr(m.BaseType)
	vs.BondMode = strToPtr(m.BondMode)
	vs.Mtu = int64ToInt32Ptr(m.Mtu)

	for _, s := range segs {
		seg := sdk.AddVirtualSwitchRequestVirtualSwitchSegmentsInner{
			TrafficType: s.trafficType,
			VlanId:      s.vlanID,
			IpAddress:   s.ipAddress,
			Netmask:     s.netmask,
			Gateway:     s.gateway,
			Mtu:         s.mtu,
		}

		for _, h := range s.hostIps {
			serverID := h.serverID
			ip := h.ip
			seg.HostIps = append(seg.HostIps,
				sdk.AddVirtualSwitchRequestVirtualSwitchSegmentsInnerHostIpsInner{
					ServerId: &serverID,
					Ip:       &ip,
				})
		}

		vs.Segments = append(vs.Segments, seg)
	}

	for _, u := range ups {
		up := sdk.AddVirtualSwitchRequestVirtualSwitchUplinksInner{
			ServerId: u.serverID,
			Nic1:     u.nic1,
		}
		if u.nic2 != nil {
			up.Nic2.Set(u.nic2)
		}

		vs.Uplinks = append(vs.Uplinks, up)
	}

	return sdk.AddVirtualSwitchRequest{VirtualSwitch: vs}
}

// buildCompositeEdit assembles the PUT .../edit body from an update plan.
func buildCompositeEdit(p updatePlan) sdk.CompositeEditVirtualSwitchRequest {
	var body sdk.CompositeEditVirtualSwitchRequest

	if p.bondModeChanged || p.mtuChanged || p.uplinksChanged {
		su := &sdk.CompositeEditVirtualSwitchRequestSwitchUpdate{}
		if p.bondModeChanged {
			su.BondMode = p.bondMode
		}

		if p.mtuChanged {
			su.Mtu = p.mtu
		}

		if p.uplinksChanged {
			for _, u := range p.uplinks {
				up := sdk.CompositeEditVirtualSwitchRequestSwitchUpdateUplinksInner{
					ServerId: u.serverID,
					Nic1:     u.nic1,
				}
				if u.nic2 != nil {
					up.Nic2.Set(u.nic2)
				}

				su.Uplinks = append(su.Uplinks, up)
			}
		}

		body.SwitchUpdate = su
	}

	for _, s := range p.newSegments {
		seg := sdk.CompositeEditVirtualSwitchRequestNewSegmentsInner{
			TrafficType: s.trafficType,
			VlanId:      s.vlanID,
			IpAddress:   s.ipAddress,
			Netmask:     s.netmask,
			Gateway:     s.gateway,
			Mtu:         s.mtu,
		}

		for _, h := range s.hostIps {
			serverID := h.serverID
			ip := h.ip
			seg.HostIps = append(seg.HostIps,
				sdk.CompositeEditVirtualSwitchRequestNewSegmentsInnerHostIpsInner{
					ServerId: &serverID,
					Ip:       &ip,
				})
		}

		body.NewSegments = append(body.NewSegments, seg)
	}

	for _, e := range p.editedSegments {
		segID := e.segID
		// Morpheus skips an edited segment that carries no trafficType before it
		// even resolves segId, yet still starts a process that never completes.
		seg := sdk.CompositeEditVirtualSwitchRequestEditedSegmentsInner{
			SegId:       &segID,
			TrafficType: e.segment.trafficType,
		}
		seg.VlanId.Set(e.segment.vlanID)
		seg.IpAddress.Set(e.segment.ipAddress)
		seg.Netmask.Set(e.segment.netmask)
		seg.Gateway.Set(e.segment.gateway)

		// An edited segment carries its whole host_ips set; an empty (non-nil)
		// slice serialises as "hostIps": [], which Morpheus's containsKey check
		// treats as "clear the per-host addresses". A nil slice would be
		// omitted and the old addresses kept.
		seg.HostIps = make([]sdk.CompositeEditVirtualSwitchRequestEditedSegmentsInnerHostIpsInner, 0, len(e.segment.hostIps))
		for _, h := range e.segment.hostIps {
			serverID := h.serverID
			ip := h.ip
			seg.HostIps = append(seg.HostIps,
				sdk.CompositeEditVirtualSwitchRequestEditedSegmentsInnerHostIpsInner{
					ServerId: &serverID,
					Ip:       &ip,
				})
		}

		body.EditedSegments = append(body.EditedSegments, seg)
	}

	return body
}

// buildRename assembles the PUT .../rename body.
func buildRename(name string) sdk.RenameVirtualSwitchRequest {
	return sdk.RenameVirtualSwitchRequest{Name: name}
}

// buildDescriptionUpdate assembles a PUT /{id} body carrying only description
// (compositeEdit ignores description).
func buildDescriptionUpdate(description *string) sdk.UpdateVirtualSwitchRequest {
	vs := &sdk.UpdateVirtualSwitchRequestVirtualSwitch{Description: description}
	if description == nil {
		// Description is an omitempty pointer, so a nil would drop the key and
		// Morpheus (which only writes the field when the key is present) would
		// keep the old value. Send an explicit null to clear it.
		vs.AdditionalProperties = map[string]interface{}{"description": nil}
	}

	return sdk.UpdateVirtualSwitchRequest{VirtualSwitch: vs}
}
