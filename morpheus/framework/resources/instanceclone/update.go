// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package instanceclone

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/cenkalti/backoff/v5"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"

	"github.com/HPE/terraform-provider-hpe/morpheus/utils/errfmt"
	"github.com/HPE/terraform-provider-hpe/morpheus/utils/wodigest"
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

	// Only the plan and the volumes and network interfaces are updated in
	// place, by one resize request, and only when one of them changed. Every
	// other update the framework plans - the write-only bookkeeping that
	// removes a deprecated marker, say - has nothing to send: a resize that
	// repeats the current shape is not a no-op on the appliance (it reconciles
	// the whole disk list against what it was sent), so it is not issued.
	resizeReq, resizing, rdiags := resizeRequestForUpdate(ctx, plan, state)
	resp.Diagnostics.Append(rdiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Carry over immutable fields from prior state; the status is the prior
	// one unless a resize moves it.
	plan.Id = state.Id
	plan.Status = state.Status

	if resizing {
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

		statusResult, wdiags := waitForResize(ctx, client, instanceID, updateTimeout)
		resp.Diagnostics.Append(wdiags...)
		if resp.Diagnostics.HasError() {
			return
		}

		plan.Status = types.StringValue(statusResult)
	} else {
		tflog.Debug(ctx, "instance clone update has no plan, volume or network interface change; no resize issued",
			map[string]any{"instance_id": instanceID})
	}

	// Refresh actual volumes/interfaces from API
	refreshDiags := refreshStateFromAPI(ctx, client, instanceID, &plan)
	resp.Diagnostics.Append(refreshDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// source_instance_id is create-only (a change replaces the clone); an
	// update planned through write_only_revision only records the digests.
	plan.WriteOnlyRevision = state.WriteOnlyRevision
	if wodigest.Resend(ctx, req.Plan) {
		plan.WriteOnlyRevision = wodigest.NewRevision()
	}
	resp.Diagnostics.Append(wo.Store(ctx, resp.Private, req.Config)...)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// resizeRequestForUpdate builds the resize request that takes the clone from
// its prior state to the plan, and reports whether there is anything in it to
// send. The plan is carried whenever it is set, as the appliance expects; the
// volumes and the network interfaces are carried only when they changed, each
// paired with its current counterpart so that existing disks and NICs keep
// their ids (see resizeVolumes).
func resizeRequestForUpdate(
	ctx context.Context,
	plan, state InstanceCloneModel,
) (sdk.ResizeInstanceRequest, bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	req := sdk.ResizeInstanceRequest{}

	if !plan.PlanId.IsNull() && !plan.PlanId.IsUnknown() {
		pid := plan.PlanId.ValueInt64()
		req.Instance = &sdk.ResizeInstanceRequestInstance{
			Plan: &sdk.ResizeInstanceRequestInstancePlan{Id: &pid},
		}
	}

	resizing := planIDChanged(plan.PlanId, state.PlanId)

	stateVolumes, d := listElements[VolumesValue](ctx, state.Volumes)
	diags.Append(d...)

	volumes, volumesChanged, d := resizeVolumes(ctx, plan.Volumes, stateVolumes)
	diags.Append(d...)

	if volumesChanged {
		req.Volumes = volumes
		resizing = true
	}

	stateIfaces, d := listElements[NetworkInterfacesValue](ctx, state.NetworkInterfaces)
	diags.Append(d...)

	ifaces, ifacesChanged, d := resizeNetworkInterfaces(ctx, plan.NetworkInterfaces, stateIfaces)
	diags.Append(d...)

	if ifacesChanged {
		req.NetworkInterfaces = ifaces
		resizing = true
	}

	return req, resizing && !diags.HasError(), diags
}

// planIDChanged reports whether the planned plan_id asks for a plan the clone
// does not have. plan_id is optional and not computed: unset in the plan
// means "leave the plan alone", not "change it".
func planIDChanged(plan, state types.Int64) bool {
	if plan.IsNull() || plan.IsUnknown() {
		return false
	}

	return state.IsNull() || state.IsUnknown() || state.ValueInt64() != plan.ValueInt64()
}

// listElements decodes a list attribute into its element values; a null or
// unknown list has none.
func listElements[T any](ctx context.Context, list types.List) ([]T, diag.Diagnostics) {
	var diags diag.Diagnostics

	if list.IsNull() || list.IsUnknown() {
		return nil, diags
	}

	var elems []T
	diags.Append(list.ElementsAs(ctx, &elems, false)...)

	return elems, diags
}

// waitForResize polls the instance until it settles after a resize and
// returns the status it settled in.
func waitForResize(
	ctx context.Context,
	client *sdk.APIClient,
	instanceID int64,
	updateTimeout time.Duration,
) (string, diag.Diagnostics) {
	var diags diag.Diagnostics

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
		diags.AddError(
			"update instance clone",
			fmt.Sprintf("instance %d failed or timed out after resize: %v",
				instanceID, errors.Unwrap(err)),
		)

		return "", diags
	}

	return *statusResult, diags
}

// cloneNeedsResize reports whether any requested volume is larger than the
// one the clone actually got (matched positionally, see resizeVolumes),
// meaning a post-clone resize is required. Used to grow data disks that some
// clouds leave at the source size. A smaller request is not a change: the
// clone copies the source disk and a disk is never shrunk.
func cloneNeedsResize(
	ctx context.Context,
	planVolumes types.List,
	actual []VolumesValue,
) (bool, diag.Diagnostics) {
	pv, diags := listElements[VolumesValue](ctx, planVolumes)

	for i, v := range pv {
		if v.Size.IsNull() || v.Size.IsUnknown() || i >= len(actual) {
			continue
		}

		var actualSize int64
		if !actual[i].Size.IsNull() && !actual[i].Size.IsUnknown() {
			actualSize = actual[i].Size.ValueInt64()
		}

		if v.Size.ValueInt64() > actualSize {
			return true, diags
		}
	}

	return false, diags
}

// volumesFromAPI maps the instance's volumes, less the removable media the
// read side also drops, into values the resize helpers can pair with the plan.
func volumesFromAPI(
	apiVolumes []sdk.AddInstance200ResponseAllOfOneOfInstanceVolumesInner,
) []VolumesValue {
	volumes := make([]VolumesValue, 0, len(apiVolumes))

	for _, v := range apiVolumes {
		if isCDROMVolume(v) {
			continue
		}

		volumes = append(volumes, volumeValueFromAPI(v))
	}

	return volumes
}

// interfacesFromAPI maps the instance's network interfaces into values the
// resize helpers can pair with the plan.
func interfacesFromAPI(
	ctx context.Context,
	apiInterfaces []sdk.AddInstance200ResponseAllOfOneOfInstanceInterfacesInner,
) ([]NetworkInterfacesValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	ifaces := make([]NetworkInterfacesValue, 0, len(apiInterfaces))

	for _, iface := range apiInterfaces {
		ni, d := interfaceValueFromAPI(ctx, iface)
		diags.Append(d...)
		ifaces = append(ifaces, ni)
	}

	return ifaces, diags
}

// resizeVolumes builds the resize request's volume list from the planned
// volumes, each paired with the clone's current counterpart (the prior state
// on update, the freshly cloned instance on create), and reports whether the
// pairing found a change.
//
// The appliance keys a resize on volume ids
// (AbstractBoxProvisionService.buildVolumeResizeSyncLists): a volume with an
// id > 0 is kept or resized, one with id -1 is created, and the pass that
// detects removed disks dereferences every submitted volume's id, so a volume
// sent without one is not "unchanged" - on an appliance whose clone has any
// non-root disk with an external id (a CD-ROM, a data disk) it is a 500. The
// id is therefore always sent: the plan's when it is known (it is computed,
// so on an update it is not), the paired current volume's otherwise, and -1
// for a volume with no counterpart, which is how a new disk is asked for.
//
// Pairing follows the instance resource (buildResizeVolumes there): when the
// plan and the current list hold the same number of volumes they are paired
// positionally - which is safe only because ModifyPlan has already refused
// any such plan whose names differ from the current ones by position, a
// rename or a reorder alike (see volumeNameChanges): paired by position, a
// reordered list would send one disk's id under another disk's name and
// size. When the counts differ a disk was added or removed, and each planned
// volume is paired by name with the first current volume of that name not
// yet taken - so dropping the middle disk of three keeps the third disk's id
// with the third disk, and the dropped one, unpaired, is omitted from the
// request and so deleted. Names are required on a clone's volumes and the
// appliance never changes them, which makes them the stable key (a volume
// renamed in the same change as a removal or an addition therefore reads as
// a new disk).
//
// A change is a volume added or removed, or a paired volume whose size,
// datastore, storage type, storage profile or controller mount point is set
// in the plan and differs from the current one. Unknown planned values are
// computed-optional fields the practitioner did not set and are not changes.
func resizeVolumes(
	ctx context.Context,
	planVolumes types.List,
	current []VolumesValue,
) ([]sdk.ResizeInstanceRequestVolumesInner, bool, diag.Diagnostics) {
	pv, diags := listElements[VolumesValue](ctx, planVolumes)
	if planVolumes.IsNull() || planVolumes.IsUnknown() {
		return nil, false, diags
	}

	changed := len(pv) != len(current)
	volumes := make([]sdk.ResizeInstanceRequestVolumesInner, 0, len(pv))
	pairs := pairVolumes(pv, current)

	for i, v := range pv {
		vol := resizeVolume(v)
		ci := pairs[i]

		switch {
		case !v.Id.IsNull() && !v.Id.IsUnknown():
			vol.Id = v.Id.ValueInt64Pointer()
		case ci >= 0 && !current[ci].Id.IsNull() && !current[ci].Id.IsUnknown():
			vol.Id = current[ci].Id.ValueInt64Pointer()
		default:
			vol.Id = sdk.PtrInt64(-1)
		}

		if ci >= 0 && volumeChanged(v, current[ci]) {
			changed = true
		}

		volumes = append(volumes, vol)
	}

	return volumes, changed, diags
}

// pairVolumes returns, for each planned volume, the index of its current
// counterpart, or -1 for a new volume: positional when the counts are equal,
// by name otherwise (see resizeVolumes).
func pairVolumes(plan, current []VolumesValue) []int {
	pairs := make([]int, len(plan))

	if len(plan) == len(current) {
		for i := range plan {
			pairs[i] = i
		}

		return pairs
	}

	taken := make([]bool, len(current))

	for i, v := range plan {
		pairs[i] = -1

		name := v.Name.ValueString()
		if name == "" {
			continue
		}

		for ci := range current {
			if !taken[ci] && current[ci].Name.ValueString() == name {
				pairs[i] = ci
				taken[ci] = true

				break
			}
		}
	}

	return pairs
}

// volumeNameChanges reports, when the plan and the state hold the same
// number of volumes, the positions whose planned name differs from the
// state's, as "volumes[i]: old -> new", and whether the change is a pure
// reorder (the same names, in another order) rather than a rename. Either is
// refused by ModifyPlan: a volume's name is fixed when the appliance creates
// it - a resize never renames one - so a rename could not converge; and the
// clone's volumes are paired with its disks by position, so a reordered list
// would send one disk's id under another's name and size. With differing
// counts the volumes pair by name instead, and a rename reads as a removal
// plus an addition, which is what it would be.
func volumeNameChanges(plan, state []VolumesValue) (changes []string, reorder bool) {
	if len(plan) != len(state) {
		return nil, false
	}

	planNames := make([]string, 0, len(plan))
	stateNames := make([]string, 0, len(state))

	for i := range plan {
		if plan[i].Name.IsNull() || plan[i].Name.IsUnknown() ||
			state[i].Name.IsNull() || state[i].Name.IsUnknown() {
			continue
		}

		planNames = append(planNames, plan[i].Name.ValueString())
		stateNames = append(stateNames, state[i].Name.ValueString())

		if plan[i].Name.ValueString() != state[i].Name.ValueString() {
			changes = append(changes, fmt.Sprintf("volumes[%d]: %q -> %q",
				i, state[i].Name.ValueString(), plan[i].Name.ValueString()))
		}
	}

	if len(changes) == 0 {
		return nil, false
	}

	sort.Strings(planNames)
	sort.Strings(stateNames)

	return changes, slices.Equal(planNames, stateNames)
}

// rejectVolumeNameChanges adds an error to the plan when a volume is renamed
// or the volumes are reordered in place (see volumeNameChanges). Only on an
// update: a create has no prior names and a destroy no plan.
func rejectVolumeNameChanges(
	ctx context.Context,
	req resource.ModifyPlanRequest,
	resp *resource.ModifyPlanResponse,
) {
	if req.State.Raw.IsNull() || req.Plan.Raw.IsNull() {
		return
	}

	var plan, state InstanceCloneModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	planVolumes, d := listElements[VolumesValue](ctx, plan.Volumes)
	resp.Diagnostics.Append(d...)
	stateVolumes, d := listElements[VolumesValue](ctx, state.Volumes)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	changes, reorder := volumeNameChanges(planVolumes, stateVolumes)
	if len(changes) == 0 {
		return
	}

	if reorder {
		resp.Diagnostics.AddAttributeError(
			path.Root("volumes"),
			"Instance clone volumes cannot be reordered",
			"The volumes are the clone's disks in a different order. They are "+
				"matched to the disks by position, so a reordered list would resize "+
				"the wrong disk. List them in the clone's order (as in state; a new "+
				"disk goes last).\n\n"+strings.Join(changes, "\n"),
		)

		return
	}

	resp.Diagnostics.AddAttributeError(
		path.Root("volumes"),
		"Instance clone volumes cannot be renamed",
		"A volume's name is fixed when Morpheus creates it; a resize never "+
			"renames one, so the change would never converge. Keep the "+
			"current name, or replace the clone (terraform apply -replace) to "+
			"re-clone with the new one.\n\n"+strings.Join(changes, "\n"),
	)
}

// volumeChanged reports whether a planned volume differs from its current
// counterpart in a way a resize acts on.
func volumeChanged(plan, current VolumesValue) bool {
	return int64Differs(plan.Size, current.Size) ||
		int64Differs(plan.DatastoreId, current.DatastoreId) ||
		int64Differs(plan.StorageType, current.StorageType) ||
		stringDiffers(plan.StorageProfile, current.StorageProfile) ||
		stringDiffers(plan.ControllerMountPoint, current.ControllerMountPoint)
}

// int64Differs reports whether a planned value is set and differs from the
// current one. An unknown or null planned value expresses no change; a null
// current value is a change if the plan sets one.
func int64Differs(plan, current types.Int64) bool {
	if plan.IsNull() || plan.IsUnknown() {
		return false
	}

	return current.IsNull() || current.IsUnknown() || current.ValueInt64() != plan.ValueInt64()
}

// stringDiffers is int64Differs for strings.
func stringDiffers(plan, current types.String) bool {
	if plan.IsNull() || plan.IsUnknown() {
		return false
	}

	return current.IsNull() || current.IsUnknown() || current.ValueString() != plan.ValueString()
}

// resizeVolume converts one planned volume for a resize request, id excepted
// (see resizeVolumes).
func resizeVolume(v VolumesValue) sdk.ResizeInstanceRequestVolumesInner {
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

	if !v.StorageProfile.IsNull() && !v.StorageProfile.IsUnknown() {
		spVal := v.StorageProfile.ValueString()
		vol.StorageProfile = *sdk.NewNullableString(&spVal)
	}

	if !v.ControllerMountPoint.IsNull() && !v.ControllerMountPoint.IsUnknown() {
		vol.ControllerMountPoint = v.ControllerMountPoint.ValueStringPointer()
	}

	return vol
}

// resizeNetworkInterfaces builds the resize request's network interface list
// from the planned interfaces, each paired with the clone's current
// counterpart, and reports whether the pairing found a change: an interface
// added or removed, a paired interface whose network, type, IP mode, IP
// address or MAC address is set in the plan and differs, or whose child
// virtual networks differ. An existing interface keeps its current id when
// the plan does not know it, so the appliance updates it rather than
// replacing it; an interface with no counterpart is sent without an id,
// which is how a new NIC is asked for.
//
// Pairing follows resizeVolumes: positional when the counts are equal, and
// when they differ each planned interface takes the first current interface
// on the same network not yet taken, so removing the middle NIC of three
// leaves the third its id. Interfaces have no name; the network is the only
// stable key a clone's configuration carries.
func resizeNetworkInterfaces(
	ctx context.Context,
	ifaceList types.List,
	current []NetworkInterfacesValue,
) ([]sdk.InstancesNetworkInterfaces4, bool, diag.Diagnostics) {
	planIfaces, diags := listElements[NetworkInterfacesValue](ctx, ifaceList)
	if ifaceList.IsNull() || ifaceList.IsUnknown() {
		return nil, false, diags
	}

	changed := len(planIfaces) != len(current)
	sdkIfaces := make([]sdk.InstancesNetworkInterfaces4, 0, len(planIfaces))
	pairs := pairInterfaces(planIfaces, current)

	for i, ni := range planIfaces {
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

		ci := pairs[i]

		switch {
		case !ni.Id.IsNull() && !ni.Id.IsUnknown():
			iface.Id = ni.Id.ValueInt64Pointer()
		case ci >= 0 && !current[ci].Id.IsNull() && !current[ci].Id.IsUnknown():
			iface.Id = current[ci].Id.ValueInt64Pointer()
		}

		children, cd := buildChildInterfaces(ctx, ni.ChildVirtualNetworks)
		diags.Append(cd...)
		iface.NetworkInterfaces = children

		if ci >= 0 {
			ifaceChanged, d := networkInterfaceChanged(ctx, ni, current[ci])
			diags.Append(d...)

			if ifaceChanged {
				changed = true
			}
		}

		sdkIfaces = append(sdkIfaces, iface)
	}

	return sdkIfaces, changed, diags
}

// pairInterfaces returns, for each planned interface, the index of its
// current counterpart, or -1 for a new interface: positional when the counts
// are equal, by network otherwise (see resizeNetworkInterfaces).
func pairInterfaces(plan, current []NetworkInterfacesValue) []int {
	pairs := make([]int, len(plan))

	if len(plan) == len(current) {
		for i := range plan {
			pairs[i] = i
		}

		return pairs
	}

	taken := make([]bool, len(current))

	for i, ni := range plan {
		pairs[i] = -1

		if ni.NetworkId.IsNull() || ni.NetworkId.IsUnknown() {
			continue
		}

		for ci := range current {
			if !taken[ci] && current[ci].NetworkId.Equal(ni.NetworkId) {
				pairs[i] = ci
				taken[ci] = true

				break
			}
		}
	}

	return pairs
}

// networkInterfaceChanged reports whether a planned interface differs from its
// current counterpart in a way a resize acts on.
func networkInterfaceChanged(
	ctx context.Context,
	plan, current NetworkInterfacesValue,
) (bool, diag.Diagnostics) {
	if int64Differs(plan.NetworkId, current.NetworkId) ||
		int64Differs(plan.NetworkInterfaceTypeId, current.NetworkInterfaceTypeId) ||
		stringDiffers(plan.IpMode, current.IpMode) ||
		stringDiffers(plan.IpAddress, current.IpAddress) ||
		stringDiffers(plan.MacAddress, current.MacAddress) {
		return true, nil
	}

	return childVirtualNetworksChanged(ctx, plan.ChildVirtualNetworks, current.ChildVirtualNetworks)
}

// childVirtualNetworksChanged reports whether the planned child virtual
// networks of an interface differ from the current ones: a child added or
// removed, or a paired child on a different network. An unknown planned list
// is one the practitioner did not set and is not a change.
func childVirtualNetworksChanged(
	ctx context.Context,
	plan, current types.List,
) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	if plan.IsUnknown() {
		return false, diags
	}

	planChildren, d := listElements[ChildVirtualNetworksValue](ctx, plan)
	diags.Append(d...)

	currentChildren, d := listElements[ChildVirtualNetworksValue](ctx, current)
	diags.Append(d...)

	if len(planChildren) != len(currentChildren) {
		return true, diags
	}

	for i := range planChildren {
		if int64Differs(planChildren[i].NetworkId, currentChildren[i].NetworkId) {
			return true, diags
		}
	}

	return false, diags
}
