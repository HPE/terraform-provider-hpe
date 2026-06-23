// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package instanceclone

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"
	"github.com/cenkalti/backoff/v5"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
)

func (r *Resource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	var plan, state InstanceCloneModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateTimeout, diags := plan.Timeouts.Update(ctx, 15*time.Minute)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, updateTimeout)
	defer cancel()

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"update instance clone",
			"failed to create client: "+err.Error(),
		)

		return
	}

	instanceID := state.Id.ValueInt64()

	// Build resize request
	resizeReq := sdk.ResizeInstanceRequest{}

	if !plan.PlanId.IsNull() && !plan.PlanId.IsUnknown() {
		pid := plan.PlanId.ValueInt64()
		resizeReq.Instance = &sdk.ResizeInstanceRequestInstance{
			Plan: &sdk.ResizeInstanceRequestInstancePlan{Id: &pid},
		}
	}

	resizeVolumes, volDiags := buildResizeVolumes(ctx, plan.Volumes)
	resp.Diagnostics.Append(volDiags...)
	resizeIfaces, ifaceDiags := buildResizeNetworkInterfaces(ctx, plan.NetworkInterfaces)
	resp.Diagnostics.Append(ifaceDiags...)
	if resp.Diagnostics.HasError() {
		return
	}
	resizeReq.Volumes = resizeVolumes
	resizeReq.NetworkInterfaces = resizeIfaces

	_, hresp, err := client.InstancesAPI.ResizeInstance(ctx, instanceID).
		ResizeInstanceRequest(resizeReq).
		Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"update instance clone",
			fmt.Sprintf("resize request failed for instance %d: %s",
				instanceID, errfmt.ErrMsg(err, hresp)),
		)

		return
	}

	// Wait for stable status after resize
	waitForStatus := func() (*string, error) {
		getResp, hresp, err := client.InstancesAPI.GetInstance(ctx, instanceID).Execute()
		if err != nil {
			return nil, pollAPIError(
				fmt.Sprintf("failed to get instance %d", instanceID), err, hresp,
			)
		}

		if getResp == nil || getResp.Instance == nil {
			return nil, fmt.Errorf("get instance %d returned nil", instanceID)
		}

		status := ""
		if getResp.Instance.Status != nil {
			status = *getResp.Instance.Status
		}

		if err := checkStatusDone(status, cloneTargetStatuses, cloneErrorStatuses); err != nil {
			return nil, err
		}

		return &status, nil
	}

	statusResult, err := backoff.Retry(
		ctx,
		waitForStatus,
		backoff.WithBackOff(backoff.NewConstantBackOff(10*time.Second)),
		backoff.WithMaxElapsedTime(updateTimeout),
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"update instance clone",
			fmt.Sprintf("instance %d failed or timed out after resize: %v",
				instanceID, errors.Unwrap(err)),
		)

		return
	}

	// Carry over immutable fields from prior state
	plan.Id = state.Id
	plan.Status = types.StringValue(*statusResult)

	// Refresh actual volumes/interfaces from API
	refreshDiags := refreshStateFromAPI(ctx, client, instanceID, &plan)
	resp.Diagnostics.Append(refreshDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// cloneNeedsResize reports whether any requested volume size exceeds the actual
// provisioned volume size (matched positionally), meaning a post-clone resize is
// required. Used to grow data disks that some clouds leave at the source size.
func cloneNeedsResize(
	ctx context.Context,
	planVolumes types.List,
	actual []sdk.AddInstance200ResponseAllOfOneOfInstanceVolumesInner,
) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	if planVolumes.IsNull() || planVolumes.IsUnknown() {
		return false, diags
	}

	var pv []VolumesValue
	diags.Append(planVolumes.ElementsAs(ctx, &pv, false)...)

	for i, v := range pv {
		if v.Size.IsNull() || v.Size.IsUnknown() || i >= len(actual) {
			continue
		}

		var actualSize int64
		if actual[i].Size != nil {
			actualSize = *actual[i].Size
		}

		if v.Size.ValueInt64() > actualSize {
			return true, diags
		}
	}

	return false, diags
}

