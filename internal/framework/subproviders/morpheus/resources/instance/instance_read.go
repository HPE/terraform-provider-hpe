package instance

import (
	"context"
	"fmt"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/convert"
	errfmt "github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/errors"
)

// Read implements resource.Resource.
func (g *Resource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var data InstanceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get timeout from HCL if set, the default is 45 minutes
	createTimeout, diags := data.Timeouts.Read(ctx, 45*time.Minute)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, createTimeout)
	defer cancel()

	client, err := g.NewClient(ctx)
	if err != nil {
		resp.Diagnostics.AddError("creating client failed", err.Error())

		return
	}

	state, diag := getInstanceAsState(ctx, data.Id.ValueInt64(), client, data)
	if resp.Diagnostics.Append(diag...); resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func getInstanceAsState(
	ctx context.Context,
	id int64,
	client *sdk.APIClient,
	plan InstanceModel,
) (InstanceModel, diag.Diagnostics) {
	var state InstanceModel
	var diags diag.Diagnostics

	resp, hresp, err := client.InstancesAPI.GetInstance(ctx, id).Execute()
	if err != nil || hresp.StatusCode != http.StatusOK {
		diags.AddError(
			"populate instance resource",
			fmt.Sprintf("instance %d GET failed: ", id)+errfmt.ErrMsg(err, hresp),
		)

		return state, diags
	}

	instance := resp.GetInstance()

	// cloud_id
	state.CloudId = convert.Int64ToType(instance.Cloud.Id)

	// config
	state.Config = types.DynamicNull()

	if !plan.Config.IsNull() && !plan.Config.IsUnknown() {
		state.Config = plan.Config
	}

	// connection_info
	cInfo, dc := getConnectionInfo(instance)
	diags.Append(dc...)
	if diags.HasError() {
		return state, diags
	}

	state.ConnectionInfo = cInfo

	// evars
	// API may respond with more evars than what the user set so we need to
	// check the /instance/{id}/envs endpoint which gives us the user specified
	// evars separately.
	envVars, diags := getInstanceEnvVars(ctx, id, client)
	if diags.HasError() {
		return state, diags
	}

	evars, d := convert.ToSetType(
		ctx,
		envVars,
		func(
			in sdk.GetEnvVariables200ResponseEnvsInner,
		) EvarsValue {
			return EvarsValue{
				Name:  types.StringValue(in.Name),
				Value: types.StringValue(in.Value),
				state: attr.ValueStateKnown,
			}
		},
	)
	diags.Append(d...)
	state.Evars = evars

	// group_id
	state.GroupId = convert.Int64ToType(instance.GetGroup().Id)

	// id
	state.Id = convert.Int64ToType(instance.Id)

	// instance_context
	state.InstanceContext = convert.StrToType(instance.InstanceContext.Get())

	// instance_type_id
	state.InstanceTypeId = convert.Int64ToType(instance.InstanceType.Id)

	// instance_type_id
	state.InstanceTypeId = convert.Int64ToType(instance.InstanceType.Id)

	// layout_id
	state.LayoutId = convert.Int64ToType(instance.Layout.Id)

	// layout_size - from Config
	if config, ok := instance.GetConfigOk(); ok {
		if layoutSize, ok := config.GetLayoutSizeOk(); ok {
			state.LayoutSize = convert.Int64ToType(layoutSize)
		}
	} else if !plan.LayoutSize.IsNull() && !plan.LayoutSize.IsUnknown() {
		// fallback to instance.layoutSize
		state.LayoutSize = plan.LayoutSize
	}

	// name
	state.Name = convert.StrToType(instance.Name)

	// interfaces
	ifaces, ifDiags := getStateInterfaces(ctx, instance, plan)
	diags = append(diags, ifDiags...)
	if diags.HasError() {
		return state, diags
	}

	networkInterfacesList, d := types.ListValueFrom(ctx, NetworkInterfacesValue{}.Type(ctx), ifaces)
	diags.Append(d...)

	if diags.HasError() {
		tflog.Error(ctx, "cannot convert network interfaces")

		return state, diags
	}

	state.NetworkInterfaces = networkInterfacesList

	// plan_id
	state.PlanId = convert.Int64ToType(instance.Plan.Id)

	// ports
	// assume ports always match the plan. Ports are a bit complicated in the
	// API. They are a part of container_details, which is an array which may
	// be of arbitrary size. The ports are contained within every element of the
	// array. The array of ports also may contain non-user controlled ports.
	// This should probably be replaced by a plan modifier.
	state.Ports = plan.Ports

	// timeouts
	state.Timeouts = plan.Timeouts

	// tags
	tags, d := convert.ToSetType(
		ctx,
		resp.GetInstance().Tags,
		func(
			in sdk.AddInstance200ResponseAllOfOneOfInstanceTagsInner,
		) TagsValue {
			return TagsValue{
				Name:  convert.StrToType(in.Name),
				Value: convert.StrToType(in.Value),
				state: attr.ValueStateKnown,
			}
		},
	)
	diags.Append(d...)
	state.Tags = tags

	// task_set_id
	// task_set_id is not included in the API response, it is a write only value
	// we might need modify the schema to add some plan modifiers to it
	state.TaskSetId = plan.TaskSetId

	// volumes
	volumes, d := getVolumes(ctx, instance, plan)
	diags.Append(d...)
	state.Volumes = volumes

	return state, diags
}

