// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package cluster

import (
	"strconv"

	"github.com/HewlettPackard/hpe-morpheus-go-sdk/oapigen/sdk"
)

func createVolumeMapper(
	vol VolumesValue,
) sdk.AddClusterRequestClusterServerVolumesInner {
	volume := sdk.AddClusterRequestClusterServerVolumesInner{}

	if !vol.Id.IsNull() && !vol.Id.IsUnknown() {
		volume.Id = vol.Id.ValueInt64Pointer()
	} else {
		// set to auto
		id := int64(-1)
		volume.Id = &id
	}

	if !vol.DatastoreId.IsNull() && !vol.DatastoreId.IsUnknown() {
		// convert to expected format of API (string)
		id := strconv.Itoa(int(vol.DatastoreId.ValueInt64()))
		volume.DatastoreId = *sdk.NewNullableString(&id)
	}

	if !vol.DatastoreAutoSelection.IsNull() && !vol.DatastoreAutoSelection.IsUnknown() {
		volume.DatastoreId = *sdk.NewNullableString(vol.DatastoreAutoSelection.ValueStringPointer())
	}

	if !vol.Name.IsNull() && !vol.Name.IsUnknown() {
		volume.Name = vol.Name.ValueString()
	}

	if !vol.RootVolume.IsNull() && !vol.RootVolume.IsUnknown() {
		volume.RootVolume = vol.RootVolume.ValueBoolPointer()
	}

	if !vol.Size.IsNull() && !vol.Size.IsUnknown() {
		volume.Size = vol.Size.ValueInt64Pointer()
	}

	if !vol.StorageTypeId.IsNull() && !vol.StorageTypeId.IsUnknown() {
		volume.StorageType = vol.StorageTypeId.ValueInt64Pointer()
	}

	return volume
}

func createSSHHostsMapper(
	sshHost SshHostsValue,
) sdk.AddClusterRequestClusterServerSshHostsInner {
	host := sdk.AddClusterRequestClusterServerSshHostsInner{}

	if !sshHost.Ip.IsNull() && !sshHost.Ip.IsUnknown() {
		host.Ip = sshHost.Ip.ValueStringPointer()
	}

	if !sshHost.Name.IsNull() && !sshHost.Name.IsUnknown() {
		host.Name = sshHost.Name.ValueStringPointer()
	}

	return host
}

func createNetworkInterfacesMapper(
	networkInterface NetworkInterfacesValue,
) sdk.AddClusterRequestClusterServerNetworkInterfacesInner {
	ni := sdk.AddClusterRequestClusterServerNetworkInterfacesInner{}

	if !networkInterface.NetworkId.IsNull() && !networkInterface.NetworkId.IsUnknown() {
		ni.Network.Id = sdk.AddClusterRequestClusterServerNetworkInterfacesInnerNetworkId{
			Int64: networkInterface.NetworkId.ValueInt64Pointer(),
		}
	}

	if !networkInterface.IpMode.IsNull() && !networkInterface.IpMode.IsUnknown() {
		ni.IpMode = networkInterface.IpMode.ValueStringPointer()
	}

	return ni
}
