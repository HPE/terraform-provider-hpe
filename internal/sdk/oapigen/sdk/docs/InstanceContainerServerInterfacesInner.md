# InstanceContainerServerInterfacesInner

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**UniqueId** | Pointer to **string** |  | [optional] 
**PublicIpAddress** | Pointer to **string** |  | [optional] 
**IpAddress** | Pointer to **string** |  | [optional] 
**Dhcp** | Pointer to **bool** |  | [optional] 
**Active** | Pointer to **bool** |  | [optional] 
**PoolAssigned** | Pointer to **bool** |  | [optional] 
**PrimaryInterface** | Pointer to **bool** |  | [optional] 
**Network** | Pointer to [**NullableInstanceContainerServerInterfacesInnerNetwork**](InstanceContainerServerInterfacesInnerNetwork.md) |  | [optional] 
**Subnet** | Pointer to [**NullableInstanceContainerServerInterfacesInnerSubnet**](InstanceContainerServerInterfacesInnerSubnet.md) |  | [optional] 
**NetworkGroup** | Pointer to [**NullableInstanceContainerServerInterfacesInnerNetworkGroup**](InstanceContainerServerInterfacesInnerNetworkGroup.md) |  | [optional] 
**NetworkPool** | Pointer to [**NullableInstanceContainerServerInterfacesInnerNetworkPool**](InstanceContainerServerInterfacesInnerNetworkPool.md) |  | [optional] 
**IpMode** | Pointer to **string** |  | [optional] 
**MacAddress** | Pointer to **string** |  | [optional] 
**Interfaces** | Pointer to [**[]InstanceContainerServerInstancesInnerInner**](InstanceContainerServerInstancesInnerInner.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &InstanceContainerServerInterfacesInner{
    // Set fields directly
}
```

### Network (Nullable)

Use the Nullable wrapper methods:
- `obj.Network.IsSet()` — check if set
- `obj.Network.Get()` — get the inner value (returns pointer)
- `obj.Network.Set(&val)` — set the value
- `obj.Network.Unset()` — clear the value
### Subnet (Nullable)

Use the Nullable wrapper methods:
- `obj.Subnet.IsSet()` — check if set
- `obj.Subnet.Get()` — get the inner value (returns pointer)
- `obj.Subnet.Set(&val)` — set the value
- `obj.Subnet.Unset()` — clear the value
### NetworkGroup (Nullable)

Use the Nullable wrapper methods:
- `obj.NetworkGroup.IsSet()` — check if set
- `obj.NetworkGroup.Get()` — get the inner value (returns pointer)
- `obj.NetworkGroup.Set(&val)` — set the value
- `obj.NetworkGroup.Unset()` — clear the value
### NetworkPool (Nullable)

Use the Nullable wrapper methods:
- `obj.NetworkPool.IsSet()` — check if set
- `obj.NetworkPool.Get()` — get the inner value (returns pointer)
- `obj.NetworkPool.Set(&val)` — set the value
- `obj.NetworkPool.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


