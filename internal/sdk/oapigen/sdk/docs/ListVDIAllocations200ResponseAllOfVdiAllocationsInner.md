# ListVDIAllocations200ResponseAllOfVdiAllocationsInner

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**Pool** | Pointer to [**NullableListVDIAllocations200ResponseAllOfVdiAllocationsInnerPool**](ListVDIAllocations200ResponseAllOfVdiAllocationsInnerPool.md) |  | [optional] 
**Instance** | Pointer to [**ListVDIAllocations200ResponseAllOfVdiAllocationsInnerInstance**](ListVDIAllocations200ResponseAllOfVdiAllocationsInnerInstance.md) |  | [optional] 
**User** | Pointer to [**ListVDIAllocations200ResponseAllOfVdiAllocationsInnerUser**](ListVDIAllocations200ResponseAllOfVdiAllocationsInnerUser.md) |  | [optional] 
**LocalUserCreated** | Pointer to **bool** |  | [optional] 
**Persistent** | Pointer to **bool** |  | [optional] 
**Status** | Pointer to **string** |  | [optional] 
**DateCreated** | Pointer to **time.Time** |  | [optional] 
**LastUpdated** | Pointer to **NullableTime** |  | [optional] 
**LastReserved** | Pointer to **NullableTime** |  | [optional] 
**ReleaseDate** | Pointer to **NullableTime** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &ListVDIAllocations200ResponseAllOfVdiAllocationsInner{
    // Set fields directly
}
```

### Pool (Nullable)

Use the Nullable wrapper methods:
- `obj.Pool.IsSet()` — check if set
- `obj.Pool.Get()` — get the inner value (returns pointer)
- `obj.Pool.Set(&val)` — set the value
- `obj.Pool.Unset()` — clear the value
### LastUpdated (Nullable)

Use the Nullable wrapper methods:
- `obj.LastUpdated.IsSet()` — check if set
- `obj.LastUpdated.Get()` — get the inner value (returns pointer)
- `obj.LastUpdated.Set(&val)` — set the value
- `obj.LastUpdated.Unset()` — clear the value
### LastReserved (Nullable)

Use the Nullable wrapper methods:
- `obj.LastReserved.IsSet()` — check if set
- `obj.LastReserved.Get()` — get the inner value (returns pointer)
- `obj.LastReserved.Set(&val)` — set the value
- `obj.LastReserved.Unset()` — clear the value
### ReleaseDate (Nullable)

Use the Nullable wrapper methods:
- `obj.ReleaseDate.IsSet()` — check if set
- `obj.ReleaseDate.Get()` — get the inner value (returns pointer)
- `obj.ReleaseDate.Set(&val)` — set the value
- `obj.ReleaseDate.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