func getInstanceEnvVars(
	ctx context.Context,
	id int64,
	client *sdk.APIClient,
) ([]sdk.GetEnvVariables200ResponseEnvsInner, diag.Diagnostics) {
	var diags diag.Diagnostics
	resp, _, _ := client.InstancesAPI.GetEnvVariables(ctx, id).Execute()
	// ignoring errors for now, the sdk can't parse some of the unused fields
	// due to polymorphic values

	// if err != nil || hresp.StatusCode != http.StatusOK {
	// diags.AddError(
	// 	"populate instance resource",
	// 	fmt.Sprintf("instance %d GET failed: ", id)+errfmt.ErrMsg(err, hresp),
	// )

	// return nil, diags
	// }

	return resp.GetEnvs(), diags
}

// getVolumes builds the volumes list from instance.containerDetails.server.volumes
func getVolumes(
	ctx context.Context,
	instance sdk.AddInstance200ResponseAllOfOneOfInstance,
	plan InstanceModel,
) (basetypes.ListValue, diag.Diagnostics) {
	diags := diag.Diagnostics{}

	// Get volumes from instance.containerDetails.server.volumes
	contDetails, ok := instance.GetContainerDetailsOk()
	if !ok || len(contDetails) == 0 {
		diags.AddError(
			"cannot get instance containerDetails",
			fmt.Sprintf("instance %d GET containerDetails failed", instance.GetId()))

		return basetypes.NewListNull(VolumesValue{}.Type(ctx)), diags
	}

	server, ok := contDetails[0].GetServerOk()
	if !ok {
		diags.AddError(
			"cannot get instance containerDetails server",
			fmt.Sprintf("instance %d GET containerDetails.server failed", instance.GetId()))

		return basetypes.NewListNull(VolumesValue{}.Type(ctx)), diags
	}

	serverVolumes, ok := server.GetVolumesOk()
	if !ok {
		diags.AddError(
			"cannot get instance containerDetails server volumes",
			fmt.Sprintf("instance %d GET containerDetails.server.volumes failed", instance.GetId()))

		return basetypes.NewListNull(VolumesValue{}.Type(ctx)), diags
	}

	// Remove any CD ROM volumes from the list
	apiVolumes := slices.DeleteFunc(
		serverVolumes,
		func(v sdk.InstanceContainerServerVolume1) bool {
			if v.Name == nil {
				return false
			}

			if strings.HasPrefix(*v.Name, "CD ROM") {
				return true
			}

			return false
		},
	)

	// Import
	if plan.Name.IsNull() || plan.Name.IsUnknown() {
		nonRaidVolumes := removeRaidDisks(apiVolumes)
		bootVolumeInFirst := bootVolumeFirst(nonRaidVolumes)

		return convertAPIVolumesToStateVolumes(ctx, bootVolumeInFirst)
	}

	// If the number of volumes is the same as the plan, we can do a direct conversion
	if len(apiVolumes) == len(plan.Volumes.Elements()) {
		autoselectVolumes := setDatastoreAutoSelectionAndSize(apiVolumes, plan)

		return convertAPIVolumesToStateVolumes(ctx, autoselectVolumes)
	}

	// The number of volumes is different to the plan
	nonRaidVolumes := removeRaidDisks(apiVolumes)
	reorderedVolumes := reorderVolumes(nonRaidVolumes, plan)
	filledVolumes := fillVolumeFieldsFromPlan(reorderedVolumes, plan)
	autoselectVolumes := setDatastoreAutoSelectionAndSize(filledVolumes, plan)

	return convertAPIVolumesToStateVolumes(ctx, autoselectVolumes)
}

