# AddStorageVolumesRequestStorageVolume

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | **string** | A unique name scoped to your account for the storage volume | 
**Type** | **string** | Storage Type Code or ID | 
**Config** | Pointer to [**AddStorageVolumesRequestStorageVolumeConfig**](AddStorageVolumesRequestStorageVolumeConfig.md) |  | [optional] 
**ProvisionType** | Pointer to **string** | Provision type for storage volume types that support it (for example 3Par - FULL, TPVV, SNP, PEER, TDVV). | [optional] 
**StorageServer** | [**AddStorageVolumesRequestStorageVolumeStorageServer**](AddStorageVolumesRequestStorageVolumeStorageServer.md) |  | 
**StorageGroup** | Pointer to [**AddStorageVolumesRequestStorageVolumeStorageGroup**](AddStorageVolumesRequestStorageVolumeStorageGroup.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &AddStorageVolumesRequestStorageVolume{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


