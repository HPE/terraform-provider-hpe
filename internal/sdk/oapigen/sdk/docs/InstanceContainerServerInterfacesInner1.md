# InstanceContainerServerInterfacesInner1

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
**Network** | Pointer to [**InstanceContainerServerInterfacesInner1Network**](InstanceContainerServerInterfacesInner1Network.md) |  | [optional] 
**Subnet** | Pointer to [**InstanceContainerServerInterfacesInner1Subnet**](InstanceContainerServerInterfacesInner1Subnet.md) |  | [optional] 
**NetworkGroup** | Pointer to [**InstanceContainerServerInterfacesInner1NetworkGroup**](InstanceContainerServerInterfacesInner1NetworkGroup.md) |  | [optional] 
**NetworkPool** | Pointer to [**InstanceContainerServerInterfacesInner1NetworkPool**](InstanceContainerServerInterfacesInner1NetworkPool.md) |  | [optional] 
**IpMode** | Pointer to **string** |  | [optional] 
**MacAddress** | Pointer to **string** |  | [optional] 
**Interfaces** | Pointer to [**[]InstanceContainerServerInstancesInnerInner1**](InstanceContainerServerInstancesInnerInner1.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &InstanceContainerServerInterfacesInner1{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


