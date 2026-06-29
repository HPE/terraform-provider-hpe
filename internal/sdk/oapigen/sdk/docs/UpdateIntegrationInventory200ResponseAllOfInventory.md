# UpdateIntegrationInventory200ResponseAllOfInventory

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**Description** | Pointer to **NullableString** |  | [optional] 
**ExternalId** | Pointer to **NullableString** |  | [optional] 
**Owner** | Pointer to [**UpdateIntegrationInventory200ResponseAllOfInventoryOwner**](UpdateIntegrationInventory200ResponseAllOfInventoryOwner.md) |  | [optional] 
**Tenants** | Pointer to [**[]UpdateIntegrationInventory200ResponseAllOfInventoryTenantsInner**](UpdateIntegrationInventory200ResponseAllOfInventoryTenantsInner.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &UpdateIntegrationInventory200ResponseAllOfInventory{
    // Set fields directly
}
```

### Description (Nullable)

Use the Nullable wrapper methods:
- `obj.Description.IsSet()` — check if set
- `obj.Description.Get()` — get the inner value (returns pointer)
- `obj.Description.Set(&val)` — set the value
- `obj.Description.Unset()` — clear the value
### ExternalId (Nullable)

Use the Nullable wrapper methods:
- `obj.ExternalId.IsSet()` — check if set
- `obj.ExternalId.Get()` — get the inner value (returns pointer)
- `obj.ExternalId.Set(&val)` — set the value
- `obj.ExternalId.Unset()` — clear the value

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


