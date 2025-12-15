package instance

import (
	"context"
	"fmt"
	"net/http"
	"slices"
	"strconv"
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

	// name
	state.Name = convert.StrToType(instance.Name)

	// interfaces
	ifaces, ifDiags := getStateInterfacesFromInstance(ctx, instance)
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
	apiVolumes := slices.DeleteFunc(
		resp.GetInstance().Volumes,
		func(v sdk.AddInstance200ResponseAllOfOneOfInstanceVolumesInner) bool {
			if v.Name == nil {
				return false
			}

			if strings.HasPrefix(*v.Name, "CD ROM") {
				return true
			}

			return false
		},
	)

	volumes, d := convert.ToListType(
		ctx,
		apiVolumes,
		func(
			in sdk.AddInstance200ResponseAllOfOneOfInstanceVolumesInner,
		) VolumesValue {
			v := VolumesValue{}
			v.Id = convert.Int64ToType(in.Id)
			v.RootVolume = convert.BoolToType(in.RootVolume)
			v.Name = convert.StrToType(in.Name)
			v.Size = convert.Int64ToType(in.Size)
			v.StorageTypeId = convert.Int64ToType(in.StorageType)

			if in.DatastoreId != nil {
				datastore, err := strconv.ParseInt(in.GetDatastoreId(), 10, 64)
				if err != nil {
					v.DatastoreId = basetypes.NewInt64Unknown()
				}

				v.DatastoreId = types.Int64Value(datastore)
			}

			v.ControllerMountPoint = convert.StrToType(in.ControllerMountPoint)
			v.state = attr.ValueStateKnown

			return v
		},
	)
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

// getStateInterfacesFromInstance get the []NetworkInterfacesValue from containerDetails.server.interfaces
func getStateInterfacesFromInstance(
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
