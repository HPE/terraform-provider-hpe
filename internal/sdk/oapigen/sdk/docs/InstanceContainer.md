# InstanceContainer

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**Uuid** | Pointer to **string** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**Ip** | Pointer to **string** |  | [optional] 
**InternalIp** | Pointer to **string** |  | [optional] 
**InternalHostname** | Pointer to **string** |  | [optional] 
**ExternalHostname** | Pointer to **string** |  | [optional] 
**ExternalDomain** | Pointer to **string** |  | [optional] 
**ExternalFqdn** | Pointer to **string** |  | [optional] 
**AccountId** | Pointer to **int64** |  | [optional] 
**Instance** | Pointer to [**InstanceContainerInstance**](InstanceContainerInstance.md) |  | [optional] 
**ContainerType** | Pointer to [**InstanceContainerContainerType**](InstanceContainerContainerType.md) |  | [optional] 
**ContainerTypeSet** | Pointer to [**InstanceContainerContainerTypeSet**](InstanceContainerContainerTypeSet.md) |  | [optional] 
**Server** | Pointer to [**InstanceContainerServer**](InstanceContainerServer.md) |  | [optional] 
**Cloud** | Pointer to [**InstanceContainerCloud**](InstanceContainerCloud.md) |  | [optional] 
**Ports** | Pointer to [**[]InstanceContainerPortsInner**](InstanceContainerPortsInner.md) |  | [optional] 
**Plan** | Pointer to [**InstanceContainerPlan**](InstanceContainerPlan.md) |  | [optional] 
**DateCreated** | Pointer to **time.Time** |  | [optional] 
**LastUpdated** | Pointer to **time.Time** |  | [optional] 
**Hostname** | Pointer to **string** |  | [optional] 
**VolumeCreated** | Pointer to **bool** |  | [optional] 
**ContainerCreated** | Pointer to **bool** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &InstanceContainer{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


