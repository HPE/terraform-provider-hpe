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
		volume.SetId(vol.Id.ValueInt64())
	} else {
		// set to auto
		volume.SetId(-1)
	}

	if !vol.DatastoreId.IsNull() && !vol.DatastoreId.IsUnknown() {
		// convert to expected format of API (string)
		id := strconv.Itoa(int(vol.DatastoreId.ValueInt64()))
		volume.SetDatastoreId(id)
	}

	if !vol.DatastoreAutoSelection.IsNull() && !vol.DatastoreAutoSelection.IsUnknown() {
		volume.SetDatastoreId(vol.DatastoreAutoSelection.ValueString())
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

	return volume
}

func createSSHHostsMapper(
	sshHost SshHostsValue,
) sdk.AddClusterRequestClusterServerSshHostsInner {
	host := sdk.AddClusterRequestClusterServerSshHostsInner{}

	if !sshHost.Ip.IsNull() && !sshHost.Ip.IsUnknown() {
		host.SetIp(sshHost.Ip.ValueString())
	}

	if !sshHost.Name.IsNull() && !sshHost.Name.IsUnknown() {
		host.SetName(sshHost.Name.ValueString())
	}

	return host
}

func createNetworkInterfacesMapper(
	networkInterface NetworkInterfacesValue,
) sdk.AddClusterRequestClusterServerNetworkInterfacesInner {
	ni := sdk.AddClusterRequestClusterServerNetworkInterfacesInner{}

	if !networkInterface.NetworkId.IsNull() && !networkInterface.NetworkId.IsUnknown() {
		ni.Network.SetId(
			sdk.AddClusterRequestClusterServerNetworkInterfacesInnerNetworkId{
				Int64: networkInterface.NetworkId.ValueInt64Pointer(),
			},
		)
	}

	if !networkInterface.IpMode.IsNull() && !networkInterface.IpMode.IsUnknown() {
		ni.SetIpMode(networkInterface.IpMode.ValueString())
	}

	return ni
}