// setDatastoreAutoSelectionAndSize sets the AdditionalProperties field for volumes with
// DatastoreAutoSelection set to that in the plan
// Also sets the MaxStorage field from the plan size
// Handles cases where the number of API volumes differs from plan volumes
func setDatastoreAutoSelectionAndSize(
	apiVolumes []sdk.InstanceContainerServerVolume1,
	plan InstanceModel,
) []sdk.InstanceContainerServerVolume1 {
	autoSelection := make([]sdk.InstanceContainerServerVolume1, 0, len(apiVolumes))
	planVolumes := plan.Volumes.Elements()

	// Determine how many volumes we can safely access from the plan
	maxIndex := len(apiVolumes)
	if len(planVolumes) < maxIndex {
		maxIndex = len(planVolumes)
	}

	for i, apiVol := range apiVolumes {
		// Only set fields from plan if we have a corresponding plan volume
		if i < maxIndex {
			planVol := planVolumes[i].(VolumesValue)

			// Initialize AdditionalProperties map if it doesn't exist
			if apiVol.AdditionalProperties == nil {
				apiVol.AdditionalProperties = make(map[string]interface{})
			}

			if !planVol.DatastoreAutoSelection.IsNull() && !planVol.DatastoreAutoSelection.IsUnknown() {
				// Set AdditionalProperties to indicate auto-selection
				apiVol.AdditionalProperties["DatastoreAutoSelection"] = planVol.DatastoreAutoSelection.ValueString()
			}

			apiVol.MaxStorage = planVol.Size.ValueInt64Pointer()
			// We set this flag to indicate that Terraform set the MaxStorage value
			apiVol.AdditionalProperties["TerraformSetMaxStorage"] = true
		}
		// If i >= maxIndex, just append the apiVol as-is (unmatched volumes)

		autoSelection = append(autoSelection, apiVol)
	}

	return autoSelection
}

// fillVolumeFieldsFromPlan fills in some missing fields in the API volumes from the plan
// This is needed because some fields are not returned by the API
// for certain volume types (e.g. Metal RAID volumes)
// We assume that the volumes are in the same order as the plan, but handle cases where
// the number of volumes differs
func fillVolumeFieldsFromPlan(
	apiVolumes []sdk.InstanceContainerServerVolume1,
	plan InstanceModel,
) []sdk.InstanceContainerServerVolume1 {
	// Now fill in any missing fields from the plan
	filledVolumes := make([]sdk.InstanceContainerServerVolume1, 0, len(apiVolumes))
	planVolumes := plan.Volumes.Elements()

	// Determine how many volumes we can safely fill from the plan
	maxIndex := len(apiVolumes)
	if len(planVolumes) < maxIndex {
		maxIndex = len(planVolumes)
	}

	for i, apiVol := range apiVolumes {
		// Only fill from plan if we have a corresponding plan volume
		if i < maxIndex {
			planVol := planVolumes[i].(VolumesValue)
			if apiVol.DatastoreId == nil {
				apiVol.DatastoreId = planVol.DatastoreId.ValueInt64Pointer()
			}
			// We set TypeId for all volume types since for some types (e.g. Metal RAID) the API TypeId is different
			apiVol.TypeId = planVol.StorageTypeId.ValueInt64Pointer()
		}
		// If i >= maxIndex, just append the apiVol as-is (unmatched volumes)

		filledVolumes = append(filledVolumes, apiVol)
	}

	return filledVolumes
}

// bootVolumeFirst puts the boot volume first in the list of volumes
// We hope to remove this function in future when the API names RAID volumes correctly
func bootVolumeFirst(
	apiVolumes []sdk.InstanceContainerServerVolume1,
) []sdk.InstanceContainerServerVolume1 {
	// If there are no volumes, return the original list
	if len(apiVolumes) == 0 {
		return apiVolumes
	}

	// Put boot volume first
	bootVolumes := make([]sdk.InstanceContainerServerVolume1, 0, len(apiVolumes))
	otherVolumes := make([]sdk.InstanceContainerServerVolume1, 0, len(apiVolumes))
	for _, v := range apiVolumes {
		if v.RootVolume != nil && *v.RootVolume {
			bootVolumes = append(bootVolumes, v)
		} else {
			otherVolumes = append(otherVolumes, v)
		}
	}

	// Combine boot volumes first, then other volumes
	result := append(bootVolumes, otherVolumes...)

	return result
}

