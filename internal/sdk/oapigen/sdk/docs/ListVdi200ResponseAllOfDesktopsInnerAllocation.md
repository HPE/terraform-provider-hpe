# ListVdi200ResponseAllOfDesktopsInnerAllocation

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**PoolId** | Pointer to **int64** |  | [optional] 
**Pool** | Pointer to [**NullableListVdi200ResponseAllOfDesktopsInnerAllocationPool**](ListVdi200ResponseAllOfDesktopsInnerAllocationPool.md) |  | [optional] 
**Instance** | Pointer to [**ListVdi200ResponseAllOfDesktopsInnerAllocationInstance**](ListVdi200ResponseAllOfDesktopsInnerAllocationInstance.md) |  | [optional] 
**User** | Pointer to [**ListVdi200ResponseAllOfDesktopsInnerAllocationUser**](ListVdi200ResponseAllOfDesktopsInnerAllocationUser.md) |  | [optional] 
**LocalUserCreated** | Pointer to **bool** |  | [optional] 
**Persistent** | Pointer to **bool** |  | [optional] 
**Recyclable** | Pointer to **bool** |  | [optional] 
**Status** | Pointer to **string** |  | [optional] 
**DateCreated** | Pointer to **time.Time** |  | [optional] 
**LastUpdated** | Pointer to **time.Time** |  | [optional] 
**LastReserved** | Pointer to **NullableTime** |  | [optional] 
**ReleaseDate** | Pointer to **time.Time** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &ListVdi200ResponseAllOfDesktopsInnerAllocation{
    // Set fields directly
}
```

### Pool (Nullable)

Use the Nullable wrapper methods:
- `obj.Pool.IsSet()` — check if set
- `obj.Pool.Get()` — get the inner value (returns pointer)
- `obj.Pool.Set(&val)` — set the value
- `obj.Pool.Unset()` — clear the value
### LastReserved (Nullable)

Use the Nullable wrapper methods:
- `obj.LastReserved.IsSet()` — check if set
- `obj.LastReserved.Get()` — get the inner value (returns pointer)
- `obj.LastReserved.Set(&val)` — set the value
- `obj.LastReserved.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


