// (C) Copyright 2026 Hewlett Packard Enterprise Development LP

package instancenode

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"

	sdk "github.com/HPE/terraform-provider-hpe/internal/sdk/oapigen"
)

func ptr[T any](v T) *T { return &v }

func TestUnitResolvePrimaryMacAddress_PrimaryFlagged(t *testing.T) {
	t.Parallel()

	server := &sdk.InstanceContainerServer2{
		Interfaces: []sdk.InstanceContainerServerInterfacesInner1{
			{MacAddress: ptr("aa:bb:cc:dd:ee:01"), PrimaryInterface: ptr(false)},
			{MacAddress: ptr("aa:bb:cc:dd:ee:02"), PrimaryInterface: ptr(true)},
			{MacAddress: ptr("aa:bb:cc:dd:ee:03"), PrimaryInterface: ptr(false)},
		},
	}

	result := resolvePrimaryMacAddress(server)
	assert.Equal(t, types.StringValue("aa:bb:cc:dd:ee:02"), result)
}

func TestUnitResolvePrimaryMacAddress_NoPrimaryFlag(t *testing.T) {
	t.Parallel()

	server := &sdk.InstanceContainerServer2{
		Interfaces: []sdk.InstanceContainerServerInterfacesInner1{
			{MacAddress: ptr("aa:bb:cc:dd:ee:01")},
			{MacAddress: ptr("aa:bb:cc:dd:ee:02")},
		},
	}

	result := resolvePrimaryMacAddress(server)
	assert.Equal(t, types.StringValue("aa:bb:cc:dd:ee:01"), result)
}

func TestUnitResolvePrimaryMacAddress_NoInterfaces(t *testing.T) {
	t.Parallel()

	server := &sdk.InstanceContainerServer2{
		Interfaces: []sdk.InstanceContainerServerInterfacesInner1{},
	}

	result := resolvePrimaryMacAddress(server)
	assert.True(t, result.IsNull())
}

func TestUnitResolvePrimaryMacAddress_NilServer(t *testing.T) {
	t.Parallel()

	result := resolvePrimaryMacAddress(nil)
	assert.True(t, result.IsNull())
}

func TestUnitRefreshNodeState_AllMetadataPopulated(t *testing.T) {
	t.Parallel()

	state := &instanceNodeModel{}
	cd := &sdk.InstanceContainer2{
		Id:           ptr(int64(100)),
		Hostname:     ptr("node-1"),
		InternalIp:   ptr("10.0.0.5"),
		ExternalFqdn: ptr("node-1.example.com"),
		Ip:           ptr("192.168.1.10"),
		Server: &sdk.InstanceContainerServer2{
			Id: ptr(int64(200)),
			Interfaces: []sdk.InstanceContainerServerInterfacesInner1{
				{MacAddress: ptr("52:54:00:04:2e:e5"), PrimaryInterface: ptr(true)},
			},
		},
	}

	populateInstanceNodeMetadata(state, cd)

	assert.Equal(t, types.Int64Value(100), state.ContainerID)
	assert.Equal(t, types.Int64Value(200), state.ServerID)
	assert.Equal(t, types.StringValue("192.168.1.10"), state.IPAddress)
	assert.Equal(t, types.StringValue("node-1"), state.Hostname)
	assert.Equal(t, types.StringValue("10.0.0.5"), state.InternalIP)
	assert.Equal(t, types.StringValue("node-1.example.com"), state.ExternalFQDN)
	assert.Equal(t, types.StringValue("52:54:00:04:2e:e5"), state.MacAddress)
}

func TestUnitRefreshNodeState_HostnameAbsent(t *testing.T) {
	t.Parallel()

	state := &instanceNodeModel{}
	cd := &sdk.InstanceContainer2{
		Id:           ptr(int64(100)),
		InternalIp:   ptr("10.0.0.5"),
		ExternalFqdn: ptr("node-1.example.com"),
		Server:       &sdk.InstanceContainerServer2{Id: ptr(int64(200))},
	}

	populateInstanceNodeMetadata(state, cd)

	assert.True(t, state.Hostname.IsNull())
	assert.Equal(t, types.StringValue("10.0.0.5"), state.InternalIP)
	assert.Equal(t, types.StringValue("node-1.example.com"), state.ExternalFQDN)
	assert.True(t, state.MacAddress.IsNull())
}