// buildResizeVolumes converts plan volumes for resize.
func buildResizeVolumes(
	ctx context.Context, volumesList types.List,
) ([]sdk.ResizeInstanceRequestVolumesInner, diag.Diagnostics) {
	var diags diag.Diagnostics

	if volumesList.IsNull() || volumesList.IsUnknown() {
		return nil, diags
	}

	var planVolumes []VolumesValue
	diags.Append(volumesList.ElementsAs(ctx, &planVolumes, false)...)

	sdkVolumes := make([]sdk.ResizeInstanceRequestVolumesInner, 0, len(planVolumes))
	for _, v := range planVolumes {
		vol := sdk.ResizeInstanceRequestVolumesInner{
			Name: v.Name.ValueStringPointer(),
			Size: v.Size.ValueInt64Pointer(),
		}

		if !v.RootVolume.IsNull() && !v.RootVolume.IsUnknown() {
			vol.RootVolume = v.RootVolume.ValueBoolPointer()
		}

		if !v.StorageType.IsNull() && !v.StorageType.IsUnknown() {
			stVal := v.StorageType.ValueInt64()
			vol.StorageType = *sdk.NewNullableInt64(&stVal)
		}

		if !v.DatastoreId.IsNull() && !v.DatastoreId.IsUnknown() {
			dsID := v.DatastoreId.ValueInt64()
			vol.DatastoreId = &sdk.CloneInstanceRequestVolumesInnerDatastoreId{
				Int64: &dsID,
			}
		}

		if !v.SizeId.IsNull() && !v.SizeId.IsUnknown() {
			siVal := v.SizeId.ValueInt64()
			vol.SizeId = *sdk.NewNullableInt64(&siVal)
		}

		if !v.ControllerMountPoint.IsNull() && !v.ControllerMountPoint.IsUnknown() {
			vol.ControllerMountPoint = v.ControllerMountPoint.ValueStringPointer()
		}

		if !v.Id.IsNull() && !v.Id.IsUnknown() {
			vol.Id = v.Id.ValueInt64Pointer()
		}

		sdkVolumes = append(sdkVolumes, vol)
	}

	return sdkVolumes, diags
}

// buildResizeNetworkInterfaces converts plan network interfaces for resize.
func buildResizeNetworkInterfaces(
	ctx context.Context, ifaceList types.List,
) ([]sdk.InstancesNetworkInterfaces4, diag.Diagnostics) {
	var diags diag.Diagnostics

	if ifaceList.IsNull() || ifaceList.IsUnknown() {
		return nil, diags
	}

	var planIfaces []NetworkInterfacesValue
	diags.Append(ifaceList.ElementsAs(ctx, &planIfaces, false)...)

	sdkIfaces := make([]sdk.InstancesNetworkInterfaces4, 0, len(planIfaces))
	for _, ni := range planIfaces {
		iface := sdk.InstancesNetworkInterfaces4{
			Network: sdk.InstancesNetworkInterfaces3Network{
				Id: strconv.FormatInt(ni.NetworkId.ValueInt64(), 10),
			},
		}

		if !ni.NetworkInterfaceTypeId.IsNull() && !ni.NetworkInterfaceTypeId.IsUnknown() {
			iface.NetworkInterfaceTypeId = ni.NetworkInterfaceTypeId.ValueInt64Pointer()
		}

		if !ni.IpMode.IsNull() && !ni.IpMode.IsUnknown() {
			iface.IpMode = ni.IpMode.ValueStringPointer()
		}

		if !ni.IpAddress.IsNull() && !ni.IpAddress.IsUnknown() {
			iface.IpAddress = ni.IpAddress.ValueStringPointer()
		}

		if !ni.MacAddress.IsNull() && !ni.MacAddress.IsUnknown() {
			iface.MacAddress = ni.MacAddress.ValueStringPointer()
		}

		if !ni.Id.IsNull() && !ni.Id.IsUnknown() {
			iface.Id = ni.Id.ValueInt64Pointer()
		}

		children, cd := buildChildInterfaces(ctx, ni.ChildVirtualNetworks)
		diags.Append(cd...)
		iface.NetworkInterfaces = children

		sdkIfaces = append(sdkIfaces, iface)
	}

	return sdkIfaces, diags
}
