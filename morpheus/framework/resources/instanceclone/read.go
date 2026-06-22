// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package instanceclone

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

func (r *Resource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var state InstanceCloneModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"read instance clone",
			"failed to create client: "+err.Error(),
		)

		return
	}

	instanceID := state.Id.ValueInt64()

	getResp, hresp, err := client.InstancesAPI.GetInstance(ctx, instanceID).Execute()
	if err != nil {
		if hresp != nil && hresp.StatusCode == http.StatusNotFound {
			tflog.Warn(ctx, fmt.Sprintf("Instance clone %d not found, removing from state",
				instanceID))
			resp.State.RemoveResource(ctx)

			return
		}

		resp.Diagnostics.AddError(
			"read instance clone",
			fmt.Sprintf("instance %d GET failed: %s", instanceID, errfmt.ErrMsg(err, hresp)),
		)

		return
	}

	if getResp == nil || getResp.Instance == nil {
		tflog.Warn(ctx, fmt.Sprintf("Instance clone %d returned nil, removing from state",
			instanceID))
		resp.State.RemoveResource(ctx)

		return
	}

	inst := getResp.Instance

	// On import the source_instance_id is unknown because the clone API does not
	// return it as a first-class field. Morpheus does stamp the source id onto
	// the clone's config as "cloneInstanceId", so recover it best-effort when the
	// current state has no source_instance_id (i.e. after import). On a normal
	// refresh the value is already set and is left untouched (it is
	// RequiresReplace, so it must never be clobbered here).
	if state.SourceInstanceId.IsNull() || state.SourceInstanceId.IsUnknown() {
		if srcID, ok := sourceInstanceIDFromConfig(inst); ok {
			state.SourceInstanceId = types.Int64Value(srcID)
		}
	}

	// Update computed fields
	state.Status = convert.StrToType(inst.Status)
	state.Name = convert.StrToType(inst.Name)

	// Refresh volumes and network interfaces from the API for drift detection,
	// preserving the write-only inputs (size_id, mac_address) that the read API
	// does not return.
	volumes, vdiags := mergeVolumesForRead(ctx, state.Volumes, inst.Volumes)
	resp.Diagnostics.Append(vdiags...)
	state.Volumes = volumes

	interfaces, idiags := mergeInterfacesForRead(ctx, state.NetworkInterfaces, inst.Interfaces)
	resp.Diagnostics.Append(idiags...)
	state.NetworkInterfaces = interfaces

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func refreshStateFromAPI(
	ctx context.Context,
	client *sdk.APIClient,
	instanceID int64,
	state *InstanceCloneModel,
) diag.Diagnostics {
	var diags diag.Diagnostics

	getResp, hresp, err := client.InstancesAPI.GetInstance(ctx, instanceID).Execute()
	if err != nil {
		diags.AddError(
			"read instance clone after create",
			fmt.Sprintf("instance %d GET failed: %s", instanceID, errfmt.ErrMsg(err, hresp)),
		)

		return diags
	}

	if getResp == nil || getResp.Instance == nil {
		diags.AddError(
			"read instance clone after create",
			fmt.Sprintf("instance %d returned nil on read", instanceID),
		)

		return diags
	}

	inst := getResp.Instance

	volumes, vdiags := mergeVolumesFromAPI(ctx, state.Volumes, inst.Volumes)
	diags.Append(vdiags...)
	state.Volumes = volumes

	interfaces, idiags := mergeInterfacesFromAPI(ctx, state.NetworkInterfaces, inst.Interfaces)
	diags.Append(idiags...)
	state.NetworkInterfaces = interfaces

	// config / config_* are clone-time overrides carried from the plan; they are
	// merged into the source instance's configuration server-side and are not
	// reconciled back from the API. The dynamic config attribute is
	// computed_optional, so resolve it to null when no generic config block was
	// supplied to avoid leaving an unknown value in state after apply.
	if state.Config.IsUnknown() {
		state.Config = types.DynamicNull()
	}

	return diags
}

