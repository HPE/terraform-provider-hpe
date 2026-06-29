# UpdateVDIPools200ResponseAnyOfVdiPoolConfigVolumesInner

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**VolumeCustomizable** | Pointer to **bool** |  | [optional] 
**VId** | Pointer to **int64** |  | [optional] 
**ReadonlyName** | Pointer to **bool** |  | [optional] 
**Size** | Pointer to **int64** |  | [optional] 
**MaxIOPS** | Pointer to **NullableString** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**RootVolume** | Pointer to **bool** |  | [optional] 
**StorageType** | Pointer to **int64** |  | [optional] 
**DatastoreId** | Pointer to **string** |  | [optional] 
**MaxStorage** | Pointer to **int64** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &UpdateVDIPools200ResponseAnyOfVdiPoolConfigVolumesInner{
    // Set fields directly
}
```

### MaxIOPS (Nullable)

Use the Nullable wrapper methods:
- `obj.MaxIOPS.IsSet()` — check if set
- `obj.MaxIOPS.Get()` — get the inner value (returns pointer)
- `obj.MaxIOPS.Set(&val)` — set the value
- `obj.MaxIOPS.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


