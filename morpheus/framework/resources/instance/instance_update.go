// (C) Copyright 2025-2026 Hewlett Packard Enterprise Development LP

package instance

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/cenkalti/backoff/v5"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"

	errfmt "github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/compare"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

var (
	UpdateTargetStatuses = []string{
		"running",
	}

	UpdateErrorStatuses = []string{
		"denied",
		"cancelled",
		"failed",
		"suspended",
		"removing",
		"pendingRemoval",
	}
)

// Update implements resource.Resource.
func (g *Resource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	client, err := g.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError("creating client failed", err.Error())

		return
	}

	var plan InstanceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state InstanceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get timeout from HCL if set, the default is 45 minutes
	updateTimeout, diags := plan.Timeouts.Update(ctx, 45*time.Minute)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, updateTimeout)
	defer cancel()

	var servicePlanOptions ServicePlanOptionsValue
	if isAPIUpdateNeeded(plan, state) {
		servicePlanOptions = makeUpdateAPIcalls(ctx, client, plan, state, updateTimeout, resp)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	tflog.Info(ctx, fmt.Sprintf("Instance update state: %v", state.Volumes.Elements()))
	tflog.Info(ctx, fmt.Sprintf("Instance update plan: %v", plan.Volumes.Elements()))

	newState, found, diag := getInstanceAsState(ctx, state.Id.ValueInt64(), client, plan, false)
	if resp.Diagnostics.Append(diag...); resp.Diagnostics.HasError() {
		return
	}

	if !found {
		resp.Diagnostics.AddError(
			"update instance",
			fmt.Sprintf("instance %d not found after update", state.Id.ValueInt64()),
		)

		return
	}

	newState.ServicePlanOptions = servicePlanOptions

	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// makeUpdateAPIcalls makes any Update API calls that are needed
func makeUpdateAPIcalls(
	ctx context.Context,
	client *sdk.APIClient,
	plan, state InstanceModel,
	updateTimeout time.Duration,
	resp *resource.UpdateResponse,
) ServicePlanOptionsValue {
	updateInstance := client.InstancesAPI.UpdateInstance(ctx, plan.Id.ValueInt64())
	updateRequest := &sdk.UpdateInstanceRequest{}
	instanceUpdateRequest := &sdk.UpdateInstanceRequestInstance{}
	updateConfig := &sdk.UpdateInstanceRequestConfig{}
	hasConfigUpdate := false

	// name
	if !plan.Name.IsNull() && !plan.Name.IsUnknown() {
		instanceUpdateRequest.Name = plan.Name.ValueStringPointer()
		instanceUpdateRequest.DisplayName = plan.Name.ValueStringPointer()
	}

	// description
	if !plan.Description.IsNull() {
		instanceUpdateRequest.Description = plan.Description.ValueStringPointer()
	}

	// instance_context
	if !plan.InstanceContext.IsNull() && !plan.InstanceContext.IsUnknown() {
		instanceUpdateRequest.InstanceContext = plan.InstanceContext.ValueStringPointer()
	}

	// group_id
	if !plan.GroupId.IsNull() && !plan.GroupId.IsUnknown() {
		site := &sdk.UpdateInstanceRequestInstanceSite{}
		site.Id = plan.GroupId.ValueInt64Pointer()
		instanceUpdateRequest.Site = site
	}

	// tags
	if !plan.Tags.IsNull() && !plan.Tags.IsUnknown() {
		tags, diags := convert.FromSetType(ctx, plan.Tags, updateTagMapper)
		if diags.HasError() {
			tflog.Error(ctx, "cannot convert tags")
			resp.Diagnostics.Append(diags...)

			return NewServicePlanOptionsValueNull()
		}
		instanceUpdateRequest.Tags = tags
	} else {
		// Morpheus will only remove tags if explicitly passed an empty list
		instanceUpdateRequest.Tags = []sdk.UpdateInstanceRequestInstanceTagsInner{}
	}

	// labels
	if !plan.Labels.IsNull() && !plan.Labels.IsUnknown() {
		var labels []string
		resp.Diagnostics.Append(plan.Labels.ElementsAs(ctx, &labels, false)...)
		if resp.Diagnostics.HasError() {
			return NewServicePlanOptionsValueNull()
		}
		instanceUpdateRequest.Labels = labels
	}

	// Config update handling. Currently only Azure uses map-based config updates
	// via AdditionalProperties. Other config types (HVM, VMware, AWS) use
	// requires_replace on their config blocks so no update logic is needed.
	switch {
	case !plan.ConfigAzure.IsNull() && !plan.ConfigAzure.IsUnknown():
		updateConfig.AdditionalProperties = make(map[string]interface{})

		if !plan.ConfigAzure.ResourcePoolId.IsNull() && !plan.ConfigAzure.ResourcePoolId.IsUnknown() {
			updateConfig.AdditionalProperties["resourcePoolId"] = plan.ConfigAzure.ResourcePoolId.ValueString()
		}

		if !plan.ConfigAzure.CreateUser.IsNull() && !plan.ConfigAzure.CreateUser.IsUnknown() {
			updateConfig.AdditionalProperties["createUser"] = plan.ConfigAzure.CreateUser.ValueBool()
		}

		if !plan.ConfigAzure.AzureRegion.IsNull() && !plan.ConfigAzure.AzureRegion.IsUnknown() {
			updateConfig.AdditionalProperties["azureRegion"] = plan.ConfigAzure.AzureRegion.ValueString()
		}

		if !plan.ConfigAzure.AzuresecurityGroupId.IsNull() && !plan.ConfigAzure.AzuresecurityGroupId.IsUnknown() {
			updateConfig.AdditionalProperties["azuresecurityGroupId"] = plan.ConfigAzure.AzuresecurityGroupId.ValueString()
		}

		if !plan.ConfigAzure.AvailabilityOptions.IsNull() && !plan.ConfigAzure.AvailabilityOptions.IsUnknown() {
			updateConfig.AdditionalProperties["availabilityOptions"] = plan.ConfigAzure.AvailabilityOptions.ValueString()
		}

		if !plan.ConfigAzure.AvailabilitySet.IsNull() && !plan.ConfigAzure.AvailabilitySet.IsUnknown() {
			updateConfig.AdditionalProperties["availabilitySet"] = plan.ConfigAzure.AvailabilitySet.ValueString()
		}

		if !plan.ConfigAzure.AvailabilityZone.IsNull() && !plan.ConfigAzure.AvailabilityZone.IsUnknown() {
			updateConfig.AdditionalProperties["availabilityZone"] = plan.ConfigAzure.AvailabilityZone.ValueString()
		}

		if !plan.ConfigAzure.AzurefloatingIp.IsNull() && !plan.ConfigAzure.AzurefloatingIp.IsUnknown() {
			updateConfig.AdditionalProperties["azurefloatingIp"] = plan.ConfigAzure.AzurefloatingIp.ValueString()
		}

		if !plan.ConfigAzure.BootDiagnostics.IsNull() && !plan.ConfigAzure.BootDiagnostics.IsUnknown() {
			updateConfig.AdditionalProperties["bootDiagnostics"] = plan.ConfigAzure.BootDiagnostics.ValueString()
		}

		if !plan.ConfigAzure.OsGuestDiagnostics.IsNull() && !plan.ConfigAzure.OsGuestDiagnostics.IsUnknown() {
			updateConfig.AdditionalProperties["osGuestDiagnostics"] = plan.ConfigAzure.OsGuestDiagnostics.ValueString()
		}

		if !plan.ConfigAzure.DiagnosticsStorageAccount.IsNull() && !plan.ConfigAzure.DiagnosticsStorageAccount.IsUnknown() {
			updateConfig.AdditionalProperties["diagnosticsStorageAccount"] = plan.ConfigAzure.DiagnosticsStorageAccount.ValueString()
		}

		hasConfigUpdate = len(updateConfig.AdditionalProperties) > 0
	}

	updateRequest.Instance = instanceUpdateRequest
	if hasConfigUpdate {
		updateRequest.Config = updateConfig
	}
	_, httpResp, err := updateInstance.UpdateInstanceRequest(*updateRequest).Execute()
	if err != nil || httpResp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError("error updating instance", errfmt.ErrMsg(err, httpResp))

		return NewServicePlanOptionsValueNull()
	}

	// if we need to update volumes, network interfaces or service plan options, then we need to make the
	// resize API call
	servicePlanOptions, diags := makeResizeAPICalls(ctx, client, plan, state, updateTimeout)
	resp.Diagnostics.Append(diags...)

	return servicePlanOptions
}

