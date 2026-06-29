# UpdateTenantRequestAccount

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** | Name | [optional] 
**Description** | Pointer to **NullableString** | Description | [optional] 
**Role** | Pointer to [**UpdateTenantRequestAccountRole**](UpdateTenantRequestAccountRole.md) |  | [optional] 
**Subdomain** | Pointer to **NullableString** | The subdomain. This will be part of the login URL and username for sub tenant users. | [optional] 
**Currency** | Pointer to **string** | Currency Code (ISO 4217) | [optional] [default to "USD"]

## Usage

Instantiate with a Go composite literal:

```go
obj := &UpdateTenantRequestAccount{
    // Set fields directly
}
```

### Description (Nullable)

Use the Nullable wrapper methods:
- `obj.Description.IsSet()` — check if set
- `obj.Description.Get()` — get the inner value (returns pointer)
- `obj.Description.Set(&val)` — set the value
- `obj.Description.Unset()` — clear the value
### Subdomain (Nullable)

Use the Nullable wrapper methods:
- `obj.Subdomain.IsSet()` — check if set
- `obj.Subdomain.Get()` — get the inner value (returns pointer)
- `obj.Subdomain.Set(&val)` — set the value
- `obj.Subdomain.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