func TestUnitRefreshNodeState_InternalIPAbsent(t *testing.T) {
	t.Parallel()

	state := &instanceNodeModel{}
	cd := &sdk.InstanceContainer2{
		Id:       ptr(int64(100)),
		Hostname: ptr("node-1"),
		Server:   &sdk.InstanceContainerServer2{Id: ptr(int64(200))},
	}

	populateInstanceNodeMetadata(state, cd)

	assert.Equal(t, types.StringValue("node-1"), state.Hostname)
	assert.True(t, state.InternalIP.IsNull())
	assert.True(t, state.ExternalFQDN.IsNull())
}

func TestUnitRefreshNodeState_ExternalFQDNAbsent(t *testing.T) {
	t.Parallel()

	state := &instanceNodeModel{}
	cd := &sdk.InstanceContainer2{
		Id:         ptr(int64(100)),
		Hostname:   ptr("node-1"),
		InternalIp: ptr("10.0.0.5"),
		Server:     &sdk.InstanceContainerServer2{Id: ptr(int64(200))},
	}

	populateInstanceNodeMetadata(state, cd)

	assert.True(t, state.ExternalFQDN.IsNull())
}

func TestUnitRefreshNodeState_MacAddressAbsent(t *testing.T) {
	t.Parallel()

	state := &instanceNodeModel{}
	cd := &sdk.InstanceContainer2{
		Id:     ptr(int64(100)),
		Server: &sdk.InstanceContainerServer2{Id: ptr(int64(200))},
	}

	populateInstanceNodeMetadata(state, cd)

	assert.True(t, state.MacAddress.IsNull())
}

// TestUnitServerResourcePoolID_PopulatedFromServer verifies that
// server_resource_pool_id is populated from the server's pool.
func TestUnitServerResourcePoolID_PopulatedFromServer(t *testing.T) {
	t.Parallel()

	state := &instanceNodeModel{}
	cd := &sdk.InstanceContainer2{
		Id: ptr(int64(100)),
		Server: &sdk.InstanceContainerServer2{
			Id:             ptr(int64(200)),
			ResourcePoolId: ptr(int64(42)),
		},
	}

	populateInstanceNodeMetadata(state, cd)

	assert.Equal(t, types.Int64Value(42), state.ServerResourcePoolID)
}

// TestUnitServerResourcePoolID_NullWhenAbsent verifies that
// server_resource_pool_id is null when the server has no pool.
func TestUnitServerResourcePoolID_NullWhenAbsent(t *testing.T) {
	t.Parallel()

	state := &instanceNodeModel{}
	cd := &sdk.InstanceContainer2{
		Id:     ptr(int64(100)),
		Server: &sdk.InstanceContainerServer2{Id: ptr(int64(200))},
	}

	populateInstanceNodeMetadata(state, cd)

	assert.True(t, state.ServerResourcePoolID.IsNull())
}

// TestUnitServerResourcePoolID_NullWhenNoServer verifies that
// server_resource_pool_id is null when the server is nil.
func TestUnitServerResourcePoolID_NullWhenNoServer(t *testing.T) {
	t.Parallel()

	state := &instanceNodeModel{}
	cd := &sdk.InstanceContainer2{
		Id: ptr(int64(100)),
	}

	populateInstanceNodeMetadata(state, cd)

	assert.True(t, state.ServerResourcePoolID.IsNull())
}

// TestUnitResourcePoolID_NotReadBack verifies that resource_pool_id is
// never written by populateInstanceNodeMetadata — it remains whatever
// the caller set it to (the regression was that a read-back caused
// ModifyPlan to fire against virtual instances).
func TestUnitResourcePoolID_NotReadBack(t *testing.T) {
	t.Parallel()

	state := &instanceNodeModel{
		ResourcePoolID: types.Int64Null(),
	}
	cd := &sdk.InstanceContainer2{
		Id: ptr(int64(100)),
		Server: &sdk.InstanceContainerServer2{
			Id:             ptr(int64(200)),
			ResourcePoolId: ptr(int64(99)),
		},
	}

	populateInstanceNodeMetadata(state, cd)

	// resource_pool_id must remain null — never populated from API.
	assert.True(t, state.ResourcePoolID.IsNull(),
		"resource_pool_id must not be read back from the API")
	// But server_resource_pool_id must reflect the server's pool.
	assert.Equal(t, types.Int64Value(99), state.ServerResourcePoolID)
}