// makeResizeAPICalls makes any Resize API calls that are needed.
func makeResizeAPICalls(
	ctx context.Context,
	client *sdk.APIClient,
	plan, state InstanceModel,
	updateTimeout time.Duration,
) (ServicePlanOptionsValue, diag.Diagnostics) {
	var d diag.Diagnostics

	resizeRequest := createResizeRequest(plan, state)

	d.Append(addVolumesToResizeRequest(ctx, plan, state, resizeRequest)...)
	if d.HasError() {
		return state.ServicePlanOptions, d
	}

	d.Append(addNetworkInterfacesToResizeRequest(ctx, plan, state, resizeRequest)...)
	if d.HasError() {
		return state.ServicePlanOptions, d
	}

	addServicePlanOptionsToResizeRequest(plan, state, resizeRequest)

	if resizeRequest.Volumes != nil || resizeRequest.NetworkInterfaces != nil ||
		resizeRequest.ServicePlanOptions != nil {
		d.Append(makeResizeRequestAndWaitForComplete(ctx, client, resizeRequest, updateTimeout)...)
		if d.HasError() {
			return state.ServicePlanOptions, d
		}

	}

	if resizeRequest.ServicePlanOptions == nil {
		return state.ServicePlanOptions, d
	}

	return plan.ServicePlanOptions, d
}