// reorderVolumes re-orders the list of volumes to match the plan
func reorderVolumes(
	apiVolumes []sdk.InstanceContainerServerVolume1,
	plan InstanceModel,
) []sdk.InstanceContainerServerVolume1 {
	// If there are no volumes, return empty list
	if len(apiVolumes) == 0 {
		return apiVolumes
	}

	// Track which API volumes have been matched to avoid duplicates
	matchedVolumes := make(map[int]bool)

	// Now re-order to match plan
	orderedVolumes := make([]sdk.InstanceContainerServerVolume1, 0, len(apiVolumes))

	// Match remaining volumes with plan volumes by name
	for _, planVol := range plan.Volumes.Elements() {
		planVolTyped := planVol.(VolumesValue)
		for i, apiVol := range apiVolumes {
			if matchedVolumes[i] {
				continue // Skip already matched volumes
			}

			if apiVol.Name != nil && planVolTyped.Name.ValueString() == *apiVol.Name {
				orderedVolumes = append(orderedVolumes, apiVol)
				matchedVolumes[i] = true

				break
			}
		}
	}

	// Append any unmatched volumes at the end to avoid data loss
	for i, apiVol := range apiVolumes {
		if !matchedVolumes[i] {
			orderedVolumes = append(orderedVolumes, apiVol)
		}
	}

	return orderedVolumes
}

// removeRaidDisks removes any RAID disks from the list of volumes
func removeRaidDisks(
	apiVolumes []sdk.InstanceContainerServerVolume1,
) []sdk.InstanceContainerServerVolume1 {
	// build a map of device-counts for the API volumes
	deviceCount := make(map[string]int)
	for _, volume := range apiVolumes {
		if volume.DeviceName != nil {
			deviceCount[*volume.DeviceName]++
		}
	}

	// remove the RAID disks from the apiVolumes list
	nonRaidDiskVolumes := slices.DeleteFunc(
		apiVolumes,
		func(v sdk.InstanceContainerServerVolume1) bool {
			// Skip volumes without a device name
			if v.DeviceName == nil {
				return false
			}

			if deviceCount[*v.DeviceName] > 1 {
				// We're going to remove volumes which have a diskType
				if v.DiskType != nil && v.DiskMode == nil {
					return true
				}
			}

			return false
		},
	)

	return nonRaidDiskVolumes
}

// convertBytesPtrToGBBytes converts a pointer to int64 bytes to a pointer to int64 GB bytes
func convertBytesPtrToGBBytes(b *int64) *int64 {
	if b == nil {
		return nil
	}
	gb := *b / (1 << 30)

	return &gb
}

func convertAPIVolumesToStateVolumes(
	ctx context.Context,
	apiVolumes []sdk.InstanceContainerServerVolume1,
) (basetypes.ListValue, diag.Diagnostics) {
	volumes, d := convert.ToListType(
		ctx,
		apiVolumes,
		func(
			in sdk.InstanceContainerServerVolume1,
		) VolumesValue {
			v := VolumesValue{}
			v.Id = convert.Int64ToType(in.Id)
			v.RootVolume = convert.BoolToType(in.RootVolume)
			v.Name = convert.StrToType(in.Name)
			v.StorageTypeId = convert.Int64ToType(in.TypeId)
			v.DatastoreId = convert.Int64ToType(in.DatastoreId)
			v.ControllerMountPoint = convert.StrToType(in.ControllerMountPoint)

			// Handle DatastoreAutoSelection and TerraformSetMaxStorage from AdditionalProperties
			// TerraformSetMaxStorage flag indicates that MaxStorage was set from plan (already in GB)
			// and should not be converted from bytes
			terraformSetMaxStorage := false
			if in.AdditionalProperties != nil {
				if dsAutoSel, ok := in.AdditionalProperties["DatastoreAutoSelection"]; ok {
					if dsAutoSelStr, ok := dsAutoSel.(string); ok {
						v.DatastoreAutoSelection = convert.StrToType(&dsAutoSelStr)
					}
				}

				if tsms, ok := in.AdditionalProperties["TerraformSetMaxStorage"]; ok {
					if tsmsBool, ok := tsms.(bool); ok {
						terraformSetMaxStorage = tsmsBool
					}
				}
			}

			// Set Size: if TerraformSetMaxStorage is true, MaxStorage is already in GB from plan
			// Otherwise, convert from bytes to GB
			if terraformSetMaxStorage {
				v.Size = convert.Int64ToType(in.MaxStorage)
			} else {
				v.Size = convert.Int64ToType(convertBytesPtrToGBBytes(in.MaxStorage))
			}

			v.state = attr.ValueStateKnown

			return v
		},
	)

	return volumes, d
}