// volumeValueFromAPI maps a single API volume into a VolumesValue.
func volumeValueFromAPI(
	v sdk.AddInstance200ResponseAllOfOneOfInstanceVolumesInner,
) VolumesValue {
	vol := VolumesValue{
		state: attr.ValueStateKnown,
	}

	vol.Id = convert.Int64ToType(v.Id)
	vol.Name = convert.StrToType(v.Name)
	vol.Size = convert.Int64ToType(v.Size)
	vol.RootVolume = convert.BoolToType(v.RootVolume)
	vol.StorageType = convert.Int64ToType(v.StorageType)

	// DatastoreId in this response is *string; parse to int64 if present.
	vol.DatastoreId = types.Int64Null()
	if v.DatastoreId != nil {
		if dsInt, parseErr := strconv.ParseInt(*v.DatastoreId, 10, 64); parseErr == nil {
			vol.DatastoreId = types.Int64Value(dsInt)
		}
	}

	vol.ControllerMountPoint = convert.StrToType(v.ControllerMountPoint)
	// size_id is a write-only input and is not returned by the read API.
	vol.SizeId = types.Int64Null()

	return vol
}

// mergeVolumesFromAPI builds the post-apply volume state: it preserves the
// user-configured Required fields (name, size) from the plan and fills the
// computed / computed-optional fields from the API, matched positionally. This
// keeps the final state structurally identical to the plan - avoiding
// "Provider produced inconsistent result after apply" when the platform reports
// a value that differs from the request (e.g. a silently-ignored shrink) - while
// still resolving the computed fields.
func mergeVolumesFromAPI(
	ctx context.Context,
	planVolumes types.List,
	apiVolumes []sdk.AddInstance200ResponseAllOfOneOfInstanceVolumesInner,
) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics

	if planVolumes.IsNull() || planVolumes.IsUnknown() {
		objs, d := mapVolumeObjs(ctx, apiVolumes)
		diags.Append(d...)
		listVal, ld := types.ListValue(VolumesValue{}.Type(ctx), objs)
		diags.Append(ld...)

		return listVal, diags
	}

	var pv []VolumesValue
	diags.Append(planVolumes.ElementsAs(ctx, &pv, false)...)

	values := make([]attr.Value, 0, len(pv))
	for i := range pv {
		api := VolumesValue{
			state:                attr.ValueStateKnown,
			Id:                   types.Int64Null(),
			RootVolume:           types.BoolNull(),
			DatastoreId:          types.Int64Null(),
			StorageType:          types.Int64Null(),
			ControllerMountPoint: types.StringNull(),
			SizeId:               types.Int64Null(),
		}
		if i < len(apiVolumes) {
			api = volumeValueFromAPI(apiVolumes[i])
		}

		vol := VolumesValue{
			state: attr.ValueStateKnown,
			// Required fields come from the plan so the final state matches it.
			Name: pv[i].Name,
			Size: pv[i].Size,
			// size_id is a write-only input; keep the plan value.
			SizeId: pv[i].SizeId,
			// Computed id always comes from the API.
			Id: api.Id,
			// Computed-optional fields keep the user's value when set.
			RootVolume:           preferKnownBool(pv[i].RootVolume, api.RootVolume),
			DatastoreId:          preferKnownInt64(pv[i].DatastoreId, api.DatastoreId),
			StorageType:          preferKnownInt64(pv[i].StorageType, api.StorageType),
			ControllerMountPoint: preferKnownString(pv[i].ControllerMountPoint, api.ControllerMountPoint),
		}

		objVal, d := vol.ToObjectValue(ctx)
		diags.Append(d...)
		values = append(values, objVal)
	}

	listVal, d := types.ListValue(VolumesValue{}.Type(ctx), values)
	diags.Append(d...)

	return listVal, diags
}

func mapVolumeObjs(
	ctx context.Context,
	apiVolumes []sdk.AddInstance200ResponseAllOfOneOfInstanceVolumesInner,
) ([]attr.Value, diag.Diagnostics) {
	var diags diag.Diagnostics

	values := make([]attr.Value, 0, len(apiVolumes))
	for _, v := range apiVolumes {
		objVal, d := volumeValueFromAPI(v).ToObjectValue(ctx)
		diags.Append(d...)
		values = append(values, objVal)
	}

	return values, diags
}

