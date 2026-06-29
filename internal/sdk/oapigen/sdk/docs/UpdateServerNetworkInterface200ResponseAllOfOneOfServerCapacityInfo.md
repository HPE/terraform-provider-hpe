# UpdateServerNetworkInterface200ResponseAllOfOneOfServerCapacityInfo

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**MaxMemory** | Pointer to **int64** |  | [optional] 
**Server** | Pointer to [**UpdateServerNetworkInterface200ResponseAllOfOneOfServerCapacityInfoServer**](UpdateServerNetworkInterface200ResponseAllOfOneOfServerCapacityInfoServer.md) |  | [optional] 
**UsedStorage** | Pointer to **int64** |  | [optional] 
**Version** | Pointer to **NullableString** |  | [optional] 
**MaxCpu** | Pointer to **NullableString** |  | [optional] 
**UsedCores** | Pointer to **int64** |  | [optional] 
**UsedMemory** | Pointer to **int64** |  | [optional] 
**MaxCores** | Pointer to **int64** |  | [optional] 
**MaxStorage** | Pointer to **int64** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &UpdateServerNetworkInterface200ResponseAllOfOneOfServerCapacityInfo{
    // Set fields directly
}
```

### Version (Nullable)

Use the Nullable wrapper methods:
- `obj.Version.IsSet()` — check if set
- `obj.Version.Get()` — get the inner value (returns pointer)
- `obj.Version.Set(&val)` — set the value
- `obj.Version.Unset()` — clear the value
### MaxCpu (Nullable)

Use the Nullable wrapper methods:
- `obj.MaxCpu.IsSet()` — check if set
- `obj.MaxCpu.Get()` — get the inner value (returns pointer)
- `obj.MaxCpu.Set(&val)` — set the value
- `obj.MaxCpu.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


