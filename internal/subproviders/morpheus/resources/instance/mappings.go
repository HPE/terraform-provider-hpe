package instance

import (
	"strconv"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/sdk"
)

// Map Terraform volume value into an API request struct
func volumeMapper(vol VolumesValue) sdk.AddCatalogItemTypeRequestCatalogItemTypeOneOfConfigVolumesInner {
	volume := sdk.AddCatalogItemTypeRequestCatalogItemTypeOneOfConfigVolumesInner{}
	if !vol.Id.IsNull() {
		volume.SetId(vol.Id.ValueInt64())
	}

	if !vol.Name.IsNull() {
		volume.SetName(vol.Name.ValueString())
	}

	if !vol.RootVolume.IsNull() {
		volume.SetRootVolume(vol.RootVolume.ValueBool())
	}

	if !vol.Size.IsNull() {
		volume.SetSize(vol.Size.ValueInt64())
	}

	if !vol.StorageTypeId.IsNull() {
		volume.SetStorageType(vol.StorageTypeId.ValueInt64())
	}

	if !vol.DatastoreId.IsNull() {
		volume.DatastoreId = &sdk.AddCatalogItemTypeRequestCatalogItemTypeOneOfConfigVolumesInnerDatastoreId{}

		id := strconv.Itoa(int(vol.DatastoreId.ValueInt64()))
		volume.DatastoreId.String = &id
	}

	if !vol.DatastoreAutoSelection.IsNull() {
		volume.DatastoreId = &sdk.AddCatalogItemTypeRequestCatalogItemTypeOneOfConfigVolumesInnerDatastoreId{}
		volume.DatastoreId.String = vol.DatastoreAutoSelection.ValueStringPointer()
	}

	return volume
}

// Map Terraform network interface value into an API request struct
func networkInterfaceMapper(in NetworkInterfacesValue) sdk.AddCatalogItemTypeRequestCatalogItemTypeOneOfConfigNetworkInterfacesInner {
	var id string
	if !in.NetworkGroupId.IsNull() {
		id = "networkGroup-" + in.NetworkGroupId.String()
	}

	if !in.NetworkId.IsNull() {
		id = in.NetworkId.String()
	}

	return sdk.AddCatalogItemTypeRequestCatalogItemTypeOneOfConfigNetworkInterfacesInner{
		Network: sdk.AddCatalogItemTypeRequestCatalogItemTypeOneOfConfigNetworkInterfacesInnerNetwork{
			Id: id,
		},
		IpMode:    in.IpMode.ValueStringPointer(),
		IpAddress: in.IpAddress.ValueStringPointer(),
	}
}

// Map Terraform tag value into an API request struct
func tagMapper(in TagsValue) sdk.AddCatalogItemTypeRequestCatalogItemTypeOneOfConfigEvarsInner {
	return sdk.AddCatalogItemTypeRequestCatalogItemTypeOneOfConfigEvarsInner{
		Name:  in.Name.ValueStringPointer(),
		Value: in.Value.ValueStringPointer(),
	}
}

// Map Terraform evar value into an API request struct
func evarMapper(in EvarsValue) sdk.AddCatalogItemTypeRequestCatalogItemTypeOneOfConfigEvarsInner {
	return sdk.AddCatalogItemTypeRequestCatalogItemTypeOneOfConfigEvarsInner{
		Name:  in.Name.ValueStringPointer(),
		Value: in.Value.ValueStringPointer(),
	}
}
