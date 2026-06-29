# UpdateInstanceNetworkInterface200Response

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**NetworkInterface** | Pointer to [**UpdateInstanceNetworkInterface200ResponseAllOfOneOfNetworkInterface**](UpdateInstanceNetworkInterface200ResponseAllOfOneOfNetworkInterface.md) |  | [optional] 
**InterfaceType** | Pointer to **string** |  | [optional] 
**NetId** | Pointer to **int64** |  | [optional] 
**Server** | Pointer to [**UpdateInstanceNetworkInterface200ResponseAllOfOneOfServer**](UpdateInstanceNetworkInterface200ResponseAllOfOneOfServer.md) |  | [optional] 
**Success** | Pointer to **bool** |  | [optional] 
**Errors** | Pointer to **map[string]interface{}** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &UpdateInstanceNetworkInterface200Response{
    // Set fields directly
}
```

### Errors (Nullable)

Use the Nullable wrapper methods:
- `obj.Errors.IsSet()` — check if set
- `obj.Errors.Get()` — get the inner value (returns pointer)
- `obj.Errors.Set(&val)` — set the value
- `obj.Errors.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


