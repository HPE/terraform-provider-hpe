package instance

import (
	"strconv"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/sdk"
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
	in NetworkInterfacesValue,
) sdk.AddCatalogItemTypeRequestCatalogItemTypeOneOfConfigNetworkInterfacesInner {
	var id string
	if !in.NetworkGroupId.IsNull() {
		id = "networkGroup-" + in.NetworkGroupId.String()
	}

	if !in.NetworkId.IsNull() {
		id = in.NetworkId.String()
	}

	return sdk.AddCatalogItemTypeRequestCatalogItemTypeOneOfConfigNetworkInterfacesInner{
		Network: sdk.
			AddCatalogItemTypeRequestCatalogItemTypeOneOfConfigNetworkInterfacesInnerNetwork{
			Id: id,
		},
		IpMode:    in.IpMode.ValueStringPointer(),
		IpAddress: in.IpAddress.ValueStringPointer(),
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
