# GetNetworkProxies200Response

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**NetworkProxies** | Pointer to **interface{}** |  | [optional] 
**NetworkProxyCount** | Pointer to **int32** |  | [optional] 
**Meta** | Pointer to [**ListApprovals200ResponseAllOfMeta**](ListApprovals200ResponseAllOfMeta.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &GetNetworkProxies200Response{
    // Set fields directly
}
```

### NetworkProxies (Nullable)

Use the Nullable wrapper methods:
- `obj.NetworkProxies.IsSet()` — check if set
- `obj.NetworkProxies.Get()` — get the inner value (returns pointer)
- `obj.NetworkProxies.Set(&val)` — set the value
- `obj.NetworkProxies.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


