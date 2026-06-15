// (C) Copyright 2025-2026 Hewlett Packard Enterprise Development LP

package instance

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"

	"github.com/HPE/terraform-provider-hpe/utils/convert"
)

// Map Terraform volume value into an API request struct
func createVolumeMapper(
	vol VolumesValue,
) sdk.AddInstanceRequestVolumesInner {
	volume := sdk.AddInstanceRequestVolumesInner{}
	if !vol.Id.IsNull() && !vol.Id.IsUnknown() {
		volume.Id = vol.Id.ValueInt64Pointer()
	} else {
		volume.Id = sdk.PtrInt64(-1)
	}

	if !vol.Name.IsNull() && !vol.Name.IsUnknown() {
		volume.Name = vol.Name.ValueStringPointer()
	}

	if !vol.RootVolume.IsNull() && !vol.RootVolume.IsUnknown() {
		volume.RootVolume = vol.RootVolume.ValueBoolPointer()
	}

	if !vol.Size.IsNull() && !vol.Size.IsUnknown() {
		volume.Size = vol.Size.ValueInt64Pointer()
	}

	if !vol.StorageTypeId.IsNull() && !vol.StorageTypeId.IsUnknown() {
		volume.StorageType.Set(vol.StorageTypeId.ValueInt64Pointer())
	}

	if !vol.DatastoreId.IsNull() && !vol.DatastoreId.IsUnknown() {
		volume.DatastoreId = &sdk.
			AddInstanceRequestVolumesInnerDatastoreId{}

		id := strconv.Itoa(int(vol.DatastoreId.ValueInt64()))
		volume.DatastoreId.String = &id
	}

	if !vol.DatastoreAutoSelection.IsNull() && !vol.DatastoreAutoSelection.IsUnknown() {
		volume.DatastoreId = &sdk.
			AddInstanceRequestVolumesInnerDatastoreId{}
		volume.DatastoreId.String = vol.DatastoreAutoSelection.ValueStringPointer()
	}

	return volume
}

// Map Terraform volume value into an API request struct
func updateVolumeMapper(
	vol VolumesValue,
) sdk.ResizeInstanceRequestVolumesInner {
	volume := sdk.ResizeInstanceRequestVolumesInner{}
	if !vol.Id.IsNull() && !vol.Id.IsUnknown() {
		volume.Id = vol.Id.ValueInt64Pointer()
	} else {
		volume.Id = sdk.PtrInt64(-1)
	}

	if !vol.Name.IsNull() && !vol.Name.IsUnknown() {
		volume.Name = vol.Name.ValueStringPointer()
	}

	if !vol.RootVolume.IsNull() && !vol.RootVolume.IsUnknown() {
		volume.RootVolume = vol.RootVolume.ValueBoolPointer()
	}

	if !vol.Size.IsNull() && !vol.Size.IsUnknown() {
		volume.Size = vol.Size.ValueInt64Pointer()
	}

	if !vol.StorageTypeId.IsNull() && !vol.StorageTypeId.IsUnknown() {
		volume.StorageType.Set(vol.StorageTypeId.ValueInt64Pointer())
	}

	if !vol.DatastoreId.IsNull() && !vol.DatastoreId.IsUnknown() {
		volume.DatastoreId = &sdk.
			ResizeInstanceRequestVolumesInnerDatastoreId{}

		id := strconv.Itoa(int(vol.DatastoreId.ValueInt64()))
		volume.DatastoreId.String = &id
	}

	if !vol.DatastoreAutoSelection.IsNull() && !vol.DatastoreAutoSelection.IsUnknown() {
		volume.DatastoreId = &sdk.
			ResizeInstanceRequestVolumesInnerDatastoreId{}
		volume.DatastoreId.String = vol.DatastoreAutoSelection.ValueStringPointer()
	}

	return volume
}