func createResizeRequest(
	plan InstanceModel,
	state InstanceModel,
) *sdk.ResizeInstanceRequest {
	resizeRequest := &sdk.ResizeInstanceRequest{}

	// plan_id
	if !plan.PlanId.IsNull() || !plan.PlanId.IsUnknown() {
		resizeRequest.Instance = &sdk.ResizeInstanceRequestInstance{}
		resizeRequest.Instance.Id = state.Id.ValueInt64Pointer()
		resizeRequest.Instance.Plan = &sdk.ResizeInstanceRequestInstancePlan{}
		resizeRequest.Instance.Plan.Id = plan.PlanId.ValueInt64Pointer()
	}

	return resizeRequest
}

func addVolumesToResizeRequest(
	ctx context.Context,
	plan, state InstanceModel,
	resizeRequest *sdk.ResizeInstanceRequest,
) diag.Diagnostics {
	var pVolumes []VolumesValue
	pdiags := plan.Volumes.ElementsAs(ctx, &pVolumes, false)
	if pdiags.HasError() {
		tflog.Error(ctx, "cannot convert plan volumes")

		return pdiags
	}

	var sVolumes []VolumesValue
	sdiags := state.Volumes.ElementsAs(ctx, &sVolumes, false)
	if sdiags.HasError() {
		tflog.Error(ctx, "cannot convert state volumes")

		return sdiags
	}

	// Build the resize volume list (see buildResizeVolumes) and only issue a resize
	// when it reports an actual change - a volume added, removed, or resized (a
	// per-volume size change). A rename or reorder of the same set of volumes is not
	// a change that requires a resize.
	volumes, resizing := buildResizeVolumes(pVolumes, sVolumes)

	if resizing {
		volumesForRequest := make([]sdk.ResizeInstanceRequestVolumesInner, len(volumes))
		for i, v := range volumes {
			volumesForRequest[i] = updateVolumeMapper(v)
		}

		resizeRequest.Volumes = volumesForRequest
	}

	return diag.Diagnostics{}
}

