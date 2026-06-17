// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package instanceclone

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"
	"github.com/cenkalti/backoff/v5"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/HPE/terraform-provider-hpe/morpheus/configure"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/utils/cleanup"
	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &Resource{}
	_ resource.ResourceWithImportState = &Resource{}
)

var (
	cloneTargetStatuses = []string{"running", "stopped"}
	cloneErrorStatuses  = []string{"failed", "errored"}
)

func NewResource() resource.Resource {
	return &Resource{}
}

// Resource defines the resource implementation.
type Resource struct {
	configure.ResourceWithMorpheusConfigure
}

func (r *Resource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_morpheus_instance_clone"
}

func (r *Resource) Schema(
	ctx context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = InstanceCloneResourceSchema(ctx)
}

func checkStatusDone(status string, targetStatuses, errorStatuses []string) error {
	for _, s := range errorStatuses {
		if strings.EqualFold(status, s) {
			return backoff.Permanent(errors.New("reached error status: " + status))
		}
	}
	for _, s := range targetStatuses {
		if strings.EqualFold(status, s) {
			return nil
		}
	}

	return fmt.Errorf("instance status %q not yet in target set", status)
}

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

	if !plan.ResourcePoolId.IsNull() && !plan.ResourcePoolId.IsUnknown() {
		rpID := plan.ResourcePoolId.ValueString()
		cloneReq.Config = &sdk.CloneInstanceRequestConfig{
			ResourcePoolId: &rpID,
		}
	}

	// Build volumes
	cloneReq.Volumes = buildCloneVolumes(ctx, plan.Volumes)

	// Build network interfaces
	cloneReq.NetworkInterfaces = buildCloneNetworkInterfaces(ctx, plan.NetworkInterfaces)

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
			return nil, backoff.Permanent(
				fmt.Errorf("failed to list instances: %s", errfmt.ErrMsg(err, hresp)),
			)
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
			return nil, backoff.Permanent(
				fmt.Errorf("failed to get clone instance %d: %s",
					cloneID, errfmt.ErrMsg(err, hresp)),
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
		resizeReq := sdk.ResizeInstanceRequest{
			Volumes:           buildResizeVolumes(ctx, plan.Volumes),
			NetworkInterfaces: buildResizeNetworkInterfaces(ctx, plan.NetworkInterfaces),
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

	// Refresh volumes and network interfaces from API
	state.Volumes = mapVolumesFromGetInstance(ctx, inst)
	state.NetworkInterfaces = mapInterfacesFromGetInstance(ctx, inst)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

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

	resizeReq.Volumes = buildResizeVolumes(ctx, plan.Volumes)
	resizeReq.NetworkInterfaces = buildResizeNetworkInterfaces(ctx, plan.NetworkInterfaces)

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
			return nil, backoff.Permanent(
				fmt.Errorf("failed to get instance %d: %s",
					instanceID, errfmt.ErrMsg(err, hresp)),
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

func (r *Resource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	var state InstanceCloneModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	deleteTimeout, diags := state.Timeouts.Delete(ctx, 10*time.Minute)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, deleteTimeout)
	defer cancel()

	client, err := r.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"delete instance clone",
			"failed to create client: "+err.Error(),
		)

		return
	}

	instanceID := state.Id.ValueInt64()

	_, hresp, err := client.InstancesAPI.DeleteInstance(ctx, instanceID).Execute()
	if err != nil {
		if hresp != nil && hresp.StatusCode == http.StatusNotFound {
			tflog.Warn(ctx, fmt.Sprintf("Instance clone %d already gone", instanceID))

			return
		}

		resp.Diagnostics.AddError(
			"delete instance clone",
			fmt.Sprintf("instance %d DELETE failed: %s",
				instanceID, errfmt.ErrMsg(err, hresp)),
		)

		return
	}
}

func (r *Resource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	instanceID, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(
			"import instance clone",
			fmt.Sprintf("invalid instance id %q: %v", req.ID, err),
		)

		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), instanceID)...)
	// source_instance_id is not returned by the clone API as a first-class
	// field. Set it null here; Read makes a best-effort attempt to recover it
	// from the clone's config (cloneInstanceId). If the platform does not expose
	// that field, set source_instance_id in configuration after import.
	resp.Diagnostics.Append(resp.State.SetAttribute(
		ctx, path.Root("source_instance_id"), types.Int64Null())...)
}