// getConnectionInfo builds the connection_info list
func getConnectionInfo(
	instance sdk.AddInstance200ResponseAllOfOneOfInstance,
) (types.List, diag.Diagnostics) {
	diags := diag.Diagnostics{}
	cInfo, ok := instance.GetConnectionInfoOk()
	if !ok {
		diags.AddError(
			"cannot get instance connectionInfo",
			fmt.Sprintf("instance %d GET connectionInfo failed", instance.GetId()))

		return types.ListNull(types.StringType), diags
	}

	if len(cInfo) == 0 {
		return types.ListNull(types.StringType), diags
	}

	var vals []attr.Value
	for _, c := range cInfo {
		if ip, ok := c.GetIpOk(); ok {
			vals = append(vals, types.StringValue(*ip))
		}
	}

	cList, dl := types.ListValue(types.StringType, vals)
	diags = append(diags, dl...)

	return cList, diags
}

// getStateInterfaces get the interfaces to be returned as state entries
func getStateInterfaces(
	ctx context.Context,
	instance sdk.AddInstance200ResponseAllOfOneOfInstance,
	plan InstanceModel,
) ([]NetworkInterfacesValue, diag.Diagnostics) {
	diags := diag.Diagnostics{}

	// Generate []NetworkInterfacesValue from instance.interfaces
	intfsFromInstance, id := getStateInterfacesFromInstance(ctx, instance)
	diags = append(diags, id...)
	if diags.HasError() {
		return nil, diags
	}

	// If this is an import, then return intfsFromInstance
	if plan.Name.IsNull() || plan.Name.IsUnknown() {
		return intfsFromInstance, diags
	}

	// Generate []NetworkInterfacesValue from the instance.containerDetails.server.interfaces
	intfsFromServer, id := getStateInterfacesFromInstanceServer(ctx, instance)
	diags = append(diags, id...)
	if diags.HasError() {
		return nil, diags
	}

	// Compare intfsFromServer against intfsFromIntance, to see if the "shapes" are the same
	if compareServerInstanceIntfs(intfsFromServer, intfsFromInstance) {
		return intfsFromServer, diags
	}

	// "Shape" isn't the same, return intfsFromInstance
	return intfsFromInstance, diags
}

// compareServerInstanceIntfs compares the []NetworkInterfacesValues from instance.containerDetails.server.interfaces
// and instance.interfaces to see if they are the same shape
// Returns true if they are, false otherwise
func compareServerInstanceIntfs(
	intfsFromServer, intfsFromInstance []NetworkInterfacesValue,
) bool {
	// Check length of lists first
	if len(intfsFromServer) != len(intfsFromInstance) {
		return false
	}

	// Get list of lengths of child interfaces for instance.containerDetails.server.interfaces list
	serverSubIntfs := make([]int, 0, len(intfsFromServer))
	for _, serverIntf := range intfsFromServer {
		serverSubIntfs = append(serverSubIntfs, len(serverIntf.ChildVirtualNetworks.Elements()))
	}

	// Get list of lengths of child interfaces for instance.interfaces list
	instanceSubIntfs := make([]int, 0, len(intfsFromInstance))
	for _, instanceIntf := range intfsFromInstance {
		instanceSubIntfs = append(instanceSubIntfs, len(instanceIntf.ChildVirtualNetworks.Elements()))
	}

	// Compare lengths of child interfaces for "server" and "instance" lists
	for i := range serverSubIntfs {
		if serverSubIntfs[i] != instanceSubIntfs[i] {
			return false
		}
	}

	return true
}