// buildResizeVolumes builds the ordered volume list for an instance resize request
// from the planned and current (state) volumes, and reports whether a resize is
// required. The morpheus-ui resize changelist
// (CloudPluginProvisioningService.buildResizeRequest -> provider.buildVolumeChangelist)
// keys existing volumes on their id, so each request volume must carry the right id:
// an existing volume keeps its state id (> 0) to be kept/resized, an added volume is
// left with a null id (updateVolumeMapper sends -1) for the API to create, and any
// state volume with no planned counterpart is omitted so the API removes it.
//
// When the plan and state hold the same number of volumes, volumes are paired
// positionally — this preserves the existing behaviour exactly: a resize is
// triggered only by a per-volume size change, so a rename or reorder of the same set
// of volumes is not misread as an add+remove and does not force an unnecessary
// resize. Only when the counts differ (a volume was genuinely added or removed) do we
// match by name (falling back to position for unnamed volumes), because that is the
// only case where positional pairing mis-assigns ids for an insert or remove in the
// middle of the list. This mirrors the read side (getVolumes).
func buildResizeVolumes(
	planVolumes, stateVolumes []VolumesValue,
) ([]VolumesValue, bool) {
	// Equal counts: positional pairing (unchanged behaviour). A resize is triggered
	// only by a per-volume size change.
	if len(planVolumes) == len(stateVolumes) {
		volumes := make([]VolumesValue, 0, len(planVolumes))
		resizing := false

		for i, pVolume := range planVolumes {
			volForRequest, different := createVolumeFromPlanAndState(pVolume, stateVolumes[i])
			if different {
				resizing = true
			}
			volumes = append(volumes, volForRequest)
		}

		return volumes, resizing
	}

	// Counts differ: a volume was added or removed, so a resize is always required.
	// Match by name (positional fallback for unnamed volumes) so each existing
	// volume keeps its id and genuinely new volumes are identified.
	consumed := make([]bool, len(stateVolumes))
	volumes := make([]VolumesValue, 0, len(planVolumes))

	for i, pVolume := range planVolumes {
		si := matchStateVolume(pVolume, stateVolumes, consumed, i)
		if si == -1 {
			// No existing counterpart: a newly added volume. updateVolumeMapper
			// sends id=-1 when the plan volume has no id, so the API creates it.
			volumes = append(volumes, pVolume)

			continue
		}

		consumed[si] = true
		// Keep the existing volume's id and pick up any size change from the plan.
		volForRequest, _ := createVolumeFromPlanAndState(pVolume, stateVolumes[si])
		volumes = append(volumes, volForRequest)
	}

	return volumes, true
}

// matchStateVolume returns the index of the state volume corresponding to the
// planned volume, or -1 when there is none (a newly added volume). A named plan
// volume is matched to the first unconsumed state volume with the same name; a
// named plan volume with no such match is treated as new. An unnamed plan volume
// falls back to the state volume at the same position when it is still unconsumed.
func matchStateVolume(
	planVolume VolumesValue,
	stateVolumes []VolumesValue,
	consumed []bool,
	i int,
) int {
	name := planVolume.Name.ValueString()
	if name != "" {
		for s := range stateVolumes {
			if !consumed[s] && stateVolumes[s].Name.ValueString() == name {
				return s
			}
		}

		return -1
	}

	if i < len(stateVolumes) && !consumed[i] {
		return i
	}

	return -1
}

func createVolumeFromPlanAndState(
	planVolume, stateVolume VolumesValue,
) (VolumesValue, bool) {
	different := false
	volume := stateVolume

	if !planVolume.Size.IsNull() && !planVolume.Size.IsUnknown() {
		if !stateVolume.Size.IsNull() && !stateVolume.Size.IsUnknown() &&
			stateVolume.Size.ValueInt64() != planVolume.Size.ValueInt64() {
			volume.Size = planVolume.Size
			different = true
		}
	}

	return volume, different
}

