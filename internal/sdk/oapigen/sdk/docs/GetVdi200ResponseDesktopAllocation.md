# GetVdi200ResponseDesktopAllocation

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**PoolId** | Pointer to **int64** |  | [optional] 
**Pool** | Pointer to [**GetVdi200ResponseDesktopAllocationPool**](GetVdi200ResponseDesktopAllocationPool.md) |  | [optional] 
**Instance** | Pointer to [**GetVdi200ResponseDesktopAllocationInstance**](GetVdi200ResponseDesktopAllocationInstance.md) |  | [optional] 
**User** | Pointer to [**GetVdi200ResponseDesktopAllocationUser**](GetVdi200ResponseDesktopAllocationUser.md) |  | [optional] 
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
obj := &GetVdi200ResponseDesktopAllocation{
    // Set fields directly
}
```

### LastReserved (Nullable)

Use the Nullable wrapper methods:
- `obj.LastReserved.IsSet()` — check if set
- `obj.LastReserved.Get()` — get the inner value (returns pointer)
- `obj.LastReserved.Set(&val)` — set the value
- `obj.LastReserved.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


