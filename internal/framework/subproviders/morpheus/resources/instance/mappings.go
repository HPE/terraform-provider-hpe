package instance

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"

	"github.com/HPE/terraform-provider-hpe/internal/framework/subproviders/morpheus/convert"
)

// Map Terraform volume value into an API request struct
func volumeMapper(
	vol VolumesValue,
) sdk.AddCatalogItemTypeRequestCatalogItemTypeOneOfConfigVolumesInner {
	volume := sdk.AddCatalogItemTypeRequestCatalogItemTypeOneOfConfigVolumesInner{}
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
			AddCatalogItemTypeRequestCatalogItemTypeOneOfConfigVolumesInnerDatastoreId{}

		id := strconv.Itoa(int(vol.DatastoreId.ValueInt64()))
		volume.DatastoreId.String = &id
	}

	if !vol.DatastoreAutoSelection.IsNull() && !vol.DatastoreAutoSelection.IsUnknown() {
		volume.DatastoreId = &sdk.
			AddCatalogItemTypeRequestCatalogItemTypeOneOfConfigVolumesInnerDatastoreId{}
		volume.DatastoreId.String = vol.DatastoreAutoSelection.ValueStringPointer()
	}

	return volume
}

// Map Terraform network interface value into an API request struct
func networkInterfaceMapper(
	ctx context.Context,
) func(in NetworkInterfacesValue) sdk.InstancesNetworkInterfaces {
	return func(in NetworkInterfacesValue) sdk.InstancesNetworkInterfaces {
		var id string
		if !in.NetworkGroupId.IsNull() {
			id = "networkGroup-" + in.NetworkGroupId.String()
		}

		if !in.NetworkId.IsNull() {
			id = in.NetworkId.String()
		}

		ipPool := sdk.NewInstancesNetworkInterfacesNetworkPoolWithDefaults()
		if !in.IpPool.IsNull() {
			ipPool.SetId(in.IpPool.ValueInt64())
		}

		childNetworkInterfaces, diags := convert.FromListType(
			ctx,
			in.ChildVirtualNetworks,
			childNetworkInterfaceMapper,
		)
		if diags.HasError() {
			tflog.Error(ctx, "cannot convert child virtual network interfaces")
		}

		return sdk.InstancesNetworkInterfaces{
			Network: sdk.
				InstancesNetworkInterfacesNetwork{
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
func childNetworkInterfaceMapper(
	in ChildVirtualNetworksValue,
) sdk.InstancesNetworkInterfacesNetworkInterfacesInner {
	var id string
	if !in.NetworkGroupId.IsNull() {
		id = "networkGroup-" + in.NetworkGroupId.String()
	}

	if !in.NetworkId.IsNull() {
		id = in.NetworkId.String()
	}

	ipPool := sdk.NewInstancesNetworkInterfacesNetworkPoolWithDefaults()
	if !in.IpPool.IsNull() {
		ipPool.SetId(in.IpPool.ValueInt64())
	}

	return sdk.InstancesNetworkInterfacesNetworkInterfacesInner{
		Network: sdk.InstancesNetworkInterfacesNetwork{
			Id:   id,
			Pool: ipPool,
		},
		IpMode:                 in.IpMode.ValueStringPointer(),
		IpAddress:              in.IpAddress.ValueStringPointer(),
		NetworkInterfaceTypeId: in.NetworkTypeId.ValueInt64Pointer(),
	}
}

// Map Terraform tag value into an API request struct
func tagMapper(
	in TagsValue,
) sdk.AddCatalogItemTypeRequestCatalogItemTypeOneOfConfigEvarsInner {
	return sdk.AddCatalogItemTypeRequestCatalogItemTypeOneOfConfigEvarsInner{
		Name:  in.Name.ValueStringPointer(),
		Value: in.Value.ValueStringPointer(),
	}
}

// Map Terraform evar value into an API request struct
func evarMapper(
	in EvarsValue,
) sdk.AddCatalogItemTypeRequestCatalogItemTypeOneOfConfigEvarsInner {
	return sdk.AddCatalogItemTypeRequestCatalogItemTypeOneOfConfigEvarsInner{
		Name:  in.Name.ValueStringPointer(),
		Value: in.Value.ValueStringPointer(),
	}
}