func addNetworkInterfacesToResizeRequest(
	ctx context.Context,
	plan, state InstanceModel,
	resizeRequest *sdk.ResizeInstanceRequest,
) diag.Diagnostics {
	resizing := false

	var pIntfs []NetworkInterfacesValue
	pdiags := plan.NetworkInterfaces.ElementsAs(ctx, &pIntfs, false)
	if pdiags.HasError() {
		tflog.Error(ctx, "cannot convert plan network interfaces")

		return pdiags
	}

	var sIntfs []NetworkInterfacesValue
	sdiags := state.NetworkInterfaces.ElementsAs(ctx, &sIntfs, false)
	if sdiags.HasError() {
		tflog.Error(ctx, "cannot convert state network interfaces")

		return sdiags
	}

	// Always resize if the number of network interfaces in the plan and state are different
	if len(pIntfs) != len(sIntfs) {
		resizing = true
	}

	// Also check if any child virtual network lists differ in length
	for i := range pIntfs {
		if i < len(sIntfs) {
			childMatch, cdiags := compare.ListsMatch[ChildVirtualNetworksValue](
				ctx, pIntfs[i].ChildVirtualNetworks, sIntfs[i].ChildVirtualNetworks,
			)
			if cdiags.HasError() {
				return cdiags
			}

			if !childMatch {
				resizing = true
			}
		}
	}

	// Go through the lists of network interfaces from state and in the plan.  Assume that network interfaces
	// are in the same order in both lists.  If the network interface in the plan is different from the network
	// interface in the state, then we add the network interface from the plan to the list of network interfaces
	// for the resize request.  If the network interface in the plan is the same as the network interface in the
	// state, then we add the network interface from the state to the list of network interfaces for the resize
	// request.
	var intfs []NetworkInterfacesValue
	for i, pIntf := range pIntfs {
		if i < len(sIntfs) {
			sIntf := sIntfs[i]
			intfForRequest, different := createNetworkInterfaceFromPlanAndState(ctx, pIntf, sIntf)
			// Resize if the network interface for the request is different from the network interface in the state
			if different {
				resizing = true
			}
			intfs = append(intfs, intfForRequest)
		} else {
			intfs = append(intfs, pIntf)
		}
	}

	if resizing {
		mapper := updateNetworkInterfaceMapper(ctx)
		intfsForRequest := make([]sdk.InstancesNetworkInterfaces4, len(intfs))
		for i, intf := range intfs {
			intfsForRequest[i] = mapper(intf)
		}

		resizeRequest.NetworkInterfaces = intfsForRequest
	}

	return diag.Diagnostics{}
}

func createNetworkInterfaceFromPlanAndState(
	ctx context.Context,
	planIntf, stateIntf NetworkInterfacesValue,
) (NetworkInterfacesValue, bool) {
	different := false
	intf := stateIntf

	// Always preserve the id from state
	intf.Id = stateIntf.Id

	if !planIntf.NetworkId.IsNull() && !planIntf.NetworkId.IsUnknown() {
		if !stateIntf.NetworkId.IsNull() && !stateIntf.NetworkId.IsUnknown() &&
			stateIntf.NetworkId.ValueInt64() != planIntf.NetworkId.ValueInt64() {
			intf.NetworkId = planIntf.NetworkId
			different = true
		}
	}

	if !planIntf.NetworkGroupId.IsNull() && !planIntf.NetworkGroupId.IsUnknown() {
		if !stateIntf.NetworkGroupId.IsNull() && !stateIntf.NetworkGroupId.IsUnknown() &&
			stateIntf.NetworkGroupId.ValueInt64() != planIntf.NetworkGroupId.ValueInt64() {
			intf.NetworkGroupId = planIntf.NetworkGroupId
			different = true
		}
	}

	if !planIntf.IpMode.IsNull() && !planIntf.IpMode.IsUnknown() {
		if !stateIntf.IpMode.IsNull() && !stateIntf.IpMode.IsUnknown() &&
			stateIntf.IpMode.ValueString() != planIntf.IpMode.ValueString() {
			intf.IpMode = planIntf.IpMode
			different = true
		}
	}

	if !planIntf.IpAddress.IsNull() && !planIntf.IpAddress.IsUnknown() {
		if !stateIntf.IpAddress.IsNull() && !stateIntf.IpAddress.IsUnknown() &&
			stateIntf.IpAddress.ValueString() != planIntf.IpAddress.ValueString() {
			intf.IpAddress = planIntf.IpAddress
			different = true
		}
	}

	if !planIntf.IpPool.IsNull() && !planIntf.IpPool.IsUnknown() {
		if !stateIntf.IpPool.IsNull() && !stateIntf.IpPool.IsUnknown() &&
			stateIntf.IpPool.ValueInt64() != planIntf.IpPool.ValueInt64() {
			intf.IpPool = planIntf.IpPool
			different = true
		}
	}

	if !planIntf.NetworkTypeId.IsNull() && !planIntf.NetworkTypeId.IsUnknown() {
		if !stateIntf.NetworkTypeId.IsNull() && !stateIntf.NetworkTypeId.IsUnknown() &&
			stateIntf.NetworkTypeId.ValueInt64() != planIntf.NetworkTypeId.ValueInt64() {
			intf.NetworkTypeId = planIntf.NetworkTypeId
			different = true
		}
	}

	// Process child virtual networks
	var pChildren []ChildVirtualNetworksValue

	pdiags := planIntf.ChildVirtualNetworks.ElementsAs(ctx, &pChildren, false)
	if pdiags.HasError() {
		tflog.Error(ctx, "cannot convert plan child virtual networks")

		return intf, different
	}

	var sChildren []ChildVirtualNetworksValue

	sdiags := stateIntf.ChildVirtualNetworks.ElementsAs(ctx, &sChildren, false)
	if sdiags.HasError() {
		tflog.Error(ctx, "cannot convert state child virtual networks")

		return intf, different
	}

	if len(pChildren) != len(sChildren) {
		different = true
	}

	var children []ChildVirtualNetworksValue
	for i, pChild := range pChildren {
		if i < len(sChildren) {
			sChild := sChildren[i]
			childForRequest, childDifferent := createChildVirtualNetworkFromPlanAndState(pChild, sChild)
			if childDifferent {
				different = true
			}
			children = append(children, childForRequest)
		} else {
			children = append(children, pChild)
		}
	}

	if len(children) > 0 {
		childList, cdiags := convert.ToListType(ctx, children, func(in ChildVirtualNetworksValue) ChildVirtualNetworksValue {
			return in
		})
		if !cdiags.HasError() {
			intf.ChildVirtualNetworks = childList
		}
	}

	return intf, different
}