// mergeVolumesForRead builds Read state: it refreshes the API-backed fields (for
// drift detection) while preserving the write-only input (size_id) that the read
// API does not return, taken from the prior state positionally.
func mergeVolumesForRead(
	ctx context.Context,
	priorVolumes types.List,
	apiVolumes []sdk.AddInstance200ResponseAllOfOneOfInstanceVolumesInner,
) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics

	if len(apiVolumes) == 0 {
		return types.ListNull(VolumesValue{}.Type(ctx)), diags
	}

	var prior []VolumesValue
	if !priorVolumes.IsNull() && !priorVolumes.IsUnknown() {
		diags.Append(priorVolumes.ElementsAs(ctx, &prior, false)...)
	}

	values := make([]attr.Value, 0, len(apiVolumes))
	for i, v := range apiVolumes {
		vol := volumeValueFromAPI(v)
		if i < len(prior) {
			vol.SizeId = prior[i].SizeId
		}

		objVal, d := vol.ToObjectValue(ctx)
		diags.Append(d...)
		values = append(values, objVal)
	}

	listVal, d := types.ListValue(VolumesValue{}.Type(ctx), values)
	diags.Append(d...)

	return listVal, diags
}

// mergeInterfacesForRead builds Read state: it refreshes the API-backed fields
// while preserving the write-only input (mac_address) from prior state.
func mergeInterfacesForRead(
	ctx context.Context,
	priorInterfaces types.List,
	apiInterfaces []sdk.AddInstance200ResponseAllOfOneOfInstanceInterfacesInner,
) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics

	if len(apiInterfaces) == 0 {
		return types.ListNull(NetworkInterfacesValue{}.Type(ctx)), diags
	}

	var prior []NetworkInterfacesValue
	if !priorInterfaces.IsNull() && !priorInterfaces.IsUnknown() {
		diags.Append(priorInterfaces.ElementsAs(ctx, &prior, false)...)
	}

	values := make([]attr.Value, 0, len(apiInterfaces))
	for i, iface := range apiInterfaces {
		ni, d := interfaceValueFromAPI(ctx, iface)
		diags.Append(d...)

		var priorChildren types.List
		if i < len(prior) {
			ni.MacAddress = prior[i].MacAddress
			priorChildren = prior[i].ChildVirtualNetworks
		}
		children, cd := mergeChildInterfacesForRead(ctx, priorChildren, iface.NetworkInterfaces)
		diags.Append(cd...)
		ni.ChildVirtualNetworks = children

		objVal, od := ni.ToObjectValue(ctx)
		diags.Append(od...)
		values = append(values, objVal)
	}

	listVal, d := types.ListValue(NetworkInterfacesValue{}.Type(ctx), values)
	diags.Append(d...)

	return listVal, diags
}

// mergeChildInterfacesForRead refreshes child interface API fields for Read while
// preserving the write-only input (mac_address) from prior state, matched
// positionally.
func mergeChildInterfacesForRead(
	ctx context.Context,
	priorChildren types.List,
	apiChildren []sdk.InstanceInterfacesNetworkInterfacesInner1,
) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics

	if len(apiChildren) == 0 {
		return types.ListNull(ChildVirtualNetworksValue{}.Type(ctx)), diags
	}

	var prior []ChildVirtualNetworksValue
	if !priorChildren.IsNull() && !priorChildren.IsUnknown() {
		diags.Append(priorChildren.ElementsAs(ctx, &prior, false)...)
	}

	values := make([]attr.Value, 0, len(apiChildren))
	for i, child := range apiChildren {
		c := childInterfaceValueFromAPI(child)
		if i < len(prior) {
			c.MacAddress = prior[i].MacAddress
		}

		objVal, d := c.ToObjectValue(ctx)
		diags.Append(d...)
		values = append(values, objVal)
	}

	listVal, d := types.ListValue(ChildVirtualNetworksValue{}.Type(ctx), values)
	diags.Append(d...)

	return listVal, diags
}

// preferKnownInt64 / preferKnownBool / preferKnownString return the plan value
// when the user supplied one (known), otherwise the value resolved from the API.
// This keeps user-configured values stable in state while still resolving
// computed-optional fields the user left unset.
func preferKnownInt64(plan, api types.Int64) types.Int64 {
	if !plan.IsUnknown() {
		return plan
	}

	return api
}

func preferKnownBool(plan, api types.Bool) types.Bool {
	if !plan.IsUnknown() {
		return plan
	}

	return api
}

func preferKnownString(plan, api types.String) types.String {
	if !plan.IsUnknown() {
		return plan
	}

	return api
}