// getStateInterfacesFromInstance build []NetworkInterfacesValue from interfaces, used on import
func getStateInterfacesFromInstance(
	ctx context.Context,
	instance sdk.AddInstance200ResponseAllOfOneOfInstance,
) ([]NetworkInterfacesValue, diag.Diagnostics) {
	diags := diag.Diagnostics{}

	instIntfs, ok := instance.GetInterfacesOk()
	if !ok {
		diags.AddError(
			"instance GetInterfaces failed",
			fmt.Sprintf("instance %d GET interfaces failed", instance.GetId()))

		return nil, diags
	}

	var ifaces []NetworkInterfacesValue
	for _, instIntf := range instIntfs {
		ifaceVal := NetworkInterfacesValue{}
		ifaceVal.IpAddress = convert.StrToType(instIntf.IpAddress)
		ifaceVal.IpMode = convert.StrToType(instIntf.IpMode)
		ifaceVal.PrimaryInterface = types.BoolNull()
		ifaceVal.Name = types.StringNull()
		ifaceVal.NetworkId = types.Int64Null()
		if net, ok := instIntf.GetNetworkOk(); ok {
			ifaceVal.NetworkId = convert.Int64ToType(net.Id)
			ifaceVal.IpPool = types.Int64Null()
			if pool, ok := net.GetPoolOk(); ok {
				ifaceVal.IpPool = convert.Int64ToType(pool.Id)
			}
			ifaceVal.NetworkGroupId = convert.Int64ToType(net.Group)
		}
		ifaceVal.NetworkTypeId = networkTypeId(instIntf.NetworkInterfaceTypeId)
		ifaceVal.ChildVirtualNetworks = types.ListNull(ChildVirtualNetworksValue{}.Type(ctx))
		if cnets, ok := instIntf.GetNetworkInterfacesOk(); ok {
			if len(cnets) > 0 {
				childNetworks, cd := getInstanceInterfacesChildNetworks(ctx, cnets)
				ifaceVal.ChildVirtualNetworks = childNetworks
				diags = append(diags, cd...)
			}
		}
		ifaceVal.state = attr.ValueStateKnown

		ifaces = append(ifaces, ifaceVal)
	}

	return ifaces, diags
}

// networkTypeId helper function to handle NetworkInterfaceTypeId
// We only use this function on import, and it looks as if NetworkInterfaceTypeId will have a value of 0 instead
// of null.  In this case we return a null value to avoid inadvertent instance recreations when the HCL has been
// generated on import and a plan/apply following the initial import is performed
func networkTypeId(i *int64) basetypes.Int64Value {
	if i != nil && *i == 0 {
		return basetypes.NewInt64Null()
	}

	return convert.Int64ToType(i)
}

// getInstanceInterfacesChildNetworks returns child_networks from interfaces.networkInterfaces, used on import
func getInstanceInterfacesChildNetworks(
	ctx context.Context,
	nets []sdk.InstanceInterfacesNetworkInterfacesInner1,
) (basetypes.ListValue, diag.Diagnostics) {
	children := make([]ChildVirtualNetworksValue, 0)
	for _, instIntf := range nets {
		ifaceVal := ChildVirtualNetworksValue{}
		ifaceVal.IpAddress = convert.StrToType(instIntf.IpAddress)
		ifaceVal.IpMode = convert.StrToType(instIntf.IpMode)
		ifaceVal.PrimaryInterface = types.BoolNull()
		ifaceVal.Name = types.StringNull()
		ifaceVal.NetworkId = types.Int64Null()
		if net, ok := instIntf.GetNetworkOk(); ok {
			ifaceVal.NetworkId = convert.Int64ToType(net.Id)
			ifaceVal.IpPool = types.Int64Null()
			if pool, ok := net.GetPoolOk(); ok {
				ifaceVal.IpPool = convert.Int64ToType(pool.Id)
			}
			ifaceVal.NetworkGroupId = convert.Int64ToType(net.Group)
		}
		ifaceVal.NetworkTypeId = networkTypeId(instIntf.NetworkInterfaceTypeId)

		ifaceVal.state = attr.ValueStateKnown
		children = append(children, ifaceVal)
	}

	return types.ListValueFrom(ctx, ChildVirtualNetworksValue{}.Type(ctx), children)
}

