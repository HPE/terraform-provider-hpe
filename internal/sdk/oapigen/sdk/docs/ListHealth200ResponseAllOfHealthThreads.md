# ListHealth200ResponseAllOfHealthThreads

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Success** | Pointer to **bool** |  | [optional] 
**ThreadList** | Pointer to **[]map[string]interface{}** |  | [optional] 
**BusyThreads** | Pointer to [**[]ListHealth200ResponseAllOfHealthThreadsBusyThreadsInner**](ListHealth200ResponseAllOfHealthThreadsBusyThreadsInner.md) |  | [optional] 
**BlockedThreads** | Pointer to **[]map[string]interface{}** |  | [optional] 
**RunningThreads** | Pointer to **[]map[string]interface{}** |  | [optional] 
**TotalCpuTime** | Pointer to **int64** |  | [optional] 
**TotalThreads** | Pointer to **int64** |  | [optional] 
**RunningWebThreads** | Pointer to **int64** |  | [optional] 
**Status** | Pointer to **string** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &ListHealth200ResponseAllOfHealthThreads{
    // Set fields directly
}
```

### ThreadList (Nullable)

Use the Nullable wrapper methods:
- `obj.ThreadList.IsSet()` — check if set
- `obj.ThreadList.Get()` — get the inner value (returns pointer)
- `obj.ThreadList.Set(&val)` — set the value
- `obj.ThreadList.Unset()` — clear the value
### BlockedThreads (Nullable)

Use the Nullable wrapper methods:
- `obj.BlockedThreads.IsSet()` — check if set
- `obj.BlockedThreads.Get()` — get the inner value (returns pointer)
- `obj.BlockedThreads.Set(&val)` — set the value
- `obj.BlockedThreads.Unset()` — clear the value
### RunningThreads (Nullable)

Use the Nullable wrapper methods:
- `obj.RunningThreads.IsSet()` — check if set
- `obj.RunningThreads.Get()` — get the inner value (returns pointer)
- `obj.RunningThreads.Set(&val)` — set the value
- `obj.RunningThreads.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


