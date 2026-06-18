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
	"github.com/HPE/terraform-provider-hpe/utils/cleanup"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

func (r *Resource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	var plan InstanceCloneModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createTimeout, diags := plan.Timeouts.Create(ctx, 15*time.Minute)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, createTimeout)
	defer cancel()

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"create instance clone",
			"failed to create client: "+err.Error(),
		)

		return
	}

	sourceID := plan.SourceInstanceId.ValueInt64()
	cloneName := plan.Name.ValueString()

	// Build CloneInstanceRequest
	cloneReq := sdk.CloneInstanceRequest{
		Name: &cloneName,
	}

	if !plan.GroupId.IsNull() && !plan.GroupId.IsUnknown() {
		gid := plan.GroupId.ValueInt64()
		cloneReq.Group = &sdk.CloneInstanceRequestGroup{Id: &gid}
	}

	if !plan.CloudId.IsNull() && !plan.CloudId.IsUnknown() {
		cid := plan.CloudId.ValueInt64()
		cloneReq.Cloud = &sdk.CloneInstanceRequestCloud{Id: &cid}
	}

	if !plan.PlanId.IsNull() && !plan.PlanId.IsUnknown() {
		pid := plan.PlanId.ValueInt64()
		cloneReq.Plan = &sdk.CloneInstanceRequestPlan{Id: &pid}
	}

	cloneConfig, cfgDiags := buildCloneConfig(ctx, plan)
	resp.Diagnostics.Append(cfgDiags...)
	if resp.Diagnostics.HasError() {
		return
	}
	cloneReq.Config = cloneConfig

	// Build volumes
	cloneVolumes, volDiags := buildCloneVolumes(ctx, plan.Volumes)
	resp.Diagnostics.Append(volDiags...)
	cloneReq.Volumes = cloneVolumes

	// Build network interfaces
	cloneIfaces, ifaceDiags := buildCloneNetworkInterfaces(ctx, plan.NetworkInterfaces)
	resp.Diagnostics.Append(ifaceDiags...)
	if resp.Diagnostics.HasError() {
		return
	}
	cloneReq.NetworkInterfaces = cloneIfaces

	// Fire clone request
	cloneResp, hresp, err := client.InstancesAPI.CloneInstance(ctx, sourceID).
		CloneInstanceRequest(cloneReq).
		Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"create instance clone",
			fmt.Sprintf("clone request failed for instance %d: %s",
				sourceID, errfmt.ErrMsg(err, hresp)),
		)

		return
	}

	if cloneResp != nil && cloneResp.Success != nil && !*cloneResp.Success {
		resp.Diagnostics.AddError(
			"create instance clone",
			fmt.Sprintf("clone request returned success=false for instance %d", sourceID),
		)

		return
	}

	// Poll for clone instance to appear by name
	type pollResult struct {
		id     int64
		status string
	}

	findClone := func() (*pollResult, error) {
		listResp, hresp, err := client.InstancesAPI.ListInstances(ctx).
			Name(cloneName).Execute()
		if err != nil {
			return nil, pollAPIError("failed to list instances", err, hresp)
		}

		if listResp == nil {
			return nil, fmt.Errorf("list instances returned nil response")
		}

		for _, inst := range listResp.Instances {
			if inst.Id == nil || inst.Name == nil {
				continue
			}
			if *inst.Name != cloneName {
				continue
			}

			status := ""
			if inst.Status != nil {
				status = *inst.Status
			}

			return &pollResult{id: *inst.Id, status: status}, nil
		}

		return nil, fmt.Errorf("clone %q not yet found in instance list", cloneName)
	}

	found, err := backoff.Retry(
		ctx,
		findClone,
		backoff.WithBackOff(backoff.NewConstantBackOff(5*time.Second)),
		backoff.WithMaxElapsedTime(createTimeout),
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"create instance clone",
			fmt.Sprintf("clone %q of instance %d failed to appear: %v",
				cloneName, sourceID, errors.Unwrap(err)),
		)

		return
	}

	cloneID := found.id

	// Wait for clone to reach a stable status
	waitForStatus := func() (*pollResult, error) {
		getResp, hresp, err := client.InstancesAPI.GetInstance(ctx, cloneID).Execute()
		if err != nil {
			return nil, pollAPIError(
				fmt.Sprintf("failed to get clone instance %d", cloneID), err, hresp,
			)
		}

		if getResp == nil || getResp.Instance == nil {
			return nil, fmt.Errorf("get instance %d returned nil", cloneID)
		}

		status := ""
		if getResp.Instance.Status != nil {
			status = *getResp.Instance.Status
		}

		if err := checkStatusDone(status, cloneTargetStatuses, cloneErrorStatuses); err != nil {
			return nil, err
		}

		return &pollResult{id: cloneID, status: status}, nil
	}

	result, err := backoff.Retry(
		ctx,
		waitForStatus,
		backoff.WithBackOff(backoff.NewConstantBackOff(10*time.Second)),
		backoff.WithMaxElapsedTime(createTimeout),
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"create instance clone",
			fmt.Sprintf("clone %q (id=%d) failed or timed out waiting for stable status: %v",
				cloneName, cloneID, errors.Unwrap(err)),
		)

		cleanup.TaintResourceState(ctx, cleanup.TaintResourceStateConfig{
			ResourceType: "instance_clone",
			ResourceID:   cloneID,
			StateWriter:  &resp.State,
			Diagnostics:  &resp.Diagnostics,
		})

		return
	}

	// Post-clone resize: some clouds (notably VMware) only resize the root
	// disk during clone, leaving data disks at the source size. If any
	// requested volume is larger than the actual provisioned volume, issue a
	// resize so the clone matches the configuration.
	if getResp, _, gErr := client.InstancesAPI.GetInstance(ctx, cloneID).Execute(); gErr == nil &&
		getResp != nil && getResp.Instance != nil &&
		cloneNeedsResize(ctx, plan.Volumes, getResp.Instance.Volumes) {
		resizeVolumes, rvDiags := buildResizeVolumes(ctx, plan.Volumes)
		resp.Diagnostics.Append(rvDiags...)
		resizeIfaces, riDiags := buildResizeNetworkInterfaces(ctx, plan.NetworkInterfaces)
		resp.Diagnostics.Append(riDiags...)
		if resp.Diagnostics.HasError() {
			return
		}
		resizeReq := sdk.ResizeInstanceRequest{
			Volumes:           resizeVolumes,
			NetworkInterfaces: resizeIfaces,
		}

		if _, hresp, rErr := client.InstancesAPI.ResizeInstance(ctx, cloneID).
			ResizeInstanceRequest(resizeReq).Execute(); rErr != nil {
			resp.Diagnostics.AddError(
				"create instance clone",
				fmt.Sprintf("post-clone resize failed for instance %d: %s",
					cloneID, errfmt.ErrMsg(rErr, hresp)),
			)

			cleanup.TaintResourceState(ctx, cleanup.TaintResourceStateConfig{
				ResourceType: "instance_clone",
				ResourceID:   cloneID,
				StateWriter:  &resp.State,
				Diagnostics:  &resp.Diagnostics,
			})

			return
		}

		if result, err = backoff.Retry(
			ctx,
			waitForStatus,
			backoff.WithBackOff(backoff.NewConstantBackOff(10*time.Second)),
			backoff.WithMaxElapsedTime(createTimeout),
		); err != nil {
			resp.Diagnostics.AddError(
				"create instance clone",
				fmt.Sprintf("clone %q (id=%d) failed or timed out after post-clone resize: %v",
					cloneName, cloneID, errors.Unwrap(err)),
			)

			cleanup.TaintResourceState(ctx, cleanup.TaintResourceStateConfig{
				ResourceType: "instance_clone",
				ResourceID:   cloneID,
				StateWriter:  &resp.State,
				Diagnostics:  &resp.Diagnostics,
			})

			return
		}
	}

	// Read final state
	plan.Id = types.Int64Value(cloneID)
	plan.Status = types.StringValue(result.status)

	// Refresh volumes and network interfaces from actual state
	refreshDiags := refreshStateFromAPI(ctx, client, cloneID, &plan)
	resp.Diagnostics.Append(refreshDiags...)
	if resp.Diagnostics.HasError() {
		cleanup.TaintResourceState(ctx, cleanup.TaintResourceStateConfig{
			ResourceType: "instance_clone",
			ResourceID:   cloneID,
			StateWriter:  &resp.State,
			Diagnostics:  &resp.Diagnostics,
		})

		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// buildCloneConfig converts the plan's config / config_* attributes into the
// clone request config. The config is a clone-time override that the platform
// merges key-by-key over the source instance's configuration; it never mutates
// the source instance. Exactly one variant is populated, in the same precedence
// order as the instance resource. The variants are mutually exclusive at plan
// time (the dynamic config conflicts with the typed blocks).
func buildCloneConfig(
	ctx context.Context, plan InstanceCloneModel,
) (*sdk.CloneInstanceRequestConfig, diag.Diagnostics) {
	var diags diag.Diagnostics

	switch {
	// AWS config
	case !plan.ConfigAws.IsNull() && !plan.ConfigAws.IsUnknown():
		noAgent := plan.ConfigAws.NoAgent.ValueBool()
		isEC2 := convert.BoolToStringTrueFalse(plan.ConfigAws.IsEc2.ValueBool()).ValueString()
		publicIpType := plan.ConfigAws.PublicIpType.ValueString()
		resourcePoolId := plan.ConfigAws.ResourcePoolId.ValueString()

		configAWS := &sdk.AmazonInstanceConfiguration3{
			NoAgent:         *sdk.NewNullableBool(&noAgent),
			ResourcePoolId:  &resourcePoolId,
			IsEC2:           &isEC2,
			KmsKeyId:        plan.ConfigAws.KmsKeyId.ValueStringPointer(),
			InstanceProfile: plan.ConfigAws.InstanceProfile.ValueStringPointer(),
			PublicIpType:    &publicIpType,
			AvailabilityId:  plan.ConfigAws.AvailabilityZoneId.ValueStringPointer(),
		}

		// The clone request has no top-level security groups field (unlike
		// instance creation), so AWS security groups are carried inside the
		// config and merged over the source instance's configuration.
		if !plan.ConfigAws.SecurityGroups.IsNull() && !plan.ConfigAws.SecurityGroups.IsUnknown() {
			var sgs []SecurityGroupsValue
			diags.Append(plan.ConfigAws.SecurityGroups.ElementsAs(ctx, &sgs, false)...)
			if diags.HasError() {
				return nil, diags
			}
			if len(sgs) > 0 {
				sgList := make([]map[string]interface{}, 0, len(sgs))
				for _, sg := range sgs {
					sgList = append(sgList, map[string]interface{}{"id": sg.Id.ValueString()})
				}
				configAWS.AdditionalProperties = map[string]interface{}{"securityGroups": sgList}
			}
		}

		return &sdk.CloneInstanceRequestConfig{AmazonInstanceConfiguration3: configAWS}, diags

	// HVM config
	case !plan.ConfigHvm.IsNull() && !plan.ConfigHvm.IsUnknown():
		createUser := plan.ConfigHvm.CreateUser.ValueBool()
		nestedVirtualization := plan.ConfigHvm.NestedVirtualization.ValueString()
		noAgent := plan.ConfigHvm.NoAgent.ValueBool()
		resourcePoolId := plan.ConfigHvm.ResourcePoolId.ValueString()

		configHvm := &sdk.HVMInstanceConfiguration1{
			CreateUser:           *sdk.NewNullableBool(&createUser),
			NestedVirtualization: &nestedVirtualization,
			NoAgent:              *sdk.NewNullableBool(&noAgent),
			ResourcePoolId:       &resourcePoolId,
		}

		if !plan.ConfigHvm.KvmHostId.IsNull() {
			configHvm.KvmHostId = plan.ConfigHvm.KvmHostId.ValueInt64Pointer()
		}

		return &sdk.CloneInstanceRequestConfig{HVMInstanceConfiguration1: configHvm}, diags

	// VMware config
	case !plan.ConfigVmware.IsNull() && !plan.ConfigVmware.IsUnknown():
		nestedVirtualization := plan.ConfigVmware.NestedVirtualization.ValueString()
		createUser := plan.ConfigVmware.CreateUser.ValueBool()
		noAgent := plan.ConfigVmware.NoAgent.ValueBool()
		resourcePoolId := plan.ConfigVmware.ResourcePoolId.ValueString()

		configVMware := &sdk.VMWareInstanceConfiguration3{
			NestedVirtualization: &nestedVirtualization,
			CreateUser:           *sdk.NewNullableBool(&createUser),
			NoAgent:              *sdk.NewNullableBool(&noAgent),
			ResourcePoolId:       &resourcePoolId,
			VmwareFolderId:       plan.ConfigVmware.VmwareFolderId.ValueStringPointer(),
		}

		return &sdk.CloneInstanceRequestConfig{VMWareInstanceConfiguration3: configVMware}, diags

	// Azure config
	case !plan.ConfigAzure.IsNull() && !plan.ConfigAzure.IsUnknown():
		createUser := plan.ConfigAzure.CreateUser.ValueBool()
		resourcePoolId := plan.ConfigAzure.ResourcePoolId.ValueString()

		configAzure := &sdk.AzureInstanceConfiguration3{
			CreateUser:     &createUser,
			ResourcePoolId: &resourcePoolId,
		}

		if !plan.ConfigAzure.AzureRegion.IsNull() && !plan.ConfigAzure.AzureRegion.IsUnknown() {
			configAzure.AzureRegion = plan.ConfigAzure.AzureRegion.ValueStringPointer()
		}
		if !plan.ConfigAzure.AzuresecurityGroupId.IsNull() && !plan.ConfigAzure.AzuresecurityGroupId.IsUnknown() {
			configAzure.AzuresecurityGroupId = plan.ConfigAzure.AzuresecurityGroupId.ValueStringPointer()
		}
		if !plan.ConfigAzure.AvailabilityOptions.IsNull() && !plan.ConfigAzure.AvailabilityOptions.IsUnknown() {
			configAzure.AvailabilityOptions = plan.ConfigAzure.AvailabilityOptions.ValueStringPointer()
		}
		if !plan.ConfigAzure.AvailabilitySet.IsNull() && !plan.ConfigAzure.AvailabilitySet.IsUnknown() {
			configAzure.AvailabilitySet = plan.ConfigAzure.AvailabilitySet.ValueStringPointer()
		}
		if !plan.ConfigAzure.AvailabilityZone.IsNull() && !plan.ConfigAzure.AvailabilityZone.IsUnknown() {
			if configAzure.AdditionalProperties == nil {
				configAzure.AdditionalProperties = make(map[string]interface{})
			}
			configAzure.AdditionalProperties["availabilityZone"] = plan.ConfigAzure.AvailabilityZone.ValueString()
		}
		if !plan.ConfigAzure.AzurefloatingIp.IsNull() && !plan.ConfigAzure.AzurefloatingIp.IsUnknown() {
			configAzure.AzurefloatingIp = plan.ConfigAzure.AzurefloatingIp.ValueStringPointer()
		}
		if !plan.ConfigAzure.BootDiagnostics.IsNull() && !plan.ConfigAzure.BootDiagnostics.IsUnknown() {
			configAzure.BootDiagnostics = plan.ConfigAzure.BootDiagnostics.ValueStringPointer()
		}
		if !plan.ConfigAzure.OsGuestDiagnostics.IsNull() && !plan.ConfigAzure.OsGuestDiagnostics.IsUnknown() {
			configAzure.OsGuestDiagnostics = plan.ConfigAzure.OsGuestDiagnostics.ValueStringPointer()
		}
		if !plan.ConfigAzure.DiagnosticsStorageAccount.IsNull() && !plan.ConfigAzure.DiagnosticsStorageAccount.IsUnknown() {
			configAzure.DiagnosticsStorageAccount = plan.ConfigAzure.DiagnosticsStorageAccount.ValueStringPointer()
		}

		return &sdk.CloneInstanceRequestConfig{AzureInstanceConfiguration3: configAzure}, diags

	// Generic / dynamic config
	case !plan.Config.IsNull() && !plan.Config.IsUnknown():
		configAny, err := convert.ValueToAny(ctx, plan.Config.UnderlyingValue())
		if err != nil {
			diags.AddError(
				"clone instance resource",
				"instance_clone: failed to convert config: "+err.Error(),
			)

			return nil, diags
		}

		configMap, ok := configAny.(map[string]interface{})
		if !ok {
			diags.AddError(
				"clone instance resource",
				"instance_clone: could not parse config value",
			)

			return nil, diags
		}

		return &sdk.CloneInstanceRequestConfig{
			GenericInstanceConfiguration3: &sdk.GenericInstanceConfiguration3{
				AdditionalProperties: configMap,
			},
		}, diags
	}

	return nil, diags
}

// buildCloneVolumes converts plan volumes to SDK CloneInstanceRequestVolumesInner.
func buildCloneVolumes(
	ctx context.Context, volumesList types.List,
) ([]sdk.CloneInstanceRequestVolumesInner, diag.Diagnostics) {
	var diags diag.Diagnostics

	if volumesList.IsNull() || volumesList.IsUnknown() {
		return nil, diags
	}

	var planVolumes []VolumesValue
	diags.Append(volumesList.ElementsAs(ctx, &planVolumes, false)...)

	sdkVolumes := make([]sdk.CloneInstanceRequestVolumesInner, 0, len(planVolumes))
	for _, v := range planVolumes {
		vol := sdk.CloneInstanceRequestVolumesInner{
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

// buildCloneNetworkInterfaces converts plan to SDK InstancesNetworkInterfaces3.
func buildCloneNetworkInterfaces(
	ctx context.Context, ifaceList types.List,
) ([]sdk.InstancesNetworkInterfaces3, diag.Diagnostics) {
	var diags diag.Diagnostics

	if ifaceList.IsNull() || ifaceList.IsUnknown() {
		return nil, diags
	}

	var planIfaces []NetworkInterfacesValue
	diags.Append(ifaceList.ElementsAs(ctx, &planIfaces, false)...)

	sdkIfaces := make([]sdk.InstancesNetworkInterfaces3, 0, len(planIfaces))
	for _, ni := range planIfaces {
		iface := sdk.InstancesNetworkInterfaces3{
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

// buildChildInterfaces converts plan child virtual networks into the SDK child
// interface slice. The same element type is shared by clone
// (InstancesNetworkInterfaces3) and resize (InstancesNetworkInterfaces4).
func buildChildInterfaces(
	ctx context.Context, childList types.List,
) ([]sdk.InstancesNetworkInterfaces3NetworkInterfacesInner, diag.Diagnostics) {
	var diags diag.Diagnostics

	if childList.IsNull() || childList.IsUnknown() {
		return nil, diags
	}

	var planChildren []ChildVirtualNetworksValue
	diags.Append(childList.ElementsAs(ctx, &planChildren, false)...)

	children := make([]sdk.InstancesNetworkInterfaces3NetworkInterfacesInner, 0, len(planChildren))
	for _, c := range planChildren {
		child := sdk.InstancesNetworkInterfaces3NetworkInterfacesInner{
			Network: sdk.InstancesNetworkInterfaces3NetworkInterfacesInnerNetwork{
				Id: strconv.FormatInt(c.NetworkId.ValueInt64(), 10),
			},
		}

		if !c.NetworkInterfaceTypeId.IsNull() && !c.NetworkInterfaceTypeId.IsUnknown() {
			child.NetworkInterfaceTypeId = c.NetworkInterfaceTypeId.ValueInt64Pointer()
		}

		if !c.IpMode.IsNull() && !c.IpMode.IsUnknown() {
			child.IpMode = c.IpMode.ValueStringPointer()
		}

		if !c.IpAddress.IsNull() && !c.IpAddress.IsUnknown() {
			child.IpAddress = c.IpAddress.ValueStringPointer()
		}

		if !c.MacAddress.IsNull() && !c.MacAddress.IsUnknown() {
			child.MacAddress = c.MacAddress.ValueStringPointer()
		}

		if !c.Id.IsNull() && !c.Id.IsUnknown() {
			child.Id = c.Id.ValueInt64Pointer()
		}

		children = append(children, child)
	}

	return children, diags
}
