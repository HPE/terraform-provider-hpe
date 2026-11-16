// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package virtualswitch

import (
	"context"
	"fmt"
	"net/http"

	"github.com/HPE/terraform-provider-hpe/utils/schemadefaults"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

func (r *Resource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var data VirtualSwitchModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError("creating client failed", err.Error())

		return
	}

	vsResp, httpResp, err := client.VirtualSwitchesAPI.
		GetVirtualSwitch(ctx, data.Id.ValueInt64()).Execute()
	if err != nil || httpResp == nil || httpResp.StatusCode != http.StatusOK {
		if httpResp != nil && httpResp.StatusCode == http.StatusNotFound {
			resp.State.RemoveResource(ctx)

			return
		}

		resp.Diagnostics.AddError(
			"error reading virtual switch",
			fmt.Sprintf("virtual switch %d GET failed: ", data.Id.ValueInt64())+
				errfmt.ErrMsg(err, httpResp),
		)

		return
	}

	// A 200 without a body is a malformed response, not a deleted switch:
	// removing the resource here would orphan it. Only a 404 removes.
	if vsResp == nil || vsResp.VirtualSwitch == nil {
		resp.Diagnostics.AddError(
			"error reading virtual switch",
			fmt.Sprintf("virtual switch %d: the API returned 200 with no virtualSwitch body", data.Id.ValueInt64()),
		)

		return
	}

	state, diags := mapDetailToState(ctx, vsResp.VirtualSwitch, data)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)

	// Fill any schema-declared default the API omitted (null in state) so an
	// imported resource does not plan a change nobody made -- base_type in
	// particular, whose drift would force a replacement.
	resp.Diagnostics.Append(
		schemadefaults.Apply(ctx, VirtualSwitchResourceSchema(ctx), &resp.State)...,
	)
}

// mapDetailToState maps a GET virtual switch detail into Terraform state,
// preserving values that live only in Terraform: the timeouts block, the
// cluster_id (never echoed by GET detail for a switch already read), and the
// user-chosen segment map keys (see carrySegmentKeys).
func mapDetailToState(
	ctx context.Context,
	vs *sdk.GetVirtualSwitch200ResponseVirtualSwitch,
	prior VirtualSwitchModel,
) (VirtualSwitchModel, diag.Diagnostics) {
	var state VirtualSwitchModel
	var diags diag.Diagnostics

	state.Id = convert.Int64ToType(vs.Id)
	state.Name = convert.StrToType(vs.Name)
	state.BaseType = convert.StrToType(vs.BaseType)
	state.Code = convert.StrToType(vs.Code.Get())
	state.Description = convert.StrToType(vs.Description.Get())
	state.BondMode = convert.StrToType(vs.BondMode.Get())
	state.Status = convert.StrToType(vs.Status)
	state.StatusMessage = convert.StrToType(vs.StatusMessage.Get())
	state.Active = convert.BoolToType(vs.Active)
	state.Mtu = int32PtrToInt64Type(vs.Mtu)
	state.NicCount = int32PtrToInt64Type(vs.NicCount)
	state.NetworkCount = int32PtrToInt64Type(vs.NetworkCount)

	// cluster_id is required and RequiresReplace; GET detail may omit it, so
	// carry the prior value forward, falling back to the response.
	if !prior.ClusterId.IsNull() && !prior.ClusterId.IsUnknown() {
		state.ClusterId = prior.ClusterId
	} else {
		state.ClusterId = convert.Int64ToType(vs.ClusterId)
	}

	// timeouts live only in configuration/state.
	state.Timeouts = prior.Timeouts

	segments, segDiags := buildSegmentsMap(ctx, vs.Segments, prior)
	if diags.Append(segDiags...); diags.HasError() {
		return state, diags
	}

	state.Segments = segments

	uplinks, upDiags := buildUplinksSet(ctx, vs.Uplinks)
	if diags.Append(upDiags...); diags.HasError() {
		return state, diags
	}

	state.Uplinks = uplinks

	return state, diags
}

func int32PtrToInt64Type(v *int32) types.Int64 {
	if v == nil {
		return types.Int64Null()
	}

	return types.Int64Value(int64(*v))
}

func nullableInt32ToType(getter func() *int32) types.Int64 {
	v := getter()
	if v == nil {
		return types.Int64Null()
	}

	return types.Int64Value(int64(*v))
}

