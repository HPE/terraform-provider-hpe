# ListHealth200ResponseAllOfHealthThreadsBusyThreadsInner

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**ThreadId** | Pointer to **int64** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**CpuTime** | Pointer to **int64** |  | [optional] 
**BlockedTime** | Pointer to **int64** |  | [optional] 
**LockName** | Pointer to **NullableString** |  | [optional] 
**LockOwnerId** | Pointer to **int64** |  | [optional] 
**LockOwnerName** | Pointer to **NullableString** |  | [optional] 
**State** | Pointer to **string** |  | [optional] 
**WaitedCount** | Pointer to **int64** |  | [optional] 
**WaitedTime** | Pointer to **int64** |  | [optional] 
**IsInNative** | Pointer to **bool** |  | [optional] 
**IsSuspended** | Pointer to **bool** |  | [optional] 
**LockedMonitors** | Pointer to **[]map[string]interface{}** |  | [optional] 
**LockedSynchronizers** | Pointer to **[]map[string]interface{}** |  | [optional] 
**LockInfo** | Pointer to **NullableString** |  | [optional] 
**CurrentLines** | Pointer to **string** |  | [optional] 
**CpuPercent** | Pointer to **float32** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &ListHealth200ResponseAllOfHealthThreadsBusyThreadsInner{
    // Set fields directly
}
```

### LockName (Nullable)

Use the Nullable wrapper methods:
- `obj.LockName.IsSet()` — check if set
- `obj.LockName.Get()` — get the inner value (returns pointer)
- `obj.LockName.Set(&val)` — set the value
- `obj.LockName.Unset()` — clear the value
### LockOwnerName (Nullable)

Use the Nullable wrapper methods:
- `obj.LockOwnerName.IsSet()` — check if set
- `obj.LockOwnerName.Get()` — get the inner value (returns pointer)
- `obj.LockOwnerName.Set(&val)` — set the value
- `obj.LockOwnerName.Unset()` — clear the value
### LockedMonitors (Nullable)

Use the Nullable wrapper methods:
- `obj.LockedMonitors.IsSet()` — check if set
- `obj.LockedMonitors.Get()` — get the inner value (returns pointer)
- `obj.LockedMonitors.Set(&val)` — set the value
- `obj.LockedMonitors.Unset()` — clear the value
### LockedSynchronizers (Nullable)

Use the Nullable wrapper methods:
- `obj.LockedSynchronizers.IsSet()` — check if set
- `obj.LockedSynchronizers.Get()` — get the inner value (returns pointer)
- `obj.LockedSynchronizers.Set(&val)` — set the value
- `obj.LockedSynchronizers.Unset()` — clear the value
### LockInfo (Nullable)

Use the Nullable wrapper methods:
- `obj.LockInfo.IsSet()` — check if set
- `obj.LockInfo.Get()` — get the inner value (returns pointer)
- `obj.LockInfo.Set(&val)` — set the value
- `obj.LockInfo.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