// getStateInterfacesFromInstanceServer get the []NetworkInterfacesValue from containerDetails.server.interfaces
func getStateInterfacesFromInstanceServer(
	ctx context.Context,
	instance sdk.AddInstance200ResponseAllOfOneOfInstance,
) ([]NetworkInterfacesValue, diag.Diagnostics) {
	// network_interfaces
	// We are going to read network interface information from containerDetails.server.interfaces
	// Note that, at present, all network IP addresses will not be available to us when we reach
	// this stage on instance creation, we will have enough in the state-file that a plan will
	// be a no-op, and that when all IP addresses are available (this can be seen in the UI) an
	// apply will update the state-file with the IP addresses etc.
	procIntfs := getAllServerInterfaces(instance)

	var ifaces []NetworkInterfacesValue
	var childInterfaces basetypes.ListValue
	var diags diag.Diagnostics

	for _, iface := range procIntfs.serverIntfsList {
		// Skip sub-interfaces
		if _, ok := procIntfs.isSubIntf[*iface.Id]; ok {
			continue
		}
		ifaceVal := NetworkInterfacesValue{}

		ifaceVal.IpAddress = convert.StrToType(iface.IpAddress)
		ifaceVal.IpMode = convert.StrToType(iface.IpMode)
		ifaceVal.NetworkGroupId = types.Int64Null()
		if group, ok := iface.GetNetworkGroupOk(); ok {
			ifaceVal.NetworkGroupId = convert.Int64ToType(group.Id)
		}
		if pool, ok := iface.GetNetworkPoolOk(); ok {
			ifaceVal.IpPool = convert.Int64ToType(pool.Id)
		}
		if net, ok := iface.GetNetworkOk(); ok {
			ifaceVal.NetworkId = convert.Int64ToType(net.Id)
		}
		ifaceVal.Name = convert.StrToType(iface.Name)
		ifaceVal.PrimaryInterface = convert.BoolToType(iface.PrimaryInterface)

		childInterfaces, diags = getChildNetworks(ctx, iface.Id, procIntfs.subIntfsMap, procIntfs.serverIntfsMap)

		ifaceVal.ChildVirtualNetworks = childInterfaces

		ifaceVal.state = attr.ValueStateKnown

		ifaces = append(ifaces, ifaceVal)
	}

	return ifaces, diags
}

// processedServerInterfaces struct that contains maps and a list produced from containerDetails.server.interfaces
type processedServerInterfaces struct {
	// subIntfsMap a map of interface-ids with a list of the ids of any sub-interfaces
	// note that it is possible for two (or maybe more) interfaces to have the same sub-interfaces (bonds)
	subIntfsMap map[int64][]int64
	// isSubIntf a map of interface-ids with a boolean saying if they are sub-interfaces
	isSubIntf map[int64]bool
	// serverIntfsMap a map of interface-ids with the corresponding interface information
	serverIntfsMap map[int64]sdk.InstanceContainerServerInterfacesInner1
	// serverIntfsList a list of the interfaces, which should (hopefully be in the same order as those specified
	// in network_interfaces
	serverIntfsList []sdk.InstanceContainerServerInterfacesInner1
}

