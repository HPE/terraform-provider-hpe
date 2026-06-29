# ClusterCreateServerVolumesInner

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** | The id for the LV configuration being created | [optional] [default to -1]
**RootVolume** | Pointer to **bool** | If set to false then a non-root LV will be created | [optional] [default to true]
**Name** | **string** | Name/type of the LV being created | [default to "root"]
**Size** | Pointer to **int64** | Size of the LV to be created in GBs  Default is from the service plan  | [optional] 
**SizeId** | Pointer to **NullableString** | Can be used to select pre-existing LV choices from Morpheus | [optional] 
**StorageType** | Pointer to **int64** | Identifier for LV type | [optional] 
**DatastoreId** | **NullableString** | The ID of the specific datastore. Auto selection can be specified as auto or &#x60;autoCluster&#x60; (for clusters). | 

## Usage

Instantiate with a Go composite literal:

```go
obj := &ClusterCreateServerVolumesInner{
    // Set fields directly
}
```

### SizeId (Nullable)

Use the Nullable wrapper methods:
- `obj.SizeId.IsSet()` — check if set
- `obj.SizeId.Get()` — get the inner value (returns pointer)
- `obj.SizeId.Set(&val)` — set the value
- `obj.SizeId.Unset()` — clear the value
### DatastoreId (Nullable)

Use the Nullable wrapper methods:
- `obj.DatastoreId.IsSet()` — check if set
- `obj.DatastoreId.Get()` — get the inner value (returns pointer)
- `obj.DatastoreId.Set(&val)` — set the value
- `obj.DatastoreId.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