// interfaceValueFromAPI maps a single API interface into a NetworkInterfacesValue.
func interfaceValueFromAPI(
	ctx context.Context,
	iface sdk.AddInstance200ResponseAllOfOneOfInstanceInterfacesInner,
) (NetworkInterfacesValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	ni := NetworkInterfacesValue{
		state: attr.ValueStateKnown,
	}

	// Interface Id is an anyOf type (Int64 | String).
	ni.Id = types.Int64Null()
	if iface.Id != nil && iface.Id.Int64 != nil {
		ni.Id = types.Int64Value(*iface.Id.Int64)
	}

	// Network ID from the nested network object.
	ni.NetworkId = types.Int64Null()
	if iface.Network != nil && iface.Network.Id != nil {
		ni.NetworkId = types.Int64Value(*iface.Network.Id)
	}

	ni.NetworkInterfaceTypeId = convert.Int64ToType(iface.NetworkInterfaceTypeId)
	ni.IpMode = convert.StrToType(iface.IpMode)
	ni.IpAddress = convert.StrToType(iface.IpAddress)
	// mac_address is a write-only input and is not returned by the read API.
	ni.MacAddress = types.StringNull()

	children, d := childInterfacesFromAPI(ctx, iface.NetworkInterfaces)
	diags.Append(d...)
	ni.ChildVirtualNetworks = children

	return ni, diags
}

// childInterfaceValueFromAPI maps a single API child interface into a
// ChildVirtualNetworksValue.
func childInterfaceValueFromAPI(
	child sdk.InstanceInterfacesNetworkInterfacesInner1,
) ChildVirtualNetworksValue {
	c := ChildVirtualNetworksValue{
		state: attr.ValueStateKnown,
	}

	// Interface Id is an anyOf type (Int64 | String).
	c.Id = types.Int64Null()
	if child.Id != nil && child.Id.Int64 != nil {
		c.Id = types.Int64Value(*child.Id.Int64)
	}

	// Network ID from the nested network object.
	c.NetworkId = types.Int64Null()
	if child.Network != nil && child.Network.Id != nil {
		c.NetworkId = types.Int64Value(*child.Network.Id)
	}

	c.NetworkInterfaceTypeId = convert.Int64ToType(child.NetworkInterfaceTypeId)
	c.IpMode = convert.StrToType(child.IpMode)
	c.IpAddress = convert.StrToType(child.IpAddress)
	// mac_address is a write-only input and is not returned by the read API.
	c.MacAddress = types.StringNull()

	return c
}

// childInterfacesFromAPI maps the API child interfaces of a single parent into a
// typed list.
func childInterfacesFromAPI(
	ctx context.Context,
	apiChildren []sdk.InstanceInterfacesNetworkInterfacesInner1,
) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics

	if len(apiChildren) == 0 {
		return types.ListNull(ChildVirtualNetworksValue{}.Type(ctx)), diags
	}

	values := make([]attr.Value, 0, len(apiChildren))
	for _, child := range apiChildren {
		objVal, d := childInterfaceValueFromAPI(child).ToObjectValue(ctx)
		diags.Append(d...)
		values = append(values, objVal)
	}

	listVal, d := types.ListValue(ChildVirtualNetworksValue{}.Type(ctx), values)
	diags.Append(d...)

	return listVal, diags
}

