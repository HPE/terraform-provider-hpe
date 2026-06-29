# ListBillingInstances200ResponseAllOfBillingInfoInstancesInnerContainersInnerUsagesInnerVolumesInnerDatastore

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**Type** | Pointer to **string** |  | [optional] 
**ExternalId** | Pointer to **string** |  | [optional] 
**InternalId** | Pointer to **NullableString** |  | [optional] 
**ExternalPath** | Pointer to **NullableString** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &ListBillingInstances200ResponseAllOfBillingInfoInstancesInnerContainersInnerUsagesInnerVolumesInnerDatastore{
    // Set fields directly
}
```

### InternalId (Nullable)

Use the Nullable wrapper methods:
- `obj.InternalId.IsSet()` — check if set
- `obj.InternalId.Get()` — get the inner value (returns pointer)
- `obj.InternalId.Set(&val)` — set the value
- `obj.InternalId.Unset()` — clear the value
### ExternalPath (Nullable)

Use the Nullable wrapper methods:
- `obj.ExternalPath.IsSet()` — check if set
- `obj.ExternalPath.Get()` — get the inner value (returns pointer)
- `obj.ExternalPath.Set(&val)` — set the value
- `obj.ExternalPath.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


