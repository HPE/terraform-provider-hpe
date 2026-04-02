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
		volume.SetId(vol.Id.ValueInt64())
	} else {
		volume.SetId(-1)
	}

	if !vol.Name.IsNull() && !vol.Name.IsUnknown() {
		volume.SetName(vol.Name.ValueString())
	}

	if !vol.RootVolume.IsNull() && !vol.RootVolume.IsUnknown() {
		volume.SetRootVolume(vol.RootVolume.ValueBool())
	}

	if !vol.Size.IsNull() && !vol.Size.IsUnknown() {
		volume.SetSize(vol.Size.ValueInt64())
	}

	if !vol.StorageTypeId.IsNull() && !vol.StorageTypeId.IsUnknown() {
		volume.SetStorageType(vol.StorageTypeId.ValueInt64())
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
		volume.SetId(vol.Id.ValueInt64())
	} else {
		volume.SetId(-1)
	}

	if !vol.Name.IsNull() && !vol.Name.IsUnknown() {
		volume.SetName(vol.Name.ValueString())
	}

	if !vol.RootVolume.IsNull() && !vol.RootVolume.IsUnknown() {
		volume.SetRootVolume(vol.RootVolume.ValueBool())
	}

	if !vol.Size.IsNull() && !vol.Size.IsUnknown() {
		volume.SetSize(vol.Size.ValueInt64())
	}

	if !vol.StorageTypeId.IsNull() && !vol.StorageTypeId.IsUnknown() {
		volume.SetStorageType(vol.StorageTypeId.ValueInt64())
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
		if !in.NetworkGroupId.IsNull() {
			id = "networkGroup-" + in.NetworkGroupId.String()
		}

		if !in.NetworkId.IsNull() {
			id = in.NetworkId.String()
		}

		ipPool := sdk.NewInstancesNetworkInterfaces2NetworkPoolWithDefaults()
		if !in.IpPool.IsNull() {
			ipPool.SetId(in.IpPool.ValueInt64())
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
		if !in.NetworkGroupId.IsNull() {
			id = "networkGroup-" + in.NetworkGroupId.String()
		}

		if !in.NetworkId.IsNull() {
			id = in.NetworkId.String()
		}

		ipPool := sdk.NewInstancesNetworkInterfaces3NetworkPoolWithDefaults()
		if !in.IpPool.IsNull() {
			ipPool.SetId(in.IpPool.ValueInt64())
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
	if !in.NetworkGroupId.IsNull() {
		id = "networkGroup-" + in.NetworkGroupId.String()
	}

	if !in.NetworkId.IsNull() {
		id = in.NetworkId.String()
	}

	ipPool := sdk.NewInstancesNetworkInterfaces2NetworkInterfacesInnerNetworkPoolWithDefaults()
	if !in.IpPool.IsNull() {
		ipPool.SetId(in.IpPool.ValueInt64())
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
	if !in.NetworkGroupId.IsNull() {
		id = "networkGroup-" + in.NetworkGroupId.String()
	}

	if !in.NetworkId.IsNull() {
		id = in.NetworkId.String()
	}

	ipPool := sdk.NewInstancesNetworkInterfaces3NetworkInterfacesInnerNetworkPoolWithDefaults()
	if !in.IpPool.IsNull() {
		ipPool.SetId(in.IpPool.ValueInt64())
	}

	return sdk.InstancesNetworkInterfaces3NetworkInterfacesInner{
		Network: sdk.InstancesNetworkInterfaces3NetworkInterfacesInnerNetwork{
			Id:   id,
			Pool: ipPool,
		},
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