// Map Terraform network interface value into an API request struct
func createNetworkInterfaceMapper(
	ctx context.Context,
) func(in NetworkInterfacesValue) sdk.InstancesNetworkInterfaces2 {
	return func(in NetworkInterfacesValue) sdk.InstancesNetworkInterfaces2 {
		var id string
		if !in.NetworkGroupId.IsNull() && !in.NetworkGroupId.IsUnknown() {
			id = "networkGroup-" + strconv.FormatInt(in.NetworkGroupId.ValueInt64(), 10)
		}

		if !in.NetworkId.IsNull() && !in.NetworkId.IsUnknown() {
			id = strconv.FormatInt(in.NetworkId.ValueInt64(), 10)
		}

		if !in.SubnetId.IsNull() && !in.SubnetId.IsUnknown() {
			id = "subnet-" + strconv.FormatInt(in.SubnetId.ValueInt64(), 10)
		}

		ipPool := &sdk.InstancesNetworkInterfaces2NetworkPool{}
		if !in.IpPool.IsNull() {
			ipPool.Id = in.IpPool.ValueInt64Pointer()
		}

		childNetworkInterfaces, diags := convert.FromListType(
			ctx,
			in.ChildVirtualNetworks,
			createChildNetworkInterfaceMapper,
		)
		if diags.HasError() {
			tflog.Error(ctx, "cannot convert child virtual network interfaces")
		}

		return sdk.InstancesNetworkInterfaces2{
			Network: sdk.
				InstancesNetworkInterfaces2Network{
				Id:   id,
				Pool: ipPool,
			},
			IpMode:                 in.IpMode.ValueStringPointer(),
			IpAddress:              in.IpAddress.ValueStringPointer(),
			NetworkInterfaceTypeId: in.NetworkTypeId.ValueInt64Pointer(),
			NetworkInterfaces:      childNetworkInterfaces,
		}
	}
}

// Map Terraform network interface value into an API request struct
func updateNetworkInterfaceMapper(
	ctx context.Context,
) func(in NetworkInterfacesValue) sdk.InstancesNetworkInterfaces3 {
	return func(in NetworkInterfacesValue) sdk.InstancesNetworkInterfaces3 {
		var id string
		if !in.NetworkGroupId.IsNull() && !in.NetworkGroupId.IsUnknown() {
			id = "networkGroup-" + strconv.FormatInt(in.NetworkGroupId.ValueInt64(), 10)
		}

		if !in.NetworkId.IsNull() && !in.NetworkId.IsUnknown() {
			id = strconv.FormatInt(in.NetworkId.ValueInt64(), 10)
		}

		if !in.SubnetId.IsNull() && !in.SubnetId.IsUnknown() {
			id = "subnet-" + strconv.FormatInt(in.SubnetId.ValueInt64(), 10)
		}

		ipPool := &sdk.InstancesNetworkInterfaces3NetworkPool{}
		if !in.IpPool.IsNull() {
			ipPool.Id = in.IpPool.ValueInt64Pointer()
		}

		var intfIdPtr *int64
		if !in.Id.IsNull() && !in.Id.IsUnknown() {
			intfIdPtr = in.Id.ValueInt64Pointer()
		}

		childNetworkInterfaces, diags := convert.FromListType(
			ctx,
			in.ChildVirtualNetworks,
			updateChildNetworkInterfaceMapper,
		)
		if diags.HasError() {
			tflog.Error(ctx, "cannot convert child virtual network interfaces")
		}

		return sdk.InstancesNetworkInterfaces3{
			Network: sdk.
				InstancesNetworkInterfaces3Network{
				Id:   id,
				Pool: ipPool,
			},
			Id:                     intfIdPtr,
			IpMode:                 in.IpMode.ValueStringPointer(),
			IpAddress:              in.IpAddress.ValueStringPointer(),
			NetworkInterfaceTypeId: in.NetworkTypeId.ValueInt64Pointer(),
			NetworkInterfaces:      childNetworkInterfaces,
		}
	}
}

