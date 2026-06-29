# ListImageBuilds200ResponseAllOfImageBuildsInnerConfigNetworkInterfacesInnerNetwork

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**IdName** | Pointer to **string** |  | [optional] 
**Pool** | Pointer to [**NullableListImageBuilds200ResponseAllOfImageBuildsInnerConfigNetworkInterfacesInnerNetworkPool**](ListImageBuilds200ResponseAllOfImageBuildsInnerConfigNetworkInterfacesInnerNetworkPool.md) |  | [optional] 
**Id** | Pointer to **string** |  | [optional] 
**HasPool** | Pointer to **bool** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &ListImageBuilds200ResponseAllOfImageBuildsInnerConfigNetworkInterfacesInnerNetwork{
    // Set fields directly
}
```

### Pool (Nullable)

Use the Nullable wrapper methods:
- `obj.Pool.IsSet()` — check if set
- `obj.Pool.Get()` — get the inner value (returns pointer)
- `obj.Pool.Set(&val)` — set the value
- `obj.Pool.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