// mergeInterfacesFromAPI builds the post-apply interface state: it preserves the
// user-configured Required field (network_id) from the plan and fills the
// computed / computed-optional fields from the API, matched positionally.
func mergeInterfacesFromAPI(
	ctx context.Context,
	planInterfaces types.List,
	apiInterfaces []sdk.AddInstance200ResponseAllOfOneOfInstanceInterfacesInner,
) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics

	if planInterfaces.IsNull() || planInterfaces.IsUnknown() {
		values := make([]attr.Value, 0, len(apiInterfaces))
		for _, iface := range apiInterfaces {
			ni, d := interfaceValueFromAPI(ctx, iface)
			diags.Append(d...)
			objVal, od := ni.ToObjectValue(ctx)
			diags.Append(od...)
			values = append(values, objVal)
		}
		listVal, d := types.ListValue(NetworkInterfacesValue{}.Type(ctx), values)
		diags.Append(d...)

		return listVal, diags
	}

	var pi []NetworkInterfacesValue
	diags.Append(planInterfaces.ElementsAs(ctx, &pi, false)...)

	values := make([]attr.Value, 0, len(pi))
	for i := range pi {
		api := NetworkInterfacesValue{
			state:                  attr.ValueStateKnown,
			Id:                     types.Int64Null(),
			IpMode:                 types.StringNull(),
			IpAddress:              types.StringNull(),
			MacAddress:             types.StringNull(),
			NetworkInterfaceTypeId: types.Int64Null(),
			ChildVirtualNetworks:   types.ListNull(ChildVirtualNetworksValue{}.Type(ctx)),
		}

		var apiChildren []sdk.InstanceInterfacesNetworkInterfacesInner1
		if i < len(apiInterfaces) {
			apiNi, d := interfaceValueFromAPI(ctx, apiInterfaces[i])
			diags.Append(d...)
			api = apiNi
			apiChildren = apiInterfaces[i].NetworkInterfaces
		}

		children, cd := mergeChildInterfacesFromAPI(ctx, pi[i].ChildVirtualNetworks, apiChildren)
		diags.Append(cd...)

		ni := NetworkInterfacesValue{
			state: attr.ValueStateKnown,
			// Required field comes from the plan so the final state matches it.
			NetworkId: pi[i].NetworkId,
			// mac_address is a write-only input; keep the plan value.
			MacAddress: pi[i].MacAddress,
			// Computed id always comes from the API.
			Id: api.Id,
			// Computed-optional fields keep the user's value when set.
			IpMode:                 preferKnownString(pi[i].IpMode, api.IpMode),
			IpAddress:              preferKnownString(pi[i].IpAddress, api.IpAddress),
			NetworkInterfaceTypeId: preferKnownInt64(pi[i].NetworkInterfaceTypeId, api.NetworkInterfaceTypeId),
			ChildVirtualNetworks:   children,
		}

		objVal, od := ni.ToObjectValue(ctx)
		diags.Append(od...)
		values = append(values, objVal)
	}

	listVal, d := types.ListValue(NetworkInterfacesValue{}.Type(ctx), values)
	diags.Append(d...)

	return listVal, diags
}

// mergeChildInterfacesFromAPI builds the post-apply child interface state: it
// preserves the user-configured Required field (network_id) and the write-only
// input (mac_address) from the plan, filling the computed / computed-optional
// fields from the API, matched positionally. Iterating over the plan children
// keeps the final state structurally identical to the plan.
func mergeChildInterfacesFromAPI(
	ctx context.Context,
	planChildren types.List,
	apiChildren []sdk.InstanceInterfacesNetworkInterfacesInner1,
) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics

	if planChildren.IsNull() || planChildren.IsUnknown() {
		return childInterfacesFromAPI(ctx, apiChildren)
	}

	var pc []ChildVirtualNetworksValue
	diags.Append(planChildren.ElementsAs(ctx, &pc, false)...)

	values := make([]attr.Value, 0, len(pc))
	for i := range pc {
		api := ChildVirtualNetworksValue{
			state:                  attr.ValueStateKnown,
			Id:                     types.Int64Null(),
			IpMode:                 types.StringNull(),
			IpAddress:              types.StringNull(),
			NetworkInterfaceTypeId: types.Int64Null(),
		}
		if i < len(apiChildren) {
			api = childInterfaceValueFromAPI(apiChildren[i])
		}

		c := ChildVirtualNetworksValue{
			state: attr.ValueStateKnown,
			// Required field comes from the plan so the final state matches it.
			NetworkId: pc[i].NetworkId,
			// mac_address is a write-only input; keep the plan value.
			MacAddress: pc[i].MacAddress,
			// Computed id always comes from the API.
			Id: api.Id,
			// Computed-optional fields keep the user's value when set.
			IpMode:                 preferKnownString(pc[i].IpMode, api.IpMode),
			IpAddress:              preferKnownString(pc[i].IpAddress, api.IpAddress),
			NetworkInterfaceTypeId: preferKnownInt64(pc[i].NetworkInterfaceTypeId, api.NetworkInterfaceTypeId),
		}

		objVal, d := c.ToObjectValue(ctx)
		diags.Append(d...)
		values = append(values, objVal)
	}

	listVal, d := types.ListValue(ChildVirtualNetworksValue{}.Type(ctx), values)
	diags.Append(d...)

	return listVal, diags
}
