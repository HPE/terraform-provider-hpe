# InstanceContainer2

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
**Instance** | Pointer to [**InstanceContainer2Instance**](InstanceContainer2Instance.md) |  | [optional] 
**ContainerType** | Pointer to [**InstanceContainer2ContainerType**](InstanceContainer2ContainerType.md) |  | [optional] 
**ContainerTypeSet** | Pointer to [**InstanceContainer2ContainerTypeSet**](InstanceContainer2ContainerTypeSet.md) |  | [optional] 
**Server** | Pointer to [**InstanceContainerServer2**](InstanceContainerServer2.md) |  | [optional] 
**Cloud** | Pointer to [**InstanceContainer2Cloud**](InstanceContainer2Cloud.md) |  | [optional] 
**Ports** | Pointer to [**[]InstanceContainer1PortsInner**](InstanceContainer1PortsInner.md) |  | [optional] 
**Plan** | Pointer to [**InstanceContainer2Plan**](InstanceContainer2Plan.md) |  | [optional] 
**DateCreated** | Pointer to **time.Time** |  | [optional] 
**LastUpdated** | Pointer to **time.Time** |  | [optional] 
**Hostname** | Pointer to **string** |  | [optional] 
**VolumeCreated** | Pointer to **bool** |  | [optional] 
**ContainerCreated** | Pointer to **bool** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &InstanceContainer2{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


