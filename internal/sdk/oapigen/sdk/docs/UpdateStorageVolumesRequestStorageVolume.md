# UpdateStorageVolumesRequestStorageVolume

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** | A unique name scoped to your account for the storage volume | [optional] 
**Type** | Pointer to **string** | Storage Type Code or ID | [optional] 
**Config** | Pointer to **map[string]interface{}** | Configuration object with parameters that vary by &#x60;type&#x60;. | [optional] 
**StorageServer** | Pointer to [**UpdateStorageVolumesRequestStorageVolumeStorageServer**](UpdateStorageVolumesRequestStorageVolumeStorageServer.md) |  | [optional] 
**StorageGroup** | Pointer to [**UpdateStorageVolumesRequestStorageVolumeStorageGroup**](UpdateStorageVolumesRequestStorageVolumeStorageGroup.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &UpdateStorageVolumesRequestStorageVolume{
    // Set fields directly
}
```


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


