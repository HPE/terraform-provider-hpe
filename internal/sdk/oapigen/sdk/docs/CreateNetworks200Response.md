# CreateNetworks200Response

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Network** | Pointer to [**CreateNetworks200ResponseAllOfNetwork**](CreateNetworks200ResponseAllOfNetwork.md) |  | [optional] 
**Errors** | Pointer to **map[string]interface{}** |  | [optional] 
**Success** | Pointer to **bool** |  | [optional] 
**Msg** | Pointer to **NullableString** |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &CreateNetworks200Response{
    // Set fields directly
}
```

### Msg (Nullable)

Use the Nullable wrapper methods:
- `obj.Msg.IsSet()` — check if set
- `obj.Msg.Get()` — get the inner value (returns pointer)
- `obj.Msg.Set(&val)` — set the value
- `obj.Msg.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