func createChildVirtualNetworkFromPlanAndState(
	planChild, stateChild ChildVirtualNetworksValue,
) (ChildVirtualNetworksValue, bool) {
	different := false
	child := stateChild

	// Always preserve the id from state
	child.Id = stateChild.Id

	if !planChild.NetworkId.IsNull() && !planChild.NetworkId.IsUnknown() {
		if !stateChild.NetworkId.IsNull() && !stateChild.NetworkId.IsUnknown() &&
			stateChild.NetworkId.ValueInt64() != planChild.NetworkId.ValueInt64() {
			child.NetworkId = planChild.NetworkId
			different = true
		}
	}

	if !planChild.NetworkGroupId.IsNull() && !planChild.NetworkGroupId.IsUnknown() {
		if !stateChild.NetworkGroupId.IsNull() && !stateChild.NetworkGroupId.IsUnknown() &&
			stateChild.NetworkGroupId.ValueInt64() != planChild.NetworkGroupId.ValueInt64() {
			child.NetworkGroupId = planChild.NetworkGroupId
			different = true
		}
	}

	if !planChild.IpMode.IsNull() && !planChild.IpMode.IsUnknown() {
		if !stateChild.IpMode.IsNull() && !stateChild.IpMode.IsUnknown() &&
			stateChild.IpMode.ValueString() != planChild.IpMode.ValueString() {
			child.IpMode = planChild.IpMode
			different = true
		}
	}

	if !planChild.IpAddress.IsNull() && !planChild.IpAddress.IsUnknown() {
		if !stateChild.IpAddress.IsNull() && !stateChild.IpAddress.IsUnknown() &&
			stateChild.IpAddress.ValueString() != planChild.IpAddress.ValueString() {
			child.IpAddress = planChild.IpAddress
			different = true
		}
	}

	if !planChild.IpPool.IsNull() && !planChild.IpPool.IsUnknown() {
		if !stateChild.IpPool.IsNull() && !stateChild.IpPool.IsUnknown() &&
			stateChild.IpPool.ValueInt64() != planChild.IpPool.ValueInt64() {
			child.IpPool = planChild.IpPool
			different = true
		}
	}

	if !planChild.NetworkTypeId.IsNull() && !planChild.NetworkTypeId.IsUnknown() {
		if !stateChild.NetworkTypeId.IsNull() && !stateChild.NetworkTypeId.IsUnknown() &&
			stateChild.NetworkTypeId.ValueInt64() != planChild.NetworkTypeId.ValueInt64() {
			child.NetworkTypeId = planChild.NetworkTypeId
			different = true
		}
	}

	return child, different
}