// Process the set of interfaces in an instance
// This function takes an "instance" input and returns processedServerInterfaces
func getAllServerInterfaces(
	instance sdk.AddInstance200ResponseAllOfOneOfInstance,
) processedServerInterfaces {
	subIntfsMap := make(map[int64][]int64)
	isSubIntf := make(map[int64]bool)
	serverIntfsMap := make(map[int64]sdk.InstanceContainerServerInterfacesInner1)
	serverIntfsList := make([]sdk.InstanceContainerServerInterfacesInner1, 0)

	// First clean-up the server interface list
	// The list of interfaces is malleable, and changes after the instance has been created and returned
	// to us for reading.  In early stages the interface "name" (eth0, eth1 etc) will have repeated
	// entries.
	// Key here is the "UniqueId".  If an interface doesn't have a value for that or for Network then all
	// it has is an IP address that will be assigned to the interface with the same name (eth0, eth1, etc)
	for _, container := range instance.ContainerDetails {
		server, _ := container.GetServerOk()
		serverIntfList, _ := server.GetInterfacesOk()
		serverIntfsNameMap := make(map[string][]sdk.InstanceContainerServerInterfacesInner1)
		serverIntfsNameListPosition := make([]string, 0)
		serverIntfsNameListMap := make(map[string]struct{})
		serverIntfsMergedNameMap := make(map[string]sdk.InstanceContainerServerInterfacesInner1)
		for _, serverIntf := range serverIntfList {
			// Skip this list entry if it doesn't have a name
			if _, ok := serverIntf.GetNameOk(); !ok {
				continue
			}

			serverIntfsNameMap[*serverIntf.Name] = append(serverIntfsNameMap[*serverIntf.Name], serverIntf)
			// Keep a record of the order of the interface name ("eth0" etc) in the input list
			// We allow for duplicate name entries, so we're looking for the first entry
			if _, ok := serverIntfsNameListMap[*serverIntf.Name]; !ok {
				serverIntfsNameListPosition = append(serverIntfsNameListPosition, *serverIntf.Name)
				serverIntfsNameListMap[*serverIntf.Name] = struct{}{}
			}
		}

		for intfName, v := range serverIntfsNameMap {
			if len(v) == 1 {
				serverIntfsMergedNameMap[intfName] = v[0]

				continue
			}

			// Hopefully there will only be two entries in the other lists
			// What we are going to do here is use the entry that has Network information
			// as the base of the cumulative interface, and then hunt through the rest for
			// an ip-address
			var cumulativeIntf sdk.InstanceContainerServerInterfacesInner1
			var ipAddress *string
			for _, serverIntf := range v {
				if _, ok := serverIntf.GetNetworkOk(); ok {
					cumulativeIntf = serverIntf

					break
				}
			}
			for _, serverIntf := range v {
				if ip, ok := serverIntf.GetIpAddressOk(); ok {
					ipAddress = ip

					break
				}
			}

			cumulativeIntf.IpAddress = ipAddress
			serverIntfsMergedNameMap[intfName] = cumulativeIntf
		}

		// Order serverIntfsList by order that the entries appeared in the original list
		for _, intfName := range serverIntfsNameListPosition {
			serverIntfsList = append(serverIntfsList, serverIntfsMergedNameMap[intfName])
		}
	}

	// Build the maps that we're going to return
	for _, serverInterface := range serverIntfsList {
		serverIntfsMap[serverInterface.GetId()] = serverInterface
		if subIntfs, ok := serverInterface.GetInterfacesOk(); ok {
			intfList := make([]int64, 0)
			for _, subIntf := range subIntfs {
				intfList = append(intfList, subIntf.GetId())
				isSubIntf[subIntf.GetId()] = true
			}
			if len(intfList) > 0 {
				subIntfsMap[serverInterface.GetId()] = intfList
			}
		}
	}

	ret := processedServerInterfaces{}
	ret.subIntfsMap = subIntfsMap
	ret.isSubIntf = isSubIntf
	ret.serverIntfsMap = serverIntfsMap
	ret.serverIntfsList = serverIntfsList

	return ret
}

// Get the child virtual network interface values
func getChildNetworks(
	ctx context.Context,
	id *int64,
	subIntfMap map[int64][]int64,
	serverIntfsMap map[int64]sdk.InstanceContainerServerInterfacesInner1,
) (basetypes.ListValue, diag.Diagnostics) {
	if id == nil {
		return types.ListNull(ChildVirtualNetworksValue{}.Type(ctx)), nil
	}

	if len(subIntfMap[*id]) == 0 {
		return types.ListNull(ChildVirtualNetworksValue{}.Type(ctx)), nil
	}

	children := make([]ChildVirtualNetworksValue, 0)
	for _, subIntf := range subIntfMap[*id] {
		ifaceVal := ChildVirtualNetworksValue{}
		iface := serverIntfsMap[subIntf]
		ifaceVal.IpAddress = convert.StrToType(iface.IpAddress)
		ifaceVal.IpMode = convert.StrToType(iface.IpMode)
		ifaceVal.NetworkGroupId = types.Int64Null()
		if group, ok := iface.GetNetworkGroupOk(); ok {
			ifaceVal.NetworkGroupId = convert.Int64ToType(group.Id)
		}
		if pool, ok := iface.GetNetworkPoolOk(); ok {
			ifaceVal.IpPool = convert.Int64ToType(pool.Id)
		}
		ifaceVal.NetworkId = convert.Int64ToType(iface.Network.Id)
		ifaceVal.Name = convert.StrToType(iface.Name)
		ifaceVal.PrimaryInterface = convert.BoolToType(iface.PrimaryInterface)
		ifaceVal.state = attr.ValueStateKnown
		children = append(children, ifaceVal)
	}

	return types.ListValueFrom(ctx, ChildVirtualNetworksValue{}.Type(ctx), children)
}
