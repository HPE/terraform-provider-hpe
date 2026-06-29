# GetIntegrationInventory200ResponseInventory

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **int64** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 
**Description** | Pointer to **NullableString** |  | [optional] 
**ExternalId** | Pointer to **NullableString** |  | [optional] 
**Owner** | Pointer to [**GetIntegrationInventory200ResponseInventoryOwner**](GetIntegrationInventory200ResponseInventoryOwner.md) |  | [optional] 
**Tenants** | Pointer to [**[]GetIntegrationInventory200ResponseInventoryTenantsInner**](GetIntegrationInventory200ResponseInventoryTenantsInner.md) |  | [optional] 

## Usage

Instantiate with a Go composite literal:

```go
obj := &GetIntegrationInventory200ResponseInventory{
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


