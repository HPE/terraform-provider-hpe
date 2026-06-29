# CreateNetworkTransportZoneRequestNetworkScope

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | **string** | Network transport zone name | 
**Description** | Pointer to **NullableString** | Network transport zone description | [optional] 
**Visibility** | Pointer to **string** | private or public | [optional] 
**Tenants** | Pointer to [**[]CreateNetworkTransportZoneRequestNetworkScopeTenantsInner**](CreateNetworkTransportZoneRequestNetworkScopeTenantsInner.md) | Array of tenant account ids that are allowed access | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &CreateNetworkTransportZoneRequestNetworkScope{
    // Set fields directly
}
```

### Description (Nullable)

Use the Nullable wrapper methods:
- `obj.Description.IsSet()` — check if set
- `obj.Description.Get()` — get the inner value (returns pointer)
- `obj.Description.Set(&val)` — set the value
- `obj.Description.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


