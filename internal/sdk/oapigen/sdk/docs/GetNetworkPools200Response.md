# GetNetworkPools200Response

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**NetworkPools** | Pointer to **interface{}** |  | [optional] 
**NetworkPoolCount** | Pointer to **int64** |  | [optional] 
**Meta** | Pointer to [**ListApprovals200ResponseAllOfMeta**](ListApprovals200ResponseAllOfMeta.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &GetNetworkPools200Response{
    // Set fields directly
}
```

### NetworkPools (Nullable)

Use the Nullable wrapper methods:
- `obj.NetworkPools.IsSet()` — check if set
- `obj.NetworkPools.Get()` — get the inner value (returns pointer)
- `obj.NetworkPools.Set(&val)` — set the value
- `obj.NetworkPools.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


