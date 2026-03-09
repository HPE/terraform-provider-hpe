// (C) Copyright 2025-2026 Hewlett Packard Enterprise Development LP

package instance

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/cenkalti/backoff/v5"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"

	errfmt "github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
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

	newState, diag := getInstanceAsState(ctx, state.Id.ValueInt64(), client, plan)
	if resp.Diagnostics.Append(diag...); resp.Diagnostics.HasError() {
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
	updateRequest := sdk.NewUpdateInstanceRequest()
	instanceUpdateRequest := sdk.NewUpdateInstanceRequestInstance()

	// name
	if !plan.Name.IsNull() && !plan.Name.IsUnknown() {
		instanceUpdateRequest.Name = plan.Name.ValueStringPointer()
		instanceUpdateRequest.DisplayName = plan.Name.ValueStringPointer()
	}

	// TODO: DESCRIPTION IS MISSING FROM THE SCHEMA??
	// if !plan.Description.IsNull() {
	// 	instanceUpdateRequest.Description = plan.Description.ValueStringPointer()
	// }

	// instance_context
	if !plan.InstanceContext.IsNull() && !plan.InstanceContext.IsUnknown() {
		instanceUpdateRequest.InstanceContext = plan.InstanceContext.ValueStringPointer()
	}

	// group_id
	if !plan.GroupId.IsNull() && !plan.GroupId.IsUnknown() {
		site := sdk.NewUpdateInstanceRequestInstanceSite()
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
		instanceUpdateRequest.SetTags(tags)
	}

	updateRequest.SetInstance(*instanceUpdateRequest)
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
	resizeRequest := sdk.NewResizeInstanceRequestWithDefaults()

	// plan_id
	if !plan.PlanId.IsNull() || !plan.PlanId.IsUnknown() {
		resizeRequest.Instance = sdk.NewResizeInstanceRequestInstance()
		resizeRequest.Instance.Id = state.Id.ValueInt64Pointer()
		resizeRequest.Instance.Plan = sdk.NewResizeInstanceRequestInstancePlan()
		resizeRequest.Instance.Plan.SetId(plan.PlanId.ValueInt64())
	}

	return resizeRequest
}

func addVolumesToResizeRequest(
	ctx context.Context,
	plan, state InstanceModel,
	resizeRequest *sdk.ResizeInstanceRequest,
) diag.Diagnostics {
	resizing := false

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

	// Always resize if the number of volumes in the plan and state are different
	if len(pVolumes) != len(sVolumes) {
		resizing = true
	}

	// Go through the lists of volumes from state and in the plan.  Assume that volumes are in the same order
	// in both lists.  If the volume in the plan is different from the volume in the state, then we add the volume
	// from the plan to the list of volumes for the resize request.  If the volume in the plan is the same as the
	// volume in the state, then we add the volume from the state to the list of volumes for the resize request.
	var volumes []VolumesValue
	for i, pVolume := range pVolumes {
		if i < len(sVolumes) {
			sVolume := sVolumes[i]
			volForRequest, different := createVolumeFromPlanAndState(pVolume, sVolume)
			// Resize if the volume for the request is different from the volume in the state
			if different {
				resizing = true
			}
			volumes = append(volumes, volForRequest)
		} else {
			volumes = append(volumes, pVolume)
		}
	}

	if resizing {
		volumesForRequest := make([]sdk.ResizeInstanceRequestVolumesInner, len(volumes))
		for i, v := range volumes {
			volumesForRequest[i] = updateVolumeMapper(v)
		}

		resizeRequest.SetVolumes(volumesForRequest)
	}

	return diag.Diagnostics{}
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
	// compare state and plan network_interfaces so we only resize if required
	// TODO: make this compare each network interface rather than just length
	// TODO: if we're resizing network interfaces, then we need to remove the PlanModifier and RequiresReplace from Schema.
	// TODO: and add the id of the interface to the Schema
	// TODO: resize doesn't work for network interfaces at present, need to revisit this when the API supports it
	intfsMatch, ndiags := listsMatch[NetworkInterfacesValue](ctx, plan.NetworkInterfaces, state.NetworkInterfaces)
	if ndiags.HasError() {
		tflog.Error(ctx, "cannot compare network interfaces")

		return ndiags
	}

	if !intfsMatch {
		networkInterfaces, idiags := convert.FromListType(
			ctx,
			plan.NetworkInterfaces,
			updateNetworkInterfaceMapper(ctx),
		)
		if idiags.HasError() {
			tflog.Error(ctx, "cannot convert network interfaces")

			return idiags
		}

		resizeRequest.SetNetworkInterfaces(networkInterfaces)
	}

	return diag.Diagnostics{}
}

func addServicePlanOptionsToResizeRequest(
	plan, state InstanceModel,
	resizeRequest *sdk.ResizeInstanceRequest,
) {
	// compare state and plan service_plan_options so we only resize if required
	if !plan.ServicePlanOptions.Equal(state.ServicePlanOptions) {
		servicePlanOptions := sdk.NewResizeInstanceRequestServicePlanOptions()
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

		resizeRequest.SetServicePlanOptions(*servicePlanOptions)
	}
}

func makeResizeRequestAndWaitForComplete(
	ctx context.Context,
	client *sdk.APIClient,
	resizeRequest *sdk.ResizeInstanceRequest,
	updateTimeout time.Duration,
) diag.Diagnostics {
	var d diag.Diagnostics
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

		// Get instance
		inst, ok := resp.GetInstanceOk()
		if !ok || inst == nil {
			return "", backoff.Permanent(fmt.Errorf("instance %d: GET returned empty instance", *resizeRequest.Instance.Id))
		}

		// Get status
		status, ok := inst.GetStatusOk()
		if !ok || status == nil {
			return "", backoff.Permanent(fmt.Errorf("instance %d: GET returned empty status", *resizeRequest.Instance.Id))
		}

		return *status, checkStatusDone(
			*status,
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

	// TODO add description here
	// description

	// instance_context
	if !plan.InstanceContext.Equal(state.InstanceContext) {
		return true
	}

	// group_id
	if !plan.GroupId.Equal(state.GroupId) {
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

// listsMatch is a generic that compares lists from plan and state to see if they are the same
// Returns true if they are, false otherwise
func listsMatch[S attr.Value](
	ctx context.Context,
	planList, stateList basetypes.ListValue,
) (bool, diag.Diagnostics) {
	var planVals, stateVals []S

	diags := planList.ElementsAs(ctx, &planVals, false)
	if diags.HasError() {
		tflog.Error(ctx, fmt.Sprintf("cannot convert list values to type %T", planVals))

		return false, diags
	}

	diags = stateList.ElementsAs(ctx, &stateVals, false)
	if diags.HasError() {
		tflog.Error(ctx, fmt.Sprintf("cannot convert list values to type %T", stateVals))

		return false, diags
	}

	// Check length of lists first
	if len(planVals) != len(stateVals) {
		return false, nil
	}

	// Compare each element in the lists to see if they are the same
	for i, planVal := range planVals {
		stateVal := stateVals[i]

		if !planVal.Equal(stateVal) {
			return false, nil
		}
	}

	return true, nil
}