// Map Child Virtual Network interface if it exists
func createChildNetworkInterfaceMapper(
	in ChildVirtualNetworksValue,
) sdk.InstancesNetworkInterfaces2NetworkInterfacesInner {
	var id string
	if !in.NetworkGroupId.IsNull() && !in.NetworkGroupId.IsUnknown() {
		id = "networkGroup-" + strconv.FormatInt(in.NetworkGroupId.ValueInt64(), 10)
	}

	if !in.NetworkId.IsNull() && !in.NetworkId.IsUnknown() {
		id = strconv.FormatInt(in.NetworkId.ValueInt64(), 10)
	}

	if !in.SubnetId.IsNull() && !in.SubnetId.IsUnknown() {
		id = "subnet-" + strconv.FormatInt(in.SubnetId.ValueInt64(), 10)
	}

	ipPool := &sdk.InstancesNetworkInterfaces2NetworkInterfacesInnerNetworkPool{}
	if !in.IpPool.IsNull() {
		ipPool.Id = in.IpPool.ValueInt64Pointer()
	}

	return sdk.InstancesNetworkInterfaces2NetworkInterfacesInner{
		Network: sdk.InstancesNetworkInterfaces2NetworkInterfacesInnerNetwork{
			Id:   id,
			Pool: ipPool,
		},
		IpMode:                 in.IpMode.ValueStringPointer(),
		IpAddress:              in.IpAddress.ValueStringPointer(),
		NetworkInterfaceTypeId: in.NetworkTypeId.ValueInt64Pointer(),
	}
}

// Map Child Virtual Network interface if it exists
func updateChildNetworkInterfaceMapper(
	in ChildVirtualNetworksValue,
) sdk.InstancesNetworkInterfaces3NetworkInterfacesInner {
	var id string
	if !in.NetworkGroupId.IsNull() && !in.NetworkGroupId.IsUnknown() {
		id = "networkGroup-" + strconv.FormatInt(in.NetworkGroupId.ValueInt64(), 10)
	}

	if !in.NetworkId.IsNull() && !in.NetworkId.IsUnknown() {
		id = strconv.FormatInt(in.NetworkId.ValueInt64(), 10)
	}

	if !in.SubnetId.IsNull() && !in.SubnetId.IsUnknown() {
		id = "subnet-" + strconv.FormatInt(in.SubnetId.ValueInt64(), 10)
	}

	var intfIdPtr *int64
	if !in.Id.IsNull() && !in.Id.IsUnknown() {
		intfIdPtr = in.Id.ValueInt64Pointer()
	}

	ipPool := &sdk.InstancesNetworkInterfaces3NetworkInterfacesInnerNetworkPool{}
	if !in.IpPool.IsNull() {
		ipPool.Id = in.IpPool.ValueInt64Pointer()
	}

	return sdk.InstancesNetworkInterfaces3NetworkInterfacesInner{
		Network: sdk.InstancesNetworkInterfaces3NetworkInterfacesInnerNetwork{
			Id:   id,
			Pool: ipPool,
		},
		Id:                     intfIdPtr,
		IpMode:                 in.IpMode.ValueStringPointer(),
		IpAddress:              in.IpAddress.ValueStringPointer(),
		NetworkInterfaceTypeId: in.NetworkTypeId.ValueInt64Pointer(),
	}
}

// Map Terraform tag value into an API request struct
func createTagMapper(
	in TagsValue,
) sdk.AddInstanceRequestTagsInner {
	return sdk.AddInstanceRequestTagsInner{
		Name:  in.Name.ValueStringPointer(),
		Value: in.Value.ValueStringPointer(),
	}
}

// Map Terraform tag value into an API request struct
func updateTagMapper(
	in TagsValue,
) sdk.UpdateInstanceRequestInstanceTagsInner {
	return sdk.UpdateInstanceRequestInstanceTagsInner{
		Name:  in.Name.ValueStringPointer(),
		Value: in.Value.ValueStringPointer(),
	}
}

// Map Terraform evar value into an API request struct
func evarMapper(
	in EvarsValue,
) sdk.AddInstanceRequestEvarsInner {
	return sdk.AddInstanceRequestEvarsInner{
		Name:  in.Name.ValueStringPointer(),
		Value: in.Value.ValueStringPointer(),
	}
}
