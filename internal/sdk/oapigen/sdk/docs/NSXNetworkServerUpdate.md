# NSXNetworkServerUpdate

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** | Name | [optional] 
**ServiceUrl** | Pointer to **NullableString** | URL | [optional] 
**ServiceUsername** | Pointer to **NullableString** | Username | [optional] 
**ServicePassword** | Pointer to **NullableString** | Password | [optional] 
**Visibility** | Pointer to **string** | Visibility - Set to public to allow all tenants | [optional] [default to "private"]
**Tenants** | Pointer to [**[]NSXNetworkServerUpdateTenantsInner**](NSXNetworkServerUpdateTenantsInner.md) | Array of tenant account ids that are allowed access | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &NSXNetworkServerUpdate{
    // Set fields directly
}
```

### ServiceUrl (Nullable)

Use the Nullable wrapper methods:
- `obj.ServiceUrl.IsSet()` — check if set
- `obj.ServiceUrl.Get()` — get the inner value (returns pointer)
- `obj.ServiceUrl.Set(&val)` — set the value
- `obj.ServiceUrl.Unset()` — clear the value
### ServiceUsername (Nullable)

Use the Nullable wrapper methods:
- `obj.ServiceUsername.IsSet()` — check if set
- `obj.ServiceUsername.Get()` — get the inner value (returns pointer)
- `obj.ServiceUsername.Set(&val)` — set the value
- `obj.ServiceUsername.Unset()` — clear the value
### ServicePassword (Nullable)

Use the Nullable wrapper methods:
- `obj.ServicePassword.IsSet()` — check if set
- `obj.ServicePassword.Get()` — get the inner value (returns pointer)
- `obj.ServicePassword.Set(&val)` — set the value
- `obj.ServicePassword.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