// cloneNeedsResize reports whether any requested volume size exceeds the actual
// provisioned volume size (matched positionally), meaning a post-clone resize is
// required. Used to grow data disks that some clouds leave at the source size.
func cloneNeedsResize(
	ctx context.Context,
	planVolumes types.List,
	actual []sdk.AddInstance200ResponseAllOfOneOfInstanceVolumesInner,
) bool {
	if planVolumes.IsNull() || planVolumes.IsUnknown() {
		return false
	}

	var pv []VolumesValue
	planVolumes.ElementsAs(ctx, &pv, false)

	for i, v := range pv {
		if v.Size.IsNull() || v.Size.IsUnknown() || i >= len(actual) {
			continue
		}

		var actualSize int64
		if actual[i].Size != nil {
			actualSize = *actual[i].Size
		}

		if v.Size.ValueInt64() > actualSize {
			return true
		}
	}

	return false
}

// sourceInstanceIDFromConfig returns cloneInstanceId - the source instance id
// that Morpheus stamps onto a clone's config during cloning - from the instance
// read response. It is absent for instances that were not created via the clone
// endpoint, in which case the second return value is false.
func sourceInstanceIDFromConfig(inst *sdk.GetInstance200ResponseInstance) (int64, bool) {
	if inst == nil || inst.Config == nil || inst.Config.CloneInstanceId == nil {
		return 0, false
	}

	return *inst.Config.CloneInstanceId, true
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
	state.Volumes = mapVolumesFromGetInstance(ctx, inst)
	state.NetworkInterfaces = mapInterfacesFromGetInstance(ctx, inst)

	return diags
}

// buildCloneVolumes converts plan volumes to SDK CloneInstanceRequestVolumesInner.
func buildCloneVolumes(ctx context.Context, volumesList types.List) []sdk.CloneInstanceRequestVolumesInner {
	if volumesList.IsNull() || volumesList.IsUnknown() {
		return nil
	}

	var planVolumes []VolumesValue
	volumesList.ElementsAs(ctx, &planVolumes, false)

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
			vol.StorageType.Set(&stVal)
		}

		if !v.DatastoreId.IsNull() && !v.DatastoreId.IsUnknown() {
			dsID := v.DatastoreId.ValueInt64()
			vol.DatastoreId = &sdk.CloneInstanceRequestVolumesInnerDatastoreId{
				Int64: &dsID,
			}
		}

		if !v.Id.IsNull() && !v.Id.IsUnknown() {
			vol.Id = v.Id.ValueInt64Pointer()
		}

		sdkVolumes = append(sdkVolumes, vol)
	}

	return sdkVolumes
}

// buildCloneNetworkInterfaces converts plan to SDK InstancesNetworkInterfaces3.
func buildCloneNetworkInterfaces(
	ctx context.Context, ifaceList types.List,
) []sdk.InstancesNetworkInterfaces3 {
	if ifaceList.IsNull() || ifaceList.IsUnknown() {
		return nil
	}

	var planIfaces []NetworkInterfacesValue
	ifaceList.ElementsAs(ctx, &planIfaces, false)

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

		if !ni.Id.IsNull() && !ni.Id.IsUnknown() {
			iface.Id = ni.Id.ValueInt64Pointer()
		}

		sdkIfaces = append(sdkIfaces, iface)
	}

	return sdkIfaces
}

// buildResizeVolumes converts plan volumes for resize.
func buildResizeVolumes(
	ctx context.Context, volumesList types.List,
) []sdk.ResizeInstanceRequestVolumesInner {
	if volumesList.IsNull() || volumesList.IsUnknown() {
		return nil
	}

	var planVolumes []VolumesValue
	volumesList.ElementsAs(ctx, &planVolumes, false)

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
			vol.StorageType.Set(&stVal)
		}

		if !v.DatastoreId.IsNull() && !v.DatastoreId.IsUnknown() {
			dsID := v.DatastoreId.ValueInt64()
			vol.DatastoreId = &sdk.CloneInstanceRequestVolumesInnerDatastoreId{
				Int64: &dsID,
			}
		}

		if !v.Id.IsNull() && !v.Id.IsUnknown() {
			vol.Id = v.Id.ValueInt64Pointer()
		}

		sdkVolumes = append(sdkVolumes, vol)
	}

	return sdkVolumes
}