// buildSegmentsMap maps API segment rows into the mnemonic-keyed map, carrying
// keys forward from prior state (see carrySegmentKeys).
func buildSegmentsMap(
	ctx context.Context,
	rows []sdk.GetVirtualSwitch200ResponseVirtualSwitchAllOfSegmentsInner1,
	prior VirtualSwitchModel,
) (types.Map, diag.Diagnostics) {
	var diags diag.Diagnostics

	segType := SegmentsValue{}.Type(ctx)

	keys := carrySegmentKeys(ctx, rows, prior)

	elements := make(map[string]attr.Value, len(rows))

	// The prior segments, to keep a configured empty host_ips as an empty set
	// (the API returns none; host_ips is not computed, so null would be an
	// inconsistent result).
	priorSegs := map[string]SegmentsValue{}
	if !prior.Segments.IsNull() && !prior.Segments.IsUnknown() {
		_ = prior.Segments.ElementsAs(ctx, &priorSegs, false)
	}

	for i, row := range rows {
		hostIps, d := buildHostIps(ctx, row.HostIps)
		if diags.Append(d...); diags.HasError() {
			return types.MapNull(segType), diags
		}

		if ps, ok := priorSegs[keys[i]]; ok && hostIps.IsNull() &&
			!ps.HostIps.IsNull() && !ps.HostIps.IsUnknown() && len(ps.HostIps.Elements()) == 0 {
			hostIps = ps.HostIps
		}

		seg := SegmentsValue{
			Gateway:           convert.StrToType(row.Gateway.Get()),
			HostIps:           hostIps,
			Id:                convert.Int64ToType(row.Id),
			InterfaceName:     convert.StrToType(row.InterfaceName.Get()),
			IpAddress:         convert.StrToType(row.IpAddress.Get()),
			Mtu:               nullableInt32ToType(row.Mtu.Get),
			Netmask:           convert.StrToType(row.Netmask.Get()),
			OvsBridge:         convert.StrToType(row.OvsBridge.Get()),
			OvsLibvirtNetwork: convert.StrToType(row.OvsLibvirtNetwork.Get()),
			OvsUuid:           convert.StrToType(row.OvsUuid.Get()),
			Status:            convert.StrToType(row.Status),
			StatusMessage:     convert.StrToType(row.StatusMessage.Get()),
			TrafficType:       convert.StrToType(row.TrafficType),
			VlanId:            nullableInt32ToType(row.VlanId.Get),
			state:             attr.ValueStateKnown,
		}

		elements[keys[i]] = seg
	}

	mapVal, d := types.MapValue(segType, elements)
	diags.Append(d...)

	return mapVal, diags
}

func buildHostIps(
	ctx context.Context,
	rows []sdk.GetVirtualSwitch200ResponseVirtualSwitchAllOfSegmentsInner1HostIpsInner,
) (types.Set, diag.Diagnostics) {
	var diags diag.Diagnostics

	hiType := HostIpsValue{}.Type(ctx)

	if len(rows) == 0 {
		return types.SetNull(hiType), diags
	}

	elems := make([]attr.Value, 0, len(rows))

	for _, row := range rows {
		hi := HostIpsValue{
			Ip:       convert.StrToType(row.Ip),
			ServerId: convert.Int64ToType(row.ServerId),
			state:    attr.ValueStateKnown,
		}

		elems = append(elems, hi)
	}

	setVal, d := types.SetValue(hiType, elems)
	diags.Append(d...)

	return setVal, diags
}

func buildUplinksSet(
	ctx context.Context,
	rows []sdk.GetVirtualSwitch200ResponseVirtualSwitchAllOfUplinksInner,
) (types.Set, diag.Diagnostics) {
	var diags diag.Diagnostics

	upType := UplinksValue{}.Type(ctx)

	elems := make([]attr.Value, 0, len(rows))

	for _, row := range rows {
		up := UplinksValue{
			Id:          convert.Int64ToType(row.Id),
			Nic1:        convert.StrToType(row.Nic1),
			Nic2:        convert.StrToType(row.Nic2.Get()),
			ServerId:    convert.Int64ToType(row.ServerId),
			ServerName:  convert.StrToType(row.ServerName.Get()),
			SyncMessage: convert.StrToType(row.SyncMessage.Get()),
			SyncStatus:  convert.StrToType(row.SyncStatus),
			state:       attr.ValueStateKnown,
		}

		elems = append(elems, up)
	}

	setVal, d := types.SetValue(upType, elems)
	diags.Append(d...)

	return setVal, diags
}

// readSwitchState performs a GET and maps it into state, preserving the prior
// model's Terraform-only values (timeouts, cluster_id, segment keys).
func readSwitchState(
	ctx context.Context,
	client *sdk.APIClient,
	id int64,
	prior VirtualSwitchModel,
) (VirtualSwitchModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	vsResp, httpResp, err := client.VirtualSwitchesAPI.GetVirtualSwitch(ctx, id).Execute()
	if err != nil || httpResp == nil || httpResp.StatusCode != http.StatusOK {
		diags.AddError(
			"error reading virtual switch",
			fmt.Sprintf("virtual switch %d GET failed: ", id)+errfmt.ErrMsg(err, httpResp),
		)

		return prior, diags
	}

	if vsResp == nil || vsResp.VirtualSwitch == nil {
		diags.AddError(
			"error reading virtual switch",
			fmt.Sprintf("virtual switch %d GET returned no switch", id),
		)

		return prior, diags
	}

	return mapDetailToState(ctx, vsResp.VirtualSwitch, prior)
}
