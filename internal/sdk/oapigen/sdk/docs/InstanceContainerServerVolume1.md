# InstanceContainerServerVolume1

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**ControllerId** | Pointer to **int64** |  | [optional] 
**ControllerMountPoint** | Pointer to **string** |  | [optional] 
**Resizeable** | Pointer to **bool** |  | [optional] 
**RootVolume** | Pointer to **bool** |  | [optional] 
**DeviceName** | Pointer to **string** |  | [optional] 
**DeviceDisplayName** | Pointer to **string** |  | [optional] 
**DiskMode** | Pointer to **string** |  | [optional] 
**DiskType** | Pointer to **string** |  | [optional] 
**Category** | Pointer to **string** |  | [optional] 
**StorageProfile** | Pointer to **string** | Storage Profile Code for the volume storage profile assignment. eg. &#x60;\&quot;kvm-cache-none\&quot;&#x60; or &#x60;\&quot;kvm-cache-directsync\&quot;&#x60;. Use &#x60;/api/provision-types?code&#x3D;kvm&#x60; to see the available &#x60;storageProfiles&#x60; for HVM and KVM. | [optional] 
**StorageServer** | Pointer to [**InstanceContainerServerVolumeStorageServer1**](InstanceContainerServerVolumeStorageServer1.md) |  | [optional] 
**ZoneId** | Pointer to **int64** |  | [optional] 
**Zone** | Pointer to [**InstanceContainerServerVolumeZone1**](InstanceContainerServerVolumeZone1.md) |  | [optional] 
**Datastore** | Pointer to [**InstanceContainerServerVolumeDatastore1**](InstanceContainerServerVolumeDatastore1.md) |  | [optional] 
**UnitNumber** | Pointer to **string** |  | [optional] 
**TypeId** | Pointer to **int64** |  | [optional] 
**ConfigurableIOPS** | Pointer to **bool** |  | [optional] 
**DatastoreId** | Pointer to **int64** |  | [optional] 
**MaxStorage** | Pointer to **int64** |  | [optional] 
**MaxIOPS** | Pointer to **string** |  | [optional] 
**Uuid** | Pointer to **string** |  | [optional] 
**UniqueId** | Pointer to **string** |  | [optional] 
**ExternalId** | Pointer to **string** |  | [optional] 
**InternalId** | Pointer to **string** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &InstanceContainerServerVolume1{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