// buildResizeNetworkInterfaces converts plan network interfaces for resize.
func buildResizeNetworkInterfaces(
	ctx context.Context, ifaceList types.List,
) []sdk.InstancesNetworkInterfaces4 {
	if ifaceList.IsNull() || ifaceList.IsUnknown() {
		return nil
	}

	var planIfaces []NetworkInterfacesValue
	ifaceList.ElementsAs(ctx, &planIfaces, false)

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

		if !ni.Id.IsNull() && !ni.Id.IsUnknown() {
			iface.Id = ni.Id.ValueInt64Pointer()
		}

		sdkIfaces = append(sdkIfaces, iface)
	}

	return sdkIfaces
}

// mapVolumesFromGetInstance maps volumes from GetInstance response.
func mapVolumesFromGetInstance(
	ctx context.Context, inst *sdk.GetInstance200ResponseInstance,
) types.List {
	if len(inst.Volumes) == 0 {
		return types.ListNull(VolumesValue{}.Type(ctx))
	}

	volumeObjValues := make([]attr.Value, 0, len(inst.Volumes))
	for _, v := range inst.Volumes {
		vol := VolumesValue{
			state: attr.ValueStateKnown,
		}

		vol.Id = convert.Int64ToType(v.Id)
		vol.Name = convert.StrToType(v.Name)
		vol.Size = convert.Int64ToType(v.Size)
		vol.RootVolume = convert.BoolToType(v.RootVolume)
		vol.StorageType = convert.Int64ToType(v.StorageType)

		// DatastoreId in this response is *string; parse to int64 if present
		if v.DatastoreId != nil {
			if dsInt, parseErr := strconv.ParseInt(*v.DatastoreId, 10, 64); parseErr == nil {
				vol.DatastoreId = types.Int64Value(dsInt)
			} else {
				vol.DatastoreId = types.Int64Null()
			}
		} else {
			vol.DatastoreId = types.Int64Null()
		}

		objVal, _ := vol.ToObjectValue(ctx)
		volumeObjValues = append(volumeObjValues, objVal)
	}

	listVal, _ := types.ListValue(
		VolumesValue{}.Type(ctx),
		volumeObjValues,
	)

	return listVal
}

// mapInterfacesFromGetInstance maps network interfaces from GetInstance response.
func mapInterfacesFromGetInstance(
	ctx context.Context, inst *sdk.GetInstance200ResponseInstance,
) types.List {
	if len(inst.Interfaces) == 0 {
		return types.ListNull(NetworkInterfacesValue{}.Type(ctx))
	}

	ifaceObjValues := make([]attr.Value, 0, len(inst.Interfaces))
	for _, iface := range inst.Interfaces {
		ni := NetworkInterfacesValue{
			state: attr.ValueStateKnown,
		}

		// Interface Id is an anyOf type (Int64 | String)
		if iface.Id != nil && iface.Id.Int64 != nil {
			ni.Id = types.Int64Value(*iface.Id.Int64)
		} else {
			ni.Id = types.Int64Null()
		}

		// Network ID from nested network object
		if iface.Network != nil && iface.Network.Id != nil {
			ni.NetworkId = types.Int64Value(*iface.Network.Id)
		} else {
			ni.NetworkId = types.Int64Null()
		}

		// NetworkInterfaceTypeId is *int64 in this response type
		ni.NetworkInterfaceTypeId = convert.Int64ToType(iface.NetworkInterfaceTypeId)

		// IpMode is *string in this response type
		ni.IpMode = convert.StrToType(iface.IpMode)

		// IpAddress is *string in this response type
		ni.IpAddress = convert.StrToType(iface.IpAddress)

		objVal, _ := ni.ToObjectValue(ctx)
		ifaceObjValues = append(ifaceObjValues, objVal)
	}

	listVal, _ := types.ListValue(
		NetworkInterfacesValue{}.Type(ctx),
		ifaceObjValues,
	)

	return listVal
}