func addServicePlanOptionsToResizeRequest(
	plan, state InstanceModel,
	resizeRequest *sdk.ResizeInstanceRequest,
) {
	// compare state and plan service_plan_options so we only resize if required
	servicePlanOptions := &sdk.ResizeInstanceRequestServicePlanOptions{}
	if !plan.ServicePlanOptions.MaxCores.IsNull() && !plan.ServicePlanOptions.MaxCores.IsUnknown() {
		servicePlanOptions.MaxCores = plan.ServicePlanOptions.MaxCores.ValueInt64Pointer()
	}
	if !plan.ServicePlanOptions.CoresPerSocket.IsNull() && !plan.ServicePlanOptions.CoresPerSocket.IsUnknown() {
		servicePlanOptions.CoresPerSocket = plan.ServicePlanOptions.CoresPerSocket.ValueInt64Pointer()
	}
	if !plan.ServicePlanOptions.MaxMemory.IsNull() && !plan.ServicePlanOptions.MaxMemory.IsUnknown() {
		memoryInBytes := *plan.ServicePlanOptions.MaxMemory.ValueInt64Pointer() << 20
		servicePlanOptions.MaxMemory = &memoryInBytes
	}

	resizeRequest.ServicePlanOptions = servicePlanOptions
}

func makeResizeRequestAndWaitForComplete(
	ctx context.Context,
	client *sdk.APIClient,
	resizeRequest *sdk.ResizeInstanceRequest,
	updateTimeout time.Duration,
) diag.Diagnostics {
	var d diag.Diagnostics

	if resizeRequest.Instance == nil {
		d.AddError("instance resize failed", "resize request Instance is nil")

		return d
	}

	if resizeRequest.Instance.Id == nil {
		d.AddError("instance resize failed", "resize request Instance ID is nil")

		return d
	}

	resizeInstance := client.InstancesAPI.ResizeInstance(ctx, *resizeRequest.Instance.Id)
	resizeResp, httpResp, err := resizeInstance.ResizeInstanceRequest(*resizeRequest).Execute()
	if err != nil || httpResp.StatusCode != http.StatusOK {
		d.AddError("instance resize failed", errfmt.ErrMsg(err, httpResp))

		return d
	}

	tflog.Info(ctx, fmt.Sprintln(resizeResp))

	waitForReady := func() (string, error) {
		resp, hresp, err := client.InstancesAPI.GetInstance(ctx, *resizeRequest.Instance.Id).Execute()
		if err != nil {
			if hresp == nil || hresp.StatusCode != http.StatusOK {
				return "", backoff.Permanent(err)
			}
		}

		inst := resp.Instance
		if inst == nil {
			return "", backoff.Permanent(fmt.Errorf("instance %d: GET returned empty instance", *resizeRequest.Instance.Id))
		}

		if inst.Status == nil {
			return "", backoff.Permanent(fmt.Errorf("instance %d: GET returned empty status", *resizeRequest.Instance.Id))
		}

		return *inst.Status, checkStatusDone(
			*inst.Status,
			UpdateTargetStatuses,
			UpdateErrorStatuses,
		)
	}

	if status, err := backoff.Retry(
		ctx,
		waitForReady,
		backoff.WithBackOff(backoff.NewConstantBackOff(5*time.Second)),
		backoff.WithMaxElapsedTime(updateTimeout),
	); err != nil {
		d.AddError(
			"resize instance resource",
			fmt.Sprintf(
				"instance %d: resizing failed current status is: %s",
				*resizeRequest.Instance.Id,
				status,
			),
		)
	}

	return d
}

// isAPIUpdateNeeded is a function that will compare the plan and state of attributes
// that can be updated, and if the only attribute is Timeouts then it returns false
func isAPIUpdateNeeded(plan, state InstanceModel) bool {
	// name
	if plan.Name != state.Name {
		return true
	}

	// description
	if plan.Description != state.Description {
		return true
	}

	// instance_context
	if !plan.InstanceContext.Equal(state.InstanceContext) {
		return true
	}

	// group_id
	if !plan.GroupId.Equal(state.GroupId) {
		return true
	}

	// config_azure
	if !plan.ConfigAzure.Equal(state.ConfigAzure) {
		return true
	}

	// tags
	if !plan.Tags.Equal(state.Tags) {
		return true
	}

	// volumes
	if !plan.Volumes.Equal(state.Volumes) {
		return true
	}

	// network-interfaces
	if !plan.NetworkInterfaces.Equal(state.NetworkInterfaces) {
		return true
	}

	// service_plan_options
	if !plan.ServicePlanOptions.Equal(state.ServicePlanOptions) {
		return true
	}

	// timeouts - this should be the last comparison
	if !plan.Timeouts.Equal(state.Timeouts) {
		return false
	}

	// For safety's sake we will return true by default
	return true
}
